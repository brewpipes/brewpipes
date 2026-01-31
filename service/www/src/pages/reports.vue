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
  background: rgba(253, 251, 247, 0.92);
  border: 1px solid rgba(28, 26, 22, 0.08);
  box-shadow: 0 10px 24px rgba(28, 26, 22, 0.05);
}
</style>
