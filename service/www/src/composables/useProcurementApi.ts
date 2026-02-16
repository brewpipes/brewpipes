import type {
  CreatePurchaseOrderFeeRequest,
  CreatePurchaseOrderLineRequest,
  CreatePurchaseOrderRequest,
  PurchaseOrder,
  PurchaseOrderFee,
  PurchaseOrderLine,
  Supplier,
  UpdatePurchaseOrderFeeRequest,
  UpdatePurchaseOrderLineRequest,
  UpdatePurchaseOrderRequest,
  UpdateSupplierRequest,
} from '@/types'
import { useApiClient } from '@/composables/useApiClient'

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
  const createPurchaseOrderLine = (data: CreatePurchaseOrderLineRequest) =>
    request<PurchaseOrderLine>('/purchase-order-lines', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updatePurchaseOrderLine = (uuid: string, data: UpdatePurchaseOrderLineRequest) =>
    request<PurchaseOrderLine>(`/purchase-order-lines/${uuid}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  const deletePurchaseOrderLine = (uuid: string) =>
    request<void>(`/purchase-order-lines/${uuid}`, {
      method: 'DELETE',
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
  const createPurchaseOrderFee = (data: CreatePurchaseOrderFeeRequest) =>
    request<PurchaseOrderFee>('/purchase-order-fees', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updatePurchaseOrderFee = (uuid: string, data: UpdatePurchaseOrderFeeRequest) =>
    request<PurchaseOrderFee>(`/purchase-order-fees/${uuid}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  const deletePurchaseOrderFee = (uuid: string) =>
    request<void>(`/purchase-order-fees/${uuid}`, {
      method: 'DELETE',
    })

  return {
    apiBase: procurementApiBase,
    request,
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
    updatePurchaseOrderLine,
    deletePurchaseOrderLine,
    // Purchase Order Fees
    getPurchaseOrderFees,
    createPurchaseOrderFee,
    updatePurchaseOrderFee,
    deletePurchaseOrderFee,
  }
}
