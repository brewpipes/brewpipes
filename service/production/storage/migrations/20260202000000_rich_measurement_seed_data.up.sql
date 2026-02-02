-- Rich measurement seed data for sparkline visualization
-- Adds time-series measurements for temperature, gravity, and pH across multiple batches

-- ==============================================================================
-- IPA 24-07 (batch 1) - Full fermentation temperature and gravity curve
-- Fermentation: Jan 10-20, 2026
-- ==============================================================================

INSERT INTO measurement (uuid, batch_id, occupancy_id, kind, value, unit, observed_at, notes)
VALUES
    -- Temperature readings during fermentation (twice daily)
    ('98100000-0000-0000-0000-000000000001', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 18.0, 'c', '2026-01-10 14:00:00+00', 'Pitch temperature.'),
    ('98100000-0000-0000-0000-000000000002', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 18.5, 'c', '2026-01-10 20:00:00+00', 'Evening check.'),
    ('98100000-0000-0000-0000-000000000003', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 19.0, 'c', '2026-01-11 08:00:00+00', 'Morning check - fermentation starting.'),
    ('98100000-0000-0000-0000-000000000004', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 19.5, 'c', '2026-01-11 20:00:00+00', 'Active fermentation.'),
    ('98100000-0000-0000-0000-000000000005', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 20.0, 'c', '2026-01-12 08:00:00+00', 'Peak fermentation temp.'),
    ('98100000-0000-0000-0000-000000000006', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 20.5, 'c', '2026-01-12 20:00:00+00', 'Holding steady.'),
    ('98100000-0000-0000-0000-000000000007', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 20.0, 'c', '2026-01-13 08:00:00+00', 'Slight cooling.'),
    ('98100000-0000-0000-0000-000000000008', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 19.5, 'c', '2026-01-13 20:00:00+00', 'Fermentation slowing.'),
    ('98100000-0000-0000-0000-000000000009', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 19.0, 'c', '2026-01-14 08:00:00+00', 'Ramping down.'),
    ('98100000-0000-0000-0000-000000000010', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 18.5, 'c', '2026-01-14 20:00:00+00', 'Conditioning temp.'),
    ('98100000-0000-0000-0000-000000000011', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 18.0, 'c', '2026-01-15 08:00:00+00', 'Stable.'),
    ('98100000-0000-0000-0000-000000000012', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 17.5, 'c', '2026-01-16 08:00:00+00', 'Cold crash starting.'),
    ('98100000-0000-0000-0000-000000000013', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 10.0, 'c', '2026-01-17 08:00:00+00', 'Cold crash day 1.'),
    ('98100000-0000-0000-0000-000000000014', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 4.0, 'c', '2026-01-18 08:00:00+00', 'Cold crash complete.'),
    ('98100000-0000-0000-0000-000000000015', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 3.5, 'c', '2026-01-19 08:00:00+00', 'Holding cold.'),
    ('98100000-0000-0000-0000-000000000016', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'temperature', 3.0, 'c', '2026-01-20 08:00:00+00', 'Ready for transfer.'),
    
    -- Gravity readings during fermentation (daily)
    ('98100000-0000-0000-0000-000000000017', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'gravity', 1.054, 'sg', '2026-01-10 14:00:00+00', 'OG at pitch.'),
    ('98100000-0000-0000-0000-000000000018', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'gravity', 1.048, 'sg', '2026-01-11 08:00:00+00', 'Day 1.'),
    ('98100000-0000-0000-0000-000000000019', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'gravity', 1.038, 'sg', '2026-01-12 08:00:00+00', 'Day 2 - active fermentation.'),
    ('98100000-0000-0000-0000-000000000020', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'gravity', 1.028, 'sg', '2026-01-13 08:00:00+00', 'Day 3.'),
    ('98100000-0000-0000-0000-000000000021', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'gravity', 1.020, 'sg', '2026-01-14 08:00:00+00', 'Day 4.'),
    ('98100000-0000-0000-0000-000000000022', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'gravity', 1.016, 'sg', '2026-01-15 08:00:00+00', 'Day 5.'),
    ('98100000-0000-0000-0000-000000000023', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'gravity', 1.014, 'sg', '2026-01-16 08:00:00+00', 'Day 6.'),
    ('98100000-0000-0000-0000-000000000024', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'gravity', 1.013, 'sg', '2026-01-17 08:00:00+00', 'Day 7.'),
    ('98100000-0000-0000-0000-000000000025', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'gravity', 1.012, 'sg', '2026-01-18 08:00:00+00', 'FG stable.'),
    ('98100000-0000-0000-0000-000000000026', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'gravity', 1.012, 'sg', '2026-01-19 08:00:00+00', 'FG confirmed.'),
    
    -- pH readings
    ('98100000-0000-0000-0000-000000000027', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'ph', 5.35, NULL, '2026-01-10 09:00:00+00', 'Mash pH.'),
    ('98100000-0000-0000-0000-000000000028', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'ph', 5.20, NULL, '2026-01-10 11:00:00+00', 'Post-boil pH.'),
    ('98100000-0000-0000-0000-000000000029', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'ph', 4.60, NULL, '2026-01-14 08:00:00+00', 'Mid-fermentation pH.'),
    ('98100000-0000-0000-0000-000000000030', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000001'), NULL, 'ph', 4.35, NULL, '2026-01-18 08:00:00+00', 'Final pH.');

-- ==============================================================================
-- Stout 24-09 (batch 4) - Temperature and gravity curve for stout fermentation
-- Fermentation: Jan 22 - Feb 2, 2026
-- ==============================================================================

INSERT INTO measurement (uuid, batch_id, occupancy_id, kind, value, unit, observed_at, notes)
VALUES
    -- Temperature readings (stout ferments warmer, longer)
    ('98100000-0000-0000-0000-000000000031', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 18.0, 'c', '2026-01-22 18:00:00+00', 'Pitch temperature.'),
    ('98100000-0000-0000-0000-000000000032', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 18.5, 'c', '2026-01-23 08:00:00+00', 'Day 1 morning.'),
    ('98100000-0000-0000-0000-000000000033', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 19.0, 'c', '2026-01-23 20:00:00+00', 'Day 1 evening.'),
    ('98100000-0000-0000-0000-000000000034', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 19.5, 'c', '2026-01-24 08:00:00+00', 'Day 2.'),
    ('98100000-0000-0000-0000-000000000035', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 20.0, 'c', '2026-01-24 20:00:00+00', 'Active fermentation.'),
    ('98100000-0000-0000-0000-000000000036', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 20.5, 'c', '2026-01-25 08:00:00+00', 'Peak activity.'),
    ('98100000-0000-0000-0000-000000000037', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 21.0, 'c', '2026-01-25 20:00:00+00', 'Raised for diacetyl rest.'),
    ('98100000-0000-0000-0000-000000000038', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 21.0, 'c', '2026-01-26 08:00:00+00', 'Diacetyl rest.'),
    ('98100000-0000-0000-0000-000000000039', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 20.5, 'c', '2026-01-27 08:00:00+00', 'Cooling down.'),
    ('98100000-0000-0000-0000-000000000040', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 19.0, 'c', '2026-01-28 08:00:00+00', 'Continued cooling.'),
    ('98100000-0000-0000-0000-000000000041', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 15.0, 'c', '2026-01-29 08:00:00+00', 'Cold conditioning start.'),
    ('98100000-0000-0000-0000-000000000042', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 10.0, 'c', '2026-01-30 08:00:00+00', 'Cold conditioning.'),
    ('98100000-0000-0000-0000-000000000043', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 5.0, 'c', '2026-01-31 08:00:00+00', 'Cold crash.'),
    ('98100000-0000-0000-0000-000000000044', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 3.0, 'c', '2026-02-01 08:00:00+00', 'Final cold crash.'),
    ('98100000-0000-0000-0000-000000000045', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'temperature', 3.0, 'c', '2026-02-02 08:00:00+00', 'Ready for packaging.'),
    
    -- Gravity readings (high OG stout)
    ('98100000-0000-0000-0000-000000000046', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'gravity', 1.060, 'sg', '2026-01-22 18:00:00+00', 'OG.'),
    ('98100000-0000-0000-0000-000000000047', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'gravity', 1.055, 'sg', '2026-01-23 08:00:00+00', 'Day 1.'),
    ('98100000-0000-0000-0000-000000000048', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'gravity', 1.045, 'sg', '2026-01-24 08:00:00+00', 'Day 2.'),
    ('98100000-0000-0000-0000-000000000049', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'gravity', 1.035, 'sg', '2026-01-25 08:00:00+00', 'Day 3.'),
    ('98100000-0000-0000-0000-000000000050', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'gravity', 1.028, 'sg', '2026-01-26 08:00:00+00', 'Day 4.'),
    ('98100000-0000-0000-0000-000000000051', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'gravity', 1.022, 'sg', '2026-01-27 08:00:00+00', 'Day 5.'),
    ('98100000-0000-0000-0000-000000000052', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'gravity', 1.018, 'sg', '2026-01-28 08:00:00+00', 'Day 6.'),
    ('98100000-0000-0000-0000-000000000053', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'gravity', 1.016, 'sg', '2026-01-29 08:00:00+00', 'Day 7.'),
    ('98100000-0000-0000-0000-000000000054', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'gravity', 1.015, 'sg', '2026-01-30 08:00:00+00', 'FG approaching.'),
    ('98100000-0000-0000-0000-000000000055', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'gravity', 1.014, 'sg', '2026-02-01 08:00:00+00', 'FG stable.'),
    
    -- pH readings
    ('98100000-0000-0000-0000-000000000056', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'ph', 5.40, NULL, '2026-01-22 09:00:00+00', 'Mash pH.'),
    ('98100000-0000-0000-0000-000000000057', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'ph', 5.25, NULL, '2026-01-22 12:00:00+00', 'Post-boil pH.'),
    ('98100000-0000-0000-0000-000000000058', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'ph', 4.50, NULL, '2026-01-26 08:00:00+00', 'Mid-fermentation.'),
    ('98100000-0000-0000-0000-000000000059', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000004'), NULL, 'ph', 4.25, NULL, '2026-02-01 08:00:00+00', 'Final pH.');

-- ==============================================================================
-- Saison 24-11 (batch 8) - Warm fermentation profile
-- Fermentation: Jan 26 - Feb 5, 2026
-- ==============================================================================

INSERT INTO measurement (uuid, batch_id, occupancy_id, kind, value, unit, observed_at, notes)
VALUES
    -- Temperature readings (saison ferments very warm)
    ('98100000-0000-0000-0000-000000000060', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 20.0, 'c', '2026-01-26 19:00:00+00', 'Pitch temperature.'),
    ('98100000-0000-0000-0000-000000000061', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 22.0, 'c', '2026-01-27 08:00:00+00', 'Rising for saison character.'),
    ('98100000-0000-0000-0000-000000000062', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 24.0, 'c', '2026-01-27 20:00:00+00', 'Continuing to rise.'),
    ('98100000-0000-0000-0000-000000000063', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 26.0, 'c', '2026-01-28 08:00:00+00', 'Target warm fermentation.'),
    ('98100000-0000-0000-0000-000000000064', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 27.5, 'c', '2026-01-28 20:00:00+00', 'Peak temperature.'),
    ('98100000-0000-0000-0000-000000000065', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 28.0, 'c', '2026-01-29 08:00:00+00', 'Holding warm.'),
    ('98100000-0000-0000-0000-000000000066', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 28.5, 'c', '2026-01-29 20:00:00+00', 'Maximum warmth.'),
    ('98100000-0000-0000-0000-000000000067', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 28.0, 'c', '2026-01-30 08:00:00+00', 'Stable.'),
    ('98100000-0000-0000-0000-000000000068', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 27.0, 'c', '2026-01-31 08:00:00+00', 'Slight cooling.'),
    ('98100000-0000-0000-0000-000000000069', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 25.0, 'c', '2026-02-01 08:00:00+00', 'Fermentation slowing.'),
    ('98100000-0000-0000-0000-000000000070', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 22.0, 'c', '2026-02-02 08:00:00+00', 'Cooling for conditioning.'),
    ('98100000-0000-0000-0000-000000000071', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 18.0, 'c', '2026-02-03 08:00:00+00', 'Conditioning.'),
    ('98100000-0000-0000-0000-000000000072', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 15.0, 'c', '2026-02-04 08:00:00+00', 'Cold conditioning.'),
    ('98100000-0000-0000-0000-000000000073', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'temperature', 12.0, 'c', '2026-02-05 08:00:00+00', 'Ready for packaging.'),
    
    -- Gravity readings (saison attenuates very dry)
    ('98100000-0000-0000-0000-000000000074', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'gravity', 1.052, 'sg', '2026-01-26 19:00:00+00', 'OG.'),
    ('98100000-0000-0000-0000-000000000075', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'gravity', 1.045, 'sg', '2026-01-27 08:00:00+00', 'Day 1.'),
    ('98100000-0000-0000-0000-000000000076', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'gravity', 1.035, 'sg', '2026-01-28 08:00:00+00', 'Day 2 - rapid fermentation.'),
    ('98100000-0000-0000-0000-000000000077', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'gravity', 1.022, 'sg', '2026-01-29 08:00:00+00', 'Day 3.'),
    ('98100000-0000-0000-0000-000000000078', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'gravity', 1.015, 'sg', '2026-01-30 08:00:00+00', 'Day 4.'),
    ('98100000-0000-0000-0000-000000000079', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'gravity', 1.010, 'sg', '2026-01-31 08:00:00+00', 'Day 5.'),
    ('98100000-0000-0000-0000-000000000080', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'gravity', 1.006, 'sg', '2026-02-01 08:00:00+00', 'Day 6.'),
    ('98100000-0000-0000-0000-000000000081', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'gravity', 1.004, 'sg', '2026-02-02 08:00:00+00', 'Day 7 - very dry.'),
    ('98100000-0000-0000-0000-000000000082', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'gravity', 1.003, 'sg', '2026-02-03 08:00:00+00', 'FG stable.'),
    ('98100000-0000-0000-0000-000000000083', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'gravity', 1.002, 'sg', '2026-02-04 08:00:00+00', 'FG confirmed - bone dry.'),
    
    -- pH readings
    ('98100000-0000-0000-0000-000000000084', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'ph', 5.30, NULL, '2026-01-26 09:00:00+00', 'Mash pH.'),
    ('98100000-0000-0000-0000-000000000085', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'ph', 5.15, NULL, '2026-01-26 12:00:00+00', 'Post-boil.'),
    ('98100000-0000-0000-0000-000000000086', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'ph', 4.30, NULL, '2026-01-30 08:00:00+00', 'Mid-fermentation.'),
    ('98100000-0000-0000-0000-000000000087', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000008'), NULL, 'ph', 4.05, NULL, '2026-02-04 08:00:00+00', 'Final pH - tart finish.');

-- ==============================================================================
-- Kolsch 24-10 (batch 5) - Clean lager-like fermentation
-- Fermentation: Jan 24 - Feb 5, 2026
-- ==============================================================================

INSERT INTO measurement (uuid, batch_id, occupancy_id, kind, value, unit, observed_at, notes)
VALUES
    -- Temperature readings (cool fermentation for clean profile)
    ('98100000-0000-0000-0000-000000000088', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 15.0, 'c', '2026-01-24 13:00:00+00', 'Pitch temperature - cool for Kolsch.'),
    ('98100000-0000-0000-0000-000000000089', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 15.5, 'c', '2026-01-25 08:00:00+00', 'Day 1.'),
    ('98100000-0000-0000-0000-000000000090', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 16.0, 'c', '2026-01-25 20:00:00+00', 'Slight rise.'),
    ('98100000-0000-0000-0000-000000000091', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 16.5, 'c', '2026-01-26 08:00:00+00', 'Day 2.'),
    ('98100000-0000-0000-0000-000000000092', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 17.0, 'c', '2026-01-26 20:00:00+00', 'Peak fermentation temp.'),
    ('98100000-0000-0000-0000-000000000093', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 17.0, 'c', '2026-01-27 08:00:00+00', 'Holding steady.'),
    ('98100000-0000-0000-0000-000000000094', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 16.5, 'c', '2026-01-28 08:00:00+00', 'Cooling slightly.'),
    ('98100000-0000-0000-0000-000000000095', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 16.0, 'c', '2026-01-29 08:00:00+00', 'Day 5.'),
    ('98100000-0000-0000-0000-000000000096', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 15.0, 'c', '2026-01-30 08:00:00+00', 'Day 6.'),
    ('98100000-0000-0000-0000-000000000097', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 12.0, 'c', '2026-01-31 08:00:00+00', 'Lagering start.'),
    ('98100000-0000-0000-0000-000000000098', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 8.0, 'c', '2026-02-01 08:00:00+00', 'Lagering.'),
    ('98100000-0000-0000-0000-000000000099', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 4.0, 'c', '2026-02-02 08:00:00+00', 'Cold lagering.'),
    ('98100000-0000-0000-0000-000000000100', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 2.0, 'c', '2026-02-03 08:00:00+00', 'Cold conditioning.'),
    ('98100000-0000-0000-0000-000000000101', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 1.0, 'c', '2026-02-04 08:00:00+00', 'Final lagering.'),
    ('98100000-0000-0000-0000-000000000102', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'temperature', 1.0, 'c', '2026-02-05 08:00:00+00', 'Ready for packaging.'),
    
    -- Gravity readings
    ('98100000-0000-0000-0000-000000000103', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'gravity', 1.048, 'sg', '2026-01-24 13:00:00+00', 'OG.'),
    ('98100000-0000-0000-0000-000000000104', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'gravity', 1.044, 'sg', '2026-01-25 08:00:00+00', 'Day 1.'),
    ('98100000-0000-0000-0000-000000000105', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'gravity', 1.038, 'sg', '2026-01-26 08:00:00+00', 'Day 2.'),
    ('98100000-0000-0000-0000-000000000106', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'gravity', 1.030, 'sg', '2026-01-27 08:00:00+00', 'Day 3.'),
    ('98100000-0000-0000-0000-000000000107', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'gravity', 1.024, 'sg', '2026-01-28 08:00:00+00', 'Day 4.'),
    ('98100000-0000-0000-0000-000000000108', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'gravity', 1.018, 'sg', '2026-01-29 08:00:00+00', 'Day 5.'),
    ('98100000-0000-0000-0000-000000000109', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'gravity', 1.014, 'sg', '2026-01-30 08:00:00+00', 'Day 6.'),
    ('98100000-0000-0000-0000-000000000110', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'gravity', 1.012, 'sg', '2026-02-01 08:00:00+00', 'Lagering.'),
    ('98100000-0000-0000-0000-000000000111', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'gravity', 1.010, 'sg', '2026-02-03 08:00:00+00', 'FG stable.'),
    ('98100000-0000-0000-0000-000000000112', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'gravity', 1.010, 'sg', '2026-02-05 08:00:00+00', 'FG confirmed.'),
    
    -- pH readings
    ('98100000-0000-0000-0000-000000000113', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'ph', 5.35, NULL, '2026-01-24 09:00:00+00', 'Mash pH.'),
    ('98100000-0000-0000-0000-000000000114', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'ph', 5.20, NULL, '2026-01-24 12:00:00+00', 'Post-boil.'),
    ('98100000-0000-0000-0000-000000000115', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'ph', 4.45, NULL, '2026-01-28 08:00:00+00', 'Mid-fermentation.'),
    ('98100000-0000-0000-0000-000000000116', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000005'), NULL, 'ph', 4.30, NULL, '2026-02-05 08:00:00+00', 'Final pH.');

-- ==============================================================================
-- Pilsner 24-08 (batch 3) - Cold lager fermentation
-- Fermentation: Jan 18 - Feb 1, 2026
-- ==============================================================================

INSERT INTO measurement (uuid, batch_id, occupancy_id, kind, value, unit, observed_at, notes)
VALUES
    -- Temperature readings (coldest fermentation)
    ('98100000-0000-0000-0000-000000000117', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 10.0, 'c', '2026-01-18 12:00:00+00', 'Pitch temperature - cold for lager.'),
    ('98100000-0000-0000-0000-000000000118', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 10.5, 'c', '2026-01-19 08:00:00+00', 'Day 1.'),
    ('98100000-0000-0000-0000-000000000119', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 11.0, 'c', '2026-01-19 20:00:00+00', 'Slight rise.'),
    ('98100000-0000-0000-0000-000000000120', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 11.5, 'c', '2026-01-20 08:00:00+00', 'Day 2.'),
    ('98100000-0000-0000-0000-000000000121', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 12.0, 'c', '2026-01-20 20:00:00+00', 'Peak lager temp.'),
    ('98100000-0000-0000-0000-000000000122', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 12.0, 'c', '2026-01-21 08:00:00+00', 'Holding steady.'),
    ('98100000-0000-0000-0000-000000000123', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 11.5, 'c', '2026-01-22 08:00:00+00', 'Day 4.'),
    ('98100000-0000-0000-0000-000000000124', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 11.0, 'c', '2026-01-23 08:00:00+00', 'Day 5.'),
    ('98100000-0000-0000-0000-000000000125', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 10.5, 'c', '2026-01-24 08:00:00+00', 'Day 6.'),
    ('98100000-0000-0000-0000-000000000126', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 10.0, 'c', '2026-01-25 08:00:00+00', 'Day 7.'),
    ('98100000-0000-0000-0000-000000000127', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 8.0, 'c', '2026-01-26 08:00:00+00', 'Lagering start.'),
    ('98100000-0000-0000-0000-000000000128', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 5.0, 'c', '2026-01-27 08:00:00+00', 'Cold lagering.'),
    ('98100000-0000-0000-0000-000000000129', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 2.0, 'c', '2026-01-28 08:00:00+00', 'Deep lagering.'),
    ('98100000-0000-0000-0000-000000000130', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 1.0, 'c', '2026-01-29 08:00:00+00', 'Near freezing.'),
    ('98100000-0000-0000-0000-000000000131', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 0.5, 'c', '2026-01-30 08:00:00+00', 'Final lagering.'),
    ('98100000-0000-0000-0000-000000000132', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'temperature', 0.5, 'c', '2026-02-01 08:00:00+00', 'Ready for packaging.'),
    
    -- Gravity readings
    ('98100000-0000-0000-0000-000000000133', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'gravity', 1.044, 'sg', '2026-01-18 12:00:00+00', 'OG.'),
    ('98100000-0000-0000-0000-000000000134', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'gravity', 1.042, 'sg', '2026-01-19 08:00:00+00', 'Day 1 - slow start.'),
    ('98100000-0000-0000-0000-000000000135', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'gravity', 1.038, 'sg', '2026-01-20 08:00:00+00', 'Day 2.'),
    ('98100000-0000-0000-0000-000000000136', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'gravity', 1.032, 'sg', '2026-01-21 08:00:00+00', 'Day 3.'),
    ('98100000-0000-0000-0000-000000000137', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'gravity', 1.026, 'sg', '2026-01-22 08:00:00+00', 'Day 4.'),
    ('98100000-0000-0000-0000-000000000138', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'gravity', 1.020, 'sg', '2026-01-23 08:00:00+00', 'Day 5.'),
    ('98100000-0000-0000-0000-000000000139', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'gravity', 1.016, 'sg', '2026-01-24 08:00:00+00', 'Day 6.'),
    ('98100000-0000-0000-0000-000000000140', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'gravity', 1.012, 'sg', '2026-01-25 08:00:00+00', 'Day 7.'),
    ('98100000-0000-0000-0000-000000000141', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'gravity', 1.010, 'sg', '2026-01-27 08:00:00+00', 'Lagering.'),
    ('98100000-0000-0000-0000-000000000142', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'gravity', 1.008, 'sg', '2026-01-30 08:00:00+00', 'FG stable.'),
    ('98100000-0000-0000-0000-000000000143', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'gravity', 1.008, 'sg', '2026-02-01 08:00:00+00', 'FG confirmed.'),
    
    -- pH readings
    ('98100000-0000-0000-0000-000000000144', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'ph', 5.40, NULL, '2026-01-18 09:00:00+00', 'Mash pH.'),
    ('98100000-0000-0000-0000-000000000145', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'ph', 5.25, NULL, '2026-01-18 11:00:00+00', 'Post-boil.'),
    ('98100000-0000-0000-0000-000000000146', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'ph', 4.55, NULL, '2026-01-24 08:00:00+00', 'Mid-fermentation.'),
    ('98100000-0000-0000-0000-000000000147', (SELECT id FROM batch WHERE uuid = '90000000-0000-0000-0000-000000000003'), NULL, 'ph', 4.40, NULL, '2026-02-01 08:00:00+00', 'Final pH.');
    