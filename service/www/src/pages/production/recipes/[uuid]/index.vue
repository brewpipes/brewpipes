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
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useRouteUuid } from '@/composables/useRouteUuid'

  const router = useRouter()
  const { getRecipe } = useProductionApi()
  const { uuid: routeUuid } = useRouteUuid()

  const loading = ref(true)
  const error = ref<string | null>(null)
  const recipe = ref<Recipe | null>(null)

  async function loadRecipe () {
    const uuid = routeUuid.value
    if (!uuid) {
      error.value = 'Invalid recipe UUID'
      loading.value = false
      return
    }

    try {
      loading.value = true
      error.value = null

      recipe.value = await getRecipe(uuid)
    } catch (error_) {
      console.error('Failed to load recipe:', error_)
      error.value = error_ instanceof Error && error_.message.includes('404') ? 'Recipe not found' : 'Failed to load recipe. Please try again.'
    } finally {
      loading.value = false
    }
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
