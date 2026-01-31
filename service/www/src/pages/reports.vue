<template>
  <v-container class="simple-page" fluid>
    <v-card class="section-card">
      <v-card-title>Reports</v-card-title>
      <v-card-text>
        <div class="text-body-1">Reporting is on deck.</div>
        <div class="text-body-2 text-medium-emphasis">
          This view will summarize batch progress and quality checks.
        </div>
        <v-btn class="mt-3" color="primary" to="/batches">View batches</v-btn>
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
