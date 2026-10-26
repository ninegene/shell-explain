import { defineConfig } from 'astro/config';
import { contentWatch } from './scripts/content-watch.mjs';

export default defineConfig({
  site: 'https://shell-explain.ninegene.com',
  output: 'static',
  trailingSlash: 'always',
  devToolbar: { enabled: false },
  vite: { plugins: [contentWatch()] },
});
