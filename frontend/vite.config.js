import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      // Dev me /api calls Go backend pe forward honge
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})