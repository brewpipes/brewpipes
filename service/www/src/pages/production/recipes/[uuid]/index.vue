<template>
  <v-container class="pa-4" fluid>
    <v-alert
      v-if="loading"
      density="comfortable"
      type="info"
      variant="tonal"
    >
      Loading recipe...
    </v-alert>

    <v-alert
      v-else-if="error"
      density="comfortable"
      type="error"
      variant="tonal"
    >
      {{ error }}
    </v-alert>

    <RecipeDetails
      v-else-if="recipe"
      :recipe="recipe"
      @back="handleBack"
      @deleted="handleDeleted"
      @updated="handleUpdated"
    />
  </v-container>
</template>

<script lang="ts" setup>
  import type { Recipe } from '@/types'
  import { onMounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import RecipeDetails from '@/components/recipe/RecipeDetails.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useRouteUuid } from '@/composables/useRouteUuid'

  const router = useRouter()
  const { getRecipe } = useProductionApi()
  const { uuid: routeUuid } = useRouteUuid()

  const error = ref('')
  const { execute: executeLoad, loading } = useAsyncAction({
    onError: (message) => {
      error.value = message.includes('404') ? 'Recipe not found' : 'Failed to load recipe. Please try again.'
    },
  })
  loading.value = true

  const recipe = ref<Recipe | null>(null)

  async function loadRecipe () {
    const uuid = routeUuid.value
    if (!uuid) {
      error.value = 'Invalid recipe UUID'
      loading.value = false
      return
    }

    await executeLoad(async () => {
      recipe.value = await getRecipe(uuid)
    })
  }

  function handleBack () {
    router.push('/production/recipes')
  }

  function handleDeleted () {
    router.push('/production/recipes')
  }

  function handleUpdated (updatedRecipe: Recipe) {
    recipe.value = updatedRecipe
  }

  onMounted(() => {
    loadRecipe()
  })
</script>
