/**
 * Inventory domain types for ingredients, lots, stock locations, and inventory operations.
 */

// ============================================================================
// Ingredient Types
// ============================================================================

/** An ingredient definition (malt, hop, yeast, adjunct, etc.) */
export interface Ingredient {
  uuid: string
  name: string
  category: string
  default_unit: string
  description: string
  created_at: string
  updated_at: string
}

// ============================================================================
// Ingredient Lot Types
// ============================================================================

/** A specific lot of an ingredient received into inventory */
export interface IngredientLot {
  uuid: string
  ingredient_uuid: string
  receipt_uuid: string | null
  received_amount: number
  received_unit: string
  best_by_at: string
  expires_at: string
  supplier_uuid: string
  brewery_lot_code: string | null
  originator_lot_code: string
  originator_name: string
  originator_type: string
  received_at: string
  notes: string
  stock_location_uuid: string
  current_amount: number
  current_unit: string
}

// ============================================================================
// Ingredient Lot Detail Types
// ============================================================================

/** Malt-specific details for an ingredient lot */
export interface IngredientLotMaltDetail {
  uuid: string
  ingredient_lot_uuid: string
  moisture_percent: number | null
}

/** Hop-specific details for an ingredient lot */
export interface IngredientLotHopDetail {
  uuid: string
  ingredient_lot_uuid: string
  alpha_acid: number | null
  beta_acid: number | null
}

/** Yeast-specific details for an ingredient lot */
export interface IngredientLotYeastDetail {
  uuid: string
  ingredient_lot_uuid: string
  viability_percent: number | null
  generation: number | null
}

// ============================================================================
// Beer Lot Types
// ============================================================================

/** A lot of finished beer product */
export interface BeerLot {
  uuid: string
  production_batch_uuid: string
  lot_code: string | null
  packaged_at: string
  notes: string
  created_at: string
  updated_at: string
  stock_location_uuid: string
  volume: number
  volume_unit: string
}

// ============================================================================
// Stock Location Types
// ============================================================================

/** A physical storage location for inventory */
export interface StockLocation {
  uuid: string
  name: string
  location_type: string
  description: string
  created_at: string
  updated_at: string
}

// ============================================================================
// Inventory Receipt Types
// ============================================================================

/** A record of inventory received from a supplier */
export interface InventoryReceipt {
  uuid: string
  supplier_uuid: string | null
  reference_code: string | null
  received_at: string
  notes: string | null
  created_at: string
  updated_at: string
}

// ============================================================================
// Inventory Usage Types
// ============================================================================

/** A record of inventory used in production */
export interface InventoryUsage {
  uuid: string
  production_ref_uuid: string | null
  used_at: string
  notes: string | null
  created_at: string
  updated_at: string
}

// ============================================================================
// Inventory Movement Types
// ============================================================================

/** A record of inventory movement (in or out of a location) */
export interface InventoryMovement {
  uuid: string
  ingredient_lot_uuid: string | null
  beer_lot_uuid: string | null
  stock_location_uuid: string
  direction: string
  reason: string
  amount: number
  amount_unit: string
  occurred_at: string
  receipt_uuid: string | null
  usage_uuid: string | null
  adjustment_uuid: string | null
  transfer_uuid: string | null
  notes: string | null
}

// ============================================================================
// Inventory Adjustment Types
// ============================================================================

/** A manual inventory adjustment (count correction, spoilage, etc.) */
export interface InventoryAdjustment {
  uuid: string
  reason: string
  adjusted_at: string
  notes: string | null
}

// ============================================================================
// Inventory Transfer Types
// ============================================================================

/** A transfer of inventory between stock locations */
export interface InventoryTransfer {
  uuid: string
  source_location_uuid: string
  dest_location_uuid: string
  transferred_at: string
  notes: string | null
}
