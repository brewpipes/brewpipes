-- Enhanced Batch Tracking: style, recipe, brew_session tables; occupancy status; volume_id on addition/measurement; BBL unit

-- Style table for beer styles (case-insensitive unique names)
CREATE TABLE IF NOT EXISTS style (
    id          serial PRIMARY KEY,
    uuid        uuid NOT NULL DEFAULT gen_random_uuid(),
    name        varchar(255) NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at  timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS style_uuid_idx ON style(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS style_name_lower_idx ON style(lower(name)) WHERE deleted_at IS NULL;

-- Recipe table for beer formulations
CREATE TABLE IF NOT EXISTS recipe (
    id          serial PRIMARY KEY,
    uuid        uuid NOT NULL DEFAULT gen_random_uuid(),
    name        varchar(255) NOT NULL,
    style_id    int REFERENCES style(id),
    style_name  varchar(255),
    notes       text,
    created_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at  timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS recipe_uuid_idx ON recipe(uuid);
CREATE INDEX IF NOT EXISTS recipe_style_id_idx ON recipe(style_id) WHERE deleted_at IS NULL;

-- Brew session table for hot-side wort production
CREATE TABLE IF NOT EXISTS brew_session (
    id              serial PRIMARY KEY,
    uuid            uuid NOT NULL DEFAULT gen_random_uuid(),
    batch_id        int REFERENCES batch(id),
    wort_volume_id  int REFERENCES volume(id),
    mash_vessel_id  int REFERENCES vessel(id),
    boil_vessel_id  int REFERENCES vessel(id),
    brewed_at       timestamptz NOT NULL,
    notes           text,
    created_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at      timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at      timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS brew_session_uuid_idx ON brew_session(uuid);
CREATE INDEX IF NOT EXISTS brew_session_batch_id_idx ON brew_session(batch_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS brew_session_wort_volume_id_idx ON brew_session(wort_volume_id) WHERE deleted_at IS NULL;

-- Add recipe_id to batch
ALTER TABLE batch ADD COLUMN recipe_id int REFERENCES recipe(id);
CREATE INDEX IF NOT EXISTS batch_recipe_id_idx ON batch(recipe_id) WHERE deleted_at IS NULL;

-- Add status to occupancy
ALTER TABLE occupancy ADD COLUMN status varchar(32);
ALTER TABLE occupancy ADD CONSTRAINT occupancy_status_check 
    CHECK (status IS NULL OR status IN (
        'fermenting', 'conditioning', 'cold_crashing', 
        'dry_hopping', 'carbonating', 'holding', 'packaging'
    ));

-- Add volume_id to addition
ALTER TABLE addition ADD COLUMN volume_id int REFERENCES volume(id);
CREATE INDEX IF NOT EXISTS addition_volume_id_idx ON addition(volume_id) WHERE deleted_at IS NULL;

-- Update addition target constraint to include volume_id option
ALTER TABLE addition DROP CONSTRAINT addition_target_check;
ALTER TABLE addition ADD CONSTRAINT addition_target_check CHECK (
    (batch_id IS NOT NULL AND occupancy_id IS NULL AND volume_id IS NULL) OR
    (batch_id IS NULL AND occupancy_id IS NOT NULL AND volume_id IS NULL) OR
    (batch_id IS NULL AND occupancy_id IS NULL AND volume_id IS NOT NULL)
);

-- Add volume_id to measurement
ALTER TABLE measurement ADD COLUMN volume_id int REFERENCES volume(id);
CREATE INDEX IF NOT EXISTS measurement_volume_id_idx ON measurement(volume_id) WHERE deleted_at IS NULL;

-- Update measurement target constraint to include volume_id option
ALTER TABLE measurement DROP CONSTRAINT measurement_target_check;
ALTER TABLE measurement ADD CONSTRAINT measurement_target_check CHECK (
    (batch_id IS NOT NULL AND occupancy_id IS NULL AND volume_id IS NULL) OR
    (batch_id IS NULL AND occupancy_id IS NOT NULL AND volume_id IS NULL) OR
    (batch_id IS NULL AND occupancy_id IS NULL AND volume_id IS NOT NULL)
);

-- Seed data for enhanced batch tracking tables.

-- Beer styles
INSERT INTO style (uuid, name)
VALUES
    ('a0000000-0000-0000-0000-000000000001', 'American IPA'),
    ('a0000000-0000-0000-0000-000000000002', 'German Pilsner'),
    ('a0000000-0000-0000-0000-000000000003', 'Irish Stout'),
    ('a0000000-0000-0000-0000-000000000004', 'Kölsch'),
    ('a0000000-0000-0000-0000-000000000005', 'Belgian Saison');

-- Recipes (linked to styles)
INSERT INTO recipe (uuid, name, style_id, style_name, notes)
VALUES
    ('a1000000-0000-0000-0000-000000000001', 'Flagship IPA', (SELECT id FROM style WHERE uuid = 'a0000000-0000-0000-0000-000000000001'), 'American IPA', 'House IPA with Citra and Mosaic.'),
    ('a1000000-0000-0000-0000-000000000002', 'Crisp Pilsner', (SELECT id FROM style WHERE uuid = 'a0000000-0000-0000-0000-000000000002'), 'German Pilsner', 'Classic German-style Pilsner with Saaz hops.'),
    ('a1000000-0000-0000-0000-000000000003', 'Robust Stout', (SELECT id FROM style WHERE uuid = 'a0000000-0000-0000-0000-000000000003'), 'Irish Stout', 'Rich and roasty with coffee notes.'),
    ('a1000000-0000-0000-0000-000000000004', 'Cologne Classic', (SELECT id FROM style WHERE uuid = 'a0000000-0000-0000-0000-000000000004'), 'Kölsch', 'Clean and crisp Kölsch-style ale.'),
    ('a1000000-0000-0000-0000-000000000005', 'Farmhouse Saison', (SELECT id FROM style WHERE uuid = 'a0000000-0000-0000-0000-000000000005'), 'Belgian Saison', 'Spicy and dry with coriander.');

-- Link batches to recipes
UPDATE batch SET recipe_id = (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000001')
WHERE uuid IN ('90000000-0000-0000-0000-000000000001', '90000000-0000-0000-0000-000000000002');

UPDATE batch SET recipe_id = (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000002')
WHERE uuid = '90000000-0000-0000-0000-000000000003';

UPDATE batch SET recipe_id = (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000003')
WHERE uuid = '90000000-0000-0000-0000-000000000004';

UPDATE batch SET recipe_id = (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000004')
WHERE uuid IN ('90000000-0000-0000-0000-000000000005', '90000000-0000-0000-0000-000000000006', '90000000-0000-0000-0000-000000000007');

UPDATE batch SET recipe_id = (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000005')
WHERE uuid = '90000000-0000-0000-0000-000000000008';

-- Update occupancy statuses for existing seed data
-- IPA 24-07 split A (fermenter A) - still fermenting
UPDATE occupancy SET status = 'fermenting'
WHERE uuid = '93100000-0000-0000-0000-000000000003';

-- IPA 24-07 split B (fermenter B) - dry hopping
UPDATE occupancy SET status = 'dry_hopping'
WHERE uuid = '93100000-0000-0000-0000-000000000004';

-- Pilsner 24-08 (fermenter C) - conditioning
UPDATE occupancy SET status = 'conditioning'
WHERE uuid = '93100000-0000-0000-0000-000000000005';

-- Stout 24-09 (fermenter D) - completed, out_at is set
-- No status update needed for completed occupancies

-- Saison 24-11 (fermenter D) - fermenting
UPDATE occupancy SET status = 'fermenting'
WHERE uuid = '93100000-0000-0000-0000-000000000007';

-- Kolsch blend (brite tank 2) - carbonating
UPDATE occupancy SET status = 'carbonating'
WHERE uuid = '93100000-0000-0000-0000-000000000008';

-- Brew sessions for each batch's hot-side production
INSERT INTO brew_session (uuid, batch_id, wort_volume_id, mash_vessel_id, boil_vessel_id, brewed_at, notes)
VALUES
    -- IPA 24-07 brew session (parent batch)
    ('a2000000-0000-0000-0000-000000000001',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000001'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'),
     '2026-01-10 08:00:00+00',
     'Single infusion mash at 152°F, 60-minute boil.'),
    -- Pilsner 24-08 brew session
    ('a2000000-0000-0000-0000-000000000002',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000004'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'),
     '2026-01-18 08:00:00+00',
     'Step mash for crisp finish, 90-minute boil.'),
    -- Stout 24-09 brew session
    ('a2000000-0000-0000-0000-000000000003',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000005'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'),
     '2026-01-22 08:00:00+00',
     'Thick mash with roasted barley, 60-minute boil.'),
    -- Kolsch 24-10 brew session (parent batch)
    ('a2000000-0000-0000-0000-000000000004',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000007'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'),
     '2026-01-24 08:00:00+00',
     'Single infusion at 148°F for dry finish.'),
    -- Saison 24-11 brew session
    ('a2000000-0000-0000-0000-000000000005',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000011'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'),
     '2026-01-26 08:00:00+00',
     'Step mash with pilsner base, 75-minute boil.');

-- Additional additions using volume_id target option (volume-based tracking)
INSERT INTO addition (
    uuid,
    batch_id,
    occupancy_id,
    volume_id,
    addition_type,
    stage,
    inventory_lot_uuid,
    amount,
    amount_unit,
    added_at,
    notes
)
VALUES
    -- Priming sugar added to Kolsch blend volume (volume-targeted)
    ('97000000-0000-0000-0000-000000000011', NULL, NULL, (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000010'), 'adjunct', 'packaging', '80000000-0000-0000-0000-000000000013', 2, 'kg', '2026-02-06 08:00:00+00', 'Priming sugar for bottle conditioning.'),
    -- Finings added to Pilsner wort volume (volume-targeted)
    ('97000000-0000-0000-0000-000000000012', NULL, NULL, (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000004'), 'other', 'boil', NULL, 100, 'g', '2026-01-18 11:30:00+00', 'Whirlfloc tablet for clarity.');

-- Additional measurements using volume_id target option (volume-based tracking)
INSERT INTO measurement (
    uuid,
    batch_id,
    occupancy_id,
    volume_id,
    kind,
    value,
    unit,
    observed_at,
    notes
)
VALUES
    -- Volume-targeted gravity reading on IPA split A
    ('98000000-0000-0000-0000-000000000014', NULL, NULL, (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000002'), 'gravity', 1.0120, 'sg', '2026-01-18 10:00:00+00', 'IPA split A final gravity.'),
    -- Volume-targeted gravity reading on IPA split B
    ('98000000-0000-0000-0000-000000000015', NULL, NULL, (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000003'), 'gravity', 1.0110, 'sg', '2026-01-18 10:15:00+00', 'IPA split B final gravity.'),
    -- Volume-targeted ABV on Stout conditioned beer
    ('98000000-0000-0000-0000-000000000016', NULL, NULL, (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000006'), 'abv', 5.2, 'pct', '2026-02-02 12:00:00+00', 'Stout final ABV.'),
    -- Volume-targeted pH on Saison beer
    ('98000000-0000-0000-0000-000000000017', NULL, NULL, (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000012'), 'ph', 4.10, NULL, '2026-01-30 08:00:00+00', 'Saison post-fermentation pH.');

-- ==============================================================================
-- FIX: Add missing occupancies for volumes that have batch_volume records
-- ==============================================================================

-- Stout wort in mash tun (volume 5)
INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at, status)
VALUES
    ('93100000-0000-0000-0000-000000000009', 
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'), 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000005'), 
     '2026-01-22 08:00:00+00', '2026-01-22 10:30:00+00', NULL);

-- Stout wort in kettle (volume 5)
INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at, status)
VALUES
    ('93100000-0000-0000-0000-000000000010', 
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'), 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000005'), 
     '2026-01-22 11:00:00+00', '2026-01-22 17:30:00+00', NULL);

-- Kolsch wort in mash tun (volume 7)
INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at, status)
VALUES
    ('93100000-0000-0000-0000-000000000011', 
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'), 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000007'), 
     '2026-01-24 08:00:00+00', '2026-01-24 10:00:00+00', NULL);

-- Kolsch wort in kettle (volume 7)
INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at, status)
VALUES
    ('93100000-0000-0000-0000-000000000012', 
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'), 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000007'), 
     '2026-01-24 10:30:00+00', '2026-01-24 13:30:00+00', NULL);

-- Kolsch split A in fermenter A (volume 8) - available after IPA moved to conditioning
INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at, status)
VALUES
    ('93100000-0000-0000-0000-000000000013', 
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000003'), 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000008'), 
     '2026-01-24 14:00:00+00', '2026-02-04 10:00:00+00', 'conditioning');

-- Kolsch split B in fermenter B (volume 9) - available after IPA moved to conditioning
INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at, status)
VALUES
    ('93100000-0000-0000-0000-000000000014', 
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000004'), 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000009'), 
     '2026-01-24 14:15:00+00', '2026-02-04 10:15:00+00', 'conditioning');

-- Saison wort in mash tun (volume 11)
INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at, status)
VALUES
    ('93100000-0000-0000-0000-000000000015', 
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'), 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000011'), 
     '2026-01-26 08:00:00+00', '2026-01-26 10:00:00+00', NULL);

-- Saison wort in kettle (volume 11)
INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at, status)
VALUES
    ('93100000-0000-0000-0000-000000000016', 
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'), 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000011'), 
     '2026-01-26 10:30:00+00', '2026-01-26 18:30:00+00', NULL);

-- Stout conditioned beer in brite tank 1 (volume 6) after leaving fermenter D
INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at, status)
VALUES
    ('93100000-0000-0000-0000-000000000017', 
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000006'), 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000006'), 
     '2026-01-26 10:30:00+00', '2026-02-05 08:00:00+00', 'conditioning');

