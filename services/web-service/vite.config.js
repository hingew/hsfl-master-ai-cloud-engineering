import { defineConfig } from "vite"
import { plugin as elmPlugin } from "vite-plugin-elm"

export default defineConfig({
    plugins: [
        elmPlugin()
    ],
    build: {
        outDir: "../public/"
    },
    server: {
        port: 3005,
        proxy: {
            '/api': 'http://localhost:3000',
            '/auth': 'http://localhost:3000'
        }
    }
})
