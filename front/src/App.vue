<script setup>
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const showHeader = computed(() => route.name !== 'login')
const authButtonLabel = computed(() => (auth.isAuthenticated ? 'Log out' : 'Login'))

const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 12) return 'Good morning'
  if (h < 18) return 'Good afternoon'
  return 'Good evening'
})

async function handleAuthAction() {
  if (auth.isAuthenticated) {
    await auth.logout()
    router.push({ name: 'home' })
    return
  }

  router.push({ name: 'login' })
}

onMounted(async () => {
  // Verify if an existing session is still valid when the app loads
  await auth.fetchCurrentUser()

  if (auth.isAuthenticated && route.name === 'login') {
    router.replace({ name: 'home' })
  }
})

watch(
  () => auth.isAuthenticated,
  (isAuthenticated) => {
    if (isAuthenticated && route.name === 'login') {
      router.replace({ name: 'home' })
    }
  },
)
</script>

<template>
  <header v-if="showHeader" class="app-header">
    <p v-if="auth.isAuthenticated" class="greeting">{{ greeting }}, {{ auth.user?.username }} 👋</p>
    <button type="button" class="btn btn--ghost auth-btn" @click="handleAuthAction">
      {{ authButtonLabel }}
    </button>
  </header>
  <RouterView />
</template>

<style scoped>
.app-header {
  max-width: 960px;
  margin: 0 auto;
  padding: 1.5rem 2rem 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.greeting {
  font-size: .95rem;
  opacity: .7;
  margin: 0;
}

.auth-btn {
  margin-left: auto;
}
</style>
