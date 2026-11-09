import { defineConfig, devices } from '@playwright/test';
// `npm test` picks a free port; the default avoids Astro's dev port (4321).
const port = process.env.WEB_TEST_PORT ?? '4322';
const origin = `http://127.0.0.1:${port}`;
export default defineConfig({
  testDir: './tests',
  fullyParallel: true,
  workers: 2,
  reporter: 'list',
  use: { baseURL: origin, trace: 'retain-on-failure' },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
  webServer: {
    command: `npm run build && npm run preview -- --port ${port} --ignore-lock`,
    url: origin,
    timeout: 120000,
    reuseExistingServer: false,
  },
});
