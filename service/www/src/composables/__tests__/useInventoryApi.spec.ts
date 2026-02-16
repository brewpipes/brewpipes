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

  describe('utility functions', () => {
    it('normalizeText trims whitespace and returns null for empty strings', () => {
      const { normalizeText } = useInventoryApi()

      expect(normalizeText('  hello  ')).toBe('hello')
      expect(normalizeText('test')).toBe('test')
      expect(normalizeText('   ')).toBeNull()
      expect(normalizeText('')).toBeNull()
    })

    it('normalizeDateTime converts to ISO string', () => {
      const { normalizeDateTime } = useInventoryApi()

      const result = normalizeDateTime('2024-01-15T10:30:00')
      expect(result).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)

      expect(normalizeDateTime('')).toBeNull()
    })

    it('toNumber parses numeric values correctly', () => {
      const { toNumber } = useInventoryApi()

      expect(toNumber('42')).toBe(42)
      expect(toNumber(42)).toBe(42)
      expect(toNumber('3.14')).toBe(3.14)
      expect(toNumber('0')).toBe(0)
      expect(toNumber(0)).toBe(0)
      expect(toNumber('')).toBeNull()
      expect(toNumber(null)).toBeNull()
      expect(toNumber('not a number')).toBeNull()
      expect(toNumber('NaN')).toBeNull()
      expect(toNumber('Infinity')).toBeNull()
    })

    it('formatDateTime formats dates correctly', () => {
      const { formatDateTime } = useInventoryApi()

      // Test with a valid date
      const result = formatDateTime('2024-01-15T10:30:00Z')
      expect(result).toContain('2024')

      expect(formatDateTime(null)).toBe('Unknown')
      expect(formatDateTime(undefined)).toBe('Unknown')
      expect(formatDateTime('')).toBe('Unknown')
    })

    it('formatAmount formats amount with unit', () => {
      const { formatAmount } = useInventoryApi()

      expect(formatAmount(10, 'kg')).toBe('10 kg')
      expect(formatAmount(5.5, 'lb')).toBe('5.5 lb')
      expect(formatAmount(100, null)).toBe('100')
      expect(formatAmount(100, undefined)).toBe('100')
      expect(formatAmount(null, 'kg')).toBe('n/a')
      expect(formatAmount(undefined, 'kg')).toBe('n/a')
    })

    it('formatAmount handles zero correctly', () => {
      const { formatAmount } = useInventoryApi()

      expect(formatAmount(0, 'kg')).toBe('0 kg')
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

    it('exposes all utility functions', () => {
      const api = useInventoryApi()

      expect(typeof api.normalizeText).toBe('function')
      expect(typeof api.normalizeDateTime).toBe('function')
      expect(typeof api.toNumber).toBe('function')
      expect(typeof api.formatDateTime).toBe('function')
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

  describe('edge cases', () => {
    it('normalizeText handles strings with only whitespace characters', () => {
      const { normalizeText } = useInventoryApi()

      expect(normalizeText('\t\n\r ')).toBeNull()
      expect(normalizeText('  \t  ')).toBeNull()
    })

    it('toNumber handles edge numeric cases', () => {
      const { toNumber } = useInventoryApi()

      expect(toNumber('-5')).toBe(-5)
      expect(toNumber('0.001')).toBe(0.001)
      expect(toNumber('1e10')).toBe(1e10)
    })

    it('formatAmount handles decimal amounts', () => {
      const { formatAmount } = useInventoryApi()

      expect(formatAmount(10.5, 'oz')).toBe('10.5 oz')
      expect(formatAmount(0.25, 'lb')).toBe('0.25 lb')
    })
  })
})
