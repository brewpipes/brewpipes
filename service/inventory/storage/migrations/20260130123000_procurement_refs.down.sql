DROP INDEX IF EXISTS ingredient_lot_purchase_order_line_uuid_idx;
ALTER TABLE ingredient_lot
    DROP COLUMN IF EXISTS purchase_order_line_uuid;

DROP INDEX IF EXISTS inventory_receipt_purchase_order_uuid_idx;
ALTER TABLE inventory_receipt
    DROP COLUMN IF EXISTS purchase_order_uuid;
