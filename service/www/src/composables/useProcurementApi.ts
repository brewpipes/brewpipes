import type {
  CreatePurchaseOrderRequest,
  PurchaseOrder,
  PurchaseOrderFee,
  PurchaseOrderLine,
  Supplier,
  UpdatePurchaseOrderRequest,
  UpdateSupplierRequest,
} from '@/types'
import { useApiClient } from '@/composables/useApiClient'
import { formatDateTime } from '@/composables/useFormatters'
import { normalizeDateTime, normalizeText, toNumber } from '@/utils/normalize'

const procurementApiBase = import.meta.env.VITE_PROCUREMENT_API_URL ?? '/api'

export function useProcurementApi () {
  const { request } = useApiClient(procurementApiBase)

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

  // Purchase Order Lines API
  const getPurchaseOrderLines = (purchaseOrderUuid?: string) => {
    const query = new URLSearchParams()
    if (purchaseOrderUuid) {
      query.set('purchase_order_uuid', purchaseOrderUuid)
    }
    const path = query.toString() ? `/purchase-order-lines?${query.toString()}` : '/purchase-order-lines'
    return request<PurchaseOrderLine[]>(path)
  }
  const createPurchaseOrderLine = (data: Record<string, unknown>) =>
    request<PurchaseOrderLine>('/purchase-order-lines', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Purchase Order Fees API
  const getPurchaseOrderFees = (purchaseOrderUuid?: string) => {
    const query = new URLSearchParams()
    if (purchaseOrderUuid) {
      query.set('purchase_order_uuid', purchaseOrderUuid)
    }
    const path = query.toString() ? `/purchase-order-fees?${query.toString()}` : '/purchase-order-fees'
    return request<PurchaseOrderFee[]>(path)
  }
  const createPurchaseOrderFee = (data: Record<string, unknown>) =>
    request<PurchaseOrderFee>('/purchase-order-fees', {
      method: 'POST',
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
    // Purchase Order Lines
    getPurchaseOrderLines,
    createPurchaseOrderLine,
    // Purchase Order Fees
    getPurchaseOrderFees,
    createPurchaseOrderFee,
  }
}
