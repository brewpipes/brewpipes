# BrewPipes Project Brief

## Product vision

BrewPipes is an open source brewery management system focused on day-to-day production operations. It aims to be the operational source of truth for procurement, inventory, and production, with clear traceability from ingredients and receipts through batches, vessels, and finished lots.

## Core modules

- Identity: access tokens and authentication for all services.
- Production: recipes, styles, batches, brew sessions, volumes, vessels, occupancies, transfers, additions, measurements, and process phases.
- Inventory: ingredients, lots, receipts, stock locations, usage, adjustments, transfers, and movements.
- Procurement: suppliers, purchase orders, order lines, and fees.
- Web UI: Vue 3 + Vuetify app under `service/www/` that drives the primary operator experience.

## Primary user journeys

- Sign in and access the operational dashboard.
- Set up production basics (vessels and stock locations).
- Manage suppliers and create purchase orders with line items and fees.
- Procurement workflow starts with suppliers, then purchase orders; line items and fees are managed from the purchase order context.
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
- Production: style, recipe, batch, brew_session, volume, volume_relation, vessel, occupancy, transfer, batch_volume, batch_process_phase, batch_relation, addition, measurement.

## Change posture

- Be conservative with changes that alter data integrity, migrations, or cross-service links.
- Iterate faster on DTOs, response shaping, UI affordances, and non-critical endpoints.

## Open questions

- Bulk batch import: confirm required CSV columns beyond `short_name` (default: `short_name` required, `brew_date` and `notes` optional) and whether duplicate `short_name` values should be rejected.
- Procurement purchase order numbers: confirm default generated format when a custom ID is not provided (recommend `PO-YYYYMMDD-###`).
- Procurement updates: confirm whether to add backend support for updating purchase orders and modifying/removing line items and fees.

## Implemented: Vessels Navigation Restructure

### Overview

The Vessels section now mirrors the Batches navigation pattern with two sub-views: Active (master/detail) and All Vessels (data table).

### Navigation structure

- **Vessels** (collapsible group)
  - **Active** (`/vessels/active`) - Master/detail view of active vessels
  - **All Vessels** (`/vessels/all`) - Data table view of all vessels

### Active Vessels page

- Shows only vessels with `status === 'active'`
- Sorted with occupied vessels first, then unoccupied, alphabetically by name within each group
- Master/detail layout with vessel list on left, details panel on right
- Details panel shows vessel info, metadata, and current occupancy with status management

### All Vessels page

- Data table with columns: ID, Name, Type, Capacity, Status, Occupancy, Updated
- Search filtering across all columns
- Status chips with colors: active (green), inactive (grey), retired (red)
- Occupancy column shows batch short_name as link when occupied, "Unoccupied" otherwise
- Default sorting: active first, then inactive, then retired; occupied before unoccupied within each status
- "New Vessel" button with create dialog
- Double-click row to navigate to vessel detail page

### Vessel detail page

- Barebones page at `/vessels/:uuid` showing vessel information
- Displays: name, type, status, capacity, make/model, timestamps
- Shows current occupancy section if vessel is occupied
- Foundation for future enhancements (occupancy history, cleaning logs, etc.)

### Shared utilities

- `useFormatters.ts` composable provides shared formatting functions:
  - `useFormatters()` - date/time formatting
  - `useVesselStatusFormatters()` - vessel status formatting and colors
  - `useOccupancyStatusFormatters()` - occupancy status formatting, colors, and icons

## Planned: Enhanced Batch Tracking

### Overview

Enhance batch tracking to capture the full brewing lifecycle from recipe through packaging, supporting spreadsheet-style views with derived metrics (days in fermentation, loss percentages, ABV calculations, etc.).

### New entities

**Style** - Reference table for beer styles (case-insensitive unique names). Enables autocomplete while allowing free-form entry.

**Recipe** - Beer formulation with name and style. Batches reference recipes. Future: target ABV/IBU, ingredient bills, process steps.

**Brew Session** - Captures hot-side wort production (mash → boil → knockout). Points to the wort volume produced. A batch can have multiple brew sessions (double batching). References mash and boil vessels for traceability.

### Entity modifications

**Batch** - Add `recipe_id` (nullable FK to recipe). The `brew_date` field becomes a convenience/derived field (earliest brew session date).

