<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="card-title-responsive">
        <div class="d-flex align-center">
          <v-icon class="mr-2" icon="mdi-barley" />
          Ingredients
        </div>
        <v-btn
          color="primary"
          :icon="$vuetify.display.xs"
          :loading="ingredientLoading"
          :prepend-icon="$vuetify.display.xs ? undefined : 'mdi-plus'"
          size="small"
          variant="text"
          @click="openIngredientDialog"
        >
          <v-icon v-if="$vuetify.display.xs" icon="mdi-plus" />
          <span v-else>New ingredient</span>
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-alert
          v-if="ingredientErrorMessage"
          class="mb-3"
          density="compact"
          type="error"
          variant="tonal"
        >
          {{ ingredientErrorMessage }}
        </v-alert>
        <v-tabs v-model="activeTab" class="inventory-tabs" color="primary" show-arrows>
          <v-tab value="malt">Malt</v-tab>
          <v-tab value="hops">Hops</v-tab>
          <v-tab value="yeast">Yeast</v-tab>
          <v-tab value="other">Other</v-tab>
          <v-tab value="usage">Usage</v-tab>
          <v-tab value="received">Received</v-tab>
        </v-tabs>

        <v-window v-model="activeTab" class="mt-4">
          <!-- Malt Tab -->
          <v-window-item value="malt">
            <v-card variant="outlined">
              <v-card-title class="d-flex align-center">
                Malt lots
                <v-spacer />
                <v-btn :loading="lotLoading" size="small" variant="text" @click="loadLots">
                  Refresh
                </v-btn>
                <v-btn
                  class="ml-2"
                  color="primary"
                  prepend-icon="mdi-plus"
                  size="small"
                  variant="text"
                  @click="openLotDialog('fermentable')"
                >
                  Create lot
                </v-btn>
              </v-card-title>
              <v-card-text>
                <v-alert
                  v-if="lotErrorMessage"
                  class="mb-3"
                  density="compact"
                  type="error"
                  variant="tonal"
                >
                  {{ lotErrorMessage }}
                </v-alert>
                <v-data-table
                  class="data-table lot-table"
                  density="compact"
                  :headers="lotHeaders"
                  hover
                  item-value="uuid"
                  :items="maltLots"
                  :loading="lotLoading"
                  @click:row="(_event: Event, row: any) => openLotDetails(row.item.uuid)"
                >
                  <template #item.ingredient_uuid="{ item }">
                    {{ ingredientName(item.ingredient_uuid) }}
                  </template>
                  <template #item.received_amount="{ item }">
                    {{ formatAmountPreferred(item.received_amount, item.received_unit) }}
                  </template>
                  <template #item.received_at="{ item }">
                    {{ formatDateTime(item.received_at) }}
                  </template>
                  <template #item.best_by_at="{ item }">
                    {{ formatDateTime(item.best_by_at) }}
                  </template>
                  <template #item.expires_at="{ item }">
                    {{ formatDateTime(item.expires_at) }}
                  </template>
                  <template #no-data>
                    <div class="text-center py-4 text-medium-emphasis">No malt lots yet.</div>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-window-item>

          <!-- Hops Tab -->
          <v-window-item value="hops">
            <v-card variant="outlined">
              <v-card-title class="d-flex align-center">
                Hop lots
                <v-spacer />
                <v-btn :loading="lotLoading" size="small" variant="text" @click="loadLots">
                  Refresh
                </v-btn>
                <v-btn
                  class="ml-2"
                  color="primary"
                  prepend-icon="mdi-plus"
                  size="small"
                  variant="text"
                  @click="openLotDialog('hop')"
                >
                  Create lot
                </v-btn>
              </v-card-title>
              <v-card-text>
                <v-alert
                  v-if="lotErrorMessage"
                  class="mb-3"
                  density="compact"
                  type="error"
                  variant="tonal"
                >
                  {{ lotErrorMessage }}
                </v-alert>
                <v-data-table
                  class="data-table lot-table"
                  density="compact"
                  :headers="lotHeaders"
                  hover
                  item-value="uuid"
                  :items="hopLots"
                  :loading="lotLoading"
                  @click:row="(_event: Event, row: any) => openLotDetails(row.item.uuid)"
                >
                  <template #item.ingredient_uuid="{ item }">
                    {{ ingredientName(item.ingredient_uuid) }}
                  </template>
                  <template #item.received_amount="{ item }">
                    {{ formatAmountPreferred(item.received_amount, item.received_unit) }}
                  </template>
                  <template #item.received_at="{ item }">
                    {{ formatDateTime(item.received_at) }}
                  </template>
                  <template #item.best_by_at="{ item }">
                    {{ formatDateTime(item.best_by_at) }}
                  </template>
                  <template #item.expires_at="{ item }">
                    {{ formatDateTime(item.expires_at) }}
                  </template>
                  <template #no-data>
                    <div class="text-center py-4 text-medium-emphasis">No hop lots yet.</div>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-window-item>

          <!-- Yeast Tab -->
          <v-window-item value="yeast">
            <v-card variant="outlined">
              <v-card-title class="d-flex align-center">
                Yeast lots
                <v-spacer />
                <v-btn :loading="lotLoading" size="small" variant="text" @click="loadLots">
                  Refresh
                </v-btn>
                <v-btn
                  class="ml-2"
                  color="primary"
                  prepend-icon="mdi-plus"
                  size="small"
                  variant="text"
                  @click="openLotDialog('yeast')"
                >
                  Create lot
                </v-btn>
              </v-card-title>
              <v-card-text>
                <v-alert
                  v-if="lotErrorMessage"
                  class="mb-3"
                  density="compact"
                  type="error"
                  variant="tonal"
                >
                  {{ lotErrorMessage }}
                </v-alert>
                <v-data-table
                  class="data-table lot-table"
                  density="compact"
                  :headers="lotHeaders"
                  hover
                  item-value="uuid"
                  :items="yeastLots"
                  :loading="lotLoading"
                  @click:row="(_event: Event, row: any) => openLotDetails(row.item.uuid)"
                >
                  <template #item.ingredient_uuid="{ item }">
                    {{ ingredientName(item.ingredient_uuid) }}
                  </template>
                  <template #item.received_amount="{ item }">
                    {{ formatAmountPreferred(item.received_amount, item.received_unit) }}
                  </template>
                  <template #item.received_at="{ item }">
                    {{ formatDateTime(item.received_at) }}
                  </template>
                  <template #item.best_by_at="{ item }">
                    {{ formatDateTime(item.best_by_at) }}
                  </template>
                  <template #item.expires_at="{ item }">
                    {{ formatDateTime(item.expires_at) }}
                  </template>
                  <template #no-data>
                    <div class="text-center py-4 text-medium-emphasis">No yeast lots yet.</div>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-window-item>

          <!-- Other Tab -->
          <v-window-item value="other">
            <v-card variant="outlined">
              <v-card-title class="d-flex align-center">
                Other lots
                <v-spacer />
                <v-btn :loading="lotLoading" size="small" variant="text" @click="loadLots">
                  Refresh
                </v-btn>
                <v-btn
                  class="ml-2"
                  color="primary"
                  prepend-icon="mdi-plus"
                  size="small"
                  variant="text"
                  @click="openLotDialog('other')"
                >
                  Create lot
                </v-btn>
              </v-card-title>
              <v-card-text>
                <v-alert
                  v-if="lotErrorMessage"
                  class="mb-3"
                  density="compact"
                  type="error"
                  variant="tonal"
                >
                  {{ lotErrorMessage }}
                </v-alert>
                <v-data-table
                  class="data-table lot-table"
                  density="compact"
                  :headers="otherLotHeaders"
                  hover
                  item-value="uuid"
                  :items="otherLots"
                  :loading="lotLoading"
                  @click:row="(_event: Event, row: any) => openLotDetails(row.item.uuid)"
                >
                  <template #item.ingredient_uuid="{ item }">
                    {{ ingredientName(item.ingredient_uuid) }}
                  </template>
                  <template #item.category="{ item }">
                    <v-chip density="compact" size="small" variant="tonal">
                      {{ ingredientCategory(item.ingredient_uuid) }}
                    </v-chip>
                  </template>
                  <template #item.received_amount="{ item }">
                    {{ formatAmountPreferred(item.received_amount, item.received_unit) }}
                  </template>
                  <template #item.received_at="{ item }">
                    {{ formatDateTime(item.received_at) }}
                  </template>
                  <template #item.best_by_at="{ item }">
                    {{ formatDateTime(item.best_by_at) }}
                  </template>
                  <template #item.expires_at="{ item }">
                    {{ formatDateTime(item.expires_at) }}
                  </template>
                  <template #no-data>
                    <div class="text-center py-4 text-medium-emphasis">No other lots yet.</div>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-window-item>

          <!-- Usage Tab -->
          <v-window-item value="usage">
            <v-card variant="outlined">
              <v-card-title class="d-flex align-center">
                Usage log
                <v-spacer />
                <v-btn :loading="usageLoading" size="small" variant="text" @click="loadUsage">
                  Refresh
                </v-btn>
                <v-btn
                  class="ml-2"
                  color="primary"
                  prepend-icon="mdi-plus"
                  size="small"
                  variant="text"
                  @click="openUsageDialog"
                >
                  Log usage
                </v-btn>
              </v-card-title>
              <v-card-text>
                <v-alert
                  v-if="usageErrorMessage"
                  class="mb-3"
                  density="compact"
                  type="error"
                  variant="tonal"
                >
                  {{ usageErrorMessage }}
                </v-alert>
                <v-data-table
                  class="data-table"
                  density="compact"
                  :headers="usageHeaders"
                  item-value="uuid"
                  :items="usages"
                  :loading="usageLoading"
                >
                  <template #item.production_ref_uuid="{ item }">
                    {{ item.production_ref_uuid || 'n/a' }}
                  </template>
                  <template #item.used_at="{ item }">
                    {{ formatDateTime(item.used_at) }}
                  </template>
                  <template #item.notes="{ item }">
                    {{ item.notes || '' }}
                  </template>
                  <template #no-data>
                    <div class="text-center py-4 text-medium-emphasis">No usage records yet.</div>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-window-item>

          <!-- Received Tab -->
          <v-window-item value="received">
            <v-card variant="outlined">
              <v-card-title class="d-flex align-center">
                Inventory receipts
                <v-spacer />
                <v-btn :loading="receiptLoading" size="small" variant="text" @click="loadReceipts">
                  Refresh
                </v-btn>
                <v-btn
                  class="ml-2"
                  color="primary"
                  prepend-icon="mdi-plus"
                  size="small"
                  variant="text"
                  @click="openReceiptDialog"
                >
                  Create receipt
                </v-btn>
              </v-card-title>
              <v-card-text>
                <v-alert
                  v-if="receiptErrorMessage"
                  class="mb-3"
                  density="compact"
                  type="error"
                  variant="tonal"
                >
                  {{ receiptErrorMessage }}
                </v-alert>
                <v-data-table
                  class="data-table"
                  density="compact"
                  :headers="receiptHeaders"
                  item-value="uuid"
                  :items="receipts"
                  :loading="receiptLoading"
                >
                  <template #item.reference_code="{ item }">
                    {{ item.reference_code || 'n/a' }}
                  </template>
                  <template #item.supplier_uuid="{ item }">
                    {{ item.supplier_uuid || 'n/a' }}
                  </template>
                  <template #item.received_at="{ item }">
                    {{ formatDateTime(item.received_at) }}
                  </template>
                  <template #item.notes="{ item }">
                    {{ item.notes || '' }}
                  </template>
                  <template #no-data>
                    <div class="text-center py-4 text-medium-emphasis">No receipts yet.</div>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>

  <!-- Create Lot Dialog -->
  <v-dialog v-model="lotDialog" :max-width="$vuetify.display.xs ? '100%' : 600" persistent>
    <v-card>
      <v-card-title class="text-h6">Create ingredient lot</v-card-title>
      <v-card-text>
        <v-select
          v-model="lotForm.ingredient_uuid"
          density="comfortable"
          :items="filteredIngredientSelectItems"
          label="Ingredient"
          :rules="[rules.required]"
        />
        <v-select
          v-model="lotForm.receipt_uuid"
          clearable
          density="comfortable"
          :items="receiptSelectItems"
          label="Receipt (optional)"
        />
        <v-text-field
          v-model="lotForm.supplier_uuid"
          density="comfortable"
          label="Supplier UUID"
        />
        <v-text-field
          v-model="lotForm.brewery_lot_code"
          density="comfortable"
          label="Brewery lot code"
        />
        <v-text-field
          v-model="lotForm.originator_lot_code"
          density="comfortable"
          label="Originator lot code"
        />
        <v-text-field
          v-model="lotForm.originator_name"
          density="comfortable"
          label="Originator name"
        />
        <v-text-field
          v-model="lotForm.originator_type"
          density="comfortable"
          label="Originator type"
        />
        <v-text-field
          v-model="lotForm.received_at"
          density="comfortable"
          label="Received at"
          type="datetime-local"
        />
        <v-text-field
          v-model="lotForm.received_amount"
          density="comfortable"
          label="Received amount"
          :rules="[rules.required]"
          type="number"
        />
        <v-combobox
          v-model="lotForm.received_unit"
          density="comfortable"
          :items="unitOptions"
          label="Received unit"
          :rules="[rules.required]"
        />
        <v-text-field
          v-model="lotForm.best_by_at"
          density="comfortable"
          label="Best by"
          type="datetime-local"
        />
        <v-text-field
          v-model="lotForm.expires_at"
          density="comfortable"
          label="Expires at"
          type="datetime-local"
        />
        <v-textarea
          v-model="lotForm.notes"
          auto-grow
          density="comfortable"
          label="Notes"
          rows="2"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="lotSaving" variant="text" @click="closeLotDialog">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isLotFormValid"
          :loading="lotSaving"
          @click="createLot"
        >
          Create lot
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Create Usage Dialog -->
  <v-dialog v-model="usageDialog" :max-width="$vuetify.display.xs ? '100%' : 500" persistent>
    <v-card>
      <v-card-title class="text-h6">Log ingredient usage</v-card-title>
      <v-card-text>
        <v-text-field
          v-model="usageForm.production_ref_uuid"
          density="comfortable"
          label="Batch reference UUID"
        />
        <v-text-field
          v-model="usageForm.used_at"
          density="comfortable"
          label="Used at"
          type="datetime-local"
        />
        <v-textarea
          v-model="usageForm.notes"
          auto-grow
          density="comfortable"
          label="Notes"
          rows="2"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="usageSaving" variant="text" @click="closeUsageDialog">Cancel</v-btn>
        <v-btn
          color="primary"
          :loading="usageSaving"
          @click="createUsage"
        >
          Log usage
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Create Receipt Dialog -->
  <v-dialog v-model="receiptDialog" :max-width="$vuetify.display.xs ? '100%' : 500" persistent>
    <v-card>
      <v-card-title class="text-h6">Create receipt</v-card-title>
      <v-card-text>
        <v-text-field
          v-model="receiptForm.reference_code"
          density="comfortable"
          label="Reference code"
        />
        <v-text-field
          v-model="receiptForm.supplier_uuid"
          density="comfortable"
          label="Supplier UUID"
        />
        <v-text-field
          v-model="receiptForm.received_at"
          density="comfortable"
          label="Received at"
          type="datetime-local"
        />
        <v-textarea
          v-model="receiptForm.notes"
          auto-grow
          density="comfortable"
          label="Notes"
          rows="2"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="receiptSaving" variant="text" @click="closeReceiptDialog">Cancel</v-btn>
        <v-btn
          color="primary"
          :loading="receiptSaving"
          @click="createReceipt"
        >
          Create receipt
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Create Ingredient Dialog -->
  <v-dialog v-model="ingredientDialog" :max-width="$vuetify.display.xs ? '100%' : 500" persistent>
    <v-card>
      <v-card-title class="text-h6">Create ingredient</v-card-title>
      <v-card-text>
        <v-text-field
          v-model="ingredientForm.name"
          density="comfortable"
          label="Name"
          :rules="[rules.required]"
        />
        <v-select
          v-model="ingredientForm.category"
          density="comfortable"
          :items="ingredientCategoryOptions"
          label="Category"
          :rules="[rules.required]"
        />
        <v-combobox
          v-model="ingredientForm.default_unit"
          density="comfortable"
          :items="unitOptions"
          label="Default unit"
          :rules="[rules.required]"
        />
        <v-textarea
          v-model="ingredientForm.description"
          auto-grow
          density="comfortable"
          label="Description"
          rows="2"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="ingredientSaving" variant="text" @click="closeIngredientDialog">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isIngredientFormValid"
          :loading="ingredientSaving"
          @click="createIngredient"
        >
          Create ingredient
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  type Ingredient = {
    uuid: string
    name: string
    category: string
    default_unit: string
    description: string
    created_at: string
    updated_at: string
  }

  type InventoryReceipt = {
    uuid: string
    supplier_uuid: string
    reference_code: string
    received_at: string
    notes: string
    created_at: string
    updated_at: string
  }

  type IngredientLot = {
    uuid: string
    ingredient_uuid: string
    receipt_uuid: string | null
    received_amount: number
    received_unit: string
    best_by_at: string
    expires_at: string
    supplier_uuid: string
    brewery_lot_code: string
    originator_lot_code: string
    originator_name: string
    originator_type: string
    received_at: string
    notes: string
  }

  type InventoryUsage = {
    uuid: string
    production_ref_uuid: string
    used_at: string
    notes: string
    created_at: string
    updated_at: string
  }

  const { request, normalizeText, normalizeDateTime, toNumber, formatDateTime } = useInventoryApi()
  const { formatAmountPreferred } = useUnitPreferences()
  const router = useRouter()

  const activeTab = ref('malt')

  // Data
  const ingredients = ref<Ingredient[]>([])
  const receipts = ref<InventoryReceipt[]>([])
  const lots = ref<IngredientLot[]>([])
  const usages = ref<InventoryUsage[]>([])

  // Loading states
  const ingredientLoading = ref(false)
  const receiptLoading = ref(false)
  const lotLoading = ref(false)
  const usageLoading = ref(false)

  // Saving states
  const ingredientSaving = ref(false)
  const receiptSaving = ref(false)
  const lotSaving = ref(false)
  const usageSaving = ref(false)

  // Error messages
  const ingredientErrorMessage = ref('')
  const receiptErrorMessage = ref('')
  const lotErrorMessage = ref('')
  const usageErrorMessage = ref('')

  // Dialog states
  const ingredientDialog = ref(false)
  const receiptDialog = ref(false)
  const lotDialog = ref(false)
  const usageDialog = ref(false)

  // Category filter for lot creation
  const lotDialogCategory = ref<string>('fermentable')

  // Options
  const ingredientCategoryOptions = [
    { title: 'Fermentable', value: 'fermentable' },
    { title: 'Hop', value: 'hop' },
    { title: 'Yeast', value: 'yeast' },
    { title: 'Adjunct', value: 'adjunct' },
    { title: 'Salt', value: 'salt' },
    { title: 'Chemical', value: 'chemical' },
    { title: 'Gas', value: 'gas' },
    { title: 'Other', value: 'other' },
  ]
  const unitOptions = ['kg', 'g', 'lb', 'oz', 'l', 'ml', 'gal', 'bbl']

  // Forms
  const ingredientForm = reactive({
    name: '',
    category: '',
    default_unit: '',
    description: '',
  })

  const receiptForm = reactive({
    supplier_uuid: '',
    reference_code: '',
    received_at: '',
    notes: '',
  })

  const lotForm = reactive({
    ingredient_uuid: null as string | null,
    receipt_uuid: null as string | null,
    supplier_uuid: '',
    brewery_lot_code: '',
    originator_lot_code: '',
    originator_name: '',
    originator_type: '',
    received_at: '',
    received_amount: '',
    received_unit: '',
    best_by_at: '',
    expires_at: '',
    notes: '',
  })

  const usageForm = reactive({
    production_ref_uuid: '',
    used_at: '',
    notes: '',
  })

  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
  })

  const rules = {
    required: (v: string | null) => (v !== null && v !== '' && String(v).trim() !== '') || 'Required',
  }

  // Table headers
  const lotHeaders = [
    { title: 'Lot Code', key: 'brewery_lot_code', sortable: true },
    { title: 'Ingredient', key: 'ingredient_uuid', sortable: true },
    { title: 'Received Amount', key: 'received_amount', sortable: true },
    { title: 'Received Date', key: 'received_at', sortable: true },
    { title: 'Best By', key: 'best_by_at', sortable: true },
    { title: 'Expires', key: 'expires_at', sortable: true },
    { title: 'Notes', key: 'notes', sortable: false },
  ]

  const otherLotHeaders = [
    { title: 'Lot Code', key: 'brewery_lot_code', sortable: true },
    { title: 'Ingredient', key: 'ingredient_uuid', sortable: true },
    { title: 'Category', key: 'category', sortable: true },
    { title: 'Received Amount', key: 'received_amount', sortable: true },
    { title: 'Received Date', key: 'received_at', sortable: true },
    { title: 'Best By', key: 'best_by_at', sortable: true },
    { title: 'Expires', key: 'expires_at', sortable: true },
    { title: 'Notes', key: 'notes', sortable: false },
  ]

  const usageHeaders = [
    { title: 'Batch Reference', key: 'production_ref_uuid', sortable: true },
    { title: 'Used At', key: 'used_at', sortable: true },
    { title: 'Notes', key: 'notes', sortable: false },
  ]

  const receiptHeaders = [
    { title: 'Reference Code', key: 'reference_code', sortable: true },
    { title: 'Supplier', key: 'supplier_uuid', sortable: true },
    { title: 'Received At', key: 'received_at', sortable: true },
    { title: 'Notes', key: 'notes', sortable: false },
  ]

  // Computed: Filter lots by category
  const maltLots = computed(() => {
    return lots.value.filter(lot => {
      const ingredient = ingredients.value.find(i => i.uuid === lot.ingredient_uuid)
      return ingredient?.category === 'fermentable'
    })
  })

  const hopLots = computed(() => {
    return lots.value.filter(lot => {
      const ingredient = ingredients.value.find(i => i.uuid === lot.ingredient_uuid)
      return ingredient?.category === 'hop'
    })
  })

  const yeastLots = computed(() => {
    return lots.value.filter(lot => {
      const ingredient = ingredients.value.find(i => i.uuid === lot.ingredient_uuid)
      return ingredient?.category === 'yeast'
    })
  })

  const otherCategories = ['adjunct', 'salt', 'chemical', 'gas', 'other']

  const otherLots = computed(() => {
    return lots.value.filter(lot => {
      const ingredient = ingredients.value.find(i => i.uuid === lot.ingredient_uuid)
      return ingredient && otherCategories.includes(ingredient.category)
    })
  })

  // Computed: Select items
  const filteredIngredientSelectItems = computed(() => {
    // For "other" tab, include all non-core categories
    const categoriesToInclude = lotDialogCategory.value === 'other'
      ? otherCategories
      : [lotDialogCategory.value]

    return ingredients.value
      .filter(ingredient => categoriesToInclude.includes(ingredient.category))
      .map(ingredient => ({
        title: ingredient.name,
        value: ingredient.uuid,
      }))
  })

  const receiptSelectItems = computed(() =>
    receipts.value.map(receipt => ({
      title: receipt.reference_code || 'Unknown Receipt',
      value: receipt.uuid,
    })),
  )

  // Computed: Form validation
  const isLotFormValid = computed(() => {
    return lotForm.ingredient_uuid && lotForm.received_amount && lotForm.received_unit
  })

  const isIngredientFormValid = computed(() => {
    return ingredientForm.name.trim() && ingredientForm.category && ingredientForm.default_unit
  })

  onMounted(async () => {
    await refreshAll()
  })

  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  async function refreshAll () {
    await Promise.allSettled([loadIngredients(), loadReceipts(), loadLots(), loadUsage()])
  }

  async function loadIngredients () {
    ingredientLoading.value = true
    ingredientErrorMessage.value = ''
    try {
      ingredients.value = await request<Ingredient[]>('/ingredients')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load ingredients'
      ingredientErrorMessage.value = message
    } finally {
      ingredientLoading.value = false
    }
  }

  async function loadReceipts () {
    receiptLoading.value = true
    receiptErrorMessage.value = ''
    try {
      receipts.value = await request<InventoryReceipt[]>('/inventory-receipts')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load receipts'
      receiptErrorMessage.value = message
    } finally {
      receiptLoading.value = false
    }
  }

  async function loadLots () {
    lotLoading.value = true
    lotErrorMessage.value = ''
    try {
      lots.value = await request<IngredientLot[]>('/ingredient-lots')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load lots'
      lotErrorMessage.value = message
    } finally {
      lotLoading.value = false
    }
  }

  async function loadUsage () {
    usageLoading.value = true
    usageErrorMessage.value = ''
    try {
      usages.value = await request<InventoryUsage[]>('/inventory-usage')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load usage'
      usageErrorMessage.value = message
    } finally {
      usageLoading.value = false
    }
  }

  // Dialog openers
  function openIngredientDialog () {
    ingredientForm.name = ''
    ingredientForm.category = ''
    ingredientForm.default_unit = ''
    ingredientForm.description = ''
    ingredientDialog.value = true
  }

  function closeIngredientDialog () {
    ingredientDialog.value = false
  }

  function openReceiptDialog () {
    receiptForm.supplier_uuid = ''
    receiptForm.reference_code = ''
    receiptForm.received_at = ''
    receiptForm.notes = ''
    receiptDialog.value = true
  }

  function closeReceiptDialog () {
    receiptDialog.value = false
  }

  function openLotDialog (category: string) {
    lotDialogCategory.value = category
    lotForm.ingredient_uuid = null
    lotForm.receipt_uuid = null
    lotForm.supplier_uuid = ''
    lotForm.brewery_lot_code = ''
    lotForm.originator_lot_code = ''
    lotForm.originator_name = ''
    lotForm.originator_type = ''
    lotForm.received_at = ''
    lotForm.received_amount = ''
    lotForm.received_unit = ''
    lotForm.best_by_at = ''
    lotForm.expires_at = ''
    lotForm.notes = ''
    lotDialog.value = true
  }

  function closeLotDialog () {
    lotDialog.value = false
  }

  function openUsageDialog () {
    usageForm.production_ref_uuid = ''
    usageForm.used_at = ''
    usageForm.notes = ''
    usageDialog.value = true
  }

  function closeUsageDialog () {
    usageDialog.value = false
  }

  // Create functions
  async function createIngredient () {
    if (!isIngredientFormValid.value) {
      return
    }

    ingredientSaving.value = true
    ingredientErrorMessage.value = ''

    try {
      const payload = {
        name: ingredientForm.name.trim(),
        category: ingredientForm.category,
        default_unit: ingredientForm.default_unit,
        description: normalizeText(ingredientForm.description),
      }
      await request<Ingredient>('/ingredients', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      closeIngredientDialog()
      await loadIngredients()
      showNotice('Ingredient created')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to create ingredient'
      ingredientErrorMessage.value = message
      showNotice(message, 'error')
    } finally {
      ingredientSaving.value = false
    }
  }

  async function createReceipt () {
    receiptSaving.value = true
    receiptErrorMessage.value = ''

    try {
      const payload = {
        supplier_uuid: normalizeText(receiptForm.supplier_uuid),
        reference_code: normalizeText(receiptForm.reference_code),
        received_at: normalizeDateTime(receiptForm.received_at),
        notes: normalizeText(receiptForm.notes),
      }
      await request<InventoryReceipt>('/inventory-receipts', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      closeReceiptDialog()
      await loadReceipts()
      showNotice('Receipt created')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to create receipt'
      receiptErrorMessage.value = message
      showNotice(message, 'error')
    } finally {
      receiptSaving.value = false
    }
  }

  async function createLot () {
    if (!isLotFormValid.value) {
      return
    }

    lotSaving.value = true
    lotErrorMessage.value = ''

    try {
      const payload = {
        ingredient_uuid: lotForm.ingredient_uuid,
        receipt_uuid: lotForm.receipt_uuid,
        supplier_uuid: normalizeText(lotForm.supplier_uuid),
        brewery_lot_code: normalizeText(lotForm.brewery_lot_code),
        originator_lot_code: normalizeText(lotForm.originator_lot_code),
        originator_name: normalizeText(lotForm.originator_name),
        originator_type: normalizeText(lotForm.originator_type),
        received_at: normalizeDateTime(lotForm.received_at),
        received_amount: toNumber(lotForm.received_amount),
        received_unit: lotForm.received_unit,
        best_by_at: normalizeDateTime(lotForm.best_by_at),
        expires_at: normalizeDateTime(lotForm.expires_at),
        notes: normalizeText(lotForm.notes),
      }
      await request<IngredientLot>('/ingredient-lots', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      closeLotDialog()
      await loadLots()
      showNotice('Lot created')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to create lot'
      lotErrorMessage.value = message
      showNotice(message, 'error')
    } finally {
      lotSaving.value = false
    }
  }

  async function createUsage () {
    usageSaving.value = true
    usageErrorMessage.value = ''

    try {
      const payload = {
        production_ref_uuid: normalizeText(usageForm.production_ref_uuid),
        used_at: normalizeDateTime(usageForm.used_at),
        notes: normalizeText(usageForm.notes),
      }
      await request<InventoryUsage>('/inventory-usage', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      closeUsageDialog()
      await loadUsage()
      showNotice('Usage logged')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to log usage'
      usageErrorMessage.value = message
      showNotice(message, 'error')
    } finally {
      usageSaving.value = false
    }
  }

  function ingredientName (ingredientUuid: string) {
    return ingredients.value.find(ingredient => ingredient.uuid === ingredientUuid)?.name ?? 'Unknown Ingredient'
  }

  function ingredientCategory (ingredientUuid: string) {
    const category = ingredients.value.find(ingredient => ingredient.uuid === ingredientUuid)?.category
    // Return a display-friendly label
    const labels: Record<string, string> = {
      adjunct: 'Adjunct',
      salt: 'Salt',
      chemical: 'Chemical',
      gas: 'Gas',
      other: 'Other',
    }
    return labels[category ?? ''] ?? category ?? 'Unknown'
  }

  function openLotDetails (lotUuid: string) {
    router.push({
      path: '/inventory/lot-details',
      query: { lot_uuid: lotUuid },
    })
  }
</script>

<style scoped>
.inventory-page {
  position: relative;
}

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.card-title-responsive {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.data-table {
  overflow-x: auto;
}

.data-table :deep(.v-table__wrapper) {
  overflow-x: auto;
}

.data-table :deep(th) {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: rgba(var(--v-theme-on-surface), 0.55);
  white-space: nowrap;
}

.data-table :deep(td) {
  font-size: 0.85rem;
}

.lot-table :deep(tr) {
  cursor: pointer;
}
</style>
