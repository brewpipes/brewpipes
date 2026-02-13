import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useProductionApi } from '@/composables/useProductionApi'

// Mock useApiClient
const mockRequest = vi.fn()

vi.mock('@/composables/useApiClient', () => ({
  useApiClient: () => ({
    request: mockRequest,
  }),
}))

describe('useProductionApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('utility functions', () => {
    it('normalizeText trims whitespace and returns null for empty strings', () => {
      const { normalizeText } = useProductionApi()

      expect(normalizeText('  hello  ')).toBe('hello')
      expect(normalizeText('test')).toBe('test')
      expect(normalizeText('   ')).toBeNull()
      expect(normalizeText('')).toBeNull()
    })

    it('normalizeDateTime converts to ISO string', () => {
      const { normalizeDateTime } = useProductionApi()

      const result = normalizeDateTime('2024-01-15T10:30:00')
      expect(result).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)

      expect(normalizeDateTime('')).toBeNull()
    })

    it('toNumber parses numeric values correctly', () => {
      const { toNumber } = useProductionApi()

      expect(toNumber('42')).toBe(42)
      expect(toNumber(42)).toBe(42)
      expect(toNumber('3.14')).toBe(3.14)
      expect(toNumber('')).toBeNull()
      expect(toNumber(null)).toBeNull()
      expect(toNumber('not a number')).toBeNull()
    })

    it('formatDateTime formats dates correctly', () => {
      const { formatDateTime } = useProductionApi()

      // Test with a valid date
      const result = formatDateTime('2024-01-15T10:30:00Z')
      expect(result).toContain('2024')

      expect(formatDateTime(null)).toBe('n/a')
      expect(formatDateTime(undefined)).toBe('n/a')
      expect(formatDateTime('')).toBe('n/a')
    })
  })

  describe('Styles API', () => {
    it('getStyles calls correct endpoint', async () => {
      mockRequest.mockResolvedValue([{ uuid: 'style-uuid-1', name: 'IPA' }])

      const { getStyles } = useProductionApi()
      const result = await getStyles()

      expect(mockRequest).toHaveBeenCalledWith('/styles')
      expect(result).toEqual([{ uuid: 'style-uuid-1', name: 'IPA' }])
    })

    it('getStyle calls correct endpoint with uuid', async () => {
      mockRequest.mockResolvedValue({ uuid: 'style-uuid-1', name: 'IPA' })

      const { getStyle } = useProductionApi()
      const result = await getStyle('style-uuid-1')

      expect(mockRequest).toHaveBeenCalledWith('/styles/style-uuid-1')
      expect(result).toEqual({ uuid: 'style-uuid-1', name: 'IPA' })
    })

    it('createStyle sends POST with correct body', async () => {
      mockRequest.mockResolvedValue({ uuid: 'style-uuid-1', name: 'Stout' })

      const { createStyle } = useProductionApi()
      const result = await createStyle({ name: 'Stout' })

      expect(mockRequest).toHaveBeenCalledWith('/styles', {
        method: 'POST',
        body: JSON.stringify({ name: 'Stout' }),
      })
      expect(result).toEqual({ uuid: 'style-uuid-1', name: 'Stout' })
    })
  })

  describe('Recipes API', () => {
    it('getRecipes calls correct endpoint', async () => {
      mockRequest.mockResolvedValue([{ uuid: 'abc-123', name: 'House IPA' }])

      const { getRecipes } = useProductionApi()
      const result = await getRecipes()

      expect(mockRequest).toHaveBeenCalledWith('/recipes')
      expect(result).toEqual([{ uuid: 'abc-123', name: 'House IPA' }])
    })

    it('getRecipe calls correct endpoint with uuid', async () => {
      mockRequest.mockResolvedValue({ uuid: 'abc-123', name: 'House IPA' })

      const { getRecipe } = useProductionApi()
      const result = await getRecipe('abc-123')

      expect(mockRequest).toHaveBeenCalledWith('/recipes/abc-123')
      expect(result).toEqual({ uuid: 'abc-123', name: 'House IPA' })
    })

    it('createRecipe sends POST with correct body', async () => {
      const recipeData = { name: 'New Recipe', style_uuid: 'style-uuid-1', notes: 'Test notes' }
      mockRequest.mockResolvedValue({ uuid: 'def-456', ...recipeData })

      const { createRecipe } = useProductionApi()
      const result = await createRecipe(recipeData)

      expect(mockRequest).toHaveBeenCalledWith('/recipes', {
        method: 'POST',
        body: JSON.stringify(recipeData),
      })
      expect(result).toEqual({ uuid: 'def-456', ...recipeData })
    })

    it('updateRecipe sends PATCH with correct body', async () => {
      const updateData = { name: 'Updated Recipe', style_uuid: 'style-uuid-2' }
      mockRequest.mockResolvedValue({ uuid: 'abc-123', ...updateData })

      const { updateRecipe } = useProductionApi()
      const result = await updateRecipe('abc-123', updateData)

      expect(mockRequest).toHaveBeenCalledWith('/recipes/abc-123', {
        method: 'PATCH',
        body: JSON.stringify(updateData),
      })
      expect(result).toEqual({ uuid: 'abc-123', ...updateData })
    })

    it('deleteRecipe sends DELETE to correct endpoint', async () => {
      mockRequest.mockResolvedValue(undefined)

      const { deleteRecipe } = useProductionApi()
      await deleteRecipe('abc-123')

      expect(mockRequest).toHaveBeenCalledWith('/recipes/abc-123', {
        method: 'DELETE',
      })
    })
  })

  describe('Vessels API', () => {
    it('getVessels calls correct endpoint', async () => {
      mockRequest.mockResolvedValue([{ uuid: 'vessel-uuid-1', name: 'Fermenter 1' }])

      const { getVessels } = useProductionApi()
      const result = await getVessels()

      expect(mockRequest).toHaveBeenCalledWith('/vessels')
      expect(result).toEqual([{ uuid: 'vessel-uuid-1', name: 'Fermenter 1' }])
    })

    it('getVessel calls correct endpoint with uuid', async () => {
      mockRequest.mockResolvedValue({ uuid: 'vessel-uuid-1', name: 'Fermenter 1' })

      const { getVessel } = useProductionApi()
      const result = await getVessel('vessel-uuid-1')

      expect(mockRequest).toHaveBeenCalledWith('/vessels/vessel-uuid-1')
      expect(result).toEqual({ uuid: 'vessel-uuid-1', name: 'Fermenter 1' })
    })
  })

  describe('Volumes API', () => {
    it('getVolumes calls correct endpoint', async () => {
      mockRequest.mockResolvedValue([{ uuid: 'vol-uuid-1', amount: 10, amount_unit: 'bbl' }])

      const { getVolumes } = useProductionApi()
      const result = await getVolumes()

      expect(mockRequest).toHaveBeenCalledWith('/volumes')
      expect(result).toEqual([{ uuid: 'vol-uuid-1', amount: 10, amount_unit: 'bbl' }])
    })

    it('getVolume calls correct endpoint with uuid', async () => {
      mockRequest.mockResolvedValue({ uuid: 'vol-uuid-1', amount: 10, amount_unit: 'bbl' })

      const { getVolume } = useProductionApi()
      const result = await getVolume('vol-uuid-1')

      expect(mockRequest).toHaveBeenCalledWith('/volumes/vol-uuid-1')
      expect(result).toEqual({ uuid: 'vol-uuid-1', amount: 10, amount_unit: 'bbl' })
    })

    it('createVolume sends POST with correct body', async () => {
      const volumeData = { amount: 15, amount_unit: 'bbl' as const }
      mockRequest.mockResolvedValue({ uuid: 'vol-uuid-2', ...volumeData })

      const { createVolume } = useProductionApi()
      const result = await createVolume(volumeData)

      expect(mockRequest).toHaveBeenCalledWith('/volumes', {
        method: 'POST',
        body: JSON.stringify(volumeData),
      })
      expect(result).toEqual({ uuid: 'vol-uuid-2', ...volumeData })
    })
  })

  describe('Brew Sessions API', () => {
    it('getBrewSessions calls correct endpoint with batch_uuid query param', async () => {
      mockRequest.mockResolvedValue([{ uuid: 'session-uuid-1', batch_uuid: 'batch-uuid-5' }])

      const { getBrewSessions } = useProductionApi()
      const result = await getBrewSessions('batch-uuid-5')

      expect(mockRequest).toHaveBeenCalledWith('/brew-sessions?batch_uuid=batch-uuid-5')
      expect(result).toEqual([{ uuid: 'session-uuid-1', batch_uuid: 'batch-uuid-5' }])
    })

    it('getBrewSession calls correct endpoint with uuid', async () => {
      mockRequest.mockResolvedValue({ uuid: 'session-uuid-1', batch_uuid: 'batch-uuid-5' })

      const { getBrewSession } = useProductionApi()
      const result = await getBrewSession('session-uuid-1')

      expect(mockRequest).toHaveBeenCalledWith('/brew-sessions/session-uuid-1')
      expect(result).toEqual({ uuid: 'session-uuid-1', batch_uuid: 'batch-uuid-5' })
    })

    it('createBrewSession sends POST with correct body', async () => {
      const sessionData = { batch_uuid: 'batch-uuid-5', brewed_at: '2024-01-15T10:00:00Z' }
      mockRequest.mockResolvedValue({ uuid: 'session-uuid-1', ...sessionData })

      const { createBrewSession } = useProductionApi()
      const result = await createBrewSession(sessionData)

      expect(mockRequest).toHaveBeenCalledWith('/brew-sessions', {
        method: 'POST',
        body: JSON.stringify(sessionData),
      })
      expect(result).toEqual({ uuid: 'session-uuid-1', ...sessionData })
    })

    it('updateBrewSession sends PUT with correct body', async () => {
      const updateData = { batch_uuid: 'batch-uuid-5', brewed_at: '2024-01-16T10:00:00Z', notes: 'Updated' }
      mockRequest.mockResolvedValue({ uuid: 'session-uuid-1', ...updateData })

      const { updateBrewSession } = useProductionApi()
      const result = await updateBrewSession('session-uuid-1', updateData)

      expect(mockRequest).toHaveBeenCalledWith('/brew-sessions/session-uuid-1', {
        method: 'PUT',
        body: JSON.stringify(updateData),
      })
      expect(result).toEqual({ uuid: 'session-uuid-1', ...updateData })
    })
  })

  describe('Additions API', () => {
    it('getAdditionsByVolume calls correct endpoint with volume_uuid query param', async () => {
      mockRequest.mockResolvedValue([{ uuid: 'add-uuid-1', volume_uuid: 'vol-uuid-3', addition_type: 'hop' }])

      const { getAdditionsByVolume } = useProductionApi()
      const result = await getAdditionsByVolume('vol-uuid-3')

      expect(mockRequest).toHaveBeenCalledWith('/additions?volume_uuid=vol-uuid-3')
      expect(result).toEqual([{ uuid: 'add-uuid-1', volume_uuid: 'vol-uuid-3', addition_type: 'hop' }])
    })

    it('createAddition sends POST with correct body', async () => {
      const additionData = {
        volume_uuid: 'vol-uuid-3',
        addition_type: 'hop' as const,
        amount: 5,
        amount_unit: 'ml' as const,
      }
      mockRequest.mockResolvedValue({ uuid: 'add-uuid-1', ...additionData })

      const { createAddition } = useProductionApi()
      const result = await createAddition(additionData)

      expect(mockRequest).toHaveBeenCalledWith('/additions', {
        method: 'POST',
        body: JSON.stringify(additionData),
      })
      expect(result).toEqual({ uuid: 'add-uuid-1', ...additionData })
    })
  })

  describe('Measurements API', () => {
    it('getMeasurementsByVolume calls correct endpoint with volume_uuid query param', async () => {
      mockRequest.mockResolvedValue([{ uuid: 'meas-uuid-1', volume_uuid: 'vol-uuid-3', kind: 'gravity' }])

      const { getMeasurementsByVolume } = useProductionApi()
      const result = await getMeasurementsByVolume('vol-uuid-3')

      expect(mockRequest).toHaveBeenCalledWith('/measurements?volume_uuid=vol-uuid-3')
      expect(result).toEqual([{ uuid: 'meas-uuid-1', volume_uuid: 'vol-uuid-3', kind: 'gravity' }])
    })

    it('createMeasurement sends POST with correct body', async () => {
      const measurementData = {
        volume_uuid: 'vol-uuid-3',
        kind: 'gravity',
        value: 1.05,
        unit: 'SG',
      }
      mockRequest.mockResolvedValue({ uuid: 'meas-uuid-1', ...measurementData })

      const { createMeasurement } = useProductionApi()
      const result = await createMeasurement(measurementData)

      expect(mockRequest).toHaveBeenCalledWith('/measurements', {
        method: 'POST',
        body: JSON.stringify(measurementData),
      })
      expect(result).toEqual({ uuid: 'meas-uuid-1', ...measurementData })
    })
  })

  describe('Batch Summary API', () => {
    it('getBatchSummary calls correct endpoint with uuid', async () => {
      const summary = {
        uuid: 'batch-uuid-1',
        short_name: 'Batch 001',
        original_gravity: 1.055,
        final_gravity: 1.012,
      }
      mockRequest.mockResolvedValue(summary)

      const { getBatchSummary } = useProductionApi()
      const result = await getBatchSummary('batch-uuid-1')

      expect(mockRequest).toHaveBeenCalledWith('/batches/batch-uuid-1/summary')
      expect(result).toEqual(summary)
    })
  })

  describe('Occupancy API', () => {
    it('getActiveOccupancies calls correct endpoint with active query param', async () => {
      mockRequest.mockResolvedValue([{ uuid: 'occ-uuid-1', status: 'fermenting' }])

      const { getActiveOccupancies } = useProductionApi()
      const result = await getActiveOccupancies()

      expect(mockRequest).toHaveBeenCalledWith('/occupancies?active=true')
      expect(result).toEqual([{ uuid: 'occ-uuid-1', status: 'fermenting' }])
    })

    it('getOccupancy calls correct endpoint with uuid', async () => {
      mockRequest.mockResolvedValue({ uuid: 'occ-uuid-1', status: 'fermenting' })

      const { getOccupancy } = useProductionApi()
      const result = await getOccupancy('occ-uuid-1')

      expect(mockRequest).toHaveBeenCalledWith('/occupancies/occ-uuid-1')
      expect(result).toEqual({ uuid: 'occ-uuid-1', status: 'fermenting' })
    })

    it('updateOccupancyStatus sends PATCH with correct body', async () => {
      mockRequest.mockResolvedValue({ uuid: 'occ-uuid-1', status: 'conditioning' })

      const { updateOccupancyStatus } = useProductionApi()
      const result = await updateOccupancyStatus('occ-uuid-1', 'conditioning')

      expect(mockRequest).toHaveBeenCalledWith('/occupancies/occ-uuid-1/status', {
        method: 'PATCH',
        body: JSON.stringify({ status: 'conditioning' }),
      })
      expect(result).toEqual({ uuid: 'occ-uuid-1', status: 'conditioning' })
    })
  })

  describe('error propagation', () => {
    it('propagates errors from request', async () => {
      const error = new Error('Network error')
      mockRequest.mockRejectedValue(error)

      const { getStyles } = useProductionApi()

      await expect(getStyles()).rejects.toThrow('Network error')
    })

    it('propagates errors from POST requests', async () => {
      const error = new Error('Validation failed')
      mockRequest.mockRejectedValue(error)

      const { createStyle } = useProductionApi()

      await expect(createStyle({ name: '' })).rejects.toThrow('Validation failed')
    })
  })

  describe('exposed properties', () => {
    it('exposes apiBase', () => {
      const api = useProductionApi()
      expect(api.apiBase).toBeDefined()
    })

    it('exposes request function', () => {
      const api = useProductionApi()
      expect(api.request).toBeDefined()
      expect(typeof api.request).toBe('function')
    })
  })
})
