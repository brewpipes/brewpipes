import type { Ref } from 'vue'
import type VesselEditDialog from '@/components/vessel/VesselEditDialog.vue'
import type VesselRetireDialog from '@/components/vessel/VesselRetireDialog.vue'
import type { UpdateVesselRequest, Vessel } from '@/types'
import { useProductionApi } from '@/composables/useProductionApi'
import { useSnackbar } from '@/composables/useSnackbar'

/**
 * Shared vessel action handlers with 409 conflict detection.
 * Used by vessel pages that include the VesselEditDialog and VesselRetireDialog.
 */
export function useVesselActions () {
  const { getVessel, updateVessel } = useProductionApi()
  const { showNotice } = useSnackbar()

  async function saveVessel (
    vesselUuid: string,
    data: UpdateVesselRequest,
    dialogRef: Ref<InstanceType<typeof VesselEditDialog> | null>,
  ): Promise<Vessel | null> {
    dialogRef.value?.setSaving(true)
    dialogRef.value?.clearError()

    try {
      const updated = await updateVessel(vesselUuid, data)
      showNotice('Vessel updated successfully')
      return updated
    } catch (error_) {
      // Check for 409 Conflict (vessel has active occupancy)
      if (error_ instanceof Error && error_.message.includes('409')) {
        dialogRef.value?.setError('Cannot change status: vessel has an active occupancy')
      } else {
        const message = error_ instanceof Error ? error_.message : 'Failed to update vessel'
        dialogRef.value?.setError(message)
      }
      return null
    } finally {
      dialogRef.value?.setSaving(false)
    }
  }

  /**
   * Retire a vessel by fetching its current data and PATCHing with status "retired".
   * Uses the VesselRetireDialog ref for loading/error state.
   */
  async function retireVessel (
    vesselUuid: string,
    dialogRef: Ref<InstanceType<typeof VesselRetireDialog> | null>,
  ): Promise<Vessel | null> {
    dialogRef.value?.setRetiring(true)

    try {
      // Fetch current vessel data to build a complete update request
      const current = await getVessel(vesselUuid)
      const data: UpdateVesselRequest = {
        name: current.name,
        type: current.type,
        capacity: current.capacity,
        capacity_unit: current.capacity_unit,
        make: current.make,
        model: current.model,
        status: 'retired',
      }

      const updated = await updateVessel(vesselUuid, data)
      showNotice('Vessel retired successfully')
      return updated
    } catch (error_) {
      if (error_ instanceof Error && error_.message.includes('409')) {
        dialogRef.value?.setError('Cannot retire vessel: it has an active occupancy. Remove the occupancy first.')
      } else {
        const message = error_ instanceof Error ? error_.message : 'Failed to retire vessel'
        dialogRef.value?.setError(message)
      }
      return null
    } finally {
      dialogRef.value?.setRetiring(false)
    }
  }

  return { retireVessel, saveVessel }
}
