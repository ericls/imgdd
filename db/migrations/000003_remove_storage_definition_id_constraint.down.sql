-- down.sql
DROP INDEX IF EXISTS stored_image_table_storage_definition_id_idx;

ALTER TABLE stored_image_table
    ADD CONSTRAINT stored_image_table_storage_definition_id_fkey
        FOREIGN KEY (storage_definition_id)
        REFERENCES storage_definition (id)
        ON DELETE SET NULL;
