import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0',
    proxy: { 
      '/session': 'http://localhost:8080',
      '/user': 'http://localhost:8080',
      '/conversations': 'http://localhost:8080',
      '/messages': 'http://localhost:8080',
      '/groups': 'http://localhost:8080'
    }
  }
})