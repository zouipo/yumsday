import { useAuthStore } from '@/stores/auth'

/**
 * Redirects authenticated users away from the login route
 */
export function redirectIfAuthenticated(_to, _from, next) {
  const auth = useAuthStore()
  auth.isAuthenticated ? next({ name: 'dashboard' }) : next()
}

/**
 * Redirects unauthenticated users to /login, preserving the intended URL.
 */
export function requireAuth(to, _from, next) {
  const auth = useAuthStore()
  auth.isAuthenticated ? next() : next({ name: 'login', query: { redirect: to.fullPath } })
}

/**
 * Allows access only to authenticated admin users.
 * - Not logged in → redirected to /login
 * - Logged in but not admin → redirected to /dashboard (or any appropriate fallback)
 */
export function requireAdmin(to, _from, next) {
  const auth = useAuthStore()
  if (!auth.isAuthenticated) return next({ name: 'login', query: { redirect: to.fullPath } })
  if (!auth.isAdmin) return next({ name: 'dashboard' })
  next()
}
