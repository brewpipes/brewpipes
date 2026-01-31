-- Production service schema.
-- Cross-domain references are stored as opaque UUIDs without foreign keys.
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS batch (
    id           serial PRIMARY KEY,
    uuid         uuid NOT NULL DEFAULT gen_random_uuid(),

    short_name   varchar(255) NOT NULL,
    brew_date    date,
    notes        text,

    created_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at   timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS batch_uuid_idx ON batch(uuid);

CREATE TABLE IF NOT EXISTS volume (
    id               serial PRIMARY KEY,
    uuid             uuid NOT NULL DEFAULT gen_random_uuid(),

    name             varchar(255),
    description      text,
    amount           bigint NOT NULL,
    amount_unit      varchar(7) NOT NULL,

    created_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at       timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS volume_uuid_idx ON volume(uuid);

CREATE TABLE IF NOT EXISTS volume_relation (
    id               serial PRIMARY KEY,
    uuid             uuid NOT NULL DEFAULT gen_random_uuid(),

    parent_volume_id int NOT NULL REFERENCES volume(id),
    child_volume_id  int NOT NULL REFERENCES volume(id),
    relation_type    varchar(16) NOT NULL,
    amount           bigint NOT NULL,
    amount_unit      varchar(7) NOT NULL,

    created_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at       timestamptz,
    CONSTRAINT volume_relation_parent_child_check CHECK (parent_volume_id <> child_volume_id),
    CONSTRAINT volume_relation_type_check CHECK (relation_type IN ('split', 'blend')),
    CONSTRAINT volume_relation_amount_check CHECK (amount > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS volume_relation_uuid_idx ON volume_relation(uuid);
CREATE INDEX IF NOT EXISTS volume_relation_parent_volume_id_idx ON volume_relation(parent_volume_id);
CREATE INDEX IF NOT EXISTS volume_relation_child_volume_id_idx ON volume_relation(child_volume_id);

CREATE TABLE IF NOT EXISTS vessel (
    id             serial PRIMARY KEY,
    uuid           uuid NOT NULL DEFAULT gen_random_uuid(),

    type           varchar(255) NOT NULL,
    name           varchar(255) NOT NULL,
    capacity       bigint NOT NULL,
    capacity_unit  varchar(7) NOT NULL,
    make           varchar(255),
    model          varchar(255),
    status         varchar(32) NOT NULL DEFAULT 'active',

    created_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at     timestamptz,
    CONSTRAINT vessel_status_check CHECK (status IN ('active', 'inactive', 'retired'))
);

CREATE UNIQUE INDEX IF NOT EXISTS vessel_uuid_idx ON vessel(uuid);

CREATE TABLE IF NOT EXISTS occupancy (
    id          serial PRIMARY KEY,
    uuid        uuid NOT NULL DEFAULT gen_random_uuid(),

    vessel_id   int NOT NULL REFERENCES vessel(id),
    volume_id   int NOT NULL REFERENCES volume(id),
    in_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    out_at      timestamptz,

    created_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at  timestamptz,
    CONSTRAINT occupancy_out_after_in_check CHECK (out_at IS NULL OR out_at >= in_at)
);

CREATE UNIQUE INDEX IF NOT EXISTS occupancy_uuid_idx ON occupancy(uuid);
CREATE INDEX IF NOT EXISTS occupancy_vessel_id_idx ON occupancy(vessel_id);
CREATE INDEX IF NOT EXISTS occupancy_volume_id_idx ON occupancy(volume_id);
CREATE UNIQUE INDEX IF NOT EXISTS occupancy_current_vessel_idx ON occupancy(vessel_id) WHERE out_at IS NULL;

CREATE TABLE IF NOT EXISTS transfer (
    id                   serial PRIMARY KEY,
    uuid                 uuid NOT NULL DEFAULT gen_random_uuid(),

    source_occupancy_id  int NOT NULL REFERENCES occupancy(id),
    dest_occupancy_id    int NOT NULL REFERENCES occupancy(id),
    amount               bigint NOT NULL,
    amount_unit          varchar(7) NOT NULL,
    loss_amount          bigint,
    loss_unit            varchar(7),
    started_at           timestamptz NOT NULL DEFAULT timezone('utc', now()),
    ended_at             timestamptz,

    created_at           timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at           timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at           timestamptz,
    CONSTRAINT transfer_amount_check CHECK (amount > 0),
    CONSTRAINT transfer_source_dest_check CHECK (source_occupancy_id <> dest_occupancy_id),
    CONSTRAINT transfer_end_after_start_check CHECK (ended_at IS NULL OR ended_at >= started_at)
);

CREATE UNIQUE INDEX IF NOT EXISTS transfer_uuid_idx ON transfer(uuid);
CREATE INDEX IF NOT EXISTS transfer_source_occupancy_id_idx ON transfer(source_occupancy_id);
CREATE INDEX IF NOT EXISTS transfer_dest_occupancy_id_idx ON transfer(dest_occupancy_id);

CREATE TABLE IF NOT EXISTS batch_volume (
    id          serial PRIMARY KEY,
    uuid        uuid NOT NULL DEFAULT gen_random_uuid(),

    batch_id    int NOT NULL REFERENCES batch(id),
    volume_id   int NOT NULL REFERENCES volume(id),
    liquid_phase varchar(16) NOT NULL,
    phase_at    timestamptz NOT NULL DEFAULT timezone('utc', now()),

    created_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at  timestamptz,
    CONSTRAINT batch_volume_liquid_phase_check CHECK (liquid_phase IN (
        'water',
        'wort',
        'beer'
    ))
);

CREATE UNIQUE INDEX IF NOT EXISTS batch_volume_uuid_idx ON batch_volume(uuid);
CREATE INDEX IF NOT EXISTS batch_volume_batch_id_idx ON batch_volume(batch_id);
CREATE INDEX IF NOT EXISTS batch_volume_volume_id_idx ON batch_volume(volume_id);

CREATE TABLE IF NOT EXISTS batch_process_phase (
    id            serial PRIMARY KEY,
    uuid          uuid NOT NULL DEFAULT gen_random_uuid(),

    batch_id      int NOT NULL REFERENCES batch(id),
    process_phase varchar(32) NOT NULL,
    phase_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),

    created_at    timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at    timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at    timestamptz,
    CONSTRAINT batch_process_phase_check CHECK (process_phase IN (
        'planning',
        'mashing',
        'heating',
        'boiling',
        'cooling',
        'fermenting',
        'conditioning',
        'packaging',
        'finished'
    ))
);

CREATE UNIQUE INDEX IF NOT EXISTS batch_process_phase_uuid_idx ON batch_process_phase(uuid);
CREATE INDEX IF NOT EXISTS batch_process_phase_batch_id_idx ON batch_process_phase(batch_id);
CREATE INDEX IF NOT EXISTS batch_process_phase_phase_at_idx ON batch_process_phase(phase_at);

CREATE TABLE IF NOT EXISTS batch_relation (
    id               serial PRIMARY KEY,
    uuid             uuid NOT NULL DEFAULT gen_random_uuid(),

    parent_batch_id  int NOT NULL REFERENCES batch(id),
    child_batch_id   int NOT NULL REFERENCES batch(id),
    relation_type    varchar(32) NOT NULL,
    volume_id        int REFERENCES volume(id),

    created_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at       timestamptz,
    CONSTRAINT batch_relation_parent_child_check CHECK (parent_batch_id <> child_batch_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS batch_relation_uuid_idx ON batch_relation(uuid);
CREATE INDEX IF NOT EXISTS batch_relation_parent_batch_id_idx ON batch_relation(parent_batch_id);
CREATE INDEX IF NOT EXISTS batch_relation_child_batch_id_idx ON batch_relation(child_batch_id);

CREATE TABLE IF NOT EXISTS addition (
    id                  serial PRIMARY KEY,
    uuid                uuid NOT NULL DEFAULT gen_random_uuid(),

    batch_id            int REFERENCES batch(id),
    occupancy_id        int REFERENCES occupancy(id),
    addition_type       varchar(32) NOT NULL,
    stage               varchar(32),
    inventory_lot_uuid  uuid,
    amount              bigint NOT NULL,
    amount_unit         varchar(7) NOT NULL,
    added_at            timestamptz NOT NULL DEFAULT timezone('utc', now()),
    notes               text,

    created_at          timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at          timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at          timestamptz,
    CONSTRAINT addition_target_check CHECK (
        (batch_id IS NOT NULL AND occupancy_id IS NULL) OR
        (batch_id IS NULL AND occupancy_id IS NOT NULL)
    ),
    CONSTRAINT addition_type_check CHECK (addition_type IN (
        'malt',
        'hop',
        'yeast',
        'adjunct',
        'water_chem',
        'gas',
        'other'
    )),
    CONSTRAINT addition_inventory_lot_check CHECK (
        (addition_type IN ('malt', 'hop', 'yeast', 'adjunct') AND inventory_lot_uuid IS NOT NULL)
        OR (addition_type NOT IN ('malt', 'hop', 'yeast', 'adjunct'))
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS addition_uuid_idx ON addition(uuid);
CREATE INDEX IF NOT EXISTS addition_batch_id_idx ON addition(batch_id);
CREATE INDEX IF NOT EXISTS addition_occupancy_id_idx ON addition(occupancy_id);
CREATE INDEX IF NOT EXISTS addition_inventory_lot_uuid_idx ON addition(inventory_lot_uuid);

CREATE TABLE IF NOT EXISTS measurement (
    id           serial PRIMARY KEY,
    uuid         uuid NOT NULL DEFAULT gen_random_uuid(),

    batch_id     int REFERENCES batch(id),
    occupancy_id int REFERENCES occupancy(id),
    kind         varchar(32) NOT NULL,
    value        numeric(12,4) NOT NULL,
    unit         varchar(16),
    observed_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    notes        text,

    created_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at   timestamptz,
    CONSTRAINT measurement_target_check CHECK (
        (batch_id IS NOT NULL AND occupancy_id IS NULL) OR
        (batch_id IS NULL AND occupancy_id IS NOT NULL)
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS measurement_uuid_idx ON measurement(uuid);
CREATE INDEX IF NOT EXISTS measurement_batch_id_idx ON measurement(batch_id);
CREATE INDEX IF NOT EXISTS measurement_occupancy_id_idx ON measurement(occupancy_id);
CREATE INDEX IF NOT EXISTS measurement_observed_at_idx ON measurement(observed_at);
