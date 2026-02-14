-- Seed data enhancement: new ingredients, lots, receipts, adjustments, movements.
-- Resolves cross-service recipe references, fills category/form/reason gaps.
BEGIN;

-- ==============================================================================
-- GROUP 1: New Ingredients (cross-service refs + category gaps)
-- ==============================================================================

-- New fermentables for recipe cross-service resolution
INSERT INTO ingredient (uuid, name, category, default_unit, description) VALUES
    ('70000000-0000-0000-0000-000000000032', 'Munich Malt', 'fermentable', 'kg', 'Bready, malty base/specialty malt.'),
    ('70000000-0000-0000-0000-000000000033', 'Chocolate Malt', 'fermentable', 'kg', 'Dark roasted malt for color and chocolate notes.'),
    ('70000000-0000-0000-0000-000000000034', 'Flaked Barley', 'fermentable', 'kg', 'Unmalted barley for body and head retention.'),
    ('70000000-0000-0000-0000-000000000035', 'Vienna Malt', 'fermentable', 'kg', 'Lightly kilned malt with biscuit character.'),
    ('70000000-0000-0000-0000-000000000036', 'Wheat Malt', 'fermentable', 'kg', 'Malted wheat for haze and head retention.');

-- New hops for recipe cross-service resolution (with form variety)
INSERT INTO ingredient (uuid, name, category, default_unit, description) VALUES
    ('70000000-0000-0000-0000-000000000025', 'Columbus Hops', 'hop', 'kg', 'High-alpha bittering hop with earthy character.'),
    ('70000000-0000-0000-0000-000000000026', 'Mosaic Hops', 'hop', 'kg', 'Complex tropical and stone fruit aroma hop.'),
    ('70000000-0000-0000-0000-000000000027', 'East Kent Goldings', 'hop', 'kg', 'Classic English aroma hop with floral character.'),
    ('70000000-0000-0000-0000-000000000028', 'Hallertau Mittelfruh', 'hop', 'kg', 'Noble German hop with mild spice and floral notes.'),
    ('70000000-0000-0000-0000-000000000029', 'Styrian Goldings', 'hop', 'kg', 'Slovenian hop with mild spice and earthy character.');

-- New yeasts for recipe cross-service resolution (with slurry form)
INSERT INTO ingredient (uuid, name, category, default_unit, description) VALUES
    ('70000000-0000-0000-0000-000000000020', 'WLP001 California Ale (Pitch)', 'yeast', 'kg', 'Clean American ale yeast — pitch-ready pack.'),
    ('70000000-0000-0000-0000-000000000021', 'W34/70 Lager (Pitch)', 'yeast', 'kg', 'Clean lager yeast — pitch-ready pack.'),
    ('70000000-0000-0000-0000-000000000022', 'WLP004 Irish Ale Yeast', 'yeast', 'kg', 'Slightly fruity, clean Irish ale strain.'),
    ('70000000-0000-0000-0000-000000000023', 'WLP029 Kolsch Yeast', 'yeast', 'kg', 'Clean, lager-like ale yeast for Kolsch.'),
    ('70000000-0000-0000-0000-000000000024', 'Belle Saison Yeast', 'yeast', 'kg', 'High-attenuating Belgian saison strain.');

-- New adjuncts for recipe cross-service resolution
INSERT INTO ingredient (uuid, name, category, default_unit, description) VALUES
    ('70000000-0000-0000-0000-000000000030', 'Coriander Seeds', 'adjunct', 'kg', 'Crushed coriander for Belgian and farmhouse styles.'),
    ('70000000-0000-0000-0000-000000000031', 'Orange Peel', 'adjunct', 'kg', 'Dried bitter orange peel for Belgian styles.');

-- New gas ingredient
INSERT INTO ingredient (uuid, name, category, default_unit, description) VALUES
    ('70000000-0000-0000-0000-000000000017', 'Nitrogen', 'gas', 'lb', 'Nitrogen gas for stout pours and nitrogenation.');

-- New "other" category ingredients (fills missing category gap)
INSERT INTO ingredient (uuid, name, category, default_unit, description) VALUES
    ('70000000-0000-0000-0000-000000000018', 'PBW Cleaner', 'other', 'kg', 'Alkaline brewery cleaner for CIP and soaking.'),
    ('70000000-0000-0000-0000-000000000019', 'Star San', 'other', 'l', 'Acid-based no-rinse sanitizer.');

