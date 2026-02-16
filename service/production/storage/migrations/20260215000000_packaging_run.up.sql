-- Packaging run tables: package_format, packaging_run, packaging_run_line.
-- Production service owns packaging runs. Inventory service owns the resulting beer lots.
-- Production calls Inventory (inter-service) to create beer lots when a packaging run is recorded.
BEGIN;

-- ==============================================================================
-- Package Format: reference table for container types (follows style table pattern)
-- ==============================================================================

CREATE TABLE IF NOT EXISTS package_format (
    id                    serial PRIMARY KEY,
    uuid                  uuid NOT NULL DEFAULT gen_random_uuid(),

    name                  varchar(255) NOT NULL,
    container             varchar(32) NOT NULL,
    volume_per_unit       bigint NOT NULL,
    volume_per_unit_unit  varchar(7) NOT NULL,
    is_active             boolean NOT NULL DEFAULT true,

    created_at            timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at            timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at            timestamptz,

    CONSTRAINT package_format_volume_per_unit_check CHECK (volume_per_unit > 0),
    CONSTRAINT package_format_container_check CHECK (container IN (
        'keg', 'can', 'bottle', 'cask', 'growler', 'other'
    ))
);

CREATE UNIQUE INDEX IF NOT EXISTS package_format_uuid_idx ON package_format(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS package_format_name_lower_idx ON package_format(lower(name)) WHERE deleted_at IS NULL;

-- ==============================================================================
-- Packaging Run: a packaging event for a batch (follows transfer table pattern)
-- ==============================================================================

CREATE TABLE IF NOT EXISTS packaging_run (
    id              serial PRIMARY KEY,
    uuid            uuid NOT NULL DEFAULT gen_random_uuid(),

    batch_id        int NOT NULL REFERENCES batch(id),
    occupancy_id    int NOT NULL REFERENCES occupancy(id),
    started_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),
    ended_at        timestamptz,
    loss_amount     bigint,
    loss_unit       varchar(7),
    notes           text,

    created_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at      timestamptz,

    CONSTRAINT packaging_run_end_after_start_check CHECK (ended_at IS NULL OR ended_at >= started_at),
    CONSTRAINT packaging_run_loss_pair_check CHECK (
        (loss_amount IS NULL AND loss_unit IS NULL) OR
        (loss_amount IS NOT NULL AND loss_unit IS NOT NULL)
    ),
    CONSTRAINT packaging_run_loss_amount_check CHECK (loss_amount IS NULL OR loss_amount >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS packaging_run_uuid_idx ON packaging_run(uuid);
CREATE INDEX IF NOT EXISTS packaging_run_batch_id_idx ON packaging_run(batch_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS packaging_run_occupancy_id_idx ON packaging_run(occupancy_id) WHERE deleted_at IS NULL;

-- ==============================================================================
-- Packaging Run Line: one line per format in a packaging run
-- ==============================================================================

CREATE TABLE IF NOT EXISTS packaging_run_line (
    id                  serial PRIMARY KEY,
    uuid                uuid NOT NULL DEFAULT gen_random_uuid(),

    packaging_run_id    int NOT NULL REFERENCES packaging_run(id) ON DELETE CASCADE,
    package_format_id   int NOT NULL REFERENCES package_format(id),
    quantity            int NOT NULL,

    created_at          timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at          timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at          timestamptz,

    CONSTRAINT packaging_run_line_quantity_check CHECK (quantity > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS packaging_run_line_uuid_idx ON packaging_run_line(uuid);
CREATE INDEX IF NOT EXISTS packaging_run_line_packaging_run_id_idx ON packaging_run_line(packaging_run_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS packaging_run_line_package_format_id_idx ON packaging_run_line(package_format_id) WHERE deleted_at IS NULL;

-- Prevent duplicate format per run (one line per format per packaging run)
CREATE UNIQUE INDEX IF NOT EXISTS packaging_run_line_run_format_idx
    ON packaging_run_line(packaging_run_id, package_format_id) WHERE deleted_at IS NULL;

-- ==============================================================================
-- Seed Data: Package Formats (deterministic UUIDs starting with b0000000-)
-- ==============================================================================

INSERT INTO package_format (uuid, name, container, volume_per_unit, volume_per_unit_unit) VALUES
    ('b0000000-0000-0000-0000-000000000001', '1/2 BBL Keg',  'keg',     1984, 'usfloz'),
    ('b0000000-0000-0000-0000-000000000002', '1/4 BBL Keg',  'keg',      992, 'usfloz'),
    ('b0000000-0000-0000-0000-000000000003', '1/6 BBL Keg',  'keg',    19558, 'ml'),
    ('b0000000-0000-0000-0000-000000000004', '16oz Can',      'can',       16, 'usfloz'),
    ('b0000000-0000-0000-0000-000000000005', '12oz Can',      'can',       12, 'usfloz'),
    ('b0000000-0000-0000-0000-000000000006', '32oz Crowler',  'can',       32, 'usfloz'),
    ('b0000000-0000-0000-0000-000000000007', '12oz Bottle',   'bottle',    12, 'usfloz');

COMMIT;
