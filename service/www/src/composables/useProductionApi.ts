import type {
  Addition,
  Batch,
  BatchSummary,
  BrewSession,
  CreateAdditionRequest,
  CreateBrewSessionRequest,
  CreateMeasurementRequest,
  CreateRecipeIngredientRequest,
  CreateRecipeRequest,
  CreateStyleRequest,
  CreateVolumeRequest,
  Measurement,
  Occupancy,
  OccupancyStatus,
  Recipe,
  RecipeIngredient,
  Style,
  UpdateBatchRequest,
  UpdateBrewSessionRequest,
  UpdateRecipeIngredientRequest,
  UpdateRecipeRequest,
  UpdateVesselRequest,
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
  CreateRecipeIngredientRequest,
  CreateRecipeRequest,
  CreateStyleRequest,
  CreateVolumeRequest,
  Measurement,
  Occupancy,
  OccupancyStatus,
  Recipe,
  RecipeIngredient,
  RecipeIngredientType,
  RecipeUseStage,
  RecipeUseType,
  Style,
  UpdateBatchRequest,
  UpdateBrewSessionRequest,
  UpdateRecipeIngredientRequest,
  UpdateRecipeRequest,
  UpdateVesselRequest,
  Vessel,
  VesselStatus,
  VesselType,
  Volume,
  VolumeUnit,
} from '@/types'

export {
  OCCUPANCY_STATUS_VALUES,
  RECIPE_INGREDIENT_TYPE_VALUES,
  RECIPE_USE_STAGE_VALUES,
  RECIPE_USE_TYPE_VALUES,
  VESSEL_STATUS_VALUES,
  VESSEL_TYPE_VALUES,
} from '@/types'

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
  const getRecipe = (uuid: string) => request<Recipe>(`/recipes/${uuid}`)
  const createRecipe = (data: CreateRecipeRequest) =>
    request<Recipe>('/recipes', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updateRecipe = (uuid: string, data: UpdateRecipeRequest) =>
    request<Recipe>(`/recipes/${uuid}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  const deleteRecipe = (uuid: string) =>
    request<void>(`/recipes/${uuid}`, {
      method: 'DELETE',
    })

  // Recipe Ingredients API
  const getRecipeIngredients = (recipeUuid: string) =>
    request<RecipeIngredient[]>(`/recipes/${recipeUuid}/ingredients`)
  const getRecipeIngredient = (recipeUuid: string, ingredientUuid: string) =>
    request<RecipeIngredient>(`/recipes/${recipeUuid}/ingredients/${ingredientUuid}`)
  const createRecipeIngredient = (recipeUuid: string, data: CreateRecipeIngredientRequest) =>
    request<RecipeIngredient>(`/recipes/${recipeUuid}/ingredients`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updateRecipeIngredient = (recipeUuid: string, ingredientUuid: string, data: UpdateRecipeIngredientRequest) =>
    request<RecipeIngredient>(`/recipes/${recipeUuid}/ingredients/${ingredientUuid}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  const deleteRecipeIngredient = (recipeUuid: string, ingredientUuid: string) =>
    request<void>(`/recipes/${recipeUuid}/ingredients/${ingredientUuid}`, {
      method: 'DELETE',
    })

  // Vessels API
  const getVessels = () => request<Vessel[]>('/vessels')
  const getVessel = (id: number) => request<Vessel>(`/vessels/${id}`)
  const getVesselByUUID = (uuid: string) => request<Vessel>(`/vessels/uuid/${uuid}`)
  const updateVessel = (id: number, data: UpdateVesselRequest) =>
    request<Vessel>(`/vessels/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })

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

  // Batches API
  const getBatch = (id: number) =>
    request<Batch>(`/batches/${id}`)
  const updateBatch = (id: number, data: UpdateBatchRequest) =>
    request<Batch>(`/batches/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  const deleteBatch = (id: number) =>
    request<void>(`/batches/${id}`, {
      method: 'DELETE',
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
    deleteRecipe,
    // Recipe Ingredients
    getRecipeIngredients,
    getRecipeIngredient,
    createRecipeIngredient,
    updateRecipeIngredient,
    deleteRecipeIngredient,
    // Vessels
    getVessels,
    getVessel,
    getVesselByUUID,
    updateVessel,
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
    // Batches
    getBatch,
    updateBatch,
    deleteBatch,
    // Batch Summary
    getBatchSummary,
    // Occupancy
    getActiveOccupancies,
    getOccupancy,
    updateOccupancyStatus,
  }
}
