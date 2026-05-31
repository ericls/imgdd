BEGIN;
CREATE INDEX image_table_root_id_idx ON image_table(root_id);
CREATE INDEX image_table_parent_id_idx ON image_table(parent_id);
COMMIT;
