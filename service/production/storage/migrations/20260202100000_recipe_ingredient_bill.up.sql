-- Recipe Ingredient Bill: Add target specifications to recipe and create recipe_ingredient table.
-- This migration extends recipes with brewing parameters and ingredient bills for formulation.

-- ==============================================================================
-- PART 1: Add target specifications to recipe table
-- ==============================================================================

-- Batch size reference for scaling calculations
ALTER TABLE recipe ADD COLUMN batch_size numeric(10,2);
ALTER TABLE recipe ADD COLUMN batch_size_unit varchar(7);

-- Original gravity targets (specific gravity, e.g., 1.050)
ALTER TABLE recipe ADD COLUMN target_og numeric(5,4);
ALTER TABLE recipe ADD COLUMN target_og_min numeric(5,4);
ALTER TABLE recipe ADD COLUMN target_og_max numeric(5,4);

-- Final gravity targets
ALTER TABLE recipe ADD COLUMN target_fg numeric(5,4);
ALTER TABLE recipe ADD COLUMN target_fg_min numeric(5,4);
ALTER TABLE recipe ADD COLUMN target_fg_max numeric(5,4);

-- Bitterness targets (IBU)
ALTER TABLE recipe ADD COLUMN target_ibu numeric(5,1);
ALTER TABLE recipe ADD COLUMN target_ibu_min numeric(5,1);
ALTER TABLE recipe ADD COLUMN target_ibu_max numeric(5,1);

-- Color targets (SRM)
ALTER TABLE recipe ADD COLUMN target_srm numeric(5,1);
ALTER TABLE recipe ADD COLUMN target_srm_min numeric(5,1);
ALTER TABLE recipe ADD COLUMN target_srm_max numeric(5,1);

-- Carbonation target (volumes CO2)
ALTER TABLE recipe ADD COLUMN target_carbonation numeric(4,2);

-- IBU calculation method
ALTER TABLE recipe ADD COLUMN ibu_method varchar(32);

-- Brewhouse efficiency (percentage, e.g., 75.0 for 75%)
ALTER TABLE recipe ADD COLUMN brewhouse_efficiency numeric(5,2);

-- Add constraints for valid ranges and values
ALTER TABLE recipe ADD CONSTRAINT recipe_batch_size_check 
    CHECK (batch_size IS NULL OR batch_size > 0);

ALTER TABLE recipe ADD CONSTRAINT recipe_og_range_check 
    CHECK (
        (target_og_min IS NULL AND target_og_max IS NULL) OR
        (target_og_min IS NULL OR target_og_max IS NULL) OR
        (target_og_min <= target_og_max)
    );

ALTER TABLE recipe ADD CONSTRAINT recipe_fg_range_check 
    CHECK (
        (target_fg_min IS NULL AND target_fg_max IS NULL) OR
        (target_fg_min IS NULL OR target_fg_max IS NULL) OR
        (target_fg_min <= target_fg_max)
    );

ALTER TABLE recipe ADD CONSTRAINT recipe_ibu_range_check 
    CHECK (
        (target_ibu_min IS NULL AND target_ibu_max IS NULL) OR
        (target_ibu_min IS NULL OR target_ibu_max IS NULL) OR
        (target_ibu_min <= target_ibu_max)
    );

ALTER TABLE recipe ADD CONSTRAINT recipe_srm_range_check 
    CHECK (
        (target_srm_min IS NULL AND target_srm_max IS NULL) OR
        (target_srm_min IS NULL OR target_srm_max IS NULL) OR
        (target_srm_min <= target_srm_max)
    );

ALTER TABLE recipe ADD CONSTRAINT recipe_ibu_method_check 
    CHECK (ibu_method IS NULL OR ibu_method IN (
        'tinseth',   -- Most common, accounts for gravity and time
        'rager',     -- Slightly higher estimates than Tinseth
        'garetz',    -- Accounts for more factors, lower estimates
        'daniels'    -- Simple percentage-based
    ));

ALTER TABLE recipe ADD CONSTRAINT recipe_efficiency_check 
    CHECK (brewhouse_efficiency IS NULL OR (brewhouse_efficiency >= 0 AND brewhouse_efficiency <= 100));

ALTER TABLE recipe ADD CONSTRAINT recipe_carbonation_check 
    CHECK (target_carbonation IS NULL OR (target_carbonation >= 0 AND target_carbonation <= 5));

-- ==============================================================================
-- PART 2: Create recipe_ingredient table
-- ==============================================================================

