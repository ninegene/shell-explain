import assert from 'node:assert/strict';
import { spawn, spawnSync } from 'node:child_process';
import { once } from 'node:events';
import {
  cp,
  mkdtemp,
  readdir,
  readFile,
  rm,
  symlink,
  writeFile,
  rename,
} from 'node:fs/promises';
import { createServer } from 'node:net';
import os from 'node:os';
import path from 'node:path';
import { setTimeout as delay } from 'node:timers/promises';
import { fileURLToPath } from 'node:url';
import test from 'node:test';

const root = fileURLToPath(new URL('../../', import.meta.url));

test(
  'Astro dev reflects Markdown edits/additions/removals and recovers from invalid content',
  { timeout: 90000 },
  async (t) => {
    const fixture = await mkdtemp(path.join(os.tmpdir(), 'shell-explain-dev-'));
    let server;
    t.after(async () => {
      if (server && server.exitCode === null && server.signalCode === null) {
        const exited = once(server, 'exit');
        server.kill('SIGTERM');
        await exited;
      }
      await rm(fixture, { recursive: true, force: true });
    });
    for (const name of [
      'go.mod',
      'go.sum',
      'cmd',
      'internal',
      'content',
      'web',
    ]) {
      await cp(path.join(root, name), path.join(fixture, name), {
        recursive: true,
        filter: (source) =>
          !['node_modules', 'dist', '.astro', 'test-results'].includes(
            path.basename(source),
          ),
      });
    }
    await symlink(
      path.join(root, 'web/node_modules'),
      path.join(fixture, 'web/node_modules'),
    );
    const initial = spawnSync(
      'go',
      [
        'run',
        './cmd/content-export',
        '--source',
        'content',
        '--output',
        'web/src/data/catalog.json',
      ],
      { cwd: fixture, encoding: 'utf8' },
    );
    assert.equal(initial.status, 0, initial.stderr);
    const socket = createServer();
    socket.listen(0, '127.0.0.1');
    await once(socket, 'listening');
    const port = socket.address().port;
    await new Promise((resolve) => socket.close(resolve));
    server = spawn(
      process.execPath,
      [
        path.join(root, 'web/node_modules/astro/bin/astro.mjs'),
        'dev',
        '--host',
        '127.0.0.1',
        '--port',
        String(port),
        '--ignore-lock',
      ],
      { cwd: path.join(fixture, 'web'), stdio: ['ignore', 'pipe', 'pipe'] },
    );
    let logs = '';
    server.stdout.on('data', (data) => {
      logs += data;
    });
    server.stderr.on('data', (data) => {
      logs += data;
    });
    const until = async (check) => {
      const deadline = Date.now() + 20000;
      while (Date.now() < deadline) {
        try {
          if (await check()) return;
        } catch {
          /* Server may be starting/reloading. */
        }
        await delay(100);
      }
      assert.fail(`Development site did not update.\n${logs}`);
    };
    const page = async (url = '/commands/grep/') =>
      fetch(`http://127.0.0.1:${port}${url}`);
    await until(async () => (await page()).ok);
    const source = path.join(fixture, 'content/en/grep.md');
    const original = await readFile(source, 'utf8');
    const updated = original.replace(
      'Find lines of text that match a pattern.',
      'Live Markdown edit regression marker.',
    );
    await writeFile(source, updated);
    await until(async () =>
      (await (await page()).text()).includes(
        'Live Markdown edit regression marker.',
      ),
    );
    const catalog = path.join(fixture, 'web/src/data/catalog.json');
    const lastGood = await readFile(catalog, 'utf8');
    await writeFile(source, '---\ninvalid content');
    await until(() => logs.includes('Content export failed'));
    assert.equal(await readFile(catalog, 'utf8'), lastGood);
    await writeFile(source, original);
    await until(async () =>
      (await (await page()).text()).includes(
        'Find lines of text that match a pattern.',
      ),
    );
    // These are valid shared documents: neither code-block count nor -n is required.
    for (const [label, body] of [
      [
        'one',
        "A single example:\n\n```sh\ngrep 'one-block-marker' app.log\n```\n",
      ],
      ['zero', 'Examples described in prose.\n'],
    ]) {
      const summary = `Homepage ${label} block regression marker.`;
      await writeFile(
        source,
        original
          .replace(/### Examples\n[\s\S]*?(?=\n### )/, `### Examples\n${body}`)
          .replace('### Option: -n', '### Line numbers')
          .replace(
            '## GNU',
            '## GNU\n### Option: -n\nVariant-only homepage regression marker.\n',
          )
          .replace('Find lines of text that match a pattern.', summary),
      );
      await until(async () => {
        const response = await page('/');
        const html = await response.text();
        return (
          response.ok &&
          html.includes(summary) &&
          !html.includes('Variant-only homepage regression marker.') &&
          (label === 'zero' || html.includes('one-block-marker'))
        );
      });
    }
    await writeFile(source, original);
    await until(async () =>
      (await (await page('/')).text()).includes(
        'Find lines of text that match a pattern.',
      ),
    );
    const extra = path.join(fixture, 'content/en/demo.md');
    await writeFile(
      extra,
      original
        .replace('id: grep', 'id: demo')
        .replace('title: grep', 'title: demo')
        .replace('aliases: [ggrep, egrep, fgrep]', 'aliases: []'),
    );
    await until(async () => (await page('/commands/demo/')).ok);
    // Guides name grep as a related command, and the export refuses a name with
    // no page. Drop it from the guides first, as an author removing the page would.
    const guides = path.join(fixture, 'content/en/guides');
    for (const name of await readdir(guides)) {
      const file = path.join(guides, name);
      const text = await readFile(file, 'utf8');
      await writeFile(
        file,
        text.replace(/^related_commands: \[(.*)\]$/m, (_, list) => {
          const kept = list
            .split(',')
            .map((s) => s.trim())
            .filter((s) => s && s !== 'grep');
          return `related_commands: [${kept.join(', ')}]`;
        }),
      );
    }
    await rm(source);
    await until(async () => (await page('/commands/grep/')).status === 404);
    // The homepage falls back to the first command alphabetically.
    // Only command pages count; guides and the glossary live in subdirectories.
    const fallback = (await readdir(path.join(fixture, 'content/en')))
      .filter((name) => name.endsWith('.md'))
      .map((name) => name.replace(/\.md$/, ''))
      .sort((a, b) => a.localeCompare(b))[0];
    await until(async () => {
      const response = await page('/');
      return (
        response.ok && (await response.text()).includes(`explain ${fallback}`)
      );
    });
    await writeFile(source, original);
    await until(async () => (await page('/commands/grep/')).ok);
    const renamed = path.join(fixture, 'content/en/renamed.md');
    await rename(extra, renamed);
    await until(async () => (await page('/commands/demo/')).ok);
    await rm(renamed);
    await until(async () => (await page('/commands/demo/')).status === 404);
    // Build only after releasing the development server's content-store lock.
    const exited = once(server, 'exit');
    server.kill('SIGTERM');
    await exited;
    await writeFile(
      source,
      original
        .replace(
          /### Examples\n[\s\S]*?(?=\n### )/,
          '### Examples\nExamples described in prose.\n',
        )
        .replace('### Option: -n', '### Line numbers'),
    );
    const exported = spawnSync(
      'go',
      [
        'run',
        './cmd/content-export',
        '--source',
        'content',
        '--output',
        'web/src/data/catalog.json',
      ],
      { cwd: fixture, encoding: 'utf8', timeout: 30000 },
    );
    assert.equal(exported.status, 0, exported.stdout + exported.stderr);
    const build = spawnSync(
      process.execPath,
      [
        path.join(root, 'web/node_modules/astro/bin/astro.mjs'),
        'build',
        '--ignore-lock',
      ],
      { cwd: path.join(fixture, 'web'), encoding: 'utf8', timeout: 30000 },
    );
    assert.equal(build.status, 0, build.stdout + build.stderr);
  },
);
