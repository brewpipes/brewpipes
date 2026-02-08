-- Reverse Recipe Ingredient Bill migration

-- ==============================================================================
-- PART 1: Remove seed data
-- ==============================================================================

-- Remove recipe ingredient seed data
DELETE FROM recipe_ingredient WHERE uuid IN (
    'b0000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000002',
    'b0000000-0000-0000-0000-000000000003',
    'b0000000-0000-0000-0000-000000000004',
    'b0000000-0000-0000-0000-000000000005',
    'b0000000-0000-0000-0000-000000000006',
    'b0000000-0000-0000-0000-000000000007',
    'b0000000-0000-0000-0000-000000000008',
    'b0000000-0000-0000-0000-000000000009',
    'b0000000-0000-0000-0000-000000000010',
    'b0000000-0000-0000-0000-000000000011',
    'b0000000-0000-0000-0000-000000000012',
    'b0000000-0000-0000-0000-000000000013',
    'b0000000-0000-0000-0000-000000000014',
    'b0000000-0000-0000-0000-000000000015',
    'b0000000-0000-0000-0000-000000000016',
    'b0000000-0000-0000-0000-000000000017',
    'b0000000-0000-0000-0000-000000000018',
    'b0000000-0000-0000-0000-000000000019',
    'b0000000-0000-0000-0000-000000000020',
    'b0000000-0000-0000-0000-000000000021',
    'b0000000-0000-0000-0000-000000000022',
    'b0000000-0000-0000-0000-000000000023',
    'b0000000-0000-0000-0000-000000000024',
    'b0000000-0000-0000-0000-000000000025',
    'b0000000-0000-0000-0000-000000000026',
    'b0000000-0000-0000-0000-000000000027',
    'b0000000-0000-0000-0000-000000000028',
    'b0000000-0000-0000-0000-000000000029',
    'b0000000-0000-0000-0000-000000000030',
    'b0000000-0000-0000-0000-000000000031',
    'b0000000-0000-0000-0000-000000000032',
    'b0000000-0000-0000-0000-000000000033'
);

-- Clear target specs from recipes
UPDATE recipe SET
    batch_size = NULL,
    batch_size_unit = NULL,
    target_og = NULL,
    target_og_min = NULL,
    target_og_max = NULL,
    target_fg = NULL,
    target_fg_min = NULL,
    target_fg_max = NULL,
    target_ibu = NULL,
    target_ibu_min = NULL,
    target_ibu_max = NULL,
    target_srm = NULL,
    target_srm_min = NULL,
    target_srm_max = NULL,
    target_carbonation = NULL,
    ibu_method = NULL,
    brewhouse_efficiency = NULL
WHERE uuid IN (
    'a1000000-0000-0000-0000-000000000001',
    'a1000000-0000-0000-0000-000000000002',
    'a1000000-0000-0000-0000-000000000003',
    'a1000000-0000-0000-0000-000000000004',
    'a1000000-0000-0000-0000-000000000005'
);

-- ==============================================================================
-- PART 2: Drop recipe_ingredient table
-- ==============================================================================

DROP INDEX IF EXISTS recipe_ingredient_type_idx;
DROP INDEX IF EXISTS recipe_ingredient_ingredient_uuid_idx;
DROP INDEX IF EXISTS recipe_ingredient_recipe_id_idx;
DROP INDEX IF EXISTS recipe_ingredient_uuid_idx;
DROP TABLE IF EXISTS recipe_ingredient;

-- ==============================================================================
-- PART 3: Remove target specification columns from recipe
-- ==============================================================================

-- Drop constraints first
ALTER TABLE recipe DROP CONSTRAINT IF EXISTS recipe_carbonation_check;
ALTER TABLE recipe DROP CONSTRAINT IF EXISTS recipe_efficiency_check;
ALTER TABLE recipe DROP CONSTRAINT IF EXISTS recipe_ibu_method_check;
ALTER TABLE recipe DROP CONSTRAINT IF EXISTS recipe_srm_range_check;
ALTER TABLE recipe DROP CONSTRAINT IF EXISTS recipe_ibu_range_check;
ALTER TABLE recipe DROP CONSTRAINT IF EXISTS recipe_fg_range_check;
ALTER TABLE recipe DROP CONSTRAINT IF EXISTS recipe_og_range_check;
ALTER TABLE recipe DROP CONSTRAINT IF EXISTS recipe_batch_size_check;

-- Drop columns
ALTER TABLE recipe DROP COLUMN IF EXISTS brewhouse_efficiency;
ALTER TABLE recipe DROP COLUMN IF EXISTS ibu_method;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_carbonation;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_srm_max;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_srm_min;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_srm;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_ibu_max;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_ibu_min;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_ibu;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_fg_max;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_fg_min;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_fg;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_og_max;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_og_min;
ALTER TABLE recipe DROP COLUMN IF EXISTS target_og;
ALTER TABLE recipe DROP COLUMN IF EXISTS batch_size_unit;
ALTER TABLE recipe DROP COLUMN IF EXISTS batch_size;
