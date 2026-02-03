import type {
  Addition,
  BatchSummary,
  BrewSession,
  CreateAdditionRequest,
  CreateBrewSessionRequest,
  CreateMeasurementRequest,
  CreateRecipeRequest,
  CreateStyleRequest,
  CreateVolumeRequest,
  Measurement,
  Occupancy,
  OccupancyStatus,
  Recipe,
  Style,
  UpdateBrewSessionRequest,
  UpdateRecipeRequest,
  Vessel,
  Volume,
} from '@/types'
import { useApiClient } from '@/composables/useApiClient'

// Re-export types for backward compatibility
export type {
  Addition,
  AdditionType,
  Batch,
  BatchSummary,
  BatchSummaryBrewSession,
  BrewSession,
  CreateAdditionRequest,
  CreateBrewSessionRequest,
  CreateMeasurementRequest,
  CreateRecipeRequest,
  CreateStyleRequest,
  CreateVolumeRequest,
  Measurement,
  Occupancy,
  OccupancyStatus,
  Recipe,
  Style,
  UpdateBrewSessionRequest,
  UpdateRecipeRequest,
  Vessel,
  VesselStatus,
  Volume,
  VolumeUnit,
} from '@/types'

export { OCCUPANCY_STATUS_VALUES } from '@/types'

const productionApiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'

export function useProductionApi () {
  const { request } = useApiClient(productionApiBase)

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

  // Styles API
  const getStyles = () => request<Style[]>('/styles')
  const getStyle = (id: number) => request<Style>(`/styles/${id}`)
  const createStyle = (data: CreateStyleRequest) =>
    request<Style>('/styles', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Recipes API
  const getRecipes = () => request<Recipe[]>('/recipes')
  const getRecipe = (id: number) => request<Recipe>(`/recipes/${id}`)
  const createRecipe = (data: CreateRecipeRequest) =>
    request<Recipe>('/recipes', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updateRecipe = (id: number, data: UpdateRecipeRequest) =>
    request<Recipe>(`/recipes/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })

  // Vessels API
  const getVessels = () => request<Vessel[]>('/vessels')
  const getVessel = (id: number) => request<Vessel>(`/vessels/${id}`)
  const getVesselByUUID = (uuid: string) => request<Vessel>(`/vessels/uuid/${uuid}`)

  // Volumes API
  const getVolumes = () => request<Volume[]>('/volumes')
  const getVolume = (id: number) => request<Volume>(`/volumes/${id}`)
  const createVolume = (data: CreateVolumeRequest) =>
    request<Volume>('/volumes', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Brew Sessions API
  const getBrewSessions = (batchId: number) =>
    request<BrewSession[]>(`/brew-sessions?batch_id=${batchId}`)
  const getBrewSession = (id: number) =>
    request<BrewSession>(`/brew-sessions/${id}`)
  const createBrewSession = (data: CreateBrewSessionRequest) =>
    request<BrewSession>('/brew-sessions', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updateBrewSession = (id: number, data: UpdateBrewSessionRequest) =>
    request<BrewSession>(`/brew-sessions/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })

  // Additions API (volume-targeted)
  const getAdditionsByVolume = (volumeId: number) =>
    request<Addition[]>(`/additions?volume_id=${volumeId}`)
  const createAddition = (data: CreateAdditionRequest) =>
    request<Addition>('/additions', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Measurements API (volume-targeted)
  const getMeasurementsByVolume = (volumeId: number) =>
    request<Measurement[]>(`/measurements?volume_id=${volumeId}`)
  const createMeasurement = (data: CreateMeasurementRequest) =>
    request<Measurement>('/measurements', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Batch Summary API
  const getBatchSummary = (id: number) =>
    request<BatchSummary>(`/batches/${id}/summary`)

  // Occupancy API
  const getActiveOccupancies = () =>
    request<Occupancy[]>('/occupancies?active=true')
  const getOccupancy = (id: number) =>
    request<Occupancy>(`/occupancies/${id}`)
  const updateOccupancyStatus = (id: number, status: OccupancyStatus) =>
    request<Occupancy>(`/occupancies/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    })

  return {
    apiBase: productionApiBase,
    request,
    normalizeText,
    normalizeDateTime,
    toNumber,
    formatDateTime,
    // Styles
    getStyles,
    getStyle,
    createStyle,
    // Recipes
    getRecipes,
    getRecipe,
    createRecipe,
    updateRecipe,
    // Vessels
    getVessels,
    getVessel,
    getVesselByUUID,
    // Volumes
    getVolumes,
    getVolume,
    createVolume,
    // Brew Sessions
    getBrewSessions,
    getBrewSession,
    createBrewSession,
    updateBrewSession,
    // Additions
    getAdditionsByVolume,
    createAddition,
    // Measurements
    getMeasurementsByVolume,
    createMeasurement,
    // Batch Summary
    getBatchSummary,
    // Occupancy
    getActiveOccupancies,
    getOccupancy,
    updateOccupancyStatus,
  }
}
