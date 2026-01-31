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
  --bp-ink: #1c1a16;
  --bp-paper: #f5f2ea;
  --bp-cream: #fdfbf7;
  --bp-amber: #c4753c;
  --bp-moss: #2f5d50;
  --bp-mist: rgba(24, 22, 18, 0.06);
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
    radial-gradient(600px 300px at 12% 20%, rgba(47, 93, 80, 0.18), transparent 70%),
    radial-gradient(520px 260px at 88% 10%, rgba(196, 117, 60, 0.18), transparent 70%),
    linear-gradient(120deg, rgba(255, 255, 255, 0.5), rgba(245, 242, 234, 0.9));
}

.app-shell::after {
  top: auto;
  bottom: -20vh;
  height: 60vh;
  background:
    radial-gradient(520px 320px at 20% 80%, rgba(45, 95, 122, 0.12), transparent 70%),
    radial-gradient(620px 360px at 80% 70%, rgba(47, 122, 74, 0.12), transparent 70%);
}
</style>
