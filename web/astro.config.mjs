import { defineConfig } from 'astro/config';
import { contentWatch } from './scripts/content-watch.mjs';

export default defineConfig({
  site: 'https://shell-explain.ninegene.com',
  output: 'static',
  trailingSlash: 'always',
  devToolbar: { enabled: false },
  vite: {
    plugins: [contentWatch()],
    // Never inline scripts into pages: the site's Content-Security-Policy allows only
    // script-src 'self' (infra/pulumi/index.ts).
    build: { assetsInlineLimit: 0 },
  },
});
