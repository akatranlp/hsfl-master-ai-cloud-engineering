import path from "path";
import { UserConfig, defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import fs from "fs";

let https: UserConfig["server"]["https"] = undefined;
if (process.env.HTTPS) {
  https = {
    key: fs.readFileSync("certs/localhost.key"),
    cert: fs.readFileSync("certs/localhost.crt"),
  };
}

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "0.0.0.0",
    proxy: {
      "/api": {
        target: "http://localhost:8080",
      },
    },
    https,
  },
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
