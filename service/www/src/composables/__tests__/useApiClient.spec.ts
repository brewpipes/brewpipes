import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { useApiClient } from '@/composables/useApiClient'

// Mock vue-router
const mockPush = vi.fn().mockResolvedValue(undefined)
const mockCurrentRoute = {
  value: {
    path: '/dashboard',
    fullPath: '/dashboard?tab=1',
  },
}

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush,
    currentRoute: mockCurrentRoute,
  }),
}))

// Mock auth store
const mockAuthStore = {
  accessToken: null as string | null,
  refresh: vi.fn(),
}

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => mockAuthStore,
}))

describe('useApiClient', () => {
  const baseUrl = 'http://localhost:8080/api'
  let originalFetch: typeof global.fetch

  beforeEach(() => {
    originalFetch = global.fetch
    vi.clearAllMocks()
    mockPush.mockResolvedValue(undefined)
    mockAuthStore.accessToken = null
    mockAuthStore.refresh.mockReset()
    mockCurrentRoute.value = {
      path: '/dashboard',
      fullPath: '/dashboard?tab=1',
    }
  })

  afterEach(() => {
    global.fetch = originalFetch
  })

  describe('request headers', () => {
    it('includes Authorization header when token exists', async () => {
      mockAuthStore.accessToken = 'test-token-123'

      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: () => Promise.resolve({ data: 'test' }),
      })

      const { request } = useApiClient(baseUrl)
      await request('/test')

      expect(global.fetch).toHaveBeenCalledWith(
        `${baseUrl}/test`,
        expect.objectContaining({
          headers: expect.any(Headers),
        }),
      )

      const callArgs = vi.mocked(global.fetch).mock.calls[0]
      const headers = callArgs[1]?.headers as Headers
      expect(headers.get('Authorization')).toBe('Bearer test-token-123')
    })

    it('does not include Authorization header when token is null', async () => {
      mockAuthStore.accessToken = null

      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: () => Promise.resolve({ data: 'test' }),
      })

      const { request } = useApiClient(baseUrl)
      await request('/test')

      const callArgs = vi.mocked(global.fetch).mock.calls[0]
      const headers = callArgs[1]?.headers as Headers
      expect(headers.get('Authorization')).toBeNull()
    })

    it('sets Content-Type to application/json for requests with body', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: () => Promise.resolve({ id: 1 }),
      })

      const { request } = useApiClient(baseUrl)
      await request('/test', {
        method: 'POST',
        body: JSON.stringify({ name: 'test' }),
      })

      const callArgs = vi.mocked(global.fetch).mock.calls[0]
      const headers = callArgs[1]?.headers as Headers
      expect(headers.get('Content-Type')).toBe('application/json')
    })

    it('does not override existing Content-Type header', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'text/plain' }),
        text: () => Promise.resolve('ok'),
      })

      const { request } = useApiClient(baseUrl)
      await request('/test', {
        method: 'POST',
        body: 'plain text',
        headers: { 'Content-Type': 'text/plain' },
      })

      const callArgs = vi.mocked(global.fetch).mock.calls[0]
      const headers = callArgs[1]?.headers as Headers
      expect(headers.get('Content-Type')).toBe('text/plain')
    })

    it('does not set Content-Type for FormData', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: () => Promise.resolve({ success: true }),
      })

      const formData = new FormData()
      formData.append('file', new Blob(['test']), 'test.txt')

      const { request } = useApiClient(baseUrl)
      await request('/upload', {
        method: 'POST',
        body: formData,
      })

      const callArgs = vi.mocked(global.fetch).mock.calls[0]
      const headers = callArgs[1]?.headers as Headers
      // FormData should not have Content-Type set (browser sets it with boundary)
      expect(headers.has('Content-Type')).toBe(false)
    })
  })

  describe('401 response handling', () => {
    it('attempts token refresh on 401 response', async () => {
      mockAuthStore.accessToken = 'expired-token'
      mockAuthStore.refresh.mockResolvedValue(true)

      let callCount = 0
      global.fetch = vi.fn().mockImplementation(() => {
        callCount++
        if (callCount === 1) {
          return Promise.resolve({
            ok: false,
            status: 401,
            headers: new Headers({ 'content-type': 'text/plain' }),
            text: () => Promise.resolve('Unauthorized'),
          })
        }
        // After refresh, return success
        mockAuthStore.accessToken = 'new-token'
        return Promise.resolve({
          ok: true,
          status: 200,
          headers: new Headers({ 'content-type': 'application/json' }),
          json: () => Promise.resolve({ data: 'success' }),
        })
      })

      const { request } = useApiClient(baseUrl)
      const result = await request('/protected')

      expect(mockAuthStore.refresh).toHaveBeenCalledTimes(1)
      expect(global.fetch).toHaveBeenCalledTimes(2)
      expect(result).toEqual({ data: 'success' })
    })

    it('retries request with new token after successful refresh', async () => {
      mockAuthStore.accessToken = 'expired-token'
      mockAuthStore.refresh.mockImplementation(async () => {
        mockAuthStore.accessToken = 'refreshed-token'
        return true
      })

      let callCount = 0
      global.fetch = vi.fn().mockImplementation(() => {
        callCount++
        if (callCount === 1) {
          return Promise.resolve({
            ok: false,
            status: 401,
            headers: new Headers({ 'content-type': 'text/plain' }),
            text: () => Promise.resolve('Unauthorized'),
          })
        }
        return Promise.resolve({
          ok: true,
          status: 200,
          headers: new Headers({ 'content-type': 'application/json' }),
          json: () => Promise.resolve({ data: 'retried' }),
        })
      })

      const { request } = useApiClient(baseUrl)
      await request('/protected')

      // Verify second call uses refreshed token
      const secondCallArgs = vi.mocked(global.fetch).mock.calls[1]
      const headers = secondCallArgs[1]?.headers as Headers
      expect(headers.get('Authorization')).toBe('Bearer refreshed-token')
    })

    it('redirects to login when refresh fails', async () => {
      mockAuthStore.accessToken = 'expired-token'
      mockAuthStore.refresh.mockResolvedValue(false)

      global.fetch = vi.fn().mockResolvedValue({
        ok: false,
        status: 401,
        headers: new Headers({ 'content-type': 'text/plain' }),
        text: () => Promise.resolve('Unauthorized'),
      })

      const { request } = useApiClient(baseUrl)

      await expect(request('/protected')).rejects.toThrow('Unauthorized')

      expect(mockPush).toHaveBeenCalledWith({
        path: '/login',
        query: { redirect: '/dashboard?tab=1' },
      })
    })

    it('does not redirect if already on login page', async () => {
      mockAuthStore.accessToken = 'expired-token'
      mockAuthStore.refresh.mockResolvedValue(false)
      mockCurrentRoute.value = {
        path: '/login',
        fullPath: '/login',
      }

      global.fetch = vi.fn().mockResolvedValue({
        ok: false,
        status: 401,
        headers: new Headers({ 'content-type': 'text/plain' }),
        text: () => Promise.resolve('Unauthorized'),
      })

      const { request } = useApiClient(baseUrl)

      await expect(request('/protected')).rejects.toThrow('Unauthorized')

      expect(mockPush).not.toHaveBeenCalled()
    })

    it('throws error when retry after refresh also fails', async () => {
      mockAuthStore.accessToken = 'expired-token'
      mockAuthStore.refresh.mockImplementation(async () => {
        mockAuthStore.accessToken = 'refreshed-token'
        return true
      })

      global.fetch = vi.fn().mockResolvedValue({
        ok: false,
        status: 401,
        headers: new Headers({ 'content-type': 'text/plain' }),
        text: () => Promise.resolve('Still unauthorized'),
      })

      const { request } = useApiClient(baseUrl)

      await expect(request('/protected')).rejects.toThrow('Still unauthorized')
    })
  })

  describe('response parsing', () => {
    it('returns null for 204 No Content responses', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 204,
        headers: new Headers(),
      })

      const { request } = useApiClient(baseUrl)
      const result = await request('/delete')

      expect(result).toBeNull()
    })

    it('parses JSON responses', async () => {
      const responseData = { id: 1, name: 'Test' }
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: () => Promise.resolve(responseData),
      })

      const { request } = useApiClient(baseUrl)
      const result = await request('/data')

      expect(result).toEqual(responseData)
    })

    it('parses JSON responses with charset', async () => {
      const responseData = { message: 'Hello' }
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'application/json; charset=utf-8' }),
        json: () => Promise.resolve(responseData),
      })

      const { request } = useApiClient(baseUrl)
      const result = await request('/data')

      expect(result).toEqual(responseData)
    })

    it('returns text for non-JSON responses', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'text/plain' }),
        text: () => Promise.resolve('Plain text response'),
      })

      const { request } = useApiClient(baseUrl)
      const result = await request('/text')

      expect(result).toBe('Plain text response')
    })
  })

  describe('error handling', () => {
    it('throws error with response text for non-ok responses', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: false,
        status: 400,
        headers: new Headers({ 'content-type': 'text/plain' }),
        text: () => Promise.resolve('Bad request: invalid input'),
      })

      const { request } = useApiClient(baseUrl)

      await expect(request('/bad')).rejects.toThrow('Bad request: invalid input')
    })

    it('throws error with status code when response text is empty', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        headers: new Headers({ 'content-type': 'text/plain' }),
        text: () => Promise.resolve(''),
      })

      const { request } = useApiClient(baseUrl)

      await expect(request('/error')).rejects.toThrow('Request failed with 500')
    })

    it('propagates network errors', async () => {
      global.fetch = vi.fn().mockRejectedValue(new Error('Network error'))

      const { request } = useApiClient(baseUrl)

      await expect(request('/network-fail')).rejects.toThrow('Network error')
    })

    it('handles fetch abort errors', async () => {
      const abortError = new DOMException('The operation was aborted', 'AbortError')
      global.fetch = vi.fn().mockRejectedValue(abortError)

      const { request } = useApiClient(baseUrl)

      await expect(request('/aborted')).rejects.toThrow('The operation was aborted')
    })
  })

  describe('request configuration', () => {
    it('passes through request init options', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: () => Promise.resolve({}),
      })

      const { request } = useApiClient(baseUrl)
      await request('/test', {
        method: 'PUT',
        body: JSON.stringify({ update: true }),
      })

      expect(global.fetch).toHaveBeenCalledWith(
        `${baseUrl}/test`,
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({ update: true }),
        }),
      )
    })

    it('constructs correct URL with base URL and path', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: () => Promise.resolve([]),
      })

      const { request } = useApiClient('https://api.example.com')
      await request('/users/123')

      expect(global.fetch).toHaveBeenCalledWith(
        'https://api.example.com/users/123',
        expect.any(Object),
      )
    })
  })
})
