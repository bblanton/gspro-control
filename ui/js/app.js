(() => {
  const API_BASE = (() => {
    // Use current origin by default when http(s); fall back to localhost when opened from file://
    const url = new URL(window.location.href);
    const override = url.searchParams.get('api');
    const origin = window.location.origin;
    const defaultBase = origin.startsWith('http') ? origin : 'http://localhost:8080';
    return override || defaultBase;
  })();

  const statusEl = document.getElementById('connection-status');
  const actionsContainer = document.getElementById('actions');

  async function fetchActions() {
    try {
      const res = await fetch(`${API_BASE}/actions`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      return data; // object: { actionName: ["ctrl", "m"], ... }
    } catch (err) {
      console.error('Failed to fetch actions', err);
      throw err;
    }
  }

  function renderActions(actionsMap) {
    actionsContainer.innerHTML = '';
    const entries = Object.entries(actionsMap);
    if (entries.length === 0) {
      const empty = document.createElement('div');
      empty.className = 'col-span-full text-center text-slate-400';
      empty.textContent = 'No actions available';
      actionsContainer.appendChild(empty);
      return;
    }
    for (const [name, combo] of entries) {
      const btn = document.createElement('button');
      btn.className = [
        'h-24 rounded-xl',
        'bg-slate-800 hover:bg-slate-700 active:bg-slate-600',
        'transition-colors',
        'px-4 py-3 text-left',
        'shadow-md ring-1 ring-slate-700/50',
      ].join(' ');

      const label = document.createElement('div');
      label.className = 'text-lg font-medium';
      label.textContent = name.replace(/_/g, ' ');

      const sub = document.createElement('div');
      sub.className = 'mt-1 text-xs text-slate-400';
      sub.textContent = Array.isArray(combo) ? combo.join(' + ') : '';

      btn.appendChild(label);
      btn.appendChild(sub);

      btn.addEventListener('click', async () => {
        btn.disabled = true;
        const original = btn.className;
        try {
          const res = await fetch(`${API_BASE}/command?action=${encodeURIComponent(name)}`);
          if (!res.ok) throw new Error(`HTTP ${res.status}`);
          // Success pulse
          btn.className = original + ' ring-2 ring-emerald-400';
          setTimeout(() => (btn.className = original), 250);
        } catch (err) {
          console.error('Command failed', err);
          btn.className = original + ' ring-2 ring-rose-400';
          setTimeout(() => (btn.className = original), 500);
          alert(`Failed to run: ${name}`);
        } finally {
          btn.disabled = false;
        }
      });

      actionsContainer.appendChild(btn);
    }
  }

  async function init() {
    statusEl.textContent = 'Loading actions…';
    try {
      const actions = await fetchActions();
      statusEl.textContent = 'Connected';
      renderActions(actions);
    } catch {
      statusEl.textContent = 'Disconnected';
      renderActions({});
    }
  }

  window.addEventListener('load', init);
})();


