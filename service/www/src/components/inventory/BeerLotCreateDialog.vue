<template>
  <v-dialog
    :fullscreen="xs"
    :max-width="xs ? '100%' : 600"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-icon class="mr-2" icon="mdi-package-variant-closed" />
        <span class="text-h6">Create Beer Lot</span>
        <v-spacer />
        <v-btn
          aria-label="Close"
          :disabled="saving"
          icon="mdi-close"
          size="small"
          variant="text"
          @click="handleCancel"
        />
      </v-card-title>

      <v-divider />

      <v-card-text>
        <v-autocomplete
          v-model="form.production_batch_uuid"
          density="comfortable"
          hint="Search by batch name"
          item-title="title"
          item-value="value"
          :items="batchSelectItems"
          label="Batch"
          :loading="batchesLoading"
          persistent-hint
          :rules="[rules.required]"
        >
          <template #item="{ props: itemProps, item }">
            <v-list-item v-bind="itemProps">
              <template v-if="item.raw.subtitle" #subtitle>
                <span>{{ item.raw.subtitle }}</span>
              </template>
            </v-list-item>
          </template>
          <template #no-data>
            <v-list-item>
              <v-list-item-title>No batches found</v-list-item-title>
              <v-list-item-subtitle>Create a batch in Production first</v-list-item-subtitle>
            </v-list-item>
          </template>
        </v-autocomplete>

        <v-text-field
          v-model="form.lot_code"
          density="comfortable"
          label="Lot code"
        />

        <v-select
          v-model="form.container"
          clearable
          density="comfortable"
          :items="containerOptions"
          label="Container type"
        />

        <v-text-field
          v-model="form.package_format_name"
          density="comfortable"
          label="Package format name"
        />

        <v-row dense>
          <v-col cols="7">
            <v-text-field
              v-model="form.volume_per_unit"
              density="comfortable"
              inputmode="decimal"
              label="Volume per unit"
              type="number"
            />
          </v-col>
          <v-col cols="5">
            <v-combobox
              v-model="form.volume_per_unit_unit"
              density="comfortable"
              :items="volumeUnitOptions"
              label="Unit"
            />
          </v-col>
        </v-row>

        <v-text-field
          v-model="form.quantity"
          density="comfortable"
          inputmode="numeric"
          label="Quantity"
          type="number"
        />

        <v-select
          v-model="form.stock_location_uuid"
          clearable
          density="comfortable"
          :items="stockLocationSelectItems"
          label="Stock location"
        >
          <template #no-data>
            <v-list-item>
              <v-list-item-title>No stock locations</v-list-item-title>
              <v-list-item-subtitle>Create stock locations in Inventory settings</v-list-item-subtitle>
            </v-list-item>
          </template>
        </v-select>

        <v-text-field
          v-model="form.packaged_at"
          density="comfortable"
          label="Packaged at"
          type="datetime-local"
        />

        <v-text-field
          v-model="form.best_by"
          density="comfortable"
          label="Best by"
          type="date"
        />

        <v-textarea
          v-model="form.notes"
          auto-grow
          density="comfortable"
          label="Notes"
          rows="2"
        />
      </v-card-text>

      <v-divider />

      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="handleSubmit"
        >
          Create lot
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Batch, CreateBeerLotRequest, StockLocation } from '@/types'
  import { computed, reactive, ref, watch } from 'vue'
  import { useDisplay } from 'vuetify'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { normalizeDateTime, normalizeDateOnly, normalizeText, toNumber } from '@/utils/normalize'

  const containerOptions = ['keg', 'can', 'bottle', 'cask', 'growler', 'other']
  const volumeUnitOptions = ['bbl', 'gal', 'l', 'ml', 'oz']

  const props = defineProps<{
    modelValue: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'created': []
  }>()

  const { xs } = useDisplay()
  const { getBatches } = useProductionApi()
  const { createBeerLot, getStockLocations } = useInventoryApi()
  const { showNotice } = useSnackbar()

  // Reference data
  const batches = ref<Batch[]>([])
  const stockLocations = ref<StockLocation[]>([])
  const batchesLoading = ref(false)

  // Form state
  const saving = ref(false)

  const form = reactive({
    production_batch_uuid: null as string | null,
    lot_code: '',
    container: null as string | null,
    package_format_name: '',
    volume_per_unit: '',
    volume_per_unit_unit: '',
    quantity: '',
    stock_location_uuid: null as string | null,
    packaged_at: '',
    best_by: '',
    notes: '',
  })

  const rules = {
    required: (v: string | null) => (v !== null && v !== '' && String(v).trim() !== '') || 'Required',
  }

  // Computed: select items
  const batchSelectItems = computed(() =>
    batches.value.map(batch => ({
      title: batch.short_name,
      value: batch.uuid,
      subtitle: batch.recipe_name ?? undefined,
    })),
  )

  const stockLocationSelectItems = computed(() =>
    stockLocations.value.map(loc => ({
      title: loc.name,
      value: loc.uuid,
    })),
  )

  const isFormValid = computed(() => {
    return !!form.production_batch_uuid
  })

  // Watch dialog open to load data and reset form
  watch(
    () => props.modelValue,
    async isOpen => {
      if (isOpen) {
        resetForm()
        await loadReferenceData()
      }
    },
  )

  function resetForm () {
    form.production_batch_uuid = null
    form.lot_code = ''
    form.container = null
    form.package_format_name = ''
    form.volume_per_unit = ''
    form.volume_per_unit_unit = ''
    form.quantity = ''
    form.stock_location_uuid = null
    form.packaged_at = ''
    form.best_by = ''
    form.notes = ''
  }

  async function loadReferenceData () {
    batchesLoading.value = true
    try {
      const [batchResult, locationResult] = await Promise.allSettled([
        getBatches(),
        getStockLocations(),
      ])

      if (batchResult.status === 'fulfilled') {
        batches.value = batchResult.value
      }
      if (locationResult.status === 'fulfilled') {
        stockLocations.value = locationResult.value
      }
    } finally {
      batchesLoading.value = false
    }
  }

  async function handleSubmit () {
    if (!isFormValid.value || !form.production_batch_uuid) return

    saving.value = true
    try {
      const payload: CreateBeerLotRequest = {
        production_batch_uuid: form.production_batch_uuid,
        lot_code: normalizeText(form.lot_code),
        container: form.container,
        package_format_name: normalizeText(form.package_format_name),
        volume_per_unit: toNumber(form.volume_per_unit),
        volume_per_unit_unit: normalizeText(form.volume_per_unit_unit),
        quantity: toNumber(form.quantity) as number | null,
        stock_location_uuid: form.stock_location_uuid,
        packaged_at: normalizeDateTime(form.packaged_at),
        best_by: normalizeDateOnly(form.best_by),
        notes: normalizeText(form.notes),
      }

      await createBeerLot(payload)
      showNotice('Beer lot created')
      emit('created')
      emit('update:modelValue', false)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to create beer lot'
      showNotice(message, 'error')
    } finally {
      saving.value = false
    }
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
