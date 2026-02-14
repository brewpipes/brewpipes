<template>
  <div>
    <!-- Yeast Section -->
    <v-card class="sub-card mb-4" variant="outlined">
      <v-card-title class="text-subtitle-1 d-flex align-center">
        Yeast
        <v-spacer />
        <v-btn
          aria-label="Add yeast"
          icon="mdi-plus"
          size="small"
          variant="text"
          @click="emit('create', 'yeast')"
        />
      </v-card-title>
      <v-card-text>
        <v-progress-linear
          v-if="loading"
          class="mb-4"
          color="primary"
          indeterminate
        />

        <template v-if="!loading && yeasts.length > 0">
          <v-card
            v-for="yeast in yeasts"
            :key="yeast.uuid"
            class="mb-2 ingredient-card"
            variant="tonal"
          >
            <v-card-text class="pa-3">
              <div class="d-flex align-center justify-space-between mb-1">
                <span class="font-weight-medium">{{ yeast.name }}</span>
                <v-chip size="x-small" variant="tonal">
                  {{ formatUseType(yeast.use_type) }}
                </v-chip>
              </div>
              <div class="d-flex flex-wrap ga-3 text-body-2">
                <span>{{ formatAmount(yeast.amount, yeast.amount_unit) }}</span>
                <span class="text-medium-emphasis">{{ formatStage(yeast.use_stage) }}</span>
              </div>
              <div v-if="yeast.notes" class="text-caption text-medium-emphasis mt-1">
                {{ yeast.notes }}
              </div>
              <div class="d-flex justify-end mt-2">
                <v-btn
                  aria-label="Edit"
                  icon="mdi-pencil"
                  size="x-small"
                  variant="text"
                  @click="emit('edit', yeast)"
                />
                <v-btn
                  aria-label="Delete"
                  color="error"
                  icon="mdi-delete"
                  size="x-small"
                  variant="text"
                  @click="emit('delete', yeast)"
                />
              </div>
            </v-card-text>
          </v-card>
        </template>

        <v-alert
          v-if="!loading && yeasts.length === 0"
          density="compact"
          type="info"
          variant="tonal"
        >
          <div class="d-flex align-center justify-space-between flex-wrap ga-2">
            <span>No yeast added yet.</span>
            <v-btn
              color="primary"
              size="small"
              variant="text"
              @click="emit('create', 'yeast')"
            >
              Add yeast
            </v-btn>
          </div>
        </v-alert>
      </v-card-text>
    </v-card>

    <!-- Adjuncts Section -->
    <v-card class="sub-card mb-4" variant="outlined">
      <v-card-title class="text-subtitle-1 d-flex align-center">
        Adjuncts
        <v-spacer />
        <v-btn
          aria-label="Add adjunct"
          icon="mdi-plus"
          size="small"
          variant="text"
          @click="emit('create', 'adjunct')"
        />
      </v-card-title>
      <v-card-text>
        <template v-if="!loading && adjuncts.length > 0">
          <v-card
            v-for="adjunct in adjuncts"
            :key="adjunct.uuid"
            class="mb-2 ingredient-card"
            variant="tonal"
          >
            <v-card-text class="pa-3">
              <div class="d-flex align-center justify-space-between mb-1">
                <span class="font-weight-medium">{{ adjunct.name }}</span>
                <v-chip size="x-small" variant="tonal">
                  {{ formatStage(adjunct.use_stage) }}
                </v-chip>
              </div>
              <div class="text-body-2">
                {{ formatAmount(adjunct.amount, adjunct.amount_unit) }}
              </div>
              <div v-if="adjunct.notes" class="text-caption text-medium-emphasis mt-1">
                {{ adjunct.notes }}
              </div>
              <div class="d-flex justify-end mt-2">
                <v-btn
                  aria-label="Edit"
                  icon="mdi-pencil"
                  size="x-small"
                  variant="text"
                  @click="emit('edit', adjunct)"
                />
                <v-btn
                  aria-label="Delete"
                  color="error"
                  icon="mdi-delete"
                  size="x-small"
                  variant="text"
                  @click="emit('delete', adjunct)"
                />
              </div>
            </v-card-text>
          </v-card>
        </template>

        <v-alert
          v-if="!loading && adjuncts.length === 0"
          density="compact"
          type="info"
          variant="tonal"
        >
          <div class="d-flex align-center justify-space-between flex-wrap ga-2">
            <span>No adjuncts added yet.</span>
            <v-btn
              color="primary"
              size="small"
              variant="text"
              @click="emit('create', 'adjunct')"
            >
              Add adjunct
            </v-btn>
          </div>
        </v-alert>
      </v-card-text>
    </v-card>

    <!-- Water Chemistry Section -->
    <v-card class="sub-card" variant="outlined">
      <v-card-title class="text-subtitle-1 d-flex align-center">
        Water Chemistry
        <v-spacer />
        <v-menu>
          <template #activator="{ props }">
            <v-btn
              v-bind="props"
              aria-label="Add water chemistry"
              icon="mdi-plus"
              size="small"
              variant="text"
            />
          </template>
          <v-list density="compact" nav>
            <v-list-item @click="emit('create', 'salt')">
              <v-list-item-title>Add Salt</v-list-item-title>
            </v-list-item>
            <v-list-item @click="emit('create', 'chemical')">
              <v-list-item-title>Add Chemical</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-card-title>
      <v-card-text>
        <template v-if="!loading && waterChemistry.length > 0">
          <v-card
            v-for="item in waterChemistry"
            :key="item.uuid"
            class="mb-2 ingredient-card"
            variant="tonal"
          >
            <v-card-text class="pa-3">
              <div class="d-flex align-center justify-space-between mb-1">
                <span class="font-weight-medium">{{ item.name }}</span>
                <v-chip size="x-small" variant="tonal">
                  {{ formatIngredientType(item.ingredient_type) }}
                </v-chip>
              </div>
              <div class="d-flex flex-wrap ga-3 text-body-2">
                <span>{{ formatAmount(item.amount, item.amount_unit) }}</span>
                <span class="text-medium-emphasis">{{ formatStage(item.use_stage) }}</span>
              </div>
              <div v-if="item.notes" class="text-caption text-medium-emphasis mt-1">
                {{ item.notes }}
              </div>
              <div class="d-flex justify-end mt-2">
                <v-btn
                  aria-label="Edit"
                  icon="mdi-pencil"
                  size="x-small"
                  variant="text"
                  @click="emit('edit', item)"
                />
                <v-btn
                  aria-label="Delete"
                  color="error"
                  icon="mdi-delete"
                  size="x-small"
                  variant="text"
                  @click="emit('delete', item)"
                />
              </div>
            </v-card-text>
          </v-card>
        </template>

        <v-alert
          v-if="!loading && waterChemistry.length === 0"
          density="compact"
          type="info"
          variant="tonal"
        >
          <div class="d-flex align-center justify-space-between flex-wrap ga-2">
            <span>No water chemistry additions yet.</span>
            <v-btn
              color="primary"
              size="small"
              variant="text"
              @click="emit('create', 'salt')"
            >
              Add salt
            </v-btn>
          </div>
        </v-alert>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts" setup>
  import type { RecipeIngredient, RecipeIngredientType } from '@/types'

  defineProps<{
    yeasts: RecipeIngredient[]
    adjuncts: RecipeIngredient[]
    waterChemistry: RecipeIngredient[]
    loading: boolean
  }>()

  const emit = defineEmits<{
    create: [type: RecipeIngredientType]
    edit: [ingredient: RecipeIngredient]
    delete: [ingredient: RecipeIngredient]
  }>()

  function formatAmount (amount: number, unit: string): string {
    return `${amount.toFixed(2)} ${unit}`
  }

  function formatStage (stage: string): string {
    const labels: Record<string, string> = {
      mash: 'Mash',
      boil: 'Boil',
      whirlpool: 'Whirlpool',
      fermentation: 'Fermentation',
      packaging: 'Packaging',
    }
    return labels[stage] ?? stage.charAt(0).toUpperCase() + stage.slice(1)
  }

  function formatUseType (useType: string | null): string {
    if (!useType) return 'Primary'
    const labels: Record<string, string> = {
      primary: 'Primary',
      secondary: 'Secondary',
      bottle: 'Bottle',
      other: 'Other',
    }
    return labels[useType] ?? useType.charAt(0).toUpperCase() + useType.slice(1)
  }

  function formatIngredientType (type: string): string {
    const labels: Record<string, string> = {
      salt: 'Salt',
      chemical: 'Chemical',
    }
    return labels[type] ?? type.charAt(0).toUpperCase() + type.slice(1)
  }
</script>

<style scoped>
.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.ingredient-card {
  background: rgba(var(--v-theme-surface), 0.5);
}
</style>
