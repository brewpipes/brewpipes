<template>
  <v-container class="procurement-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Suppliers
        <v-spacer />
        <v-btn size="small" variant="text" :loading="loading" @click="loadSuppliers">
          Refresh
        </v-btn>
        <v-btn class="ml-2" color="primary" size="small" variant="text" @click="createSupplierDialog = true">
          New supplier
        </v-btn>
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
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="supplier in suppliers" :key="supplier.id">
                      <td>{{ supplier.name }}</td>
                      <td>{{ supplier.contact_name || 'n/a' }}</td>
                      <td>{{ supplier.email || 'n/a' }}</td>
                      <td>{{ formatDateTime(supplier.updated_at) }}</td>
                    </tr>
                    <tr v-if="suppliers.length === 0">
                      <td colspan="4">No suppliers yet.</td>
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

  <v-dialog v-model="createSupplierDialog" max-width="720">
    <v-card>
      <v-card-title class="text-h6">Create supplier</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field v-model="supplierForm.name" label="Name" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="supplierForm.contact_name" label="Contact name" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="supplierForm.email" label="Email" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="supplierForm.phone" label="Phone" />
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
        <v-btn variant="text" @click="createSupplierDialog = false">Cancel</v-btn>
        <v-btn color="primary" :disabled="!supplierForm.name.trim()" @click="createSupplier">
          Add supplier
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue'
import { useProcurementApi } from '@/composables/useProcurementApi'

type Supplier = {
  id: number
  uuid: string
  name: string
  contact_name: string | null
  email: string | null
  phone: string | null
  address_line1: string | null
  address_line2: string | null
  city: string | null
  region: string | null
  postal_code: string | null
  country: string | null
  created_at: string
  updated_at: string
}

const { request, normalizeText, formatDateTime } = useProcurementApi()

const suppliers = ref<Supplier[]>([])
const loading = ref(false)
const errorMessage = ref('')
const createSupplierDialog = ref(false)

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

onMounted(async () => {
  await loadSuppliers()
})

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

async function loadSuppliers() {
  loading.value = true
  errorMessage.value = ''
  try {
    suppliers.value = await request<Supplier[]>('/suppliers')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load suppliers'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}

async function createSupplier() {
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
    await request<Supplier>('/suppliers', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
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
    await loadSuppliers()
    createSupplierDialog.value = false
    showNotice('Supplier created')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create supplier'
    errorMessage.value = message
    showNotice(message, 'error')
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

.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.data-table :deep(th) {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: rgba(var(--v-theme-on-surface), 0.55);
}

.data-table :deep(td) {
  font-size: 0.85rem;
}
</style>