-- ==============================================================================
-- GROUP 2: Ingredient Detail Records
-- ==============================================================================

-- Malt details for new fermentables
INSERT INTO ingredient_malt_detail (ingredient_id, maltster_name, variety, lovibond, srm, diastatic_power) VALUES
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000032'), 'Great Lakes Malt Co', 'Munich', 10.0, 10.0, 50.0),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000033'), 'Midwest Maltings', 'Chocolate', 350.0, 350.0, 0.0),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000034'), 'Great Lakes Malt Co', 'Flaked Barley', 2.0, 2.0, 0.0),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000035'), 'Great Lakes Malt Co', 'Vienna', 4.0, 4.0, 50.0),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000036'), 'Great Lakes Malt Co', 'Wheat', 2.0, 2.0, 60.0);

-- Hop details (with form variety: whole_leaf, cryo, pellet)
INSERT INTO ingredient_hop_detail (ingredient_id, producer_name, variety, crop_year, form, alpha_acid, beta_acid) VALUES
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000025'), 'Northwest Hop Farms', 'Columbus', 2025, 'whole_leaf', 14.0, 5.5),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000026'), 'Northwest Hop Farms', 'Mosaic', 2025, 'cryo', 24.0, 3.5),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000027'), 'Charles Faram', 'East Kent Goldings', 2025, 'pellet', 5.5, 3.5),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000028'), 'Hopsteiner', 'Hallertau Mittelfruh', 2025, 'pellet', 4.5, 4.0),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000029'), 'Styrian Hops', 'Styrian Goldings', 2025, 'pellet', 5.0, 3.0);

-- Yeast details (with slurry form for Belle Saison)
INSERT INTO ingredient_yeast_detail (ingredient_id, lab_name, strain, form) VALUES
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000020'), 'White Labs', 'WLP001', 'liquid'),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000021'), 'Fermentis', 'W34/70', 'dry'),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000022'), 'White Labs', 'WLP004', 'liquid'),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000023'), 'White Labs', 'WLP029', 'liquid'),
    ((SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000024'), 'Lallemand', 'Belle Saison', 'slurry');

-- ==============================================================================
-- GROUP 3: New Receipts
-- ==============================================================================

INSERT INTO inventory_receipt (uuid, supplier_uuid, purchase_order_uuid, reference_code, received_at, notes) VALUES
    ('82000000-0000-0000-0000-000000000010', '11111111-1111-1111-1111-111111111111', '50000000-0000-0000-0000-000000000001', 'RCV-1010', '2026-02-01 14:00:00+00', 'Second pale malt delivery for spring brews.'),
    ('82000000-0000-0000-0000-000000000011', '22222222-2222-2222-2222-222222222222', '50000000-0000-0000-0000-000000000002', 'RCV-1011', '2026-02-02 09:00:00+00', 'Specialty hops delivery — Columbus, Mosaic, EKG, Hallertau, Styrian.'),
    ('82000000-0000-0000-0000-000000000012', '33333333-3333-3333-3333-333333333333', '50000000-0000-0000-0000-000000000003', 'RCV-1012', '2026-02-03 08:30:00+00', 'Yeast delivery — WLP001, WLP004, WLP029.'),
    ('82000000-0000-0000-0000-000000000013', '77777777-7777-7777-7777-777777777777', '50000000-0000-0000-0000-000000000006', 'RCV-1013', '2026-02-03 10:00:00+00', 'Yeast delivery — W34/70 pitch, Belle Saison.'),
    ('82000000-0000-0000-0000-000000000014', '66666666-6666-6666-6666-666666666666', NULL, 'RCV-1014', '2026-02-04 13:00:00+00', 'Specialty malts and adjuncts — Munich, Chocolate, Flaked Barley, Vienna, Wheat, Coriander, Orange Peel.'),
    ('82000000-0000-0000-0000-000000000015', '55555555-5555-5555-5555-555555555555', NULL, 'RCV-1015', '2026-02-05 09:00:00+00', 'Cleaning supplies — PBW.'),
    ('82000000-0000-0000-0000-000000000020', NULL, NULL, 'RCV-PKG-001', '2026-01-25 16:30:00+00', 'IPA 24-07 packaging run into finished goods.'),
    ('82000000-0000-0000-0000-000000000021', NULL, NULL, 'RCV-PKG-002', '2026-02-06 12:30:00+00', 'Stout 24-09 packaging run into finished goods.');

-- ==============================================================================
-- GROUP 4: New Ingredient Lots
-- ==============================================================================

-- Lots for cross-service recipe ingredients: yeasts
INSERT INTO ingredient_lot (uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, notes) VALUES
    ('80000000-0000-0000-0000-000000000020', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000020'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000012'), '33333333-3333-3333-3333-333333333333', 'YEAST-WLP001P-2402', 'WL-001P-2402', 'White Labs', 'yeast_lab', '2026-02-03 08:30:00+00', 8, 'kg', '2026-04-01 00:00:00+00', 'WLP001 pitch-ready packs.'),
    ('80000000-0000-0000-0000-000000000021', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000021'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000013'), '77777777-7777-7777-7777-777777777777', 'YEAST-W3470P-2402', 'FER-W3470P-2402', 'Fermentis', 'yeast_lab', '2026-02-03 10:00:00+00', 6, 'kg', '2026-05-01 00:00:00+00', 'W34/70 pitch-ready packs.'),
    ('80000000-0000-0000-0000-000000000022', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000022'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000012'), '33333333-3333-3333-3333-333333333333', 'YEAST-WLP004-2402', 'WL-004-2402', 'White Labs', 'yeast_lab', '2026-02-03 08:30:00+00', 5, 'kg', '2026-04-01 00:00:00+00', 'Irish ale yeast.'),
    ('80000000-0000-0000-0000-000000000023', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000023'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000012'), '33333333-3333-3333-3333-333333333333', 'YEAST-WLP029-2402', 'WL-029-2402', 'White Labs', 'yeast_lab', '2026-02-03 08:30:00+00', 5, 'kg', '2026-04-01 00:00:00+00', 'Kolsch yeast.'),
    ('80000000-0000-0000-0000-000000000024', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000024'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000013'), '77777777-7777-7777-7777-777777777777', 'YEAST-BSAIS-2402', 'LAL-BSAIS-2402', 'Lallemand', 'yeast_lab', '2026-02-03 10:00:00+00', 4, 'kg', '2026-05-01 00:00:00+00', 'Belle Saison slurry.');

