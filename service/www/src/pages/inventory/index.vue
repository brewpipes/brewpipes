<template>
  <v-container class="inventory-page" fluid>
    <v-row align="stretch">
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Ingredients</v-card-title>
          <v-card-text>
            Manage ingredient definitions and default units.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/ingredients">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Ingredient details</v-card-title>
          <v-card-text>
            Track malt, hop, and yeast attributes per ingredient.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/ingredient-details">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Stock locations</v-card-title>
          <v-card-text>
            Configure storage locations for inventory tracking.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/locations">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Receipts</v-card-title>
          <v-card-text>
            Log incoming deliveries and supplier receipts.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/receipts">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Ingredient lots</v-card-title>
          <v-card-text>
            Track received lots, quantities, and originator data.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/lots">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Lot details</v-card-title>
          <v-card-text>
            Capture malt moisture, hop alpha, and yeast viability.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/lot-details">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Movements</v-card-title>
          <v-card-text>
            Record inventory inflow, usage, and adjustments.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/movements">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Usage</v-card-title>
          <v-card-text>
            Tie ingredient usage to production references.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/usage">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Adjustments</v-card-title>
          <v-card-text>
            Capture corrections, shrink, or cycle count updates.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/adjustments">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Transfers</v-card-title>
          <v-card-text>
            Move inventory between storage locations.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/transfers">Open</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" lg="4">
        <v-card class="section-card">
          <v-card-title>Beer lots</v-card-title>
          <v-card-text>
            Track finished product lots tied to production batches.
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" to="/inventory/beer-lots">Open</v-btn>
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
