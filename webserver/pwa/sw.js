const CACHE_NAME = 'bubble-cache-v1';
const ASSETS = [
    '/',
    '/sw.js'
];

// 1. Evenimentul de Instalare: Se creează cache-ul și se salvează rutele de bază
self.addEventListener('install', (event) => {
    event.waitUntil(
        caches.open(CACHE_NAME).then((cache) => {
            return cache.addAll(ASSETS);
        })
    );
});

// 2. Evenimentul de Activare: Se curăță cache-urile vechi dacă schimbi versiunea
self.addEventListener('activate', (event) => {
    event.waitUntil(
        caches.keys().then((keys) => {
            return Promise.all(
                keys.map((key) => {
                    if (key !== CACHE_NAME) {
                        return caches.delete(key);
                    }
                })
            );
        })
    );
});

// 3. Evenimentul Fetch: Interceptează cererile. Dacă rețeaua pică, servește din Cache
self.addEventListener('fetch', (event) => {
    event.respondWith(
        fetch(event.request)
            .then((response) => {
                // Dacă rețeaua e disponibilă, facem o copie a răspunsului în cache
                if (event.request.method === 'GET' && response.status === 200) {
                    const responseClone = response.clone();
                    caches.open(CACHE_NAME).then((cache) => {
                        cache.put(event.request, responseClone);
                    });
                }
                return response;
            })
            .catch(() => {
                // Dacă ești offline (Network Fail), caută în Cache
                return caches.match(event.request).then((cachedResponse) => {
                    if (cachedResponse) {
                        return cachedResponse;
                    }
                    // Dacă nu e în cache ruta exactă (ex: un ID nou), trimitem rădăcina "/" ca fallback
                    if (event.request.mode === 'navigate') {
                        return caches.match('/');
                    }
                });
            })
    );
});