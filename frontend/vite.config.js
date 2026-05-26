import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    tailwindcss(),
    react(),
  ],
  base: '/platform/',
  build: {
    rollupOptions: {
      output: {
        manualChunks: (id) => {
          if (id.includes('@monaco-editor') || id.includes('monaco-editor')) {
            return 'vendor-monaco';
          }
          if (id.includes('echarts')) {
            return 'vendor-echarts';
          }
        },
      },
    },
  },
});
