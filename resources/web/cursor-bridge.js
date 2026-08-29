(() => {
  const style = document.createElement('style');
  style.textContent = `
html.rustdesk-cursor-fallback,
html.rustdesk-cursor-fallback * {
  cursor: default !important;
}
html.rustdesk-remote-cursor,
html.rustdesk-remote-cursor * {
  cursor: var(--rustdesk-remote-cursor) !important;
}`;
  document.head.appendChild(style);
  document.documentElement.classList.add('rustdesk-cursor-fallback');
  window.__rustdeskCursorBridgeStatus = {
    installed: false,
    cursorCommands: 0,
    customCursorActive: false,
    uppercaseTranslations: 0,
    twoFactorPrompts: 0,
    twoFactorCodesSent: 0,
  };

  let twoFactorPromptOpen = false;

  const requestTwoFactorCode = (showOriginalError, text) => {
    if (twoFactorPromptOpen) return;
    twoFactorPromptOpen = true;
    window.__rustdeskCursorBridgeStatus.twoFactorPrompts += 1;

    window.setTimeout(() => {
      const code = window.prompt(
        'Authentification à deux facteurs RustDesk\n\nSaisissez le code à 6 chiffres :',
        ''
      );
      twoFactorPromptOpen = false;

      if (code === null) {
        showOriginalError('error', 'Login Error', text);
        return;
      }
      const normalizedCode = code.trim();
      if (!/^\d{6}$/.test(normalizedCode)) {
        window.alert('Le code 2FA doit contenir exactement 6 chiffres.');
        requestTwoFactorCode(showOriginalError, text);
        return;
      }

      window.__rustdeskCursorBridgeStatus.twoFactorCodesSent += 1;
      window.setByName('send_2fa', JSON.stringify({ code: normalizedCode }));
    }, 0);
  };

  let bridgedConnection;
  const installConnectionBridge = () => {
    const connection = window.curConn;
    if (!connection || connection === bridgedConnection || typeof connection.msgbox !== 'function') {
      return;
    }

    const originalMsgbox = connection.msgbox.bind(connection);
    connection.msgbox = (type, title, message) => {
      if (
        type === 'input-2fa' ||
        message === '2FA Required' ||
        message === 'Wrong 2FA Code'
      ) {
        requestTwoFactorCode(originalMsgbox, message);
        return;
      }
      return originalMsgbox(type, title, message);
    };
    bridgedConnection = connection;
  };

  window.setInterval(installConnectionBridge, 25);

  const resetCursor = () => {
    const root = document.documentElement;
    root.classList.remove('rustdesk-remote-cursor');
    root.style.removeProperty('--rustdesk-remote-cursor');
    window.__rustdeskCursorBridgeStatus.customCursorActive = false;
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
      window.__rustdeskCursorBridgeStatus.customCursorActive = true;
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
      if (name === 'input_key') {
        try {
          const key = JSON.parse(value);
          const match = /^VK_([A-Z])$/.exec(key.name);
          if (key.shift === 'true' && match) {
            key.name = match[1];
            key.shift = 'false';
            value = JSON.stringify(key);
            arguments[1] = value;
            window.__rustdeskCursorBridgeStatus.uppercaseTranslations += 1;
          }
        } catch (error) {
          console.error('Failed to normalize the keyboard event', error);
        }
      }
      if (name === 'cursor') {
        window.__rustdeskCursorBridgeStatus.cursorCommands += 1;
        applyCursor(value);
        return;
      }
      return original.apply(this, arguments);
    }

    Object.defineProperty(bridgedSetByName, '__rustdeskCursorBridge', {
      value: true,
    });
    window.setByName = bridgedSetByName;
    window.__rustdeskCursorBridgeStatus.installed = true;
    console.info('RustDesk cursor bridge installed');
  };

  installBridge();
})();
