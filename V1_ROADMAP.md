# BrewPipes V1 Product Roadmap

**Last Updated:** 2026-02-14  
**Status:** Phase 4 Complete

---

## Overview

This document defines the V1 product roadmap for BrewPipes, an open source brewery management system. It serves as the living plan for delivering a production-ready V1 release.

**Target Users:** Small craft breweries where 1-2 people wear multiple hats, managing both the business and performing brewery work.

**V1 Goal:** Enable a brewery to manage the complete production lifecycle from procurement through packaging, with full ingredient traceability, accurate batch costing, and a great experience on mobile, tablet, and desktop.

**Out of Scope for V1:** Taproom management, merchandise, sales, multi-user/roles, multi-tenancy.

---

## Key Decisions

| Decision | Answer |
|----------|--------|
| **Primary Persona** | Small brewery owner/operator wearing multiple hats |
| **Mobile/Tablet** | **Critical** - must be excellent on all device sizes |
| **Recipe Depth** | Full ingredient bills with traceability to POs |
| **Packaging** | In scope - finished goods inventory tracking |
| **Cost Tracking** | In scope - accurate batch costing |
| **Brewhouse Removals** | In scope - track dumps, waste, spillage |
| **TTB Compliance** | Design for it, but full reporting is V2 |
| **Multi-user** | Defer to V2 (single user acceptable) |
| **Technical Debt** | **Priority #1** - address before new features |

---

## User Journeys

These are the core user journeys that must work flawlessly for V1. Each journey represents a complete workflow that delivers value to the user.

### Journey 1: Procurement & Receiving

**Goal:** Order ingredients from suppliers and receive them into inventory with full traceability.

**Flow:**
1. Create or select a supplier
2. Create a purchase order with line items (ingredients, quantities, costs)
3. Submit/confirm the PO
4. When delivery arrives, open the PO
5. Receive items against PO line items:
   - Create ingredient lots with supplier lot codes, brewery lot codes, best-by dates
   - Assign storage locations
   - Record actual quantities received
6. PO status updates automatically (partially received → received)
7. Inventory levels reflect new stock
8. Full traceability: PO → Receipt → Lot → (future) Batch

**Current Gaps:**
- No PO → Receipt linkage workflow in UI
- No current stock level display
- PO status doesn't auto-update on receipt
- Receiving workflow is fragmented

**Success Criteria:**
- User can create PO, receive against it, and see inventory updated
- Ingredient lots link back to PO line items
- Stock levels visible at a glance

---

### Journey 2: Brew Day Execution

**Goal:** Execute and record a complete brew day from grain to fermenter.

**Flow:**
1. Select batch scheduled for today (or create new batch from recipe)
2. View recipe with full ingredient bill and target specs
3. Pick ingredients from inventory:
   - See available lots for each ingredient
   - Select lots to use (FIFO suggested)
   - Deduct from inventory
4. Create brew session and assign mash/boil vessels
5. Record mash parameters (strike temp, mash pH, rest temps)
6. Record boil additions with timestamps (hops, finings, etc.)
7. Record knockout measurements (OG, volume to fermenter)
8. Create wort volume and assign to fermenter vessel (create occupancy)
9. Pitch yeast (record as addition)
10. Batch transitions to "fermenting" phase

**Current Gaps:**
- Recipe lacks ingredient bills
- No ingredient pick list generation
- No automatic inventory deduction on pick
- Occupancy creation requires backend knowledge
- Mobile UX not optimized for brew day

**Success Criteria:**
- User can execute entire brew day on a tablet
- Ingredients deducted from correct lots
- Batch in fermenter with OG recorded
- Full traceability from ingredient lots to batch

---

### Journey 3: Fermentation Management

**Goal:** Monitor fermentation, make additions, and manage transfers including splits and blends.

**Flow:**
1. View active fermentations (all batches currently fermenting)
2. Select a batch to monitor
3. Record daily measurements:
   - Gravity (track attenuation)
   - Temperature
   - pH (optional)
4. View fermentation curve (gravity over time)
5. Record cold-side additions:
   - Dry hops (with lot traceability)
   - Fruit, adjuncts, finings, etc.
