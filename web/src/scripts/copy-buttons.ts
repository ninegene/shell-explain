// Progressive enhancement shared by command and guide pages. Everything here
// is optional: without JavaScript the pages still show all text and code.

const copyStatus = () => document.querySelector('.copy-status');

// A diagram is read as one picture, described by the paragraph after it.
function describeDiagram(pre: HTMLPreElement) {
  const description = pre.nextElementSibling?.textContent?.trim();
  pre.setAttribute('role', 'img');
  pre.setAttribute('aria-label', description || 'Diagram');
}

for (const pre of document.querySelectorAll<HTMLPreElement>('.prose pre')) {
  const code = pre.querySelector('code');
  if (!code) continue;
  if (code.classList.contains('language-diagram')) {
    describeDiagram(pre);
    continue; // a picture is not something to copy
  }
  const button = document.createElement('button');
  button.className = 'copy-button';
  button.type = 'button';
  button.textContent = 'Copy';
  const options = pre.dataset.copyLabel === 'Copy options';
  button.setAttribute(
    'aria-label',
    options ? 'Copy options' : 'Copy code example',
  );
  button.addEventListener('click', async () => {
    const live = copyStatus()!;
    try {
      await navigator.clipboard.writeText(code.textContent ?? '');
      button.textContent = 'Copied';
      live.textContent = options
        ? 'Options copied.'
        : 'Example copied. Review it before running.';
    } catch {
      live.textContent =
        'Clipboard unavailable. Select and copy the example manually.';
    }
  });
  pre.append(button);
}

// A block that scrolls sideways must be reachable with the keyboard.
for (const el of document.querySelectorAll<HTMLElement>(
  '.prose pre, .prose table',
)) {
  if (el.scrollWidth > el.clientWidth) el.tabIndex = 0;
}

export {};
