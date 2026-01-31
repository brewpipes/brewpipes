-- Inventory service schema.
-- Cross-domain references are stored as opaque UUIDs without foreign keys.
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS ingredient (
    id           serial PRIMARY KEY,
    uuid         uuid NOT NULL DEFAULT gen_random_uuid(),

    name         varchar(255) NOT NULL,
    category     varchar(32) NOT NULL,
    default_unit varchar(7) NOT NULL,
    description  text,

    created_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at   timestamptz,
    CONSTRAINT ingredient_category_check CHECK (category IN (
        'fermentable',
        'hop',
        'yeast',
        'adjunct',
        'salt',
        'chemical',
        'gas',
        'other'
    ))
);

CREATE UNIQUE INDEX IF NOT EXISTS ingredient_uuid_idx ON ingredient(uuid);
CREATE INDEX IF NOT EXISTS ingredient_category_idx ON ingredient(category);

CREATE TABLE IF NOT EXISTS ingredient_malt_detail (
    id               serial PRIMARY KEY,
    uuid             uuid NOT NULL DEFAULT gen_random_uuid(),

    ingredient_id    int NOT NULL REFERENCES ingredient(id),
    maltster_name    varchar(255),
    variety          varchar(255),
    lovibond         numeric(6,2),
    srm              numeric(6,2),
    diastatic_power  numeric(6,2),

    created_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at       timestamptz,
    CONSTRAINT ingredient_malt_detail_lovibond_check CHECK (lovibond IS NULL OR lovibond >= 0),
    CONSTRAINT ingredient_malt_detail_srm_check CHECK (srm IS NULL OR srm >= 0),
    CONSTRAINT ingredient_malt_detail_diastatic_power_check CHECK (diastatic_power IS NULL OR diastatic_power >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS ingredient_malt_detail_uuid_idx ON ingredient_malt_detail(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS ingredient_malt_detail_ingredient_id_idx ON ingredient_malt_detail(ingredient_id);

CREATE TABLE IF NOT EXISTS ingredient_hop_detail (
    id               serial PRIMARY KEY,
    uuid             uuid NOT NULL DEFAULT gen_random_uuid(),

    ingredient_id    int NOT NULL REFERENCES ingredient(id),
    producer_name    varchar(255),
    variety          varchar(255),
    crop_year        int,
    form             varchar(16),
    alpha_acid       numeric(6,2),
    beta_acid        numeric(6,2),

    created_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at       timestamptz,
    CONSTRAINT ingredient_hop_detail_crop_year_check CHECK (crop_year IS NULL OR crop_year >= 1900),
    CONSTRAINT ingredient_hop_detail_form_check CHECK (form IS NULL OR form IN (
        'pellet',
        'whole_leaf',
        'cryo',
        'extract',
        'other'
    )),
    CONSTRAINT ingredient_hop_detail_alpha_acid_check CHECK (alpha_acid IS NULL OR (alpha_acid >= 0 AND alpha_acid <= 100)),
    CONSTRAINT ingredient_hop_detail_beta_acid_check CHECK (beta_acid IS NULL OR (beta_acid >= 0 AND beta_acid <= 100))
);

CREATE UNIQUE INDEX IF NOT EXISTS ingredient_hop_detail_uuid_idx ON ingredient_hop_detail(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS ingredient_hop_detail_ingredient_id_idx ON ingredient_hop_detail(ingredient_id);

CREATE TABLE IF NOT EXISTS ingredient_yeast_detail (
    id               serial PRIMARY KEY,
    uuid             uuid NOT NULL DEFAULT gen_random_uuid(),

    ingredient_id    int NOT NULL REFERENCES ingredient(id),
    lab_name         varchar(255),
    strain           varchar(255),
    form             varchar(16),

    created_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at       timestamptz,
    CONSTRAINT ingredient_yeast_detail_form_check CHECK (form IS NULL OR form IN (
        'liquid',
        'dry',
        'slurry',
        'propagated',
        'other'
    ))
);

CREATE UNIQUE INDEX IF NOT EXISTS ingredient_yeast_detail_uuid_idx ON ingredient_yeast_detail(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS ingredient_yeast_detail_ingredient_id_idx ON ingredient_yeast_detail(ingredient_id);

CREATE TABLE IF NOT EXISTS stock_location (
    id             serial PRIMARY KEY,
    uuid           uuid NOT NULL DEFAULT gen_random_uuid(),

    name           varchar(255) NOT NULL,
    location_type  varchar(32),
    description    text,

    created_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at     timestamptz,
    CONSTRAINT stock_location_type_check CHECK (location_type IS NULL OR location_type IN (
        'dry',
        'cold',
        'gas',
        'bulk',
        'packaging',
        'other'
    ))
);

CREATE UNIQUE INDEX IF NOT EXISTS stock_location_uuid_idx ON stock_location(uuid);

CREATE TABLE IF NOT EXISTS inventory_receipt (
    id             serial PRIMARY KEY,
    uuid           uuid NOT NULL DEFAULT gen_random_uuid(),

    supplier_uuid  uuid,
    reference_code varchar(64),
    received_at    timestamptz NOT NULL DEFAULT timezone('utc', now()),
    notes          text,

    created_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at     timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS inventory_receipt_uuid_idx ON inventory_receipt(uuid);
CREATE INDEX IF NOT EXISTS inventory_receipt_supplier_uuid_idx ON inventory_receipt(supplier_uuid);

CREATE TABLE IF NOT EXISTS ingredient_lot (
    id                       serial PRIMARY KEY,
    uuid                     uuid NOT NULL DEFAULT gen_random_uuid(),

    ingredient_id            int NOT NULL REFERENCES ingredient(id),
    receipt_id               int REFERENCES inventory_receipt(id),
    supplier_uuid            uuid,
    brewery_lot_code          varchar(64),
    originator_lot_code      varchar(64),
    originator_name          varchar(255),
    originator_type          varchar(32),
    received_at              timestamptz NOT NULL DEFAULT timezone('utc', now()),
    received_amount          bigint NOT NULL,
    received_unit            varchar(7) NOT NULL,
    best_by_at               timestamptz,
    expires_at               timestamptz,
    notes                    text,

    created_at               timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at               timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at               timestamptz,
    CONSTRAINT ingredient_lot_received_amount_check CHECK (received_amount > 0),
    CONSTRAINT ingredient_lot_originator_type_check CHECK (originator_type IS NULL OR originator_type IN (
        'maltster',
        'hop_producer',
        'yeast_lab',
        'gas_vendor',
        'other'
    )),
    CONSTRAINT ingredient_lot_expires_after_best_by_check CHECK (expires_at IS NULL OR best_by_at IS NULL OR expires_at >= best_by_at)
);

CREATE UNIQUE INDEX IF NOT EXISTS ingredient_lot_uuid_idx ON ingredient_lot(uuid);
CREATE INDEX IF NOT EXISTS ingredient_lot_ingredient_id_idx ON ingredient_lot(ingredient_id);
CREATE INDEX IF NOT EXISTS ingredient_lot_receipt_id_idx ON ingredient_lot(receipt_id);
CREATE INDEX IF NOT EXISTS ingredient_lot_supplier_uuid_idx ON ingredient_lot(supplier_uuid);
CREATE INDEX IF NOT EXISTS ingredient_lot_brewery_lot_code_idx ON ingredient_lot(brewery_lot_code);
CREATE INDEX IF NOT EXISTS ingredient_lot_originator_lot_code_idx ON ingredient_lot(originator_lot_code);

CREATE TABLE IF NOT EXISTS ingredient_lot_malt_detail (
    id                 serial PRIMARY KEY,
    uuid               uuid NOT NULL DEFAULT gen_random_uuid(),

    ingredient_lot_id  int NOT NULL REFERENCES ingredient_lot(id),
    moisture_percent   numeric(5,2),

    created_at         timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at         timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at         timestamptz,
    CONSTRAINT ingredient_lot_malt_detail_moisture_check CHECK (moisture_percent IS NULL OR (moisture_percent >= 0 AND moisture_percent <= 100))
);

CREATE UNIQUE INDEX IF NOT EXISTS ingredient_lot_malt_detail_uuid_idx ON ingredient_lot_malt_detail(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS ingredient_lot_malt_detail_lot_id_idx ON ingredient_lot_malt_detail(ingredient_lot_id);

CREATE TABLE IF NOT EXISTS ingredient_lot_hop_detail (
    id                 serial PRIMARY KEY,
    uuid               uuid NOT NULL DEFAULT gen_random_uuid(),

    ingredient_lot_id  int NOT NULL REFERENCES ingredient_lot(id),
    alpha_acid         numeric(6,2),
    beta_acid          numeric(6,2),

    created_at         timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at         timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at         timestamptz,
    CONSTRAINT ingredient_lot_hop_detail_alpha_acid_check CHECK (alpha_acid IS NULL OR (alpha_acid >= 0 AND alpha_acid <= 100)),
    CONSTRAINT ingredient_lot_hop_detail_beta_acid_check CHECK (beta_acid IS NULL OR (beta_acid >= 0 AND beta_acid <= 100))
);

CREATE UNIQUE INDEX IF NOT EXISTS ingredient_lot_hop_detail_uuid_idx ON ingredient_lot_hop_detail(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS ingredient_lot_hop_detail_lot_id_idx ON ingredient_lot_hop_detail(ingredient_lot_id);

CREATE TABLE IF NOT EXISTS ingredient_lot_yeast_detail (
    id                      serial PRIMARY KEY,
    uuid                    uuid NOT NULL DEFAULT gen_random_uuid(),

    ingredient_lot_id       int NOT NULL REFERENCES ingredient_lot(id),
    viability_percent       numeric(5,2),
    generation              int,

    created_at              timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at              timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at              timestamptz,
    CONSTRAINT ingredient_lot_yeast_detail_viability_check CHECK (viability_percent IS NULL OR (viability_percent >= 0 AND viability_percent <= 100)),
    CONSTRAINT ingredient_lot_yeast_detail_generation_check CHECK (generation IS NULL OR generation >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS ingredient_lot_yeast_detail_uuid_idx ON ingredient_lot_yeast_detail(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS ingredient_lot_yeast_detail_lot_id_idx ON ingredient_lot_yeast_detail(ingredient_lot_id);

CREATE TABLE IF NOT EXISTS inventory_usage (
    id                   serial PRIMARY KEY,
    uuid                 uuid NOT NULL DEFAULT gen_random_uuid(),

    production_ref_uuid  uuid,
    used_at              timestamptz NOT NULL DEFAULT timezone('utc', now()),
    notes                text,

    created_at           timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at           timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at           timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS inventory_usage_uuid_idx ON inventory_usage(uuid);
CREATE INDEX IF NOT EXISTS inventory_usage_production_ref_uuid_idx ON inventory_usage(production_ref_uuid);

CREATE TABLE IF NOT EXISTS inventory_adjustment (
    id             serial PRIMARY KEY,
    uuid           uuid NOT NULL DEFAULT gen_random_uuid(),

    reason         varchar(32) NOT NULL,
    adjusted_at    timestamptz NOT NULL DEFAULT timezone('utc', now()),
    notes          text,

    created_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at     timestamptz,
    CONSTRAINT inventory_adjustment_reason_check CHECK (reason IN (
        'cycle_count',
        'spoilage',
        'shrink',
        'damage',
        'correction',
        'other'
    ))
);

CREATE UNIQUE INDEX IF NOT EXISTS inventory_adjustment_uuid_idx ON inventory_adjustment(uuid);

CREATE TABLE IF NOT EXISTS inventory_transfer (
    id                 serial PRIMARY KEY,
    uuid               uuid NOT NULL DEFAULT gen_random_uuid(),

    source_location_id int NOT NULL REFERENCES stock_location(id),
    dest_location_id   int NOT NULL REFERENCES stock_location(id),
    transferred_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    notes              text,

    created_at         timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at         timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at         timestamptz,
    CONSTRAINT inventory_transfer_source_dest_check CHECK (source_location_id <> dest_location_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS inventory_transfer_uuid_idx ON inventory_transfer(uuid);
CREATE INDEX IF NOT EXISTS inventory_transfer_source_location_id_idx ON inventory_transfer(source_location_id);
CREATE INDEX IF NOT EXISTS inventory_transfer_dest_location_id_idx ON inventory_transfer(dest_location_id);

CREATE TABLE IF NOT EXISTS beer_lot (
    id                    serial PRIMARY KEY,
    uuid                  uuid NOT NULL DEFAULT gen_random_uuid(),

    production_batch_uuid uuid NOT NULL,
    lot_code              varchar(64),
    packaged_at           timestamptz NOT NULL DEFAULT timezone('utc', now()),
    notes                 text,

    created_at            timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at            timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at            timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS beer_lot_uuid_idx ON beer_lot(uuid);
CREATE INDEX IF NOT EXISTS beer_lot_production_batch_uuid_idx ON beer_lot(production_batch_uuid);

CREATE TABLE IF NOT EXISTS inventory_movement (
    id                serial PRIMARY KEY,
    uuid              uuid NOT NULL DEFAULT gen_random_uuid(),

    ingredient_lot_id int REFERENCES ingredient_lot(id),
    beer_lot_id       int REFERENCES beer_lot(id),
    stock_location_id int NOT NULL REFERENCES stock_location(id),
    direction         varchar(8) NOT NULL,
    reason            varchar(16) NOT NULL,
    amount            bigint NOT NULL,
    amount_unit       varchar(7) NOT NULL,
    occurred_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),
    receipt_id        int REFERENCES inventory_receipt(id),
    usage_id          int REFERENCES inventory_usage(id),
    adjustment_id     int REFERENCES inventory_adjustment(id),
    transfer_id       int REFERENCES inventory_transfer(id),
    notes             text,

    created_at        timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at        timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at        timestamptz,
    CONSTRAINT inventory_movement_target_check CHECK (
        (ingredient_lot_id IS NOT NULL AND beer_lot_id IS NULL) OR
        (ingredient_lot_id IS NULL AND beer_lot_id IS NOT NULL)
    ),
    CONSTRAINT inventory_movement_direction_check CHECK (direction IN (
        'in',
        'out'
    )),
    CONSTRAINT inventory_movement_reason_check CHECK (reason IN (
        'receive',
        'use',
        'transfer',
        'adjust',
        'waste'
    )),
    CONSTRAINT inventory_movement_amount_check CHECK (amount > 0),
    CONSTRAINT inventory_movement_reason_reference_check CHECK (
        (reason = 'receive' AND receipt_id IS NOT NULL) OR
        (reason = 'use' AND usage_id IS NOT NULL) OR
        (reason = 'transfer' AND transfer_id IS NOT NULL) OR
        (reason IN ('adjust', 'waste') AND adjustment_id IS NOT NULL)
    ),
    CONSTRAINT inventory_movement_reference_check CHECK (num_nonnulls(receipt_id, usage_id, adjustment_id, transfer_id) <= 1)
);

CREATE UNIQUE INDEX IF NOT EXISTS inventory_movement_uuid_idx ON inventory_movement(uuid);
CREATE INDEX IF NOT EXISTS inventory_movement_ingredient_lot_id_idx ON inventory_movement(ingredient_lot_id);
CREATE INDEX IF NOT EXISTS inventory_movement_beer_lot_id_idx ON inventory_movement(beer_lot_id);
CREATE INDEX IF NOT EXISTS inventory_movement_stock_location_id_idx ON inventory_movement(stock_location_id);
CREATE INDEX IF NOT EXISTS inventory_movement_receipt_id_idx ON inventory_movement(receipt_id);
CREATE INDEX IF NOT EXISTS inventory_movement_usage_id_idx ON inventory_movement(usage_id);
CREATE INDEX IF NOT EXISTS inventory_movement_adjustment_id_idx ON inventory_movement(adjustment_id);
CREATE INDEX IF NOT EXISTS inventory_movement_transfer_id_idx ON inventory_movement(transfer_id);
CREATE INDEX IF NOT EXISTS inventory_movement_occurred_at_idx ON inventory_movement(occurred_at);
