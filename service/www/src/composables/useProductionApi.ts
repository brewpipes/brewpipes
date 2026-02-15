import type {
  Addition,
  Batch,
  BatchProcessPhase,
  BatchRelation,
  BatchSummary,
  BatchVolume,
  BrewSession,
  CreateAdditionRequest,
  CreateBatchProcessPhaseRequest,
  CreateBatchRelationRequest,
  CreateBatchVolumeRequest,
  CreateBrewSessionRequest,
  CreateMeasurementRequest,
  CreateOccupancyRequest,
  CreateRecipeIngredientRequest,
  CreateRecipeRequest,
  CreateStyleRequest,
  CreateTransferRequest,
  CreateVolumeRelationRequest,
  CreateVolumeRequest,
  Measurement,
  Occupancy,
  OccupancyStatus,
  Recipe,
  RecipeIngredient,
  Style,
  Transfer,
  TransferRecordResponse,
  UpdateBatchRequest,
  UpdateBrewSessionRequest,
  UpdateRecipeIngredientRequest,
  UpdateRecipeRequest,
  UpdateVesselRequest,
  Vessel,
  Volume,
  VolumeRelation,
} from '@/types'
import { useApiClient } from '@/composables/useApiClient'
import { formatDateTime } from '@/composables/useFormatters'
import { normalizeDateTime, normalizeText, toNumber } from '@/utils/normalize'

const productionApiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'

