BEGIN;

CREATE TABLE album_table (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- unguessable public identifier, used in URLs and the API; the id is never exposed
    identifier CHARACTER VARYING(64) NOT NULL,
    title CHARACTER VARYING(255) NOT NULL DEFAULT '',
    created_by_id UUID REFERENCES organization_user_table(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,

    CONSTRAINT album_table_unique_identifier UNIQUE(identifier)
);

CREATE INDEX album_table_created_by_id_idx ON album_table(created_by_id);

-- Albums and images are many-to-many. The caption describes the image's place
-- in this album, so the same image can carry a different caption elsewhere.
-- Removing a row removes the image from the album; it never deletes the image.
CREATE TABLE album_image_table (
    album_id UUID NOT NULL REFERENCES album_table(id) ON DELETE CASCADE,
    image_id UUID NOT NULL REFERENCES image_table(id) ON DELETE CASCADE,
    caption CHARACTER VARYING(2000) NULL DEFAULT NULL,
    added_by_id UUID REFERENCES organization_user_table(id) ON DELETE SET NULL,
    added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    PRIMARY KEY (album_id, image_id)
);

-- Album pages list images in the order they were added.
CREATE INDEX album_image_table_album_id_added_at_idx ON album_image_table(album_id, added_at, image_id);
-- Lets image deletes find their memberships.
CREATE INDEX album_image_table_image_id_idx ON album_image_table(image_id);

COMMIT;
