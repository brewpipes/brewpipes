-- Reverse Enhanced Batch Tracking changes

-- Remove seed data added by this migration

-- Remove transfers added by this migration
DELETE FROM transfer WHERE uuid IN (
    '93200000-0000-0000-0000-000000000004',
    '93200000-0000-0000-0000-000000000005',
    '93200000-0000-0000-0000-000000000006',
    '93200000-0000-0000-0000-000000000007',
    '93200000-0000-0000-0000-000000000008',
    '93200000-0000-0000-0000-000000000009',
    '93200000-0000-0000-0000-000000000010',
    '93200000-0000-0000-0000-000000000011',
    '93200000-0000-0000-0000-000000000012',
    '93200000-0000-0000-0000-000000000013'
);

-- Remove volume relations added by this migration
DELETE FROM volume_relation WHERE uuid IN (
    '92000000-0000-0000-0000-000000000007',
    '92000000-0000-0000-0000-000000000008'
);

-- Remove occupancies added by this migration
DELETE FROM occupancy WHERE uuid IN (
    '93100000-0000-0000-0000-000000000009',
    '93100000-0000-0000-0000-000000000010',
    '93100000-0000-0000-0000-000000000011',
    '93100000-0000-0000-0000-000000000012',
    '93100000-0000-0000-0000-000000000013',
    '93100000-0000-0000-0000-000000000014',
    '93100000-0000-0000-0000-000000000015',
    '93100000-0000-0000-0000-000000000016',
    '93100000-0000-0000-0000-000000000017'
);

-- Remove volume-targeted measurements
DELETE FROM measurement WHERE uuid IN (
    '98000000-0000-0000-0000-000000000014',
    '98000000-0000-0000-0000-000000000015',
    '98000000-0000-0000-0000-000000000016',
    '98000000-0000-0000-0000-000000000017'
);

-- Remove volume-targeted additions
DELETE FROM addition WHERE uuid IN (
    '97000000-0000-0000-0000-000000000011',
    '97000000-0000-0000-0000-000000000012'
);

-- Remove brew sessions
DELETE FROM brew_session WHERE uuid IN (
    'a2000000-0000-0000-0000-000000000001',
    'a2000000-0000-0000-0000-000000000002',
    'a2000000-0000-0000-0000-000000000003',
    'a2000000-0000-0000-0000-000000000004',
    'a2000000-0000-0000-0000-000000000005'
);

-- Clear occupancy statuses
UPDATE occupancy SET status = NULL;

-- Clear recipe_id from batches
UPDATE batch SET recipe_id = NULL;

-- Remove recipes
DELETE FROM recipe WHERE uuid IN (
    'a1000000-0000-0000-0000-000000000001',
    'a1000000-0000-0000-0000-000000000002',
    'a1000000-0000-0000-0000-000000000003',
    'a1000000-0000-0000-0000-000000000004',
    'a1000000-0000-0000-0000-000000000005'
);

-- Remove styles
DELETE FROM style WHERE uuid IN (
    'a0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000002',
    'a0000000-0000-0000-0000-000000000003',
    'a0000000-0000-0000-0000-000000000004',
    'a0000000-0000-0000-0000-000000000005'
);

-- Revert measurement constraint and remove volume_id
ALTER TABLE measurement DROP CONSTRAINT measurement_target_check;
ALTER TABLE measurement ADD CONSTRAINT measurement_target_check CHECK (
    (batch_id IS NOT NULL AND occupancy_id IS NULL) OR
    (batch_id IS NULL AND occupancy_id IS NOT NULL)
);
DROP INDEX IF EXISTS measurement_volume_id_idx;
ALTER TABLE measurement DROP COLUMN volume_id;

-- Revert addition constraint and remove volume_id
ALTER TABLE addition DROP CONSTRAINT addition_target_check;
ALTER TABLE addition ADD CONSTRAINT addition_target_check CHECK (
    (batch_id IS NOT NULL AND occupancy_id IS NULL) OR
    (batch_id IS NULL AND occupancy_id IS NOT NULL)
);
DROP INDEX IF EXISTS addition_volume_id_idx;
ALTER TABLE addition DROP COLUMN volume_id;

-- Remove status from occupancy
ALTER TABLE occupancy DROP CONSTRAINT IF EXISTS occupancy_status_check;
ALTER TABLE occupancy DROP COLUMN status;

-- Remove recipe_id from batch
DROP INDEX IF EXISTS batch_recipe_id_idx;
ALTER TABLE batch DROP COLUMN recipe_id;

-- Drop brew_session table
DROP INDEX IF EXISTS brew_session_wort_volume_id_idx;
DROP INDEX IF EXISTS brew_session_batch_id_idx;
DROP INDEX IF EXISTS brew_session_uuid_idx;
DROP TABLE IF EXISTS brew_session;

-- Drop recipe table
DROP INDEX IF EXISTS recipe_style_id_idx;
DROP INDEX IF EXISTS recipe_uuid_idx;
DROP TABLE IF EXISTS recipe;

-- Drop style table
DROP INDEX IF EXISTS style_name_lower_idx;
DROP INDEX IF EXISTS style_uuid_idx;
DROP TABLE IF EXISTS style;
