<template>
  <v-container class="procurement-page" fluid>
    <v-row align="stretch">
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Suppliers</v-card-title>
          <v-card-text>
            Manage vendors, contacts, and address details.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/procurement/suppliers">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Purchase orders</v-card-title>
          <v-card-text>
            Track orders, status, and expected arrivals.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/procurement/purchase-orders">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Order lines</v-card-title>
          <v-card-text>
            Capture line items, quantities, and unit costs.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/procurement/purchase-order-lines">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Order fees</v-card-title>
          <v-card-text>
            Record shipping, freight, and other fees.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/procurement/purchase-order-fees">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.procurement-page {
  position: relative;
}

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}
</style>
