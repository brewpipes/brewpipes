import { describe, expect, it } from 'vitest'
import { ref } from 'vue'
import { useRecipeScaling } from '../useRecipeScaling'

describe('useRecipeScaling', () => {
  describe('initial state', () => {
    it('starts with scaling inactive', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { isScaling, scaleFactor, targetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      expect(isScaling.value).toBe(false)
      expect(scaleFactor.value).toBe(1)
      expect(targetBatchSize.value).toBeNull()
    })

    it('returns scale factor 1 when recipe has no batch size', () => {
      const recipeBatchSize = ref<number | null>(null)
      const recipeBatchSizeUnit = ref<string | null>(null)
      const { scaleFactor } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      expect(scaleFactor.value).toBe(1)
    })
  })

  describe('setTargetBatchSize', () => {
    it('activates scaling when target differs from recipe', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { isScaling, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(15)
      expect(isScaling.value).toBe(true)
    })

    it('does not activate scaling when target equals recipe', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { isScaling, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(10)
      expect(isScaling.value).toBe(false)
    })

    it('does not activate scaling when cross-unit target is equivalent volume', () => {
      const recipeBatchSize = ref<number | null>(1)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { isScaling, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      // 1 bbl ≈ 31 gal — should not be considered scaling
      setTargetBatchSize(31, 'gal')
      expect(isScaling.value).toBe(false)
    })

    it('sets the unit when provided', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { targetBatchSizeUnit, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(100, 'gal')
      expect(targetBatchSizeUnit.value).toBe('gal')
    })
  })

  describe('scaleFactor', () => {
    it('computes correct factor for same-unit scaling', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleFactor, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(15)
      expect(scaleFactor.value).toBe(1.5)
    })

    it('computes correct factor for scaling down', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleFactor, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(5)
      expect(scaleFactor.value).toBe(0.5)
    })

    it('handles cross-unit scaling (bbl to gal)', () => {
      const recipeBatchSize = ref<number | null>(1)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleFactor, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      // 1 bbl = 31 gal, so 31 gal target should give factor ~1.0
      setTargetBatchSize(31, 'gal')
      expect(scaleFactor.value).toBeCloseTo(1.0, 2)
    })

    it('handles cross-unit scaling (gal to bbl)', () => {
      const recipeBatchSize = ref<number | null>(31)
      const recipeBatchSizeUnit = ref<string | null>('gal')
      const { scaleFactor, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      // 31 gal recipe, target 2 bbl = 62 gal → factor ~2.0
      setTargetBatchSize(2, 'bbl')
      expect(scaleFactor.value).toBeCloseTo(2.0, 2)
    })

    it('returns 1 when target batch size is null', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleFactor } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      expect(scaleFactor.value).toBe(1)
    })

    it('returns 1 when recipe batch size is zero', () => {
      const recipeBatchSize = ref<number | null>(0)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleFactor, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(10)
      expect(scaleFactor.value).toBe(1)
    })

    it('returns 1 when recipe batch size is negative', () => {
      const recipeBatchSize = ref<number | null>(-5)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleFactor, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(10)
      expect(scaleFactor.value).toBe(1)
    })

    it('returns 1 when target batch size is negative', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleFactor, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(-5)
      expect(scaleFactor.value).toBe(1)
    })
  })

  describe('scaleAmount', () => {
    it('scales amount by the scale factor with default scaling factor', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleAmount, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(15)
      expect(scaleAmount(100)).toBe(150)
    })

    it('applies ingredient scaling factor', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleAmount, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(20) // 2x scale
      // With ingredient scaling factor 0.5:
      // fullyScaled = 100 * 2 = 200
      // result = 100 + (200 - 100) * 0.5 = 150
      expect(scaleAmount(100, 0.5)).toBe(150)
    })

    it('returns original amount when ingredient scaling factor is 0', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleAmount, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(20) // 2x scale
      expect(scaleAmount(100, 0)).toBe(100)
    })

    it('returns unscaled amount when no target is set', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleAmount } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      expect(scaleAmount(100)).toBe(100)
    })
  })

  describe('scaleFactorDisplay', () => {
    it('shows percentage', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleFactorDisplay, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(15)
      expect(scaleFactorDisplay.value).toBe('150%')
    })

    it('shows 100% when not scaling', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleFactorDisplay } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      expect(scaleFactorDisplay.value).toBe('100%')
    })
  })

  describe('resetScaling', () => {
    it('clears target and deactivates scaling', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { isScaling, targetBatchSize, setTargetBatchSize, resetScaling } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(15)
      expect(isScaling.value).toBe(true)

      resetScaling()
      expect(isScaling.value).toBe(false)
      expect(targetBatchSize.value).toBeNull()
    })
  })

  describe('reactivity', () => {
    it('recomputes when recipe batch size changes', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleFactor, setTargetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      setTargetBatchSize(20)
      expect(scaleFactor.value).toBe(2)

      recipeBatchSize.value = 20
      expect(scaleFactor.value).toBe(1)
    })

    it('recomputes when target batch size changes directly', () => {
      const recipeBatchSize = ref<number | null>(10)
      const recipeBatchSizeUnit = ref<string | null>('bbl')
      const { scaleFactor, targetBatchSize } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

      targetBatchSize.value = 30
      expect(scaleFactor.value).toBe(3)
    })
  })
})
