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

    supplier_uuid uuid NOT NULL,
    order_number  varchar(64) NOT NULL,
    status        varchar(32) NOT NULL DEFAULT 'draft',
    ordered_at    timestamptz,
    expected_at   timestamptz,
    notes         text,

    created_at    timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at    timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at    timestamptz,
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
CREATE INDEX IF NOT EXISTS purchase_order_supplier_uuid_idx ON purchase_order(supplier_uuid);

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
