BEGIN;

-- ==============================================================================
-- 1. inventory_removal table
-- ==============================================================================

CREATE TABLE IF NOT EXISTS inventory_removal (
    id               serial PRIMARY KEY,
    uuid             uuid NOT NULL DEFAULT gen_random_uuid(),

    -- Classification
    category         varchar(32) NOT NULL,
    reason           varchar(32) NOT NULL,

    -- What was removed (at least one must be set)
    batch_uuid       uuid,
    beer_lot_id      int REFERENCES beer_lot(id),
    occupancy_uuid   uuid,

    -- Quantity
    amount           bigint NOT NULL,
    amount_unit      varchar(7) NOT NULL,
    amount_bbl       numeric(10,5),

    -- TTB compliance
    is_taxable       boolean NOT NULL DEFAULT false,

    -- Operational metadata
    reference_code   varchar(64),
    performed_by     varchar(255),
    removed_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    destination      text,
    notes            text,

    -- Movement back-reference (for beer lot removals)
    movement_id      int,

    -- Standard timestamps
    created_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at       timestamptz,

    CONSTRAINT inventory_removal_category_check CHECK (category IN (
        'dump', 'waste', 'sample', 'expired', 'other'
    )),
    CONSTRAINT inventory_removal_reason_check CHECK (reason IN (
        'infection', 'off_flavor', 'failed_fermentation', 'equipment_failure',
        'quality_reject', 'past_date', 'damaged_package', 'spillage',
        'cleaning', 'qc_sample', 'tasting', 'competition', 'other'
    )),
    CONSTRAINT inventory_removal_reference_required_check CHECK (
        batch_uuid IS NOT NULL OR beer_lot_id IS NOT NULL OR occupancy_uuid IS NOT NULL
    ),
    CONSTRAINT inventory_removal_amount_check CHECK (amount > 0),
    CONSTRAINT inventory_removal_amount_bbl_check CHECK (amount_bbl IS NULL OR amount_bbl > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS inventory_removal_uuid_idx ON inventory_removal(uuid);
CREATE INDEX IF NOT EXISTS inventory_removal_batch_uuid_idx ON inventory_removal(batch_uuid) WHERE batch_uuid IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS inventory_removal_occupancy_uuid_idx ON inventory_removal(occupancy_uuid) WHERE occupancy_uuid IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS inventory_removal_beer_lot_id_idx ON inventory_removal(beer_lot_id) WHERE beer_lot_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS inventory_removal_category_idx ON inventory_removal(category) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS inventory_removal_removed_at_idx ON inventory_removal(removed_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS inventory_removal_reference_code_idx ON inventory_removal(reference_code) WHERE reference_code IS NOT NULL AND deleted_at IS NULL;

-- ==============================================================================
-- 2. inventory_movement modifications
-- ==============================================================================

ALTER TABLE inventory_movement ADD COLUMN removal_id int REFERENCES inventory_removal(id);
CREATE INDEX IF NOT EXISTS inventory_movement_removal_id_idx ON inventory_movement(removal_id);

ALTER TABLE inventory_movement DROP CONSTRAINT inventory_movement_reason_check;
ALTER TABLE inventory_movement ADD CONSTRAINT inventory_movement_reason_check CHECK (reason IN (
    'receive', 'use', 'transfer', 'adjust', 'waste', 'package', 'removal'
));

ALTER TABLE inventory_movement DROP CONSTRAINT inventory_movement_reason_reference_check;
ALTER TABLE inventory_movement ADD CONSTRAINT inventory_movement_reason_reference_check CHECK (
    (reason = 'receive' AND receipt_id IS NOT NULL) OR
    (reason = 'use' AND usage_id IS NOT NULL) OR
    (reason = 'transfer' AND transfer_id IS NOT NULL) OR
    (reason IN ('adjust', 'waste') AND adjustment_id IS NOT NULL) OR
    (reason = 'package') OR
    (reason = 'removal' AND removal_id IS NOT NULL)
);

ALTER TABLE inventory_movement DROP CONSTRAINT inventory_movement_reference_check;
ALTER TABLE inventory_movement ADD CONSTRAINT inventory_movement_reference_check CHECK (
    num_nonnulls(receipt_id, usage_id, adjustment_id, transfer_id, removal_id) <= 1
);

-- ==============================================================================
-- 3. Add deferred FK from removal.movement_id → inventory_movement.id
-- (Must be after movement table is altered)
-- ==============================================================================

ALTER TABLE inventory_removal ADD CONSTRAINT inventory_removal_movement_id_fk
    FOREIGN KEY (movement_id) REFERENCES inventory_movement(id);

-- ==============================================================================
-- 4. Seed data
-- ==============================================================================

-- Removal 1: Batch dump — infected Kolsch (in-process, no movement)
INSERT INTO inventory_removal (
    uuid, category, reason, batch_uuid, occupancy_uuid,
    amount, amount_unit, amount_bbl, is_taxable,
    reference_code, performed_by, removed_at, notes
) VALUES (
    '87000000-0000-0000-0000-000000000001',
    'dump', 'infection',
    '90000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000010',
    800, 'l', 6.81555, false,
    'RMV-0001', 'Jake (Head Brewer)', '2026-02-10 09:00:00+00',
    'Kolsch 24-10 dumped — lactobacillus contamination detected at day 3.'
);

-- Removal 2: QC sample from packaged IPA (creates movement)
INSERT INTO inventory_removal (
    uuid, category, reason, batch_uuid, beer_lot_id,
    amount, amount_unit, amount_bbl, is_taxable,
    reference_code, performed_by, removed_at, notes
) VALUES (
    '87000000-0000-0000-0000-000000000002',
    'sample', 'qc_sample',
    '90000000-0000-0000-0000-000000000001',
    (SELECT id FROM beer_lot WHERE uuid = '86000000-0000-0000-0000-000000000001'),
    2, 'l', 0.01704, false,
    'RMV-0002', 'Maria (QC Lead)', '2026-02-11 14:30:00+00',
    'QC retention samples from IPA 24-07 packaging run.'
);

-- Movement for Removal 2
INSERT INTO inventory_movement (
    beer_lot_id, stock_location_id, direction, reason,
    amount, amount_unit, occurred_at, removal_id
) VALUES (
    (SELECT id FROM beer_lot WHERE uuid = '86000000-0000-0000-0000-000000000001'),
    (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000006'),
    'out', 'removal', 2, 'l', '2026-02-11 14:30:00+00',
    (SELECT id FROM inventory_removal WHERE uuid = '87000000-0000-0000-0000-000000000002')
);

-- Back-reference movement on removal 2
UPDATE inventory_removal SET movement_id = (
    SELECT m.id FROM inventory_movement m
    WHERE m.removal_id = (SELECT id FROM inventory_removal WHERE uuid = '87000000-0000-0000-0000-000000000002')
    LIMIT 1
) WHERE uuid = '87000000-0000-0000-0000-000000000002';

-- Removal 3: Waste from tank cleaning (in-process, no movement)
INSERT INTO inventory_removal (
    uuid, category, reason, batch_uuid, occupancy_uuid,
    amount, amount_unit, amount_bbl, is_taxable,
    reference_code, performed_by, removed_at, notes
) VALUES (
    '87000000-0000-0000-0000-000000000003',
    'waste', 'cleaning',
    '90000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000011',
    15, 'l', 0.12783, false,
    'RMV-0003', 'Jake (Head Brewer)', '2026-02-08 16:00:00+00',
    'Residual beer in brite tank lines after Stout 24-09 transfer. CIP flush waste.'
);

-- Removal 4: Expired stout destruction (creates movement)
INSERT INTO inventory_removal (
    uuid, category, reason, batch_uuid, beer_lot_id,
    amount, amount_unit, amount_bbl, is_taxable,
    reference_code, performed_by, removed_at, notes
) VALUES (
    '87000000-0000-0000-0000-000000000004',
    'expired', 'past_date',
    '90000000-0000-0000-0000-000000000004',
    (SELECT id FROM beer_lot WHERE uuid = '86000000-0000-0000-0000-000000000002'),
    24, 'l', 0.20449, false,
    'RMV-0004', 'Maria (QC Lead)', '2026-02-14 10:00:00+00',
    'One case of Stout 24-09 cans past best-by. Drain pour witnessed by shift lead.'
);

-- Movement for Removal 4
INSERT INTO inventory_movement (
    beer_lot_id, stock_location_id, direction, reason,
    amount, amount_unit, occurred_at, removal_id
) VALUES (
    (SELECT id FROM beer_lot WHERE uuid = '86000000-0000-0000-0000-000000000002'),
    (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000006'),
    'out', 'removal', 24, 'l', '2026-02-14 10:00:00+00',
    (SELECT id FROM inventory_removal WHERE uuid = '87000000-0000-0000-0000-000000000004')
);

-- Back-reference movement on removal 4
UPDATE inventory_removal SET movement_id = (
    SELECT m.id FROM inventory_movement m
    WHERE m.removal_id = (SELECT id FROM inventory_removal WHERE uuid = '87000000-0000-0000-0000-000000000004')
    LIMIT 1
) WHERE uuid = '87000000-0000-0000-0000-000000000004';

COMMIT;
