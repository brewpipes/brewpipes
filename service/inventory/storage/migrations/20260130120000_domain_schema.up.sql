-- Inventory service schema.
-- Cross-domain references are stored as opaque UUIDs without foreign keys.
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS ingredient (
    id             serial PRIMARY KEY,
    uuid           uuid NOT NULL DEFAULT gen_random_uuid(),

    name           varchar(255) NOT NULL,
    type           varchar(32) NOT NULL,
    variety        varchar(64),
    origin_country varchar(64),
    notes          text,

    created_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at     timestamptz,
    CONSTRAINT ingredient_type_check CHECK (type IN (
        'malt',
        'hop',
        'yeast',
        'adjunct',
        'chemical',
        'gas',
        'other'
    ))
);

CREATE UNIQUE INDEX IF NOT EXISTS ingredient_uuid_idx ON ingredient(uuid);

CREATE TABLE IF NOT EXISTS ingredient_lot (
    id              serial PRIMARY KEY,
    uuid            uuid NOT NULL DEFAULT gen_random_uuid(),

    ingredient_id   int NOT NULL REFERENCES ingredient(id),
    supplier_uuid   uuid,
    lot_code        varchar(64),
    harvest_year    int,
    received_at     timestamptz,
    amount_received bigint NOT NULL,
    amount_unit     varchar(7) NOT NULL,

    created_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at      timestamptz,
    CONSTRAINT ingredient_lot_amount_check CHECK (amount_received > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS ingredient_lot_uuid_idx ON ingredient_lot(uuid);
CREATE INDEX IF NOT EXISTS ingredient_lot_ingredient_id_idx ON ingredient_lot(ingredient_id);
CREATE INDEX IF NOT EXISTS ingredient_lot_supplier_uuid_idx ON ingredient_lot(supplier_uuid);

CREATE TABLE IF NOT EXISTS ingredient_lot_transaction (
    id                serial PRIMARY KEY,
    uuid              uuid NOT NULL DEFAULT gen_random_uuid(),

    ingredient_lot_id int NOT NULL REFERENCES ingredient_lot(id),
    direction         varchar(8) NOT NULL,
    amount            bigint NOT NULL,
    amount_unit       varchar(7) NOT NULL,
    occurred_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    reference_type    varchar(32),
    reference_uuid    uuid,
    notes             text,

    created_at        timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at        timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at        timestamptz,
    CONSTRAINT ingredient_lot_transaction_direction_check CHECK (direction IN (
        'in',
        'out',
        'adjust'
    )),
    CONSTRAINT ingredient_lot_transaction_amount_check CHECK (amount > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS ingredient_lot_transaction_uuid_idx ON ingredient_lot_transaction(uuid);
CREATE INDEX IF NOT EXISTS ingredient_lot_transaction_lot_id_idx ON ingredient_lot_transaction(ingredient_lot_id);
CREATE INDEX IF NOT EXISTS ingredient_lot_transaction_reference_uuid_idx ON ingredient_lot_transaction(reference_uuid);
CREATE INDEX IF NOT EXISTS ingredient_lot_transaction_occurred_at_idx ON ingredient_lot_transaction(occurred_at);

CREATE TABLE IF NOT EXISTS ingredient_usage (
    id                       serial PRIMARY KEY,
    uuid                     uuid NOT NULL DEFAULT gen_random_uuid(),

    ingredient_lot_id         int NOT NULL REFERENCES ingredient_lot(id),
    production_addition_uuid  uuid NOT NULL,
    amount                   bigint NOT NULL,
    amount_unit              varchar(7) NOT NULL,
    used_at                  timestamptz NOT NULL DEFAULT timezone('utc', now()),

    created_at               timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at               timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at               timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS ingredient_usage_uuid_idx ON ingredient_usage(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS ingredient_usage_production_addition_uuid_idx ON ingredient_usage(production_addition_uuid);
CREATE INDEX IF NOT EXISTS ingredient_usage_lot_id_idx ON ingredient_usage(ingredient_lot_id);

CREATE TABLE IF NOT EXISTS package_sku (
    id              serial PRIMARY KEY,
    uuid            uuid NOT NULL DEFAULT gen_random_uuid(),

    name            varchar(255) NOT NULL,
    package_type    varchar(32) NOT NULL,
    volume_per_unit bigint NOT NULL,
    volume_unit     varchar(7) NOT NULL,

    created_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at      timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS package_sku_uuid_idx ON package_sku(uuid);

CREATE TABLE IF NOT EXISTS packaged_lot (
    id                    serial PRIMARY KEY,
    uuid                  uuid NOT NULL DEFAULT gen_random_uuid(),

    package_sku_id         int NOT NULL REFERENCES package_sku(id),
    production_batch_uuid  uuid NOT NULL,
    lot_code               varchar(64),
    packaged_at            timestamptz NOT NULL DEFAULT timezone('utc', now()),
    units_produced         int NOT NULL,

    created_at             timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at             timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at             timestamptz,
    CONSTRAINT packaged_lot_units_check CHECK (units_produced >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS packaged_lot_uuid_idx ON packaged_lot(uuid);
CREATE INDEX IF NOT EXISTS packaged_lot_package_sku_id_idx ON packaged_lot(package_sku_id);
CREATE INDEX IF NOT EXISTS packaged_lot_production_batch_uuid_idx ON packaged_lot(production_batch_uuid);

CREATE TABLE IF NOT EXISTS packaged_inventory_transaction (
    id               serial PRIMARY KEY,
    uuid             uuid NOT NULL DEFAULT gen_random_uuid(),

    packaged_lot_id  int NOT NULL REFERENCES packaged_lot(id),
    direction        varchar(8) NOT NULL,
    units            int NOT NULL,
    occurred_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),
    reference_type   varchar(32),
    reference_uuid   uuid,
    notes            text,

    created_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at       timestamptz,
    CONSTRAINT packaged_inventory_direction_check CHECK (direction IN (
        'in',
        'out',
        'adjust'
    )),
    CONSTRAINT packaged_inventory_units_check CHECK (units > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS packaged_inventory_transaction_uuid_idx ON packaged_inventory_transaction(uuid);
CREATE INDEX IF NOT EXISTS packaged_inventory_transaction_lot_id_idx ON packaged_inventory_transaction(packaged_lot_id);
CREATE INDEX IF NOT EXISTS packaged_inventory_transaction_reference_uuid_idx ON packaged_inventory_transaction(reference_uuid);
CREATE INDEX IF NOT EXISTS packaged_inventory_transaction_occurred_at_idx ON packaged_inventory_transaction(occurred_at);
