import { resolve } from "node:path"
import tailwindcss from "@tailwindcss/vite"
import react from "@vitejs/plugin-react"
import { defineConfig, loadEnv } from "vite"

// Set BASE_PATH in CI to the repository's Pages base path, such as /ha-go.
// VITE_SITE_URL is the absolute URL used by canonical and social metadata.
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "")
  const base = env.BASE_PATH
    ? `${env.BASE_PATH.replace(/\/$/, "")}/`
    : "/"

  return {
    base,
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        "@": resolve(import.meta.dirname, "./src"),
      },
    },
    build: {
      outDir: "dist",
      assetsDir: "assets",
    },
  }
})
