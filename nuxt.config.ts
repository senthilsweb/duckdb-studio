import { defineNuxtConfig } from "nuxt/config";
// https://v3.nuxtjs.org/docs/directory-structure/nuxt.config
export default defineNuxtConfig({
  app: {
    pageTransition: { name: "page", mode: "out-in" },
    layoutTransition: { name: "layout", mode: "out-in" },
  },
  ssr: false,
  sourcemap: false,
  tailwindcss: {
    cssPath: "~/assets/css/tailwind.css",
    configPath: "tailwind.config",
    exposeConfig: false,
    config: {},
    injectPosition: 0,
    viewer: true,
  },
  routeRules: {
    // Static page generated on-demand once
    //'/': { static: true },
  },
  router: {},
  build: {
    transpile: ["@heroicons/vue", "primevue"],
    postcss: {
      plugins: {
        tailwindcss: {},
        autoprefixer: {},
      },
    },
  },
  buildModules: [
    "@nuxt/postcss8",
    // ...
  ],
  runtimeConfig: {
    public: {
      API_BASE_URL: process.env.API_BASE_URL || "",
      CONFIG_BASE_URL: process.env.CONFIG_BASE_URL || "",
      LOGO_WEB: process.env.LOGO_WEB || 'logo.svg',
      API_GATEWAY_BASE_URL: process.env.API_GATEWAY_BASE_URL || '',
      DUCKDB_DATA_API_BASE_PATH: process.env.DUCKDB_DATA_API_BASE_PATH || '',
    },
  },
  modules: [
    "@nuxt/content",
    "@tailvue/nuxt",
    "@nuxtjs/tailwindcss",
    "nuxt-icon",
    "nuxt-lodash",
    "@nuxtjs/robots",
  ],
  mdc: {
    highlight: {
      theme: {
        default: "vitesse-light",
        dark: "material-theme-palenight",
      },
    },
  },
  devtools: {
    enabled: true,
  },
  content: {
    documentDriven: true,
  }
});
