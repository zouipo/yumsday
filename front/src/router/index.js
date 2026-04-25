import { createRouter, createWebHistory } from 'vue-router'
import { routes } from './routeConfig'
import { requireAuth } from './guards'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach(requireAuth)

export default router