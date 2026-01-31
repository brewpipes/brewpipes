-- Link inventory receipts and lots to procurement records.
-- Cross-domain references are stored as opaque UUIDs without foreign keys.

ALTER TABLE inventory_receipt
    ADD COLUMN IF NOT EXISTS purchase_order_uuid uuid;

CREATE INDEX IF NOT EXISTS inventory_receipt_purchase_order_uuid_idx
    ON inventory_receipt(purchase_order_uuid);

ALTER TABLE ingredient_lot
    ADD COLUMN IF NOT EXISTS purchase_order_line_uuid uuid;

CREATE INDEX IF NOT EXISTS ingredient_lot_purchase_order_line_uuid_idx
    ON ingredient_lot(purchase_order_line_uuid);
