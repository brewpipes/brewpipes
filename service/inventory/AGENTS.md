# BrewPipes "Inventory" Service Agent Guide

This guide is for agentic coding tools working in this particular service in the BrewPipes repo.
It captures commands and conventions observed in the current codebase.
When making changes to this service, its data models, its logic, or any other aspect, this document MUST be kept-up-to-date do that it accurately reflects the actual implementation.

## Service Domain

The Inventory service models ingredients and consumables inventory: malt/fermentables, hop lots, yeast, adjuncts, salts, chemicals, gases, and finished beer lots.

## Overview

The big picture: the system tracks ingredients from receiving, into storage, through usage, and into adjustments and transfers.

## Auth

- All HTTP routes require `Authorization: Bearer <access token>`.
- Service startup fails if `BREWPIPES_SECRET_KEY` is missing.

Ingredient
- An ingredient is the master definition (e.g., "Pale Malt", "Cascade 2024", "WLP001", "Lactic Acid", "CO2").
- It defines the category, default unit, and any critical attributes used in inventory calculations.
- Category-specific attributes live in one-to-one detail records (e.g., malt, hop, yeast) keyed by ingredient UUID.

Ingredient Lot
- A lot is a specific received batch of an ingredient with its own quantity, dates, supplier reference, and originator details.
- Lots can optionally reference the procurement purchase order line UUID for traceability.
- Lots are the unit of traceability for usage and quality (e.g., alpha acids for hops, yeast generation, moisture for grain).
- Lot-specific quality attributes live in lot detail records keyed by ingredient lot UUID.

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
- Receipts can optionally reference the originating procurement purchase order UUID.
- It creates one or more lots and the corresponding movement records.

Usage
- Usage records consumption of inventory for brewing or operations.
- It references the lot and can link to production actions via opaque UUIDs.
- Batch usage deduction (`POST /inventory-usage/batch`) atomically creates a usage record and one movement per ingredient pick in a single transaction, with stock validation per lot/location.

Adjustment
- Adjustments capture corrections, shrink, spoilage, or cycle count deltas.
- They are always explicit movements with a reason and optional notes.
- `POST /inventory-adjustments` atomically creates the adjustment record and a corresponding inventory movement in a single transaction.
- The request accepts `ingredient_lot_uuid` or `beer_lot_uuid` (exactly one), `stock_location_uuid`, `amount` (positive for increase, negative for decrease), `amount_unit`, `reason`, and optional `notes` and `adjusted_at`.
- The movement uses `reason='adjust'`, `direction='in'` for positive amounts or `'out'` for negative, and `amount=abs(amount)`.

Transfer
- Transfers move inventory between locations without changing the total on-hand.
- The transfer is represented by paired movements (out from source, in to destination).
- `POST /inventory-transfers` atomically creates the transfer record and two inventory movements (out from source, in to destination) in a single transaction.
- The request accepts `ingredient_lot_uuid` or `beer_lot_uuid` (exactly one), `source_location_uuid`, `dest_location_uuid`, `amount` (positive), `amount_unit`, and optional `notes` and `transferred_at`.

Beer Lot
- A beer lot tracks finished product inventory and ties back to a production batch UUID.
- It supports traceability from packaged product back to the source batch.

Batch Ingredient Lot Lookup
- `GET /ingredient-lots/batch?production_ref_uuid={uuid}` returns all ingredient lots consumed by a specific production batch.
- Joins through `inventory_usage` → `inventory_movement` → `ingredient_lot` → `ingredient` to resolve lot and ingredient details.
- Returns `ingredient_lot_uuid`, `ingredient_uuid`, `ingredient_name`, `ingredient_category`, `brewery_lot_code`, `purchase_order_line_uuid`, and `received_unit`.
- Uses `DISTINCT ON (il.uuid)` to deduplicate when multiple movements reference the same lot.
- Returns empty array (not 404) when no matching lots are found.
- Used by the Production service for batch cost calculation.

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
- Receipts and lots can store procurement purchase order and line UUIDs for cross-service traceability.
- Lots support both a brewery-assigned lot number and an originator's lot number.
- Lots store a supplier reference UUID (Procurement) and originator details for traceability.
- The system can record usage against a specific lot with amount, unit, time, and an optional production reference UUID.
- Transfers between locations are captured with source, destination, amount, and time, and atomically create paired movements (out/in) without changing total on-hand.
- Adjustments can be recorded with a reason and notes, and they atomically create a corresponding inventory movement to update on-hand balances.
- Inventory balances can be derived from movements and reconciled by location and by lot.
- Inventory records retain traceability to production and procurement via opaque UUIDs, without shared tables or foreign keys.
- Beer lots can be created for finished product inventory and related back to production batch UUIDs.
- Beer lots support packaging metadata: packaging run UUID, best-by date, package format name, container type, volume per unit, and quantity.
- Beer lots can be created with an initial inventory movement atomically when a stock location is provided (`POST /beer-lots` with `stock_location_uuid`).
- Beer lot stock levels can be queried via `GET /beer-lot-stock-levels`, showing current volume and derived quantity per lot per location.
- A brewer can atomically deduct inventory for a batch's ingredient picks via `POST /inventory-usage/batch`, with per-pick stock validation and descriptive error messages on insufficient stock.
- The system can return all ingredient lots consumed by a production batch via `GET /ingredient-lots/batch?production_ref_uuid={uuid}`, joining through usage records and movements to resolve lot and ingredient details.

## API Convention: UUID-Only

All inventory API endpoints use UUIDs exclusively for resource identification:
- Path parameters use `{uuid}` (e.g., `/ingredients/{uuid}`, `/stock-locations/{uuid}`).
- Request bodies use `_uuid` suffix for FK references (e.g., `ingredient_uuid`, `receipt_uuid`, `ingredient_lot_uuid`).
- Response bodies include `uuid` fields, never internal `id` fields.
- Query parameters use `_uuid` suffix (e.g., `ingredient_uuid`, `ingredient_lot_uuid`, `receipt_uuid`).
- Internal storage models retain both `ID int64` and UUID fields; int64 is for DB operations only.
- Handlers resolve UUID→internal ID for creates via small lookup queries before INSERT.
