-- Reverse packaging run tables and seed data.
BEGIN;

-- Remove seed data
DELETE FROM package_format WHERE uuid IN (
    'b0000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000002',
    'b0000000-0000-0000-0000-000000000003',
    'b0000000-0000-0000-0000-000000000004',
    'b0000000-0000-0000-0000-000000000005',
    'b0000000-0000-0000-0000-000000000006',
    'b0000000-0000-0000-0000-000000000007'
);

-- Drop packaging_run_line
DROP INDEX IF EXISTS packaging_run_line_run_format_idx;
DROP INDEX IF EXISTS packaging_run_line_package_format_id_idx;
DROP INDEX IF EXISTS packaging_run_line_packaging_run_id_idx;
DROP INDEX IF EXISTS packaging_run_line_uuid_idx;
DROP TABLE IF EXISTS packaging_run_line;

-- Drop packaging_run
DROP INDEX IF EXISTS packaging_run_occupancy_id_idx;
DROP INDEX IF EXISTS packaging_run_batch_id_idx;
DROP INDEX IF EXISTS packaging_run_uuid_idx;
DROP TABLE IF EXISTS packaging_run;

-- Drop package_format
DROP INDEX IF EXISTS package_format_name_lower_idx;
DROP INDEX IF EXISTS package_format_uuid_idx;
DROP TABLE IF EXISTS package_format;

COMMIT;
