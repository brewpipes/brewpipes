<template>
  <v-app class="app-shell">
    <router-view />
  </v-app>
</template>

<script lang="ts" setup>
  //
</script>

<style>
:root {
  --bp-ink: rgb(var(--v-theme-on-background));
  --bp-paper: rgb(var(--v-theme-background));
  --bp-cream: rgb(var(--v-theme-surface));
  --bp-amber: rgb(var(--v-theme-secondary));
  --bp-moss: rgb(var(--v-theme-primary));
  --bp-mist: rgba(var(--v-theme-on-background), 0.08);
}

body {
  background-color: var(--bp-paper);
  color: var(--bp-ink);
}

.app-shell {
  background: var(--bp-paper);
  min-height: 100vh;
  position: relative;
  overflow: hidden;
}

.app-shell::before,
.app-shell::after {
  content: '';
  position: fixed;
  inset: -10% 0 auto 0;
  height: 55vh;
  pointer-events: none;
  z-index: 0;
}

.app-shell::before {
  background:
    radial-gradient(600px 300px at 12% 20%, rgba(var(--v-theme-primary), 0.2), transparent 70%),
    radial-gradient(520px 260px at 88% 10%, rgba(var(--v-theme-secondary), 0.18), transparent 70%),
    linear-gradient(120deg, rgba(var(--v-theme-surface), 0.65), rgba(var(--v-theme-background), 0.95));
}

.app-shell::after {
  top: auto;
  bottom: -20vh;
  height: 60vh;
  background:
    radial-gradient(520px 320px at 20% 80%, rgba(var(--v-theme-info), 0.16), transparent 70%),
    radial-gradient(620px 360px at 80% 70%, rgba(var(--v-theme-success), 0.16), transparent 70%);
}
</style>
