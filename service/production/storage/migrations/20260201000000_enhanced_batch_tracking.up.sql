-- Enhanced Batch Tracking: style, recipe, brew_session tables; occupancy status; volume_id on addition/measurement; BBL unit

-- Style table for beer styles (case-insensitive unique names)
CREATE TABLE IF NOT EXISTS style (
    id          serial PRIMARY KEY,
    uuid        uuid NOT NULL DEFAULT gen_random_uuid(),
    name        varchar(255) NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at  timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS style_uuid_idx ON style(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS style_name_lower_idx ON style(lower(name)) WHERE deleted_at IS NULL;

-- Recipe table for beer formulations
CREATE TABLE IF NOT EXISTS recipe (
    id          serial PRIMARY KEY,
    uuid        uuid NOT NULL DEFAULT gen_random_uuid(),
    name        varchar(255) NOT NULL,
    style_id    int REFERENCES style(id),
    style_name  varchar(255),
    notes       text,
    created_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at  timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS recipe_uuid_idx ON recipe(uuid);
CREATE INDEX IF NOT EXISTS recipe_style_id_idx ON recipe(style_id) WHERE deleted_at IS NULL;

-- Brew session table for hot-side wort production
CREATE TABLE IF NOT EXISTS brew_session (
    id              serial PRIMARY KEY,
    uuid            uuid NOT NULL DEFAULT gen_random_uuid(),
    batch_id        int REFERENCES batch(id),
    wort_volume_id  int REFERENCES volume(id),
    mash_vessel_id  int REFERENCES vessel(id),
    boil_vessel_id  int REFERENCES vessel(id),
    brewed_at       timestamptz NOT NULL,
    notes           text,
    created_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at      timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS brew_session_uuid_idx ON brew_session(uuid);
CREATE INDEX IF NOT EXISTS brew_session_batch_id_idx ON brew_session(batch_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS brew_session_wort_volume_id_idx ON brew_session(wort_volume_id) WHERE deleted_at IS NULL;

-- Add recipe_id to batch
ALTER TABLE batch ADD COLUMN recipe_id int REFERENCES recipe(id);
CREATE INDEX IF NOT EXISTS batch_recipe_id_idx ON batch(recipe_id) WHERE deleted_at IS NULL;

-- Add status to occupancy
ALTER TABLE occupancy ADD COLUMN status varchar(32);
ALTER TABLE occupancy ADD CONSTRAINT occupancy_status_check 
    CHECK (status IS NULL OR status IN (
        'fermenting', 'conditioning', 'cold_crashing', 
        'dry_hopping', 'carbonating', 'holding', 'packaging'
    ));

-- Add volume_id to addition
ALTER TABLE addition ADD COLUMN volume_id int REFERENCES volume(id);
CREATE INDEX IF NOT EXISTS addition_volume_id_idx ON addition(volume_id) WHERE deleted_at IS NULL;

-- Update addition target constraint to include volume_id option
ALTER TABLE addition DROP CONSTRAINT addition_target_check;
ALTER TABLE addition ADD CONSTRAINT addition_target_check CHECK (
    (batch_id IS NOT NULL AND occupancy_id IS NULL AND volume_id IS NULL) OR
    (batch_id IS NULL AND occupancy_id IS NOT NULL AND volume_id IS NULL) OR
    (batch_id IS NULL AND occupancy_id IS NULL AND volume_id IS NOT NULL)
);

-- Add volume_id to measurement
ALTER TABLE measurement ADD COLUMN volume_id int REFERENCES volume(id);
CREATE INDEX IF NOT EXISTS measurement_volume_id_idx ON measurement(volume_id) WHERE deleted_at IS NULL;

-- Update measurement target constraint to include volume_id option
ALTER TABLE measurement DROP CONSTRAINT measurement_target_check;
ALTER TABLE measurement ADD CONSTRAINT measurement_target_check CHECK (
    (batch_id IS NOT NULL AND occupancy_id IS NULL AND volume_id IS NULL) OR
    (batch_id IS NULL AND occupancy_id IS NOT NULL AND volume_id IS NULL) OR
    (batch_id IS NULL AND occupancy_id IS NULL AND volume_id IS NOT NULL)
);
