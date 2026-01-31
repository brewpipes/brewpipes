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
    purchase_order_uuid uuid,
    reference_code varchar(64),
    received_at    timestamptz NOT NULL DEFAULT timezone('utc', now()),
    notes          text,

    created_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at     timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS inventory_receipt_uuid_idx ON inventory_receipt(uuid);
CREATE INDEX IF NOT EXISTS inventory_receipt_supplier_uuid_idx ON inventory_receipt(supplier_uuid);
CREATE INDEX IF NOT EXISTS inventory_receipt_purchase_order_uuid_idx ON inventory_receipt(purchase_order_uuid);

CREATE TABLE IF NOT EXISTS ingredient_lot (
    id                       serial PRIMARY KEY,
    uuid                     uuid NOT NULL DEFAULT gen_random_uuid(),

    ingredient_id            int NOT NULL REFERENCES ingredient(id),
    receipt_id               int REFERENCES inventory_receipt(id),
    supplier_uuid            uuid,
    purchase_order_line_uuid uuid,
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
CREATE INDEX IF NOT EXISTS ingredient_lot_purchase_order_line_uuid_idx ON ingredient_lot(purchase_order_line_uuid);
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

-- Seed data for early development.
INSERT INTO ingredient (uuid, name, category, default_unit, description)
VALUES
    ('70000000-0000-0000-0000-000000000001', 'Pale Malt', 'fermentable', 'kg', 'Base malt for ales.'),
    ('70000000-0000-0000-0000-000000000002', 'Pilsner Malt', 'fermentable', 'kg', 'Clean, light base malt.'),
    ('70000000-0000-0000-0000-000000000003', 'Crystal 60L', 'fermentable', 'kg', 'Medium crystal malt for color and caramel notes.'),
    ('70000000-0000-0000-0000-000000000004', 'Cascade Hops', 'hop', 'kg', 'Citrus-forward aroma hops.'),
    ('70000000-0000-0000-0000-000000000005', 'Citra Hops', 'hop', 'kg', 'Tropical fruit and citrus aroma hops.'),
    ('70000000-0000-0000-0000-000000000006', 'WLP001 California Ale Yeast', 'yeast', 'kg', 'Neutral American ale yeast.'),
    ('70000000-0000-0000-0000-000000000007', 'Lactic Acid', 'chemical', 'l', 'Mash pH adjustment.'),
    ('70000000-0000-0000-0000-000000000008', 'Gypsum', 'salt', 'g', 'Calcium sulfate for water chemistry.'),
    ('70000000-0000-0000-0000-000000000009', 'CO2', 'gas', 'lb', 'Packaging and carbonation gas.'),
    ('70000000-0000-0000-0000-000000000010', 'Dextrose', 'adjunct', 'kg', 'Fermentable sugar adjunct.');

INSERT INTO ingredient (uuid, name, category, default_unit, description)
VALUES
    ('70000000-0000-0000-0000-000000000011', 'Roasted Barley', 'fermentable', 'kg', 'Deep roast for stout color.'),
    ('70000000-0000-0000-0000-000000000012', 'Saaz Hops', 'hop', 'kg', 'Traditional noble hop.'),
    ('70000000-0000-0000-0000-000000000013', 'W34/70 Lager Yeast', 'yeast', 'kg', 'Clean lager strain.'),
    ('70000000-0000-0000-0000-000000000014', 'Irish Moss', 'adjunct', 'kg', 'Kettle fining for clarity.'),
    ('70000000-0000-0000-0000-000000000015', 'Oxygen', 'gas', 'lb', 'Wort oxygenation gas.'),
    ('70000000-0000-0000-0000-000000000016', 'Coriander', 'adjunct', 'kg', 'Spice for farmhouse saisons.');

INSERT INTO ingredient_malt_detail (ingredient_id, maltster_name, variety, lovibond, srm, diastatic_power)
VALUES
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000001'), 'Great Lakes Malt Co', '2-row', 2.0, 2.0, 120.0),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000002'), 'Great Lakes Malt Co', 'Pilsner', 1.6, 1.6, 110.0),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000003'), 'Great Lakes Malt Co', 'Crystal', 60.0, 60.0, 0.0);

