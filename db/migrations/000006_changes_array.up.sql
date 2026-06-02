-- Migrate single-change objects to single-element arrays
UPDATE image_table
  SET changes = jsonb_build_array(changes)
  WHERE changes != '{}' AND changes != '[]';

-- Convert legacy empty-object default to empty array
UPDATE image_table SET changes = '[]' WHERE changes = '{}';

ALTER TABLE image_table ALTER COLUMN changes SET DEFAULT '[]';