6. Update occupancy status as needed (fermenting → conditioning → carbonating)
7. When ready, initiate transfer:
   - Simple transfer to brite tank
   - Split batch across multiple vessels
   - Blend with another batch
8. Record transfer volumes and losses
9. End source occupancy, create destination occupancy(ies)

**Current Gaps:**
- No dedicated active fermentations view
- No fermentation curve visualization
- Volume split/blend UI doesn't exist
- Transfer workflow is complex

**Success Criteria:**
- User can monitor all active fermentations at a glance
- Fermentation curves show progress visually
- Splits and blends are intuitive to create
- Volume losses tracked accurately

---

### Journey 4: Packaging & Finished Goods

**Goal:** Package beer, record final specs, create finished goods inventory, and track packaging loss.

**Flow:**
1. Select batch ready for packaging
2. Record final measurements:
   - Final gravity (FG)
   - ABV (calculated or measured)
   - IBU
   - Color (SRM)
   - Carbonation level
3. Initiate packaging run:
   - Select package types (kegs, cans, bottles)
   - Record quantities packaged
   - Record packaging date
4. Create beer lots (finished goods) linked to batch
5. Calculate and record packaging loss
6. Move finished goods to inventory locations
7. Batch transitions to "finished" phase
8. End final occupancy

**Current Gaps:**
- Beer lot creation not exposed in UI
- No packaging run recording
- No finished goods inventory view
- Packaging loss not calculated

**Success Criteria:**
- User can record complete packaging run
- Beer lots created with batch traceability
- Packaging loss calculated and recorded
- Finished goods visible in inventory

---

### Journey 5: Batch Costing & Review

**Goal:** Review batch performance, verify quality, and understand true costs.

**Flow:**
1. Select any batch (in-progress or finished)
2. View batch summary dashboard:
   - Recipe and style
   - Key metrics: OG, FG, ABV, IBU, color
   - Duration: days in FV, days in brite, grain-to-glass
   - Volumes: starting, current, total loss, loss %
3. View cost breakdown:
   - Ingredient costs (from PO line items via lots)
   - Cost per barrel
   - Cost per package (if packaged)
4. Compare actual vs. recipe targets
5. Review full timeline:
   - All measurements
   - All additions
   - All transfers
   - Phase changes
6. Review ingredient traceability:
   - Which lots were used
   - Link back to suppliers and POs
7. Add notes and observations

**Current Gaps:**
- Recipe lacks target specs for comparison
- Cost calculation not implemented
- Ingredient traceability UI incomplete

**Success Criteria:**
- User can see complete batch story at a glance
- Costs calculated accurately from PO data
- Actual vs. target comparison available
- Full traceability visible

---

### Journey 6: Inventory Management & Removals

**Goal:** Maintain accurate inventory, track stock levels, and record brewhouse removals.

**Flow:**
1. View current stock levels by ingredient category
2. See low stock alerts
3. Adjust inventory as needed:
   - Cycle count corrections
   - Spoilage/damage
   - Transfer between locations
4. Record brewhouse removals:
   - Batch dumps (failed batches)
   - Tank cleaning waste
   - Sample pulls
   - Spillage
5. View inventory activity history
6. Trace any lot back to its origin (PO, supplier)

**Current Gaps:**
- No current stock level display
- No low stock alerts
- Removal tracking for beer not implemented
- Beer lot inventory management incomplete

**Success Criteria:**
- User can see what's in stock at a glance
- Low stock items highlighted
- All removals tracked for future TTB compliance
- Full audit trail available

---

## Features

Features are organized into phases. **Phase 0 (Technical Debt) must be completed before any feature work.**

Each feature is sized to be deliverable in a focused session:
- **S (Small):** 1-2 hours, single component or endpoint
- **M (Medium):** 2-4 hours, multiple components or cross-cutting
- **L (Large):** 4-8 hours, significant new capability

---

### Phase 0: Technical Debt (PRIORITY #1)

These items must be completed before any new feature work.

