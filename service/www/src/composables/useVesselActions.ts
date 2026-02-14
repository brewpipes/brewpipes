import type { UpdateVesselRequest, Vessel } from '@/types'
import type VesselEditDialog from '@/components/vessel/VesselEditDialog.vue'
import { useProductionApi } from '@/composables/useProductionApi'
import { useSnackbar } from '@/composables/useSnackbar'
import { type Ref } from 'vue'

/**
 * Shared vessel save handler with 409 conflict detection.
 * Used by vessel pages that include the VesselEditDialog.
 */
export function useVesselActions () {
  const { updateVessel } = useProductionApi()
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

  return { saveVessel }
}
