-- Reverse Enhanced Batch Tracking changes

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