| ID | Feature | Size | Status |
|----|---------|------|--------|
| TD-01 | Decompose BatchDetails.vue into smaller components | L | **Complete** |
| TD-02 | Add frontend unit tests for API composables | M | **Complete** |
| TD-03 | Add frontend unit tests for utility composables | M | **Complete** |
| TD-04 | Consolidate duplicate TypeScript types | S | **Complete** |
| TD-05 | Audit and fix mobile responsiveness issues | M | **Complete** |
| TD-06 | Add frontend component tests for critical flows | L | **Complete** |

---

### Phase 1: Core CRUD & UX Completion

Complete basic operations that are currently missing.

| ID | Feature | Size | Journey | Status |
|----|---------|------|---------|--------|
| CRUD-01 | Add batch editing (name, brew date, recipe, notes) | M | 2, 5 | **Complete** |
| CRUD-02 | Add batch deletion with confirmation and checks | S | 2 | **Complete** |
| CRUD-03 | Add vessel editing | S | 2, 3 | **Complete** |
| CRUD-04 | Add vessel retirement workflow | S | - | **Complete** |
| CRUD-05 | Add recipe deletion with batch reference check | S | 2 | **Complete** |
| CRUD-06 | Add supplier editing | S | 1 | **Complete** |
| CRUD-07 | Add purchase order editing | M | 1 | **Complete** |

---

### Phase 2: Recipe Enhancement

Enable full recipe management with ingredient bills.

| ID | Feature | Size | Journey | Status |
|----|---------|------|---------|--------|
| RCP-01 | Design recipe ingredient bill data model | M | 2 | **Complete** |
| RCP-02 | Add recipe target specifications (OG, FG, ABV, IBU, color) | M | 2, 5 | **Complete** |
| RCP-03 | Backend: Recipe ingredient bill CRUD endpoints | M | 2 | **Complete** |
| RCP-04 | Frontend: Recipe ingredient bill management UI | L | 2 | **Complete** |
| RCP-05 | Frontend: Recipe detail view with full specs | M | 2, 5 | **Complete** |
| RCP-06 | Recipe scaling calculator | M | 2 | Moved to Phase 4 |

---

### Phase 3: Procurement & Receiving Workflow

Complete the procurement → inventory flow.

| ID | Feature | Size | Journey | Status |
|----|---------|------|---------|--------|
| PROC-01 | Backend: Link receipts to PO line items | M | 1 | **Complete** |
| PROC-02 | Backend: Auto-update PO status on receipt | S | 1 | **Complete** |
| PROC-03 | Frontend: PO detail view with line items inline | M | 1 | **Complete** |
| PROC-04 | Frontend: Receiving workflow (receive against PO) | L | 1 | **Complete** |
| PROC-05 | Frontend: Current stock levels display | M | 1, 6 | **Complete** |
| PROC-06 | Frontend: Low stock alerts | S | 6 | **Complete** |

---

### Phase 4: Brew Day & Occupancy

Enable complete brew day recording with inventory integration.

| ID | Feature | Size | Journey | Status |
|----|---------|------|---------|--------|
| BREW-04 | Frontend: Assign batch to fermenter | M | 2, 3 | **Complete** |
| BREW-05 | Backend + Frontend: Mark vessel as empty | S | 3, 4 | **Complete** |
| BREW-01 | Backend + Frontend: Ingredient pick list from recipe | M | 2 | **Complete** |
| BREW-03 | Backend: Inventory deduction on ingredient use | M | 2 | **Complete** |
| BREW-02 | Frontend: Lot selection for ingredient picks | M | 2 | **Complete** |
| RCP-06 | Recipe scaling calculator (from Phase 2) | M | 2 | **Complete** |
| BREW-06 | Mobile-optimized brew day wizard | L | 2 | **Complete** |

---

### Phase 5: Fermentation & Transfers

Enable fermentation monitoring and complex transfers.

| ID | Feature | Size | Journey | Status |
|----|---------|------|---------|--------|
| FERM-01 | Frontend: Active fermentations dashboard | M | 3 | Not Started |
| FERM-02 | Frontend: Fermentation curve visualization | M | 3 | Not Started |
| FERM-03 | Frontend: Transfer wizard (simple transfer) | M | 3 | Not Started |
| FERM-04 | Frontend: Volume split UI | L | 3 | Not Started |
| FERM-05 | Frontend: Volume blend UI | L | 3 | Not Started |
| FERM-06 | Frontend: Quick measurement entry (mobile-friendly) | M | 3 | Not Started |

