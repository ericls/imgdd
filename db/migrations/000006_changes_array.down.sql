-- Restore empty arrays to empty objects
UPDATE image_table SET changes = '{}' WHERE changes = '[]';

-- Unwrap single-element arrays back to bare objects
UPDATE image_table
  SET changes = changes->0
  WHERE jsonb_array_length(changes) = 1;

ALTER TABLE image_table ALTER COLUMN changes SET DEFAULT '{}';
