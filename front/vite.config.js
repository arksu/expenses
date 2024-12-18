import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

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
    host: '0.0.0.0',
    proxy: {
      '/api': {
        target: 'http://localhost:9101', // Replace with your backend server URL
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''), // Rewrite the URL path

        // configure: (proxy) => {
        //   proxy.on('proxyReq', (proxyReq, req) => {
        //     console.log(`[Proxy Request] ${req.method} ${req.url}`);
        //   });
        //   proxy.on('proxyRes', (proxyRes, req) => {
        //     console.log(`[Proxy Response] ${proxyRes.statusCode} ${req.url}`);
        //   });
        // },
      },
      
    },
  },
})
