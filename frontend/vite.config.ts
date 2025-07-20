import react from '@vitejs/plugin-react'
import { resolve } from 'path'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: '../backend/internal/app/embed_static',
    emptyOutDir: true,
  },
  css: {
    postcss: './postcss.config.js',
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3000,
    hmr: {
      port: 3001,
    },
    proxy: {
      '/api': {
        target: `http://localhost:${process.env.PORT || '8080'}`,
        changeOrigin: true,
      },
    },
  },
})
