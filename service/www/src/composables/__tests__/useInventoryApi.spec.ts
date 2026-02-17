import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useInventoryApi } from '@/composables/useInventoryApi'

// Mock useApiClient
const mockRequest = vi.fn()

vi.mock('@/composables/useApiClient', () => ({
  useApiClient: () => ({
    request: mockRequest,
  }),
}))

describe('useInventoryApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('formatAmount', () => {
    it('formats amount with unit', () => {
      const { formatAmount } = useInventoryApi()

      expect(formatAmount(10, 'kg')).toBe('10 kg')
      expect(formatAmount(5.5, 'lb')).toBe('5.5 lb')
      expect(formatAmount(100, null)).toBe('100')
      expect(formatAmount(100, undefined)).toBe('100')
      expect(formatAmount(null, 'kg')).toBe('n/a')
      expect(formatAmount(undefined, 'kg')).toBe('n/a')
    })

    it('handles zero correctly', () => {
      const { formatAmount } = useInventoryApi()

      expect(formatAmount(0, 'kg')).toBe('0 kg')
    })

    it('handles decimal amounts', () => {
      const { formatAmount } = useInventoryApi()

      expect(formatAmount(10.5, 'oz')).toBe('10.5 oz')
      expect(formatAmount(0.25, 'lb')).toBe('0.25 lb')
    })
  })

  describe('exposed properties', () => {
    it('exposes apiBase', () => {
      const api = useInventoryApi()
      expect(api.apiBase).toBeDefined()
    })

    it('exposes request function', () => {
      const api = useInventoryApi()
      expect(api.request).toBeDefined()
      expect(typeof api.request).toBe('function')
    })

    it('exposes formatAmount function', () => {
      const api = useInventoryApi()

      expect(typeof api.formatAmount).toBe('function')
    })
  })

  describe('request function usage', () => {
    it('request can be used for GET requests', async () => {
      mockRequest.mockResolvedValue([{ uuid: 'ing-uuid-1', name: 'Pale Malt' }])

      const { request } = useInventoryApi()
      const result = await request('/ingredients')

      expect(mockRequest).toHaveBeenCalledWith('/ingredients')
      expect(result).toEqual([{ uuid: 'ing-uuid-1', name: 'Pale Malt' }])
    })

    it('request can be used for POST requests', async () => {
      const newItem = { name: 'Cascade Hops', quantity: 50 }
      mockRequest.mockResolvedValue({ uuid: 'ing-uuid-2', ...newItem })

      const { request } = useInventoryApi()
      const result = await request('/ingredients', {
        method: 'POST',
        body: JSON.stringify(newItem),
      })

      expect(mockRequest).toHaveBeenCalledWith('/ingredients', {
        method: 'POST',
        body: JSON.stringify(newItem),
      })
      expect(result).toEqual({ uuid: 'ing-uuid-2', ...newItem })
    })

    it('request can be used for PUT requests', async () => {
      const updateData = { name: 'Updated Malt', quantity: 100 }
      mockRequest.mockResolvedValue({ uuid: 'ing-uuid-1', ...updateData })

      const { request } = useInventoryApi()
      const result = await request('/ingredients/ing-uuid-1', {
        method: 'PUT',
        body: JSON.stringify(updateData),
      })

      expect(mockRequest).toHaveBeenCalledWith('/ingredients/ing-uuid-1', {
        method: 'PUT',
        body: JSON.stringify(updateData),
      })
      expect(result).toEqual({ uuid: 'ing-uuid-1', ...updateData })
    })

    it('request can be used for DELETE requests', async () => {
      mockRequest.mockResolvedValue(null)

      const { request } = useInventoryApi()
      const result = await request('/ingredients/ing-uuid-1', {
        method: 'DELETE',
      })

      expect(mockRequest).toHaveBeenCalledWith('/ingredients/ing-uuid-1', {
        method: 'DELETE',
      })
      expect(result).toBeNull()
    })
  })

  describe('error propagation', () => {
    it('propagates errors from request', async () => {
      const error = new Error('Network error')
      mockRequest.mockRejectedValue(error)

      const { request } = useInventoryApi()

      await expect(request('/ingredients')).rejects.toThrow('Network error')
    })

    it('propagates validation errors', async () => {
      const error = new Error('Invalid quantity')
      mockRequest.mockRejectedValue(error)

      const { request } = useInventoryApi()

      await expect(
        request('/ingredients', {
          method: 'POST',
          body: JSON.stringify({ quantity: -1 }),
        }),
      ).rejects.toThrow('Invalid quantity')
    })
  })

  describe('stock levels API', () => {
    it('getStockLevels fetches stock levels', async () => {
      const mockStockLevels = [
        {
          ingredient_uuid: 'ing-uuid-1',
          ingredient_name: 'Pale Malt 2-Row',
          category: 'fermentable',
          default_unit: 'kg',
          total_on_hand: 500,
          locations: [
            { location_uuid: 'loc-1', location_name: 'Grain Room', quantity: 500 },
          ],
        },
        {
          ingredient_uuid: 'ing-uuid-2',
          ingredient_name: 'Cascade Hops',
          category: 'hop',
          default_unit: 'kg',
          total_on_hand: 25,
          locations: [
            { location_uuid: 'loc-2', location_name: 'Cold Storage', quantity: 25 },
          ],
        },
      ]
      mockRequest.mockResolvedValue(mockStockLevels)

      const { getStockLevels } = useInventoryApi()
      const result = await getStockLevels()

      expect(mockRequest).toHaveBeenCalledWith('/stock-levels')
      expect(result).toEqual(mockStockLevels)
    })

    it('getStockLevels returns empty array when no stock', async () => {
      mockRequest.mockResolvedValue([])

      const { getStockLevels } = useInventoryApi()
      const result = await getStockLevels()

      expect(mockRequest).toHaveBeenCalledWith('/stock-levels')
      expect(result).toEqual([])
    })
  })

  describe('beer lot stock levels API', () => {
    it('getBeerLotStockLevels fetches beer lot stock levels', async () => {
      const mockLevels = [
        {
          beer_lot_uuid: 'bl-1',
          production_batch_uuid: 'batch-1',
          lot_code: 'LOT-001',
          package_format_name: '1/2 BBL Keg',
          container: 'keg',
          stock_location_uuid: 'loc-1',
          stock_location_name: 'Cold Room',
          current_volume: 58674,
          current_volume_unit: 'ml',
          current_quantity: 4,
        },
      ]
      mockRequest.mockResolvedValue(mockLevels)

      const { getBeerLotStockLevels } = useInventoryApi()
      const result = await getBeerLotStockLevels()

      expect(mockRequest).toHaveBeenCalledWith('/beer-lot-stock-levels')
      expect(result).toEqual(mockLevels)
    })

    it('getBeerLotStockLevels returns empty array when no stock', async () => {
      mockRequest.mockResolvedValue([])

      const { getBeerLotStockLevels } = useInventoryApi()
      const result = await getBeerLotStockLevels()

      expect(mockRequest).toHaveBeenCalledWith('/beer-lot-stock-levels')
      expect(result).toEqual([])
    })
  })

  describe('ingredient lots API', () => {
    it('getIngredientLots fetches all lots without filters', async () => {
      const mockLots = [{ uuid: 'lot-1', ingredient_uuid: 'ing-1' }]
      mockRequest.mockResolvedValue(mockLots)

      const { getIngredientLots } = useInventoryApi()
      const result = await getIngredientLots()

      expect(mockRequest).toHaveBeenCalledWith('/ingredient-lots')
      expect(result).toEqual(mockLots)
    })

    it('getIngredientLots filters by purchase_order_line_uuid', async () => {
      mockRequest.mockResolvedValue([])

      const { getIngredientLots } = useInventoryApi()
      await getIngredientLots({ purchase_order_line_uuid: 'po-line-1' })

      expect(mockRequest).toHaveBeenCalledWith('/ingredient-lots?purchase_order_line_uuid=po-line-1')
    })

    it('getIngredientLots filters by ingredient_uuid', async () => {
      mockRequest.mockResolvedValue([])

      const { getIngredientLots } = useInventoryApi()
      await getIngredientLots({ ingredient_uuid: 'ing-uuid-1' })

      expect(mockRequest).toHaveBeenCalledWith('/ingredient-lots?ingredient_uuid=ing-uuid-1')
    })

    it('getIngredientLots supports both filters simultaneously', async () => {
      mockRequest.mockResolvedValue([])

      const { getIngredientLots } = useInventoryApi()
      await getIngredientLots({ purchase_order_line_uuid: 'po-1', ingredient_uuid: 'ing-1' })

      expect(mockRequest).toHaveBeenCalledWith(
        '/ingredient-lots?purchase_order_line_uuid=po-1&ingredient_uuid=ing-1',
      )
    })
  })

  describe('batch usage API', () => {
    it('createBatchUsage sends POST request with correct body', async () => {
      const mockResponse = {
        usage_uuid: 'usage-1',
        movements: [{ uuid: 'mov-1', amount: 50, amount_unit: 'kg' }],
      }
      mockRequest.mockResolvedValue(mockResponse)

      const { createBatchUsage } = useInventoryApi()
      const requestData = {
        production_ref_uuid: 'batch-uuid-1',
        used_at: '2026-02-14T08:00:00Z',
        picks: [
          {
            ingredient_lot_uuid: 'lot-1',
            stock_location_uuid: 'loc-1',
            amount: 50,
            amount_unit: 'kg',
          },
        ],
        notes: 'Brew day pick',
      }

      const result = await createBatchUsage(requestData)

      expect(mockRequest).toHaveBeenCalledWith('/inventory-usage/batch', {
        method: 'POST',
        body: JSON.stringify(requestData),
      })
      expect(result).toEqual(mockResponse)
    })

    it('createBatchUsage supports multiple picks', async () => {
      const mockResponse = { usage_uuid: 'usage-2', movements: [] }
      mockRequest.mockResolvedValue(mockResponse)

      const { createBatchUsage } = useInventoryApi()
      const requestData = {
        production_ref_uuid: 'batch-uuid-1',
        used_at: '2026-02-14T08:00:00Z',
        picks: [
          {
            ingredient_lot_uuid: 'lot-1',
            stock_location_uuid: 'loc-1',
            amount: 30,
            amount_unit: 'kg',
          },
          {
            ingredient_lot_uuid: 'lot-2',
            stock_location_uuid: 'loc-1',
            amount: 20,
            amount_unit: 'kg',
          },
        ],
      }

      await createBatchUsage(requestData)

      expect(mockRequest).toHaveBeenCalledWith('/inventory-usage/batch', {
        method: 'POST',
        body: JSON.stringify(requestData),
      })
    })

    it('createBatchUsage propagates errors', async () => {
      const error = new Error('Insufficient stock for lot LOT-001 at Grain Room: requested 100 kg, available 50 kg')
      mockRequest.mockRejectedValue(error)

      const { createBatchUsage } = useInventoryApi()

      await expect(
        createBatchUsage({
          used_at: '2026-02-14T08:00:00Z',
          picks: [
            {
              ingredient_lot_uuid: 'lot-1',
              stock_location_uuid: 'loc-1',
              amount: 100,
              amount_unit: 'kg',
            },
          ],
        }),
      ).rejects.toThrow('Insufficient stock')
    })
  })

  describe('stock locations API', () => {
    it('getStockLocations fetches all locations', async () => {
      const mockLocations = [{ uuid: 'loc-1', name: 'Main Warehouse' }]
      mockRequest.mockResolvedValue(mockLocations)

      const { getStockLocations } = useInventoryApi()
      const result = await getStockLocations()

      expect(mockRequest).toHaveBeenCalledWith('/stock-locations')
      expect(result).toEqual(mockLocations)
    })

    it('createStockLocation sends POST request with correct body', async () => {
      const mockLocation = { uuid: 'loc-new', name: 'Cold Storage', location_type: 'refrigerated', description: null }
      mockRequest.mockResolvedValue(mockLocation)

      const { createStockLocation } = useInventoryApi()
      const requestData = { name: 'Cold Storage', location_type: 'refrigerated', description: null }

      const result = await createStockLocation(requestData)

      expect(mockRequest).toHaveBeenCalledWith('/stock-locations', {
        method: 'POST',
        body: JSON.stringify(requestData),
      })
      expect(result).toEqual(mockLocation)
    })

    it('updateStockLocation sends PATCH request with correct body', async () => {
      const mockLocation = { uuid: 'loc-1', name: 'Updated Warehouse', location_type: 'warehouse', description: 'Updated' }
      mockRequest.mockResolvedValue(mockLocation)

      const { updateStockLocation } = useInventoryApi()
      const updateData = { name: 'Updated Warehouse', description: 'Updated' }

      const result = await updateStockLocation('loc-1', updateData)

      expect(mockRequest).toHaveBeenCalledWith('/stock-locations/loc-1', {
        method: 'PATCH',
        body: JSON.stringify(updateData),
      })
      expect(result).toEqual(mockLocation)
    })

    it('updateStockLocation supports partial updates', async () => {
      mockRequest.mockResolvedValue({ uuid: 'loc-1', name: 'Renamed' })

      const { updateStockLocation } = useInventoryApi()
      await updateStockLocation('loc-1', { name: 'Renamed' })

      expect(mockRequest).toHaveBeenCalledWith('/stock-locations/loc-1', {
        method: 'PATCH',
        body: JSON.stringify({ name: 'Renamed' }),
      })
    })

    it('deleteStockLocation sends DELETE request', async () => {
      mockRequest.mockResolvedValue(null)

      const { deleteStockLocation } = useInventoryApi()
      const result = await deleteStockLocation('loc-1')

      expect(mockRequest).toHaveBeenCalledWith('/stock-locations/loc-1', {
        method: 'DELETE',
      })
      expect(result).toBeNull()
    })

    it('deleteStockLocation propagates conflict errors', async () => {
      const error = new Error('stock location has inventory and cannot be deleted')
      mockRequest.mockRejectedValue(error)

      const { deleteStockLocation } = useInventoryApi()

      await expect(deleteStockLocation('loc-1')).rejects.toThrow('stock location has inventory and cannot be deleted')
    })
  })

  describe('removals API', () => {
    it('listRemovals fetches all removals without filters', async () => {
      const mockRemovals = [{ uuid: 'rem-1', category: 'dump', reason: 'infection' }]
      mockRequest.mockResolvedValue(mockRemovals)

      const { listRemovals } = useInventoryApi()
      const result = await listRemovals()

      expect(mockRequest).toHaveBeenCalledWith('/removals')
      expect(result).toEqual(mockRemovals)
    })

    it('listRemovals applies batch_uuid filter', async () => {
      mockRequest.mockResolvedValue([])

      const { listRemovals } = useInventoryApi()
      await listRemovals({ batch_uuid: 'batch-1' })

      expect(mockRequest).toHaveBeenCalledWith('/removals?batch_uuid=batch-1')
    })

    it('listRemovals applies category filter', async () => {
      mockRequest.mockResolvedValue([])

      const { listRemovals } = useInventoryApi()
      await listRemovals({ category: 'dump' })

      expect(mockRequest).toHaveBeenCalledWith('/removals?category=dump')
    })

    it('listRemovals applies date range filters', async () => {
      mockRequest.mockResolvedValue([])

      const { listRemovals } = useInventoryApi()
      await listRemovals({ from: '2026-01-01T00:00:00Z', to: '2026-02-01T00:00:00Z' })

      expect(mockRequest).toHaveBeenCalledWith('/removals?from=2026-01-01T00%3A00%3A00Z&to=2026-02-01T00%3A00%3A00Z')
    })

    it('listRemovals supports multiple filters simultaneously', async () => {
      mockRequest.mockResolvedValue([])

      const { listRemovals } = useInventoryApi()
      await listRemovals({ batch_uuid: 'batch-1', category: 'sample' })

      expect(mockRequest).toHaveBeenCalledWith('/removals?batch_uuid=batch-1&category=sample')
    })

    it('getRemoval fetches a single removal', async () => {
      const mockRemoval = { uuid: 'rem-1', category: 'dump', reason: 'infection', amount: 800, amount_unit: 'l' }
      mockRequest.mockResolvedValue(mockRemoval)

      const { getRemoval } = useInventoryApi()
      const result = await getRemoval('rem-1')

      expect(mockRequest).toHaveBeenCalledWith('/removals/rem-1')
      expect(result).toEqual(mockRemoval)
    })

    it('createRemoval sends POST request with correct body', async () => {
      const mockRemoval = { uuid: 'rem-new', category: 'dump', reason: 'infection', amount: 800, amount_unit: 'l' }
      mockRequest.mockResolvedValue(mockRemoval)

      const { createRemoval } = useInventoryApi()
      const requestData = {
        category: 'dump' as const,
        reason: 'infection' as const,
        amount: 800,
        amount_unit: 'l',
        batch_uuid: 'batch-1',
        notes: 'Infected batch',
      }

      const result = await createRemoval(requestData)

      expect(mockRequest).toHaveBeenCalledWith('/removals', {
        method: 'POST',
        body: JSON.stringify(requestData),
      })
      expect(result).toEqual(mockRemoval)
    })

    it('updateRemoval sends PATCH request with correct body', async () => {
      const mockRemoval = { uuid: 'rem-1', category: 'dump', reason: 'off_flavor', amount: 800, amount_unit: 'l' }
      mockRequest.mockResolvedValue(mockRemoval)

      const { updateRemoval } = useInventoryApi()
      const updateData = { reason: 'off_flavor' as const, notes: 'Updated reason' }

      const result = await updateRemoval('rem-1', updateData)

      expect(mockRequest).toHaveBeenCalledWith('/removals/rem-1', {
        method: 'PATCH',
        body: JSON.stringify(updateData),
      })
      expect(result).toEqual(mockRemoval)
    })

    it('deleteRemoval sends DELETE request', async () => {
      mockRequest.mockResolvedValue(null)

      const { deleteRemoval } = useInventoryApi()
      const result = await deleteRemoval('rem-1')

      expect(mockRequest).toHaveBeenCalledWith('/removals/rem-1', {
        method: 'DELETE',
      })
      expect(result).toBeNull()
    })

    it('getRemovalSummary fetches summary without filters', async () => {
      const mockSummary = {
        total_bbl: 7.16,
        taxable_bbl: 0,
        tax_free_bbl: 7.16,
        total_count: 4,
        by_category: [
          { category: 'dump', total_bbl: 6.82, count: 1 },
        ],
      }
      mockRequest.mockResolvedValue(mockSummary)

      const { getRemovalSummary } = useInventoryApi()
      const result = await getRemovalSummary()

      expect(mockRequest).toHaveBeenCalledWith('/removal-summary')
      expect(result).toEqual(mockSummary)
    })

    it('getRemovalSummary applies date range filters', async () => {
      mockRequest.mockResolvedValue({ total_bbl: 0, taxable_bbl: 0, tax_free_bbl: 0, total_count: 0, by_category: [] })

      const { getRemovalSummary } = useInventoryApi()
      await getRemovalSummary({ from: '2026-01-01T00:00:00Z' })

      expect(mockRequest).toHaveBeenCalledWith('/removal-summary?from=2026-01-01T00%3A00%3A00Z')
    })

    it('createRemoval propagates errors', async () => {
      const error = new Error('Validation failed')
      mockRequest.mockRejectedValue(error)

      const { createRemoval } = useInventoryApi()

      await expect(
        createRemoval({
          category: 'dump',
          reason: 'infection',
          amount: 0,
          amount_unit: 'l',
        }),
      ).rejects.toThrow('Validation failed')
    })
  })

})
