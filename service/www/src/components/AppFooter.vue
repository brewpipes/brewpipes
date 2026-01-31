<template>
  <v-footer app class="app-footer" height="56">
    <v-container class="d-flex align-center">
      <div class="text-caption text-medium-emphasis">
        BrewPipes Production UI
      </div>
      <v-spacer />
      <div class="d-flex align-center ga-2">
        <v-btn
          class="text-caption"
          href="https://github.com/brewpipes/brewpipes"
          rel="noopener noreferrer"
          size="small"
          target="_blank"
          variant="text"
        >
          Repo
        </v-btn>
        <v-btn
          class="text-caption"
          href="http://localhost:8080"
          rel="noopener noreferrer"
          size="small"
          target="_blank"
          variant="text"
        >
          API
        </v-btn>
      </div>
    </v-container>
  </v-footer>
</template>

<script setup lang="ts">
  //
</script>

<style scoped>
.app-footer {
  border-top: 1px solid rgba(28, 26, 22, 0.08);
  background: rgba(253, 251, 247, 0.8);
  backdrop-filter: blur(10px);
}
</style>