---

### Phase 6: Packaging & Finished Goods

Enable packaging recording and finished goods inventory.

| ID | Feature | Size | Journey | Status |
|----|---------|------|---------|--------|
| PKG-01 | Design packaging run data model | M | 4 | Not Started |
| PKG-02 | Backend: Packaging run CRUD endpoints | M | 4 | Not Started |
| PKG-03 | Frontend: Packaging run recording UI | L | 4 | Not Started |
| PKG-04 | Backend: Beer lot creation from packaging | M | 4 | Not Started |
| PKG-05 | Frontend: Beer lot / finished goods inventory view | M | 4, 6 | Not Started |
| PKG-06 | Packaging loss calculation | S | 4 | Not Started |

---

### Phase 7: Cost Tracking

Enable batch costing from procurement data.

| ID | Feature | Size | Journey | Status |
|----|---------|------|---------|--------|
| COST-01 | Backend: Calculate ingredient cost per batch | M | 5 | Not Started |
| COST-02 | Backend: Calculate cost per barrel | S | 5 | Not Started |
| COST-03 | Frontend: Batch cost breakdown view | M | 5 | Not Started |
| COST-04 | Frontend: Actual vs. target comparison | M | 5 | Not Started |

---

### Phase 8: Removals & Compliance Prep

Track brewhouse removals for future TTB compliance.

| ID | Feature | Size | Journey | Status |
|----|---------|------|---------|--------|
| REM-01 | Design removal tracking data model | M | 6 | Not Started |
| REM-02 | Backend: Removal CRUD endpoints | M | 6 | Not Started |
| REM-03 | Frontend: Record batch dump | S | 6 | Not Started |
| REM-04 | Frontend: Record waste/spillage | S | 6 | Not Started |
| REM-05 | Frontend: Record sample pulls | S | 6 | Not Started |
| REM-06 | Removal history and reporting | M | 6 | Not Started |

---

## Progress Tracking

### Milestones

| Milestone | Description | Status |
|-----------|-------------|--------|
| **M0: Roadmap Complete** | User journeys and feature backlog finalized | **Complete** |
| **M1: Tech Debt Clear** | Phase 0 complete, codebase ready for features | **Complete** |
| **M2: Core Complete** | Phases 1-2 complete (CRUD, recipes) | **Complete** (RCP-06 delivered in Phase 4) |
| **M3: Procurement Flow** | Phase 3 complete (PO → inventory) | **Complete** |
| **M4: Brew Day Flow** | Phase 4 complete (brew day recording) | **Complete** |
| **M5: Fermentation Flow** | Phase 5 complete (monitoring, transfers) | Not Started |
| **M6: Packaging Flow** | Phase 6 complete (packaging, finished goods) | Not Started |
| **M7: Costing Complete** | Phase 7 complete (batch costing) | Not Started |
| **M8: V1 Alpha** | All phases complete, internal testing | Not Started |
| **M9: V1 Beta** | External testing, bug fixes | Not Started |
| **M10: V1 Release** | Production ready | Not Started |

### Completed Features

