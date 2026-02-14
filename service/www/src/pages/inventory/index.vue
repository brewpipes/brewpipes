<template>
  <v-container class="inventory-page" fluid>
    <v-row align="stretch">
      <v-col cols="12" lg="4" md="6">
        <v-card class="section-card">
          <v-card-title>
            <v-icon class="mr-2" color="success" icon="mdi-package-down" />
            Receive Inventory
          </v-card-title>
          <v-card-text>
            Record incoming inventory without a purchase order.
          </v-card-text>
          <v-card-actions>
            <v-btn color="success" @click="receiveDialogOpen = true">Receive</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" lg="4" md="6">
        <v-card class="section-card">
          <v-card-title>Activity</v-card-title>
          <v-card-text>
            Review inventory movement history and corrections.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/activity">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" lg="4" md="6">
        <v-card class="section-card">
          <v-card-title>Product</v-card-title>
          <v-card-text>
            Track packaged product lots tied to batches.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/product">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" lg="4" md="6">
        <v-card class="section-card">
          <v-card-title>Ingredients</v-card-title>
          <v-card-text>
            Manage ingredient stock, usage, and types.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/ingredients">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" lg="4" md="6">
        <v-card class="section-card">
          <v-card-title>Adjustments & Transfers</v-card-title>
          <v-card-text>
            Capture corrections and stock moves between locations.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/adjustments-transfers">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" lg="4" md="6">
        <v-card class="section-card">
          <v-card-title>Locations</v-card-title>
          <v-card-text>
            Configure storage locations for inventory tracking.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/locations">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- Receive Without PO Dialog -->
    <ReceiveWithoutPODialog
      v-model="receiveDialogOpen"
      @received="handleReceived"
    />
  </v-container>
</template>

<script lang="ts" setup>
  import { ref } from 'vue'
  import ReceiveWithoutPODialog from '@/components/procurement/ReceiveWithoutPODialog.vue'
  import { useSnackbar } from '@/composables/useSnackbar'

  const { showNotice } = useSnackbar()

  const receiveDialogOpen = ref(false)

  function handleReceived () {
    showNotice('Inventory received successfully', 'success')
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
</style>