-- ==============================================================================
-- FIX: Add missing volume_relation records for wort → beer progressions
-- These track the liquid transformation (same batch, different phase)
-- ==============================================================================

INSERT INTO volume_relation (uuid, parent_volume_id, child_volume_id, relation_type, amount, amount_unit)
VALUES
    -- Stout wort (900l) → Stout beer (870l) - 30l loss during fermentation
    ('92000000-0000-0000-0000-000000000007', 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000005'), 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000006'), 
     'split', 870, 'l'),
    -- Saison wort (950l) → Saison beer (930l) - 20l loss during fermentation
    ('92000000-0000-0000-0000-000000000008', 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000011'), 
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000012'), 
     'split', 930, 'l');

-- ==============================================================================
-- FIX: Add missing transfers for mash → kettle → fermenter workflows
-- ==============================================================================

-- Stout transfers: mash tun → kettle → fermenter D → brite tank 1
INSERT INTO transfer (uuid, source_occupancy_id, dest_occupancy_id, amount, amount_unit, loss_amount, loss_unit, started_at, ended_at)
VALUES
    ('93200000-0000-0000-0000-000000000004', 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000009'), 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000010'), 
     900, 'l', 5, 'l', '2026-01-22 10:35:00+00', '2026-01-22 10:55:00+00'),
    ('93200000-0000-0000-0000-000000000005', 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000010'), 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000006'), 
     895, 'l', 5, 'l', '2026-01-22 17:35:00+00', '2026-01-22 18:00:00+00'),
    ('93200000-0000-0000-0000-000000000013', 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000006'), 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000017'), 
     870, 'l', 10, 'l', '2026-01-26 10:00:00+00', '2026-01-26 10:30:00+00');

