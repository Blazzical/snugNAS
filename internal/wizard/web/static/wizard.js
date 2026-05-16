(async function () {
    const form = document.getElementById('wizard-form');
    const formSection = document.getElementById('form-section');
    const provisionSection = document.getElementById('provision-section');
    const log = document.getElementById('log');
    const provisionTitle = document.getElementById('provision-title');
    const provisionError = document.getElementById('provision-error');
    const dashboardLink = document.getElementById('dashboard-link');
    const formError = document.getElementById('form-error');

    // Pre-fill from existing config if any.
    try {
        const r = await fetch('/api/defaults');
        if (r.ok) {
            const d = await r.json();
            if (d.storage_root) document.getElementById('storage_root').value = d.storage_root;
            if (d.hostname) document.getElementById('hostname').value = d.hostname;
            if (d.admin_email) document.getElementById('admin_email').value = d.admin_email;
        }
    } catch (_) { /* ignore */ }

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        formError.hidden = true;
        const submit = form.querySelector('button[type="submit"]');
        submit.disabled = true;

        const payload = {
            storage_root: document.getElementById('storage_root').value.trim(),
            hostname: document.getElementById('hostname').value.trim(),
            admin_email: document.getElementById('admin_email').value.trim(),
        };

        const r = await fetch('/api/commit', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
        });
        if (!r.ok) {
            const msg = await r.text();
            formError.textContent = msg || `Request failed (${r.status}).`;
            formError.hidden = false;
            submit.disabled = false;
            return;
        }

        formSection.hidden = true;
        provisionSection.hidden = false;
        poll();
    });

    async function poll() {
        let lastLen = 0;
        while (true) {
            let st;
            try {
                const r = await fetch('/api/state');
                st = await r.json();
            } catch (_) {
                await sleep(1500);
                continue;
            }

            // Append any new log lines.
            for (let i = lastLen; i < st.logs.length; i++) {
                log.textContent += st.logs[i] + '\n';
            }
            lastLen = st.logs.length;
            log.scrollTop = log.scrollHeight;

            if (st.phase === 'done') {
                provisionTitle.textContent = 'snugNAS is ready';
                dashboardLink.hidden = false;
                if (st.post_up_url) dashboardLink.href = st.post_up_url;
                setTimeout(() => { window.location.href = dashboardLink.href; }, 3000);
                return;
            }
            if (st.phase === 'failed') {
                provisionTitle.textContent = 'Setup failed';
                provisionError.textContent = st.error || 'Unknown error';
                provisionError.hidden = false;
                return;
            }
            await sleep(1500);
        }
    }

    function sleep(ms) { return new Promise(r => setTimeout(r, ms)); }
})();
