<template>
  <v-container class="removals-page" fluid>
    <v-card class="section-card">
      <v-card-title class="card-title-responsive">
        <div class="d-flex align-center">
          <v-icon class="mr-2" icon="mdi-delete-variant" />
          Removals
        </div>
        <div class="card-title-actions">
          <v-btn
            color="error"
            prepend-icon="mdi-plus"
            size="small"
            variant="tonal"
            @click="showRemovalDialog = true"
          >
            <span class="d-none d-sm-inline">Record Removal</span>
            <span class="d-sm-none">New</span>
          </v-btn>
          <v-btn
            :loading="loading"
            size="small"
            variant="text"
            @click="refreshAll"
          >
            <v-icon class="mr-1" icon="mdi-refresh" size="small" />
            <span class="d-none d-sm-inline">Refresh</span>
          </v-btn>
        </div>
      </v-card-title>

      <v-card-text>
        <v-alert
          v-if="errorMessage"
          class="mb-3"
          density="compact"
          type="error"
          variant="tonal"
        >
          {{ errorMessage }}
        </v-alert>

        <!-- Summary cards -->
        <v-row class="mb-4" dense>
          <v-col cols="6" md="3">
            <div class="metric-card">
              <div class="metric-label">Total Removals</div>
              <div class="metric-value">{{ summary?.total_count ?? 0 }}</div>
              <div class="text-caption text-medium-emphasis">
                {{ formatBbl(summary?.total_bbl) }}
              </div>
            </div>
          </v-col>
          <v-col cols="6" md="3">
            <div class="metric-card metric-card--error">
              <div class="metric-label">Dumps</div>
              <div class="metric-value">{{ categoryStat('dump').count }}</div>
              <div class="text-caption text-medium-emphasis">
                {{ formatBbl(categoryStat('dump').bbl) }}
              </div>
            </div>
          </v-col>
          <v-col cols="6" md="3">
            <div class="metric-card metric-card--warning">
              <div class="metric-label">Waste & Spillage</div>
              <div class="metric-value">{{ categoryStat('waste').count }}</div>
              <div class="text-caption text-medium-emphasis">
                {{ formatBbl(categoryStat('waste').bbl) }}
              </div>
            </div>
          </v-col>
          <v-col cols="6" md="3">
            <div class="metric-card metric-card--info">
              <div class="metric-label">Samples</div>
              <div class="metric-value">{{ categoryStat('sample').count }}</div>
              <div class="text-caption text-medium-emphasis">
                {{ formatBbl(categoryStat('sample').bbl) }}
              </div>
            </div>
          </v-col>
        </v-row>

        <!-- Date range filter -->
        <v-chip-group
          v-model="dateRange"
          class="mb-4"
          color="primary"
          mandatory
          selected-class="text-primary"
        >
          <v-chip filter value="month">This Month</v-chip>
          <v-chip filter value="quarter">This Quarter</v-chip>
          <v-chip filter value="year">This Year</v-chip>
          <v-chip filter value="all">All Time</v-chip>
        </v-chip-group>

        <!-- Desktop table view -->
        <div class="d-none d-md-block">
          <v-data-table
            class="data-table"
            density="compact"
            :headers="tableHeaders"
            hover
            item-value="uuid"
            :items="removals"
            :loading="loading"
          >
            <template #item.removed_at="{ item }">
              {{ formatDateTime(item.removed_at) }}
            </template>

            <template #item.category="{ item }">
              <v-chip
                :color="categoryColors[item.category]"
                :prepend-icon="categoryIcons[item.category]"
                size="x-small"
                variant="tonal"
              >
                {{ categoryLabels[item.category] }}
              </v-chip>
            </template>

            <template #item.reason="{ item }">
              {{ reasonLabels[item.reason] }}
            </template>

            <template #item.batch_name="{ item }">
              <router-link
                v-if="item.batch_uuid && batchMap.get(item.batch_uuid)"
                class="batch-link"
                :to="`/batches/${item.batch_uuid}`"
              >
                {{ batchName(item.batch_uuid) }}
              </router-link>
              <span v-else-if="item.batch_uuid" class="text-medium-emphasis">
                {{ batchName(item.batch_uuid) }}
              </span>
              <span v-else class="text-medium-emphasis">—</span>
            </template>

            <template #item.amount="{ item }">
              {{ formatAmountPreferred(item.amount, item.amount_unit) }}
            </template>

            <template #item.amount_bbl="{ item }">
              {{ item.amount_bbl != null ? item.amount_bbl.toFixed(2) : '—' }}
            </template>

            <template #item.notes="{ item }">
              <span
                v-if="item.notes"
                class="text-medium-emphasis text-truncate d-inline-block"
                style="max-width: 200px;"
              >
                <v-tooltip :text="item.notes" location="top">
                  <template #activator="{ props: tooltipProps }">
                    <span v-bind="tooltipProps">{{ item.notes }}</span>
                  </template>
                </v-tooltip>
              </span>
              <span v-else class="text-medium-emphasis">—</span>
            </template>

            <template #item.actions="{ item }">
              <v-btn
                aria-label="Delete removal"
                color="error"
                density="compact"
                icon="mdi-delete-outline"
                size="small"
                variant="text"
                @click="confirmDelete(item)"
              />
            </template>

            <template #no-data>
              <div class="text-center py-8 text-medium-emphasis">
                <v-icon class="mb-2" icon="mdi-delete-variant" size="48" />
                <div class="text-body-1">No removals recorded.</div>
                <div class="text-body-2 mt-1">
                  Record a removal to start tracking losses.
                </div>
              </div>
            </template>
          </v-data-table>
        </div>

        <!-- Mobile card view -->
        <div class="d-md-none">
          <v-progress-linear v-if="loading" color="primary" indeterminate />

          <div
            v-if="!loading && removals.length === 0"
            class="text-center py-8 text-medium-emphasis"
          >
            <v-icon class="mb-2" icon="mdi-delete-variant" size="48" />
            <div class="text-body-1">No removals recorded.</div>
            <div class="text-body-2 mt-1">
              Record a removal to start tracking losses.
            </div>
          </div>

          <v-card
            v-for="item in removals"
            :key="item.uuid"
            class="mb-3"
            variant="outlined"
          >
            <v-card-title class="d-flex align-center py-2 text-body-1">
              <v-chip
                :color="categoryColors[item.category]"
                :prepend-icon="categoryIcons[item.category]"
                size="x-small"
                variant="tonal"
              >
                {{ categoryLabels[item.category] }}
              </v-chip>
              <v-spacer />
              <v-btn
                aria-label="Delete removal"
                color="error"
                density="compact"
                icon="mdi-delete-outline"
                size="small"
                variant="text"
                @click="confirmDelete(item)"
              />
            </v-card-title>

            <v-card-text class="pt-0">
              <div class="d-flex justify-space-between text-body-2 mb-1">
                <span class="text-medium-emphasis">Date</span>
                <span>{{ formatDateTime(item.removed_at) }}</span>
              </div>
              <div class="d-flex justify-space-between text-body-2 mb-1">
                <span class="text-medium-emphasis">Reason</span>
                <span>{{ reasonLabels[item.reason] }}</span>
              </div>
              <div class="d-flex justify-space-between text-body-2 mb-1">
                <span class="text-medium-emphasis">Amount</span>
                <span class="font-weight-medium">{{ formatAmountPreferred(item.amount, item.amount_unit) }}</span>
              </div>
              <div class="d-flex justify-space-between text-body-2 mb-1">
                <span class="text-medium-emphasis">BBL</span>
                <span>{{ item.amount_bbl != null ? item.amount_bbl.toFixed(2) : '—' }}</span>
              </div>
              <div v-if="item.batch_uuid" class="d-flex justify-space-between text-body-2 mb-1">
                <span class="text-medium-emphasis">Batch</span>
                <router-link
                  v-if="batchMap.get(item.batch_uuid)"
                  class="batch-link"
                  :to="`/batches/${item.batch_uuid}`"
                >
                  {{ batchName(item.batch_uuid) }}
                </router-link>
                <span v-else class="text-medium-emphasis">{{ batchName(item.batch_uuid) }}</span>
              </div>
              <div v-if="item.notes" class="text-body-2 text-medium-emphasis mt-2">
                {{ item.notes }}
              </div>
            </v-card-text>
          </v-card>
        </div>
      </v-card-text>
    </v-card>

    <!-- Record Removal dialog -->
    <RemovalDialog
      v-model="showRemovalDialog"
      @created="handleRemovalCreated"
    />

    <!-- Delete confirmation dialog -->
    <v-dialog v-model="showDeleteDialog" max-width="400" persistent>
      <v-card>
        <v-card-title class="text-h6">Delete Removal?</v-card-title>
        <v-card-text>
          This will remove this record. This action cannot be undone.
        </v-card-text>
        <v-card-actions class="justify-end">
          <v-btn :disabled="deleting" variant="text" @click="showDeleteDialog = false">Cancel</v-btn>
          <v-btn color="error" :loading="deleting" variant="tonal" @click="handleDelete">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script lang="ts" setup>
  import type { Batch, Removal, RemovalCategory, RemovalSummary } from '@/types'
  import { computed, onMounted, ref, watch } from 'vue'
  import RemovalDialog from '@/components/removal/RemovalDialog.vue'
  import { categoryColors, categoryIcons, categoryLabels, reasonLabels } from '@/components/removal/constants'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { formatDateTime } from '@/composables/useFormatters'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  // ==================== Composables ====================

  const { listRemovals, getRemovalSummary, deleteRemoval } = useInventoryApi()
  const { getBatches } = useProductionApi()
  const { showNotice } = useSnackbar()
  const { formatAmountPreferred } = useUnitPreferences()

  // ==================== State ====================

  const { execute, loading, error: errorMessage } = useAsyncAction()

  const removals = ref<Removal[]>([])
  const summary = ref<RemovalSummary | null>(null)
  const batches = ref<Batch[]>([])
  const dateRange = ref('month')

  const showRemovalDialog = ref(false)
  const showDeleteDialog = ref(false)
  const deletingRemoval = ref<Removal | null>(null)
  const deleting = ref(false)

  // ==================== Computed ====================

  const batchMap = computed(() =>
    new Map(batches.value.map(b => [b.uuid, b])),
  )

  const dateRangeParams = computed(() => {
    const now = new Date()
    let from: string | undefined

    if (dateRange.value === 'month') {
      from = new Date(now.getFullYear(), now.getMonth(), 1).toISOString()
    } else if (dateRange.value === 'quarter') {
      const quarterMonth = Math.floor(now.getMonth() / 3) * 3
      from = new Date(now.getFullYear(), quarterMonth, 1).toISOString()
    } else if (dateRange.value === 'year') {
      from = new Date(now.getFullYear(), 0, 1).toISOString()
    }

    return { from, to: undefined }
  })

  // ==================== Table ====================

  const tableHeaders = [
    { title: 'Date', key: 'removed_at', sortable: true },
    { title: 'Category', key: 'category', sortable: true },
    { title: 'Reason', key: 'reason', sortable: true },
    { title: 'Batch', key: 'batch_name', sortable: false },
    { title: 'Amount', key: 'amount', sortable: true },
    { title: 'BBL', key: 'amount_bbl', sortable: true },
    { title: 'Notes', key: 'notes', sortable: false },
    { title: '', key: 'actions', sortable: false, width: 48 },
  ]

  // ==================== Helpers ====================

  function batchName (uuid: string): string {
    return batchMap.value.get(uuid)?.short_name ?? uuid.slice(0, 8)
  }

  function formatBbl (value: number | null | undefined): string {
    if (value === null || value === undefined) return '0.00 BBL'
    return `${value.toFixed(2)} BBL`
  }

  function categoryStat (category: RemovalCategory): { count: number; bbl: number } {
    const entry = summary.value?.by_category?.find(c => c.category === category)
    return {
      count: entry?.count ?? 0,
      bbl: entry?.total_bbl ?? 0,
    }
  }

  // ==================== Data loading ====================

  onMounted(async () => {
    await refreshAll()
  })

  watch(dateRange, async () => {
    await refreshAll()
  })

  async function refreshAll () {
    await execute(async () => {
      const params = dateRangeParams.value
      const [removalsResult, summaryResult, batchResult] = await Promise.allSettled([
        listRemovals({ from: params.from, to: params.to }),
        getRemovalSummary({ from: params.from, to: params.to }),
        getBatches(),
      ])

      if (removalsResult.status === 'fulfilled') {
        removals.value = removalsResult.value
      } else {
        throw new Error('Unable to load removals')
      }

      if (summaryResult.status === 'fulfilled') {
        summary.value = summaryResult.value
      }

      // Cross-service data: graceful failure
      if (batchResult.status === 'fulfilled') {
        batches.value = batchResult.value
      }
    })
  }

  // ==================== Actions ====================

  function handleRemovalCreated () {
    refreshAll()
  }

  function confirmDelete (removal: Removal) {
    deletingRemoval.value = removal
    showDeleteDialog.value = true
  }

  async function handleDelete () {
    if (!deletingRemoval.value) return

    deleting.value = true
    try {
      await deleteRemoval(deletingRemoval.value.uuid)
      showNotice('Removal deleted')
      showDeleteDialog.value = false
      deletingRemoval.value = null
      await refreshAll()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to delete removal'
      showNotice(message, 'error')
    } finally {
      deleting.value = false
    }
  }
</script>

<style scoped>
.removals-page {
  position: relative;
}

.metric-card {
  padding: 12px 16px;
  border-radius: 8px;
  background: rgba(var(--v-theme-surface), 0.5);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  text-align: center;
}

.metric-card--error {
  background: rgba(var(--v-theme-error), 0.08);
  border-color: rgba(var(--v-theme-error), 0.2);
}

.metric-card--warning {
  background: rgba(var(--v-theme-warning), 0.08);
  border-color: rgba(var(--v-theme-warning), 0.2);
}

.metric-card--info {
  background: rgba(var(--v-theme-info), 0.08);
  border-color: rgba(var(--v-theme-info), 0.2);
}

.metric-label {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: rgba(var(--v-theme-on-surface), 0.55);
  margin-bottom: 4px;
}

.metric-value {
  font-size: 1.25rem;
  font-weight: 600;
  color: rgba(var(--v-theme-on-surface), 0.87);
}
</style>
