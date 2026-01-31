-- Create schema for this service's domain
CREATE SCHEMA IF NOT EXISTS identity;

-- A User represents a user that can log in to BrewPipes perform operations with an set of tokens.
CREATE TABLE IF NOT EXISTS identity.user (
    id           serial PRIMARY KEY,
    uuid         uuid NOT NULL DEFAULT gen_random_uuid(),

    username     text UNIQUE NOT NULL,
    password     text NOT NULL,

    created_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at   timestamptz
);

-- add some seed users for testing
INSERT INTO identity.user (username, password)
VALUES
    ('user1', 'pass1'),
    ('user2', 'pass2'),
    ('user3', 'pass3'),
    ('brewmaster', 'brewpass'),
    ('cellar', 'cellarpass'),
    ('inventory', 'inventpass'),
    ('procurement', 'procpass');
