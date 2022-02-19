const cacheName = "app-" + "ba6c5fb2c915f9c37ffaad4d9ca44ef5395defcb";

self.addEventListener("install", event => {
  console.log("installing app worker ba6c5fb2c915f9c37ffaad4d9ca44ef5395defcb");

  event.waitUntil(
    caches.open(cacheName).
      then(cache => {
        return cache.addAll([
          "/",
          "/app.css",
          "/app.js",
          "/manifest.webmanifest",
          "/wasm_exec.js",
          "/web/app.wasm",
          "/web/css/style.css",
          "/web/images/ellipse.png",
          "https://cdn.jsdelivr.net/npm/monaco-editor@0.31.1/min/vs/loader.js",
          "https://cdn.jsdelivr.net/npm/sweetalert2@11.3.10/dist/sweetalert2.all.min.js",
          
        ]);
      }).
      then(() => {
        self.skipWaiting();
      })
  );
});

self.addEventListener("activate", event => {
  event.waitUntil(
    caches.keys().then(keyList => {
      return Promise.all(
        keyList.map(key => {
          if (key !== cacheName) {
            return caches.delete(key);
          }
        })
      );
    })
  );
  console.log("app worker ba6c5fb2c915f9c37ffaad4d9ca44ef5395defcb is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