**Occupancy** - Add `status` field for granular liquid-in-vessel state: `fermenting`, `conditioning`, `cold_crashing`, `dry_hopping`, `carbonating`, `holding`, `packaging`.

**Addition** - Add `volume_id` as target option. Hot-side additions (malt, boil hops) attach to the wort volume. Cold-side additions attach to occupancy.

**Measurement** - Add `volume_id` as target option. Hot-side measurements (mash temp, OG) attach to wort volume. Cold-side measurements attach to occupancy.

**Volume units** - Add `bbl` (US barrel = 31 gallons).

### Key concepts

**Dual status tracking**: Batch process phases track high-level lifecycle (planning → finished). Occupancy status tracks granular liquid-in-vessel state (fermenting, conditioning, etc.).

**Hot-side vs cold-side tracking**: Brew session captures hot-side production; additions/measurements attach to the wort volume. Everything post-yeast-pitch is cold-side; additions/measurements attach to occupancies.

**ABV auto-calculation**: Compute from OG/FG measurements using `(OG - FG) × 131.25`, with manual override capability.

**Derived fields**: Days in fermentation, days in brite, days grain-to-glass, total loss, loss percentage - all computed from timestamps and volumes on read.

### Batch summary endpoint

New `GET /batches/{id}/summary` endpoint aggregates:
- Recipe name and style
- Brew session dates and volumes
- Current phase and vessel
- Key measurements (OG, FG, ABV, IBU)
- Duration metrics (days in FV, brite, grain-to-glass)
- Loss metrics (total loss BBL, loss percentage)

### Implementation milestones

1. Schema & Core API: style, recipe, brew_session tables; occupancy status; volume_id on addition/measurement; BBL unit
2. Batch Summary Endpoint: derived field calculations, ABV auto-calculation
3. UI - Recipe & Style Management: CRUD views, style autocomplete
4. UI - Brew Session Tracking: create/edit, hot-side additions/measurements
5. UI - Enhanced Batch Views: spreadsheet-style list, batch detail timeline
6. UI - Occupancy Status: status display and change workflow

## Implemented: User Settings (Brewery Name)

### Overview

Users can configure their brewery name, which is displayed prominently throughout the app. This setting is persisted in localStorage and defaults to "Acme Brewing" if not set.

### Display locations

- **Dashboard**: The brewery name appears as the main heading in the hero section
- **App Bar**: The brewery name appears below "BrewPipes" in the header

### Storage approach

- **Preferences**: Persisted in browser localStorage (`brewpipes:userSettings`)
- **Default**: "Acme Brewing"
- **Validation**: Empty/whitespace-only values are rejected; the field reverts to the last valid value

### Implementation

- Settings page accessible via user dropdown menu → "Settings"
- Composable: `useUserSettings.ts` (singleton state with localStorage persistence)
- Uses blur-based validation to prevent empty brewery names

## Implemented: User Display Preferences

### Overview

Users can customize their preferred units for displaying measurements throughout the app. Values are stored in their original units on the backend and converted to the user's preferred units for display.

### Supported unit types

- **Temperature**: Celsius (°C), Fahrenheit (°F)
- **Gravity**: Specific Gravity (SG), Degrees Plato (°P)
- **Volume**: mL, L, hL, US fl oz, UK fl oz, US gal, UK gal, US bbl (31 gal), UK bbl (36 gal)
- **Mass/Weight**: g, kg, oz, lb
- **Pressure**: kPa, PSI, bar
- **Color**: SRM, EBC

### Default preferences

Defaults are US-centric (first customer is US-based):
- Temperature: Fahrenheit
- Gravity: Specific Gravity
- Volume: US Barrels (bbl)
- Mass: Pounds (lb)
- Pressure: PSI
- Color: SRM

### Storage approach

- **Backend**: Stores values in their original recorded unit (preserves measurement fidelity)
- **Frontend**: Converts to user-preferred units on display
- **Preferences**: Persisted in browser localStorage (`brewpipes:unitPreferences`)

### Implementation

- Settings page accessible via user dropdown menu → "Settings"
- Composables: `useUnitConversion.ts` (pure conversion functions), `useUnitPreferences.ts` (preferences state and formatting)
- Applied across: Batches (sparklines, metrics), Vessels (capacity), Inventory (amounts), Dashboard
