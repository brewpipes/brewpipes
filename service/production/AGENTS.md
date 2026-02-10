# BrewPipes "Production" Service Agent Guide

This guide is for agentic coding tools working in this particular service in the BrewPipes repo.
It captures commands and conventions observed in the current codebase.
When making changes to this service, its data models, its logic, or any other aspect, this document MUST be kept-up-to-date do that it accurately reflects the actual implementation.

## Service Domain

The Production service models the data and workflows of producing batches of beer, from creation to fermentation to packaging.

## Overview

The big picture: the system tracks your beer as it moves through tanks, gets split or blended, has ingredients added, and gets measured over time.

## Auth

- All HTTP routes require `Authorization: Bearer <access token>`.
- Service startup fails if `BREWPIPES_SECRET_KEY` is missing.

Batch
- A batch is the overall production run you care about (e.g., "IPA 24‑07").
- It's the top‑level record you plan, brew, ferment, and eventually package.
- Batches can be updated via `PATCH /batches/{id}` with editable fields: `short_name`, `brew_date`, `recipe_id`, `notes`.
- Batches can be soft-deleted via `DELETE /batches/{id}`, which sets `deleted_at`.
- Deletion cascades soft-delete to all related records: batch volumes, process phases, brew sessions, additions, and measurements.
- The cascade is performed in a single transaction to ensure data consistency.

Volume
- A volume is a specific quantity of liquid at a point in time (e.g., “10,000 liters of wort”).
- Volumes can be split (one volume becomes two) or blended (two volumes become one). These relationships are tracked explicitly, including amounts.

Volume Relation
- Tracks liquid lineage between volumes for splits and blends.
- Each relation includes the amount and unit contributed from the parent to the child.

Modeling blends
- For a blend, create a new child volume and add one volume relation per parent with the contributed amount and unit.
- The child volume amount should match the sum of relation amounts, minus any recorded loss elsewhere.

Vessel
- A vessel is any physical container (mash tun, kettle, fermenter, brite tank).
- It has a capacity and type so the system knows what it can hold.
- Vessels can be looked up by internal ID (`GET /vessels/{id}`) or by UUID (`GET /vessels/uuid/{uuid}`).

Occupancy
- Occupancy is the idea of "this volume is in this vessel during this time."
- When you move beer from one vessel to another, the occupancy changes to reflect where the liquid is now.
- Active occupancies are listed with `GET /occupancies?active=true`; single active lookups use `GET /occupancies/active` with `active_vessel_id` or `active_volume_id`.
- Active occupancy responses include a derived `batch_id` field (from the most recent batch_volume record for that volume) to help clients display which batch is in each vessel.

Transfer
- A transfer records a move from one occupancy to another (e.g., fermenter A to brite tank B).
- It captures how much moved, when it started/ended, and any losses.

Batch Volume
- This links a batch to the volume that represents it at a given liquid phase (water, wort, beer).
- It keeps the batch identity attached to the evolving liquid as it moves and changes.

Batch Process Phase
- This records the production phase history for a batch (planning, mashing, boiling, fermenting, conditioning, packaging, finished, etc).
- Liquid phase and process phase are tracked separately so “wort” is never shown as “fermenting.”

Batch Relation
- This is how splits and blends are tracked at the batch level.
  - Example: one batch splits into two smaller batches (parent → child).
  - Example: two batches are blended into a single new batch.
- It can also reference the specific volume involved.

Addition
- An addition records anything you add to the beer: malt, hops, yeast, chemicals, gas, adjuncts.
- It’s tied to either a batch or a specific vessel occupancy.
- It references the ingredient lot by UUID for ingredient additions; non-ingredient additions can omit it.

Measurement
- Measurements capture what you observe: gravity, temperature, pH, CO2, ABV, IBU, etc.
- They can be tied to the batch overall or to a specific vessel occupancy.
- This lets you track fermentation and quality over time.

Style
- A reference table for beer styles with case-insensitive unique names.
- Enables autocomplete while allowing free-form entry.

Recipe
- A beer formulation with name and optional style reference.
- Batches can reference a recipe via `recipe_id` (internal FK).
- All API endpoints use UUIDs, not internal IDs.
- Includes target specifications: batch size, OG/FG ranges, IBU/SRM ranges, carbonation, brewhouse efficiency.
- `target_abv` is calculated from `target_og` and `target_fg` using `(OG - FG) × 131.25` (not stored in DB).
- IBU calculation method can be: tinseth, rager, garetz, or daniels.
- CRUD endpoints: `GET/POST /recipes`, `GET/PUT/PATCH/DELETE /recipes/{uuid}`.
- API responses include only `uuid`, not internal `id`.

Recipe Ingredient
- Stores the ingredient bill for a recipe.
- Each ingredient has: type, amount, unit, use stage, optional use type, timing, and notes.
- Ingredient types: fermentable, hop, yeast, adjunct, salt, chemical, gas, other.
- Use stages: mash, boil, whirlpool, fermentation, packaging.
- Use types vary by context: bittering/flavor/aroma/dry_hop (hops), base/specialty/adjunct/sugar (fermentables), primary/secondary/bottle (yeast).
- `alpha_acid_assumed` is only valid for hop ingredients.
- `scaling_factor` controls how the ingredient scales with batch size (default 1.0).
- `ingredient_uuid` is a cross-service reference to inventory.ingredient (no FK).
- CRUD endpoints: `GET/POST /recipes/{uuid}/ingredients`, `GET/PATCH/DELETE /recipes/{uuid}/ingredients/{ingredient_uuid}`.
- API responses include `uuid` and `recipe_uuid`, not internal IDs.

