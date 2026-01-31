<template>
  <v-app-bar class="app-bar" height="78" flat>
    <v-container class="d-flex align-center">
      <div class="brand-mark">
        <v-avatar color="surface" size="44">
          <v-img alt="BrewPipes" src="@/assets/logo.svg" />
        </v-avatar>
      </div>

      <div class="brand-text">
        <div class="text-overline">BrewPipes</div>
        <div class="text-h6">Production Console</div>
      </div>

      <v-spacer />

      <v-chip class="text-caption" color="secondary" variant="tonal">
        Brew Day UX
      </v-chip>
    </v-container>
  </v-app-bar>

  <v-main class="app-main">
    <router-view />
  </v-main>

  <AppFooter />
</template>

<script lang="ts" setup>
  //
</script>

<style scoped>
.app-bar {
  background: rgba(253, 251, 247, 0.86);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(28, 26, 22, 0.08);
  position: sticky;
}

.brand-mark {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
}

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.1;
}

.app-main {
  position: relative;
  z-index: 1;
  padding-top: 12px;
  padding-bottom: 40px;
}
</style>
