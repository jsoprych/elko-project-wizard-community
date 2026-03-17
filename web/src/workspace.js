// Workspace: renders the right-hand panel based on active section
import { navigate, refresh } from './app.js';

export function renderWorkspace(state) {
  const el = document.getElementById('workspace-content');
  switch (state.activeSection) {
    case 'directive':     return renderDirective(el, state);
    case 'profile':       return renderProfile(el, state);
    case 'new-profile':   return renderNewProfile(el, state);
    case 'generations':   return renderGenerations(el);
    default:              return renderHome(el, state);
  }
}

function renderHome(el, state) {
  const numDirs  = (state.directives||[]).length;
  const numCats  = countCats(state.directives);
  const numProfs = (state.profiles||[]).length;

  el.innerHTML = `
    <!-- HERO -->
    <div class="hero">
      <div class="hero-eyebrow">⚡ elko-project-wizard</div>
      <h1 class="hero-title">Stop briefing your AI agent<br/>from a blank page.</h1>
      <p class="hero-sub">
        Pick your directives. Compose a profile. Hit <strong>Generate</strong>.<br/>
        Get a <code>.zip</code> with <code>AGENTS.md</code>, <code>CLAUDE.md</code>, <code>Dockerfile</code>, and more —
        ready to <code>cd</code> into and start coding in under 60 seconds.
      </p>
      <div class="hero-actions">
        <button class="btn btn-primary btn-lg" id="hero-new-profile">+ New Profile</button>
        <button class="btn btn-secondary btn-lg" onclick="import('./app.js').then(m=>m.navigate('generations'))">View History</button>
      </div>
    </div>

    <!-- LIVE STATS -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-number">${numDirs}</div>
        <div class="stat-label">Built-in Directives</div>
      </div>
      <div class="stat-card">
        <div class="stat-number">${numCats}</div>
        <div class="stat-label">Categories</div>
      </div>
      <div class="stat-card">
        <div class="stat-number">${numProfs}</div>
        <div class="stat-label">Saved Profiles</div>
      </div>
      <div class="stat-card">
        <div class="stat-number">60s</div>
        <div class="stat-label">From Zero to Coding</div>
      </div>
    </div>

    <!-- HOW IT WORKS -->
    <div class="section-header mt-16">How It Works</div>
    <div class="flow-row">
      <div class="flow-step">
        <div class="flow-icon">📋</div>
        <div class="flow-num">1</div>
        <div class="flow-title">Browse Directives</div>
        <div class="flow-desc">Pick from ${numDirs} built-in policies across coding standards, tech stack, visibility, Docker, and workflow.</div>
      </div>
      <div class="flow-arrow">→</div>
      <div class="flow-step">
        <div class="flow-icon">🧩</div>
        <div class="flow-num">2</div>
        <div class="flow-title">Compose a Profile</div>
        <div class="flow-desc">Combine directives into a reusable profile — <em>Go + Minimal Deps + Docker + Dual Edition</em> — saved for every future project.</div>
      </div>
      <div class="flow-arrow">→</div>
      <div class="flow-step">
        <div class="flow-icon">⬇</div>
        <div class="flow-num">3</div>
        <div class="flow-title">Generate &amp; Go</div>
        <div class="flow-desc">Download a <code>.zip</code> scaffold. Unzip, <code>git init</code>, open Claude Code or OpenCode — your agent is already briefed.</div>
      </div>
    </div>

    <!-- FEATURE HIGHLIGHTS -->
    <div class="section-header mt-16">What You Get</div>
    <div class="feature-grid">
      <div class="feature-card">
        <div class="feature-icon">🤖</div>
        <div class="feature-title">AGENTS.md + CLAUDE.md</div>
        <div class="feature-desc">One source of truth. <code>AGENTS.md</code> is the full directive doc. <code>CLAUDE.md</code> is a one-line <code>@AGENTS.md</code> shim — every AI agent stays in sync, zero drift.</div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🐳</div>
        <div class="feature-title">Docker, Done Right</div>
        <div class="feature-desc">Multi-stage Dockerfile, minimal runtime image, docker-compose with health checks and named volumes — generated for your exact stack.</div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">✨</div>
        <div class="feature-title">AI Assist</div>
        <div class="feature-desc">Need a custom directive? Get a perfectly-structured prompt to paste into Claude, ChatGPT, Grok, or DeepSeek. Or set <code>ANTHROPIC_API_KEY</code> and skip the copy-paste entirely.</div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🔒</div>
        <div class="feature-title">Dual Edition Ready</div>
        <div class="feature-desc">Select the Dual Edition directive and your generated project comes with <code>publish-community.sh</code> — a one-command private→public pruned release pipeline.</div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🧪</div>
        <div class="feature-title">Tests Before Every Commit</div>
        <div class="feature-desc">The Test-First directive drops a pre-commit hook into your project — <code>go build + vet + test</code> runs automatically before every commit. No broken builds.</div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">⚡</div>
        <div class="feature-title">Zero Bloat</div>
        <div class="feature-desc">Single Go binary. Pure Go SQLite (no CGO). Vanilla JS UI. No frameworks, no runtimes, no surprise dependencies. Runs anywhere Docker runs.</div>
      </div>
    </div>

    <!-- ALL DIRECTIVES -->
    <div class="section-header mt-16">All Directives — click to explore</div>
    <div class="directive-grid">
      ${(state.directives||[]).map(d => `
        <div class="directive-card" onclick="import('./app.js').then(m=>m.navigate('directive','${d.id}'))">
          <div class="d-name">${d.name}</div>
          <div class="d-desc">${d.description||''}</div>
          <div class="d-tags">${(d.tags||[]).map(t=>`<span class="tag">${t}</span>`).join('')}</div>
        </div>`).join('')}
    </div>`;

  document.getElementById('hero-new-profile')?.addEventListener('click', () => navigate('new-profile'));
}

