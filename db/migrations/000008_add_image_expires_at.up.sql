ALTER TABLE image_table
  ADD COLUMN expires_at TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL;

CREATE INDEX image_table_expires_at_idx ON image_table(expires_at)
  WHERE expires_at IS NOT NULL AND deleted_at IS NULL;