export function useProductionApi () {
  const { request } = useApiClient(productionApiBase)

  // Styles API
  const getStyles = () => request<Style[]>('/styles')
  const getStyle = (uuid: string) => request<Style>(`/styles/${uuid}`)
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
  const getVessel = (uuid: string) => request<Vessel>(`/vessels/${uuid}`)
  const createVessel = (data: Record<string, unknown>) =>
    request<Vessel>('/vessels', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updateVessel = (uuid: string, data: UpdateVesselRequest) =>
    request<Vessel>(`/vessels/${uuid}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })

  // Volumes API
  const getVolumes = () => request<Volume[]>('/volumes')
  const getVolume = (uuid: string) => request<Volume>(`/volumes/${uuid}`)
  const createVolume = (data: CreateVolumeRequest) =>
    request<Volume>('/volumes', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Brew Sessions API
  const getBrewSessions = (batchUuid: string) =>
    request<BrewSession[]>(`/brew-sessions?batch_uuid=${batchUuid}`)
  const getBrewSession = (uuid: string) =>
    request<BrewSession>(`/brew-sessions/${uuid}`)
  const createBrewSession = (data: CreateBrewSessionRequest) =>
    request<BrewSession>('/brew-sessions', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updateBrewSession = (uuid: string, data: UpdateBrewSessionRequest) =>
    request<BrewSession>(`/brew-sessions/${uuid}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })

  // Additions API
  const getAdditionsByVolume = (volumeUuid: string) =>
    request<Addition[]>(`/additions?volume_uuid=${volumeUuid}`)
  const getAdditionsByBatch = (batchUuid: string) =>
    request<Addition[]>(`/additions?batch_uuid=${batchUuid}`)
  const getAdditionsByOccupancy = (occupancyUuid: string) =>
    request<Addition[]>(`/additions?occupancy_uuid=${occupancyUuid}`)
  const createAddition = (data: CreateAdditionRequest) =>
    request<Addition>('/additions', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Measurements API
  const getMeasurementsByVolume = (volumeUuid: string) =>
    request<Measurement[]>(`/measurements?volume_uuid=${volumeUuid}`)
  const getMeasurementsByBatch = (batchUuid: string) =>
    request<Measurement[]>(`/measurements?batch_uuid=${batchUuid}`)
  const getMeasurementsByOccupancy = (occupancyUuid: string) =>
    request<Measurement[]>(`/measurements?occupancy_uuid=${occupancyUuid}`)
  const createMeasurement = (data: CreateMeasurementRequest) =>
    request<Measurement>('/measurements', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Batches API
  const getBatches = () => request<Batch[]>('/batches')
  const getBatch = (uuid: string) =>
    request<Batch>(`/batches/${uuid}`)
  const createBatch = (data: Record<string, unknown>) =>
    request<Batch>('/batches', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updateBatch = (uuid: string, data: UpdateBatchRequest) =>
    request<Batch>(`/batches/${uuid}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  const deleteBatch = (uuid: string) =>
    request<void>(`/batches/${uuid}`, {
      method: 'DELETE',
    })

  // Batch Summary API
  const getBatchSummary = (uuid: string) =>
    request<BatchSummary>(`/batches/${uuid}/summary`)

  // Occupancy API
  const getActiveOccupancies = () =>
    request<Occupancy[]>('/occupancies?active=true')
  const getOccupancy = (uuid: string) =>
    request<Occupancy>(`/occupancies/${uuid}`)
  const createOccupancy = (data: CreateOccupancyRequest) =>
    request<Occupancy>('/occupancies', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  const updateOccupancyStatus = (uuid: string, status: OccupancyStatus) =>
    request<Occupancy>(`/occupancies/${uuid}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    })
  const closeOccupancy = (uuid: string, outAt?: string) =>
    request<Occupancy>(`/occupancies/${uuid}/close`, {
      method: 'PATCH',
      body: JSON.stringify(outAt ? { out_at: outAt } : {}),
    })
  const getActiveOccupancyByVessel = (vesselUuid: string) =>
    request<Occupancy>(`/occupancies/active?active_vessel_uuid=${vesselUuid}`)

  // Transfers API
  const getTransfers = (batchUuid: string) =>
    request<Transfer[]>(`/transfers?batch_uuid=${batchUuid}`)
  const getTransfer = (uuid: string) =>
    request<Transfer>(`/transfers/${uuid}`)
  const createTransfer = (data: CreateTransferRequest) =>
    request<TransferRecordResponse>('/transfers', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Batch Process Phases API
  const getBatchProcessPhases = (batchUuid: string) =>
    request<BatchProcessPhase[]>(`/batch-process-phases?batch_uuid=${batchUuid}`)
  const createBatchProcessPhase = (data: CreateBatchProcessPhaseRequest) =>
    request<BatchProcessPhase>('/batch-process-phases', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Batch Volumes API
  const getBatchVolumes = (params?: string) =>
    request<BatchVolume[]>(`/batch-volumes${params ? `?${params}` : ''}`)
  const createBatchVolume = (data: CreateBatchVolumeRequest) =>
    request<BatchVolume>('/batch-volumes', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Volume Relations API
  const getVolumeRelations = (params?: string) =>
    request<VolumeRelation[]>(`/volume-relations${params ? `?${params}` : ''}`)
  const createVolumeRelation = (data: CreateVolumeRelationRequest) =>
    request<VolumeRelation>('/volume-relations', {
      method: 'POST',
      body: JSON.stringify(data),
    })

  // Batch Relations API
  const getBatchRelations = (params?: string) =>
    request<BatchRelation[]>(`/batch-relations${params ? `?${params}` : ''}`)
  const createBatchRelation = (data: CreateBatchRelationRequest) =>
    request<BatchRelation>('/batch-relations', {
      method: 'POST',
      body: JSON.stringify(data),
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
    createVessel,
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
    getAdditionsByBatch,
    getAdditionsByOccupancy,
    createAddition,
    // Measurements
    getMeasurementsByVolume,
    getMeasurementsByBatch,
    getMeasurementsByOccupancy,
    createMeasurement,
    // Batches
    getBatches,
    getBatch,
    createBatch,
    updateBatch,
    deleteBatch,
    // Batch Summary
    getBatchSummary,
    // Occupancy
    getActiveOccupancies,
    getOccupancy,
    createOccupancy,
    updateOccupancyStatus,
    closeOccupancy,
    getActiveOccupancyByVessel,
    // Transfers
    getTransfers,
    getTransfer,
    createTransfer,
    // Batch Process Phases
    getBatchProcessPhases,
    createBatchProcessPhase,
    // Batch Volumes
    getBatchVolumes,
    createBatchVolume,
    // Volume Relations
    getVolumeRelations,
    createVolumeRelation,
    // Batch Relations
    getBatchRelations,
    createBatchRelation,
  }
}
