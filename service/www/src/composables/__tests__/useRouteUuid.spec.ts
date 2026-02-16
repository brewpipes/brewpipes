import { describe, expect, it, vi } from 'vitest'
import { reactive } from 'vue'

// Mock vue-router with a reactive params object so computed() tracks changes
const mockRoute = reactive({
  params: {} as Record<string, string | string[]>,
})

vi.mock('vue-router', () => ({
  useRoute: () => mockRoute,
}))

import { useRouteUuid } from '../useRouteUuid'

describe('useRouteUuid', () => {
  describe('valid UUID extraction', () => {
    it('returns the uuid param when present as a string', () => {
      mockRoute.params = { uuid: 'abc-123-def-456' }
      const { uuid } = useRouteUuid()
      expect(uuid.value).toBe('abc-123-def-456')
    })

    it('returns a full UUID format string', () => {
      mockRoute.params = { uuid: '550e8400-e29b-41d4-a716-446655440000' }
      const { uuid } = useRouteUuid()
      expect(uuid.value).toBe('550e8400-e29b-41d4-a716-446655440000')
    })

    it('returns a simple string param', () => {
      mockRoute.params = { uuid: 'vessel-1' }
      const { uuid } = useRouteUuid()
      expect(uuid.value).toBe('vessel-1')
    })
  })

  describe('missing or invalid params', () => {
    it('returns null when uuid param is missing', () => {
      mockRoute.params = {}
      const { uuid } = useRouteUuid()
      expect(uuid.value).toBeNull()
    })

    it('returns null when uuid param is empty string', () => {
      mockRoute.params = { uuid: '' }
      const { uuid } = useRouteUuid()
      expect(uuid.value).toBeNull()
    })

    it('returns null when uuid param is whitespace-only', () => {
      mockRoute.params = { uuid: '   ' }
      const { uuid } = useRouteUuid()
      expect(uuid.value).toBeNull()
    })

    it('returns null when uuid param is a tab-only string', () => {
      mockRoute.params = { uuid: '\t' }
      const { uuid } = useRouteUuid()
      expect(uuid.value).toBeNull()
    })

    it('returns null when uuid param is an array', () => {
      // vue-router params can be string | string[]
      mockRoute.params = { uuid: ['abc', 'def'] }
      const { uuid } = useRouteUuid()
      expect(uuid.value).toBeNull()
    })
  })

  describe('reactivity', () => {
    it('updates when route params change', () => {
      mockRoute.params = { uuid: 'initial-uuid' }
      const { uuid } = useRouteUuid()
      expect(uuid.value).toBe('initial-uuid')

      // Simulate route change
      mockRoute.params = { uuid: 'updated-uuid' }
      expect(uuid.value).toBe('updated-uuid')
    })

    it('returns null when uuid param is removed', () => {
      mockRoute.params = { uuid: 'some-uuid' }
      const { uuid } = useRouteUuid()
      expect(uuid.value).toBe('some-uuid')

      // Simulate navigating to a route without uuid
      mockRoute.params = {}
      expect(uuid.value).toBeNull()
    })
  })
})
