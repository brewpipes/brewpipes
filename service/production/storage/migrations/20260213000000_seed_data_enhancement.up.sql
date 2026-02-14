-- Seed data enhancement: new batches, vessels, styles, recipes, occupancies, brew sessions.
-- Covers: batch without recipe, planning-only batch, multi-brew-session batch, soft-deleted batch,
-- retired/inactive vessels, empty active vessel, cold_crashing/holding/packaging occupancy statuses,
-- gas and other addition types, blend volume relations.
BEGIN;

-- ==============================================================================
-- GROUP 1: New Styles and Recipes
-- ==============================================================================

-- New styles
INSERT INTO style (uuid, name) VALUES
    ('a0000000-0000-0000-0000-000000000006', 'American Wheat'),
    ('a0000000-0000-0000-0000-000000000007', 'English Porter'),
    ('a0000000-0000-0000-0000-000000000008', 'Double IPA');

-- Soft-deleted style
INSERT INTO style (uuid, name, deleted_at) VALUES
    ('a0000000-0000-0000-0000-000000000009', 'American Amber Ale', '2026-02-10 12:00:00+00');

-- New recipes with target specs
INSERT INTO recipe (uuid, name, style_id, style_name, notes, batch_size, batch_size_unit, target_og, target_fg, target_ibu, target_srm, target_carbonation, ibu_method, brewhouse_efficiency) VALUES
    ('a1000000-0000-0000-0000-000000000006', 'Hearty Porter',
     (SELECT id FROM style WHERE uuid = 'a0000000-0000-0000-0000-000000000007'), 'English Porter',
     'Rich English porter with chocolate and toffee notes.',
     1200, 'l', 1.0540, 1.0140, 32, 28, 2.1, 'tinseth', 73.0),
    ('a1000000-0000-0000-0000-000000000007', 'Hop Bomb DIPA',
     (SELECT id FROM style WHERE uuid = 'a0000000-0000-0000-0000-000000000008'), 'Double IPA',
     'Aggressive double IPA with massive hop character.',
     1000, 'l', 1.0750, 1.0140, 85, 7, 2.5, 'tinseth', 72.0);

-- Soft-deleted recipe
INSERT INTO recipe (uuid, name, style_id, style_name, notes, deleted_at) VALUES
    ('a1000000-0000-0000-0000-000000000008', 'Autumn Amber',
     (SELECT id FROM style WHERE uuid = 'a0000000-0000-0000-0000-000000000009'), 'American Amber Ale',
     'Balanced amber ale — recipe retired.',
     '2026-02-10 12:00:00+00');

-- ==============================================================================
-- GROUP 2: New Batches
-- ==============================================================================

-- Batch without recipe (tests conditional UI rendering)
INSERT INTO batch (uuid, short_name, brew_date, notes) VALUES
    ('90000000-0000-0000-0000-000000000009', 'Wheat 24-12', '2026-02-10', 'Experimental wheat beer — no recipe, brewed from notes.');