-- Kolsch transfers: mash tun → kettle → split into fermenters A & B
INSERT INTO transfer (uuid, source_occupancy_id, dest_occupancy_id, amount, amount_unit, loss_amount, loss_unit, started_at, ended_at)
VALUES
    ('93200000-0000-0000-0000-000000000006', 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000011'), 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000012'), 
     1000, 'l', 5, 'l', '2026-01-24 10:05:00+00', '2026-01-24 10:25:00+00'),
    ('93200000-0000-0000-0000-000000000007', 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000012'), 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000013'), 
     500, 'l', 3, 'l', '2026-01-24 13:35:00+00', '2026-01-24 14:00:00+00'),
    ('93200000-0000-0000-0000-000000000008', 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000012'), 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000014'), 
     490, 'l', 2, 'l', '2026-01-24 13:35:00+00', '2026-01-24 14:15:00+00');

-- Kolsch blend transfers: fermenters A & B → brite tank 2
INSERT INTO transfer (uuid, source_occupancy_id, dest_occupancy_id, amount, amount_unit, loss_amount, loss_unit, started_at, ended_at)
VALUES
    ('93200000-0000-0000-0000-000000000009', 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000013'), 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000008'), 
     500, 'l', 5, 'l', '2026-02-05 08:00:00+00', '2026-02-05 08:30:00+00'),
    ('93200000-0000-0000-0000-000000000010', 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000014'), 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000008'), 
     490, 'l', 5, 'l', '2026-02-05 08:00:00+00', '2026-02-05 09:00:00+00');

-- Saison transfers: mash tun → kettle → fermenter D
INSERT INTO transfer (uuid, source_occupancy_id, dest_occupancy_id, amount, amount_unit, loss_amount, loss_unit, started_at, ended_at)
VALUES
    ('93200000-0000-0000-0000-000000000011', 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000015'), 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000016'), 
     950, 'l', 5, 'l', '2026-01-26 10:05:00+00', '2026-01-26 10:25:00+00'),
    ('93200000-0000-0000-0000-000000000012', 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000016'), 
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000007'), 
     945, 'l', 5, 'l', '2026-01-26 18:35:00+00', '2026-01-26 19:00:00+00');
