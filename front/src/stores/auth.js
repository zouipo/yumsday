import {defineStore} from 'pinia'
import {ref, computed} from 'vue'

export const useAuthStore = defineStore('auth', () => {
    const user = ref(null)
    
    const isAuthenticated = computed(() => !!user.value)  // does user exists?

    async function login(credentials) {
        // The server sets the session cookie in its response headers (Set-Cookie).
        // For later requests, the browser will manage to include session cookie for this site (Same-Site: script).
        const res = await fetch('/api/login', {
            method: 'POST',
            credentials: 'include', // send and receive cookies; not necessary for same-origin requests but good practice
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(credentials),
        })

        if (!res.ok) throw new Error('Login failed!')

        user.value = await res.json()
    }

    // Verify if an existing session is still valid.
    async function fetchCurrentUser() {
        try {
            const res = await fetch('/api/auth/me', {
                credentials: 'include',
            })

            // Typically a 401 error if the cookie has expired or was cleared by the server
            if (!res.ok) {
                user.value = null
                return
            }

            user.value = await res.json()
        } catch {   // network error for example
            user.value = null
        }
    }

    // Remove the session from the server, then clear local storage.
    async function logout() {
        try {
            await fetch('/api/logout', {
                method: 'POST',
                credentials: 'include',
            })
        } finally {
            // Clear local stage regardless of whether the server call succeeded
            user.value = null
        }
    }

    return {user, isAuthenticated, login, logout, fetchCurrentUser}
})