-- Planning-only batch (no volumes/sessions/occupancies)
INSERT INTO batch (uuid, short_name, brew_date, notes, recipe_id) VALUES
    ('90000000-0000-0000-0000-000000000010', 'DIPA 24-13', NULL, 'Double IPA planned for late February. No brew date set yet.',
     (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000007'));

-- Multi-brew-session batch (double batching)
INSERT INTO batch (uuid, short_name, brew_date, notes, recipe_id) VALUES
    ('90000000-0000-0000-0000-000000000011', 'Porter 24-14', '2026-02-07', 'Double-brewed porter — two brew sessions blended into one fermenter.',
     (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000006'));

-- Soft-deleted batch
INSERT INTO batch (uuid, short_name, brew_date, notes, recipe_id, deleted_at) VALUES
    ('90000000-0000-0000-0000-000000000012', 'Amber 24-15', '2026-02-01', 'Cancelled amber ale — recipe retired.',
     (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000008'),
     '2026-02-10 12:00:00+00');

-- ==============================================================================
-- GROUP 3: New Vessels
-- ==============================================================================

INSERT INTO vessel (uuid, type, name, capacity, capacity_unit, make, model, status) VALUES
    ('93000000-0000-0000-0000-000000000009', 'fermenter', 'Fermenter E', 1200, 'l', 'SS Brewtech', 'FV-12', 'active'),
    ('93000000-0000-0000-0000-000000000010', 'fermenter', 'Old Fermenter', 600, 'l', 'Blichmann', 'F-6', 'retired'),
    ('93000000-0000-0000-0000-000000000011', 'kettle', 'Pilot Kettle', 200, 'l', 'Blichmann', 'K-2', 'inactive'),
    ('93000000-0000-0000-0000-000000000012', 'fermenter', 'Fermenter G', 1000, 'l', 'SS Brewtech', 'FV-10', 'active'),
    ('93000000-0000-0000-0000-000000000013', 'fermenter', 'Fermenter H', 1200, 'l', 'SS Brewtech', 'FV-12', 'active');

-- ==============================================================================
-- GROUP 4: New Volumes
-- ==============================================================================

INSERT INTO volume (uuid, name, description, amount, amount_unit) VALUES
    ('91000000-0000-0000-0000-000000000020', 'Porter 24-14 wort A', 'First brew session wort.', 600, 'l'),
    ('91000000-0000-0000-0000-000000000021', 'Porter 24-14 wort B', 'Second brew session wort.', 600, 'l'),
    ('91000000-0000-0000-0000-000000000022', 'Porter 24-14 combined', 'Blended porter from double brew.', 1180, 'l'),
    ('91000000-0000-0000-0000-000000000023', 'Wheat 24-12 wort', 'Experimental wheat wort.', 500, 'l');

-- ==============================================================================
-- GROUP 5: Volume Relations (blend for Porter double batch)
-- ==============================================================================

INSERT INTO volume_relation (uuid, parent_volume_id, child_volume_id, relation_type, amount, amount_unit) VALUES
    ('92000000-0000-0000-0000-000000000010',
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000020'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000022'),
     'blend', 590, 'l'),
    ('92000000-0000-0000-0000-000000000011',
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000021'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000022'),
     'blend', 590, 'l');

-- ==============================================================================
-- GROUP 6: Occupancies
-- ==============================================================================

INSERT INTO occupancy (uuid, vessel_id, volume_id, in_at, out_at, status) VALUES
    -- Porter wort A: mash tun (completed)
    ('93100000-0000-0000-0000-000000000020',
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000020'),
     '2026-02-07 08:00:00+00', '2026-02-07 10:00:00+00', NULL),
    -- Porter wort A: kettle (completed)
    ('93100000-0000-0000-0000-000000000021',
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000020'),
     '2026-02-07 10:30:00+00', '2026-02-07 12:00:00+00', NULL),
    -- Porter wort B: mash tun (completed)
    ('93100000-0000-0000-0000-000000000022',
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000021'),
     '2026-02-08 08:00:00+00', '2026-02-08 10:00:00+00', NULL),
    -- Porter wort B: kettle (completed)
    ('93100000-0000-0000-0000-000000000023',
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000021'),
     '2026-02-08 10:30:00+00', '2026-02-08 12:00:00+00', NULL),
    -- Porter combined in Fermenter E (cold_crashing, then moved to brite)
    ('93100000-0000-0000-0000-000000000024',
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000009'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000022'),
     '2026-02-08 13:00:00+00', '2026-02-12 08:00:00+00', 'cold_crashing'),
    -- Porter combined in Brite Tank 1 (packaging — current active occupancy)
    ('93100000-0000-0000-0000-000000000025',
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000006'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000022'),
     '2026-02-12 09:00:00+00', NULL, 'packaging'),
    -- Wheat wort in Fermenter G (holding — current active occupancy)
    ('93100000-0000-0000-0000-000000000026',
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000012'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000023'),
     '2026-02-10 14:00:00+00', NULL, 'holding');

-- ==============================================================================
-- GROUP 7: Transfers (Porter brew day flow)
-- ==============================================================================

INSERT INTO transfer (uuid, source_occupancy_id, dest_occupancy_id, amount, amount_unit, loss_amount, loss_unit, started_at, ended_at) VALUES
    -- Porter wort A: mash → kettle
    ('93200000-0000-0000-0000-000000000020',
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000020'),
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000021'),
     600, 'l', 3, 'l', '2026-02-07 10:05:00+00', '2026-02-07 10:25:00+00'),
    -- Porter wort A: kettle → fermenter E
    ('93200000-0000-0000-0000-000000000021',
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000021'),
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000024'),
     597, 'l', 3, 'l', '2026-02-07 12:05:00+00', '2026-02-07 12:30:00+00'),
    -- Porter wort B: mash → kettle
    ('93200000-0000-0000-0000-000000000022',
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000022'),
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000023'),
     600, 'l', 3, 'l', '2026-02-08 10:05:00+00', '2026-02-08 10:25:00+00'),
    -- Porter wort B: kettle → fermenter E (blended with wort A)
    ('93200000-0000-0000-0000-000000000023',
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000023'),
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000024'),
     597, 'l', 3, 'l', '2026-02-08 12:05:00+00', '2026-02-08 12:30:00+00'),
    -- Porter combined: fermenter E → brite tank 1
    ('93200000-0000-0000-0000-000000000024',
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000024'),
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000025'),
     1180, 'l', 10, 'l', '2026-02-12 08:05:00+00', '2026-02-12 08:45:00+00');

