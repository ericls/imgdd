DROP INDEX IF EXISTS image_table_expires_at_idx;
ALTER TABLE image_table
  DROP COLUMN expires_at;
