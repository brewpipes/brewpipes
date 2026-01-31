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

INSERT INTO batch (uuid, short_name, brew_date, notes)
VALUES
    ('90000000-0000-0000-0000-000000000004', 'Stout 24-09', '2026-01-22', 'Robust stout finishing conditioning.'),
    ('90000000-0000-0000-0000-000000000005', 'Kolsch 24-10', '2026-01-24', 'Clean Kolsch fermented cool.'),
    ('90000000-0000-0000-0000-000000000006', 'Kolsch 24-10B', '2026-01-24', 'Split portion for late hopping.'),
    ('90000000-0000-0000-0000-000000000007', 'Kolsch 24-10 Blend', '2026-02-05', 'Recombined Kolsch blend after conditioning.'),
    ('90000000-0000-0000-0000-000000000008', 'Saison 24-11', '2026-01-26', 'Farmhouse saison with warm fermentation.');

INSERT INTO volume (uuid, name, description, amount, amount_unit)
VALUES
    ('91000000-0000-0000-0000-000000000001', 'IPA 24-07 wort', 'Pre-fermentation wort.', 1200, 'l'),
    ('91000000-0000-0000-0000-000000000002', 'IPA 24-07 split A', 'Fermenter A portion.', 600, 'l'),
    ('91000000-0000-0000-0000-000000000003', 'IPA 24-07 split B', 'Fermenter B portion.', 590, 'l'),
    ('91000000-0000-0000-0000-000000000004', 'Pilsner 24-08 wort', 'Pilsner wort.', 800, 'l');

INSERT INTO volume (uuid, name, description, amount, amount_unit)
VALUES
    ('91000000-0000-0000-0000-000000000005', 'Stout 24-09 wort', 'High-gravity stout wort.', 900, 'l'),
    ('91000000-0000-0000-0000-000000000006', 'Stout 24-09 beer', 'Conditioned stout.', 870, 'l'),
    ('91000000-0000-0000-0000-000000000007', 'Kolsch 24-10 wort', 'Kolsch wort before split.', 1000, 'l'),
    ('91000000-0000-0000-0000-000000000008', 'Kolsch 24-10 split A', 'Primary Kolsch fermenter portion.', 500, 'l'),
    ('91000000-0000-0000-0000-000000000009', 'Kolsch 24-10 split B', 'Late-hop trial portion.', 490, 'l'),
    ('91000000-0000-0000-0000-000000000010', 'Kolsch 24-10 blend', 'Recombined Kolsch blend.', 990, 'l'),
    ('91000000-0000-0000-0000-000000000011', 'Saison 24-11 wort', 'Saison wort before fermentation.', 950, 'l'),
    ('91000000-0000-0000-0000-000000000012', 'Saison 24-11 beer', 'Active saison fermentation.', 930, 'l');

INSERT INTO volume_relation (uuid, parent_volume_id, child_volume_id, relation_type, amount, amount_unit)
VALUES
    ('92000000-0000-0000-0000-000000000001', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000001'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000002'), 'split', 600, 'l'),
    ('92000000-0000-0000-0000-000000000002', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000001'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000003'), 'split', 590, 'l');

INSERT INTO volume_relation (uuid, parent_volume_id, child_volume_id, relation_type, amount, amount_unit)
VALUES
    ('92000000-0000-0000-0000-000000000003', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000007'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000008'), 'split', 500, 'l'),
    ('92000000-0000-0000-0000-000000000004', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000007'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000009'), 'split', 490, 'l'),
    ('92000000-0000-0000-0000-000000000005', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000008'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000010'), 'blend', 500, 'l'),
    ('92000000-0000-0000-0000-000000000006', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000009'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000010'), 'blend', 490, 'l');

INSERT INTO vessel (uuid, type, name, capacity, capacity_unit, make, model, status)
VALUES
    ('93000000-0000-0000-0000-000000000001', 'mash_tun', 'Mash Tun 1', 1500, 'l', 'Stout Tanks', 'MT-15', 'active'),
    ('93000000-0000-0000-0000-000000000002', 'kettle', 'Kettle 1', 1500, 'l', 'Stout Tanks', 'K-15', 'active'),
    ('93000000-0000-0000-0000-000000000003', 'fermenter', 'Fermenter A', 1200, 'l', 'SS Brewtech', 'FV-12', 'active'),
    ('93000000-0000-0000-0000-000000000004', 'fermenter', 'Fermenter B', 1200, 'l', 'SS Brewtech', 'FV-12', 'active'),
    ('93000000-0000-0000-0000-000000000005', 'fermenter', 'Fermenter C', 1000, 'l', 'SS Brewtech', 'FV-10', 'active'),
    ('93000000-0000-0000-0000-000000000006', 'brite_tank', 'Brite Tank 1', 1000, 'l', 'Premier Stainless', 'BT-10', 'active');

