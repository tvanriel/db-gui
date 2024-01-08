import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  
  plugins: [vue({isProduction: false})],

  base: "/app/",
  server: {
    port: 3000,
    proxy: {
      "/api/v1/sql": {
        target: "http://127.0.0.1:8080",
        changeOrigin: true,
        ws:true,
      },
      "/api": {
        target: "http://127.0.0.1:8080",
        changeOrigin: true,
      }
    },
  },
  resolve: {
    alias: {
      "~bootstrap": path.resolve(__dirname, "node_modules/bootstrap"),
      "~bootswatch": path.resolve(__dirname, "node_modules/bootswatch"),
    }
  },
  build: {
    outDir: "../app/frontend/static",
    emptyOutDir: true,
  }
});
