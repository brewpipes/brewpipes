import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Module-level flag to prevent multiple concurrent 401 redirects
let redirecting = false

function isJsonResponse (response: Response) {
  const contentType = response.headers.get('content-type') ?? ''
  return contentType.includes('application/json')
}

async function readErrorMessage (response: Response) {
  const message = await response.text()
  return message || `Request failed with ${response.status}`
}

export function useApiClient (baseUrl: string) {
  const authStore = useAuthStore()
  const router = useRouter()

  const buildHeaders = (init: RequestInit, token: string | null) => {
    const headers = new Headers(init.headers ?? {})
    const isFormData
      = typeof FormData !== 'undefined' && init.body instanceof FormData
    if (init.body && !headers.has('Content-Type') && !isFormData) {
      headers.set('Content-Type', 'application/json')
    }
    if (token && !headers.has('Authorization')) {
      headers.set('Authorization', `Bearer ${token}`)
    }
    return headers
  }

  const executeRequest = (path: string, init: RequestInit, token: string | null) => {
    const headers = buildHeaders(init, token)
    return fetch(`${baseUrl}${path}`, {
      ...init,
      headers,
    })
  }

  const parseResponse = async <T>(response: Response): Promise<T> => {
    if (response.status === 204) {
      return null as T
    }
    if (isJsonResponse(response)) {
      return response.json() as Promise<T>
    }
    return (await response.text()) as T
  }

  const request = async <T>(path: string, init: RequestInit = {}): Promise<T> => {
    const response = await executeRequest(path, init, authStore.accessToken)

    if (response.status === 401) {
      const refreshed = await authStore.refresh()
      if (refreshed) {
        const retry = await executeRequest(path, init, authStore.accessToken)
        if (!retry.ok) {
          throw new Error(await readErrorMessage(retry))
        }
        return parseResponse<T>(retry)
      }

      if (!redirecting && router.currentRoute.value.path !== '/login') {
        redirecting = true
        router.push({
          path: '/login',
          query: { redirect: router.currentRoute.value.fullPath },
        }).catch(() => {}).finally(() => {
          redirecting = false
        })
      }
      throw new Error(await readErrorMessage(response))
    }

    if (!response.ok) {
      throw new Error(await readErrorMessage(response))
    }

    return parseResponse<T>(response)
  }

  return { request }
}
