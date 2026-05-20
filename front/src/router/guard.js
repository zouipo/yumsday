import { useAuthStore } from '@/stores/auth'

/**
 * Redirects authenticated users away from the login route
 */
export function redirectIfAuthenticated(_to, _from) {
  const auth = useAuthStore()
  return auth.isAuthenticated ? { name: 'home' } : undefined
}

/**
 * Redirects unauthenticated users to /login, preserving the intended URL.
 */
export function requireAuth(to, _from) {
  const requiresAuth = to.matched.some((record) => record.meta?.requiresAuth)
  if (!requiresAuth) return undefined

  const auth = useAuthStore()
  return auth.isAuthenticated ? undefined : { name: 'login', query: { redirect: to.fullPath } }
}

/**
 * Allows access only to authenticated admin users.
 * - Not logged in -> redirected to /login
 * - Logged in but not admin -> redirected to / (home fallback)
 */
export function requireAdmin(to, _from) {
  const auth = useAuthStore()
  if (!auth.isAuthenticated) return { name: 'login', query: { redirect: to.fullPath } }
  if (!auth.isAdmin) return { name: 'home' }
  return undefined
}
