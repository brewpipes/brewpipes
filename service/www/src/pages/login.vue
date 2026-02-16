<route lang="yaml">
meta:
  layout: auth
  public: true
</route>

<template>
  <v-card class="login-card" elevation="12">
    <v-card-text>
      <div class="d-flex align-center ga-3 mb-6">
        <v-avatar color="surface" size="42">
          <v-img alt="BrewPipes" src="@/assets/logo.svg" />
        </v-avatar>
        <div>
          <div class="text-caption text-medium-emphasis">BrewPipes</div>
          <div class="text-h5">Sign in</div>
        </div>
      </div>

      <v-form @submit.prevent="submit">
        <v-text-field
          v-model="username"
          autocomplete="username"
          autofocus
          density="comfortable"
          hide-details="auto"
          label="Username"
        />
        <v-text-field
          v-model="password"
          autocomplete="current-password"
          density="comfortable"
          hide-details="auto"
          label="Password"
          :type="showPassword ? 'text' : 'password'"
        >
          <template #append-inner>
            <v-icon
              aria-label="Toggle password visibility"
              :icon="showPassword ? 'mdi-eye-off-outline' : 'mdi-eye-outline'"
              @click="showPassword = !showPassword"
            />
          </template>
        </v-text-field>

        <v-btn
          block
          class="mt-4"
          color="primary"
          :loading="submitting"
          type="submit"
        >
          Sign in
        </v-btn>
      </v-form>

      <v-alert
        v-if="errorMessage"
        class="mt-4"
        density="compact"
        type="error"
        variant="tonal"
      >
        {{ errorMessage }}
      </v-alert>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { useAuthStore } from '@/stores/auth'

  const authStore = useAuthStore()
  const router = useRouter()
  const route = useRoute()

  const username = ref('')
  const password = ref('')
  const showPassword = ref(false)

  const { execute, loading: submitting, error: errorMessage } = useAsyncAction()

  const redirectPath = computed(() => {
    const redirect = route.query.redirect
    if (typeof redirect === 'string' && redirect.startsWith('/')) {
      return redirect
    }
    return '/'
  })

  async function submit () {
    const user = username.value.trim()
    if (!user || !password.value) {
      errorMessage.value = 'Enter both username and password to continue.'
      return
    }

    await execute(async () => {
      await authStore.login(user, password.value)
      password.value = ''
      await router.replace(redirectPath.value)
    })
  }
</script>

<style scoped>
.login-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 18px 38px rgba(0, 0, 0, 0.22);
}
</style>
