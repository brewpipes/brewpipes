import { useApiClient } from '@/composables/useApiClient'

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

// Types
export type Style = {
  id: number
  uuid: string
  name: string
  created_at: string
  updated_at: string
}

export type CreateStyleRequest = {
  name: string
}

export type Recipe = {
  id: number
  uuid: string
  name: string
  style_id: number | null
  style_name: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

export type CreateRecipeRequest = {
  name: string
  style_id?: number | null
  style_name?: string | null
  notes?: string | null
}

export type UpdateRecipeRequest = {
  name: string
  style_id?: number | null
  style_name?: string | null
  notes?: string | null
}

// Vessel types
export type VolumeUnit = 'ml' | 'usfloz' | 'ukfloz' | 'bbl'

export type Vessel = {
  id: number
  uuid: string
  type: string
  name: string
  capacity: number
  capacity_unit: VolumeUnit
  make: string | null
  model: string | null
  status: 'active' | 'inactive' | 'retired'
  created_at: string
  updated_at: string
  deleted_at: string | null
}

// Volume types
export type Volume = {
  id: number
  uuid: string
  name: string | null
  description: string | null
  amount: number
  amount_unit: VolumeUnit
  created_at: string
  updated_at: string
  deleted_at: string | null
}

export type CreateVolumeRequest = {
  name?: string | null
  description?: string | null
  amount: number
  amount_unit: VolumeUnit
}

// Brew Session types
export type BrewSession = {
  id: number
  uuid: string
  batch_id: number | null
  wort_volume_id: number | null
  mash_vessel_id: number | null
  boil_vessel_id: number | null
  brewed_at: string
  notes: string | null
  created_at: string
  updated_at: string
  deleted_at: string | null
}

export type CreateBrewSessionRequest = {
  batch_id?: number | null
  wort_volume_id?: number | null
  mash_vessel_id?: number | null
  boil_vessel_id?: number | null
  brewed_at: string
  notes?: string | null
}

export type UpdateBrewSessionRequest = {
  batch_id?: number | null
  wort_volume_id?: number | null
  mash_vessel_id?: number | null
  boil_vessel_id?: number | null
  brewed_at: string
  notes?: string | null
}

// Addition types
export type AdditionType
  = | 'malt'
    | 'hop'
    | 'yeast'
    | 'adjunct'
    | 'water_chem'
    | 'gas'
    | 'other'

export type Addition = {
  id: number
  uuid: string
  batch_id: number | null
  occupancy_id: number | null
  volume_id: number | null
  addition_type: AdditionType
  stage: string | null
  inventory_lot_uuid: string | null
  amount: number
  amount_unit: VolumeUnit
  added_at: string
  notes: string | null
  created_at: string
  updated_at: string
  deleted_at: string | null
}

export type CreateAdditionRequest = {
  batch_id?: number | null
  occupancy_id?: number | null
  volume_id?: number | null
  addition_type: AdditionType
  stage?: string | null
  inventory_lot_uuid?: string | null
  amount: number
  amount_unit: VolumeUnit
  added_at?: string | null
  notes?: string | null
}

// Measurement types
export type Measurement = {
  id: number
  uuid: string
  batch_id: number | null
  occupancy_id: number | null
  volume_id: number | null
  kind: string
  value: number
  unit: string | null
  observed_at: string
  notes: string | null
  created_at: string
  updated_at: string
  deleted_at: string | null
}

export type CreateMeasurementRequest = {
  batch_id?: number | null
  occupancy_id?: number | null
  volume_id?: number | null
  kind: string
  value: number
  unit?: string | null
  observed_at?: string | null
  notes?: string | null
}

// Batch Summary types
export type BatchSummaryBrewSession = {
  id: number
  brewed_at: string
  notes: string | null
}

export type BatchSummary = {
  id: number
  uuid: string
  short_name: string
  brew_date: string | null
  notes: string | null
  recipe_name: string | null
  style_name: string | null
  brew_sessions: BatchSummaryBrewSession[]
  current_phase: string | null
  current_vessel: string | null
  current_occupancy_status: string | null
  current_occupancy_id: number | null
  original_gravity: number | null
  final_gravity: number | null
  abv: number | null
  abv_calculated: boolean
  ibu: number | null
  days_in_fermenter: number | null
  days_in_brite: number | null
  days_grain_to_glass: number | null
  starting_volume_bbl: number | null
  current_volume_bbl: number | null
  total_loss_bbl: number | null
  loss_percentage: number | null
  created_at: string
  updated_at: string
}

// Occupancy types
export type OccupancyStatus
  = | 'fermenting'
    | 'conditioning'
    | 'cold_crashing'
    | 'dry_hopping'
    | 'carbonating'
    | 'holding'
    | 'packaging'

export const OCCUPANCY_STATUS_VALUES: OccupancyStatus[] = [
  'fermenting',
  'conditioning',
  'cold_crashing',
  'dry_hopping',
  'carbonating',
  'holding',
  'packaging',
]

export type Occupancy = {
  id: number
  uuid: string
  vessel_id: number
  volume_id: number
  batch_id: number | null
  status: OccupancyStatus | null
  in_at: string
  out_at: string | null
  created_at: string
  updated_at: string
}
