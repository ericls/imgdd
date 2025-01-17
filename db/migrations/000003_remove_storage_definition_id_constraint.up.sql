-- up.sql
ALTER TABLE stored_image_table
    DROP CONSTRAINT stored_image_table_storage_definition_id_fkey;

CREATE INDEX IF NOT EXISTS stored_image_table_storage_definition_id_idx
    ON stored_image_table (storage_definition_id);