INSERT INTO vessel (uuid, type, name, capacity, capacity_unit, make, model, status)
VALUES
    ('93000000-0000-0000-0000-000000000007', 'fermenter', 'Fermenter D', 1100, 'l', 'SS Brewtech', 'FV-11', 'active'),
    ('93000000-0000-0000-0000-000000000008', 'brite_tank', 'Brite Tank 2', 1200, 'l', 'Premier Stainless', 'BT-12', 'active');

INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at)
VALUES
    ('93100000-0000-0000-0000-000000000001', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000001'), '2026-01-10 08:00:00+00', '2026-01-10 10:00:00+00'),
    ('93100000-0000-0000-0000-000000000002', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000001'), '2026-01-10 10:30:00+00', '2026-01-10 12:00:00+00'),
    ('93100000-0000-0000-0000-000000000003', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000003'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000002'), '2026-01-10 13:00:00+00', NULL),
    ('93100000-0000-0000-0000-000000000004', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000004'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000003'), '2026-01-10 13:15:00+00', NULL),
    ('93100000-0000-0000-0000-000000000005', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000005'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000004'), '2026-01-18 12:00:00+00', NULL);

INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at)
VALUES
    ('93100000-0000-0000-0000-000000000006', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000007'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000006'), '2026-01-22 18:00:00+00', '2026-02-02 10:00:00+00'),
    ('93100000-0000-0000-0000-000000000007', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000007'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000012'), '2026-01-26 19:00:00+00', NULL),
    ('93100000-0000-0000-0000-000000000008', (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000008'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000010'), '2026-02-05 09:00:00+00', NULL);

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

INSERT INTO batch_volume (uuid, batch_id, volume_id, liquid_phase, phase_at)
VALUES
    ('94000000-0000-0000-0000-000000000005', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000005'), 'wort', '2026-01-22 09:00:00+00'),
    ('94000000-0000-0000-0000-000000000006', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000006'), 'beer', '2026-02-01 12:00:00+00'),
    ('94000000-0000-0000-0000-000000000007', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000007'), 'wort', '2026-01-24 09:00:00+00'),
    ('94000000-0000-0000-0000-000000000008', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000008'), 'beer', '2026-02-02 10:00:00+00'),
    ('94000000-0000-0000-0000-000000000009', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000006'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000009'), 'beer', '2026-02-02 10:00:00+00'),
    ('94000000-0000-0000-0000-000000000010', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000007'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000010'), 'beer', '2026-02-05 09:00:00+00'),
    ('94000000-0000-0000-0000-000000000011', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000011'), 'wort', '2026-01-26 09:00:00+00'),
    ('94000000-0000-0000-0000-000000000012', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000012'), 'beer', '2026-01-30 12:00:00+00');

INSERT INTO batch_process_phase (uuid, batch_id, process_phase, phase_at)
VALUES
    ('95000000-0000-0000-0000-000000000001', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), 'planning', '2026-01-05 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000002', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), 'mashing', '2026-01-10 08:00:00+00'),
    ('95000000-0000-0000-0000-000000000003', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), 'boiling', '2026-01-10 10:45:00+00'),
    ('95000000-0000-0000-0000-000000000004', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), 'fermenting', '2026-01-10 13:30:00+00'),
    ('95000000-0000-0000-0000-000000000005', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000002'), 'fermenting', '2026-01-10 13:45:00+00'),
    ('95000000-0000-0000-0000-000000000006', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), 'planning', '2026-01-12 09:00:00+00');

INSERT INTO batch_process_phase (uuid, batch_id, process_phase, phase_at)
VALUES
    ('95000000-0000-0000-0000-000000000007', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), 'conditioning', '2026-01-18 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000008', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), 'packaging', '2026-01-25 15:00:00+00'),
    ('95000000-0000-0000-0000-000000000009', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), 'finished', '2026-01-26 12:00:00+00'),
    ('95000000-0000-0000-0000-000000000010', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000002'), 'conditioning', '2026-01-21 10:00:00+00'),
    ('95000000-0000-0000-0000-000000000011', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000002'), 'packaging', '2026-01-26 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000012', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), 'mashing', '2026-01-18 08:00:00+00'),
    ('95000000-0000-0000-0000-000000000013', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), 'boiling', '2026-01-18 10:30:00+00'),
    ('95000000-0000-0000-0000-000000000014', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), 'fermenting', '2026-01-18 14:00:00+00'),
    ('95000000-0000-0000-0000-000000000015', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), 'conditioning', '2026-01-28 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000016', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), 'planning', '2026-01-18 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000017', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), 'mashing', '2026-01-22 08:00:00+00'),
    ('95000000-0000-0000-0000-000000000018', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), 'boiling', '2026-01-22 11:00:00+00'),
    ('95000000-0000-0000-0000-000000000019', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), 'fermenting', '2026-01-22 18:00:00+00'),
    ('95000000-0000-0000-0000-000000000020', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), 'conditioning', '2026-01-30 10:00:00+00'),
    ('95000000-0000-0000-0000-000000000021', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), 'packaging', '2026-02-05 08:00:00+00'),
    ('95000000-0000-0000-0000-000000000022', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), 'finished', '2026-02-06 12:00:00+00'),
    ('95000000-0000-0000-0000-000000000023', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), 'planning', '2026-01-20 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000024', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), 'mashing', '2026-01-24 08:00:00+00'),
    ('95000000-0000-0000-0000-000000000025', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), 'boiling', '2026-01-24 10:30:00+00'),
    ('95000000-0000-0000-0000-000000000026', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), 'fermenting', '2026-01-24 14:00:00+00'),
    ('95000000-0000-0000-0000-000000000027', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), 'conditioning', '2026-02-01 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000028', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), 'packaging', '2026-02-06 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000029', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000006'), 'fermenting', '2026-01-24 14:15:00+00'),
    ('95000000-0000-0000-0000-000000000030', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000006'), 'conditioning', '2026-02-01 09:15:00+00'),
    ('95000000-0000-0000-0000-000000000031', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000007'), 'conditioning', '2026-02-05 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000032', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000007'), 'packaging', '2026-02-06 11:00:00+00'),
    ('95000000-0000-0000-0000-000000000033', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), 'planning', '2026-01-21 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000034', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), 'mashing', '2026-01-26 08:00:00+00'),
    ('95000000-0000-0000-0000-000000000035', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), 'boiling', '2026-01-26 10:30:00+00'),
    ('95000000-0000-0000-0000-000000000036', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), 'fermenting', '2026-01-26 19:00:00+00');

