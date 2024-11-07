import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'node:path'

// https://vitejs.dev/config/
export default defineConfig({
	root: "src",
	server: {
		proxy: {
			"/pv/api": {
				target: "http://127.0.0.1:61016",
				changeOrigin: true
			},
		},
		host: true
	},
	plugins: [
		vue(),
	],
	build: {
		outDir: "../../dist/client",
		emptyOutDir: true,
		rollupOptions: {
			input: {
				index: resolve(__dirname, "src/index.html"),
				admin: resolve(__dirname, "src/admin.html"),
				datenschutz: resolve(__dirname, "src/legal/datenschutz.html"),
				impressum: resolve(__dirname, "src/legal/impressum.html"),
			},
			external: [
				"/public/banner_energiespenderin.svg"
			]
		}
	},
	resolve: {
		alias: {
			'@': fileURLToPath(new URL('./src', import.meta.url))
		}
	}
})
