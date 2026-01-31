CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create schema for this service's domain
CREATE SCHEMA IF NOT EXISTS identity;

-- A User represents a user that can log in to BrewPipes to perform operations with a set of tokens.
CREATE TABLE IF NOT EXISTS identity.user (
    id           serial PRIMARY KEY,
    uuid         uuid NOT NULL DEFAULT gen_random_uuid(),

    username     text UNIQUE NOT NULL,
    password     text NOT NULL,

    created_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at   timestamptz,
    CONSTRAINT user_uuid_unique UNIQUE (uuid)
);

-- Refresh tokens are tracked to allow rotation and revocation.
CREATE TABLE IF NOT EXISTS identity.refresh_token (
    id           serial PRIMARY KEY,
    token_id     uuid NOT NULL,
    user_uuid    uuid NOT NULL,
    expires_at   timestamptz NOT NULL,
    revoked_at   timestamptz,
    created_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    CONSTRAINT refresh_token_id_unique UNIQUE (token_id),
    CONSTRAINT refresh_token_user_uuid_fk FOREIGN KEY (user_uuid)
        REFERENCES identity.user (uuid)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS refresh_token_user_uuid_idx
    ON identity.refresh_token (user_uuid);

-- Seed data for early development.
INSERT INTO identity.user (uuid, username, password)
VALUES
    ('10000000-0000-0000-0000-000000000001', 'brewmaster', crypt('brewpipes', gen_salt('bf'))),
    ('10000000-0000-0000-0000-000000000002', 'brewer', crypt('password', gen_salt('bf')));
