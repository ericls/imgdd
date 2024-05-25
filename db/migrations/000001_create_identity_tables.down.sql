BEGIN;

DROP TABLE IF EXISTS organization_user_role_table;
DROP TABLE IF EXISTS organization_user_table;
ALTER TABLE role_table
DROP CONSTRAINT IF EXISTS unique_role_key_organization_id;
DROP TABLE IF EXISTS role_table;
ALTER TABLE user_table
DROP CONSTRAINT IF EXISTS unique_user_email;
DROP TABLE IF EXISTS user_table;
DROP TABLE IF EXISTS organization_table;
COMMIT;