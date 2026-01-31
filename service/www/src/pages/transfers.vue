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
  background: rgba(253, 251, 247, 0.92);
  border: 1px solid rgba(28, 26, 22, 0.08);
  box-shadow: 0 10px 24px rgba(28, 26, 22, 0.05);
}
</style>