INSERT INTO batch_relation (uuid, parent_batch_id, child_batch_id, relation_type, volume_id)
VALUES
    ('96000000-0000-0000-0000-000000000001', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000002'), 'split', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000003'));

INSERT INTO batch_relation (uuid, parent_batch_id, child_batch_id, relation_type, volume_id)
VALUES
    ('96000000-0000-0000-0000-000000000002', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000006'), 'split', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000009')),
    ('96000000-0000-0000-0000-000000000003', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000007'), 'blend', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000010')),
    ('96000000-0000-0000-0000-000000000004', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000006'), (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000007'), 'blend', (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000010'));

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
    ('97000000-0000-0000-0000-000000000006', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'malt', 'mash', '80000000-0000-0000-0000-000000000008', 50, 'kg', '2026-01-22 08:20:00+00', 'Roasted barley for stout color.'),
    ('97000000-0000-0000-0000-000000000007', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'hop', 'boil', '80000000-0000-0000-0000-000000000009', 4, 'kg', '2026-01-24 10:45:00+00', 'Saaz kettle addition.'),
    ('97000000-0000-0000-0000-000000000008', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'yeast', 'fermentation', '80000000-0000-0000-0000-000000000010', 3, 'kg', '2026-01-24 13:40:00+00', 'W34/70 pitch.'),
    ('97000000-0000-0000-0000-000000000009', NULL, (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000007'), 'adjunct', 'fermentation', '80000000-0000-0000-0000-000000000012', 1, 'kg', '2026-01-28 10:00:00+00', 'Coriander addition.'),
    ('97000000-0000-0000-0000-000000000010', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'water_chem', 'mash', NULL, 300, 'g', '2026-01-22 08:05:00+00', 'Calcium chloride addition.');

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
    ('98000000-0000-0000-0000-000000000004', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'gravity', 1.0600, 'sg', '2026-01-22 12:10:00+00', 'Stout original gravity.'),
    ('98000000-0000-0000-0000-000000000005', NULL, (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000006'), 'temperature', 19.0, 'c', '2026-01-23 08:00:00+00', 'Stout fermenter temp.'),
    ('98000000-0000-0000-0000-000000000006', NULL, (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000006'), 'gravity', 1.0200, 'sg', '2026-01-30 08:00:00+00', 'Stout mid-fermentation gravity.'),
    ('98000000-0000-0000-0000-000000000007', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'gravity', 1.0480, 'sg', '2026-01-24 12:00:00+00', 'Kolsch original gravity.'),
    ('98000000-0000-0000-0000-000000000008', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'ibu', 28.0, 'ibu', '2026-01-24 12:30:00+00', 'Calculated kettle IBUs.'),
    ('98000000-0000-0000-0000-000000000009', NULL, (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000008'), 'co2', 2.45, 'vol', '2026-02-05 14:00:00+00', 'Brite tank carbonation check.'),
    ('98000000-0000-0000-0000-000000000010', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'gravity', 1.0520, 'sg', '2026-01-26 12:05:00+00', 'Saison original gravity.'),
    ('98000000-0000-0000-0000-000000000011', NULL, (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000007'), 'temperature', 26.5, 'c', '2026-01-27 09:00:00+00', 'Warm saison fermentation.'),
    ('98000000-0000-0000-0000-000000000012', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'gravity', 1.0440, 'sg', '2026-01-18 12:00:00+00', 'Pilsner original gravity.'),
    ('98000000-0000-0000-0000-000000000013', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'abv', 6.5, 'pct', '2026-01-26 12:00:00+00', 'Final IPA ABV.');
