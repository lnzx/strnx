import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import AutoImport from 'unplugin-auto-import/vite'
import VueRouter from 'unplugin-vue-router/vite'
import { VueRouterAutoImports } from 'unplugin-vue-router'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    outDir: "../dist/"
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  plugins: [
    VueRouter(),
    vue(), // ⚠️ Vue must be placed after VueRouter()
    Components(),
    AutoImport({
      imports: [
        'vue',
        VueRouterAutoImports,
        {
          '@/composable/useSession': [
            ['default', 'useSession']
          ],
          'axios': [
            ['default', 'axios'], // import { default as axios } from 'axios',
          ],
          '@/composable/useApi': [
            ['useApi'],
          ],
        }
      ]
    }),
  ],
})
