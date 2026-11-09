CREATE TABLE IF NOT EXISTS api_key_table(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_user_id uuid REFERENCES organization_user_table(id) NOT NULL,
    name VARCHAR(100) NOT NULL,
    key_prefix VARCHAR(32) NOT NULL,
    key_hash VARCHAR(64) NOT NULL UNIQUE,
    last_used_at TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX api_key_table_organization_user_id_idx ON api_key_table(organization_user_id);
