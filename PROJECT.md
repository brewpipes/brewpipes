# BrewPipes Project Brief

## Product vision

BrewPipes is an open source brewery management system focused on day-to-day production operations. It aims to be the operational source of truth for procurement, inventory, and production, with clear traceability from ingredients and receipts through batches, vessels, and finished lots.

## V1 Roadmap

See [V1_ROADMAP.md](./V1_ROADMAP.md) for the complete V1 product roadmap including:
- Finalized user journeys
- Prioritized feature backlog
- Progress tracking

### V1 Target Users
Small craft breweries where 1-2 people wear multiple hats, managing both the business and performing brewery work.

### V1 Scope Decisions
- **In scope:** Procurement, brew day, fermentation, packaging, inventory, batch costing, brewhouse removals
- **Out of scope:** Taproom, merchandise, sales, multi-user/roles, multi-tenancy
- **Critical requirement:** Excellent mobile/tablet experience
- **Technical debt:** Must be addressed before new feature work

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

## API conventions

### UUID-only identifiers

All API endpoints use UUIDs as the sole external identifier for entities. Internal integer database primary keys are never exposed in API responses, request parameters, or URL paths.

- **Routes**: All entity routes use `{uuid}` path parameters (e.g., `/api/batches/{uuid}`, `/api/suppliers/{uuid}`)
- **Response DTOs**: Only the `uuid` field is included; no `id` field is returned
- **Foreign key references**: All FK fields in responses use `_uuid` suffix with string values (e.g., `batch_uuid`, `vessel_uuid`, `supplier_uuid`)
- **Request DTOs**: Create/update requests accept UUID strings for FK references (e.g., `recipe_uuid`, `supplier_uuid`)
- **Query parameters**: List endpoint filters use `_uuid` suffix (e.g., `?batch_uuid=xxx`, `?supplier_uuid=xxx`)
- **Frontend**: All entity types use `uuid: string` as the primary identifier; no numeric `id` fields exist in type definitions

## Implementation context

- Services live under `service/` and run independently or via the monolith entrypoint.
- The web app lives under `service/www/` and uses API clients that default to `/api` unless overridden by `VITE_*_API_URL`.
- Auth tokens are stored in the browser and used for Bearer requests across services.

## Deployment

### Single-container deployment

BrewPipes supports a single-container deployment model where the Go monolith serves both the API and the Vue frontend. This is the simplest deployment option for getting started.

**How it works:**
- The Vue frontend is built to `service/www/dist/`
- The Go binary embeds the frontend assets using `//go:embed`
- API routes are served at `/api/*`
- All other routes serve the embedded frontend with SPA fallback

**Building the container:**
```bash
docker build -f cmd/monolith/Dockerfile -t brewpipes .
```

**Running the container:**
```bash
docker run -p 8080:8080 \
  -e POSTGRES_DSN="postgres://user:pass@host:5432/brewpipes?sslmode=disable" \
  -e BREWPIPES_SECRET_KEY="your-secret-key" \
  brewpipes
```

**Image characteristics:**
- Multi-stage build (Node → Go → Alpine runtime)
- Final image size: ~45MB
- Static Go binary with embedded frontend assets
- No runtime dependencies beyond the binary itself

### API route structure

All backend API routes are prefixed with `/api`:
- Identity: `/api/login`, `/api/refresh`, `/api/logout`, `/api/users/*`
- Production: `/api/batches/*`, `/api/vessels/*`, `/api/recipes/*`, etc.
- Inventory: `/api/ingredients/*`, `/api/ingredient-lots/*`, etc.
- Procurement: `/api/suppliers/*`, `/api/purchase-orders/*`, etc.

The frontend is served at all non-API routes with SPA fallback (unknown paths serve `index.html` for client-side routing).

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

- Data table with columns: Name, Type, Capacity, Status, Occupancy, Updated
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

**Batch** - Add `recipe_uuid` (nullable FK to recipe). The `brew_date` field becomes a convenience/derived field (earliest brew session date).

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

New `GET /batches/{uuid}/summary` endpoint aggregates:
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

## Implemented: Inventory UX Simplification

### Overview

Simplified the Inventory module UX to better reflect how inventory data flows through the system and provide more intuitive workflows.

### Changes

**Activity Page** - Now read-only
- Removed the "Log activity" form entirely
- Inventory activity (movements) is derived from higher-level operations: ingredient usage, transfers between locations, receipt of purchase orders, etc.
- The page now displays a read-only activity log with filtering capabilities

