-- Beer lot packaging support: new columns on beer_lot, beer_lot_item table,
-- and expanded inventory_movement reason constraint for 'package'.
BEGIN;

-- ==============================================================================
-- 1. New columns on beer_lot for packaging traceability
-- ==============================================================================

ALTER TABLE beer_lot ADD COLUMN packaging_run_uuid uuid;
ALTER TABLE beer_lot ADD COLUMN best_by timestamptz;
ALTER TABLE beer_lot ADD COLUMN package_format_name varchar(255);
ALTER TABLE beer_lot ADD COLUMN container varchar(32);
ALTER TABLE beer_lot ADD COLUMN volume_per_unit bigint;
ALTER TABLE beer_lot ADD COLUMN volume_per_unit_unit varchar(7);
ALTER TABLE beer_lot ADD COLUMN quantity int;

CREATE INDEX IF NOT EXISTS beer_lot_packaging_run_uuid_idx ON beer_lot(packaging_run_uuid) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS beer_lot_container_idx ON beer_lot(container) WHERE deleted_at IS NULL;

ALTER TABLE beer_lot ADD CONSTRAINT beer_lot_container_check CHECK (container IS NULL OR container IN ('keg', 'can', 'bottle', 'cask', 'growler', 'other'));
ALTER TABLE beer_lot ADD CONSTRAINT beer_lot_volume_per_unit_check CHECK (volume_per_unit IS NULL OR volume_per_unit > 0);
ALTER TABLE beer_lot ADD CONSTRAINT beer_lot_volume_per_unit_pair_check CHECK ((volume_per_unit IS NULL AND volume_per_unit_unit IS NULL) OR (volume_per_unit IS NOT NULL AND volume_per_unit_unit IS NOT NULL));
ALTER TABLE beer_lot ADD CONSTRAINT beer_lot_quantity_check CHECK (quantity IS NULL OR quantity > 0);

-- ==============================================================================
-- 2. Beer lot item table (structural seam for future per-unit tracking)
-- ==============================================================================

CREATE TABLE IF NOT EXISTS beer_lot_item (
    id            serial PRIMARY KEY,
    uuid          uuid NOT NULL DEFAULT gen_random_uuid(),
    beer_lot_id   int NOT NULL REFERENCES beer_lot(id),
    status        varchar(32) NOT NULL DEFAULT 'available',
    identifier    varchar(255),
    notes         text,
    created_at    timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at    timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at    timestamptz,
    CONSTRAINT beer_lot_item_status_check CHECK (status IN ('available', 'reserved', 'sold', 'returned', 'damaged', 'destroyed'))
);

CREATE UNIQUE INDEX IF NOT EXISTS beer_lot_item_uuid_idx ON beer_lot_item(uuid);
CREATE INDEX IF NOT EXISTS beer_lot_item_beer_lot_id_idx ON beer_lot_item(beer_lot_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS beer_lot_item_status_idx ON beer_lot_item(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS beer_lot_item_identifier_idx ON beer_lot_item(identifier) WHERE deleted_at IS NULL AND identifier IS NOT NULL;

-- ==============================================================================
-- 3. Expand inventory_movement reason constraints to include 'package'
-- ==============================================================================

ALTER TABLE inventory_movement DROP CONSTRAINT inventory_movement_reason_check;
ALTER TABLE inventory_movement ADD CONSTRAINT inventory_movement_reason_check CHECK (reason IN (
    'receive',
    'use',
    'transfer',
    'adjust',
    'waste',
    'package'
));

ALTER TABLE inventory_movement DROP CONSTRAINT inventory_movement_reason_reference_check;
ALTER TABLE inventory_movement ADD CONSTRAINT inventory_movement_reason_reference_check CHECK (
    (reason = 'receive' AND receipt_id IS NOT NULL) OR
    (reason = 'use' AND usage_id IS NOT NULL) OR
    (reason = 'transfer' AND transfer_id IS NOT NULL) OR
    (reason IN ('adjust', 'waste') AND adjustment_id IS NOT NULL) OR
    (reason = 'package')
);

COMMIT;