Brew Session
- Captures hot-side wort production (mash → boil → knockout).
- Points to the wort volume produced via `wort_volume_id`.
- A batch can have multiple brew sessions (double batching).
- References mash and boil vessels for traceability.

## Batch Summary

The `GET /batches/{id}/summary` endpoint provides an aggregated view of a batch with derived metrics:

### Core batch info
- `id`, `uuid`, `short_name`, `brew_date`, `notes`

### Recipe and style
- `recipe_name`, `style_name` (from linked recipe)

### Brew sessions
- Array of `brew_sessions` with `id`, `brewed_at`, `notes`

### Current state
- `current_phase` (most recent process phase)
- `current_vessel` (from active occupancy)
- `current_occupancy_status` (fermenting, conditioning, etc.)

### Key measurements
- `original_gravity` (first OG measurement)
- `final_gravity` (most recent FG measurement)
- `abv` (manual measurement or calculated from OG/FG)
- `abv_calculated` (true if ABV was auto-calculated using `(OG - FG) × 131.25`)
- `ibu` (most recent IBU measurement)

### Duration metrics (in days)
- `days_in_fermenter` (total time in fermenter vessels)
- `days_in_brite` (total time in brite vessels)
- `days_grain_to_glass` (from first brew date to now or completion)

### Volume and loss metrics
- `starting_volume_bbl` (earliest batch volume, converted to BBL)
- `current_volume_bbl` (latest batch volume, converted to BBL)
- `total_loss_bbl` (sum of transfer losses, converted to BBL)
- `loss_percentage` (total loss / starting volume × 100)

### Volume unit conversions
The summary converts volumes to US barrels (BBL = 31 gallons) from:
- `bbl` (US barrels, direct)
- `ml` (milliliters)
- `usfloz` (US fluid ounces)
- `ukfloz` (UK fluid ounces)

## User Journey: Brewer

Here’s a simple "brewer's story" that follows the production records, told in brewery terms.

You plan a new batch called “IPA 24‑07.” That creates the batch record.

Brew day starts: you record a process phase of “mashing.” A volume is created for the liquid you’re working with, and it’s marked as occupying the mash tun. The liquid phase is “water” or “wort” depending on the step. You add malt — the system records an addition, tied to this batch or vessel, and it notes which ingredient lot it came from (so inventory can decrement it later). You take a mash temp and pH reading — those are measurements attached to this batch or vessel.

You transfer to the kettle. That’s a transfer: the volume leaves the mash tun occupancy and enters the kettle occupancy. During the boil, you add hops and record each hop addition. You take a gravity reading — another measurement.

You transfer to the fermenter. That’s another transfer and occupancy change. You pitch yeast (an addition). Over the next few days, you record daily gravity and temperature measurements.

If you split the batch into two fermenters, the original volume becomes two new volumes. The system records volume relations with the split amounts, and it notes a batch split (a batch relation), so it’s clear the two child batches came from the original batch.

Later, you transfer to the brite tank. Another transfer and occupancy change. You dry hop or dose CO2 if needed — additions again. You take a final gravity and CO2 measurement. You record process phases like conditioning and packaging, while the liquid phase stays “beer.”

As the beer progresses, the system updates the batch volume liquid phase from “water” to “wort” to “beer.” Process phases track the brewing workflow separately. When you package, production records the packaging phase, and inventory uses the batch and addition references to manage stock and finished goods.

In short:  
- batches describe the brew run,  
- volumes track the liquid itself,  
- vessels/occupancies show where it sits,  
- transfers show movement,  
- additions record ingredients and chemicals,  
- measurements capture quality and process data,  
- batch relations capture splits/blends.

## Acceptance Criteria

- A brewmaster can create a new production batch with a short name, and the batch appears in the production list.
- A brewmaster can bulk import batches from CSV via `POST /batches/import` with per-row results.
- The system can record a starting volume for a batch and place it into a specific vessel as an active occupancy.
- Transfers between vessels are captured with source, destination, amount, and timestamps, and the active occupancy updates accordingly.
- Ingredient, chemical, and gas additions can be recorded against a batch or a specific vessel occupancy with amount, unit, time, and an external inventory lot UUID.
- Measurements (e.g., gravity, temperature, pH, CO2) can be recorded against a batch or a specific vessel occupancy with value, unit, and time.
- A batch’s liquid phase can be updated over time (water → wort → beer) and is visible in batch history.
- A batch’s process phase can be updated over time (mashing → boiling → fermenting → conditioning → packaging → finished) and is visible in batch history.
- Splits and blends are represented so a brewmaster can see parent/child batch relationships and the related volume changes.
- Production records retain traceability to inventory/procurement via opaque UUIDs, without requiring shared tables or foreign keys.
- The full brew‑day flow can be reconstructed chronologically from additions, measurements, transfers, and phase changes for a given batch.
- A brewmaster can view a batch summary with aggregated metrics including recipe/style, brew sessions, current state, OG/FG/ABV/IBU, duration metrics, and loss calculations via `GET /batches/{id}/summary`.
- ABV is auto-calculated from OG and FG measurements using `(OG - FG) × 131.25` when no manual ABV measurement exists.
- A brewmaster can define recipe target specifications including batch size, OG/FG/IBU/SRM ranges, carbonation, and brewhouse efficiency.
- A brewmaster can manage recipe ingredient bills with full CRUD operations on `/recipes/{recipe_id}/ingredients`.