**Locations Page** - Modal dialog pattern
- Converted the inline "Create Stock Location" form to a modal dialog
- Follows the established modal pattern used in Production (recipes) and Procurement (suppliers, purchase orders)
- "Create location" button in the page header opens the dialog
- Locations table now uses full page width

**Adjustments & Transfers Page** - Intuitive browse-then-act workflow
- Complete redesign replacing separate "create adjustment" and "create transfer" forms
- New workflow:
  1. Search for lots by name/ID, OR select a stock location to browse its inventory
  2. Unified inventory table shows all ingredient and beer lots with type, lot ID, name, quantity, and location
  3. Each row has "Adjust" and "Transfer" action buttons
- **Adjustment Modal**: Shows lot info (read-only), current quantity, adjustment amount (+/-), reason (required), notes, and timestamp
- **Transfer Modal**: Shows lot info (read-only), from location (read-only), to location dropdown, quantity to transfer (with validation against available quantity), notes, and timestamp
- Transfer validation prevents transferring more than available quantity

## Implemented: Ingredients Page Restructure

### Overview

Restructured the Ingredients page with a new tab layout focused on ingredient types and converted all inline forms to modal dialogs.

### New Tab Structure

| Tab | Content |
|-----|---------|
| **Malt** | Ingredient lots where parent ingredient has `category === 'fermentable'` |
| **Hops** | Ingredient lots where parent ingredient has `category === 'hop'` |
| **Yeast** | Ingredient lots where parent ingredient has `category === 'yeast'` |
| **Other** | Ingredient lots for non-core categories: `adjunct`, `salt`, `chemical`, `gas`, `other` |
| **Usage** | Ingredient usage records (grouped by batch reference) |
| **Received** | Inventory receipts from suppliers |

### Modal Dialogs

All forms converted to modal dialogs following established patterns:

- **Create Lot Modal**: Ingredient dropdown filtered by current tab's category, receipt reference, supplier info, lot codes, amounts, dates
- **Log Usage Modal**: Batch reference UUID, used at timestamp, notes
- **Create Receipt Modal**: Reference code, supplier UUID, received at timestamp, notes
- **Create Ingredient Modal**: Accessible via "New ingredient" button in page header; name, category, default unit, description

### Category Values Fixed

Frontend now uses correct backend category values:
- `fermentable` (not `malt`)
- `hop`
- `yeast`
- `adjunct`
- `salt` (not `water_chem`)
- `chemical`
- `gas`
- `other`

### Notes

- The "Other" tab serves as a catch-all for non-core ingredient categories. This can be refined into more specific tabs (e.g., Adjuncts, Water Chemistry) as needs evolve.
- The "Other" tab includes a Category column to help distinguish between different ingredient types within that tab.

## Implemented: Inventory Activity Page Enhancements

### Overview

Enhanced the Inventory Activity page with a richer data display that provides better context for each inventory movement.

### New Columns

| Column | Description |
|--------|-------------|
| **Item** | Ingredient name (resolved via lot → ingredient) or beer lot code |
| **Lot #** | Brewery lot code for ingredient lots, or lot code for beer lots |
| **Direction** | Colored icons: green down arrow for "in" (received), orange up arrow for "out" (used/transferred/adjusted). Includes tooltips. |
| **Reason** | Context-aware display based on movement type (see below) |
| **Amount** | Formatted with user's preferred units |
| **Location** | Stock location name |

### Reason Column Logic

The Reason column displays context-aware information based on the movement type:

- **Usage (`reason === 'use'`)**: Shows "Used in [BATCH_SHORT_NAME]" with a link to the batch page. Falls back to "Used in production" if batch cannot be resolved.
- **Receipt (`reason === 'receive'`)**: Shows "Received from [SUPPLIER_NAME]" when supplier can be resolved. Falls back to "Received" otherwise.
- **Adjustment (`reason === 'adjust'` or `'waste'`)**: Shows formatted adjustment reason (Cycle count, Spoilage, Shrink, Damage, Correction, Other). Notes displayed in tooltip if available.
- **Transfer (`reason === 'transfer'`)**: Shows "Transferred to [LOCATION]" or "Transferred from [LOCATION]" based on direction. Notes displayed in tooltip if available.

### Data Fetching

The page fetches data from multiple sources:

**Inventory service:**
- Ingredients (to resolve ingredient names)
- Ingredient lots, beer lots, stock locations
- Inventory movements, receipts, usages, adjustments, transfers

