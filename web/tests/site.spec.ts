import { test, expect, type Locator } from '@playwright/test';
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

test('search also finds guides and glossary words', async ({ page }) => {
  const errors: string[] = [];
  page.on('pageerror', (e) => errors.push(e.message));
  await page.goto('/commands/');
  const query = page.getByLabel('Find a command or a use case');
  // Nothing from the learning content shows until someone types.
  await expect(page.locator('[data-learn]:visible')).toHaveCount(0);
  await expect(page.locator('[data-learn-results]')).toBeHidden();

  await query.fill('stderr');
  // A guide is found by its own text and by the glossary words that link to it.
  const guideHits = catalog.guides.filter(
    (g: {
      id: string;
      title: string;
      subtitle: string;
      summary: string;
      tags: string[];
    }) =>
      [
        g.title,
        g.subtitle,
        g.summary,
        ...g.tags,
        ...catalog.glossary.terms
          .filter((t: { guides: string[] }) => t.guides.includes(g.id))
          .map((t: { term: string }) => t.term),
      ]
        .join(' ')
        .toLowerCase()
        .includes('stderr'),
  ).length;
  expect(guideHits).toBeGreaterThan(0);
  await expect(page.locator('[data-learn="guide"]:visible')).toHaveCount(
    guideHits,
  );
  const word = page.locator('[data-learn="term"]:visible a', {
    hasText: 'stderr',
  });
  await expect(word.first()).toHaveAttribute('href', '/learn/glossary/#stderr');
  await expect(
    page.locator('[data-learn="guide"]:visible a').first(),
  ).toHaveAttribute('href', /^\/learn\/[a-z-]+\/$/);

  // A one-letter query matches everything: only a few words are listed, with
  // a link to the full glossary.
  await query.fill('e');
  await expect(page.locator('[data-learn="term"]:visible')).toHaveCount(8);
  await expect(page.locator('#learn-more')).toBeVisible();

  // Choosing an implementation hides them: they are not implementation-specific.
  await query.fill('stderr');
  await page.getByLabel('Implementation coverage').selectOption('gnu');
  await expect(page.locator('[data-learn]:visible')).toHaveCount(0);
  await page.getByLabel('Implementation coverage').selectOption('');

  // The empty state appears only when nothing at all matches.
  await query.fill('zzzzqqqq');
  await expect(page.locator('#no-results')).toBeVisible();
  await query.fill('pipe');
  await expect(page.locator('#no-results')).toBeHidden();
  await expect(page.locator('[data-learn]:visible').first()).toBeVisible();

  // A tag page lists commands for that tag only.
  await page.goto('/tags/logs/');
  await expect(page.locator('[data-learn-results]')).toHaveCount(0);
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
        expect(normalizeHTML(await renderedHTML(node))).toBe(
          normalizeHTML(section.html),
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

// Go escapes a double quote as &quot;; innerHTML serializes it raw. Same HTML.
const normalizeHTML = (s: string) =>
  s
    .replace(/&quot;/g, '"')
    .replace(/\s+/g, ' ')
    .trim();

// The rendered prose of a section, without what the page script adds.
async function renderedHTML(node: Locator) {
  return node.locator('.prose').evaluate((el) => {
    const copy = el.cloneNode(true) as HTMLElement;
    copy.querySelectorAll('button').forEach((b) => b.remove());
    copy
      .querySelectorAll('pre, table')
      .forEach((e) =>
        ['tabindex', 'role', 'aria-label'].forEach((a) => e.removeAttribute(a)),
      );
    return copy.innerHTML;
  });
}

for (const guide of catalog.guides) {
  test(`CLI/web parity: guide ${guide.id}`, async ({ page }) => {
    const cli = JSON.parse(
      execFileSync(binary, ['learn', guide.id, '--format', 'json'], {
        encoding: 'utf8',
      }),
    );
    expect(cli.content_revision).toBe(catalog.revision);
    await page.goto(`/learn/${guide.id}/`);
    await expect(page.locator('.terminal-note code')).toHaveText(
      `explain learn ${guide.id}`,
    );
    await expect(page.locator('h1')).toHaveText(cli.guide.title);
    expect(
      await page
        .locator('[data-section-id]')
        .evaluateAll((nodes) =>
          nodes.map((n) => n.getAttribute('data-section-id')),
        ),
    ).toEqual(cli.guide.sections.map((s: { id: string }) => s.id));
    for (const section of cli.guide.sections) {
      const node = page.locator(`[data-section-id="${section.id}"]`);
      await expect(node.locator('h2')).toHaveText(section.title);
      expect(await node.locator('pre code').allTextContents()).toEqual(
        section.code,
      );
      expect(normalizeHTML(await renderedHTML(node))).toBe(
        normalizeHTML(section.html),
      );
    }
    // Every diagram is described in words and is not something to copy.
    const diagrams = page.locator('pre:has(code.language-diagram)');
    for (const pre of await diagrams.all()) {
      await expect(pre.locator('button')).toHaveCount(0);
      await expect(pre).toHaveAttribute('role', 'img');
      expect(await pre.getAttribute('aria-label')).toMatch(/^Diagram: /);
    }
    // Every block that is not a diagram keeps its own copy button.
    await expect(page.locator('.prose pre button.copy-button')).toHaveCount(
      await page.locator('.prose pre:not(:has(code.language-diagram))').count(),
    );
  });
}

test('CLI/web parity: glossary terms and anchors', async ({ page }) => {
  const cli = JSON.parse(
    execFileSync(binary, ['glossary', '--format', 'json'], {
      encoding: 'utf8',
    }),
  ).glossary;
  await page.goto('/learn/glossary/');
  await expect(page.locator('h1')).toHaveText(cli.title);
  expect(
    await page
      .locator('[data-term-id]')
      .evaluateAll((nodes) => nodes.map((n) => n.getAttribute('data-term-id'))),
  ).toEqual(cli.terms.map((t: { id: string }) => t.id));
  for (const term of cli.terms) {
    const node = page.locator(`[id="${term.id}"]`);
    await expect(node).toHaveCount(1);
    await expect(node.locator('h2')).toHaveText(term.term);
    expect(normalizeHTML(await renderedHTML(node))).toBe(
      normalizeHTML(term.html),
    );
    for (const id of term.guides)
      await expect(
        node.locator(`.term-more a[href="/learn/${id}/"]`),
      ).toHaveCount(1);
  }
  await page.goto('/learn/glossary/#operating-system');
  await expect(page.locator('#operating-system')).toBeInViewport();
});

test('learn index, home page start section, navigation and learn-more links', async ({
  page,
}) => {
  const active = catalog.series.filter((s: { id: string }) =>
    catalog.guides.some((g: { series: string }) => g.series === s.id),
  );
  await page.goto('/learn/');
  await expect(page.locator('.series-section h2').first()).toHaveText(
    active[0].title,
  );
  await expect(page.locator('.guide-card')).toHaveCount(catalog.guides.length);
  for (const [i, guide] of catalog.guides.entries()) {
    const card = page.locator('.guide-card').nth(i);
    await expect(card.locator('a')).toHaveAttribute(
      'href',
      `/learn/${guide.id}/`,
    );
    await expect(card.locator('.guide-number')).toHaveText(String(i + 1));
  }
  await expect(
    page.getByRole('link', { name: 'Words You Will See' }),
  ).toHaveAttribute('href', '/learn/glossary/');
  await page.goto('/');
  await expect(
    page
      .getByRole('navigation', { name: 'Main navigation' })
      .getByRole('link', { name: 'Learn' }),
  ).toHaveAttribute('href', '/learn/');
  await expect(page.locator('#start-here .series-card')).toHaveCount(
    active.length,
  );
  const first = catalog.guides.find(
    (g: { series: string }) => g.series === active[0].id,
  );
  await expect(
    page.locator('#start-here .series-card a').first(),
  ).toHaveAttribute('href', `/learn/${first.id}/`);
  await page.getByRole('link', { name: 'Browse all guides' }).click();
  await expect(page).toHaveURL(/\/learn\/$/);
  // A command page lists the guides it names, and only those.
  await page.goto('/commands/ls/');
  const related = page.getByRole('navigation', { name: 'Related guides' });
  await expect(related.getByRole('link')).toHaveText(
    catalog.guides.find((g: { id: string }) => g.id === 'big-picture').title,
  );
  await expect(related.getByRole('link')).toHaveAttribute(
    'href',
    '/learn/big-picture/',
  );
  await page.goto('/commands/grep/');
  await expect(
    page.getByRole('navigation', { name: 'Related guides' }),
  ).toHaveCount(0);
});

test('guide pages: series navigation, related links, tables and mobile layout', async ({
  page,
}) => {
  const guide = catalog.guides.find(
    (g: { id: string }) => g.id === 'big-picture',
  );
  await page.goto(`/learn/${guide.id}/`);
  await expect(
    page
      .getByRole('navigation', { name: 'Guides in this series' })
      .getByRole('link', { name: guide.title }),
  ).toHaveAttribute('aria-current', 'page');
  const index = catalog.guides.indexOf(guide);
  await expect(page.locator('.guide-nav-next')).toHaveCount(
    index < catalog.guides.length - 1 ? 1 : 0,
  );
  await expect(page.locator('.guide-nav-prev')).toHaveCount(index > 0 ? 1 : 0);
  for (const name of guide.related_commands)
    await expect(
      page.locator(`#related a[href="/commands/${name}/"]`),
    ).toHaveCount(1);
  await expect(page.locator('.prose table')).toHaveCount(1);
  await expect(page.locator('.prose th').first()).toHaveText('Layer');
  // On a phone the page does not scroll sideways; wide blocks scroll inside.
  await page.setViewportSize({ width: 375, height: 800 });
  await page.goto(`/learn/${guide.id}/`);
  expect(
    await page.evaluate(
      () => document.documentElement.scrollWidth <= innerWidth,
    ),
  ).toBe(true);
  const diagram = page.locator('pre:has(code.language-diagram)');
  expect(await diagram.evaluate((el) => el.scrollWidth > el.clientWidth)).toBe(
    true,
  );
  await expect(diagram).toHaveAttribute('tabindex', '0');
});

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
    for (const path of [
      '/',
      '/commands/grep/gnu/',
      '/commands/ls/',
      '/tags/text/',
      '/about/',
      '/learn/',
      '/learn/big-picture/',
      '/learn/glossary/',
    ]) {
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
  // The page with guide and word results showing is scanned too.
  for (const size of [
    { width: 1440, height: 1000 },
    { width: 390, height: 844 },
  ]) {
    await page.setViewportSize(size);
    await page.goto('/commands/?q=stderr');
    await expect(page.locator('[data-learn]:visible').first()).toBeVisible();
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
  await page.goto('/learn/big-picture/');
  await expect(
    page.getByRole('heading', { name: 'Goal', exact: true }),
  ).toBeVisible();
  expect(
    await page.locator('pre code.language-diagram').textContent(),
  ).toContain('KEYBOARD and SCREEN');
  await expect(page.locator('.prose table')).toHaveCount(1);
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
