BEGIN;

DROP INDEX IF EXISTS image_table_album_id_created_at_idx;
ALTER TABLE image_table
  DROP COLUMN caption,
  DROP COLUMN album_id;
DROP TABLE IF EXISTS album_table;

COMMIT;