| ID | Feature | Date |
|----|---------|------|
| TD-01 | Decompose BatchDetails.vue into smaller components | 2026-02-02 |
| TD-02 | Add frontend unit tests for API composables | 2026-02-02 |
| TD-03 | Add frontend unit tests for utility composables | 2026-02-02 |
| TD-04 | Consolidate duplicate TypeScript types | 2026-02-02 |
| TD-05 | Audit and fix mobile responsiveness issues | 2026-02-02 |
| TD-06 | Add frontend component tests for critical flows | 2026-02-02 |
| CRUD-01 | Add batch editing (name, brew date, recipe, notes) | 2026-02-02 |
| CRUD-02 | Add batch deletion with confirmation and checks | 2026-02-02 |
| CRUD-03 | Add vessel editing | 2026-02-02 |
| CRUD-04 | Add vessel retirement workflow | 2026-02-02 |
| CRUD-05 | Add recipe deletion with batch reference check | 2026-02-02 |
| CRUD-06 | Add supplier editing | 2026-02-02 |
| CRUD-07 | Add purchase order editing | 2026-02-02 |
| RCP-01 | Design recipe ingredient bill data model | 2026-02-02 |
| RCP-02 | Add recipe target specifications (OG, FG, ABV, IBU, color) | 2026-02-02 |
| RCP-03 | Backend: Recipe ingredient bill CRUD endpoints | 2026-02-02 |
| RCP-04 | Frontend: Recipe ingredient bill management UI | 2026-02-02 |
| RCP-05 | Frontend: Recipe detail view with full specs | 2026-02-02 |
| PROC-01 | Backend: Link receipts to PO line items | 2026-02-13 |
| PROC-02 | Backend: Auto-update PO status on receipt | 2026-02-13 |
| PROC-03 | Frontend: PO detail view with line items inline | 2026-02-13 |
| PROC-04 | Frontend: Receiving workflow (receive against PO) | 2026-02-13 |
| PROC-05 | Frontend: Current stock levels display | 2026-02-13 |
| PROC-06 | Frontend: Low stock alerts | 2026-02-13 |
| BREW-04 | Frontend: Assign batch to fermenter | 2026-02-14 |
| BREW-05 | Backend + Frontend: Mark vessel as empty | 2026-02-14 |
| BREW-01 | Backend + Frontend: Ingredient pick list from recipe | 2026-02-14 |
| BREW-03 | Backend: Inventory deduction on ingredient use | 2026-02-14 |
| BREW-02 | Frontend: Lot selection for ingredient picks | 2026-02-14 |
| RCP-06 | Recipe scaling calculator | 2026-02-14 |
| BREW-06 | Mobile-optimized brew day wizard | 2026-02-14 |

**TD-01 Details:** Refactored 3,227-line component into 15 smaller components in `service/www/src/components/batch/`. Main component reduced to 1,677 lines (~48% reduction). Created 6 tab components, 7 dialog components, 1 reusable card component, shared types file, and barrel export.

**TD-02 Details:** Created 4 test files with 97 tests covering `useApiClient`, `useProductionApi`, `useInventoryApi`, `useProcurementApi`. Set up Vitest with happy-dom environment.

**TD-03 Details:** Created 4 test files with 198 tests covering `useFormatters`, `useUnitConversion`, `useUnitPreferences`, `useUserSettings`. Total test suite now has 295 passing tests.

**TD-04 Details:** Created `service/www/src/types/` directory with 6 files organizing domain types (common, units, production, settings, auth). Consolidated duplicate types including `VolumeUnit` (4→9 values). Maintained backward compatibility via re-exports from composables.

**TD-05 Details:** Modified 12 files for mobile responsiveness. Implemented master-detail mobile pattern (list OR detail on mobile), responsive dialogs, icon-only buttons on xs, 44px touch targets, table horizontal scrolling. App now works well on phones, tablets, and desktops.

**TD-06 Details:** Created 3 component test files with 36 tests covering `BatchList`, `VesselList`, and `AppFooter`. Tests cover rendering, selection, events, empty states, and sorting. Total test suite now has 331 passing tests. Updated Vitest config for Vuetify component testing.

**Phase 1 Details:** Implemented full CRUD operations for core entities:
- **Batch Edit/Delete (CRUD-01, CRUD-02):** Backend PATCH/DELETE endpoints with dependency checking (blocks deletion if batch has volumes, phases, sessions, additions, or measurements). Frontend edit dialog with recipe selector, delete confirmation dialog with 409 conflict handling. Accessible from batch detail page and batch list.
- **Vessel Edit/Retire (CRUD-03, CRUD-04):** Backend PATCH endpoint with active occupancy check (blocks retirement if vessel is occupied). Frontend edit dialog with status dropdown and retirement warning. Accessible from vessel detail, active vessels, and all vessels pages.
- **Recipe Delete (CRUD-05):** Backend DELETE endpoint with batch reference check (blocks deletion if batches use the recipe). Frontend delete action with confirmation dialog. Accessible from recipes page.
- **Supplier Edit (CRUD-06):** Backend PATCH endpoint with partial update semantics. Frontend unified create/edit dialog. Accessible from suppliers page.
- **Purchase Order Edit (CRUD-07):** Frontend edit functionality added to existing backend PATCH endpoint. Unified create/edit dialog with supplier lock (cannot change supplier after creation). Accessible from purchase orders page.

