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

-- Seed data for early development.
INSERT INTO batch (uuid, short_name, brew_date, notes)
VALUES
    ('90000000-0000-0000-0000-000000000001', 'IPA 24-07', '2026-01-10', 'Flagship IPA split into two fermenters.'),
    ('90000000-0000-0000-0000-000000000002', 'IPA 24-07B', '2026-01-10', 'Split portion for dry hop trial.'),
    ('90000000-0000-0000-0000-000000000003', 'Pilsner 24-08', '2026-01-18', 'Crisp lager pilot batch.');

INSERT INTO volume (uuid, name, description, amount, amount_unit)
VALUES
    ('91000000-0000-0000-0000-000000000001', 'IPA 24-07 wort', 'Pre-fermentation wort.', 1200, 'l'),
    ('91000000-0000-0000-0000-000000000002', 'IPA 24-07 split A', 'Fermenter A portion.', 600, 'l'),
    ('91000000-0000-0000-0000-000000000003', 'IPA 24-07 split B', 'Fermenter B portion.', 590, 'l'),
    ('91000000-0000-0000-0000-000000000004', 'Pilsner 24-08 wort', 'Pilsner wort.', 800, 'l');

INSERT INTO volume_relation (uuid, parent_volume_id, child_volume_id, relation_type, amount, amount_unit)
VALUES
    ('92000000-0000-0000-0000-000000000001', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000001'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000002'), 'split', 600, 'l'),
    ('92000000-0000-0000-0000-000000000002', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000001'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000003'), 'split', 590, 'l');

INSERT INTO vessel (uuid, type, name, capacity, capacity_unit, make, model, status)
VALUES
    ('93000000-0000-0000-0000-000000000001', 'mash_tun', 'Mash Tun 1', 1500, 'l', 'Stout Tanks', 'MT-15', 'active'),
    ('93000000-0000-0000-0000-000000000002', 'kettle', 'Kettle 1', 1500, 'l', 'Stout Tanks', 'K-15', 'active'),
    ('93000000-0000-0000-0000-000000000003', 'fermenter', 'Fermenter A', 1200, 'l', 'SS Brewtech', 'FV-12', 'active'),
    ('93000000-0000-0000-0000-000000000004', 'fermenter', 'Fermenter B', 1200, 'l', 'SS Brewtech', 'FV-12', 'active'),
    ('93000000-0000-0000-0000-000000000005', 'fermenter', 'Fermenter C', 1000, 'l', 'SS Brewtech', 'FV-10', 'active'),
    ('93000000-0000-0000-0000-000000000006', 'brite_tank', 'Brite Tank 1', 1000, 'l', 'Premier Stainless', 'BT-10', 'active');

INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at)
VALUES
    ('93100000-0000-0000-0000-000000000001', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000001'), '2026-01-10 08:00:00+00', '2026-01-10 10:00:00+00'),
    ('93100000-0000-0000-0000-000000000002', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000001'), '2026-01-10 10:30:00+00', '2026-01-10 12:00:00+00'),
    ('93100000-0000-0000-0000-000000000003', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000003'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000002'), '2026-01-10 13:00:00+00', NULL),
    ('93100000-0000-0000-0000-000000000004', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000004'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000003'), '2026-01-10 13:15:00+00', NULL),
    ('93100000-0000-0000-0000-000000000005', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000005'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000004'), '2026-01-18 12:00:00+00', NULL);

INSERT INTO transfer (
    uuid,
    source_occupancy_id,
    dest_occupancy_id,
    amount,
    amount_unit,
    loss_amount,
    loss_unit,
    started_at,
    ended_at
)
VALUES
    ('93200000-0000-0000-0000-000000000001', (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000001'), (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000002'), 1200, 'l', 5, 'l', '2026-01-10 10:05:00+00', '2026-01-10 10:25:00+00'),
    ('93200000-0000-0000-0000-000000000002', (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000002'), (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000003'), 600, 'l', 5, 'l', '2026-01-10 12:10:00+00', '2026-01-10 12:40:00+00'),
    ('93200000-0000-0000-0000-000000000003', (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000002'), (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000004'), 590, 'l', 5, 'l', '2026-01-10 12:10:00+00', '2026-01-10 12:45:00+00');