**Cross-service calls (graceful failure):**
- Production service: batches (to resolve batch names for usage movements)
- Procurement service: suppliers (to resolve supplier names for receipts)

### Implementation Notes

- Uses computed lookup maps for O(1) data resolution
- Cross-service calls use `Promise.allSettled` for non-blocking, graceful failure
- Graceful fallbacks when related data cannot be resolved
- Filter dropdowns now show ingredient name alongside lot code for better identification

## Planned: Phase 4 — Brew Day & Occupancy

### Overview

Phase 4 enables complete brew day recording with inventory integration. A brewer can pull ingredients from inventory using a recipe-driven pick list, record a brew session, assign beer to a fermenter, and manage occupancy lifecycle — all optimized for phone use on the brew floor.

### Key Design Decisions

- **New "Brew Day" tab in batch details** — The pick list and brew session recording live in a dedicated tab alongside Summary, Brew Sessions, Timeline, etc. This keeps brew-day workflows consolidated without cluttering the existing tabs.
- **Features first, wizard later** — Individual features (pick list, lot selection, occupancy management) are delivered independently before being unified into a mobile-optimized brew day wizard (BREW-06).
- **Recipe scaling (RCP-06) included** — Moved from Phase 2 to Phase 4 since scaling directly affects pick list amounts.
- **Offline resilience deferred** — WiFi in breweries is often spotty, but offline support is deferred to a later polish phase.
- **User-facing language avoids "occupancy"** — From the brewer's perspective, they assign beer to fermenters, transfer between vessels, and mark vessels as empty. The occupancy data model is an implementation detail. UI labels use brewer-centric language ("Currently in: FV-3", "Assign to Fermenter", "Mark Empty") while the backend continues to use occupancy entities internally.

### Brew Day Workflow

1. **Review & Pick** — Pick list generated from recipe ingredient bill, grouped by use stage ("Needed Today" vs. "Needed Later"). FIFO lot suggestion with manual override. Lot-specific attributes shown (alpha acid for hops, best-by dates).
2. **Confirm Picks** — Inventory deducted atomically when brewer confirms the pick list. Stock validation prevents over-deduction.
3. **Record Session** — Create brew session with mash/boil vessel assignment. Record key measurements (mash temp, OG, volume to fermenter).
4. **Assign to Fermenter** — Transfer wort to fermenter at knockout. Beer status defaults to `fermenting`. Only active, available vessels shown. (Backend creates an occupancy behind the scenes.)
5. **Mark Vessel as Empty** — Record that a vessel has been emptied (for corrections or dumps). The common case — transferring beer to another vessel or packaging — is handled in Phases 5-6 and automatically manages vessel availability.

### Inventory Deduction Strategy

- **Mash/boil ingredients:** Deducted when brewer confirms the pick list (batch-level `inventory_usage`)
- **Post-brew additions (dry hops, finings):** Deducted when recorded against the occupancy (occupancy-level, handled in Phase 5+)
- **Cross-service orchestration:** Frontend orchestrates Production (additions) and Inventory (usage/movements) calls. Backend validates stock atomically within transactions.

### Mobile UX Principles

- Phone-first, one-handed operation (brewer has one hand on equipment)
- Large tap targets (44px minimum), numeric keypads for measurements
- Under 15 seconds per data entry interaction
- Step-based cards, not long scrolling forms
- Non-linear — brewer can skip steps and return

### New/Modified API Endpoints

| Method | Path | Service | Description |
|--------|------|---------|-------------|
| `GET` | `/api/ingredient-lots?ingredient_uuid=...` | Inventory | Filter lots by ingredient (new filter) |
| `POST` | `/api/inventory-usage/batch` | Inventory | Batch-level inventory deduction (new) |
| `PATCH` | `/api/occupancies/{uuid}/close` | Production | End/close occupancy (new) |
| `POST` | `/api/occupancies` | Production | Create occupancy (exists) |

### Implementation Order

1. BREW-04 (Assign to Fermenter) + BREW-05 (Mark Vessel as Empty) — simplest, no dependencies
2. BREW-01 (Pick List) + BREW-03 (Inventory Deduction) — backend work in parallel
3. BREW-02 (Lot Selection) — depends on BREW-01 + BREW-03
4. RCP-06 (Recipe Scaling) — enhances pick list with scaled amounts
5. BREW-06 (Mobile Wizard) — capstone that unifies the workflow

## Implemented: Phase 5 — Fermentation & Transfers

### Overview