**Phase 2 Details:** Implemented recipe ingredient bills and target specifications:
- **Data Model (RCP-01):** Created `recipe_ingredient` table with support for fermentables, hops, yeast, adjuncts, water chemistry, and other additions. Each ingredient has amount, unit, use stage (mash/boil/whirlpool/fermentation/packaging), use type, timing, and notes. Hops include alpha acid for IBU calculations. Scaling factor supports non-linear scaling.
- **Target Specs (RCP-02):** Extended `recipe` table with batch_size, target_og/fg/ibu/srm (with optional min/max ranges), target_carbonation, ibu_method, and brewhouse_efficiency. ABV is calculated from OG/FG using `(OG - FG) × 131.25`.
- **Backend CRUD (RCP-03):** Full CRUD endpoints for recipe ingredients at `/recipes/{id}/ingredients`. Validation for ingredient types, use stages, alpha acid (hops only), and positive amounts.
- **Frontend UI (RCP-04, RCP-05):** Recipe detail page at `/recipes/{uuid}` with tabs for Overview, Fermentables, Hops, Yeast & Other, and Specs. Ingredient management via modal dialogs. Mobile-responsive with card layouts on small screens. SRM color preview swatch. Spec badges in recipe list.

**Phase 3 Details:** Implemented complete procurement-to-inventory receiving workflow:
- **PO Detail Page (PROC-03):** Consolidated PO view with header, line items, fees, and totals. Status management with valid transition enforcement. Mobile-responsive with card layouts.
- **Receiving Workflow (PROC-04):** 3-step wizard for receiving against PO (select lines → enter details → confirm). Ad-hoc receiving dialog for inventory without PO. Creates receipts, lots, and movements atomically. Auto-updates PO status.
- **Stock Levels (PROC-05):** New stock levels page with category tabs, per-location breakdown, and zero-stock indicators. Backend endpoint aggregates movements into current stock.
- **Low Stock Alerts (PROC-06):** Dashboard card showing out-of-stock ingredients. Navigation badge on Inventory when low stock exists.
- **Backend Enhancements (PROC-01, PROC-02):** PO linkage exposed on receipts and lots. PO status transition validation with state machine.

**Phase 4 Details:** Implemented complete brew day recording with inventory integration:
- **Assign to Fermenter (BREW-04):** Frontend dialog for creating occupancies from batch detail. Visual vessel selection with capacity display. "Beer Status" label replaces "Occupancy Status" for brewer-friendly language.
- **Mark Vessel Empty (BREW-05):** Backend PATCH endpoint for closing occupancies. Frontend dialog accessible from batch summary. Validates occupancy isn't already closed.
- **Ingredient Pick List (BREW-01):** Brew Day tab in batch details showing recipe ingredients grouped by "Needed Today"/"Needed Later" based on use stage. Cross-service stock level resolution. FIFO lot ordering with soft-delete filtering.
- **Inventory Deduction (BREW-03):** Atomic batch-level inventory deduction endpoint (POST /inventory-usage/batch). Transactional stock validation prevents over-deduction. 22 backend tests.
- **Lot Selection (BREW-02):** Per-ingredient lot selection with FIFO ordering, multi-lot splitting, confirmation dialog. Backend enhanced with `current_amount` computed from movement ledger. 7 frontend tests.
- **Recipe Scaling (RCP-06):** `useRecipeScaling` composable with cross-unit volume conversion (bbl/gal/L/hL). Scaling controls in recipe details and brew day tab. 23 tests.
- **Brew Day Wizard (BREW-06):** Three-step guided wizard (pick ingredients → record session → assign fermenter) optimized for mobile. Fullscreen on small screens, non-linear step navigation, auto-fill FIFO lots, visual fermenter card picker, completion summary. Reviewed and refined by frontend tech lead.
- **Seed Data Enhancement:** 20 new ingredients, 22 lots, 4 new batches, 5 new vessels, missing enum values, FIFO test lots, cross-service UUID fix migration.

### Deferred QA Items

Issues identified during domain-by-domain QA review passes that were deferred for future work. Items will be added here as QA passes continue.

