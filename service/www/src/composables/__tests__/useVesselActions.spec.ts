import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'
import type { UpdateVesselRequest, Vessel } from '@/types'

// Mock useProductionApi
const mockUpdateVessel = vi.fn()

vi.mock('@/composables/useProductionApi', () => ({
  useProductionApi: () => ({
    updateVessel: mockUpdateVessel,
  }),
}))

// Mock useSnackbar
const mockShowNotice = vi.fn()

vi.mock('@/composables/useSnackbar', () => ({
  useSnackbar: () => ({
    showNotice: mockShowNotice,
  }),
}))

import { useVesselActions } from '../useVesselActions'

// Helper to create a mock dialog ref
function createMockDialogRef () {
  const setSaving = vi.fn()
  const clearError = vi.fn()
  const setError = vi.fn()

  const dialogRef = ref({
    setSaving,
    clearError,
    setError,
  })

  return { dialogRef, setSaving, clearError, setError }
}

const sampleVessel: Vessel = {
  uuid: 'vessel-uuid-1',
  name: 'Fermenter 1',
  type: 'fermenter',
  status: 'active',
  capacity: 100,
  capacity_unit: 'gal',
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

const sampleRequest: UpdateVesselRequest = {
  name: 'Fermenter 1 Updated',
  type: 'fermenter',
  status: 'active',
  capacity: 150,
  capacity_unit: 'gal',
}

describe('useVesselActions', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('saveVessel', () => {
    describe('successful save', () => {
      it('returns the updated vessel on success', async () => {
        mockUpdateVessel.mockResolvedValue(sampleVessel)
        const { dialogRef } = createMockDialogRef()

        const { saveVessel } = useVesselActions()
        const result = await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(result).toEqual(sampleVessel)
      })

      it('calls updateVessel with correct uuid and data', async () => {
        mockUpdateVessel.mockResolvedValue(sampleVessel)
        const { dialogRef } = createMockDialogRef()

        const { saveVessel } = useVesselActions()
        await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(mockUpdateVessel).toHaveBeenCalledWith('vessel-uuid-1', sampleRequest)
      })

      it('shows success snackbar notification', async () => {
        mockUpdateVessel.mockResolvedValue(sampleVessel)
        const { dialogRef } = createMockDialogRef()

        const { saveVessel } = useVesselActions()
        await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(mockShowNotice).toHaveBeenCalledWith('Vessel updated successfully')
      })

      it('sets saving to true then false', async () => {
        mockUpdateVessel.mockResolvedValue(sampleVessel)
        const { dialogRef, setSaving } = createMockDialogRef()

        const { saveVessel } = useVesselActions()
        await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(setSaving).toHaveBeenCalledWith(true)
        expect(setSaving).toHaveBeenCalledWith(false)
        // true is called first, false is called last
        expect(setSaving.mock.calls[0][0]).toBe(true)
        expect(setSaving.mock.calls[setSaving.mock.calls.length - 1][0]).toBe(false)
      })

      it('clears error before saving', async () => {
        mockUpdateVessel.mockResolvedValue(sampleVessel)
        const { dialogRef, clearError } = createMockDialogRef()

        const { saveVessel } = useVesselActions()
        await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(clearError).toHaveBeenCalled()
      })
    })

    describe('conflict detection (409)', () => {
      it('sets conflict error message on 409 response', async () => {
        mockUpdateVessel.mockRejectedValue(new Error('409 Conflict'))
        const { dialogRef, setError } = createMockDialogRef()

        const { saveVessel } = useVesselActions()
        const result = await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(result).toBeNull()
        expect(setError).toHaveBeenCalledWith(
          'Cannot change status: vessel has an active occupancy',
        )
      })

      it('returns null on 409 conflict', async () => {
        mockUpdateVessel.mockRejectedValue(new Error('Request failed with 409'))
        const { dialogRef } = createMockDialogRef()

        const { saveVessel } = useVesselActions()
        const result = await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(result).toBeNull()
      })

      it('does not show success snackbar on 409', async () => {
        mockUpdateVessel.mockRejectedValue(new Error('409 Conflict'))
        const { dialogRef } = createMockDialogRef()

        const { saveVessel } = useVesselActions()
        await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(mockShowNotice).not.toHaveBeenCalled()
      })
    })

    describe('general error handling', () => {
      it('sets error message from Error instance', async () => {
        mockUpdateVessel.mockRejectedValue(new Error('Network error'))
        const { dialogRef, setError } = createMockDialogRef()

        const { saveVessel } = useVesselActions()
        const result = await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(result).toBeNull()
        expect(setError).toHaveBeenCalledWith('Network error')
      })

      it('sets fallback error message for non-Error throws', async () => {
        mockUpdateVessel.mockRejectedValue('string error')
        const { dialogRef, setError } = createMockDialogRef()

        const { saveVessel } = useVesselActions()
        const result = await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(result).toBeNull()
        expect(setError).toHaveBeenCalledWith('Failed to update vessel')
      })

      it('sets saving to false after error', async () => {
        mockUpdateVessel.mockRejectedValue(new Error('Server error'))
        const { dialogRef, setSaving } = createMockDialogRef()

        const { saveVessel } = useVesselActions()
        await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        // Last call should be false (from finally block)
        expect(setSaving).toHaveBeenLastCalledWith(false)
      })

      it('does not show success snackbar on error', async () => {
        mockUpdateVessel.mockRejectedValue(new Error('Server error'))
        const { dialogRef } = createMockDialogRef()

        const { saveVessel } = useVesselActions()
        await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(mockShowNotice).not.toHaveBeenCalled()
      })
    })

    describe('null dialog ref', () => {
      it('handles null dialog ref gracefully on success', async () => {
        mockUpdateVessel.mockResolvedValue(sampleVessel)
        const dialogRef = ref(null)

        const { saveVessel } = useVesselActions()
        const result = await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(result).toEqual(sampleVessel)
        expect(mockShowNotice).toHaveBeenCalledWith('Vessel updated successfully')
      })

      it('handles null dialog ref gracefully on error', async () => {
        mockUpdateVessel.mockRejectedValue(new Error('Server error'))
        const dialogRef = ref(null)

        const { saveVessel } = useVesselActions()
        const result = await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(result).toBeNull()
      })

      it('handles null dialog ref gracefully on 409', async () => {
        mockUpdateVessel.mockRejectedValue(new Error('409 Conflict'))
        const dialogRef = ref(null)

        const { saveVessel } = useVesselActions()
        const result = await saveVessel('vessel-uuid-1', sampleRequest, dialogRef)

        expect(result).toBeNull()
      })
    })
  })
})
