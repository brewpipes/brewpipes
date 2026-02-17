/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 */

import { setupLayouts } from 'virtual:generated-layouts'
// Composables
import { createRouter, createWebHistory } from 'vue-router'
import { routes } from 'vue-router/auto-routes'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: setupLayouts(routes),
})

router.beforeEach(to => {
  const authStore = useAuthStore()
  if (!authStore.hydrated) {
    authStore.hydrateFromStorage()
  }

  const isPublic = Boolean(to.meta.public) || to.path === '/login'
  if (!authStore.isAuthenticated && !isPublic) {
    return {
      path: '/login',
      query: { redirect: to.fullPath },
    }
  }

  if (authStore.isAuthenticated && to.path === '/login') {
    const redirect = to.query.redirect
    if (typeof redirect === 'string' && redirect.startsWith('/') && !redirect.startsWith('//')) {
      return redirect
    }
    return '/'
  }

  return true
})

// Workaround for https://github.com/vitejs/vite/issues/11804
router.onError((err, to) => {
  if (err?.message?.includes?.('Failed to fetch dynamically imported module') ||
      err?.message?.includes?.('Importing a module script failed')) {
    // Only attempt one automatic reload per browser session
    const reloadKey = 'brewpipes:chunk-reload'
    if (!sessionStorage.getItem(reloadKey)) {
      console.warn('[router] Reloading to resolve stale chunk', to.fullPath)
      sessionStorage.setItem(reloadKey, '1')
      location.assign(to.fullPath)
    } else {
      console.error('[router] Chunk load error persists after reload — not retrying', err)
    }
  } else {
    console.error(err)
  }
})

router.isReady().then(() => {
  // Clear the reload flag once the app boots successfully
  sessionStorage.removeItem('brewpipes:chunk-reload')
})

export default router
