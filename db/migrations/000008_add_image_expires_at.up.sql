-- Expired images are auto marked as deleted periodically (deleted_at is set);
-- expires_at is only read by that sweep.
ALTER TABLE image_table
  ADD COLUMN expires_at TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL;

-- Supports the periodic sweep for undeleted images that have expired.
CREATE INDEX image_table_expires_at_idx ON image_table(expires_at)
  WHERE expires_at IS NOT NULL AND deleted_at IS NULL;
