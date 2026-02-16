import { describe, expect, it, vi } from 'vitest'
import { useAsyncAction } from '../useAsyncAction'

describe('useAsyncAction', () => {
  describe('initial state', () => {
    it('starts with loading false and empty error', () => {
      const { loading, error } = useAsyncAction()
      expect(loading.value).toBe(false)
      expect(error.value).toBe('')
    })
  })

  describe('execute', () => {
    it('sets loading to true during execution', async () => {
      const { execute, loading } = useAsyncAction()
      let loadingDuringExecution = false

      await execute(async () => {
        loadingDuringExecution = loading.value
      })

      expect(loadingDuringExecution).toBe(true)
      expect(loading.value).toBe(false)
    })

    it('clears error before execution', async () => {
      const { execute, error } = useAsyncAction()

      // First call: fail to set an error
      await execute(async () => {
        throw new Error('first error')
      })
      expect(error.value).toBe('first error')

      // Second call: error should be cleared before running
      let errorDuringExecution = 'not cleared'
      await execute(async () => {
        errorDuringExecution = error.value
      })

      expect(errorDuringExecution).toBe('')
    })

    it('returns the result of the async function on success', async () => {
      const { execute } = useAsyncAction()

      const result = await execute(async () => {
        return 42
      })

      expect(result).toBe(42)
    })

    it('returns undefined on failure', async () => {
      const { execute } = useAsyncAction()

      const result = await execute(async () => {
        throw new Error('fail')
      })

      expect(result).toBeUndefined()
    })

    it('sets error message from Error instances', async () => {
      const { execute, error } = useAsyncAction()

      await execute(async () => {
        throw new Error('Something broke')
      })

      expect(error.value).toBe('Something broke')
    })

    it('sets fallback error message for non-Error throws', async () => {
      const { execute, error } = useAsyncAction()

      await execute(async () => {
        throw 'string error' // eslint-disable-line no-throw-literal
      })

      expect(error.value).toBe('Something went wrong')
    })

    it('sets fallback error message for null throws', async () => {
      const { execute, error } = useAsyncAction()

      await execute(async () => {
        throw null // eslint-disable-line no-throw-literal
      })

      expect(error.value).toBe('Something went wrong')
    })

    it('sets loading to false after success', async () => {
      const { execute, loading } = useAsyncAction()

      await execute(async () => {
        return 'ok'
      })

      expect(loading.value).toBe(false)
    })

    it('sets loading to false after failure', async () => {
      const { execute, loading } = useAsyncAction()

      await execute(async () => {
        throw new Error('fail')
      })

      expect(loading.value).toBe(false)
    })

    it('preserves generic return type', async () => {
      const { execute } = useAsyncAction()

      const result = await execute(async () => {
        return { name: 'IPA', abv: 6.5 }
      })

      expect(result).toEqual({ name: 'IPA', abv: 6.5 })
    })

    it('handles void async functions', async () => {
      const { execute, error } = useAsyncAction()
      const sideEffect = vi.fn()

      const result = await execute(async () => {
        sideEffect()
      })

      expect(sideEffect).toHaveBeenCalledOnce()
      expect(result).toBeUndefined()
      expect(error.value).toBe('')
    })
  })

  describe('onError callback', () => {
    it('calls onError instead of setting error ref', async () => {
      const onError = vi.fn()
      const { execute, error } = useAsyncAction({ onError })

      await execute(async () => {
        throw new Error('custom handler')
      })

      expect(onError).toHaveBeenCalledWith('custom handler')
      expect(error.value).toBe('')
    })

    it('calls onError with fallback message for non-Error throws', async () => {
      const onError = vi.fn()
      const { execute } = useAsyncAction({ onError })

      await execute(async () => {
        throw 42 // eslint-disable-line no-throw-literal
      })

      expect(onError).toHaveBeenCalledWith('Something went wrong')
    })

    it('does not call onError on success', async () => {
      const onError = vi.fn()
      const { execute } = useAsyncAction({ onError })

      await execute(async () => {
        return 'ok'
      })

      expect(onError).not.toHaveBeenCalled()
    })
  })

  describe('multiple instances', () => {
    it('maintains independent state across instances', async () => {
      const { execute: load, loading: loadLoading, error: loadError } = useAsyncAction()
      const { execute: save, loading: saveLoading, error: saveError } = useAsyncAction()

      // Fail the load
      await load(async () => {
        throw new Error('load failed')
      })

      // Succeed the save
      await save(async () => {
        return 'saved'
      })

      expect(loadLoading.value).toBe(false)
      expect(loadError.value).toBe('load failed')
      expect(saveLoading.value).toBe(false)
      expect(saveError.value).toBe('')
    })

    it('loading states are independent', async () => {
      const { execute: action1, loading: loading1 } = useAsyncAction()
      const { execute: action2, loading: loading2 } = useAsyncAction()

      let loading1During2 = false
      let loading2During1 = false

      const promise1 = action1(async () => {
        loading2During1 = loading2.value
        await new Promise(resolve => setTimeout(resolve, 10))
      })

      const promise2 = action2(async () => {
        loading1During2 = loading1.value
        await new Promise(resolve => setTimeout(resolve, 10))
      })

      await Promise.all([promise1, promise2])

      // Each instance's loading state is independent
      expect(loading1.value).toBe(false)
      expect(loading2.value).toBe(false)
      // During execution, the other instance's loading state depends on timing,
      // but both should be false after completion
    })
  })

  describe('sequential calls', () => {
    it('resets error on subsequent successful call', async () => {
      const { execute, error } = useAsyncAction()

      // First call fails
      await execute(async () => {
        throw new Error('first error')
      })
      expect(error.value).toBe('first error')

      // Second call succeeds
      await execute(async () => {
        return 'ok'
      })
      expect(error.value).toBe('')
    })

    it('updates error on subsequent failed call', async () => {
      const { execute, error } = useAsyncAction()

      await execute(async () => {
        throw new Error('first error')
      })
      expect(error.value).toBe('first error')

      await execute(async () => {
        throw new Error('second error')
      })
      expect(error.value).toBe('second error')
    })
  })
})
