(() => {
  const style = document.createElement('style');
  style.textContent = `
html.rustdesk-remote-cursor,
html.rustdesk-remote-cursor * {
  cursor: var(--rustdesk-remote-cursor) !important;
}`;
  document.head.appendChild(style);

  const resetCursor = () => {
    const root = document.documentElement;
    root.classList.remove('rustdesk-remote-cursor');
    root.style.removeProperty('--rustdesk-remote-cursor');
  };

  const applyCursor = (value) => {
    if (!value || value === 'auto') {
      resetCursor();
      return;
    }

    try {
      const cursor = JSON.parse(value);
      const url = String(cursor.url || '').replace(
        /^data:image\/rgba;base64,/,
        'data:image/png;base64,'
      );
      if (!url.startsWith('data:image/png;base64,')) {
        throw new Error('Unsupported cursor image');
      }
      const hotx = Number.parseInt(cursor.hotx, 10) || 0;
      const hoty = Number.parseInt(cursor.hoty, 10) || 0;
      const root = document.documentElement;
      root.style.setProperty(
        '--rustdesk-remote-cursor',
        `url("${url}") ${hotx} ${hoty}, auto`
      );
      root.classList.add('rustdesk-remote-cursor');
    } catch (error) {
      console.error('Failed to apply the remote cursor', error);
      resetCursor();
    }
  };

  let attempts = 0;
  const installBridge = () => {
    const original = window.setByName;
    if (typeof original !== 'function') {
      if (attempts++ < 1200) {
        window.setTimeout(installBridge, 25);
      }
      return;
    }
    if (original.__rustdeskCursorBridge) {
      return;
    }

    function bridgedSetByName (name, value) {
      if (name === 'cursor') {
        applyCursor(value);
        return;
      }
      return original.apply(this, arguments);
    }

    Object.defineProperty(bridgedSetByName, '__rustdeskCursorBridge', {
      value: true,
    });
    window.setByName = bridgedSetByName;
  };

  installBridge();
})();
