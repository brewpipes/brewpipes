import type {
  BatchUsageResponse,
  BeerLot,
  BeerLotStockLevel,
  CreateBatchUsageRequest,
  CreateIngredientLotRequest,
  CreateInventoryMovementRequest,
  CreateInventoryReceiptRequest,
  Ingredient,
  IngredientLot,
  IngredientLotHopDetail,
  IngredientLotMaltDetail,
  IngredientLotYeastDetail,
  InventoryAdjustment,
  InventoryMovement,
  InventoryReceipt,
  InventoryTransfer,
  InventoryUsage,
  StockLevel,
  StockLocation,
} from '@/types'
import { useApiClient } from '@/composables/useApiClient'
import { formatDateTime } from '@/composables/useFormatters'
import { normalizeDateTime, normalizeText, toNumber } from '@/utils/normalize'

const inventoryApiBase = import.meta.env.VITE_INVENTORY_API_URL ?? '/api'

export function useInventoryApi () {
  const { request } = useApiClient(inventoryApiBase)

  const formatAmount = (amount: number | null | undefined, unit: string | null | undefined) => {
    if (amount === null || amount === undefined) {
      return 'n/a'
    }
    return `${amount} ${unit ?? ''}`.trim()
  }

  // Ingredients API
  const getIngredients = () => request<Ingredient[]>('/ingredients')
  const createIngredient = (data: Record<string, unknown>) =>
    request<Ingredient>('/ingredients', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Ingredient Lots API
  const getIngredientLots = (filters?: { purchase_order_line_uuid?: string, ingredient_uuid?: string }) => {
    const query = new URLSearchParams()
    if (filters?.purchase_order_line_uuid) {
      query.set('purchase_order_line_uuid', filters.purchase_order_line_uuid)
    }
    if (filters?.ingredient_uuid) {
      query.set('ingredient_uuid', filters.ingredient_uuid)
    }
    const path = query.toString() ? `/ingredient-lots?${query.toString()}` : '/ingredient-lots'
    return request<IngredientLot[]>(path)
  }
  const getIngredientLot = (uuid: string) => request<IngredientLot>(`/ingredient-lots/${uuid}`)
  const createIngredientLot = (data: CreateIngredientLotRequest | Record<string, unknown>) =>
    request<IngredientLot>('/ingredient-lots', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Ingredient Lot Details API
  const getIngredientLotMaltDetail = (ingredientLotUuid: string) =>
    request<IngredientLotMaltDetail>(`/ingredient-lot-malt-details?ingredient_lot_uuid=${ingredientLotUuid}`)
  const createIngredientLotMaltDetail = (data: Record<string, unknown>) =>
    request<IngredientLotMaltDetail>('/ingredient-lot-malt-details', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updateIngredientLotMaltDetail = (uuid: string, data: Record<string, unknown>) =>
    request<IngredientLotMaltDetail>(`/ingredient-lot-malt-details/${uuid}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  const getIngredientLotHopDetail = (ingredientLotUuid: string) =>
    request<IngredientLotHopDetail>(`/ingredient-lot-hop-details?ingredient_lot_uuid=${ingredientLotUuid}`)
  const createIngredientLotHopDetail = (data: Record<string, unknown>) =>
    request<IngredientLotHopDetail>('/ingredient-lot-hop-details', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updateIngredientLotHopDetail = (uuid: string, data: Record<string, unknown>) =>
    request<IngredientLotHopDetail>(`/ingredient-lot-hop-details/${uuid}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  const getIngredientLotYeastDetail = (ingredientLotUuid: string) =>
    request<IngredientLotYeastDetail>(`/ingredient-lot-yeast-details?ingredient_lot_uuid=${ingredientLotUuid}`)
  const createIngredientLotYeastDetail = (data: Record<string, unknown>) =>
    request<IngredientLotYeastDetail>('/ingredient-lot-yeast-details', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updateIngredientLotYeastDetail = (uuid: string, data: Record<string, unknown>) =>
    request<IngredientLotYeastDetail>(`/ingredient-lot-yeast-details/${uuid}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })

  // Stock Locations API
  const getStockLocations = () => request<StockLocation[]>('/stock-locations')
  const createStockLocation = (data: Record<string, unknown>) =>
    request<StockLocation>('/stock-locations', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Inventory Receipts API
  const getInventoryReceipts = () => request<InventoryReceipt[]>('/inventory-receipts')
  const createInventoryReceipt = (data: CreateInventoryReceiptRequest | Record<string, unknown>) =>
    request<InventoryReceipt>('/inventory-receipts', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Inventory Usages API
  const getInventoryUsages = () => request<InventoryUsage[]>('/inventory-usage')
  const createInventoryUsage = (data: Record<string, unknown>) =>
    request<InventoryUsage>('/inventory-usage', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Batch Usage API (atomic batch ingredient deduction)
  const createBatchUsage = (data: CreateBatchUsageRequest) =>
    request<BatchUsageResponse>('/inventory-usage/batch', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Inventory Adjustments API
  const getInventoryAdjustments = () => request<InventoryAdjustment[]>('/inventory-adjustments')
  const createInventoryAdjustment = (data: Record<string, unknown>) =>
    request<InventoryAdjustment>('/inventory-adjustments', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Inventory Transfers API
  const getInventoryTransfers = () => request<InventoryTransfer[]>('/inventory-transfers')
  const createInventoryTransfer = (data: Record<string, unknown>) =>
    request<InventoryTransfer>('/inventory-transfers', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Inventory Movements API
  const getInventoryMovements = (filters?: { ingredient_lot_uuid?: string, beer_lot_uuid?: string }) => {
    const query = new URLSearchParams()
    if (filters?.ingredient_lot_uuid) {
      query.set('ingredient_lot_uuid', filters.ingredient_lot_uuid)
    }
    if (filters?.beer_lot_uuid) {
      query.set('beer_lot_uuid', filters.beer_lot_uuid)
    }
    const path = query.toString() ? `/inventory-movements?${query.toString()}` : '/inventory-movements'
    return request<InventoryMovement[]>(path)
  }
  const createInventoryMovement = (data: CreateInventoryMovementRequest) =>
    request<InventoryMovement>('/inventory-movements', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Beer Lots API
  const getBeerLots = () => request<BeerLot[]>('/beer-lots')
  const createBeerLot = (data: Record<string, unknown>) =>
    request<BeerLot>('/beer-lots', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Stock Levels API
  const getStockLevels = () => request<StockLevel[]>('/stock-levels')

  // Beer Lot Stock Levels API
  const getBeerLotStockLevels = () => request<BeerLotStockLevel[]>('/beer-lot-stock-levels')

  return {
    apiBase: inventoryApiBase,
    request,
    normalizeText,
    normalizeDateTime,
    toNumber,
    formatDateTime,
    formatAmount,
    // Ingredients
    getIngredients,
    createIngredient,
    // Ingredient Lots
    getIngredientLots,
    getIngredientLot,
    createIngredientLot,
    // Ingredient Lot Details
    getIngredientLotMaltDetail,
    createIngredientLotMaltDetail,
    updateIngredientLotMaltDetail,
    getIngredientLotHopDetail,
    createIngredientLotHopDetail,
    updateIngredientLotHopDetail,
    getIngredientLotYeastDetail,
    createIngredientLotYeastDetail,
    updateIngredientLotYeastDetail,
    // Stock Locations
    getStockLocations,
    createStockLocation,
    // Inventory Receipts
    getInventoryReceipts,
    createInventoryReceipt,
    // Inventory Usages
    getInventoryUsages,
    createInventoryUsage,
    // Batch Usage
    createBatchUsage,
    // Inventory Adjustments
    getInventoryAdjustments,
    createInventoryAdjustment,
    // Inventory Transfers
    getInventoryTransfers,
    createInventoryTransfer,
    // Inventory Movements
    getInventoryMovements,
    createInventoryMovement,
    // Beer Lots
    getBeerLots,
    createBeerLot,
    // Stock Levels
    getStockLevels,
    // Beer Lot Stock Levels
    getBeerLotStockLevels,
  }
}
