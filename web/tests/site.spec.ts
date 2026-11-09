import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';
import { execFileSync } from 'node:child_process';
import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';

const binary = fileURLToPath(new URL('../../bin/explain', import.meta.url));
const catalog = JSON.parse(
  readFileSync(new URL('../src/data/catalog.json', import.meta.url), 'utf8'),
);

// Number of code blocks a page for entry/implementation shows (common + variant).
function codeBlockCount(id: string, implementation: string) {
  const entry = catalog.entries.find((e: { id: string }) => e.id === id);
  return entry.sections
    .filter(
      (s: { implementation: string }) =>
        s.implementation === 'common' || s.implementation === implementation,
    )
    .reduce((n: number, s: { code: string[] }) => n + s.code.length, 0);
}

test('search, aliases, tags, implementation filter, and empty state', async ({
  page,
}) => {
  const errors: string[] = [];
  page.on('pageerror', (e) => errors.push(e.message));
  await page.goto('/commands/');
  await expect(page.locator('[data-entry]:visible')).toHaveCount(
    catalog.entries.length,
  );
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
  await expect(page.locator('[data-entry]:visible')).toHaveCount(
    catalog.entries.filter((e: { tags: string[] }) => e.tags.includes('logs'))
      .length,
  );
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
        `explain ${entry.id} --implementation ${implementation} --details`,
      );
      const exampleSections = response.sections.filter(
        (s: { example_text?: string }) => s.example_text !== undefined,
      );
      const exampleText =
        exampleSections
          .map((s: { example_text: string }) =>
            s.example_text.replace(/\n$/, ''),
          )
          .join('\n\n') + '\n';
      const optionSections = response.sections.filter(
        (s: { option_text?: string }) => s.option_text !== undefined,
      );
      const optionText =
        optionSections
          .map((s: { option_text: string }) => s.option_text.replace(/\n$/, ''))
          .join('\n\n') + '\n';
      const displaySections = response.sections.filter(
        (s: { id: string; example_text?: string; option_text?: string }) =>
          s.option_text === undefined || s.id === optionSections[0]?.id,
      );
      await expect(page.locator('[data-section-id]')).toHaveCount(
        displaySections.length,
      );
      expect(
        await page
          .locator('[data-section-id]')
          .evaluateAll((nodes) =>
            nodes.map((n) => n.getAttribute('data-section-id')),
          ),
      ).toEqual(displaySections.map((s: { id: string }) => s.id));
      await expect(page.locator('.prose pre')).toHaveCount(
        exampleSections.reduce(
          (count: number, s: { code: string[] }) => count + s.code.length,
          0,
        ) + (optionSections.length ? 1 : 0),
      );
      expect(
        await page.locator('.prose pre code.language-sh').allTextContents(),
      ).toEqual(exampleSections.flatMap((s: { code: string[] }) => s.code));
      if (optionSections.length) {
        expect(
          await page
            .locator('pre[data-copy-label="Copy options"] code')
            .textContent(),
        ).toBe(optionText);
        for (const line of optionText.split('\n'))
          expect(line.length).toBeLessThanOrEqual(100);
        const cliDetails = execFileSync(
          binary,
          [
            entry.id,
            '--implementation',
            implementation,
            '--details',
            '--width=0',
          ],
          { encoding: 'utf8' },
        );
        expect(cliDetails).toContain('Options\n' + optionText);
      }
      const cliExamples = execFileSync(
        binary,
        [entry.id, '--implementation', implementation, '--width=0'],
        { encoding: 'utf8' },
      );
      expect(cliExamples).toContain(exampleText);
      for (const section of response.sections) {
        await expect(page.locator(`[id="${section.id}"]`)).toHaveCount(1);
      }
      for (const section of displaySections) {
        const node = page.locator(`[data-section-id="${section.id}"]`);
        await expect(node.locator('h2')).toHaveText(
          section.option_text !== undefined ? 'Options' : section.title,
        );
        if (section.option_text !== undefined) continue;
        expect(await node.locator('pre code').allTextContents()).toEqual(
          section.code,
        );
        const actual = await node.locator('.prose').evaluate((el) => {
          const copy = el.cloneNode(true) as HTMLElement;
          copy.querySelectorAll('button').forEach((b) => b.remove());
          return copy.innerHTML;
        });
        // Go escapes a double quote as &quot;; innerHTML serializes it raw. Same HTML.
        expect(actual.replace(/\s+/g, ' ').trim()).toBe(
          section.html
            .replace(/&quot;/g, '"')
            .replace(/\s+/g, ' ')
            .trim(),
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

test('every command page shows its last updated date', async ({ page }) => {
  for (const entry of catalog.entries) {
    await page.goto(`/commands/${entry.id}/`);
    const updated = page.locator('.page-intro time');
    await expect(updated).toHaveText(entry.updated);
    await expect(updated).toHaveAttribute('datetime', entry.updated);
  }
});

test('copy preserves exact code and implementation links are shareable', async ({
  page,
  context,
}) => {
  await context.grantPermissions(['clipboard-read', 'clipboard-write']);
  await page.goto('/commands/sed/');
  await page.getByRole('link', { name: 'macOS / Apple' }).click();
  await expect(page).toHaveURL('/commands/sed/apple/');
  const buttons = page.getByRole('button', {
    name: 'Copy code example',
    exact: true,
  });
  const snippets = page.locator('pre code.language-sh');
  await expect(buttons).toHaveCount(codeBlockCount('sed', 'apple'));
  for (let index = 0; index < (await buttons.count()); index++) {
    const expected = await snippets.nth(index).textContent();
    await buttons.nth(index).click();
    expect(await page.evaluate(() => navigator.clipboard.readText())).toBe(
      expected,
    );
  }
  await expect(page.locator('.copy-status')).toContainText('Example copied');
  const options = await page
    .locator('pre[data-copy-label="Copy options"] code')
    .textContent();
  await page.getByRole('button', { name: 'Copy options', exact: true }).click();
  expect(await page.evaluate(() => navigator.clipboard.readText())).toBe(
    options,
  );
  await expect(page.locator('.copy-status')).toHaveText('Options copied.');
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
  baseURL,
}) => {
  const context = await browser.newContext({
    javaScriptEnabled: false,
    baseURL,
  });
  const page = await context.newPage();
  const external: string[] = [];
  page.on('request', (r) => {
    if (new URL(r.url()).origin !== new URL(baseURL!).origin)
      external.push(r.url());
  });
  await page.goto('/commands/find/gnu/');
  await expect(
    page.getByRole('heading', { name: 'Examples', exact: true }),
  ).toBeVisible();
  // Example blocks, plus the one grouped Options block.
  await expect(page.locator('pre code')).toHaveCount(
    codeBlockCount('find', 'gnu') + 1,
  );
  expect(
    (await page.locator('pre code.language-sh').allTextContents()).join('\n'),
  ).toContain('find . -maxdepth');
  await expect(
    page.getByRole('heading', { name: 'Options', exact: true }),
  ).toBeVisible();
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
