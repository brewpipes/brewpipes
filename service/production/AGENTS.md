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
- A batch is the overall production run you care about (e.g., “IPA 24‑07”).
- It’s the top‑level record you plan, brew, ferment, and eventually package.

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

Occupancy
- Occupancy is the idea of “this volume is in this vessel during this time.”
- When you move beer from one vessel to another, the occupancy changes to reflect where the liquid is now.
- Active occupancies are listed with `GET /occupancies?active=true`; single active lookups use `GET /occupancies/active` with `active_vessel_id` or `active_volume_id`.

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