INSERT INTO batch_volume (uuid, batch_id, volume_id, liquid_phase, phase_at)
VALUES
    ('94000000-0000-0000-0000-000000000001', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000001'), 'wort', '2026-01-10 09:00:00+00'),
    ('94000000-0000-0000-0000-000000000002', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000002'), 'beer', '2026-01-20 12:00:00+00'),
    ('94000000-0000-0000-0000-000000000003', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000002'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000003'), 'beer', '2026-01-20 12:00:00+00'),
    ('94000000-0000-0000-0000-000000000004', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000004'), 'wort', '2026-01-18 11:00:00+00');

INSERT INTO batch_process_phase (uuid, batch_id, process_phase, phase_at)
VALUES
    ('95000000-0000-0000-0000-000000000001', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), 'planning', '2026-01-05 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000002', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), 'mashing', '2026-01-10 08:00:00+00'),
    ('95000000-0000-0000-0000-000000000003', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), 'boiling', '2026-01-10 10:45:00+00'),
    ('95000000-0000-0000-0000-000000000004', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), 'fermenting', '2026-01-10 13:30:00+00'),
    ('95000000-0000-0000-0000-000000000005', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000002'), 'fermenting', '2026-01-10 13:45:00+00'),
    ('95000000-0000-0000-0000-000000000006', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), 'planning', '2026-01-12 09:00:00+00');

INSERT INTO batch_relation (uuid, parent_batch_id, child_batch_id, relation_type, volume_id)
VALUES
    ('96000000-0000-0000-0000-000000000001', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000002'), 'split', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000003'));

INSERT INTO addition (
    uuid,
    batch_id,
    occupancy_id,
    addition_type,
    stage,
    inventory_lot_uuid,
    amount,
    amount_unit,
    added_at,
    notes
)
VALUES
    ('97000000-0000-0000-0000-000000000001', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'malt', 'mash', '80000000-0000-0000-0000-000000000001', 220, 'kg', '2026-01-10 08:30:00+00', 'Base malt for IPA mash.'),
    ('97000000-0000-0000-0000-000000000002', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'hop', 'boil', '80000000-0000-0000-0000-000000000004', 5, 'kg', '2026-01-10 11:00:00+00', 'Bittering addition.'),
    ('97000000-0000-0000-0000-000000000003', NULL, (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000003'), 'yeast', 'fermentation', '80000000-0000-0000-0000-000000000006', 2, 'kg', '2026-01-10 13:20:00+00', 'Pitch WLP001.'),
    ('97000000-0000-0000-0000-000000000004', NULL, (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000004'), 'hop', 'dry_hop', '80000000-0000-0000-0000-000000000005', 3, 'kg', '2026-01-14 11:00:00+00', 'Dry hop trial in split B.'),
    ('97000000-0000-0000-0000-000000000005', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'water_chem', 'mash', NULL, 500, 'g', '2026-01-10 08:10:00+00', 'Gypsum addition.');

INSERT INTO measurement (
    uuid,
    batch_id,
    occupancy_id,
    kind,
    value,
    unit,
    observed_at,
    notes
)
VALUES
    ('98000000-0000-0000-0000-000000000001', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'gravity', 1.0540, 'sg', '2026-01-10 12:00:00+00', 'Pre-fermentation gravity.'),
    ('98000000-0000-0000-0000-000000000002', NULL, (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000003'), 'temperature', 20.5, 'c', '2026-01-12 09:00:00+00', 'Fermentation temp.'),
    ('98000000-0000-0000-0000-000000000003', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'ph', 5.20, NULL, '2026-01-10 09:15:00+00', 'Mash pH.');
