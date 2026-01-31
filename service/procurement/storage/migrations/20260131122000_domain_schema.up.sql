-- Procurement service schema.
-- Cross-domain references are stored as opaque UUIDs without foreign keys.
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS supplier (
    id             serial PRIMARY KEY,
    uuid           uuid NOT NULL DEFAULT gen_random_uuid(),

    name           varchar(255) NOT NULL,
    contact_name   varchar(255),
    email          varchar(255),
    phone          varchar(64),
    address_line1  varchar(255),
    address_line2  varchar(255),
    city           varchar(128),
    region         varchar(128),
    postal_code    varchar(32),
    country        varchar(64),

    created_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at     timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS supplier_uuid_idx ON supplier(uuid);

CREATE TABLE IF NOT EXISTS purchase_order (
    id           serial PRIMARY KEY,
    uuid         uuid NOT NULL DEFAULT gen_random_uuid(),

    supplier_id  int NOT NULL REFERENCES supplier(id),
    order_number varchar(64) NOT NULL,
    status       varchar(32) NOT NULL DEFAULT 'draft',
    ordered_at   timestamptz,
    expected_at  timestamptz,
    notes        text,

    created_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at   timestamptz,
    CONSTRAINT purchase_order_status_check CHECK (status IN (
        'draft',
        'submitted',
        'confirmed',
        'partially_received',
        'received',
        'cancelled'
    ))
);

CREATE UNIQUE INDEX IF NOT EXISTS purchase_order_uuid_idx ON purchase_order(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS purchase_order_order_number_idx ON purchase_order(order_number);
CREATE INDEX IF NOT EXISTS purchase_order_supplier_id_idx ON purchase_order(supplier_id);

CREATE TABLE IF NOT EXISTS purchase_order_line (
    id               serial PRIMARY KEY,
    uuid             uuid NOT NULL DEFAULT gen_random_uuid(),

    purchase_order_id int NOT NULL REFERENCES purchase_order(id),
    line_number      int NOT NULL,
    item_type        varchar(32) NOT NULL,
    item_name        varchar(255) NOT NULL,
    inventory_item_uuid uuid,
    quantity         bigint NOT NULL,
    quantity_unit    varchar(7) NOT NULL,
    unit_cost_cents  bigint NOT NULL,
    currency         char(3) NOT NULL,

    created_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at       timestamptz,
    CONSTRAINT purchase_order_line_item_type_check CHECK (item_type IN (
        'ingredient',
        'packaging',
        'service',
        'equipment',
        'other'
    )),
    CONSTRAINT purchase_order_line_cost_check CHECK (unit_cost_cents >= 0),
    CONSTRAINT purchase_order_line_quantity_check CHECK (quantity > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS purchase_order_line_uuid_idx ON purchase_order_line(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS purchase_order_line_unique_idx ON purchase_order_line(purchase_order_id, line_number);
CREATE INDEX IF NOT EXISTS purchase_order_line_inventory_item_uuid_idx ON purchase_order_line(inventory_item_uuid);

CREATE TABLE IF NOT EXISTS purchase_order_fee (
    id               serial PRIMARY KEY,
    uuid             uuid NOT NULL DEFAULT gen_random_uuid(),

    purchase_order_id int NOT NULL REFERENCES purchase_order(id),
    fee_type         varchar(32) NOT NULL,
    amount_cents     bigint NOT NULL,
    currency         char(3) NOT NULL,

    created_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at       timestamptz,
    CONSTRAINT purchase_order_fee_amount_check CHECK (amount_cents >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS purchase_order_fee_uuid_idx ON purchase_order_fee(uuid);
CREATE INDEX IF NOT EXISTS purchase_order_fee_purchase_order_id_idx ON purchase_order_fee(purchase_order_id);

-- Seed data for early development.
INSERT INTO supplier (uuid, name, contact_name, email, phone, address_line1, city, region, postal_code, country)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'Great Lakes Malt Co', 'Erin Foster', 'orders@greatlakesmalt.co', '+1-312-555-0144', '1200 Grain Ave', 'Chicago', 'IL', '60601', 'USA'),
    ('22222222-2222-2222-2222-222222222222', 'Northwest Hop Farms', 'Leo Vargas', 'sales@nwhops.example', '+1-503-555-0199', '88 Hop Yard Rd', 'Yakima', 'WA', '98901', 'USA'),
    ('33333333-3333-3333-3333-333333333333', 'Coastal Yeast Labs', 'Priya Shah', 'support@coastalyeast.example', '+1-619-555-0118', '405 Culture Way', 'San Diego', 'CA', '92101', 'USA'),
    ('44444444-4444-4444-4444-444444444444', 'Carbonic Gases LLC', 'Maya Reed', 'logistics@carbonicgases.example', '+1-713-555-0102', '900 Tanker Ln', 'Houston', 'TX', '77002', 'USA');

INSERT INTO purchase_order (uuid, supplier_id, order_number, status, ordered_at, expected_at, notes)
VALUES
    ('50000000-0000-0000-0000-000000000001', (SELECT id FROM supplier WHERE uuid = '11111111-1111-1111-1111-111111111111'), 'PO-1001', 'confirmed', '2026-01-02 10:00:00+00', '2026-01-09 10:00:00+00', 'Base and specialty malts for Q1 brews.'),
    ('50000000-0000-0000-0000-000000000002', (SELECT id FROM supplier WHERE uuid = '22222222-2222-2222-2222-222222222222'), 'PO-1002', 'received', '2026-01-03 11:00:00+00', '2026-01-10 11:00:00+00', 'Hop pellets for IPA and pale ale.'),
    ('50000000-0000-0000-0000-000000000003', (SELECT id FROM supplier WHERE uuid = '33333333-3333-3333-3333-333333333333'), 'PO-1003', 'received', '2026-01-04 12:00:00+00', '2026-01-11 12:00:00+00', 'House yeast for flagship fermentations.'),
    ('50000000-0000-0000-0000-000000000004', (SELECT id FROM supplier WHERE uuid = '44444444-4444-4444-4444-444444444444'), 'PO-1004', 'submitted', '2026-01-20 09:00:00+00', '2026-01-27 09:00:00+00', 'Bulk CO2 delivery for packaging line.');

INSERT INTO purchase_order_line (
    uuid,
    purchase_order_id,
    line_number,
    item_type,
    item_name,
    inventory_item_uuid,
    quantity,
    quantity_unit,
    unit_cost_cents,
    currency
)
VALUES
    ('60000000-0000-0000-0000-000000000001', (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000001'), 1, 'ingredient', 'Pale Malt', '70000000-0000-0000-0000-000000000001', 1000, 'kg', 85, 'USD'),
    ('60000000-0000-0000-0000-000000000002', (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000001'), 2, 'ingredient', 'Pilsner Malt', '70000000-0000-0000-0000-000000000002', 800, 'kg', 90, 'USD'),
    ('60000000-0000-0000-0000-000000000003', (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000001'), 3, 'ingredient', 'Crystal 60L', '70000000-0000-0000-0000-000000000003', 200, 'kg', 120, 'USD'),
    ('60000000-0000-0000-0000-000000000004', (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000002'), 1, 'ingredient', 'Cascade Hops', '70000000-0000-0000-0000-000000000004', 50, 'kg', 1600, 'USD'),
    ('60000000-0000-0000-0000-000000000005', (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000002'), 2, 'ingredient', 'Citra Hops', '70000000-0000-0000-0000-000000000005', 40, 'kg', 2000, 'USD'),
    ('60000000-0000-0000-0000-000000000006', (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000003'), 1, 'ingredient', 'WLP001 California Ale Yeast', '70000000-0000-0000-0000-000000000006', 10, 'kg', 12000, 'USD'),
    ('60000000-0000-0000-0000-000000000007', (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000004'), 1, 'ingredient', 'CO2 Gas', '70000000-0000-0000-0000-000000000009', 500, 'lb', 35, 'USD');

INSERT INTO purchase_order_fee (uuid, purchase_order_id, fee_type, amount_cents, currency)
VALUES
    ('61000000-0000-0000-0000-000000000001', (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000001'), 'shipping', 2500, 'USD'),
    ('61000000-0000-0000-0000-000000000002', (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000002'), 'freight', 1800, 'USD'),
    ('61000000-0000-0000-0000-000000000003', (SELECT id FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000004'), 'hazmat', 500, 'USD');
