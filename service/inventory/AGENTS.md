# BrewPipes "Inventory" Service Agent Guide

This guide is for agentic coding tools working in this particular service in the BrewPipes repo.
It captures commands and conventions observed in the current codebase.
When making changes to this service, its data models, its logic, or any other aspect, this document MUST be kept-up-to-date do that it accurately reflects the actual implementation.

## Service Domain

The Inventory service models ingredients and consumables inventory: malt/fermentables, hop lots, yeast, adjuncts, salts, chemicals, gases, and finished beer lots.

## Overview

The big picture: the system tracks ingredients from receiving, into storage, through usage, and into adjustments and transfers.

Ingredient
- An ingredient is the master definition (e.g., "Pale Malt", "Cascade 2024", "WLP001", "Lactic Acid", "CO2").
- It defines the category, default unit, and any critical attributes used in inventory calculations.
- Category-specific attributes live in one-to-one detail records (e.g., malt, hop, yeast) keyed by ingredient ID.

Ingredient Lot
- A lot is a specific received batch of an ingredient with its own quantity, dates, supplier reference, and originator details.
- Lots are the unit of traceability for usage and quality (e.g., alpha acids for hops, yeast generation, moisture for grain).
- Lot-specific quality attributes (e.g., measured alpha acids, moisture, viability) are stored on the lot.

Supplier Reference
- Suppliers are managed in the Procurement service; inventory records store an opaque supplier UUID.
- The supplier is who you buy from, not necessarily who produced the ingredient.

Originator
- The originator is the maltster, hop producer, yeast lab, or gas vendor who produced the ingredient.
- It can be stored as a distinct record or as explicit name fields on a lot.

Lot Identifiers
- Each lot supports a brewery-assigned lot number plus the originator's lot number.
- Both identifiers are preserved for traceability and audits.

Stock Location
- A stock location is where inventory is stored (grain room, cold room, chemical cage, bulk tank).
- Locations help track on-hand balances and support transfers between areas.

Inventory Movement
- A movement is the append-only record of inventory changes.
- Movements capture the amount, unit, timestamp, and a reason (receive, use, transfer, adjust, waste).

Receipt
- A receipt records receiving inventory from a supplier or a production transfer-in.
- It creates one or more lots and the corresponding movement records.

Usage
- Usage records consumption of inventory for brewing or operations.
- It references the lot and can link to production actions via opaque UUIDs.

Adjustment
- Adjustments capture corrections, shrink, spoilage, or cycle count deltas.
- They are always explicit movements with a reason and optional notes.

Transfer
- Transfers move inventory between locations without changing the total on-hand.
- The transfer is represented by paired movements (out from source, in to destination).

Beer Lot
- A beer lot tracks finished product inventory and ties back to a production batch UUID.
- It supports traceability from packaged product back to the source batch.

## User Journey: Inventory Manager

Here is a simple inventory story that follows the inventory records, told in brewery terms.

You add an ingredient called "Pale Malt" with a default unit of kilograms. A supplier delivers a new lot, so you record a receipt with your internal lot code, the maltster's lot code, quantity, and location (grain room). The system creates the lot and a receive movement so the on-hand balance increases.

You also receive hops and yeast, each as their own lots stored in the cold room. When brew day arrives, the brewer records ingredient usage for the batch. Each usage references a specific lot so traceability is preserved, and the lot balance decreases.

If a bag of malt gets damaged, you record an adjustment for spoilage. If hops are moved from the cold room to the brew deck, you record a transfer to the new location. Over time, on-hand inventory is the sum of movements across all lots and locations. When the batch is packaged, you create a beer lot tied to the production batch UUID for finished goods traceability.

In short:
- ingredients define what you stock,
- lots track what you received,
- movements record every change,
- receipts increase inventory,
- usage reduces inventory,
- adjustments reconcile reality,
- transfers move stock between locations.

## Acceptance Criteria

- An inventory manager can create an ingredient with a name, category, and default unit for fermentables, hops, yeast, adjuncts, and gases.
- Category-specific ingredient details (malt, hop, yeast) can be stored without fragmenting lot or movement workflows.
- The system can record a receipt that creates one or more ingredient lots and increases on-hand balances.
- Lots support both a brewery-assigned lot number and an originator's lot number.
- Lots store a supplier reference UUID (Procurement) and originator details for traceability.
- The system can record usage against a specific lot with amount, unit, time, and an optional production reference UUID.
- Transfers between locations are captured with source, destination, amount, and time, and do not change total on-hand.
- Adjustments can be recorded with a reason and notes, and they update on-hand balances.
- Inventory balances can be derived from movements and reconciled by location and by lot.
- Inventory records retain traceability to production and procurement via opaque UUIDs, without shared tables or foreign keys.
- Beer lots can be created for finished product inventory and related back to production batch UUIDs.
