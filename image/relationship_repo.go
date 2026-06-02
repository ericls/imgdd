//lint:file-ignore ST1001 Allow using dot imports following Jet's convention
package image

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/db/.gen/imgdd/public/model"
	. "github.com/ericls/imgdd/db/.gen/imgdd/public/table"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

const (
	RelationshipTypeBase    = "base"
	RelationshipTypeOverlay = "overlay"
)

type ImageRelationship struct {
	Id               string
	ImageId          string
	ParentImageId    string
	RelationshipType string
}

type ImageRelationshipRepo interface {
	CreateRelationship(imageId, parentImageId, relationshipType string) (*ImageRelationship, error)
	GetParentsByImageId(imageId string) ([]ImageRelationship, error)
	GetChildrenByImageId(imageId string) ([]ImageRelationship, error)
	HasRelationships(imageId string) (bool, error)
	// GetDescendantIds returns all transitive children of the given image (not including itself).
	GetDescendantIds(imageId string) ([]string, error)
	// GetAncestorIds returns all transitive parents of the given image (not including itself).
	GetAncestorIds(imageId string) ([]string, error)
	// AreRelated returns true if there is any path between the two images in the DAG.
	AreRelated(imageId1, imageId2 string) (bool, error)
	// IsAncestor returns true if ancestorId is a transitive parent of imageId.
	IsAncestor(imageId, ancestorId string) (bool, error)
}

type DBImageRelationshipRepo struct {
	db.RepoConn
}

func (repo *DBImageRelationshipRepo) WithTransaction(tx *sql.Tx) db.DBRepo {
	return &DBImageRelationshipRepo{
		RepoConn: repo.RepoConn.WithTransaction(tx),
	}
}

func NewDBImageRelationshipRepo(conn *sql.DB) *DBImageRelationshipRepo {
	return &DBImageRelationshipRepo{
		RepoConn: db.NewRepoConn(conn),
	}
}

func (repo *DBImageRelationshipRepo) CreateRelationship(imageId, parentImageId, relationshipType string) (*ImageRelationship, error) {
	if imageId == parentImageId {
		return nil, fmt.Errorf("an image cannot be its own parent")
	}

	// Cycle detection: if imageId is already an ancestor of parentImageId,
	// adding parentImageId as a parent of imageId would create a cycle.
	isAnc, err := repo.IsAncestor(parentImageId, imageId)
	if err != nil {
		return nil, fmt.Errorf("cycle detection failed: %w", err)
	}
	if isAnc {
		return nil, fmt.Errorf("cannot create relationship: would form a cycle")
	}

	stmt := ImageParentTable.INSERT(
		ImageParentTable.ImageID,
		ImageParentTable.ParentImageID,
		ImageParentTable.RelationshipType,
	).VALUES(
		UUID(uuid.MustParse(imageId)),
		UUID(uuid.MustParse(parentImageId)),
		relationshipType,
	).RETURNING(
		ImageParentTable.AllColumns,
	)

	dest := model.ImageParentTable{}
	err = stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to create image relationship: %w", err)
	}
	return &ImageRelationship{
		Id:               dest.ID.String(),
		ImageId:          dest.ImageID.String(),
		ParentImageId:    dest.ParentImageID.String(),
		RelationshipType: dest.RelationshipType,
	}, nil
}

func (repo *DBImageRelationshipRepo) GetParentsByImageId(imageId string) ([]ImageRelationship, error) {
	stmt := ImageParentTable.SELECT(
		ImageParentTable.AllColumns,
	).FROM(
		ImageParentTable,
	).WHERE(
		ImageParentTable.ImageID.EQ(UUID(uuid.MustParse(imageId))),
	).ORDER_BY(
		ImageParentTable.CreatedAt.ASC(),
	)

	var dest []model.ImageParentTable
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	return mapRelationships(dest), nil
}

func (repo *DBImageRelationshipRepo) GetChildrenByImageId(imageId string) ([]ImageRelationship, error) {
	stmt := ImageParentTable.SELECT(
		ImageParentTable.AllColumns,
	).FROM(
		ImageParentTable,
	).WHERE(
		ImageParentTable.ParentImageID.EQ(UUID(uuid.MustParse(imageId))),
	).ORDER_BY(
		ImageParentTable.CreatedAt.ASC(),
	)

	var dest []model.ImageParentTable
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	return mapRelationships(dest), nil
}

