// Run the Playwright suite on a free port, so it never collides with a
// development server (for example `just web-dev` on Astro's default 4321).
import { spawn } from 'node:child_process';
import { once } from 'node:events';
import { createServer } from 'node:net';
import { fileURLToPath } from 'node:url';

const web = fileURLToPath(new URL('..', import.meta.url));

const socket = createServer();
socket.listen(0, '127.0.0.1');
await once(socket, 'listening');
const port = socket.address().port;
await new Promise((resolve) => socket.close(resolve));

console.log(`Browser tests use http://127.0.0.1:${port}`);
const child = spawn(
  process.execPath,
  [
    fileURLToPath(
      new URL('../node_modules/@playwright/test/cli.js', import.meta.url),
    ),
    'test',
    ...process.argv.slice(2),
  ],
  {
    cwd: web,
    stdio: 'inherit',
    env: { ...process.env, WEB_TEST_PORT: String(port) },
  },
);
const [code, signal] = await once(child, 'exit');
process.exit(code ?? (signal ? 1 : 0));
