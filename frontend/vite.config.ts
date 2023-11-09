import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import path from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      'Stores': path.resolve(__dirname, './src/state/stores.ts'),
      'Atoms': path.resolve(__dirname, './src/state/atoms.ts'),
    }
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:3000/api',
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  },
  build: {
    outDir: '../backend/public',
    emptyOutDir: true,
  }
})
