import { execFile } from 'node:child_process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root = fileURLToPath(new URL('../../', import.meta.url));
const content = path.join(root, 'content');

// Use Vite's watcher so Markdown outside web/ shares the dev server's lifecycle.
export function contentWatch() {
  return {
    name: 'shell-explain-content',
    apply: 'serve',
    configureServer(server) {
      let timer;
      let running = false;
      let pending = false;
      let closed = false;
      let child;
      const exportContent = () => {
        if (closed) return;
        if (running) {
          pending = true;
          return;
        }
        running = true;
        child = execFile(
          'go',
          [
            'run',
            './cmd/content-export',
            '--source',
            'content',
            '--output',
            'web/src/data/catalog.json',
          ],
          { cwd: root, timeout: 30000 },
          (error, _stdout, stderr) => {
            running = false;
            if (closed) return;
            if (error)
              server.config.logger.error(
                `Content export failed; keeping the last valid catalog.\n${stderr || error.message}`,
              );
            else server.config.logger.info('Markdown catalog updated.');
            if (pending) {
              pending = false;
              exportContent();
            }
          },
        );
      };
      const onFile = (file) => {
        const relative = path.relative(content, file);
        if (
          relative.startsWith('..') ||
          path.isAbsolute(relative) ||
          !file.endsWith('.md')
        )
          return;
        clearTimeout(timer);
        timer = setTimeout(exportContent, 100);
      };
      server.watcher.add(content);
      for (const event of ['add', 'change', 'unlink'])
        server.watcher.on(event, onFile);
      server.httpServer?.once('close', () => {
        closed = true;
        clearTimeout(timer);
        child?.kill();
        for (const event of ['add', 'change', 'unlink'])
          server.watcher.off(event, onFile);
      });
    },
  };
}
