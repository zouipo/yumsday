<script setup>
import { computed, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route  = useRoute()
const auth   = useAuthStore()

const username = ref('')
const password = ref('')
const error    = ref('')
const loading  = ref(false)
const canSubmit = computed(() => username.value.trim() !== '' && password.value.trim() !== '')

async function handleLogin() {
  if (!canSubmit.value) return

  error.value   = ''
  loading.value = true
  try {
    await auth.login(username.value, password.value)
    const redirect = route.query.redirect || '/'
    router.push(redirect)
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-wrapper">

    <div class="login-card card">
      <header class="card-header">
        <span class="symbol" :class="{ 'symbol--error': !!error }">●</span>
        <h1>Feeling hungry?</h1>
        <p>Welcome back — enter your details below.</p>
      </header>

      <form @submit.prevent="handleLogin">
        <Transition name="shake">
          <p v-if="error" class="error-banner" role="alert">
            <span class="error-banner__icon">!</span> {{ error }}
          </p>
        </Transition>

        <div class="field">
          <label for="username">Username</label>
          <input
            id="username"
            v-model="username"
            type="text"
            required
            autocomplete="username"
            placeholder="your username"
          />
        </div>

        <div class="field">
          <label for="password">
            Password
          </label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            autocomplete="current-password"
            placeholder="••••••••"
          />
        </div>

        <button
          type="submit"
          class="btn btn--primary btn--full btn-submit"
          :disabled="loading || !canSubmit"
        >
          <span v-if="loading" class="spinner" aria-hidden="true" />
          {{ loading ? 'Signing in…' : 'Sign in' }}
        </button>
      </form>
    </div>

  </div>
</template>

<style scoped>
/* ── Layout ──────────────────────────────────────────────────── */
.login-wrapper {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 2rem;
}

/* ── Card ────────────────────────────────────────────────────── */
.login-card {
  width: 100%;
  max-width: 420px;
  padding: 2.5rem;
  animation: rise .35s ease both;
}

/* ── Card header ─────────────────────────────────────────────── */
.card-header {
  text-align: center;
  margin-bottom: 2rem;
}

.symbol {
  display: inline-block;
  font-size: 1.5rem;
  color: var(--color-secondary);
  margin-bottom: .75rem;
}

.symbol--error {
  color: var(--color-danger);
}

.card-header p {
  font-size: .875rem;
  color: var(--color-muted);
  margin-top: .35rem;
}

/* ── Forgot link ─────────────────────────────────────────────── */
.forgot {
  font-size: .75rem;
  text-transform: none;
  letter-spacing: 0;
}

/* ── shake transition ────────────────────────────────────────── */
.shake-enter-active { animation: shake .4s ease; }

/* ── Submit button ───────────────────────────────────────────── */
.btn-submit { margin-top: .5rem; }

/* ── Card footer ─────────────────────────────────────────────── */
.card-footer {
  margin-top: 1.75rem;
  text-align: center;
  font-size: .875rem;
  color: var(--color-muted);
}
</style>