function renderDirective(el, state) {
  const d = (state.directives||[]).find(x => x.id === state.selectedProfile);
  if (!d) return el.innerHTML = '<div class="empty-state">Directive not found</div>';
  el.innerHTML = `
    <div class="page-header">
      <div><div class="page-title">${d.name}</div>
        <div class="page-subtitle">${d.category} · ${d.builtin ? '🔒 built-in' : '✏️ custom'}</div>
      </div>
    </div>
    <div class="section-header">Description</div>
    <p class="text-muted">${d.description||''}</p>
    <div class="section-header mt-16">Tags</div>
    <div class="flex gap-8">${(d.tags||[]).map(t=>`<span class="tag">${t}</span>`).join('')}</div>
    <div class="section-header mt-16">Content Preview</div>
    <div class="preview-box">${escHtml(d.content||'')}</div>
    ${!d.builtin ? `<div class="mt-16"><button class="btn btn-danger btn-sm" onclick="deleteDirective('${d.id}')">Delete Directive</button></div>` : ''}`;

  window.deleteDirective = async (id) => {
    if (!confirm('Delete this directive?')) return;
    await fetch(`/api/directives/${id}`, {method:'DELETE'});
    await refresh();
  };
}

function renderProfile(el, state) {
  const p = (state.profiles||[]).find(x => x.id === state.selectedProfile);
  if (!p) return el.innerHTML = '<div class="empty-state">Profile not found</div>';

  const dirMap = Object.fromEntries((state.directives||[]).map(d=>[d.id,d]));
  el.innerHTML = `
    <div class="page-header">
      <div><div class="page-title">${p.name}</div>
        <div class="page-subtitle">${p.description||''}</div>
      </div>
      <div class="flex gap-8">
        <button class="btn btn-primary" id="btn-generate">⬇ Generate ZIP</button>
        <button class="btn btn-secondary" id="btn-ai-assist">🤖 AI Assist</button>
        <button class="btn btn-danger btn-sm" id="btn-delete-profile">Delete</button>
      </div>
    </div>
    <div class="section-header">Directives (${(p.directive_ids||[]).length})</div>
    <div class="directive-grid">
      ${(p.directive_ids||[]).map(id => {
        const d = dirMap[id];
        return d ? `<div class="directive-card selected"><div class="d-name">${d.name}</div><div class="d-desc">${d.description||''}</div></div>`
                 : `<div class="directive-card"><div class="d-name text-muted">${id} (not found)</div></div>`;
      }).join('')}
    </div>
    <div id="generate-result" class="mt-16"></div>
    <div id="ai-result" class="mt-16"></div>`;

  document.getElementById('btn-generate').onclick = async () => {
    const projectName = prompt('Project name:', p.project_name || 'my-project');
    if (!projectName) return;
    const res = await fetch('/api/generate', {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify({profile_id: p.id, project_name: projectName}),
    });
    if (!res.ok) {
      document.getElementById('generate-result').innerHTML = `<div class="text-error">Error: ${await res.text()}</div>`;
      return;
    }
    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url; a.download = projectName+'.zip'; a.click();
    URL.revokeObjectURL(url);
    document.getElementById('generate-result').innerHTML =
      `<div class="text-success">✓ ${projectName}.zip downloaded</div>`;
  };

  document.getElementById('btn-ai-assist').onclick = async () => {
    const req = prompt('Describe the directive you need (e.g. "enforce 12-factor principles"):');
    if (!req) return;
    const res = await fetch('/api/prompt', {
      method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify({requirement: req, project_type: p.description, tech_stack: 'go', category: 'workflow'}),
    });
    const data = await res.json();
    const aiEl = document.getElementById('ai-result');
    if (data.mode === 'sandwich') {
      aiEl.innerHTML = `<div class="sandwich-box"><div class="sandwich-hint">📋 Paste this into Claude / ChatGPT / Grok / DeepSeek:</div><div class="preview-box">${escHtml(data.prompt)}</div><button class="btn btn-secondary btn-sm mt-8" onclick="navigator.clipboard.writeText(${JSON.stringify(data.prompt)})">Copy to Clipboard</button></div>`;
    } else {
      aiEl.innerHTML = `<div class="section-header">Generated Directive</div><div class="preview-box">${escHtml(data.result)}</div>`;
    }
  };

  document.getElementById('btn-delete-profile').onclick = async () => {
    if (!confirm(`Delete profile "${p.name}"?`)) return;
    await fetch(`/api/profiles/${p.id}`, {method:'DELETE'});
    await refresh();
    navigate('directives');
  };
}

