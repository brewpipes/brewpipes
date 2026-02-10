# BrewPipes Shared Domain Context

This file is the canonical domain reference for all BrewPipes agents. Agents should reference this file rather than maintaining their own copy of domain context. For the full product brief, see `PROJECT.md` at the repo root. For the V1 roadmap, see `V1_ROADMAP.md`.

## Product summary

BrewPipes is an open source brewery management system focused on day-to-day production operations. It is the operational source of truth for procurement, inventory, and production, with clear traceability from ingredients and receipts through batches, vessels, and finished lots.

## Target users

Small craft breweries where 1-2 people wear multiple hats — managing both business operations and hands-on brewing. The primary user is often on a tablet or phone on the brewery floor, possibly with wet or gloved hands. Mobile/tablet experience is critical for V1.

## Core modules

| Module | Service path | Owns |
|--------|-------------|------|
| **Identity** | `service/identity/` | Users, authentication, JWT sessions |
| **Production** | `service/production/` | Styles, recipes, batches, brew sessions, volumes, vessels, occupancies, transfers, additions, measurements, process phases |
| **Inventory** | `service/inventory/` | Ingredients, ingredient lots, beer lots, stock locations, movements (receipts, usage, adjustments, transfers) |
| **Procurement** | `service/procurement/` | Suppliers, purchase orders, line items, fees |
| **Web UI** | `service/www/` | Vue 3 + Vuetify 3 frontend |

## Primary user journeys (V1)

1. **Procurement & Receiving** — Order ingredients, receive deliveries, update inventory with full PO traceability.
2. **Brew Day Execution** — Follow recipe, pick ingredients from inventory, record mash/boil/knockout, pitch yeast, move to fermenter.
3. **Fermentation Management** — Monitor active fermentations, record measurements/additions, manage transfers/splits/blends.
4. **Packaging & Finished Goods** — Record packaging runs, create beer lots, calculate packaging loss.
5. **Batch Costing & Review** — Review batch performance, ingredient costs, actual vs. target comparison, full traceability.
6. **Inventory Management & Removals** — Track stock levels, adjustments, transfers, dumps, waste, spillage.

## Cross-service data ownership

Each service owns its tables and is the source of truth for its entities. Cross-service references use UUIDs — not foreign keys — and are resolved at the application layer:

| Service | References from other services |
|---------|-------------------------------|
| Identity | Referenced by all services (user context) |
| Production | References ingredient lots (Inventory) by UUID |
| Inventory | References batches (Production) and suppliers (Procurement) by UUID |
| Procurement | References ingredients (Inventory) by UUID |

## Core entities

- **Procurement:** supplier, purchase_order, purchase_order_line, purchase_order_fee
- **Inventory:** ingredient, ingredient_*_detail, stock_location, inventory_receipt, ingredient_lot, inventory_usage, inventory_adjustment, inventory_transfer, inventory_movement, beer_lot
- **Production:** style, recipe, recipe_ingredient, batch, brew_session, volume, volume_relation, vessel, occupancy, transfer, batch_volume, batch_process_phase, batch_relation, addition, measurement

## Key domain concepts

- **Dual status tracking:** Batch process phases track high-level lifecycle (planning → finished). Occupancy status tracks granular liquid-in-vessel state (fermenting, conditioning, etc.).
- **Hot-side vs cold-side:** Brew session captures hot-side production; additions/measurements attach to the wort volume. Post-yeast-pitch is cold-side; additions/measurements attach to occupancies.
- **Volume model:** Supports splits (one volume → many), blends (many → one), and simple transfers.
- **Traceability chain:** Supplier → PO → Receipt → Ingredient Lot → Batch → Beer Lot.
- **ABV calculation:** `(OG - FG) × 131.25`, with manual override capability.

## V1 scope boundaries

- **In scope:** Procurement, brew day, fermentation, packaging, inventory, batch costing, brewhouse removals.
- **Out of scope:** Taproom, merchandise, sales, multi-user/roles, multi-tenancy.

## Change posture

- Be conservative with changes that alter data integrity, migrations, or cross-service links.
- Iterate faster on DTOs, response shaping, UI affordances, and non-critical endpoints.
- Prefer backward-compatible changes to schemas and APIs.