-- Lots for cross-service recipe ingredients: hops
INSERT INTO ingredient_lot (uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, notes) VALUES
    ('80000000-0000-0000-0000-000000000025', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000025'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011'), '22222222-2222-2222-2222-222222222222', 'HOP-COL-2501', 'NHF-COL-25A', 'Northwest Hop Farms', 'hop_producer', '2026-02-02 09:00:00+00', 20, 'kg', '2027-06-30 00:00:00+00', 'Columbus whole leaf crop 2025.'),
    ('80000000-0000-0000-0000-000000000026', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000026'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011'), '22222222-2222-2222-2222-222222222222', 'HOP-MOS-2501', 'NHF-MOS-25C', 'Northwest Hop Farms', 'hop_producer', '2026-02-02 09:00:00+00', 15, 'kg', '2027-06-30 00:00:00+00', 'Mosaic cryo hops crop 2025.'),
    ('80000000-0000-0000-0000-000000000027', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000027'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011'), '22222222-2222-2222-2222-222222222222', 'HOP-EKG-2501', 'CF-EKG-25A', 'Charles Faram', 'hop_producer', '2026-02-02 09:00:00+00', 10, 'kg', '2027-06-30 00:00:00+00', 'East Kent Goldings pellets.'),
    ('80000000-0000-0000-0000-000000000028', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000028'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011'), '22222222-2222-2222-2222-222222222222', 'HOP-HMF-2501', 'HS-HMF-25A', 'Hopsteiner', 'hop_producer', '2026-02-02 09:00:00+00', 12, 'kg', '2027-06-30 00:00:00+00', 'Hallertau Mittelfruh pellets.'),
    ('80000000-0000-0000-0000-000000000029', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000029'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011'), '22222222-2222-2222-2222-222222222222', 'HOP-STY-2501', 'SH-STY-25A', 'Styrian Hops', 'hop_producer', '2026-02-02 09:00:00+00', 8, 'kg', '2027-06-30 00:00:00+00', 'Styrian Goldings pellets.');

