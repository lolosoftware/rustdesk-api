(() => {
  const proxyBaseValue = window.ws_host;
  const NativeWebSocket = window.WebSocket;

  if (!proxyBaseValue || !NativeWebSocket || NativeWebSocket.__rustdeskProxied) {
    return;
  }

  const proxyBase = new URL(proxyBaseValue, window.location.origin);
  const routeNames = {
    '21118': 'id',
    '21119': 'relay',
  };

  const route = (value) => {
    const original = new URL(String(value), window.location.href);
    const routeName = routeNames[original.port];
    if (!routeName) {
      return value;
    }

    const target = new URL(proxyBase.toString());
    target.protocol = ['https:', 'wss:'].includes(proxyBase.protocol) ? 'wss:' : 'ws:';
    target.pathname = `${target.pathname.replace(/\/$/, '')}/${routeName}`;
    target.search = '';
    target.hash = '';
    return target.toString();
  };

  function RoutedWebSocket (url, protocols) {
    const target = route(url);
    return protocols === undefined
      ? new NativeWebSocket(target)
      : new NativeWebSocket(target, protocols);
  }

  RoutedWebSocket.prototype = NativeWebSocket.prototype;
  Object.setPrototypeOf(RoutedWebSocket, NativeWebSocket);
  Object.defineProperty(RoutedWebSocket, '__rustdeskProxied', { value: true });
  window.WebSocket = RoutedWebSocket;
})();
