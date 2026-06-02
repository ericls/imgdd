BEGIN;

CREATE TABLE image_parent_table (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image_id UUID NOT NULL REFERENCES image_table(id) ON DELETE RESTRICT,
    parent_image_id UUID NOT NULL REFERENCES image_table(id) ON DELETE RESTRICT,
    relationship_type CHARACTER VARYING(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT image_parent_unique_relationship UNIQUE(image_id, parent_image_id, relationship_type)
);

CREATE INDEX image_parent_table_image_id_idx ON image_parent_table(image_id);
CREATE INDEX image_parent_table_parent_image_id_idx ON image_parent_table(parent_image_id);

-- Backfill from existing parent_id column
INSERT INTO image_parent_table (image_id, parent_image_id, relationship_type)
SELECT id, parent_id, 'base'
FROM image_table
WHERE parent_id IS NOT NULL AND deleted_at IS NULL;

-- Backfill overlay relationships from changes JSON
INSERT INTO image_parent_table (image_id, parent_image_id, relationship_type)
SELECT
    it.id,
    (it.changes->'params'->>'overlay_image_id')::UUID,
    'overlay'
FROM image_table it
WHERE it.changes != '{}'
  AND it.changes->>'type' = 'watermark'
  AND it.changes->'params'->>'overlay_image_id' IS NOT NULL
  AND it.deleted_at IS NULL
  AND EXISTS (
    SELECT 1 FROM image_table ot
    WHERE ot.id = (it.changes->'params'->>'overlay_image_id')::UUID
    AND ot.deleted_at IS NULL
  );

COMMIT;