-- Lots for cross-service recipe ingredients: adjuncts
INSERT INTO ingredient_lot (uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, notes) VALUES
    ('80000000-0000-0000-0000-000000000030', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000030'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014'), '66666666-6666-6666-6666-666666666666', 'ADJ-CORI-2402', 'CORI-2402', 'Spice River Co', 'other', '2026-02-04 13:00:00+00', 5, 'kg', '2027-06-01 00:00:00+00', 'Coriander seeds for Belgian styles.'),
    ('80000000-0000-0000-0000-000000000031', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000031'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014'), '66666666-6666-6666-6666-666666666666', 'ADJ-ORPEEL-2402', 'ORPEEL-2402', 'Spice River Co', 'other', '2026-02-04 13:00:00+00', 3, 'kg', '2027-06-01 00:00:00+00', 'Dried bitter orange peel.');

-- Lots for cross-service recipe ingredients: fermentables
INSERT INTO ingredient_lot (uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, notes) VALUES
    ('80000000-0000-0000-0000-000000000032', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000032'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014'), '66666666-6666-6666-6666-666666666666', 'MALT-MUN-2402', 'GLM-MUN-2402', 'Great Lakes Malt Co', 'maltster', '2026-02-04 13:00:00+00', 300, 'kg', '2026-08-01 00:00:00+00', 'Munich malt for spring brews.'),
    ('80000000-0000-0000-0000-000000000033', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000033'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014'), '66666666-6666-6666-6666-666666666666', 'MALT-CHOC-2402', 'MM-CHOC-2402', 'Midwest Maltings', 'maltster', '2026-02-04 13:00:00+00', 100, 'kg', '2026-08-01 00:00:00+00', 'Chocolate malt for porters and stouts.'),
    ('80000000-0000-0000-0000-000000000034', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000034'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014'), '66666666-6666-6666-6666-666666666666', 'MALT-FLKB-2402', 'GLM-FLKB-2402', 'Great Lakes Malt Co', 'maltster', '2026-02-04 13:00:00+00', 50, 'kg', '2026-08-01 00:00:00+00', 'Flaked barley for head retention.'),
    ('80000000-0000-0000-0000-000000000035', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000035'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014'), '66666666-6666-6666-6666-666666666666', 'MALT-VIE-2402', 'GLM-VIE-2402', 'Great Lakes Malt Co', 'maltster', '2026-02-04 13:00:00+00', 200, 'kg', '2026-08-01 00:00:00+00', 'Vienna malt for Kolsch and amber styles.'),
    ('80000000-0000-0000-0000-000000000036', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000036'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014'), '66666666-6666-6666-6666-666666666666', 'MALT-WHT-2402', 'GLM-WHT-2402', 'Great Lakes Malt Co', 'maltster', '2026-02-04 13:00:00+00', 150, 'kg', '2026-08-01 00:00:00+00', 'Wheat malt for wheat and saison styles.');

-- Second lot of Pale Malt (multiple lots per ingredient, FIFO testing)
INSERT INTO ingredient_lot (uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, notes) VALUES
    ('80000000-0000-0000-0000-000000000040', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000001'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000010'), '11111111-1111-1111-1111-111111111111', 'MALT-PALE-2402', 'GLM-PALE-2402', 'Great Lakes Malt Co', 'maltster', '2026-02-01 14:00:00+00', 500, 'kg', '2026-08-01 00:00:00+00', 'Second pale malt lot — stored in Brite Cellar for overflow.');

-- Second lot of Citra Hops (multiple lots per ingredient, FIFO testing)
INSERT INTO ingredient_lot (uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, notes) VALUES
    ('80000000-0000-0000-0000-000000000041', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000005'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011'), '22222222-2222-2222-2222-222222222222', 'HOP-CIT-2502', 'NHF-CIT-25C', 'Northwest Hop Farms', 'hop_producer', '2026-02-02 09:00:00+00', 25, 'kg', '2027-06-30 00:00:00+00', 'Second Citra lot for spring IPA runs.');

-- Expired Crystal 60L lot (expired lot, also soft-deleted)
INSERT INTO ingredient_lot (uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, expires_at, notes, deleted_at) VALUES
    ('80000000-0000-0000-0000-000000000050', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000003'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000010'), '11111111-1111-1111-1111-111111111111', 'MALT-CR60-2301', 'Great Lakes Malt Co', 'maltster', '2025-06-15 14:00:00+00', 50, 'kg', '2025-12-15 00:00:00+00', '2026-01-15 00:00:00+00', 'Old Crystal 60L lot — expired, written off.', '2026-02-10 12:00:00+00');

