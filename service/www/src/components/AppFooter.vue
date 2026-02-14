<template>
  <v-footer app class="app-footer" height="56">
    <v-container class="footer-content">
      <div class="text-caption text-medium-emphasis d-none d-sm-block">
        BrewPipes Production UI
      </div>
      <v-spacer class="d-none d-sm-block" />
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
          :href="apiUrl"
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
  import { computed } from 'vue'

  const apiUrl = computed(() => `${window.location.origin}/api`)
</script>

<style scoped>
.app-footer {
  border-top: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.8);
  backdrop-filter: blur(10px);
}

.footer-content {
  display: flex;
  align-items: center;
}

@media (max-width: 599px) {
  .footer-content {
    justify-content: center;
  }
}
</style>
