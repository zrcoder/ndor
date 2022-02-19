const cacheName = "app-" + "f56a09fa094e9f4516f570e25f8570ef30f1c44b";

self.addEventListener("install", event => {
  console.log("installing app worker f56a09fa094e9f4516f570e25f8570ef30f1c44b");

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
  console.log("app worker f56a09fa094e9f4516f570e25f8570ef30f1c44b is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
