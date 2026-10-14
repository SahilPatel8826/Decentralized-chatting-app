import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// The Go services send no CORS headers, so the browser talks to Vite
// and Vite forwards the calls. Ports must match your services.
const USER_SERVICE = 'http://localhost:8080'
const CHAT_SERVICE = 'http://localhost:8081'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/user-api': { target: USER_SERVICE, rewrite: (p) => p.replace(/^\/user-api/, '') },
      '/chat-api': { target: CHAT_SERVICE, rewrite: (p) => p.replace(/^\/chat-api/, '') },
      '/ws': { target: CHAT_SERVICE.replace('http', 'ws'), ws: true },
    },
  },
})
