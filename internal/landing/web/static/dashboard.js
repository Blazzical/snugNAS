(function () {
    const tiles = document.querySelectorAll('a.tile');

    // Tile hrefs use whatever hostname the user reached the dashboard on, so
    // a phone scanning the QR and a desktop browser both land on a reachable
    // URL. (Localhost from this PC, snugnas.local from a phone, etc.)
    tiles.forEach(t => {
        const host = window.location.hostname;
        const port = t.dataset.port;
        const path = t.dataset.path || '/';
        t.href = `http://${host}:${port}${path}`;
    });

    async function pollHealth() {
        let services;
        try {
            const r = await fetch('/api/health');
            if (!r.ok) return;
            services = await r.json();
        } catch (_) {
            return;
        }
        const byName = {};
        services.forEach(s => { byName[s.service] = s; });

        tiles.forEach(t => {
            const svc = t.dataset.service;
            if (!svc) return;
            const s = byName[svc];
            const dot = t.querySelector('.status-dot');
            if (!dot) return;

            if (!s) {
                dot.dataset.status = 'down';
                dot.title = 'not running';
                return;
            }
            if (s.state !== 'running') {
                dot.dataset.status = 'down';
                dot.title = s.state;
                return;
            }
            switch (s.health) {
                case 'unhealthy':
                    dot.dataset.status = 'down';
                    dot.title = 'unhealthy';
                    break;
                case 'starting':
                    dot.dataset.status = 'starting';
                    dot.title = 'starting';
                    break;
                case 'healthy':
                case '':
                case undefined:
                    dot.dataset.status = 'healthy';
                    dot.title = s.health || 'running (no healthcheck)';
                    break;
                default:
                    dot.dataset.status = 'starting';
                    dot.title = s.health;
            }
        });
    }

    pollHealth();
    setInterval(pollHealth, 5000);

    async function pollInfo() {
        try {
            const r = await fetch('/api/info');
            if (!r.ok) return;
            const i = await r.json();
            const set = (id, val) => {
                const el = document.getElementById(id);
                if (el) el.textContent = val;
            };
            set('info-version', i.version || '?');
            set('info-hostname', i.hostname || '?');
            set('info-uptime', i.uptime || '?');
        } catch (_) { /* ignore */ }
    }
    pollInfo();
    setInterval(pollInfo, 10000);
})();
