import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';
import { execFileSync } from 'node:child_process';
import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';

const binary = fileURLToPath(new URL('../../bin/explain', import.meta.url));
const catalog = JSON.parse(
  readFileSync(new URL('../src/data/catalog.json', import.meta.url), 'utf8'),
);

test('search, aliases, tags, implementation filter, and empty state', async ({
  page,
}) => {
  const errors: string[] = [];
  page.on('pageerror', (e) => errors.push(e.message));
  await page.goto('/commands/');
  await expect(page.locator('[data-entry]:visible')).toHaveCount(4);
  await page.getByLabel('Find a command or a use case').fill('ggrep');
  await expect(page.locator('[data-entry]:visible')).toHaveCount(1);
  await expect(page.getByRole('status')).toHaveText('1 command');
  await expect(page).toHaveURL(/q=ggrep/);
  await page.reload();
  await expect(page.getByLabel('Find a command or a use case')).toHaveValue(
    'ggrep',
  );
  await page.getByLabel('Implementation coverage').selectOption('bash');
  await expect(page.locator('#no-results')).toBeVisible();
  await page.getByLabel('Find a command or a use case').fill('');
  await expect(page.locator('[data-entry]:visible')).toHaveCount(1);
  await page.getByLabel('Implementation coverage').selectOption('');
  await page.getByLabel('Find a command or a use case').fill('search logs');
  await expect(page.locator('[data-entry]:visible')).toHaveCount(1);
  await page.getByRole('link', { name: '#logs', exact: true }).click();
  await expect(
    page.getByRole('heading', { name: '#logs', exact: true }),
  ).toBeVisible();
  await expect(page.locator('[data-entry]:visible')).toHaveCount(1);
  expect(errors).toEqual([]);
});

for (const entry of catalog.entries) {
  for (const implementation of ['common', ...entry.implementations]) {
    test(`CLI/web parity: ${entry.id}/${implementation}`, async ({ page }) => {
      const response = JSON.parse(
        execFileSync(
          binary,
          [entry.id, '--implementation', implementation, '--format', 'json'],
          { encoding: 'utf8' },
        ),
      );
      const path = `/commands/${entry.id}/${implementation === 'common' ? '' : implementation + '/'}`;
      await page.goto(path);
      await expect(page.locator('.terminal-note code')).toHaveText(
        `explain ${entry.id} --implementation ${implementation}`,
      );
      await expect(page.locator('[data-section-id]')).toHaveCount(
        response.sections.length,
      );
      expect(
        await page
          .locator('[data-section-id]')
          .evaluateAll((nodes) =>
            nodes.map((n) => n.getAttribute('data-section-id')),
          ),
      ).toEqual(response.sections.map((s: { id: string }) => s.id));
      for (const section of response.sections) {
        const node = page.locator(`[data-section-id="${section.id}"]`);
        await expect(node.locator('h2')).toHaveText(section.title);
        expect(await node.locator('pre code').allTextContents()).toEqual(
          section.code,
        );
        const actual = await node.locator('.prose').evaluate((el) => {
          const copy = el.cloneNode(true) as HTMLElement;
          copy.querySelectorAll('button').forEach((b) => b.remove());
          return copy.innerHTML;
        });
        expect(actual.replace(/\s+/g, ' ').trim()).toBe(
          section.html.replace(/\s+/g, ' ').trim(),
        );
      }
      expect(response.content_revision).toBe(catalog.revision);
      await expect(page.locator('link[rel=canonical]')).toHaveAttribute(
        'href',
        `https://shell-explain.ninegene.com${path}`,
      );
      await expect(
        page
          .getByRole('navigation', { name: 'Implementation', exact: true })
          .locator('[aria-current=page]'),
      ).toHaveText(
        new RegExp(
          implementation === 'common'
            ? 'Portable'
            : implementation === 'apple'
              ? 'Apple'
              : implementation,
          'i',
        ),
      );
    });
  }
}

test('copy preserves exact code and implementation links are shareable', async ({
  page,
  context,
}) => {
  await context.grantPermissions(['clipboard-read', 'clipboard-write']);
  await page.goto('/commands/sed/');
  await page.getByRole('link', { name: 'macOS / Apple' }).click();
  await expect(page).toHaveURL('/commands/sed/apple/');
  const button = page
    .getByRole('button', { name: 'Copy code example' })
    .first();
  const expected = await page.locator('pre code').first().textContent();
  await button.click();
  expect(await page.evaluate(() => navigator.clipboard.readText())).toBe(
    expected,
  );
  await expect(page.locator('.copy-status')).toContainText('Example copied');
});

test('desktop and mobile accessibility, keyboard navigation, and overflow', async ({
  page,
}) => {
  for (const size of [
    { width: 1440, height: 1000 },
    { width: 390, height: 844 },
  ]) {
    await page.setViewportSize(size);
    for (const path of ['/', '/commands/grep/gnu/', '/tags/text/', '/about/']) {
      await page.goto(path);
      const results = await new AxeBuilder({ page })
        .withTags(['wcag2a', 'wcag2aa', 'wcag21aa'])
        .analyze();
      expect(results.violations).toEqual([]);
      expect(
        await page.evaluate(
          () => document.documentElement.scrollWidth <= innerWidth,
        ),
      ).toBe(true);
    }
  }
  await page.goto('/commands/');
  await page.keyboard.press('Tab');
  await expect(
    page.getByRole('link', { name: 'Skip to content' }),
  ).toBeFocused();
  await page.keyboard.press('Enter');
  await page.getByLabel('Find a command or a use case').focus();
  await page.keyboard.type('replace');
  await expect(page.locator('[data-entry]:visible')).toHaveCount(1);
});

test('static pages work without JavaScript and use no third-party resources', async ({
  browser,
}) => {
  const context = await browser.newContext({ javaScriptEnabled: false });
  const page = await context.newPage();
  const external: string[] = [];
  page.on('request', (r) => {
    if (!r.url().startsWith('http://127.0.0.1:4321')) external.push(r.url());
  });
  await page.goto('/commands/find/gnu/');
  await expect(
    page.getByRole('heading', { name: 'GNU examples' }),
  ).toBeVisible();
  await expect(page.locator('pre code')).toHaveCount(3);
  expect(external).toEqual([]);
  await context.close();
});

test('404 and visual review artifacts', async ({ page }) => {
  await page.goto('/404.html');
  await expect(
    page.getByRole('heading', { name: 'Take another path.' }),
  ).toBeVisible();
  await page.setViewportSize({ width: 1440, height: 1000 });
  await page.goto('/');
  await page.screenshot({
    path: '/tmp/shell-explain-desktop.png',
    fullPage: true,
  });
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto('/commands/sed/apple/');
  await page.screenshot({
    path: '/tmp/shell-explain-mobile.png',
    fullPage: true,
  });
});