-- ==============================================================================
-- GROUP 8: Batch Volumes
-- ==============================================================================

INSERT INTO batch_volume (uuid, batch_id, volume_id, liquid_phase, phase_at) VALUES
    ('94000000-0000-0000-0000-000000000020',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000011'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000020'),
     'wort', '2026-02-07 09:00:00+00'),
    ('94000000-0000-0000-0000-000000000021',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000011'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000021'),
     'wort', '2026-02-08 09:00:00+00'),
    ('94000000-0000-0000-0000-000000000022',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000011'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000022'),
     'beer', '2026-02-09 12:00:00+00'),
    ('94000000-0000-0000-0000-000000000023',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000009'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000023'),
     'wort', '2026-02-10 09:00:00+00');

-- ==============================================================================
-- GROUP 9: Batch Process Phases
-- ==============================================================================

INSERT INTO batch_process_phase (uuid, batch_id, process_phase, phase_at) VALUES
    -- Wheat (no recipe) — currently fermenting
    ('95000000-0000-0000-0000-000000000040',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000009'),
     'planning', '2026-02-08 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000041',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000009'),
     'mashing', '2026-02-10 08:00:00+00'),
    ('95000000-0000-0000-0000-000000000042',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000009'),
     'boiling', '2026-02-10 10:30:00+00'),
    ('95000000-0000-0000-0000-000000000043',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000009'),
     'fermenting', '2026-02-10 14:00:00+00'),
    -- DIPA (planning only — no brew sessions/volumes/occupancies)
    ('95000000-0000-0000-0000-000000000044',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000010'),
     'planning', '2026-02-12 09:00:00+00'),
    -- Porter (multi-brew) — currently conditioning
    ('95000000-0000-0000-0000-000000000045',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000011'),
     'planning', '2026-02-04 09:00:00+00'),
    ('95000000-0000-0000-0000-000000000046',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000011'),
     'mashing', '2026-02-07 08:00:00+00'),
    ('95000000-0000-0000-0000-000000000047',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000011'),
     'boiling', '2026-02-07 10:30:00+00'),
    ('95000000-0000-0000-0000-000000000048',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000011'),
     'fermenting', '2026-02-08 13:00:00+00'),
    ('95000000-0000-0000-0000-000000000049',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000011'),
     'conditioning', '2026-02-12 09:00:00+00');

-- Soft-deleted batch process phase
INSERT INTO batch_process_phase (uuid, batch_id, process_phase, phase_at, deleted_at) VALUES
    ('95000000-0000-0000-0000-000000000050',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000012'),
     'planning', '2026-02-01 09:00:00+00', '2026-02-10 12:00:00+00');

-- ==============================================================================
-- GROUP 10: Brew Sessions
-- ==============================================================================

INSERT INTO brew_session (uuid, batch_id, wort_volume_id, mash_vessel_id, boil_vessel_id, brewed_at, notes) VALUES
    -- Porter: TWO brew sessions (double batching)
    ('a2000000-0000-0000-0000-000000000006',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000011'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000020'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'),
     '2026-02-07 08:00:00+00',
     'First brew session for double-brewed porter.'),
    ('a2000000-0000-0000-0000-000000000007',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000011'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000021'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'),
     '2026-02-08 08:00:00+00',
     'Second brew session for double-brewed porter.'),
    -- Wheat: single brew session
    ('a2000000-0000-0000-0000-000000000008',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000009'),
     (SELECT id FROM volume WHERE uuid = '91000000-0000-0000-0000-000000000023'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000001'),
     (SELECT id FROM vessel WHERE uuid = '93000000-0000-0000-0000-000000000002'),
     '2026-02-10 08:00:00+00',
     'Experimental wheat brew — no recipe, brewed from notes.');

-- ==============================================================================
-- GROUP 11: Additions (gas and other types)
-- ==============================================================================

INSERT INTO addition (uuid, batch_id, occupancy_id, addition_type, stage, inventory_lot_uuid, amount, amount_unit, added_at, notes) VALUES
    -- Gas addition: CO2 force carbonation on Porter in brite tank
    ('97000000-0000-0000-0000-000000000020', NULL,
     (SELECT id FROM occupancy WHERE uuid = '93100000-0000-0000-0000-000000000025'),
     'gas', 'packaging', NULL, 50, 'l', '2026-02-12 10:00:00+00', 'CO2 force carbonation.'),
    -- Other addition: yeast nutrient for Porter fermentation
    ('97000000-0000-0000-0000-000000000021',
     (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000011'), NULL,
     'other', 'fermentation', NULL, 200, 'g', '2026-02-08 14:00:00+00', 'Yeast nutrient addition.');

COMMIT;
