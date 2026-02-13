<template>
  <v-container class="procurement-page" fluid>
    <v-card class="section-card">
      <v-card-title class="card-title-responsive">
        <span>Suppliers</span>
        <div class="card-title-actions">
          <v-btn
            :icon="$vuetify.display.xs"
            :loading="loading"
            size="small"
            variant="text"
            @click="loadSuppliers"
          >
            <v-icon v-if="$vuetify.display.xs" icon="mdi-refresh" />
            <span v-else>Refresh</span>
          </v-btn>
          <v-btn
            color="primary"
            :icon="$vuetify.display.xs"
            size="small"
            variant="text"
            @click="openCreateDialog"
          >
            <v-icon v-if="$vuetify.display.xs" icon="mdi-plus" />
            <span v-else>New supplier</span>
          </v-btn>
        </div>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12">
            <v-card class="sub-card" variant="outlined">
              <v-card-title>Supplier list</v-card-title>
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
                <v-table class="data-table" density="compact">
                  <thead>
                    <tr>
                      <th>Name</th>
                      <th>Contact</th>
                      <th>Email</th>
                      <th>Updated</th>
                      <th class="text-right">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="supplier in suppliers" :key="supplier.uuid">
                      <td>{{ supplier.name }}</td>
                      <td>{{ supplier.contact_name || 'n/a' }}</td>
                      <td>{{ supplier.email || 'n/a' }}</td>
                      <td>{{ formatDateTime(supplier.updated_at) }}</td>
                      <td class="text-right">
                        <v-btn
                          icon="mdi-pencil"
                          size="x-small"
                          variant="text"
                          @click="openEditDialog(supplier)"
                        />
                      </td>
                    </tr>
                    <tr v-if="suppliers.length === 0">
                      <td colspan="5">No suppliers yet.</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>

  <!-- Create/Edit Supplier Dialog -->
  <v-dialog v-model="supplierDialog" :max-width="$vuetify.display.xs ? '100%' : 720" persistent>
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit supplier' : 'Create supplier' }}
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field v-model="supplierForm.name" label="Name" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="supplierForm.contact_name" label="Contact name" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="supplierForm.email" label="Email" type="email" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="supplierForm.phone" label="Phone" type="tel" />
          </v-col>
          <v-col cols="12">
            <v-text-field v-model="supplierForm.address_line1" label="Address line 1" />
          </v-col>
          <v-col cols="12">
            <v-text-field v-model="supplierForm.address_line2" label="Address line 2" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="supplierForm.city" label="City" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="supplierForm.region" label="Region" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="supplierForm.postal_code" label="Postal code" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="supplierForm.country" label="Country" />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="closeSupplierDialog">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!supplierForm.name.trim()"
          :loading="saving"
          @click="saveSupplier"
        >
          {{ isEditing ? 'Save changes' : 'Add supplier' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref } from 'vue'
  import { type Supplier, useProcurementApi } from '@/composables/useProcurementApi'

  const {
    getSuppliers,
    createSupplier,
    updateSupplier,
    normalizeText,
    formatDateTime,
  } = useProcurementApi()

  // State
  const suppliers = ref<Supplier[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const errorMessage = ref('')

  // Dialog state
  const supplierDialog = ref(false)
  const editingSupplierUuid = ref<string | null>(null)

  const supplierForm = reactive({
    name: '',
    contact_name: '',
    email: '',
    phone: '',
    address_line1: '',
    address_line2: '',
    city: '',
    region: '',
    postal_code: '',
    country: '',
  })

  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
  })

  // Computed
  const isEditing = computed(() => editingSupplierUuid.value !== null)

  // Lifecycle
  onMounted(async () => {
    await loadSuppliers()
  })

  // Methods
  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  function resetForm () {
    supplierForm.name = ''
    supplierForm.contact_name = ''
    supplierForm.email = ''
    supplierForm.phone = ''
    supplierForm.address_line1 = ''
    supplierForm.address_line2 = ''
    supplierForm.city = ''
    supplierForm.region = ''
    supplierForm.postal_code = ''
    supplierForm.country = ''
  }

  function openCreateDialog () {
    editingSupplierUuid.value = null
    resetForm()
    supplierDialog.value = true
  }

  function openEditDialog (supplier: Supplier) {
    editingSupplierUuid.value = supplier.uuid
    supplierForm.name = supplier.name
    supplierForm.contact_name = supplier.contact_name ?? ''
    supplierForm.email = supplier.email ?? ''
    supplierForm.phone = supplier.phone ?? ''
    supplierForm.address_line1 = supplier.address_line1 ?? ''
    supplierForm.address_line2 = supplier.address_line2 ?? ''
    supplierForm.city = supplier.city ?? ''
    supplierForm.region = supplier.region ?? ''
    supplierForm.postal_code = supplier.postal_code ?? ''
    supplierForm.country = supplier.country ?? ''
    supplierDialog.value = true
  }

  function closeSupplierDialog () {
    supplierDialog.value = false
    editingSupplierUuid.value = null
  }

  async function loadSuppliers () {
    loading.value = true
    errorMessage.value = ''
    try {
      suppliers.value = await getSuppliers()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load suppliers'
      errorMessage.value = message
    } finally {
      loading.value = false
    }
  }

  async function saveSupplier () {
    if (!supplierForm.name.trim()) {
      return
    }

    saving.value = true
    errorMessage.value = ''

    try {
      const payload = {
        name: supplierForm.name.trim(),
        contact_name: normalizeText(supplierForm.contact_name),
        email: normalizeText(supplierForm.email),
        phone: normalizeText(supplierForm.phone),
        address_line1: normalizeText(supplierForm.address_line1),
        address_line2: normalizeText(supplierForm.address_line2),
        city: normalizeText(supplierForm.city),
        region: normalizeText(supplierForm.region),
        postal_code: normalizeText(supplierForm.postal_code),
        country: normalizeText(supplierForm.country),
      }

      if (isEditing.value && editingSupplierUuid.value) {
        await updateSupplier(editingSupplierUuid.value, payload)
        showNotice('Supplier updated')
      } else {
        await createSupplier(payload)
        showNotice('Supplier created')
      }

      closeSupplierDialog()
      await loadSuppliers()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to save supplier'
      errorMessage.value = message
      showNotice(message, 'error')
    } finally {
      saving.value = false
    }
  }
</script>

<style scoped>
.procurement-page {
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

.card-title-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
}

.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
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
</style>
