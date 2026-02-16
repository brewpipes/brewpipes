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
  supplier_uuid: string | null
  purchase_order_line_uuid: string | null
  brewery_lot_code: string | null
  originator_lot_code: string | null
  originator_name: string | null
  originator_type: string | null
  received_at: string
  received_amount: number
  received_unit: string
  current_amount: number
  current_unit: string
  best_by_at: string | null
  expires_at: string | null
  notes: string | null
  created_at: string
  updated_at: string
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
  packaging_run_uuid?: string
  lot_code: string | null
  best_by?: string
  package_format_name?: string
  container?: string
  volume_per_unit?: number
  volume_per_unit_unit?: string
  quantity?: number
  packaged_at: string
  notes: string | null
  created_at: string
  updated_at: string
}

/** Current stock level for a beer lot at a specific location */
export interface BeerLotStockLevel {
  beer_lot_uuid: string
  production_batch_uuid: string
  packaging_run_uuid?: string
  lot_code?: string
  best_by?: string
  package_format_name?: string
  container?: string
  volume_per_unit?: number
  volume_per_unit_unit?: string
  initial_quantity?: number
  packaged_at: string
  stock_location_uuid: string
  stock_location_name: string
  current_volume: number
  current_volume_unit: string
  current_quantity?: number
}

