<template>
  <v-container class="inventory-page" fluid>
    <v-row align="stretch">
      <v-col cols="12" md="6" lg="4">
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
      <v-col cols="12" md="6" lg="4">
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
      <v-col cols="12" md="6" lg="4">
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
      <v-col cols="12" md="6" lg="4">
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
      <v-col cols="12" md="6" lg="4">
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
  </v-container>
</template>

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
