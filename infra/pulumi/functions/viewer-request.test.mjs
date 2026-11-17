// Run from infra/pulumi: npm test
import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';

// Load the function the way CloudFront does: plain source text, no module system.
const source = fs.readFileSync(
  new URL('./viewer-request.js', import.meta.url),
  'utf8',
);
const handler = new Function(`${source}; return handler;`)();

const run = (uri, { querystring = {} } = {}) =>
  handler({ request: { uri, querystring, headers: {} } });

test('the root maps to index.html', () => {
  assert.equal(run('/').uri, '/index.html');
});

test('a directory URL maps to its index.html', () => {
  assert.equal(run('/commands/').uri, '/commands/index.html');
  assert.equal(run('/commands/grep/').uri, '/commands/grep/index.html');
});

test('a missing trailing slash redirects with 301 and keeps the query', () => {
  const response = run('/commands');
  assert.equal(response.statusCode, 301);
  assert.equal(response.headers.location.value, '/commands/');

  const withQuery = run('/commands', {
    querystring: {
      a: { value: '1' },
      t: { value: 'x', multiValue: [{ value: 'x' }, { value: 'y' }] },
    },
  });
  assert.equal(withQuery.headers.location.value, '/commands/?a=1&t=x&t=y');
});

test('the trailing-slash redirect never points to another host', () => {
  for (const [uri, location] of [
    ['//evil', '/evil/'],
    ['//evil.com/page', '/evil.com/page/'],
    ['/\\evil', '/evil/'],
    ['//\\evil', '/evil/'],
  ]) {
    assert.equal(run(uri).headers.location.value, location, uri);
  }
});

test('files with an extension pass through unchanged', () => {
  for (const uri of [
    '/favicon.svg',
    '/404.html',
    '/_astro/index.abc123.css',
    '/_astro/CommandPage.abc123.js',
  ]) {
    assert.equal(run(uri).uri, uri);
  }
});
