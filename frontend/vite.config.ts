import react from '@vitejs/plugin-react'
import { resolve } from 'path'
import { defineConfig } from 'vite'

export default defineConfig(({ mode }) => {
  if (mode !== 'development' && mode !== 'production') {
    throw new Error(`Unsupported mode: ${mode}. Only 'development' and 'production' are allowed.`)
  }

  const isDev = mode === 'development'

  return {
    plugins: [react()],
    build: {
      outDir: '../backend/internal/app/embed_static',
      emptyOutDir: true,
      sourcemap: isDev,
    },
    css: {
      postcss: './postcss.config.js',
    },
    resolve: {
      alias: [{ find: '@', replacement: resolve(__dirname, './src') }],
    },
  }
})
