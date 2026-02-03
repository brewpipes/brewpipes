import type { AuthTokens } from '@/types'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

const identityApiBase = import.meta.env.VITE_IDENTITY_API_URL ?? '/api'
const storageKeys = {
  access: 'brewpipes:accessToken',
  refresh: 'brewpipes:refreshToken',
  username: 'brewpipes:username',
}

let refreshInFlight: Promise<boolean> | null = null

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const username = ref<string | null>(null)
  const hydrated = ref(false)

  const isAuthenticated = computed(() => Boolean(accessToken.value))

  const readStorage = (key: string) => {
    try {
      return localStorage.getItem(key)
    } catch {
      return null
    }
  }

  const writeStorage = (key: string, value: string | null) => {
    try {
      if (value === null) {
        localStorage.removeItem(key)
      } else {
        localStorage.setItem(key, value)
      }
    } catch {
      // ignore storage errors
    }
  }

  const setTokens = (tokens: AuthTokens) => {
    accessToken.value = tokens.access_token
    refreshToken.value = tokens.refresh_token
    writeStorage(storageKeys.access, tokens.access_token)
    writeStorage(storageKeys.refresh, tokens.refresh_token)
  }

  const setUsername = (value: string | null) => {
    username.value = value
    writeStorage(storageKeys.username, value)
  }

  const clearSession = () => {
    accessToken.value = null
    refreshToken.value = null
    username.value = null
    hydrated.value = true
    writeStorage(storageKeys.access, null)
    writeStorage(storageKeys.refresh, null)
    writeStorage(storageKeys.username, null)
  }

  const hydrateFromStorage = () => {
    if (hydrated.value) {
      return
    }
    accessToken.value = readStorage(storageKeys.access)
    refreshToken.value = readStorage(storageKeys.refresh)
    username.value = readStorage(storageKeys.username)
    hydrated.value = true
  }

  const login = async (user: string, password: string) => {
    const response = await fetch(`${identityApiBase}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username: user, password }),
    })

    if (!response.ok) {
      const message = await response.text()
      throw new Error(message || 'Unable to sign in')
    }

    const data = (await response.json()) as AuthTokens
    setTokens(data)
    setUsername(user)
    hydrated.value = true
    return data
  }

  const refresh = async () => {
    if (!refreshToken.value) {
      clearSession()
      return false
    }

    if (refreshInFlight) {
      return refreshInFlight
    }

    refreshInFlight = (async () => {
      try {
        const response = await fetch(`${identityApiBase}/refresh`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ refresh_token: refreshToken.value }),
        })

        if (!response.ok) {
          clearSession()
          return false
        }

        const data = (await response.json()) as AuthTokens
        setTokens(data)
        hydrated.value = true
        return true
      } catch {
        clearSession()
        return false
      } finally {
        refreshInFlight = null
      }
    })()

    return refreshInFlight
  }

  const logout = async () => {
    const token = refreshToken.value
    clearSession()
    if (!token) {
      return
    }

    try {
      await fetch(`${identityApiBase}/logout`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ refresh_token: token }),
      })
    } catch {
      // ignore logout errors
    }
  }

  return {
    accessToken,
    refreshToken,
    username,
    hydrated,
    isAuthenticated,
    hydrateFromStorage,
    login,
    refresh,
    logout,
    clearSession,
  }
})