CREATE TABLE IF NOT EXISTS recipe_ingredient (
    id                      serial PRIMARY KEY,
    uuid                    uuid NOT NULL DEFAULT gen_random_uuid(),

    -- Recipe reference (required)
    recipe_id               int NOT NULL REFERENCES recipe(id) ON DELETE CASCADE,

    -- Cross-service reference to inventory.ingredient (no FK across service boundaries)
    ingredient_uuid         uuid,

    -- Ingredient classification
    ingredient_type         varchar(32) NOT NULL,

    -- Amount specification
    amount                  numeric(10,4) NOT NULL,
    amount_unit             varchar(16) NOT NULL,

    -- Usage context
    use_stage               varchar(32) NOT NULL,
    use_type                varchar(32),

    -- Timing parameters
    timing_duration_minutes int,
    timing_temperature_c    numeric(5,1),

    -- Hop-specific: assumed alpha acid percentage for IBU calculations
    alpha_acid_assumed      numeric(4,2),

    -- Scaling behavior (1.0 = linear scaling, <1.0 = sublinear for things like hops)
    scaling_factor          numeric(4,2) NOT NULL DEFAULT 1.0,

    -- Display ordering within recipe
    sort_order              int NOT NULL DEFAULT 0,

    -- Free-form notes
    notes                   text,

    created_at              timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at              timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at              timestamptz,

    -- Ingredient type must be valid
    CONSTRAINT recipe_ingredient_type_check CHECK (ingredient_type IN (
        'fermentable',  -- Malt, grain, extract, sugar
        'hop',          -- Hops (pellet, whole, extract)
        'yeast',        -- Yeast and bacteria
        'adjunct',      -- Fruit, spices, flavorings
        'salt',         -- Brewing salts (gypsum, calcium chloride, etc.)
        'chemical',     -- Acids, finings, enzymes
        'gas',          -- CO2, nitrogen
        'other'         -- Anything else
    )),

    -- Use stage must be valid
    CONSTRAINT recipe_ingredient_use_stage_check CHECK (use_stage IN (
        'mash',         -- Added during mashing
        'boil',         -- Added during the boil
        'whirlpool',    -- Added during whirlpool/hopstand
        'fermentation', -- Added during fermentation
        'packaging'     -- Added at packaging (priming, carbonation)
    )),

    -- Use type validation (context-dependent, so we allow common values)
    CONSTRAINT recipe_ingredient_use_type_check CHECK (use_type IS NULL OR use_type IN (
        -- Hop use types
        'bittering',    -- Early boil for bitterness
        'flavor',       -- Mid-boil for flavor
        'aroma',        -- Late boil/whirlpool for aroma
        'dry_hop',      -- Post-fermentation dry hopping
        -- Fermentable use types
        'base',         -- Base malt (majority of grist)
        'specialty',    -- Specialty/character malts
        'adjunct',      -- Non-malt fermentables
        'sugar',        -- Simple sugars
        -- Yeast use types
        'primary',      -- Primary fermentation
        'secondary',    -- Secondary fermentation
        'bottle',       -- Bottle conditioning
        -- General
        'other'
    )),

    -- Alpha acid only makes sense for hops
    CONSTRAINT recipe_ingredient_alpha_acid_check CHECK (
        alpha_acid_assumed IS NULL OR 
        (ingredient_type = 'hop' AND alpha_acid_assumed >= 0 AND alpha_acid_assumed <= 30)
    ),

    -- Scaling factor must be positive
    CONSTRAINT recipe_ingredient_scaling_factor_check CHECK (scaling_factor > 0),

    -- Amount must be positive
    CONSTRAINT recipe_ingredient_amount_check CHECK (amount > 0)
);

