import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react-swc'

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
        manualChunks: {
          'vendor-monaco': ['@monaco-editor/react'],
          'vendor-echarts': ['echarts', 'echarts-for-react'],
        },
      },
    },
  },
});