/** Current stock level for an ingredient lot at a specific location */
export interface IngredientLotStockLevel {
  ingredient_lot_uuid: string
  ingredient_uuid: string
  ingredient_name: string
  ingredient_category: string
  brewery_lot_code?: string
  received_at: string
  received_amount: number
  received_unit: string
  stock_location_uuid: string
  stock_location_name: string
  current_amount: number
  current_unit: string
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

// ============================================================================
// Ingredient Request Types
// ============================================================================

/** Request payload for creating a new ingredient */
export interface CreateIngredientRequest {
  name: string
  category: string
  default_unit: string
  description?: string | null
}

// ============================================================================
// Ingredient Lot Detail Request Types
// ============================================================================

/** Request payload for creating malt-specific lot details */
export interface CreateIngredientLotMaltDetailRequest {
  ingredient_lot_uuid: string
  moisture_percent?: number | null
}

/** Request payload for updating malt-specific lot details */
export interface UpdateIngredientLotMaltDetailRequest {
  ingredient_lot_uuid: string
  moisture_percent?: number | null
}

/** Request payload for creating hop-specific lot details */
export interface CreateIngredientLotHopDetailRequest {
  ingredient_lot_uuid: string
  alpha_acid?: number | null
  beta_acid?: number | null
}

/** Request payload for updating hop-specific lot details */
export interface UpdateIngredientLotHopDetailRequest {
  ingredient_lot_uuid: string
  alpha_acid?: number | null
  beta_acid?: number | null
}

/** Request payload for creating yeast-specific lot details */
export interface CreateIngredientLotYeastDetailRequest {
  ingredient_lot_uuid: string
  viability_percent?: number | null
  generation?: number | null
}

/** Request payload for updating yeast-specific lot details */
export interface UpdateIngredientLotYeastDetailRequest {
  ingredient_lot_uuid: string
  viability_percent?: number | null
  generation?: number | null
}

// ============================================================================
// Stock Location Request Types
// ============================================================================

/** Request payload for creating a new stock location */
export interface CreateStockLocationRequest {
  name: string
  location_type?: string | null
  description?: string | null
}

// ============================================================================
// Inventory Usage Request Types
// ============================================================================

/** Request payload for creating an inventory usage record */
export interface CreateInventoryUsageRequest {
  production_ref_uuid?: string | null
  used_at?: string | null
  notes?: string | null
}

// ============================================================================
// Inventory Adjustment Request Types
// ============================================================================

/** Request payload for creating an inventory adjustment */
export interface CreateInventoryAdjustmentRequest {
  ingredient_lot_uuid?: string | null
  beer_lot_uuid?: string | null
  stock_location_uuid: string
  amount: number
  amount_unit: string
  reason: string
  adjusted_at?: string | null
  notes?: string | null
}

// ============================================================================
// Inventory Transfer Request Types
// ============================================================================

/** Request payload for creating an inventory transfer */
export interface CreateInventoryTransferRequest {
  ingredient_lot_uuid?: string | null
  beer_lot_uuid?: string | null
  source_location_uuid: string
  dest_location_uuid: string | null
  amount: number
  amount_unit: string
  transferred_at?: string | null
  notes?: string | null
}

// ============================================================================
// Beer Lot Request Types
// ============================================================================

/** Request payload for creating a beer lot */
export interface CreateBeerLotRequest {
  production_batch_uuid: string
  lot_code?: string | null
  packaged_at?: string | null
  notes?: string | null
  packaging_run_uuid?: string | null
  best_by?: string | null
  package_format_name?: string | null
  container?: string | null
  volume_per_unit?: number | null
  volume_per_unit_unit?: string | null
  quantity?: number | null
  stock_location_uuid?: string | null
}

// ============================================================================
// Receipt Request Types
// ============================================================================

/** Request payload for creating an inventory receipt */
export interface CreateInventoryReceiptRequest {
  supplier_uuid?: string | null
  purchase_order_uuid?: string | null
  reference_code?: string | null
  received_at?: string | null
  notes?: string | null
}

/** Request payload for creating an ingredient lot */
export interface CreateIngredientLotRequest {
  ingredient_uuid: string
  receipt_uuid?: string | null
  purchase_order_line_uuid?: string | null
  supplier_uuid?: string | null
  brewery_lot_code?: string | null
  originator_lot_code?: string | null
  originator_name?: string | null
  originator_type?: string | null
  received_at?: string | null
  received_amount: number
  received_unit: string
  best_by_at?: string | null
  expires_at?: string | null
  notes?: string | null
}

/** Request payload for creating an inventory movement */
export interface CreateInventoryMovementRequest {
  ingredient_lot_uuid?: string | null
  beer_lot_uuid?: string | null
  stock_location_uuid: string
  direction: 'in' | 'out'
  reason: 'receive' | 'use' | 'transfer' | 'adjust' | 'waste'
  amount: number
  amount_unit: string
  occurred_at?: string | null
  receipt_uuid?: string | null
  usage_uuid?: string | null
  adjustment_uuid?: string | null
  transfer_uuid?: string | null
  notes?: string | null
}

// ============================================================================
// Stock Level Types
// ============================================================================

/** Stock quantity at a specific location */
export interface StockLevelLocation {
  location_uuid: string
  location_name: string
  quantity: number
}

/** Aggregated stock level for an ingredient across all locations */
export interface StockLevel {
  ingredient_uuid: string
  ingredient_name: string
  category: string
  default_unit: string
  total_on_hand: number
  locations: StockLevelLocation[]
}

// ============================================================================
// Batch Usage Types
// ============================================================================

/** A single ingredient pick for batch usage deduction */
export interface BatchUsagePick {
  ingredient_lot_uuid: string
  stock_location_uuid: string
  amount: number
  amount_unit: string
}

/** Request payload for creating a batch usage deduction */
export interface CreateBatchUsageRequest {
  production_ref_uuid?: string
  used_at: string
  picks: BatchUsagePick[]
  notes?: string
}

/** Response from a batch usage deduction */
export interface BatchUsageResponse {
  usage_uuid: string
  movements: InventoryMovement[]
}

// ============================================================================
// Receiving Workflow Types
// ============================================================================

/** Details for each line being received in a shipment */
export interface LineReceivingDetails {
  lineUuid: string
  quantity: number
  unit: string
  locationUuid: string
  breweryLotCode: string
  supplierLotCode: string
}

// ============================================================================
// Removal Types
// ============================================================================

/** Category of a removal event */
export type RemovalCategory = 'dump' | 'waste' | 'sample' | 'expired' | 'other'

/** Reason for a removal event */
export type RemovalReason =
  | 'infection'
  | 'off_flavor'
  | 'failed_fermentation'
  | 'equipment_failure'
  | 'quality_reject'
  | 'past_date'
  | 'damaged_package'
  | 'spillage'
  | 'cleaning'
  | 'qc_sample'
  | 'tasting'
  | 'competition'
  | 'other'

/** A removal record (beer dumped, wasted, sampled, expired, etc.) */
export interface Removal {
  uuid: string
  category: RemovalCategory
  reason: RemovalReason
  amount: number
  amount_unit: string
  amount_bbl: number | null
  is_taxable: boolean
  removed_at: string
  batch_uuid?: string
  beer_lot_uuid?: string
  occupancy_uuid?: string
  stock_location_uuid?: string
  movement_uuid?: string
  reference_code?: string
  performed_by?: string
  destination?: string
  notes?: string
  created_at: string
  updated_at: string
}

/** Request payload for creating a removal */
export interface CreateRemovalRequest {
  category: RemovalCategory
  reason: RemovalReason
  amount: number
  amount_unit: string
  removed_at?: string
  batch_uuid?: string
  beer_lot_uuid?: string
  occupancy_uuid?: string
  stock_location_uuid?: string
  reference_code?: string
  performed_by?: string
  destination?: string
  notes?: string
}

/** Request payload for updating a removal */
export interface UpdateRemovalRequest {
  category?: RemovalCategory
  reason?: RemovalReason
  amount?: number
  amount_unit?: string
  removed_at?: string
  destination?: string
  notes?: string
}

/** Aggregated removal summary */
export interface RemovalSummary {
  total_bbl: number
  taxable_bbl: number
  tax_free_bbl: number
  total_count: number
  by_category: RemovalCategorySummary[]
}

/** Per-category breakdown in a removal summary */
export interface RemovalCategorySummary {
  category: RemovalCategory
  total_bbl: number
  count: number
}