-- Near-expiry Saaz lot (best_by approaching)
INSERT INTO ingredient_lot (uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, notes) VALUES
    ('80000000-0000-0000-0000-000000000051', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000012'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011'), '22222222-2222-2222-2222-222222222222', 'HOP-SAA-2401', 'Bohemian Hop Co', 'hop_producer', '2025-03-01 09:00:00+00', 10, 'kg', '2026-03-01 00:00:00+00', 'Old Saaz lot — approaching best-by date.');

-- PBW Cleaner lot ("other" category)
INSERT INTO ingredient_lot (uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, notes) VALUES
    ('80000000-0000-0000-0000-000000000052', (SELECT id FROM ingredient WHERE uuid = '70000000-0000-0000-0000-000000000018'), (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000015'), '55555555-5555-5555-5555-555555555555', 'CLN-PBW-2402', 'Five Star Chemicals', 'other', '2026-02-05 09:00:00+00', 25, 'kg', 'PBW cleaner for CIP system.');

-- ==============================================================================
-- GROUP 5: Lot Detail Records for New Lots
-- ==============================================================================

-- Malt details for new lots
INSERT INTO ingredient_lot_malt_detail (ingredient_lot_id, moisture_percent) VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000032'), 4.3),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000033'), 3.5),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000034'), 5.0),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000035'), 4.1),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000036'), 4.4),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000040'), 4.0);

-- Hop details for new lots
INSERT INTO ingredient_lot_hop_detail (ingredient_lot_id, alpha_acid, beta_acid) VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000025'), 14.2, 5.3),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000026'), 23.8, 3.6),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000027'), 5.3, 3.4),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000028'), 4.6, 4.1),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000029'), 5.1, 2.9),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000041'), 12.1, 4.2),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000051'), 3.2, 3.8);

-- Yeast details for new lots
INSERT INTO ingredient_lot_yeast_detail (ingredient_lot_id, viability_percent, generation) VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000020'), 96.0, 0),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000021'), 95.0, 0),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000022'), 94.0, 0),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000023'), 95.0, 0),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000024'), 88.0, 2);

-- ==============================================================================
-- GROUP 6: New Adjustments (missing reasons: shrink, damage, correction, other)
-- ==============================================================================

INSERT INTO inventory_adjustment (uuid, reason, adjusted_at, notes) VALUES
    ('84000000-0000-0000-0000-000000000003', 'shrink', '2026-02-05 07:00:00+00', 'Grain room shrink — moisture absorption.'),
    ('84000000-0000-0000-0000-000000000004', 'damage', '2026-02-06 14:00:00+00', 'Forklift punctured a bag of Munich malt.'),
    ('84000000-0000-0000-0000-000000000005', 'correction', '2026-02-07 07:00:00+00', 'Scale calibration correction on Citra inventory.'),
    ('84000000-0000-0000-0000-000000000006', 'other', '2026-02-08 07:00:00+00', 'Donated PBW to homebrew club.');

-- Waste/spoilage adjustment for expired lot
INSERT INTO inventory_adjustment (uuid, reason, adjusted_at, notes) VALUES
    ('84000000-0000-0000-0000-000000000007', 'spoilage', '2026-02-09 07:00:00+00', 'Expired Crystal 60L lot written off as waste.');

-- ==============================================================================
-- GROUP 7: Inventory Movements
-- ==============================================================================

