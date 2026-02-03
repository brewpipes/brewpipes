import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useProcurementApi } from '@/composables/useProcurementApi'

// Mock useApiClient
const mockRequest = vi.fn()

vi.mock('@/composables/useApiClient', () => ({
  useApiClient: () => ({
    request: mockRequest,
  }),
}))

describe('useProcurementApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('utility functions', () => {
    it('normalizeText trims whitespace and returns null for empty strings', () => {
      const { normalizeText } = useProcurementApi()

      expect(normalizeText('  hello  ')).toBe('hello')
      expect(normalizeText('test')).toBe('test')
      expect(normalizeText('   ')).toBeNull()
      expect(normalizeText('')).toBeNull()
    })

    it('normalizeDateTime converts to ISO string', () => {
      const { normalizeDateTime } = useProcurementApi()

      const result = normalizeDateTime('2024-01-15T10:30:00')
      expect(result).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)

      expect(normalizeDateTime('')).toBeNull()
    })

    it('toNumber parses numeric values correctly', () => {
      const { toNumber } = useProcurementApi()

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
      const { formatDateTime } = useProcurementApi()

      // Test with a valid date
      const result = formatDateTime('2024-01-15T10:30:00Z')
      expect(result).toContain('2024')

      expect(formatDateTime(null)).toBe('n/a')
      expect(formatDateTime(undefined)).toBe('n/a')
      expect(formatDateTime('')).toBe('n/a')
    })

    it('formatCurrency formats cents to currency string', () => {
      const { formatCurrency } = useProcurementApi()

      expect(formatCurrency(1000, 'USD')).toBe('10.00 USD')
      expect(formatCurrency(1550, 'EUR')).toBe('15.50 EUR')
      expect(formatCurrency(99, 'GBP')).toBe('0.99 GBP')
      expect(formatCurrency(0, 'USD')).toBe('0.00 USD')
    })

    it('formatCurrency handles null/undefined currency', () => {
      const { formatCurrency } = useProcurementApi()

      expect(formatCurrency(1000, null)).toBe('10.00')
      expect(formatCurrency(1000, undefined)).toBe('10.00')
    })

    it('formatCurrency returns n/a for null/undefined cents', () => {
      const { formatCurrency } = useProcurementApi()

      expect(formatCurrency(null, 'USD')).toBe('n/a')
      expect(formatCurrency(undefined, 'USD')).toBe('n/a')
    })
  })

  describe('exposed properties', () => {
    it('exposes apiBase', () => {
      const api = useProcurementApi()
      expect(api.apiBase).toBeDefined()
    })

    it('exposes request function', () => {
      const api = useProcurementApi()
      expect(api.request).toBeDefined()
      expect(typeof api.request).toBe('function')
    })

    it('exposes all utility functions', () => {
      const api = useProcurementApi()

      expect(typeof api.normalizeText).toBe('function')
      expect(typeof api.normalizeDateTime).toBe('function')
      expect(typeof api.toNumber).toBe('function')
      expect(typeof api.formatDateTime).toBe('function')
      expect(typeof api.formatCurrency).toBe('function')
    })
  })

  describe('request function usage', () => {
    it('request can be used for GET requests', async () => {
      mockRequest.mockResolvedValue([{ id: 1, supplier: 'Malt Co' }])

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders')

      expect(mockRequest).toHaveBeenCalledWith('/purchase-orders')
      expect(result).toEqual([{ id: 1, supplier: 'Malt Co' }])
    })

    it('request can be used for POST requests', async () => {
      const newOrder = { supplier_id: 1, total_cents: 50_000 }
      mockRequest.mockResolvedValue({ id: 1, ...newOrder })

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders', {
        method: 'POST',
        body: JSON.stringify(newOrder),
      })

      expect(mockRequest).toHaveBeenCalledWith('/purchase-orders', {
        method: 'POST',
        body: JSON.stringify(newOrder),
      })
      expect(result).toEqual({ id: 1, ...newOrder })
    })

    it('request can be used for PUT requests', async () => {
      const updateData = { status: 'received', received_at: '2024-01-15T10:00:00Z' }
      mockRequest.mockResolvedValue({ id: 1, ...updateData })

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders/1', {
        method: 'PUT',
        body: JSON.stringify(updateData),
      })

      expect(mockRequest).toHaveBeenCalledWith('/purchase-orders/1', {
        method: 'PUT',
        body: JSON.stringify(updateData),
      })
      expect(result).toEqual({ id: 1, ...updateData })
    })

    it('request can be used for PATCH requests', async () => {
      const patchData = { status: 'shipped' }
      mockRequest.mockResolvedValue({ id: 1, ...patchData })

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders/1/status', {
        method: 'PATCH',
        body: JSON.stringify(patchData),
      })

      expect(mockRequest).toHaveBeenCalledWith('/purchase-orders/1/status', {
        method: 'PATCH',
        body: JSON.stringify(patchData),
      })
      expect(result).toEqual({ id: 1, ...patchData })
    })

    it('request can be used for DELETE requests', async () => {
      mockRequest.mockResolvedValue(null)

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders/1', {
        method: 'DELETE',
      })

      expect(mockRequest).toHaveBeenCalledWith('/purchase-orders/1', {
        method: 'DELETE',
      })
      expect(result).toBeNull()
    })
  })

  describe('error propagation', () => {
    it('propagates errors from request', async () => {
      const error = new Error('Network error')
      mockRequest.mockRejectedValue(error)

      const { request } = useProcurementApi()

      await expect(request('/purchase-orders')).rejects.toThrow('Network error')
    })

    it('propagates validation errors', async () => {
      const error = new Error('Invalid supplier')
      mockRequest.mockRejectedValue(error)

      const { request } = useProcurementApi()

      await expect(
        request('/purchase-orders', {
          method: 'POST',
          body: JSON.stringify({ supplier_id: null }),
        }),
      ).rejects.toThrow('Invalid supplier')
    })

    it('propagates authorization errors', async () => {
      const error = new Error('Unauthorized')
      mockRequest.mockRejectedValue(error)

      const { request } = useProcurementApi()

      await expect(request('/purchase-orders')).rejects.toThrow('Unauthorized')
    })
  })

  describe('edge cases', () => {
    it('normalizeText handles strings with only whitespace characters', () => {
      const { normalizeText } = useProcurementApi()

      expect(normalizeText('\t\n\r ')).toBeNull()
      expect(normalizeText('  \t  ')).toBeNull()
    })

    it('toNumber handles edge numeric cases', () => {
      const { toNumber } = useProcurementApi()

      expect(toNumber('-5')).toBe(-5)
      expect(toNumber('0.001')).toBe(0.001)
      expect(toNumber('1e10')).toBe(1e10)
    })

    it('formatCurrency handles large amounts', () => {
      const { formatCurrency } = useProcurementApi()

      expect(formatCurrency(1_000_000, 'USD')).toBe('10000.00 USD')
      expect(formatCurrency(123_456_789, 'USD')).toBe('1234567.89 USD')
    })

    it('formatCurrency handles small amounts', () => {
      const { formatCurrency } = useProcurementApi()

      expect(formatCurrency(1, 'USD')).toBe('0.01 USD')
      expect(formatCurrency(5, 'EUR')).toBe('0.05 EUR')
    })

    it('formatCurrency handles zero correctly', () => {
      const { formatCurrency } = useProcurementApi()

      expect(formatCurrency(0, 'USD')).toBe('0.00 USD')
    })
  })

  describe('typical procurement workflows', () => {
    it('can fetch suppliers', async () => {
      mockRequest.mockResolvedValue([
        { id: 1, name: 'Malt Supplier Inc' },
        { id: 2, name: 'Hops Direct' },
      ])

      const { request } = useProcurementApi()
      const result = await request('/suppliers')

      expect(mockRequest).toHaveBeenCalledWith('/suppliers')
      expect(result).toHaveLength(2)
    })

    it('can create a purchase order', async () => {
      const orderData = {
        supplier_id: 1,
        items: [
          { ingredient_id: 1, quantity: 100, unit_price_cents: 500 },
        ],
        total_cents: 50_000,
        currency: 'USD',
      }
      mockRequest.mockResolvedValue({ id: 1, status: 'pending', ...orderData })

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders', {
        method: 'POST',
        body: JSON.stringify(orderData),
      })

      expect(result).toHaveProperty('id', 1)
      expect(result).toHaveProperty('status', 'pending')
    })

    it('can update order status', async () => {
      mockRequest.mockResolvedValue({ id: 1, status: 'received' })

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders/1/status', {
        method: 'PATCH',
        body: JSON.stringify({ status: 'received' }),
      })

      expect(result).toHaveProperty('status', 'received')
    })
  })
})
