import { defineConfig } from "vite";
import { viteSingleFile } from "vite-plugin-singlefile";

export default defineConfig({
  plugins: [viteSingleFile()],
  server: { port: 5172, strictPort: true },
  preview: { port: 5172, strictPort: true },
});
