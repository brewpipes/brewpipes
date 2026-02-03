# BrewPipes V1 Product Roadmap

**Last Updated:** 2026-02-02  
**Status:** User Journeys Finalized, Feature Backlog In Progress

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
| CRUD-01 | Add batch editing (name, brew date, recipe, notes) | M | 2, 5 | Not Started |
| CRUD-02 | Add batch deletion with confirmation and checks | S | 2 | Not Started |
| CRUD-03 | Add vessel editing | S | 2, 3 | Not Started |
| CRUD-04 | Add vessel retirement workflow | S | - | Not Started |
| CRUD-05 | Add recipe deletion with batch reference check | S | 2 | Not Started |
| CRUD-06 | Add supplier editing | S | 1 | Not Started |
| CRUD-07 | Add purchase order editing | M | 1 | Not Started |

---

### Phase 2: Recipe Enhancement

Enable full recipe management with ingredient bills.

| ID | Feature | Size | Journey | Status |
|----|---------|------|---------|--------|
| RCP-01 | Design recipe ingredient bill data model | M | 2 | Not Started |
| RCP-02 | Add recipe target specifications (OG, FG, ABV, IBU, color) | M | 2, 5 | Not Started |
| RCP-03 | Backend: Recipe ingredient bill CRUD endpoints | M | 2 | Not Started |
| RCP-04 | Frontend: Recipe ingredient bill management UI | L | 2 | Not Started |
| RCP-05 | Frontend: Recipe detail view with full specs | M | 2, 5 | Not Started |
| RCP-06 | Recipe scaling calculator | M | 2 | Not Started |

---

### Phase 3: Procurement & Receiving Workflow

Complete the procurement → inventory flow.

| ID | Feature | Size | Journey | Status |
|----|---------|------|---------|--------|
| PROC-01 | Backend: Link receipts to PO line items | M | 1 | Not Started |
| PROC-02 | Backend: Auto-update PO status on receipt | S | 1 | Not Started |
| PROC-03 | Frontend: PO detail view with line items inline | M | 1 | Not Started |
| PROC-04 | Frontend: Receiving workflow (receive against PO) | L | 1 | Not Started |
| PROC-05 | Frontend: Current stock levels display | M | 1, 6 | Not Started |
| PROC-06 | Frontend: Low stock alerts | S | 6 | Not Started |

---

### Phase 4: Brew Day & Occupancy

Enable complete brew day recording with inventory integration.

| ID | Feature | Size | Journey | Status |
|----|---------|------|---------|--------|
| BREW-01 | Frontend: Ingredient pick list from recipe | M | 2 | Not Started |
| BREW-02 | Frontend: Lot selection for ingredient picks | M | 2 | Not Started |
| BREW-03 | Backend: Inventory deduction on ingredient use | M | 2 | Not Started |
| BREW-04 | Frontend: Create occupancy from batch detail | M | 2, 3 | Not Started |
| BREW-05 | Frontend: End occupancy workflow | S | 3, 4 | Not Started |
| BREW-06 | Mobile-optimized brew day recording UI | L | 2 | Not Started |

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
| **M2: Core Complete** | Phases 1-2 complete (CRUD, recipes) | Not Started |
| **M3: Procurement Flow** | Phase 3 complete (PO → inventory) | Not Started |
| **M4: Brew Day Flow** | Phase 4 complete (brew day recording) | Not Started |
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

**TD-01 Details:** Refactored 3,227-line component into 15 smaller components in `service/www/src/components/batch/`. Main component reduced to 1,677 lines (~48% reduction). Created 6 tab components, 7 dialog components, 1 reusable card component, shared types file, and barrel export.

**TD-02 Details:** Created 4 test files with 97 tests covering `useApiClient`, `useProductionApi`, `useInventoryApi`, `useProcurementApi`. Set up Vitest with happy-dom environment.

**TD-03 Details:** Created 4 test files with 198 tests covering `useFormatters`, `useUnitConversion`, `useUnitPreferences`, `useUserSettings`. Total test suite now has 295 passing tests.

**TD-04 Details:** Created `service/www/src/types/` directory with 6 files organizing domain types (common, units, production, settings, auth). Consolidated duplicate types including `VolumeUnit` (4→9 values). Maintained backward compatibility via re-exports from composables.

**TD-05 Details:** Modified 12 files for mobile responsiveness. Implemented master-detail mobile pattern (list OR detail on mobile), responsive dialogs, icon-only buttons on xs, 44px touch targets, table horizontal scrolling. App now works well on phones, tablets, and desktops.

**TD-06 Details:** Created 3 component test files with 36 tests covering `BatchList`, `VesselList`, and `AppFooter`. Tests cover rendering, selection, events, empty states, and sorting. Total test suite now has 331 passing tests. Updated Vitest config for Vuetify component testing.

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

### Key Gaps (to be addressed in Phases 1-8)
- Limited edit/delete operations
- No occupancy creation from UI
- No stock level visibility
- No PO → receipt workflow
- Recipe lacks ingredient details
- No cost tracking
- No packaging workflow
- No removal tracking

### Recently Addressed (Phase 0)
- ✅ Mobile responsiveness - now works well on all devices
- ✅ Frontend tests - 331 tests covering composables and components
- ✅ Code organization - BatchDetails decomposed, types consolidated
