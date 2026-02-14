-- Seed data enhancement: cancelled purchase order with line items and fee.
BEGIN;

-- 1. Cancelled purchase order (fills missing 'cancelled' status gap)
INSERT INTO purchase_order (uuid, supplier_id, order_number, status, ordered_at, expected_at, notes)
VALUES (
    '50000000-0000-0000-0000-000000000008',
    (SELECT id FROM supplier WHERE uuid = '66666666-6666-6666-6666-666666666666'),
    'PO-1008',
    'cancelled',
    '2026-02-01 10:00:00+00',
    '2026-02-08 10:00:00+00',
    'Cancelled — switched to different malt supplier for spring batches.'
);

-- 2. Line items for cancelled PO
INSERT INTO purchase_order_line (uuid, purchase_order_id, line_number, item_type, item_name, inventory_item_uuid, quantity, quantity_unit, unit_cost_cents, currency)
VALUES
    ('60000000-0000-0000-0000-000000000020',
     (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000008'),
     1, 'ingredient', 'Munich Malt', '70000000-0000-0000-0000-000000000032', 200, 'kg', 95, 'USD'),
    ('60000000-0000-0000-0000-000000000021',
     (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000008'),
     2, 'ingredient', 'Vienna Malt', '70000000-0000-0000-0000-000000000035', 100, 'kg', 100, 'USD');

-- 3. Fee for cancelled PO
INSERT INTO purchase_order_fee (uuid, purchase_order_id, fee_type, amount_cents, currency)
VALUES (
    '61000000-0000-0000-0000-000000000007',
    (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000008'),
    'shipping',
    1800,
    'USD'
);

COMMIT;
