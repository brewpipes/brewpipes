-- Reverse beer lot packaging support.
BEGIN;

-- ==============================================================================
-- 1. Restore inventory_movement reason constraints (remove 'package')
-- ==============================================================================

ALTER TABLE inventory_movement DROP CONSTRAINT inventory_movement_reason_reference_check;
ALTER TABLE inventory_movement ADD CONSTRAINT inventory_movement_reason_reference_check CHECK (
    (reason = 'receive' AND receipt_id IS NOT NULL) OR
    (reason = 'use' AND usage_id IS NOT NULL) OR
    (reason = 'transfer' AND transfer_id IS NOT NULL) OR
    (reason IN ('adjust', 'waste') AND adjustment_id IS NOT NULL)
);

ALTER TABLE inventory_movement DROP CONSTRAINT inventory_movement_reason_check;
ALTER TABLE inventory_movement ADD CONSTRAINT inventory_movement_reason_check CHECK (reason IN (
    'receive',
    'use',
    'transfer',
    'adjust',
    'waste'
));

-- ==============================================================================
-- 2. Drop beer_lot_item table
-- ==============================================================================

DROP TABLE IF EXISTS beer_lot_item CASCADE;

-- ==============================================================================
-- 3. Remove new columns and constraints from beer_lot
-- ==============================================================================

ALTER TABLE beer_lot DROP CONSTRAINT IF EXISTS beer_lot_quantity_check;
ALTER TABLE beer_lot DROP CONSTRAINT IF EXISTS beer_lot_volume_per_unit_pair_check;
ALTER TABLE beer_lot DROP CONSTRAINT IF EXISTS beer_lot_volume_per_unit_check;
ALTER TABLE beer_lot DROP CONSTRAINT IF EXISTS beer_lot_container_check;

DROP INDEX IF EXISTS beer_lot_container_idx;
DROP INDEX IF EXISTS beer_lot_packaging_run_uuid_idx;

ALTER TABLE beer_lot DROP COLUMN IF EXISTS quantity;
ALTER TABLE beer_lot DROP COLUMN IF EXISTS volume_per_unit_unit;
ALTER TABLE beer_lot DROP COLUMN IF EXISTS volume_per_unit;
ALTER TABLE beer_lot DROP COLUMN IF EXISTS container;
ALTER TABLE beer_lot DROP COLUMN IF EXISTS package_format_name;
ALTER TABLE beer_lot DROP COLUMN IF EXISTS best_by;
ALTER TABLE beer_lot DROP COLUMN IF EXISTS packaging_run_uuid;

COMMIT;
