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

  describe('formatCurrency', () => {
    it('formats cents to currency string', () => {
      const { formatCurrency } = useProcurementApi()

      expect(formatCurrency(1000, 'USD')).toBe('10.00 USD')
      expect(formatCurrency(1550, 'EUR')).toBe('15.50 EUR')
      expect(formatCurrency(99, 'GBP')).toBe('0.99 GBP')
      expect(formatCurrency(0, 'USD')).toBe('0.00 USD')
    })

    it('handles null/undefined currency', () => {
      const { formatCurrency } = useProcurementApi()

      expect(formatCurrency(1000, null)).toBe('10.00')
      expect(formatCurrency(1000, undefined)).toBe('10.00')
    })

    it('returns n/a for null/undefined cents', () => {
      const { formatCurrency } = useProcurementApi()

      expect(formatCurrency(null, 'USD')).toBe('n/a')
      expect(formatCurrency(undefined, 'USD')).toBe('n/a')
    })

    it('handles large amounts', () => {
      const { formatCurrency } = useProcurementApi()

      expect(formatCurrency(1_000_000, 'USD')).toBe('10000.00 USD')
      expect(formatCurrency(123_456_789, 'USD')).toBe('1234567.89 USD')
    })

    it('handles small amounts', () => {
      const { formatCurrency } = useProcurementApi()

      expect(formatCurrency(1, 'USD')).toBe('0.01 USD')
      expect(formatCurrency(5, 'EUR')).toBe('0.05 EUR')
    })

    it('handles zero correctly', () => {
      const { formatCurrency } = useProcurementApi()

      expect(formatCurrency(0, 'USD')).toBe('0.00 USD')
    })
  })

  describe('dollarsToCents', () => {
    it('converts dollar amounts to cents', () => {
      const { dollarsToCents } = useProcurementApi()

      expect(dollarsToCents(12.50)).toBe(1250)
      expect(dollarsToCents(15.00)).toBe(1500)
      expect(dollarsToCents(0.99)).toBe(99)
      expect(dollarsToCents(0)).toBe(0)
    })

    it('rounds to nearest cent', () => {
      const { dollarsToCents } = useProcurementApi()

      expect(dollarsToCents(12.555)).toBe(1256)
      expect(dollarsToCents(12.554)).toBe(1255)
    })

    it('handles null', () => {
      const { dollarsToCents } = useProcurementApi()

      expect(dollarsToCents(null)).toBeNull()
    })

    it('handles whole dollar amounts', () => {
      const { dollarsToCents } = useProcurementApi()

      expect(dollarsToCents(100)).toBe(10000)
      expect(dollarsToCents(1)).toBe(100)
    })

    it('handles small amounts', () => {
      const { dollarsToCents } = useProcurementApi()

      expect(dollarsToCents(0.01)).toBe(1)
      expect(dollarsToCents(0.05)).toBe(5)
    })
  })

  describe('centsToDollars', () => {
    it('converts cents to dollar amounts', () => {
      const { centsToDollars } = useProcurementApi()

      expect(centsToDollars(1250)).toBe(12.50)
      expect(centsToDollars(1500)).toBe(15.00)
      expect(centsToDollars(99)).toBe(0.99)
      expect(centsToDollars(0)).toBe(0)
    })

    it('handles null', () => {
      const { centsToDollars } = useProcurementApi()

      expect(centsToDollars(null)).toBeNull()
    })

    it('handles large amounts', () => {
      const { centsToDollars } = useProcurementApi()

      expect(centsToDollars(1_000_000)).toBe(10000)
      expect(centsToDollars(123_456_789)).toBe(1234567.89)
    })

    it('handles small amounts', () => {
      const { centsToDollars } = useProcurementApi()

      expect(centsToDollars(1)).toBe(0.01)
      expect(centsToDollars(5)).toBe(0.05)
    })

    it('round-trips with dollarsToCents', () => {
      const { dollarsToCents, centsToDollars } = useProcurementApi()

      expect(centsToDollars(dollarsToCents(12.50)!)).toBe(12.50)
      expect(dollarsToCents(centsToDollars(1250)!)).toBe(1250)
      expect(centsToDollars(dollarsToCents(0.01)!)).toBe(0.01)
      expect(dollarsToCents(centsToDollars(1)!)).toBe(1)
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

    it('exposes formatCurrency function', () => {
      const api = useProcurementApi()

      expect(typeof api.formatCurrency).toBe('function')
    })

    it('exposes dollarsToCents function', () => {
      const api = useProcurementApi()
      expect(typeof api.dollarsToCents).toBe('function')
    })

    it('exposes centsToDollars function', () => {
      const api = useProcurementApi()
      expect(typeof api.centsToDollars).toBe('function')
    })
  })

  describe('request function usage', () => {
    it('request can be used for GET requests', async () => {
      mockRequest.mockResolvedValue([{ uuid: 'po-uuid-1', supplier: 'Malt Co' }])

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders')

      expect(mockRequest).toHaveBeenCalledWith('/purchase-orders')
      expect(result).toEqual([{ uuid: 'po-uuid-1', supplier: 'Malt Co' }])
    })

    it('request can be used for POST requests', async () => {
      const newOrder = { supplier_uuid: 'sup-uuid-1', total_cents: 50_000 }
      mockRequest.mockResolvedValue({ uuid: 'po-uuid-1', ...newOrder })

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders', {
        method: 'POST',
        body: JSON.stringify(newOrder),
      })

      expect(mockRequest).toHaveBeenCalledWith('/purchase-orders', {
        method: 'POST',
        body: JSON.stringify(newOrder),
      })
      expect(result).toEqual({ uuid: 'po-uuid-1', ...newOrder })
    })

    it('request can be used for PUT requests', async () => {
      const updateData = { status: 'received', received_at: '2024-01-15T10:00:00Z' }
      mockRequest.mockResolvedValue({ uuid: 'po-uuid-1', ...updateData })

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders/po-uuid-1', {
        method: 'PUT',
        body: JSON.stringify(updateData),
      })

      expect(mockRequest).toHaveBeenCalledWith('/purchase-orders/po-uuid-1', {
        method: 'PUT',
        body: JSON.stringify(updateData),
      })
      expect(result).toEqual({ uuid: 'po-uuid-1', ...updateData })
    })

    it('request can be used for PATCH requests', async () => {
      const patchData = { status: 'shipped' }
      mockRequest.mockResolvedValue({ uuid: 'po-uuid-1', ...patchData })

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders/po-uuid-1/status', {
        method: 'PATCH',
        body: JSON.stringify(patchData),
      })

      expect(mockRequest).toHaveBeenCalledWith('/purchase-orders/po-uuid-1/status', {
        method: 'PATCH',
        body: JSON.stringify(patchData),
      })
      expect(result).toEqual({ uuid: 'po-uuid-1', ...patchData })
    })

    it('request can be used for DELETE requests', async () => {
      mockRequest.mockResolvedValue(null)

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders/po-uuid-1', {
        method: 'DELETE',
      })

      expect(mockRequest).toHaveBeenCalledWith('/purchase-orders/po-uuid-1', {
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
          body: JSON.stringify({ supplier_uuid: null }),
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

  describe('typical procurement workflows', () => {
    it('can fetch suppliers', async () => {
      mockRequest.mockResolvedValue([
        { uuid: 'sup-uuid-1', name: 'Malt Supplier Inc' },
        { uuid: 'sup-uuid-2', name: 'Hops Direct' },
      ])

      const { request } = useProcurementApi()
      const result = await request('/suppliers')

      expect(mockRequest).toHaveBeenCalledWith('/suppliers')
      expect(result).toHaveLength(2)
    })

    it('can create a purchase order', async () => {
      const orderData = {
        supplier_uuid: 'sup-uuid-1',
        items: [
          { ingredient_uuid: 'ing-uuid-1', quantity: 100, unit_price_cents: 500 },
        ],
        total_cents: 50_000,
        currency: 'USD',
      }
      mockRequest.mockResolvedValue({ uuid: 'po-uuid-1', status: 'pending', ...orderData })

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders', {
        method: 'POST',
        body: JSON.stringify(orderData),
      })

      expect(result).toHaveProperty('uuid', 'po-uuid-1')
      expect(result).toHaveProperty('status', 'pending')
    })

    it('can update order status', async () => {
      mockRequest.mockResolvedValue({ uuid: 'po-uuid-1', status: 'received' })

      const { request } = useProcurementApi()
      const result = await request('/purchase-orders/po-uuid-1/status', {
        method: 'PATCH',
        body: JSON.stringify({ status: 'received' }),
      })

      expect(result).toHaveProperty('status', 'received')
    })
  })
})
