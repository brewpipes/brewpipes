-- Add name column to recipe_ingredient
ALTER TABLE recipe_ingredient ADD COLUMN name text NOT NULL DEFAULT '';

-- Backfill names from notes for existing seed data
-- These are the human-readable names that were previously only in the notes field
UPDATE recipe_ingredient SET name = 'Pale Malt 2-Row' WHERE uuid = 'b0000000-0000-0000-0000-000000000001';
UPDATE recipe_ingredient SET name = 'Crystal 40L' WHERE uuid = 'b0000000-0000-0000-0000-000000000002';
UPDATE recipe_ingredient SET name = 'Munich Malt' WHERE uuid = 'b0000000-0000-0000-0000-000000000003';
UPDATE recipe_ingredient SET name = 'Columbus' WHERE uuid = 'b0000000-0000-0000-0000-000000000004';
UPDATE recipe_ingredient SET name = 'Citra' WHERE uuid = 'b0000000-0000-0000-0000-000000000005';
UPDATE recipe_ingredient SET name = 'Mosaic' WHERE uuid = 'b0000000-0000-0000-0000-000000000006';
UPDATE recipe_ingredient SET name = 'Citra' WHERE uuid = 'b0000000-0000-0000-0000-000000000007';
UPDATE recipe_ingredient SET name = 'Mosaic' WHERE uuid = 'b0000000-0000-0000-0000-000000000008';
UPDATE recipe_ingredient SET name = 'California Ale Yeast (WLP001)' WHERE uuid = 'b0000000-0000-0000-0000-000000000009';
UPDATE recipe_ingredient SET name = 'Gypsum' WHERE uuid = 'b0000000-0000-0000-0000-000000000010';
UPDATE recipe_ingredient SET name = 'Pilsner Malt' WHERE uuid = 'b0000000-0000-0000-0000-000000000011';
UPDATE recipe_ingredient SET name = 'Hallertau Mittelfrüh' WHERE uuid = 'b0000000-0000-0000-0000-000000000012';
UPDATE recipe_ingredient SET name = 'Saaz' WHERE uuid = 'b0000000-0000-0000-0000-000000000013';
UPDATE recipe_ingredient SET name = 'Lager Yeast (W-34/70)' WHERE uuid = 'b0000000-0000-0000-0000-000000000014';
UPDATE recipe_ingredient SET name = 'Pale Malt 2-Row' WHERE uuid = 'b0000000-0000-0000-0000-000000000015';
UPDATE recipe_ingredient SET name = 'Chocolate Malt' WHERE uuid = 'b0000000-0000-0000-0000-000000000016';
UPDATE recipe_ingredient SET name = 'Roasted Barley' WHERE uuid = 'b0000000-0000-0000-0000-000000000017';
UPDATE recipe_ingredient SET name = 'Flaked Oats' WHERE uuid = 'b0000000-0000-0000-0000-000000000018';
UPDATE recipe_ingredient SET name = 'East Kent Goldings' WHERE uuid = 'b0000000-0000-0000-0000-000000000019';
UPDATE recipe_ingredient SET name = 'Irish Ale Yeast (WLP004)' WHERE uuid = 'b0000000-0000-0000-0000-000000000020';
UPDATE recipe_ingredient SET name = 'Pilsner Malt' WHERE uuid = 'b0000000-0000-0000-0000-000000000021';
UPDATE recipe_ingredient SET name = 'Wheat Malt' WHERE uuid = 'b0000000-0000-0000-0000-000000000022';
UPDATE recipe_ingredient SET name = 'Hallertau Mittelfrüh' WHERE uuid = 'b0000000-0000-0000-0000-000000000023';
UPDATE recipe_ingredient SET name = 'Tettnang' WHERE uuid = 'b0000000-0000-0000-0000-000000000024';
UPDATE recipe_ingredient SET name = 'Kölsch Yeast (WLP029)' WHERE uuid = 'b0000000-0000-0000-0000-000000000025';
UPDATE recipe_ingredient SET name = 'Pilsner Malt' WHERE uuid = 'b0000000-0000-0000-0000-000000000026';
UPDATE recipe_ingredient SET name = 'Wheat Malt' WHERE uuid = 'b0000000-0000-0000-0000-000000000027';
UPDATE recipe_ingredient SET name = 'Munich Malt' WHERE uuid = 'b0000000-0000-0000-0000-000000000028';
UPDATE recipe_ingredient SET name = 'Saaz' WHERE uuid = 'b0000000-0000-0000-0000-000000000029';
UPDATE recipe_ingredient SET name = 'Styrian Goldings' WHERE uuid = 'b0000000-0000-0000-0000-000000000030';
UPDATE recipe_ingredient SET name = 'Grains of Paradise' WHERE uuid = 'b0000000-0000-0000-0000-000000000031';
UPDATE recipe_ingredient SET name = 'French Saison Yeast (WY3711)' WHERE uuid = 'b0000000-0000-0000-0000-000000000032';
UPDATE recipe_ingredient SET name = 'Coriander' WHERE uuid = 'b0000000-0000-0000-0000-000000000033';
