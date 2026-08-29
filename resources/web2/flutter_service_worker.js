// Replace the service worker installed by the incomplete Flutter build.
// Once activated, this worker removes its own registration and reloads any
// open Web Client V2 pages so that index.html can redirect them safely.
self.addEventListener('install', () => self.skipWaiting());

self.addEventListener('activate', (event) => {
  event.waitUntil((async () => {
    await self.registration.unregister();
    const windows = await self.clients.matchAll({ type: 'window' });
    await Promise.all(windows.map((client) => client.navigate(client.url)));
  })());
});
