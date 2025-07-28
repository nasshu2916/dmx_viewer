import react from '@vitejs/plugin-react'
import { resolve } from 'path'
import { defineConfig } from 'vite'

export default defineConfig(({ mode }) => ({
  plugins: [react()],
  build: {
    outDir: '../backend/internal/app/embed_static',
    emptyOutDir: true,
    sourcemap: mode === 'debug',
  },
  css: {
    postcss: './postcss.config.js',
  },
  resolve: {
    alias: [{ find: '@', replacement: resolve(__dirname, './src') }],
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
}))
