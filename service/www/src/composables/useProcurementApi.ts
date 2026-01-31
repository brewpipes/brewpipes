const procurementApiBase = import.meta.env.VITE_PROCUREMENT_API_URL ?? '/api'

export function useProcurementApi() {
  const request = async <T>(path: string, init: RequestInit = {}): Promise<T> => {
    const response = await fetch(`${procurementApiBase}${path}`, {
      ...init,
      headers: {
        'Content-Type': 'application/json',
        ...(init.headers ?? {}),
      },
    })

    if (!response.ok) {
      const message = await response.text()
      throw new Error(message || `Request failed with ${response.status}`)
    }

    return response.json() as Promise<T>
  }

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

  return {
    apiBase: procurementApiBase,
    request,
    normalizeText,
    normalizeDateTime,
    toNumber,
    formatDateTime,
    formatCurrency,
  }
}
