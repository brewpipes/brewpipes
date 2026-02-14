import { computed } from 'vue'
import { useRoute } from 'vue-router'

/**
 * Extract and validate `route.params.uuid` as a reactive computed string.
 * Returns null if the param is missing, empty, or not a string.
 */
export function useRouteUuid () {
  const route = useRoute()
  const uuid = computed(() => {
    const params = route.params
    if ('uuid' in params) {
      const param = params.uuid
      if (typeof param === 'string' && param.trim()) {
        return param
      }
    }
    return null
  })
  return { uuid }
}
