# BrewPipes Project Brief

## Product vision

BrewPipes is an open source brewery management system focused on day-to-day production operations. It aims to be the operational source of truth for procurement, inventory, and production, with clear traceability from ingredients and receipts through batches, vessels, and finished lots.

## Core modules

- Identity: access tokens and authentication for all services.
- Production: batches, volumes, vessels, occupancies, transfers, additions, measurements, and process phases.
- Inventory: ingredients, lots, receipts, stock locations, usage, adjustments, transfers, and movements.
- Procurement: suppliers, purchase orders, order lines, and fees.
- Web UI: Vue 3 + Vuetify app under `service/www/` that drives the primary operator experience.

## Primary user journeys

- Sign in and access the operational dashboard.
- Set up production basics (vessels and stock locations).
- Manage suppliers and create purchase orders with line items and fees.
- Receive ingredients into inventory, creating lots and tracking movements.
- Plan and run production batches, tracking volumes through vessels and transfers.
- Bulk import planned batches via CSV for rapid setup.
- Record additions and measurements during production phases.
- Reconcile inventory usage, adjustments, and transfers tied to real operations.
- Maintain traceability across procurement, inventory, and production using UUID references.

## Product principles

- Protect core workflows: batch lifecycle, inventory accuracy, procurement traceability, and authentication.
- Prefer backward-compatible changes to schemas and APIs.
- Keep logging and error handling stable and predictable.
- Embrace rapid iteration on non-core enhancements with clear rollback paths.

## Implementation context

- Services live under `service/` and run independently or via the monolith entrypoint.
- The web app lives under `service/www/` and uses API clients that default to `/api` unless overridden by `VITE_*_API_URL`.
- Auth tokens are stored in the browser and used for Bearer requests across services.

## Core entities (current)

- Procurement: supplier, purchase_order, purchase_order_line, purchase_order_fee.
- Inventory: ingredient, ingredient_*_detail, stock_location, inventory_receipt, ingredient_lot, inventory_usage, inventory_adjustment, inventory_transfer, inventory_movement, beer_lot.
- Production: batch, volume, volume_relation, vessel, occupancy, transfer, batch_volume, batch_process_phase, batch_relation, addition, measurement.

## Change posture

- Be conservative with changes that alter data integrity, migrations, or cross-service links.
- Iterate faster on DTOs, response shaping, UI affordances, and non-critical endpoints.

## Open questions

- Bulk batch import: confirm required CSV columns beyond `short_name` (default: `short_name` required, `brew_date` and `notes` optional) and whether duplicate `short_name` values should be rejected.
