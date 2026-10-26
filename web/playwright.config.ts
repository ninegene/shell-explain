import { defineConfig, devices } from '@playwright/test';
export default defineConfig({
  testDir: './tests',
  fullyParallel: true,
  workers: 2,
  reporter: 'list',
  use: { baseURL: 'http://127.0.0.1:4321', trace: 'retain-on-failure' },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
  webServer: {
    command: 'npm run build && npm run preview -- --port 4321 --ignore-lock',
    url: 'http://127.0.0.1:4321',
    timeout: 120000,
    reuseExistingServer: false,
  },
});
