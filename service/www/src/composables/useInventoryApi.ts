import { useApiClient } from '@/composables/useApiClient'

const inventoryApiBase = import.meta.env.VITE_INVENTORY_API_URL ?? '/api'

export function useInventoryApi () {
  const { request } = useApiClient(inventoryApiBase)

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

  const formatAmount = (amount: number | null | undefined, unit: string | null | undefined) => {
    if (amount === null || amount === undefined) {
      return 'n/a'
    }
    return `${amount} ${unit ?? ''}`.trim()
  }

  return {
    apiBase: inventoryApiBase,
    request,
    normalizeText,
    normalizeDateTime,
    toNumber,
    formatDateTime,
    formatAmount,
  }
}
