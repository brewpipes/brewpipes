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
                <v-data-table
                  class="data-table"
                  density="compact"
                  :headers="headers"
                  item-value="uuid"
                  :items="suppliers"
                  :loading="loading"
                >
                  <template #item.contact_name="{ item }">
                    {{ item.contact_name || 'n/a' }}
                  </template>
                  <template #item.email="{ item }">
                    {{ item.email || 'n/a' }}
                  </template>
                  <template #item.updated_at="{ item }">
                    {{ formatDateTime(item.updated_at) }}
                  </template>
                  <template #item.actions="{ item }">
                    <v-btn
                      icon="mdi-pencil"
                      size="x-small"
                      variant="text"
                      @click.stop="openEditDialog(item)"
                    />
                  </template>
                  <template #no-data>
                    <div class="text-center py-4 text-medium-emphasis">No suppliers yet.</div>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-container>

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
  import type { Supplier } from '@/types'
  import { computed, onMounted, reactive, ref } from 'vue'
  import { formatDateTime } from '@/composables/useFormatters'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { normalizeText } from '@/utils/normalize'

  const {
    getSuppliers,
    createSupplier,
    updateSupplier,
  } = useProcurementApi()
  const { showNotice } = useSnackbar()

  // Table configuration
  const headers = [
    { title: 'Name', key: 'name', sortable: true },
    { title: 'Contact', key: 'contact_name', sortable: true },
    { title: 'Email', key: 'email', sortable: true },
    { title: 'Updated', key: 'updated_at', sortable: true },
    { title: '', key: 'actions', sortable: false, align: 'end' as const, width: '80px' },
  ]

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

  // Computed
  const isEditing = computed(() => editingSupplierUuid.value !== null)

  // Lifecycle
  onMounted(async () => {
    await loadSuppliers()
  })

  // Methods
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
</style>
