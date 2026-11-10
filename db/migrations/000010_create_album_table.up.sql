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

ALTER TABLE image_table
  ADD COLUMN album_id UUID NULL DEFAULT NULL REFERENCES album_table(id) ON DELETE SET NULL,
  ADD COLUMN caption CHARACTER VARYING(2000) NULL DEFAULT NULL;

CREATE INDEX image_table_album_id_created_at_idx ON image_table(album_id, created_at)
  WHERE album_id IS NOT NULL;

COMMIT;
