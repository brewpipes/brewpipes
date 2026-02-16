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
  <SupplierCreateEditDialog
    v-model="supplierDialog"
    :edit-supplier="editingSupplier"
    :saving="saving"
    @submit="saveSupplier"
  />
</template>

<script lang="ts" setup>
  import type { Supplier } from '@/types'
  import { onMounted, ref } from 'vue'
  import SupplierCreateEditDialog from '@/components/procurement/SupplierCreateEditDialog.vue'
  import type { SupplierFormData } from '@/components/procurement/SupplierCreateEditDialog.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { formatDateTime } from '@/composables/useFormatters'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { useSnackbar } from '@/composables/useSnackbar'

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

  const { execute: executeLoad, loading, error: errorMessage } = useAsyncAction()
  const { execute: executeSave, loading: saving, error: saveError } = useAsyncAction()

  // Dialog state
  const supplierDialog = ref(false)
  const editingSupplier = ref<Supplier | null>(null)

  // Lifecycle
  onMounted(async () => {
    await loadSuppliers()
  })

  // Methods
  function openCreateDialog () {
    editingSupplier.value = null
    supplierDialog.value = true
  }

  function openEditDialog (supplier: Supplier) {
    editingSupplier.value = supplier
    supplierDialog.value = true
  }

  async function loadSuppliers () {
    await executeLoad(async () => {
      suppliers.value = await getSuppliers()
    })
  }

  async function saveSupplier (data: SupplierFormData) {
    await executeSave(async () => {
      if (editingSupplier.value) {
        await updateSupplier(editingSupplier.value.uuid, data)
        showNotice('Supplier updated')
      } else {
        await createSupplier(data)
        showNotice('Supplier created')
      }

      supplierDialog.value = false
      editingSupplier.value = null
      await loadSuppliers()
    })
    if (saveError.value) {
      errorMessage.value = saveError.value
      showNotice(saveError.value, 'error')
    }
  }
</script>

<style scoped>
.procurement-page {
  position: relative;
}
</style>
