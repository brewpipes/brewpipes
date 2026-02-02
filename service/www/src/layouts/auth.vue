<template>
  <v-main class="auth-main">
    <div class="auth-shell">
      <router-view />
    </div>
  </v-main>
</template>

<style scoped>
.auth-main {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.auth-shell {
  width: min(520px, 100%);
}
</style>