INSERT INTO ingredient_malt_detail (ingredient_id, maltster_name, variety, lovibond, srm, diastatic_power)
VALUES
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000011'), 'Midwest Maltings', 'Roasted', 300.0, 300.0, 0.0);

INSERT INTO ingredient_hop_detail (ingredient_id, producer_name, variety, crop_year, form, alpha_acid, beta_acid)
VALUES
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000004'), 'Northwest Hop Farms', 'Cascade', 2025, 'pellet', 6.5, 4.8),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000005'), 'Northwest Hop Farms', 'Citra', 2025, 'pellet', 12.0, 4.0);

INSERT INTO ingredient_hop_detail (ingredient_id, producer_name, variety, crop_year, form, alpha_acid, beta_acid)
VALUES
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000012'), 'Bohemian Hop Co', 'Saaz', 2025, 'pellet', 3.5, 4.2);

INSERT INTO ingredient_yeast_detail (ingredient_id, lab_name, strain, form)
VALUES
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000006'), 'White Labs', 'WLP001', 'liquid');

INSERT INTO ingredient_yeast_detail (ingredient_id, lab_name, strain, form)
VALUES
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000013'), 'Fermentis', 'W34/70', 'dry');

INSERT INTO stock_location (uuid, name, location_type, description)
VALUES
    ('81000000-0000-0000-0000-000000000001', 'Grain Room', 'dry', 'Palletized malt storage.'),
    ('81000000-0000-0000-0000-000000000002', 'Cold Room', 'cold', 'Temperature-controlled hop and yeast storage.'),
    ('81000000-0000-0000-0000-000000000003', 'Chemical Cage', 'other', 'Secure chemical storage.'),
    ('81000000-0000-0000-0000-000000000004', 'Gas Pad', 'gas', 'Bulk and cylinder gas storage.'),
    ('81000000-0000-0000-0000-000000000005', 'Brew Deck', 'other', 'Staging area for brew day picks.');

INSERT INTO stock_location (uuid, name, location_type, description)
VALUES
    ('81000000-0000-0000-0000-000000000006', 'Packaging Warehouse', 'packaging', 'Finished goods and packaging materials.'),
    ('81000000-0000-0000-0000-000000000007', 'Brite Cellar', 'cold', 'Cold conditioning storage.');

INSERT INTO inventory_receipt (uuid, supplier_uuid, purchase_order_uuid, reference_code, received_at, notes)
VALUES
    ('82000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', '50000000-0000-0000-0000-000000000001', 'RCV-1001', '2026-01-09 14:00:00+00', 'Malt delivery for IPA and lager.'),
    ('82000000-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', '50000000-0000-0000-0000-000000000002', 'RCV-1002', '2026-01-10 09:30:00+00', 'Hop pellets delivered.'),
    ('82000000-0000-0000-0000-000000000003', '33333333-3333-3333-3333-333333333333', '50000000-0000-0000-0000-000000000003', 'RCV-1003', '2026-01-11 08:30:00+00', 'House yeast received.'),
    ('82000000-0000-0000-0000-000000000004', '44444444-4444-4444-4444-444444444444', '50000000-0000-0000-0000-000000000004', 'RCV-1004', '2026-01-22 15:00:00+00', 'CO2 bulk delivery.');

INSERT INTO inventory_receipt (uuid, supplier_uuid, purchase_order_uuid, reference_code, received_at, notes)
VALUES
    ('82000000-0000-0000-0000-000000000005', '66666666-6666-6666-6666-666666666666', '50000000-0000-0000-0000-000000000005', 'RCV-1005', '2026-01-18 13:00:00+00', 'Specialty malt and Saaz hops received.'),
    ('82000000-0000-0000-0000-000000000006', '77777777-7777-7777-7777-777777777777', '50000000-0000-0000-0000-000000000006', 'RCV-1006', '2026-01-22 10:00:00+00', 'Lager yeast shipment.'),
    ('82000000-0000-0000-0000-000000000007', '66666666-6666-6666-6666-666666666666', '50000000-0000-0000-0000-000000000005', 'RCV-1007', '2026-01-18 13:10:00+00', 'Adjuncts for fining and spice.');

