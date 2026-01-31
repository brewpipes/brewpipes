<template>
  <v-app-bar class="app-bar" height="72" flat>
    <v-app-bar-nav-icon @click="toggleDrawer" />

    <div class="brand-mark">
      <v-avatar color="surface" size="38">
        <v-img alt="BrewPipes" src="@/assets/logo.svg" />
      </v-avatar>
    </div>

    <div class="brand-text">
      <div class="text-caption text-medium-emphasis">BrewPipes</div>
      <div class="text-subtitle-1">Production Console</div>
    </div>

    <v-spacer />

    <v-btn class="theme-button" icon variant="text" @click="toggleTheme">
      <v-icon :icon="themeIcon" size="22" />
    </v-btn>
    <v-menu location="bottom end" offset="8">
      <template #activator="{ props }">
        <v-btn class="profile-button" icon variant="text" v-bind="props">
          <v-icon icon="mdi-account-circle" size="26" />
        </v-btn>
      </template>
      <v-list density="compact">
        <v-list-item>
          <v-list-item-title class="text-body-2 font-weight-medium">{{ userLabel }}</v-list-item-title>
          <v-list-item-subtitle>Signed in</v-list-item-subtitle>
        </v-list-item>
        <v-divider />
        <v-list-item prepend-icon="mdi-logout" title="Logout" @click="handleLogout" />
      </v-list>
    </v-menu>
  </v-app-bar>

  <v-navigation-drawer
    v-model="drawer"
    class="app-drawer"
    :permanent="!isMobile"
    :rail="rail"
    :temporary="isMobile"
    width="260"
    rail-width="78"
  >
    <v-list class="nav-list" density="comfortable" nav>
      <v-list-item
        v-for="item in navItems"
        :key="item.title"
        :to="item.to"
        :prepend-icon="item.icon"
        :title="item.title"
      />

      <v-list-group class="nav-group" value="inventory">
        <template #activator="{ props }">
          <v-list-item v-bind="props" prepend-icon="mdi-warehouse" title="Inventory" />
        </template>
        <v-list-item to="/inventory/activity" title="Activity" />
        <v-list-item to="/inventory/product" title="Product" />
        <v-list-item to="/inventory/ingredients" title="Ingredients" />
        <v-list-item to="/inventory/adjustments-transfers" title="Adjustments & Transfers" />
        <v-list-item to="/inventory/locations" title="Locations" />
      </v-list-group>

      <v-list-group class="nav-group" value="procurement">
        <template #activator="{ props }">
          <v-list-item v-bind="props" prepend-icon="mdi-clipboard-text" title="Procurement" />
        </template>
        <v-list-item to="/procurement" title="Overview" />
        <v-list-item to="/procurement/suppliers" title="Suppliers" />
        <v-list-item to="/procurement/purchase-orders" title="Purchase orders" />
        <v-list-item to="/procurement/purchase-order-lines" title="Order lines" />
        <v-list-item to="/procurement/purchase-order-fees" title="Order fees" />
      </v-list-group>
    </v-list>
  </v-navigation-drawer>

  <v-main class="app-main">
    <router-view />
  </v-main>

  <AppFooter />
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useDisplay, useTheme } from 'vuetify'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const drawer = ref(true)
const rail = ref(false)
const theme = useTheme()
const display = useDisplay()
const isMobile = computed(() => display.smAndDown.value)
const isDark = computed(() => theme.global.current.value.dark)
const themeIcon = computed(() => (isDark.value ? 'mdi-weather-sunny' : 'mdi-weather-night'))
const themeStorageKey = 'brewpipes:theme'
const userLabel = computed(() => authStore.username ?? 'Account')

const navItems = [
  { title: 'Dashboard', icon: 'mdi-view-dashboard-outline', to: '/' },
  { title: 'Batches', icon: 'mdi-barley', to: '/batches' },
  { title: 'Vessels', icon: 'mdi-silo', to: '/vessels' },
]

function toggleDrawer() {
  if (isMobile.value) {
    drawer.value = !drawer.value
    return
  }
  rail.value = !rail.value
}

function toggleTheme() {
  theme.global.name.value = isDark.value ? 'brewLight' : 'brewDark'
}

async function handleLogout() {
  await authStore.logout()
  await router.push('/login')
}

onMounted(() => {
  const storedTheme = localStorage.getItem(themeStorageKey)
  if (storedTheme === 'brewLight' || storedTheme === 'brewDark') {
    theme.global.name.value = storedTheme
    return
  }

  const prefersDark = window.matchMedia?.('(prefers-color-scheme: dark)').matches
  theme.global.name.value = prefersDark ? 'brewDark' : 'brewLight'
})

watch(
  () => theme.global.name.value,
  (value) => {
    localStorage.setItem(themeStorageKey, value)
  },
)
</script>

<style scoped>
.app-bar {
  background: rgba(var(--v-theme-surface), 0.92);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(var(--v-theme-on-surface), 0.12);
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

.app-drawer {
  background: rgba(var(--v-theme-surface), 0.96);
  border-right: 1px solid rgba(var(--v-theme-on-surface), 0.12);
}

.nav-list :deep(.v-list-item) {
  border-radius: 8px;
}

.nav-group :deep(.v-list-group__items .v-list-item) {
  padding-inline-start: 60px !important;
}

.nav-group :deep(.v-list-group__items) {
  padding-inline-start: 4px;
}

.nav-group {
  --v-list-indent: 0px;
  --list-indent-size: 0px;
}

.app-drawer.v-navigation-drawer--rail :deep(.v-list-item__content) {
  opacity: 0;
  width: 0;
  overflow: hidden;
}

.theme-button {
  color: rgba(var(--v-theme-on-surface), 0.7);
}

.profile-button {
  color: rgba(var(--v-theme-on-surface), 0.7);
}

.app-main {
  position: relative;
  z-index: 1;
  padding-top: calc(var(--v-layout-top) + 8px);
  padding-right: calc(var(--v-layout-right) + 8px);
  padding-bottom: calc(var(--v-layout-bottom) + 8px);
  padding-left: calc(var(--v-layout-left) + 8px);
}
</style>
