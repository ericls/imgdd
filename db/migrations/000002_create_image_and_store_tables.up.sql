BEGIN;
-- Create the image tables
CREATE TABLE image_table (
    id UUID PRIMARY KEY,
    organization_id UUID REFERENCES organization_table(id) ON DELETE SET NULL,
    created_by UUID REFERENCES user_table(id) ON DELETE SET NULL,
    --
    name CHARACTER VARYING(255) NOT NULL,
    -- this is the unique identifier for the image, used for URL
    identifier CHARACTER VARYING(255) NOT NULL,
    --
    root UUID REFERENCES image_table(id) ON DELETE CASCADE,
    parent UUID REFERENCES image_table(id) ON DELETE CASCADE,
    -- changes from parent
    changes JSONB NOT NULL DEFAULT '{}',
    --
    uploader_ip INET,
    --
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE INDEX image_table_created_at_idx ON image_table(created_at);
CREATE INDEX image_table_updated_at_idx ON image_table(updated_at);
CREATE INDEX image_table_deleted_at_idx ON image_table(deleted_at);
CREATE INDEX image_table_identifier_idx ON image_table USING HASH (identifier);

ALTER TABLE
    image_table
ADD
    CONSTRAINT image_table_unique_identifier UNIQUE(identifier);

-- Storage Config table
CREATE TABLE storage_definition_table (
    id UUID PRIMARY KEY,

    type CHARACTER VARYING(255) NOT NULL,
    identifier CHARACTER VARYING(255) NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    priority INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX storage_definition_table_identifier_idx ON storage_definition_table  USING HASH (identifier);
ALTER TABLE
    storage_definition_table
ADD
    CONSTRAINT storage_definition_table_unique_identifier UNIQUE(identifier);


-- Stored Image table
CREATE TABLE stored_image_table (
    id UUID PRIMARY KEY,
    image_id UUID REFERENCES image_table(id) ON DELETE SET NULL,
    storage_definition_id UUID REFERENCES storage_definition_table(id) ON DELETE SET NULL,

    file_identifier CHARACTER VARYING(255) NOT NULL,
    copied_from UUID REFERENCES stored_image_table(id) ON DELETE SET NULL,
    --
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    is_file_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

ALTER TABLE
    stored_image_table
ADD
    CONSTRAINT stored_image_table_unique_storage_definition_file_identifier UNIQUE(storage_definition_id, file_identifier);
ALTER TABLE
    stored_image_table
ADD
    CONSTRAINT stored_image_table_unique_image_id_storage_definition_id UNIQUE(image_id, storage_definition_id);
---
COMMIT;