func (repo *DBImageRelationshipRepo) HasRelationships(imageId string) (bool, error) {
	id := UUID(uuid.MustParse(imageId))
	stmt := ImageParentTable.SELECT(
		ImageParentTable.ID,
	).FROM(
		ImageParentTable,
	).WHERE(
		ImageParentTable.ImageID.EQ(id).OR(ImageParentTable.ParentImageID.EQ(id)),
	).LIMIT(1)

	var dest []model.ImageParentTable
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return false, err
	}
	return len(dest) > 0, nil
}

// GetDescendantIds walks the DAG downward from imageId using a recursive CTE.
func (repo *DBImageRelationshipRepo) GetDescendantIds(imageId string) ([]string, error) {
	query := `
		WITH RECURSIVE descendants AS (
			SELECT image_id FROM image_parent_table WHERE parent_image_id = $1
			UNION
			SELECT ip.image_id FROM image_parent_table ip
			INNER JOIN descendants d ON ip.parent_image_id = d.image_id
		)
		SELECT image_id FROM descendants
	`
	return repo.queryIds(query, imageId)
}

// GetAncestorIds walks the DAG upward from imageId using a recursive CTE.
func (repo *DBImageRelationshipRepo) GetAncestorIds(imageId string) ([]string, error) {
	query := `
		WITH RECURSIVE ancestors AS (
			SELECT parent_image_id FROM image_parent_table WHERE image_id = $1
			UNION
			SELECT ip.parent_image_id FROM image_parent_table ip
			INNER JOIN ancestors a ON ip.image_id = a.parent_image_id
		)
		SELECT parent_image_id FROM ancestors
	`
	return repo.queryIds(query, imageId)
}

// AreRelated returns true if there is any DAG path between the two images (in either direction).
func (repo *DBImageRelationshipRepo) AreRelated(imageId1, imageId2 string) (bool, error) {
	// Check if imageId2 is an ancestor or descendant of imageId1.
	// We walk both directions from imageId1 and check for imageId2.
	query := `
		WITH RECURSIVE
		descendants AS (
			SELECT image_id AS id FROM image_parent_table WHERE parent_image_id = $1
			UNION
			SELECT ip.image_id FROM image_parent_table ip
			INNER JOIN descendants d ON ip.parent_image_id = d.id
		),
		ancestors AS (
			SELECT parent_image_id AS id FROM image_parent_table WHERE image_id = $1
			UNION
			SELECT ip.parent_image_id FROM image_parent_table ip
			INNER JOIN ancestors a ON ip.image_id = a.id
		)
		SELECT 1 WHERE EXISTS (
			SELECT 1 FROM descendants WHERE id = $2
			UNION ALL
			SELECT 1 FROM ancestors WHERE id = $2
		)
	`
	return repo.queryExists(query, imageId1, imageId2)
}

// IsAncestor returns true if ancestorId is a transitive parent of imageId.
func (repo *DBImageRelationshipRepo) IsAncestor(imageId, ancestorId string) (bool, error) {
	query := `
		WITH RECURSIVE ancestors AS (
			SELECT parent_image_id FROM image_parent_table WHERE image_id = $1
			UNION
			SELECT ip.parent_image_id FROM image_parent_table ip
			INNER JOIN ancestors a ON ip.image_id = a.parent_image_id
		)
		SELECT 1 WHERE EXISTS (SELECT 1 FROM ancestors WHERE parent_image_id = $2)
	`
	return repo.queryExists(query, imageId, ancestorId)
}

func (repo *DBImageRelationshipRepo) queryIds(query string, imageId string) ([]string, error) {
	rows, err := repo.DB.QueryContext(context.Background(), query, imageId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (repo *DBImageRelationshipRepo) queryExists(query string, args ...interface{}) (bool, error) {
	rows, err := repo.DB.QueryContext(context.Background(), query, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), rows.Err()
}

func mapRelationships(rows []model.ImageParentTable) []ImageRelationship {
	result := make([]ImageRelationship, len(rows))
	for i, r := range rows {
		result[i] = ImageRelationship{
			Id:               r.ID.String(),
			ImageId:          r.ImageID.String(),
			ParentImageId:    r.ParentImageID.String(),
			RelationshipType: r.RelationshipType,
		}
	}
	return result
}
