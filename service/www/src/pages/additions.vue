<template>
  <v-container class="simple-page" fluid>
    <v-card class="section-card">
      <v-card-title>Additions</v-card-title>
      <v-card-text>
        <div class="text-body-1">Addition logging is moving here next.</div>
        <div class="text-body-2 text-medium-emphasis">
          Until then, log ingredient and gas additions from the Batches workflow.
        </div>
        <v-btn class="mt-3" color="primary" to="/batches">Log an addition</v-btn>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<style scoped>
.simple-page {
  position: relative;
}

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}
</style>
