BEGIN;
DELETE FROM inventory_movement WHERE removal_id IS NOT NULL;
UPDATE inventory_removal SET movement_id = NULL;
DELETE FROM inventory_removal;

ALTER TABLE inventory_movement DROP CONSTRAINT inventory_movement_reference_check;
ALTER TABLE inventory_movement ADD CONSTRAINT inventory_movement_reference_check CHECK (
    num_nonnulls(receipt_id, usage_id, adjustment_id, transfer_id) <= 1
);
ALTER TABLE inventory_movement DROP CONSTRAINT inventory_movement_reason_reference_check;
ALTER TABLE inventory_movement ADD CONSTRAINT inventory_movement_reason_reference_check CHECK (
    (reason = 'receive' AND receipt_id IS NOT NULL) OR
    (reason = 'use' AND usage_id IS NOT NULL) OR
    (reason = 'transfer' AND transfer_id IS NOT NULL) OR
    (reason IN ('adjust', 'waste') AND adjustment_id IS NOT NULL) OR
    (reason = 'package')
);
ALTER TABLE inventory_movement DROP CONSTRAINT inventory_movement_reason_check;
ALTER TABLE inventory_movement ADD CONSTRAINT inventory_movement_reason_check CHECK (reason IN (
    'receive', 'use', 'transfer', 'adjust', 'waste', 'package'
));
DROP INDEX IF EXISTS inventory_movement_removal_id_idx;
ALTER TABLE inventory_movement DROP COLUMN IF EXISTS removal_id;

ALTER TABLE inventory_removal DROP CONSTRAINT IF EXISTS inventory_removal_movement_id_fk;
DROP TABLE IF EXISTS inventory_removal CASCADE;
COMMIT;
