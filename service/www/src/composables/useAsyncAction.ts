import { ref } from 'vue'

/**
 * Composable that wraps async operations with loading state, error handling,
 * and optional custom error callbacks.
 *
 * Eliminates repetitive try/catch/loading/error boilerplate across page components.
 *
 * @example
 * ```ts
 * const { execute, loading, error } = useAsyncAction()
 *
 * // Basic usage
 * await execute(async () => {
 *   batches.value = await getBatches()
 * })
 *
 * // With return value
 * const result = await execute(async () => {
 *   return await api.getBatch(uuid)
 * })
 *
 * // Multiple instances for separate loading/error states
 * const { execute: load, loading, error: loadError } = useAsyncAction()
 * const { execute: save, loading: saving, error: saveError } = useAsyncAction()
 * ```
 */
export function useAsyncAction (options?: { onError?: (message: string) => void }) {
  const loading = ref(false)
  const error = ref('')

  async function execute<T> (fn: () => Promise<T>): Promise<T | undefined> {
    loading.value = true
    error.value = ''
    try {
      const result = await fn()
      return result
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Something went wrong'
      if (options?.onError) {
        options.onError(message)
      } else {
        error.value = message
      }
      return undefined
    } finally {
      loading.value = false
    }
  }

  return { execute, loading, error }
}
