import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// If the user is still authenticated, no authentication required
const auth = useAuthStore()
auth.fetchCurrentUser().finally(() => {
  app.mount('#app')
})

