import type { VolumeUnit } from '@/types'
import { computed, ref } from 'vue'
import type { Ref } from 'vue'
import { convertVolume } from '@/composables/useUnitConversion'

/**
 * Batch size units supported by the scaling calculator.
 * Maps user-facing labels (bbl, gal, L, hL) to the canonical VolumeUnit values
 * used by the unit conversion system.
 */
const BATCH_SIZE_UNIT_MAP: Record<string, VolumeUnit> = {
  bbl: 'bbl',
  gal: 'usgal',
  l: 'l',
  hl: 'hl',
  // Also accept canonical VolumeUnit values directly
  usgal: 'usgal',
}

/**
 * Resolve a batch size unit string (which may be a user-facing label like "gal")
 * to the canonical VolumeUnit used by convertVolume.
 */
function resolveVolumeUnit (unit: string | null): VolumeUnit | null {
  if (!unit) return null
  return BATCH_SIZE_UNIT_MAP[unit.toLowerCase()] ?? null
}

/**
 * Composable for client-side recipe scaling.
 *
 * Computes a scale factor from the recipe's batch size to a user-specified
 * target batch size, with unit conversion support. Provides a `scaleAmount`
 * function that applies the scale factor along with per-ingredient scaling
 * factors (e.g., yeast that doesn't scale linearly).
 *
 * All scaling is ephemeral — it does not modify stored recipe data.
 */
export function useRecipeScaling (
  recipeBatchSize: Ref<number | null>,
  recipeBatchSizeUnit: Ref<string | null>,
) {
  /** Target batch size entered by the user. null = no scaling active. */
  const targetBatchSize = ref<number | null>(null)

  /** Unit for the target batch size. Defaults to recipe's unit. */
  const targetBatchSizeUnit = ref<string>('bbl')

  /**
   * The ratio of target batch size to recipe batch size, with unit conversion.
   * Returns 1 when scaling is not active or inputs are invalid (null, zero, negative).
   */
  const scaleFactor = computed(() => {
    if (!recipeBatchSize.value || recipeBatchSize.value <= 0) return 1
    if (!targetBatchSize.value || targetBatchSize.value <= 0) return 1

    const recipeUnit = resolveVolumeUnit(recipeBatchSizeUnit.value ?? 'bbl')
    const targetUnit = resolveVolumeUnit(targetBatchSizeUnit.value)

    if (!recipeUnit || !targetUnit) return 1

    // Convert target batch size to recipe's unit for comparison
    const targetInRecipeUnits = convertVolume(targetBatchSize.value, targetUnit, recipeUnit)
    if (targetInRecipeUnits === null || targetInRecipeUnits <= 0) return 1

    return targetInRecipeUnits / recipeBatchSize.value
  })

  /**
   * Whether scaling is currently active (scale factor differs meaningfully from 1).
   * Uses a tolerance of 0.01% to account for floating-point rounding in unit conversion.
   */
  const isScaling = computed(() => {
    if (targetBatchSize.value === null || recipeBatchSize.value === null) return false
    return Math.abs(scaleFactor.value - 1) > 1e-4
  })

  /** Human-readable scale factor display (e.g., "150%" or "× 1.5"). */
  const scaleFactorDisplay = computed(() => {
    const factor = scaleFactor.value
    const percent = Math.round(factor * 100)
    return `${percent}%`
  })

  /**
   * Scale an ingredient amount by the current scale factor and the
   * ingredient's own scaling factor.
   *
   * @param amount - The original ingredient amount
   * @param ingredientScalingFactor - Per-ingredient scaling behavior (default 1.0).
   *   A value of 0.5 means the ingredient scales at half the rate.
   *   A value of 0 means the ingredient amount is fixed regardless of batch size.
   */
  function scaleAmount (amount: number, ingredientScalingFactor: number = 1.0): number {
    if (ingredientScalingFactor === 0) return amount
    // Blend between unscaled and fully-scaled based on ingredientScalingFactor
    const fullyScaled = amount * scaleFactor.value
    return amount + (fullyScaled - amount) * ingredientScalingFactor
  }

  /** Reset scaling to inactive state, restoring the recipe's native unit. */
  function resetScaling () {
    targetBatchSize.value = null
    targetBatchSizeUnit.value = recipeBatchSizeUnit.value ?? 'bbl'
  }

  /** Set the target batch size and optionally the unit. */
  function setTargetBatchSize (size: number, unit?: string) {
    targetBatchSize.value = size
    if (unit) targetBatchSizeUnit.value = unit
  }

  return {
    targetBatchSize,
    targetBatchSizeUnit,
    isScaling,
    scaleFactor,
    scaleFactorDisplay,
    scaleAmount,
    resetScaling,
    setTargetBatchSize,
  }
}
