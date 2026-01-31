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
    <v-btn class="profile-button" icon variant="text">
      <v-icon icon="mdi-account-circle" size="26" />
    </v-btn>
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
    </v-list>

    <template #append>
      <div class="drawer-footer">
        <v-btn block size="small" variant="tonal" @click="toggleRail">
          {{ rail ? 'Expand menu' : 'Collapse menu' }}
        </v-btn>
      </div>
    </template>
  </v-navigation-drawer>

  <v-main class="app-main">
    <router-view />
  </v-main>

  <AppFooter />
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useDisplay, useTheme } from 'vuetify'

const drawer = ref(true)
const rail = ref(false)
const theme = useTheme()
const display = useDisplay()
const isMobile = computed(() => display.smAndDown.value)
const isDark = computed(() => theme.global.current.value.dark)
const themeIcon = computed(() => (isDark.value ? 'mdi-weather-sunny' : 'mdi-weather-night'))
const themeStorageKey = 'brewpipes:theme'

const navItems = [
  { title: 'Dashboard', icon: 'mdi-view-dashboard-outline', to: '/' },
  { title: 'Batches', icon: 'mdi-barley', to: '/batches' },
  { title: 'Vessels', icon: 'mdi-silo', to: '/vessels' },
  { title: 'Transfers', icon: 'mdi-truck-fast-outline', to: '/transfers' },
  { title: 'Additions', icon: 'mdi-flask-outline', to: '/additions' },
  { title: 'Measurements', icon: 'mdi-thermometer', to: '/measurements' },
  { title: 'Reports', icon: 'mdi-chart-bar', to: '/reports' },
]

function toggleDrawer() {
  if (isMobile.value) {
    drawer.value = !drawer.value
    return
  }
  rail.value = !rail.value
}

function toggleRail() {
  rail.value = !rail.value
}

function toggleTheme() {
  theme.global.name.value = isDark.value ? 'brewLight' : 'brewDark'
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

.drawer-footer {
  padding: 16px;
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
