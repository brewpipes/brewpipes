<template>
  <v-container class="simple-page" fluid>
    <v-card class="section-card">
      <v-card-title>Vessels</v-card-title>
      <v-card-text>
        <div class="text-body-1">Vessel management is being built.</div>
        <div class="text-body-2 text-medium-emphasis">
          For now, add vessels from the Batches workflow to start tracking capacity.
        </div>
        <v-btn class="mt-3" color="primary" to="/batches">Go to batches</v-btn>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<style scoped>
.simple-page {
  position: relative;
}

.section-card {
  background: rgba(253, 251, 247, 0.92);
  border: 1px solid rgba(28, 26, 22, 0.08);
  box-shadow: 0 10px 24px rgba(28, 26, 22, 0.05);
}
</style>
