<template>
  <v-card class="sub-card" variant="outlined">
    <v-card-text>
      <!-- Notes Section -->
      <div class="text-overline text-medium-emphasis mb-2">Notes</div>
      <div v-if="recipe.notes" class="text-body-1 mb-4">
        {{ recipe.notes }}
      </div>
      <div v-else class="text-body-2 text-medium-emphasis mb-4">
        No notes for this recipe.
      </div>

      <v-divider class="mb-4" />

      <!-- Recipe Info -->
      <div class="text-overline text-medium-emphasis mb-2">Recipe Info</div>
      <v-row dense>
        <v-col cols="12" md="6">
          <div class="info-item">
            <div class="info-label">Created</div>
            <div class="info-value">{{ formatDateTime(recipe.created_at) }}</div>
          </div>
        </v-col>
        <v-col cols="12" md="6">
          <div class="info-item">
            <div class="info-label">Last Updated</div>
            <div class="info-value">{{ formatDateTime(recipe.updated_at) }}</div>
          </div>
        </v-col>
        <v-col v-if="recipe.brewhouse_efficiency" cols="12" md="6">
          <div class="info-item">
            <div class="info-label">Brewhouse Efficiency</div>
            <div class="info-value">{{ recipe.brewhouse_efficiency }}%</div>
          </div>
        </v-col>
        <v-col v-if="recipe.ibu_method" cols="12" md="6">
          <div class="info-item">
            <div class="info-label">IBU Method</div>
            <div class="info-value">{{ formatIbuMethod(recipe.ibu_method) }}</div>
          </div>
        </v-col>
      </v-row>

      <v-divider class="my-4" />

      <!-- Batches Using This Recipe -->
      <div class="text-overline text-medium-emphasis mb-2">Batches</div>
      <v-alert
        density="compact"
        type="info"
        variant="tonal"
      >
        Batch history will be available in a future update.
      </v-alert>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import type { Recipe } from '@/composables/useProductionApi'
  import { useBrewingFormatters, useFormatters } from '@/composables/useFormatters'

  defineProps<{
    recipe: Recipe
  }>()

  const { formatDateTime } = useFormatters()
  const { formatIbuMethod } = useBrewingFormatters()
</script>

<style scoped>
.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.info-item {
  margin-bottom: 8px;
}

.info-label {
  font-size: 0.75rem;
  color: rgba(var(--v-theme-on-surface), 0.55);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.info-value {
  font-size: 0.95rem;
  color: rgba(var(--v-theme-on-surface), 0.87);
}
</style>
