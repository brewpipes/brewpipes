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

/* Global mobile responsiveness improvements */

/* Ensure minimum touch target size (44x44px) for interactive elements on mobile */
@media (max-width: 599px) {
  /* Increase tap targets for list items */
  .v-list-item {
    min-height: 48px;
  }

  /* Ensure buttons have adequate touch targets */
  .v-btn--size-small {
    min-width: 44px;
    min-height: 44px;
  }

  /* Improve chip tap targets */
  .v-chip {
    min-height: 32px;
  }

  /* Reduce card padding on mobile for more content space */
  .v-card-text {
    padding: 12px;
  }

  .v-card-title {
    padding: 12px;
  }

  .v-card-actions {
    padding: 8px 12px;
  }
}

/* Ensure tables are scrollable on all screen sizes */
.v-table__wrapper {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

/* Improve form field usability on mobile */
@media (max-width: 599px) {
  .v-field__input {
    font-size: 16px; /* Prevents iOS zoom on focus */
  }
}
</style>
