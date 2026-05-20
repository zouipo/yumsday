import { createRouter, createWebHistory } from 'vue-router'
import { routes } from './routesConfig'
import { requireAuth } from './guard'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach(requireAuth)

export default router