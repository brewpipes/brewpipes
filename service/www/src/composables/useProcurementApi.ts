import { useApiClient } from '@/composables/useApiClient'

const procurementApiBase = import.meta.env.VITE_PROCUREMENT_API_URL ?? '/api'

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

export interface CreatePurchaseOrderRequest {
  supplier_uuid: string
  order_number: string
  status?: string | null
  ordered_at?: string | null
  expected_at?: string | null
  notes?: string | null
}

export interface UpdatePurchaseOrderRequest {
  order_number?: string
  status?: string
  ordered_at?: string | null
  expected_at?: string | null
  notes?: string | null
}

export function useProcurementApi () {
  const { request } = useApiClient(procurementApiBase)

  const normalizeText = (value: string) => {
    const trimmed = value.trim()
    return trimmed.length > 0 ? trimmed : null
  }

  const normalizeDateTime = (value: string) => {
    return value ? new Date(value).toISOString() : null
  }

  const toNumber = (value: string | number | null) => {
    if (value === null || value === undefined || value === '') {
      return null
    }
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : null
  }

  const formatDateTime = (value: string | null | undefined) => {
    if (!value) {
      return 'n/a'
    }
    return new Intl.DateTimeFormat('en-US', {
      dateStyle: 'medium',
      timeStyle: 'short',
    }).format(new Date(value))
  }

  const formatCurrency = (cents: number | null | undefined, currency: string | null | undefined) => {
    if (cents === null || cents === undefined) {
      return 'n/a'
    }
    const amount = (cents / 100).toFixed(2)
    return `${amount} ${currency ?? ''}`.trim()
  }

  // Suppliers API
  const getSuppliers = () => request<Supplier[]>('/suppliers')
  const getSupplier = (uuid: string) => request<Supplier>(`/suppliers/${uuid}`)
  const createSupplier = (data: Omit<UpdateSupplierRequest, 'name'> & { name: string }) =>
    request<Supplier>('/suppliers', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updateSupplier = (uuid: string, data: UpdateSupplierRequest) =>
    request<Supplier>(`/suppliers/${uuid}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })

  // Purchase Orders API
  const getPurchaseOrders = (supplierUuid?: string) => {
    const query = new URLSearchParams()
    if (supplierUuid) {
      query.set('supplier_uuid', supplierUuid)
    }
    const path = query.toString() ? `/purchase-orders?${query.toString()}` : '/purchase-orders'
    return request<PurchaseOrder[]>(path)
  }
  const getPurchaseOrder = (uuid: string) => request<PurchaseOrder>(`/purchase-orders/${uuid}`)
  const createPurchaseOrder = (data: CreatePurchaseOrderRequest) =>
    request<PurchaseOrder>('/purchase-orders', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updatePurchaseOrder = (uuid: string, data: UpdatePurchaseOrderRequest) =>
    request<PurchaseOrder>(`/purchase-orders/${uuid}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })

  return {
    apiBase: procurementApiBase,
    request,
    normalizeText,
    normalizeDateTime,
    toNumber,
    formatDateTime,
    formatCurrency,
    // Suppliers
    getSuppliers,
    getSupplier,
    createSupplier,
    updateSupplier,
    // Purchase Orders
    getPurchaseOrders,
    getPurchaseOrder,
    createPurchaseOrder,
    updatePurchaseOrder,
  }
}
