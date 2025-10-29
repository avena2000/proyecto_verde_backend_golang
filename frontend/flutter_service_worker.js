'use strict';
const MANIFEST = 'flutter-app-manifest';
const TEMP = 'flutter-temp-cache';
const CACHE_NAME = 'flutter-app-cache';

const RESOURCES = {"canvaskit/canvaskit.js": "140ccb7d34d0a55065fbd422b843add6",
"canvaskit/canvaskit.js.symbols": "58832fbed59e00d2190aa295c4d70360",
"canvaskit/canvaskit.wasm": "07b9f5853202304d3b0749d9306573cc",
"canvaskit/chromium/canvaskit.js": "5e27aae346eee469027c80af0751d53d",
"canvaskit/chromium/canvaskit.js.symbols": "193deaca1a1424049326d4a91ad1d88d",
"canvaskit/chromium/canvaskit.wasm": "24c77e750a7fa6d474198905249ff506",
"canvaskit/skwasm.js": "1ef3ea3a0fec4569e5d531da25f34095",
"canvaskit/skwasm.js.symbols": "0088242d10d7e7d6d2649d1fe1bda7c1",
"canvaskit/skwasm.wasm": "264db41426307cfc7fa44b95a7772109",
"canvaskit/skwasm_heavy.js": "413f5b2b2d9345f37de148e2544f584f",
"canvaskit/skwasm_heavy.js.symbols": "3c01ec03b5de6d62c34e17014d1decd3",
"canvaskit/skwasm_heavy.wasm": "8034ad26ba2485dab2fd49bdd786837b",
"flutter.js": "888483df48293866f9f41d3d9274a779",
"flutter_bootstrap.js": "557072d92158df6e63d9df3c22ce3e74",
"index.html": "d515ceaf39d775f2db2ff8de14498ad8",
"/": "d515ceaf39d775f2db2ff8de14498ad8",
"main.dart.mjs": "9dd10c73463bcee7ffd985a9bc564550",
"main.dart.wasm": "373c624cb886f587629c28128b78a24d",
"main.dart.js": "28348bca023f48aba88dc0dd219fbfd1",
"version.json": "78e57b41b0759c9cbd2b58ac9589ce30",
"assets/assets/accesorios/barba/1.png": "4f3c985e2a97079a3fb97037e0eb1545",
"assets/assets/accesorios/barba/2.png": "988ca180caa76a0d853387737f3d2c31",
"assets/assets/accesorios/barba/3.png": "285a721362ef2f16d73055e80fd34cab",
"assets/assets/accesorios/barba/4.png": "4aa4a19a0b5108cbcb587b47484fedeb",
"assets/assets/accesorios/barba/5.png": "2c6234d97cf8cc5a28ae84d54ad9ddd3",
"assets/assets/accesorios/barba/6.png": "ffbd738c31d7710fb6849becc572d241",
"assets/assets/accesorios/cabello/1.png": "7d2cdee75b494d79a2c5ddb78afaa875",
"assets/assets/accesorios/cabello/10.png": "e193504c609d9151b6a877131dfcd97d",
"assets/assets/accesorios/cabello/2.png": "fb0e423f0a8618fcbd06d8a498c18fe0",
"assets/assets/accesorios/cabello/3.png": "2505935634f990d8aafdec850584c87b",
"assets/assets/accesorios/cabello/4.png": "563f7ec237ce80117caf0f7aab465c9e",
"assets/assets/accesorios/cabello/5.png": "ed6d7584bef989fce50fe2e555e93cf6",
"assets/assets/accesorios/cabello/6.png": "a160623724ea4690cdce116902d74b3e",
"assets/assets/accesorios/cabello/7.png": "a5fe290cdf2b7abbcac4e2e1a911d469",
"assets/assets/accesorios/cabello/8.png": "0a3a240cc2c82c247f83a75f2d7c9875",
"assets/assets/accesorios/cabello/9.png": "6ed01b39459c606b107c1d7e43c9f62b",
"assets/assets/accesorios/cabello/default.png": "32b971a9533b1ea82ee7203fbab1325e",
"assets/assets/accesorios/detalle_adicional/1.png": "444bf5daec6b4ef4553639ab7e264f51",
"assets/assets/accesorios/detalle_adicional/2.png": "3eb6281394eaa6a8fe86a6bed55a4068",
"assets/assets/accesorios/detalle_facial/1.png": "579d8654439657b57b18f2729bc7ebbd",
"assets/assets/accesorios/detalle_facial/2.png": "9952047eb81a6bc6458157d3d7e58391",
"assets/assets/accesorios/vestimenta/1.png": "5a363fc858f444c9d5d9db5ec14c17ee",
"assets/assets/accesorios/vestimenta/2.png": "54104ed0ffef7edf45d09ca984933a2a",
"assets/assets/accesorios/vestimenta/3.png": "7fea8bf89b073867933833ecb1c7e30b",
"assets/assets/accesorios/vestimenta/4.png": "821b88cb9d1850aeb779d456863c0616",
"assets/assets/accesorios/vestimenta/5.png": "bccee86d8c60e4c5a75a43a905600ccd",
"assets/assets/accesorios/vestimenta/default.png": "85c3e0780a2b0f6bd1ec59b3185a4127",
"assets/assets/google.png": "9940378f2de149de5a3c0946020f0bb7",
"assets/assets/github.png": "ec3a60c8c6539a07eb70b52f6737ea6e",
"assets/assets/facebook.png": "5c648c3c83f03bd089f5f71516b414db",
"assets/assets/personaje_base.png": "a9839f15c5ad60d0c482c6770a8beff1",
"assets/assets/acciones/accion_alerta.png": "f2a8ec1606c936eb83663f12d3621d52",
"assets/assets/acciones/accion_ayuda.png": "aee405fb117a18889fa285577c4f2a11",
"assets/assets/acciones/accion_descubrimiento.png": "eca0643ecc1c34947699f9b64bbb3fe3",
"assets/assets/logo.png": "85b5f602f06796d3f9af50ad530583cf",
"assets/assets/logo-text.png": "1e9855f6430c83a4de02df6d016308b3",
"assets/assets/fonts/YesevaOne-Regular.ttf": "5567d0bf3fe8eba4f85fbc611e8ff1ff",
"assets/packages/flutter_map/lib/assets/flutter_map_logo.png": "208d63cc917af9713fc9572bd5c09362",
"assets/packages/cupertino_icons/assets/CupertinoIcons.ttf": "33b7d9392238c04c131b6ce224e13711",
"assets/fonts/MaterialIcons-Regular.otf": "81db6b1feef369aca9dd79e96386779c",
"assets/shaders/ink_sparkle.frag": "ecc85a2e95f5e9f53123dcaf8cb9b6ce",
"assets/AssetManifest.json": "94d42339d145d6350927a3452928a5c9",
"assets/AssetManifest.bin": "56aaa8a2d86b5169f9041cf24046ca72",
"assets/AssetManifest.bin.json": "a07a884e9e027fabfe92f07aad469246",
"assets/FontManifest.json": "ab062c0635f18894e0291897857bff3b",
"assets/NOTICES": "7d35dbe347b33fb21d4aa998682165ce",
"favicon.png": "e317f32c004e6e7ea3c9acc88f5efc6c",
"icons/Icon-192.png": "9c91964b739ea5af83259aa35825cbf1",
"icons/Icon-512.png": "9c2d379c1c9c23f4e7a8a23dfb4aca26",
"icons/Icon-maskable-192.png": "9c91964b739ea5af83259aa35825cbf1",
"icons/Icon-maskable-512.png": "9c2d379c1c9c23f4e7a8a23dfb4aca26",
"manifest.json": "65ec0c22d73c0a8183f68222a91354da"};
// The application shell files that are downloaded before a service worker can
// start.
const CORE = ["main.dart.js",
"main.dart.wasm",
"main.dart.mjs",
"index.html",
"flutter_bootstrap.js",
"assets/AssetManifest.bin.json",
"assets/FontManifest.json"];