-- Receive movements for all new lots
INSERT INTO inventory_movement (ingredient_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, receipt_id) VALUES
    -- Yeasts → Cold Room
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000020'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 8, 'kg', '2026-02-03 08:35:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000012')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000021'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 6, 'kg', '2026-02-03 10:05:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000013')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000022'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 5, 'kg', '2026-02-03 08:40:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000012')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000023'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 5, 'kg', '2026-02-03 08:45:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000012')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000024'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 4, 'kg', '2026-02-03 10:10:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000013')),
    -- Hops → Cold Room
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000025'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 20, 'kg', '2026-02-02 09:10:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000026'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 15, 'kg', '2026-02-02 09:15:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000027'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 10, 'kg', '2026-02-02 09:20:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000028'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 12, 'kg', '2026-02-02 09:25:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000029'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 8, 'kg', '2026-02-02 09:30:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011')),
    -- Adjuncts → Grain Room
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000030'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 5, 'kg', '2026-02-04 13:10:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000031'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 3, 'kg', '2026-02-04 13:15:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014')),
    -- Fermentables → Grain Room
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000032'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 300, 'kg', '2026-02-04 13:20:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000033'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 100, 'kg', '2026-02-04 13:25:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000034'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 50, 'kg', '2026-02-04 13:30:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000035'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 200, 'kg', '2026-02-04 13:35:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014')),
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000036'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'in', 'receive', 150, 'kg', '2026-02-04 13:40:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000014')),
    -- Second Pale Malt lot → Brite Cellar (different location for FIFO testing)
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000040'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000007'), 'in', 'receive', 500, 'kg', '2026-02-01 14:10:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000010')),
    -- Second Citra lot → Cold Room
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000041'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 25, 'kg', '2026-02-02 09:35:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011')),
    -- Near-expiry Saaz → Cold Room
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000051'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'receive', 10, 'kg', '2025-03-01 09:10:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000011')),
    -- PBW → Chemical Cage
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000052'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000003'), 'in', 'receive', 25, 'kg', '2026-02-05 09:10:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000015'));

-- Adjustment movements (shrink, damage, correction, other)
INSERT INTO inventory_movement (ingredient_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, adjustment_id) VALUES
    -- Shrink: Pale Malt moisture absorption
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000001'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'out', 'adjust', 10, 'kg', '2026-02-05 07:05:00+00', (SELECT id FROM inventory_adjustment WHERE uuid = '84000000-0000-0000-0000-000000000003')),
    -- Damage: Munich Malt forklift puncture
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000032'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'out', 'adjust', 25, 'kg', '2026-02-06 14:05:00+00', (SELECT id FROM inventory_adjustment WHERE uuid = '84000000-0000-0000-0000-000000000004')),
    -- Correction: Citra scale calibration (positive adjustment)
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000005'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000002'), 'in', 'adjust', 2, 'kg', '2026-02-07 07:05:00+00', (SELECT id FROM inventory_adjustment WHERE uuid = '84000000-0000-0000-0000-000000000005')),
    -- Other: PBW donation
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000052'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000003'), 'out', 'adjust', 5, 'kg', '2026-02-08 07:05:00+00', (SELECT id FROM inventory_adjustment WHERE uuid = '84000000-0000-0000-0000-000000000006'));

-- Waste movement (missing 'waste' reason)
INSERT INTO inventory_movement (ingredient_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, adjustment_id) VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000050'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'out', 'waste', 50, 'kg', '2026-02-09 07:05:00+00', (SELECT id FROM inventory_adjustment WHERE uuid = '84000000-0000-0000-0000-000000000007'));

-- Beer lot movements (existing beer lots have no movements)
INSERT INTO inventory_movement (beer_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, receipt_id) VALUES
    ((SELECT id FROM beer_lot WHERE uuid = '86000000-0000-0000-0000-000000000001'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000006'), 'in', 'receive', 300, 'l', '2026-01-25 16:30:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000020')),
    ((SELECT id FROM beer_lot WHERE uuid = '86000000-0000-0000-0000-000000000002'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000006'), 'in', 'receive', 400, 'l', '2026-02-06 12:30:00+00', (SELECT id FROM inventory_receipt WHERE uuid = '82000000-0000-0000-0000-000000000021'));

-- Zero-stock scenario: Fully deplete Irish Moss lot (lot ...0011 received 20 kg, no prior usage)
INSERT INTO inventory_usage (uuid, production_ref_uuid, used_at, notes) VALUES
    ('83000000-0000-0000-0000-000000000010', '90000000-0000-0000-0000-000000000003', '2026-01-18 11:30:00+00', 'Irish Moss usage for Pilsner clarity.');

INSERT INTO inventory_movement (ingredient_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, usage_id) VALUES
    ((SELECT id FROM ingredient_lot WHERE uuid = '80000000-0000-0000-0000-000000000011'), (SELECT id FROM stock_location WHERE uuid = '81000000-0000-0000-0000-000000000001'), 'out', 'use', 20, 'kg', '2026-01-18 11:35:00+00', (SELECT id FROM inventory_usage WHERE uuid = '83000000-0000-0000-0000-000000000010'));

COMMIT;