function renderNewProfile(el, state) {
  const cats = groupByCategory(state.directives||[]);
  const selected = new Set();

  el.innerHTML = `
    <div class="page-header">
      <div><div class="page-title">New Profile</div>
        <div class="page-subtitle">Select directives to compose, then save</div>
      </div>
    </div>
    <div class="form-group">
      <label class="form-label">Profile Name</label>
      <input id="prof-name" class="form-input" placeholder="e.g. Go Private Greenfield"/>
    </div>
    <div class="form-group">
      <label class="form-label">Default Project Name</label>
      <input id="prof-project" class="form-input" placeholder="e.g. my-api"/>
    </div>
    <div class="form-group">
      <label class="form-label">Description</label>
      <input id="prof-desc" class="form-input" placeholder="Brief description"/>
    </div>
    <div class="section-header">Select Directives</div>
    <div id="directive-picker"></div>
    <div class="mt-16 flex gap-8">
      <button id="btn-save-profile" class="btn btn-primary">Save Profile</button>
      <button class="btn btn-secondary" onclick="import('./app.js').then(m=>m.navigate('directives'))">Cancel</button>
    </div>`;

  // Render directive picker grouped by category
  const picker = document.getElementById('directive-picker');
  Object.entries(cats).forEach(([cat, dirs]) => {
    const sec = document.createElement('div');
    sec.innerHTML = `<div class="section-header">${cat}</div><div class="directive-grid" id="cat-${cat}"></div>`;
    picker.appendChild(sec);
    const grid = sec.querySelector(`#cat-${cat}`);
    dirs.forEach(d => {
      const card = document.createElement('div');
      card.className = 'directive-card';
      card.innerHTML = `<div class="d-name">${d.name}</div><div class="d-desc">${d.description||''}</div>`;
      card.onclick = () => {
        if (selected.has(d.id)) { selected.delete(d.id); card.classList.remove('selected'); }
        else { selected.add(d.id); card.classList.add('selected'); }
      };
      grid.appendChild(card);
    });
  });

  document.getElementById('btn-save-profile').onclick = async () => {
    const name = document.getElementById('prof-name').value.trim();
    if (!name) { alert('Profile name required'); return; }
    const profile = {
      name,
      project_name: document.getElementById('prof-project').value.trim(),
      description:  document.getElementById('prof-desc').value.trim(),
      directive_ids: [...selected],
    };
    const res = await fetch('/api/profiles', {
      method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify(profile),
    });
    if (!res.ok) { alert('Error saving profile: ' + await res.text()); return; }
    const saved = await res.json();
    await refresh();
    navigate('profile', saved.id);
  };
}

async function renderGenerations(el) {
  const gens = await fetch('/api/generations').then(r=>r.json());
  el.innerHTML = `
    <div class="page-header"><div class="page-title">Generation History</div></div>
    ${(!gens || gens.length===0) ? '<div class="empty-state">No generations yet — create a profile and generate a project!</div>' :
      gens.map(g=>`<div class="card"><div class="card-title">${g.project_name}</div>
        <div class="card-meta">${g.created_at}</div>
        <div class="card-content text-muted">Profile: ${g.profile_id||'—'}</div></div>`).join('')}`;
}

// Utils
function groupByCategory(dirs) {
  return dirs.reduce((acc, d) => {
    (acc[d.category] = acc[d.category]||[]).push(d);
    return acc;
  }, {});
}
function countCats(dirs) { return new Set((dirs||[]).map(d=>d.category)).size; }
function escHtml(s) { return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); }