Phase 5 enables fermentation monitoring and complex vessel-to-vessel transfers. Brewers can monitor all active fermentations at a glance, visualize fermentation curves, quickly record daily readings from their phone, and transfer beer between vessels with support for splits and blends.

### New pages and navigation

- **Fermentation Dashboard** (`/fermentation`) — New top-level page showing all active fermentations with tank cards, sparklines, and attention indicators. Accessible from the main sidebar navigation.

### New components

- `FermentationCard` — Per-tank card with vessel name, batch, status, gravity/temp sparklines, attenuation, ABV, and action buttons
- `FermentationCurve` — Chart.js dual-axis line chart (gravity + temperature over time) with OG/FG reference lines
- `QuickReadingSheet` — Mobile-optimized bottom sheet for fast measurement entry (gravity, temperature, pH)
- `TransferDialog` — 2-step wizard supporting three modes: Transfer (1:1), Split (1:N), and Blend (N:1)

### New batch detail tab

- **Fermentation** tab added between "Brew Sessions" and "Timeline" — shows fermentation curve chart with toggle controls and stats summary

### Transfer workflow

- **Simple transfer:** Source vessel → destination vessel with volume, loss, and status
- **Split:** One source → 2-4 destinations with volume math validation
- **Blend:** 2-4 sources → one destination with batch identity selection
- All modes create appropriate volumes, transfers, volume relations, and batch-volume records

### Backend enhancements

- `POST /api/transfers` enhanced with `close_source` (boolean) and `dest_status` (occupancy status) parameters
- `GET /api/batches/{uuid}/summary` now includes `current_occupancy_uuid`
- New composite index on `measurement(occupancy_id, kind, observed_at DESC)` for dashboard performance

### Frontend infrastructure

- 17 new API wrapper functions in `useProductionApi`
- 18 new TypeScript types in `types/production.ts`
- Chart.js + vue-chartjs + chartjs-plugin-annotation for fermentation curve
- `BatchDetails.vue` refactored to use typed composable functions
- 83 new frontend tests (468 total across 15 test files)

## Implemented: Phase 6 — Packaging & Finished Goods

### Overview

Phase 6 enables packaging run recording and finished goods inventory tracking. Brewers can record packaging events from a brite tank, specifying container formats and quantities. Beer lots are automatically created in the inventory service via an inter-service HTTP call, with initial stock movements.

### New entities

**Package Format** (Production) — User-extensible reference table for container types (1/2 BBL Keg, 16oz Can, etc.) with standard fill volumes.

**Packaging Run** (Production) — Records a packaging event for a batch from a specific brite tank occupancy. Has line items per container format with quantities.

**Beer Lot Item** (Inventory) — Structural placeholder for future per-unit tracking (individual keg codes, barcodes). Table created but not populated by the packaging workflow.

### Beer lot enhancements

The `beer_lot` table gained: `packaging_run_uuid`, `best_by`, `package_format_name`, `container`, `volume_per_unit`, `volume_per_unit_unit`, `quantity`. One beer lot per format per packaging run.

### Inter-service communication

First backend-to-backend HTTP call in the system. Production service calls Inventory service's `POST /api/beer-lots` endpoint with JWT pass-through from the original request. Beer lot creation is best-effort — the packaging run is committed even if the inventory call fails.

### New API endpoints

| Method | Path | Service | Description |
|--------|------|---------|-------------|
| `GET` | `/api/package-formats` | Production | List package formats |
| `POST` | `/api/package-formats` | Production | Create package format |
| `GET` | `/api/package-formats/{uuid}` | Production | Get package format |
| `PATCH` | `/api/package-formats/{uuid}` | Production | Update package format |
| `DELETE` | `/api/package-formats/{uuid}` | Production | Delete package format (blocked if in use) |
| `GET` | `/api/packaging-runs` | Production | List packaging runs (optional `?batch_uuid=` filter) |
| `POST` | `/api/packaging-runs` | Production | Create packaging run with lines |
| `GET` | `/api/packaging-runs/{uuid}` | Production | Get packaging run with lines |
| `DELETE` | `/api/packaging-runs/{uuid}` | Production | Soft-delete packaging run |
| `GET` | `/api/beer-lot-stock-levels` | Inventory | Finished goods stock levels |

### Frontend components

- **PackagingDialog** — 3-step wizard for recording packaging runs (details → format lines → review)
- **BestByIndicator** — Reusable component for best-by date display with expiry indicators
- **Product page** — Redesigned with Stock Levels and All Lots tabs, container filtering, best-by indicators
