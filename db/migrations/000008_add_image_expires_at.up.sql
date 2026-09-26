ALTER TABLE image_table
  ADD COLUMN expires_at TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL;

-- Supports the periodic sweep for live images that have expired.
CREATE INDEX image_table_expires_at_idx ON image_table(expires_at)
  WHERE expires_at IS NOT NULL AND deleted_at IS NULL;
