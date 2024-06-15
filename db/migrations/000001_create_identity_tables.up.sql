BEGIN;

CREATE TABLE IF NOT EXISTS organization_table(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR (50) UNIQUE NOT NULL,
    display_name VARCHAR (50) UNIQUE NOT NULL,
    extra_attrs jsonb NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_table(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    password VARCHAR (128) NOT NULL,
    email VARCHAR (300) NOT NULL,
    extra_attrs jsonb NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE
    user_table
ADD
    CONSTRAINT unique_user_email UNIQUE (email);

CREATE TABLE IF NOT EXISTS role_table(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(50) NOT NULL,
    organization_id uuid REFERENCES organization_table(id),
    display_name VARCHAR(255) NOT NULL,
    extra_attrs jsonb NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE
    role_table
ADD
    CONSTRAINT unique_role_key_organization_id UNIQUE NULLS NOT DISTINCT (key, organization_id);

CREATE UNIQUE INDEX unique_site_owner_role_key ON role_table (key) WHERE key = 'site_owner';

CREATE TABLE IF NOT EXISTS organization_user_table(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id uuid REFERENCES organization_table(id) NOT NULL,
    user_id uuid REFERENCES user_table(id) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_user_role_table (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_user_id uuid REFERENCES organization_user_table(id) NOT NULL,
    role_id uuid REFERENCES role_table(id) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;