// During install, the TEMP cache is populated with the application shell files.
self.addEventListener("install", (event) => {
  self.skipWaiting();
  return event.waitUntil(
    caches.open(TEMP).then((cache) => {
      return cache.addAll(
        CORE.map((value) => new Request(value, {'cache': 'reload'})));
    })
  );
});
// During activate, the cache is populated with the temp files downloaded in
// install. If this service worker is upgrading from one with a saved
// MANIFEST, then use this to retain unchanged resource files.
self.addEventListener("activate", function(event) {
  return event.waitUntil(async function() {
    try {
      var contentCache = await caches.open(CACHE_NAME);
      var tempCache = await caches.open(TEMP);
      var manifestCache = await caches.open(MANIFEST);
      var manifest = await manifestCache.match('manifest');
      // When there is no prior manifest, clear the entire cache.
      if (!manifest) {
        await caches.delete(CACHE_NAME);
        contentCache = await caches.open(CACHE_NAME);
        for (var request of await tempCache.keys()) {
          var response = await tempCache.match(request);
          await contentCache.put(request, response);
        }
        await caches.delete(TEMP);
        // Save the manifest to make future upgrades efficient.
        await manifestCache.put('manifest', new Response(JSON.stringify(RESOURCES)));
        // Claim client to enable caching on first launch
        self.clients.claim();
        return;
      }
      var oldManifest = await manifest.json();
      var origin = self.location.origin;
      for (var request of await contentCache.keys()) {
        var key = request.url.substring(origin.length + 1);
        if (key == "") {
          key = "/";
        }
        // If a resource from the old manifest is not in the new cache, or if
        // the MD5 sum has changed, delete it. Otherwise the resource is left
        // in the cache and can be reused by the new service worker.
        if (!RESOURCES[key] || RESOURCES[key] != oldManifest[key]) {
          await contentCache.delete(request);
        }
      }
      // Populate the cache with the app shell TEMP files, potentially overwriting
      // cache files preserved above.
      for (var request of await tempCache.keys()) {
        var response = await tempCache.match(request);
        await contentCache.put(request, response);
      }
      await caches.delete(TEMP);
      // Save the manifest to make future upgrades efficient.
      await manifestCache.put('manifest', new Response(JSON.stringify(RESOURCES)));
      // Claim client to enable caching on first launch
      self.clients.claim();
      return;
    } catch (err) {
      // On an unhandled exception the state of the cache cannot be guaranteed.
      console.error('Failed to upgrade service worker: ' + err);
      await caches.delete(CACHE_NAME);
      await caches.delete(TEMP);
      await caches.delete(MANIFEST);
    }
  }());
});
// The fetch handler redirects requests for RESOURCE files to the service
// worker cache.
self.addEventListener("fetch", (event) => {
  if (event.request.method !== 'GET') {
    return;
  }
  var origin = self.location.origin;
  var key = event.request.url.substring(origin.length + 1);
  // Redirect URLs to the index.html
  if (key.indexOf('?v=') != -1) {
    key = key.split('?v=')[0];
  }
  if (event.request.url == origin || event.request.url.startsWith(origin + '/#') || key == '') {
    key = '/';
  }
  // If the URL is not the RESOURCE list then return to signal that the
  // browser should take over.
  if (!RESOURCES[key]) {
    return;
  }
  // If the URL is the index.html, perform an online-first request.
  if (key == '/') {
    return onlineFirst(event);
  }
  event.respondWith(caches.open(CACHE_NAME)
    .then((cache) =>  {
      return cache.match(event.request).then((response) => {
        // Either respond with the cached resource, or perform a fetch and
        // lazily populate the cache only if the resource was successfully fetched.
        return response || fetch(event.request).then((response) => {
          if (response && Boolean(response.ok)) {
            cache.put(event.request, response.clone());
          }
          return response;
        });
      })
    })
  );
});
self.addEventListener('message', (event) => {
  // SkipWaiting can be used to immediately activate a waiting service worker.
  // This will also require a page refresh triggered by the main worker.
  if (event.data === 'skipWaiting') {
    self.skipWaiting();
    return;
  }
  if (event.data === 'downloadOffline') {
    downloadOffline();
    return;
  }
});
// Download offline will check the RESOURCES for all files not in the cache
// and populate them.
async function downloadOffline() {
  var resources = [];
  var contentCache = await caches.open(CACHE_NAME);
  var currentContent = {};
  for (var request of await contentCache.keys()) {
    var key = request.url.substring(origin.length + 1);
    if (key == "") {
      key = "/";
    }
    currentContent[key] = true;
  }
  for (var resourceKey of Object.keys(RESOURCES)) {
    if (!currentContent[resourceKey]) {
      resources.push(resourceKey);
    }
  }
  return contentCache.addAll(resources);
}
// Attempt to download the resource online before falling back to
// the offline cache.
function onlineFirst(event) {
  return event.respondWith(
    fetch(event.request).then((response) => {
      return caches.open(CACHE_NAME).then((cache) => {
        cache.put(event.request, response.clone());
        return response;
      });
    }).catch((error) => {
      return caches.open(CACHE_NAME).then((cache) => {
        return cache.match(event.request).then((response) => {
          if (response != null) {
            return response;
          }
          throw error;
        });
      });
    })
  );
}