#### Dashboard (Pass 1)
- [ ] Low stock ingredient link auto-selects correct tab but doesn't highlight/scroll to the specific ingredient row
- [ ] Footer overlaps content on mobile at 375px width

#### Production — Batches (Pass 2)
- [ ] OG/FG show "—" in batch summary even when gravity measurements exist — backend `populateMeasurements` may not be extracting OG/FG from measurements correctly
- [ ] Inventory lot UUID fields in addition/measurement dialogs require raw UUID input — should be searchable comboboxes with lot codes and ingredient names
- [ ] Measurement "Kind" field inconsistency — hot-side dialog uses curated dropdown, batch-level dialog uses free text input
- [ ] Temperature unit inconsistency — sparkline converts to user preference (°F), measurement table shows raw values (°C)

#### Production — Recipes (Pass 3)
- [ ] Vue "Scroll target is not found" warnings when resizing to mobile width on recipe detail pages

#### Vessels (Pass 4)
- [ ] Retire button opens generic Edit dialog instead of a dedicated retire confirmation — should pre-select "Retired" status or use a dedicated confirmation dialog
- [ ] Vessel detail page header clips Active status chip on mobile at 375px width

#### Inventory (Pass 5)
- [ ] Adjustments & Transfers page shows "Unknown Location" for all items — backend `ingredient-lots` API doesn't return `stock_location_uuid`; needs backend change or frontend refactor to derive from stock-levels API
- [ ] Activity page "Unknown Lot" entry for soft-deleted lot — backend should include deleted lot info or show last-known name
- [ ] Transfer dialog max quantity shows wrong unit and value — uses raw `current_amount` (kg) before unit conversion instead of converted display value
- [ ] Activity page has no mobile-responsive layout at 375px — table columns cut off with no horizontal scroll; needs card layout redesign
- [ ] Product page create form requires raw batch UUID input — should be a searchable combobox listing available batches
- [ ] Lot Details page shows all 3 detail sections (Malt, Hop, Yeast) for every lot type regardless of category
- [ ] Locations page has no edit/delete actions on location rows — can create but not modify or remove
- [ ] Vue "Scroll target is not found" warnings on mobile resize across multiple inventory pages

---

## Appendix: Current System Capabilities

### Backend Services
- **Identity:** Authentication, JWT tokens, session management
- **Production:** Batches, recipes, styles, vessels, occupancies, brew sessions, volumes, transfers, additions, measurements, process phases
- **Inventory:** Ingredients, lots, receipts, usage, adjustments, transfers, movements, stock locations, beer lots
- **Procurement:** Suppliers, purchase orders, line items, fees

### Frontend Screens
- Dashboard with operational overview
- Batch management (all, in-progress, detail with 5 tabs)
- Vessel management (active, all, detail)
- Recipe management (list with basic CRUD)
- Inventory hub (ingredients by category, activity, locations, adjustments/transfers, product)
- Procurement (purchase orders, suppliers)
- Settings (brewery name, display units)

### Key Strengths
- Comprehensive batch lifecycle tracking
- Full lot traceability (supplier → ingredient → batch → beer lot)
- Flexible volume model supporting splits and blends
- Dual status tracking (batch phase + occupancy status)
- Polished UI with consistent design patterns
- User display preferences (units)

### Key Gaps (to be addressed in Phases 5-8)
- Limited edit/delete operations
- No stock level visibility
- No PO → receipt workflow
- Recipe lacks ingredient details
- No cost tracking
- No packaging workflow
- No removal tracking

### Recently Addressed (Phase 4)
- ✅ Brew day workflow - complete ingredient picking, session recording, fermenter assignment
- ✅ Inventory integration - atomic deduction with FIFO lot selection and stock validation
- ✅ Recipe scaling - cross-unit volume conversion with scaling controls
- ✅ Mobile brew day wizard - guided 3-step flow optimized for brew floor use
- ✅ Frontend tests - 366 tests covering composables, components, and scaling

### Recently Addressed (Phase 0)
- ✅ Mobile responsiveness - now works well on all devices
- ✅ Frontend tests - 331 tests covering composables and components
- ✅ Code organization - BatchDetails decomposed, types consolidated
