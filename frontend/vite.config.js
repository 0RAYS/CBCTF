import { defineConfig } from 'vite'
import transformImports from '@rolldown/plugin-transform-imports'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'

const getPackageName = (id) => {
  const normalized = id.replace(/\\/g, '/');
  const marker = '/node_modules/';
  const index = normalized.lastIndexOf(marker);

  if (index === -1) {
    return '';
  }

  const parts = normalized.slice(index + marker.length).split('/');
  return parts[0]?.startsWith('@') ? `${parts[0]}/${parts[1]}` : parts[0];
};

const isPnpmPackage = (id, packageName) => {
  const normalized = id.replace(/\\/g, '/');
  return normalized.includes(`/node_modules/.pnpm/${packageName}@`);
};

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    transformImports({
      '@tabler/icons-react': {
        transform: '@tabler/icons-react/dist/esm/icons/{{member}}.mjs',
      },
    }),
    tailwindcss(),
    react(),
  ],
  base: '/platform/',
  build: {
    modulePreload: {
      resolveDependencies: (_, deps) => deps.filter((dep) => !/vendor-(monaco|echarts)/.test(dep)),
    },
    rollupOptions: {
      output: {
        manualChunks: (id) => {
          const packageName = getPackageName(id);

          if (
            packageName === 'react' ||
            packageName === 'react-dom' ||
            packageName === 'scheduler' ||
            isPnpmPackage(id, 'react') ||
            isPnpmPackage(id, 'react-dom') ||
            isPnpmPackage(id, 'scheduler')
          ) {
            return 'vendor-react';
          }
          if (packageName === 'monaco-editor') {
            return 'vendor-monaco';
          }
          if (packageName === 'echarts' || packageName === 'zrender') {
            return 'vendor-echarts';
          }
        },
      },
    },
  },
});
