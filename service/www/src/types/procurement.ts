/**
 * Procurement domain types for suppliers, purchase orders, and related entities.
 */

// ============================================================================
// Supplier Types
// ============================================================================

/** A supplier of ingredients or materials */
export interface Supplier {
  uuid: string
  name: string
  contact_name: string | null
  email: string | null
  phone: string | null
  address_line1: string | null
  address_line2: string | null
  city: string | null
  region: string | null
  postal_code: string | null
  country: string | null
  created_at: string
  updated_at: string
}

/** Request payload for updating an existing supplier */
export interface UpdateSupplierRequest {
  name?: string
  contact_name?: string | null
  email?: string | null
  phone?: string | null
  address_line1?: string | null
  address_line2?: string | null
  city?: string | null
  region?: string | null
  postal_code?: string | null
  country?: string | null
}

// ============================================================================
// Purchase Order Types
// ============================================================================

/** A purchase order placed with a supplier */
export interface PurchaseOrder {
  uuid: string
  supplier_uuid: string
  order_number: string
  status: string
  ordered_at: string | null
  expected_at: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

/** Request payload for creating a new purchase order */
export interface CreatePurchaseOrderRequest {
  supplier_uuid: string
  order_number: string
  status?: string | null
  ordered_at?: string | null
  expected_at?: string | null
  notes?: string | null
}

/** Request payload for updating an existing purchase order */
export interface UpdatePurchaseOrderRequest {
  order_number?: string
  status?: string
  ordered_at?: string | null
  expected_at?: string | null
  notes?: string | null
}

// ============================================================================
// Purchase Order Line Types
// ============================================================================

/** A line item on a purchase order */
export interface PurchaseOrderLine {
  uuid: string
  purchase_order_uuid: string
  line_number: number
  item_type: string
  item_name: string
  inventory_item_uuid: string | null
  quantity: number
  quantity_unit: string
  unit_cost_cents: number
  currency: string
  created_at: string
  updated_at: string
}

// ============================================================================
// Purchase Order Fee Types
// ============================================================================

/** A fee associated with a purchase order (shipping, tax, etc.) */
export interface PurchaseOrderFee {
  uuid: string
  purchase_order_uuid: string
  fee_type: string
  amount_cents: number
  currency: string
  created_at: string
  updated_at: string
}
