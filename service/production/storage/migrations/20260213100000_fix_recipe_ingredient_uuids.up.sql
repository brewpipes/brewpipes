-- Fix recipe_ingredient cross-service ingredient_uuid references.
-- The original seed data used sequential UUIDs that don't match the actual
-- inventory ingredient records. This migration corrects each mismatch so
-- the pick list feature can resolve recipe ingredients to inventory.
BEGIN;

-- Flagship IPA --

-- b...02: "Crystal 40L for body" pointed to Pilsner Malt (...0002), should be Crystal 60L (...0003)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000003'
WHERE uuid = 'b0000000-0000-0000-0000-000000000002';

-- b...03: "Munich for complexity" pointed to Crystal 60L (...0003), should be Munich Malt (...0032)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000032'
WHERE uuid = 'b0000000-0000-0000-0000-000000000003';

-- b...04: "Columbus for bittering" pointed to Dextrose (...0010), should be Columbus Hops (...0025)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000025'
WHERE uuid = 'b0000000-0000-0000-0000-000000000004';

-- b...05: "Citra for flavor" pointed to Roasted Barley (...0011), should be Citra Hops (...0005)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000005'
WHERE uuid = 'b0000000-0000-0000-0000-000000000005';

-- b...06: "Mosaic whirlpool" pointed to Saaz Hops (...0012), should be Mosaic Hops (...0026)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000026'
WHERE uuid = 'b0000000-0000-0000-0000-000000000006';

-- b...07: "Citra dry hop" pointed to Roasted Barley (...0011), should be Citra Hops (...0005)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000005'
WHERE uuid = 'b0000000-0000-0000-0000-000000000007';

-- b...08: "Mosaic dry hop" pointed to Saaz Hops (...0012), should be Mosaic Hops (...0026)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000026'
WHERE uuid = 'b0000000-0000-0000-0000-000000000008';

-- Crisp Pilsner --

-- b...11: "Pilsner malt" pointed to Cascade Hops (...0004), should be Pilsner Malt (...0002)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000002'
WHERE uuid = 'b0000000-0000-0000-0000-000000000011';

-- b...12: "Saaz bittering" pointed to W34/70 Lager Yeast (...0013), should be Saaz Hops (...0012)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000012'
WHERE uuid = 'b0000000-0000-0000-0000-000000000012';

-- b...13: "Saaz late addition" pointed to W34/70 Lager Yeast (...0013), should be Saaz Hops (...0012)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000012'
WHERE uuid = 'b0000000-0000-0000-0000-000000000013';

-- Robust Stout --

-- b...16: "Roasted barley" pointed to Citra Hops (...0005), should be Roasted Barley (...0011)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000011'
WHERE uuid = 'b0000000-0000-0000-0000-000000000016';

-- b...17: "Chocolate malt" pointed to WLP001 California Ale Yeast (...0006), should be Chocolate Malt (...0033)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000033'
WHERE uuid = 'b0000000-0000-0000-0000-000000000017';

-- b...18: "Flaked barley for head" pointed to Lactic Acid (...0007), should be Flaked Barley (...0034)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000034'
WHERE uuid = 'b0000000-0000-0000-0000-000000000018';

-- b...19: "East Kent Goldings" pointed to Irish Moss (...0014), should be East Kent Goldings (...0027)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000027'
WHERE uuid = 'b0000000-0000-0000-0000-000000000019';

-- Cologne Classic (Kolsch) --

-- b...21: "Pilsner malt" pointed to Cascade Hops (...0004), should be Pilsner Malt (...0002)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000002'
WHERE uuid = 'b0000000-0000-0000-0000-000000000021';

-- b...22: "Vienna malt" pointed to Gypsum (...0008), should be Vienna Malt (...0035)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000035'
WHERE uuid = 'b0000000-0000-0000-0000-000000000022';

-- b...23: "Hallertau Mittelfruh" pointed to Oxygen (...0015), should be Hallertau Mittelfruh (...0028)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000028'
WHERE uuid = 'b0000000-0000-0000-0000-000000000023';

-- b...24: "Hallertau late" pointed to Oxygen (...0015), should be Hallertau Mittelfruh (...0028)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000028'
WHERE uuid = 'b0000000-0000-0000-0000-000000000024';

-- Farmhouse Saison --

-- b...26: "Pilsner malt" pointed to Cascade Hops (...0004), should be Pilsner Malt (...0002)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000002'
WHERE uuid = 'b0000000-0000-0000-0000-000000000026';

-- b...27: "Wheat malt" pointed to CO2 (...0009), should be Wheat Malt (...0036)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000036'
WHERE uuid = 'b0000000-0000-0000-0000-000000000027';

-- b...28: "Munich malt" pointed to Crystal 60L (...0003), should be Munich Malt (...0032)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000032'
WHERE uuid = 'b0000000-0000-0000-0000-000000000028';

-- b...29: "Styrian Goldings" pointed to Coriander (...0016), should be Styrian Goldings (...0029)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000029'
WHERE uuid = 'b0000000-0000-0000-0000-000000000029';

-- b...30: "Saaz late" pointed to W34/70 Lager Yeast (...0013), should be Saaz Hops (...0012)
UPDATE recipe_ingredient SET ingredient_uuid = '70000000-0000-0000-0000-000000000012'
WHERE uuid = 'b0000000-0000-0000-0000-000000000030';

COMMIT;