INSERT INTO ingredient_lot (
    uuid,
    ingredient_id,
    receipt_id,
    supplier_uuid,
    purchase_order_line_uuid,
    brewery_lot_code,
    originator_lot_code,
    originator_name,
    originator_type,
    received_at,
    received_amount,
    received_unit,
    best_by_at,
    notes
)
VALUES
    ('80000000-0000-0000-0000-000000000001', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000001'), (SELECT id FROM inventory_receipt WHERE reference_code = 'RCV-1001'), '11111111-1111-1111-1111-111111111111', '60000000-0000-0000-0000-000000000001', 'MALT-PALE-2401', 'GLM-PALE-2401', 'Great Lakes Malt Co', 'maltster', '2026-01-09 14:00:00+00', 1000, 'kg', '2026-07-01 00:00:00+00', 'Base malt lot for Q1.'),
    ('80000000-0000-0000-0000-000000000002', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000002'), (SELECT id FROM inventory_receipt WHERE reference_code = 'RCV-1001'), '11111111-1111-1111-1111-111111111111', '60000000-0000-0000-0000-000000000002', 'MALT-PILS-2401', 'GLM-PILS-2401', 'Great Lakes Malt Co', 'maltster', '2026-01-09 14:00:00+00', 800, 'kg', '2026-07-01 00:00:00+00', 'Pilsner base malt.'),
    ('80000000-0000-0000-0000-000000000003', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000003'), (SELECT id FROM inventory_receipt WHERE reference_code = 'RCV-1001'), '11111111-1111-1111-1111-111111111111', '60000000-0000-0000-0000-000000000003', 'MALT-CR60-2401', 'GLM-CR60-2401', 'Great Lakes Malt Co', 'maltster', '2026-01-09 14:00:00+00', 200, 'kg', '2026-07-01 00:00:00+00', 'Crystal 60L malt.'),
    ('80000000-0000-0000-0000-000000000004', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000004'), (SELECT id FROM inventory_receipt WHERE reference_code = 'RCV-1002'), '22222222-2222-2222-2222-222222222222', '60000000-0000-0000-0000-000000000004', 'HOP-CAS-2501', 'NHF-CAS-25A', 'Northwest Hop Farms', 'hop_producer', '2026-01-10 09:30:00+00', 50, 'kg', '2026-12-31 00:00:00+00', 'Cascade crop 2025.'),
    ('80000000-0000-0000-0000-000000000005', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000005'), (SELECT id FROM inventory_receipt WHERE reference_code = 'RCV-1002'), '22222222-2222-2222-2222-222222222222', '60000000-0000-0000-0000-000000000005', 'HOP-CIT-2501', 'NHF-CIT-25B', 'Northwest Hop Farms', 'hop_producer', '2026-01-10 09:30:00+00', 40, 'kg', '2026-12-31 00:00:00+00', 'Citra crop 2025.'),
    ('80000000-0000-0000-0000-000000000006', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000006'), (SELECT id FROM inventory_receipt WHERE reference_code = 'RCV-1003'), '33333333-3333-3333-3333-333333333333', '60000000-0000-0000-0000-000000000006', 'YEAST-WLP001-2401', 'CYL-WLP001-2401', 'Coastal Yeast Labs', 'yeast_lab', '2026-01-11 08:30:00+00', 10, 'kg', '2026-03-01 00:00:00+00', 'House ale yeast.'),
    ('80000000-0000-0000-0000-000000000007', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000009'), (SELECT id FROM inventory_receipt WHERE reference_code = 'RCV-1004'), '44444444-4444-4444-4444-444444444444', '60000000-0000-0000-0000-000000000007', 'GAS-CO2-2401', 'CG-CO2-2401', 'Carbonic Gases LLC', 'gas_vendor', '2026-01-22 15:00:00+00', 500, 'lb', NULL, 'Bulk CO2 delivery.');

INSERT INTO ingredient_lot (
    uuid,
    ingredient_id,
    receipt_id,
    supplier_uuid,
    purchase_order_line_uuid,
    brewery_lot_code,
    originator_lot_code,
    originator_name,
    originator_type,
    received_at,
    received_amount,
    received_unit,
    best_by_at,
    notes
)
VALUES
    ('80000000-0000-0000-0000-000000000008', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000011'), (SELECT id FROM inventory_receipt WHERE reference_code = 'RCV-1005'), '66666666-6666-6666-6666-666666666666', '60000000-0000-0000-0000-000000000008', 'MALT-ROAST-2401', 'MM-ROAST-2401', 'Midwest Maltings', 'maltster', '2026-01-18 13:00:00+00', 150, 'kg', '2026-07-15 00:00:00+00', 'Roasted barley for stouts.'),
    ('80000000-0000-0000-0000-000000000009', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000012'), (SELECT id FROM inventory_receipt WHERE reference_code = 'RCV-1005'), '66666666-6666-6666-6666-666666666666', '60000000-0000-0000-0000-000000000009', 'HOP-SAA-2501', 'BHC-SAA-25A', 'Bohemian Hop Co', 'hop_producer', '2026-01-18 13:00:00+00', 30, 'kg', '2026-12-31 00:00:00+00', 'Saaz crop 2025.'),
    ('80000000-0000-0000-0000-000000000010', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000013'), (SELECT id FROM inventory_receipt WHERE reference_code = 'RCV-1006'), '77777777-7777-7777-7777-777777777777', '60000000-0000-0000-0000-000000000010', 'YEAST-W3470-2401', 'FER-W3470-2401', 'Fermentis', 'yeast_lab', '2026-01-22 10:00:00+00', 8, 'kg', '2026-04-01 00:00:00+00', 'Lager yeast shipment.'),
    ('80000000-0000-0000-0000-000000000011', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000014'), (SELECT id FROM inventory_receipt WHERE reference_code = 'RCV-1007'), '66666666-6666-6666-6666-666666666666', '60000000-0000-0000-0000-000000000011', 'ADJ-IRISH-2401', 'IRISH-2401', 'Brew Supply Co', 'other', '2026-01-18 13:10:00+00', 20, 'kg', '2027-01-01 00:00:00+00', 'Irish moss for kettle finings.'),
    ('80000000-0000-0000-0000-000000000012', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000016'), (SELECT id FROM inventory_receipt WHERE reference_code = 'RCV-1007'), '66666666-6666-6666-6666-666666666666', '60000000-0000-0000-0000-000000000012', 'ADJ-CORI-2401', 'CORI-2401', 'Spice River Co', 'other', '2026-01-18 13:10:00+00', 10, 'kg', '2027-01-01 00:00:00+00', 'Coriander for saison spice.');

INSERT INTO ingredient_lot_malt_detail (ingredient_lot_id, moisture_percent)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000001'), 4.2),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000002'), 4.0),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000003'), 3.8);

INSERT INTO ingredient_lot_malt_detail (ingredient_lot_id, moisture_percent)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000008'), 4.5);

INSERT INTO ingredient_lot_hop_detail (ingredient_lot_id, alpha_acid, beta_acid)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000004'), 6.8, 4.9),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000005'), 12.3, 4.1);

INSERT INTO ingredient_lot_hop_detail (ingredient_lot_id, alpha_acid, beta_acid)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000009'), 3.6, 4.0);

INSERT INTO ingredient_lot_yeast_detail (ingredient_lot_id, viability_percent, generation)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000006'), 95.0, 0);

INSERT INTO ingredient_lot_yeast_detail (ingredient_lot_id, viability_percent, generation)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000010'), 94.0, 0);

INSERT INTO inventory_usage (uuid, production_ref_uuid, used_at, notes)
VALUES
    ('83000000-0000-0000-0000-000000000001', '90000000-0000-0000-0000-000000000001', '2026-01-10 08:15:00+00', 'Brew day usage for IPA 24-07.');

INSERT INTO inventory_usage (uuid, production_ref_uuid, used_at, notes)
VALUES
    ('83000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000004', '2026-01-22 08:25:00+00', 'Brew day usage for Stout 24-09.'),
    ('83000000-0000-0000-0000-000000000003', '90000000-0000-0000-0000-000000000005', '2026-01-24 08:20:00+00', 'Brew day usage for Kolsch 24-10.'),
    ('83000000-0000-0000-0000-000000000004', '90000000-0000-0000-0000-000000000008', '2026-01-26 08:10:00+00', 'Brew day usage for Saison 24-11.');

INSERT INTO inventory_adjustment (uuid, reason, adjusted_at, notes)
VALUES
    ('84000000-0000-0000-0000-000000000001', 'cycle_count', '2026-01-15 07:00:00+00', 'Monthly inventory count adjustment.');

INSERT INTO inventory_adjustment (uuid, reason, adjusted_at, notes)
VALUES
    ('84000000-0000-0000-0000-000000000002', 'spoilage', '2026-01-23 07:15:00+00', 'Coriander loss from torn bag.');

INSERT INTO inventory_transfer (uuid, source_location_id, dest_location_id, transferred_at, notes)
VALUES
    ('85000000-0000-0000-0000-000000000001',
        (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'),
        (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000005'),
        '2026-01-10 10:00:00+00',
        'Move hops to brew deck for boil.'
    );

INSERT INTO inventory_transfer (uuid, source_location_id, dest_location_id, transferred_at, notes)
VALUES
    ('85000000-0000-0000-0000-000000000002',
        (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'),
        (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000005'),
        '2026-01-24 07:45:00+00',
        'Move Saaz to brew deck for Kolsch kettle additions.'
    );

INSERT INTO beer_lot (uuid, production_batch_uuid, lot_code, packaged_at, notes)
VALUES
    ('86000000-0000-0000-0000-000000000001', '90000000-0000-0000-0000-000000000001', 'IPA24-07-01', '2026-01-25 16:00:00+00', 'First packaging run.');

INSERT INTO beer_lot (uuid, production_batch_uuid, lot_code, packaged_at, notes)
VALUES
    ('86000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000004', 'STOUT24-09-01', '2026-02-06 12:00:00+00', 'Stout packaging run.');

INSERT INTO inventory_movement (
    ingredient_lot_id,
    stock_location_id,
    direction,
    reason,
    amount,
    amount_unit,
    occurred_at,
    receipt_id
)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000001'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 1000, 'kg', '2026-01-09 14:15:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000001')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000002'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 800, 'kg', '2026-01-09 14:20:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000001')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000003'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 200, 'kg', '2026-01-09 14:25:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000001')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000004'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 50, 'kg', '2026-01-10 09:40:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000002')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000005'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 40, 'kg', '2026-01-10 09:45:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000002')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000006'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 10, 'kg', '2026-01-11 08:35:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000003')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000007'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000004'), 'in', 'receive', 500, 'lb', '2026-01-22 15:10:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000004'));

INSERT INTO inventory_movement (
    ingredient_lot_id,
    stock_location_id,
    direction,
    reason,
    amount,
    amount_unit,
    occurred_at,
    receipt_id
)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000008'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 150, 'kg', '2026-01-18 13:05:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000005')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000009'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 30, 'kg', '2026-01-18 13:10:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000005')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000010'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 8, 'kg', '2026-01-22 10:05:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000006')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000011'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 20, 'kg', '2026-01-18 13:15:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000007')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000012'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 10, 'kg', '2026-01-18 13:16:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000007'));

INSERT INTO inventory_movement (
    ingredient_lot_id,
    stock_location_id,
    direction,
    reason,
    amount,
    amount_unit,
    occurred_at,
    usage_id
)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000001'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'out', 'use', 220, 'kg', '2026-01-10 08:30:00+00', (SELECT id FROM inventory_usage WHERE uuid = '83000000-0000-0000-0000-000000000001')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000004'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'out', 'use', 5, 'kg', '2026-01-10 11:05:00+00', (SELECT id FROM inventory_usage WHERE uuid = '83000000-0000-0000-0000-000000000001'));

INSERT INTO inventory_movement (
    ingredient_lot_id,
    stock_location_id,
    direction,
    reason,
    amount,
    amount_unit,
    occurred_at,
    usage_id
)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000008'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'out', 'use', 50, 'kg', '2026-01-22 08:25:00+00', (SELECT id FROM inventory_usage WHERE uuid = '83000000-0000-0000-0000-000000000002')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000001'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'out', 'use', 180, 'kg', '2026-01-22 08:25:00+00', (SELECT id FROM inventory_usage WHERE uuid = '83000000-0000-0000-0000-000000000002')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000002'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'out', 'use', 180, 'kg', '2026-01-24 08:20:00+00', (SELECT id FROM inventory_usage WHERE uuid = '83000000-0000-0000-0000-000000000003')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000009'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'out', 'use', 4, 'kg', '2026-01-24 08:20:00+00', (SELECT id FROM inventory_usage WHERE uuid = '83000000-0000-0000-0000-000000000003')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000010'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'out', 'use', 3, 'kg', '2026-01-24 13:40:00+00', (SELECT id FROM inventory_usage WHERE uuid = '83000000-0000-0000-0000-000000000003')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000001'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'out', 'use', 190, 'kg', '2026-01-26 08:10:00+00', (SELECT id FROM inventory_usage WHERE uuid = '83000000-0000-0000-0000-000000000004')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000006'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'out', 'use', 2, 'kg', '2026-01-26 08:10:00+00', (SELECT id FROM inventory_usage WHERE uuid = '83000000-0000-0000-0000-000000000004')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000012'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'out', 'use', 1, 'kg', '2026-01-26 08:10:00+00', (SELECT id FROM inventory_usage WHERE uuid = '83000000-0000-0000-0000-000000000004'));

INSERT INTO inventory_movement (
    ingredient_lot_id,
    stock_location_id,
    direction,
    reason,
    amount,
    amount_unit,
    occurred_at,
    transfer_id
)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000005'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'out', 'transfer', 5, 'kg', '2026-01-10 10:05:00+00', (SELECT id FROM inventory_transfer WHERE uuid = '85000000-0000-0000-0000-000000000001')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000005'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000005'), 'in', 'transfer', 5, 'kg', '2026-01-10 10:05:00+00', (SELECT id FROM inventory_transfer WHERE uuid = '85000000-0000-0000-0000-000000000001'));

INSERT INTO inventory_movement (
    ingredient_lot_id,
    stock_location_id,
    direction,
    reason,
    amount,
    amount_unit,
    occurred_at,
    transfer_id
)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000009'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'out', 'transfer', 4, 'kg', '2026-01-24 07:45:00+00', (SELECT id FROM inventory_transfer WHERE uuid = '85000000-0000-0000-0000-000000000002')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000009'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000005'), 'in', 'transfer', 4, 'kg', '2026-01-24 07:45:00+00', (SELECT id FROM inventory_transfer WHERE uuid = '85000000-0000-0000-0000-000000000002'));

INSERT INTO inventory_movement (
    ingredient_lot_id,
    stock_location_id,
    direction,
    reason,
    amount,
    amount_unit,
    occurred_at,
    adjustment_id
)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000003'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'out', 'adjust', 5, 'kg', '2026-01-15 07:05:00+00', (SELECT id FROM inventory_adjustment WHERE uuid = '84000000-0000-0000-0000-000000000001'));

INSERT INTO inventory_movement (
    ingredient_lot_id,
    stock_location_id,
    direction,
    reason,
    amount,
    amount_unit,
    occurred_at,
    adjustment_id
)
VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000012'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'out', 'adjust', 1, 'kg', '2026-01-23 07:20:00+00', (SELECT id FROM inventory_adjustment WHERE uuid = '84000000-0000-0000-0000-000000000002'));