-- Indexes for recipe_ingredient
CREATE UNIQUE INDEX IF NOT EXISTS recipe_ingredient_uuid_idx ON recipe_ingredient(uuid);
CREATE INDEX IF NOT EXISTS recipe_ingredient_recipe_id_idx ON recipe_ingredient(recipe_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS recipe_ingredient_ingredient_uuid_idx ON recipe_ingredient(ingredient_uuid) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS recipe_ingredient_type_idx ON recipe_ingredient(ingredient_type) WHERE deleted_at IS NULL;

-- ==============================================================================
-- PART 3: Seed data for existing recipes
-- ==============================================================================

-- Update Flagship IPA recipe with target specs
UPDATE recipe SET
    batch_size = 1200,
    batch_size_unit = 'l',
    target_og = 1.0650,
    target_og_min = 1.0620,
    target_og_max = 1.0680,
    target_fg = 1.0120,
    target_fg_min = 1.0100,
    target_fg_max = 1.0140,
    target_ibu = 65,
    target_ibu_min = 60,
    target_ibu_max = 70,
    target_srm = 8,
    target_srm_min = 6,
    target_srm_max = 10,
    target_carbonation = 2.4,
    ibu_method = 'tinseth',
    brewhouse_efficiency = 75.0
WHERE uuid = 'a1000000-0000-0000-0000-000000000001';

-- Update Crisp Pilsner recipe with target specs
UPDATE recipe SET
    batch_size = 800,
    batch_size_unit = 'l',
    target_og = 1.0480,
    target_og_min = 1.0460,
    target_og_max = 1.0500,
    target_fg = 1.0080,
    target_fg_min = 1.0060,
    target_fg_max = 1.0100,
    target_ibu = 35,
    target_ibu_min = 30,
    target_ibu_max = 40,
    target_srm = 3,
    target_srm_min = 2,
    target_srm_max = 4,
    target_carbonation = 2.6,
    ibu_method = 'tinseth',
    brewhouse_efficiency = 78.0
WHERE uuid = 'a1000000-0000-0000-0000-000000000002';

-- Update Robust Stout recipe with target specs
UPDATE recipe SET
    batch_size = 900,
    batch_size_unit = 'l',
    target_og = 1.0600,
    target_og_min = 1.0580,
    target_og_max = 1.0620,
    target_fg = 1.0140,
    target_fg_min = 1.0120,
    target_fg_max = 1.0160,
    target_ibu = 40,
    target_ibu_min = 35,
    target_ibu_max = 45,
    target_srm = 35,
    target_srm_min = 30,
    target_srm_max = 40,
    target_carbonation = 2.0,
    ibu_method = 'tinseth',
    brewhouse_efficiency = 72.0
WHERE uuid = 'a1000000-0000-0000-0000-000000000003';

-- Update Cologne Classic (Kolsch) recipe with target specs
UPDATE recipe SET
    batch_size = 1000,
    batch_size_unit = 'l',
    target_og = 1.0480,
    target_og_min = 1.0460,
    target_og_max = 1.0500,
    target_fg = 1.0080,
    target_fg_min = 1.0060,
    target_fg_max = 1.0100,
    target_ibu = 28,
    target_ibu_min = 22,
    target_ibu_max = 32,
    target_srm = 4,
    target_srm_min = 3,
    target_srm_max = 5,
    target_carbonation = 2.5,
    ibu_method = 'tinseth',
    brewhouse_efficiency = 76.0
WHERE uuid = 'a1000000-0000-0000-0000-000000000004';

-- Update Farmhouse Saison recipe with target specs
UPDATE recipe SET
    batch_size = 950,
    batch_size_unit = 'l',
    target_og = 1.0560,
    target_og_min = 1.0520,
    target_og_max = 1.0600,
    target_fg = 1.0040,
    target_fg_min = 1.0020,
    target_fg_max = 1.0060,
    target_ibu = 30,
    target_ibu_min = 25,
    target_ibu_max = 35,
    target_srm = 5,
    target_srm_min = 4,
    target_srm_max = 7,
    target_carbonation = 2.8,
    ibu_method = 'tinseth',
    brewhouse_efficiency = 74.0
WHERE uuid = 'a1000000-0000-0000-0000-000000000005';

-- ==============================================================================
-- Recipe Ingredient Seed Data
-- ==============================================================================

-- Flagship IPA ingredients
INSERT INTO recipe_ingredient (uuid, recipe_id, ingredient_uuid, ingredient_type, amount, amount_unit, use_stage, use_type, timing_duration_minutes, alpha_acid_assumed, sort_order, notes)
VALUES
    -- Fermentables
    ('b0000000-0000-0000-0000-000000000001', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000001'), '70000000-0000-0000-0000-000000000001', 'fermentable', 200, 'kg', 'mash', 'base', NULL, NULL, 1, 'Pale malt base'),
    ('b0000000-0000-0000-0000-000000000002', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000001'), '70000000-0000-0000-0000-000000000002', 'fermentable', 15, 'kg', 'mash', 'specialty', NULL, NULL, 2, 'Crystal 40L for body'),
    ('b0000000-0000-0000-0000-000000000003', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000001'), '70000000-0000-0000-0000-000000000003', 'fermentable', 5, 'kg', 'mash', 'specialty', NULL, NULL, 3, 'Munich for complexity'),
    -- Hops
    ('b0000000-0000-0000-0000-000000000004', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000001'), '70000000-0000-0000-0000-000000000010', 'hop', 2.5, 'kg', 'boil', 'bittering', 60, 12.5, 10, 'Columbus for bittering'),
    ('b0000000-0000-0000-0000-000000000005', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000001'), '70000000-0000-0000-0000-000000000011', 'hop', 1.5, 'kg', 'boil', 'flavor', 15, 13.0, 11, 'Citra for flavor'),
    ('b0000000-0000-0000-0000-000000000006', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000001'), '70000000-0000-0000-0000-000000000012', 'hop', 2.0, 'kg', 'whirlpool', 'aroma', 20, 12.0, 12, 'Mosaic whirlpool'),
    ('b0000000-0000-0000-0000-000000000007', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000001'), '70000000-0000-0000-0000-000000000011', 'hop', 3.0, 'kg', 'fermentation', 'dry_hop', NULL, 13.0, 13, 'Citra dry hop'),
    ('b0000000-0000-0000-0000-000000000008', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000001'), '70000000-0000-0000-0000-000000000012', 'hop', 2.0, 'kg', 'fermentation', 'dry_hop', NULL, 12.0, 14, 'Mosaic dry hop'),
    -- Yeast
    ('b0000000-0000-0000-0000-000000000009', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000001'), '70000000-0000-0000-0000-000000000020', 'yeast', 2, 'kg', 'fermentation', 'primary', NULL, NULL, 20, 'WLP001 California Ale'),
    -- Water chemistry
    ('b0000000-0000-0000-0000-000000000010', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000001'), NULL, 'salt', 500, 'g', 'mash', NULL, NULL, NULL, 30, 'Gypsum for sulfate');

-- Crisp Pilsner ingredients
INSERT INTO recipe_ingredient (uuid, recipe_id, ingredient_uuid, ingredient_type, amount, amount_unit, use_stage, use_type, timing_duration_minutes, alpha_acid_assumed, sort_order, notes)
VALUES
    -- Fermentables
    ('b0000000-0000-0000-0000-000000000011', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000002'), '70000000-0000-0000-0000-000000000004', 'fermentable', 150, 'kg', 'mash', 'base', NULL, NULL, 1, 'Pilsner malt'),
    -- Hops
    ('b0000000-0000-0000-0000-000000000012', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000002'), '70000000-0000-0000-0000-000000000013', 'hop', 1.5, 'kg', 'boil', 'bittering', 60, 4.0, 10, 'Saaz bittering'),
    ('b0000000-0000-0000-0000-000000000013', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000002'), '70000000-0000-0000-0000-000000000013', 'hop', 1.0, 'kg', 'boil', 'aroma', 5, 4.0, 11, 'Saaz late addition'),
    -- Yeast
    ('b0000000-0000-0000-0000-000000000014', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000002'), '70000000-0000-0000-0000-000000000021', 'yeast', 3, 'kg', 'fermentation', 'primary', NULL, NULL, 20, 'W34/70 lager yeast');

-- Robust Stout ingredients
INSERT INTO recipe_ingredient (uuid, recipe_id, ingredient_uuid, ingredient_type, amount, amount_unit, use_stage, use_type, timing_duration_minutes, alpha_acid_assumed, sort_order, notes)
VALUES
    -- Fermentables
    ('b0000000-0000-0000-0000-000000000015', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000003'), '70000000-0000-0000-0000-000000000001', 'fermentable', 160, 'kg', 'mash', 'base', NULL, NULL, 1, 'Pale malt base'),
    ('b0000000-0000-0000-0000-000000000016', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000003'), '70000000-0000-0000-0000-000000000005', 'fermentable', 25, 'kg', 'mash', 'specialty', NULL, NULL, 2, 'Roasted barley'),
    ('b0000000-0000-0000-0000-000000000017', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000003'), '70000000-0000-0000-0000-000000000006', 'fermentable', 15, 'kg', 'mash', 'specialty', NULL, NULL, 3, 'Chocolate malt'),
    ('b0000000-0000-0000-0000-000000000018', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000003'), '70000000-0000-0000-0000-000000000007', 'fermentable', 10, 'kg', 'mash', 'adjunct', NULL, NULL, 4, 'Flaked barley for head'),
    -- Hops
    ('b0000000-0000-0000-0000-000000000019', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000003'), '70000000-0000-0000-0000-000000000014', 'hop', 2.0, 'kg', 'boil', 'bittering', 60, 5.5, 10, 'East Kent Goldings'),
    -- Yeast
    ('b0000000-0000-0000-0000-000000000020', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000003'), '70000000-0000-0000-0000-000000000022', 'yeast', 2.5, 'kg', 'fermentation', 'primary', NULL, NULL, 20, 'Irish ale yeast');

-- Cologne Classic (Kolsch) ingredients
INSERT INTO recipe_ingredient (uuid, recipe_id, ingredient_uuid, ingredient_type, amount, amount_unit, use_stage, use_type, timing_duration_minutes, alpha_acid_assumed, sort_order, notes)
VALUES
    -- Fermentables
    ('b0000000-0000-0000-0000-000000000021', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000004'), '70000000-0000-0000-0000-000000000004', 'fermentable', 170, 'kg', 'mash', 'base', NULL, NULL, 1, 'Pilsner malt'),
    ('b0000000-0000-0000-0000-000000000022', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000004'), '70000000-0000-0000-0000-000000000008', 'fermentable', 10, 'kg', 'mash', 'specialty', NULL, NULL, 2, 'Vienna malt'),
    -- Hops
    ('b0000000-0000-0000-0000-000000000023', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000004'), '70000000-0000-0000-0000-000000000015', 'hop', 1.8, 'kg', 'boil', 'bittering', 60, 4.5, 10, 'Hallertau Mittelfruh'),
    ('b0000000-0000-0000-0000-000000000024', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000004'), '70000000-0000-0000-0000-000000000015', 'hop', 0.8, 'kg', 'boil', 'aroma', 5, 4.5, 11, 'Hallertau late'),
    -- Yeast
    ('b0000000-0000-0000-0000-000000000025', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000004'), '70000000-0000-0000-0000-000000000023', 'yeast', 2.5, 'kg', 'fermentation', 'primary', NULL, NULL, 20, 'Kolsch yeast');

-- Farmhouse Saison ingredients
INSERT INTO recipe_ingredient (uuid, recipe_id, ingredient_uuid, ingredient_type, amount, amount_unit, use_stage, use_type, timing_duration_minutes, alpha_acid_assumed, sort_order, notes)
VALUES
    -- Fermentables
    ('b0000000-0000-0000-0000-000000000026', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000005'), '70000000-0000-0000-0000-000000000004', 'fermentable', 150, 'kg', 'mash', 'base', NULL, NULL, 1, 'Pilsner malt'),
    ('b0000000-0000-0000-0000-000000000027', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000005'), '70000000-0000-0000-0000-000000000009', 'fermentable', 15, 'kg', 'mash', 'specialty', NULL, NULL, 2, 'Wheat malt'),
    ('b0000000-0000-0000-0000-000000000028', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000005'), '70000000-0000-0000-0000-000000000003', 'fermentable', 10, 'kg', 'mash', 'specialty', NULL, NULL, 3, 'Munich malt'),
    -- Hops
    ('b0000000-0000-0000-0000-000000000029', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000005'), '70000000-0000-0000-0000-000000000016', 'hop', 1.5, 'kg', 'boil', 'bittering', 60, 8.0, 10, 'Styrian Goldings'),
    ('b0000000-0000-0000-0000-000000000030', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000005'), '70000000-0000-0000-0000-000000000013', 'hop', 0.5, 'kg', 'boil', 'aroma', 5, 4.0, 11, 'Saaz late'),
    -- Yeast
    ('b0000000-0000-0000-0000-000000000031', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000005'), '70000000-0000-0000-0000-000000000024', 'yeast', 2, 'kg', 'fermentation', 'primary', NULL, NULL, 20, 'Belgian saison yeast'),
    -- Adjuncts
    ('b0000000-0000-0000-0000-000000000032', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000005'), '70000000-0000-0000-0000-000000000030', 'adjunct', 0.5, 'kg', 'boil', NULL, 5, NULL, 25, 'Coriander seeds'),
    ('b0000000-0000-0000-0000-000000000033', (SELECT id FROM recipe WHERE uuid = 'a1000000-0000-0000-0000-000000000005'), '70000000-0000-0000-0000-000000000031', 'adjunct', 0.1, 'kg', 'boil', NULL, 5, NULL, 26, 'Orange peel');
