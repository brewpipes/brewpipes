<template>
  <v-container class="simple-page" fluid>
    <v-card class="section-card">
      <v-card-title>Transfers</v-card-title>
      <v-card-text>
        <div class="text-body-1">Transfer tracking is coming soon.</div>
        <div class="text-body-2 text-medium-emphasis">
          Use the Batches workflow to record transfers while this view is finalized.
        </div>
        <v-btn class="mt-3" color="primary" to="/batches">Record a transfer</v-btn>
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
