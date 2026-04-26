import { fileURLToPath, URL } from 'node:url'
import fs from 'node:fs'
import path from 'node:path'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

function readBackendAddressFromRootConfig() {
  const configPath = path.resolve(fileURLToPath(new URL('.', import.meta.url)), '../config/config.yaml')
  const fallback = { host: 'localhost', port: 8080 }

  try {
    const content = fs.readFileSync(configPath, 'utf8')
    const lines = content.split('\n')

    let host = fallback.host
    let port = fallback.port

    for (const rawLine of lines) {
      const line = rawLine.trim()
      if (!line || line.startsWith('#')) continue

      const separatorIndex = line.indexOf(':')
      if (separatorIndex === -1) continue

      const key = line.slice(0, separatorIndex).trim()
      const value = line.slice(separatorIndex + 1).trim()

      if (key === 'host' && value) host = value
      if (key === 'port') {
        const parsed = Number.parseInt(value, 10)
        if (Number.isFinite(parsed)) port = parsed
      }
    }

    return { host, port }
  } catch {
    return fallback
  }
}

const backendAddress = readBackendAddressFromRootConfig()
const backendTarget = `http://${backendAddress.host}:${backendAddress.port}`

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    proxy: {
      '/api': {
        target: backendTarget,
        changeOrigin: true,
      },
      '/auth': {
        target: backendTarget,
        changeOrigin: true,
      },
    },
  },
})
