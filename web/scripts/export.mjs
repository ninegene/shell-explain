import { spawnSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';
const root = fileURLToPath(new URL('../../', import.meta.url));
const result = spawnSync(
  'go',
  ['run', './cmd/content-export', '--output', 'web/src/data/catalog.json'],
  { cwd: root, stdio: 'inherit' },
);
if (result.error) console.error(result.error.message);
process.exit(result.status ?? 1);
