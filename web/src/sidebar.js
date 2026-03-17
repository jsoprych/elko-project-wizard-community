// Sidebar: directive category tree + profiles list
import { navigate } from './app.js';

export function renderSidebar(state) {
  const nav = document.getElementById('sidebar-nav');
  nav.innerHTML = '';

  // Group directives by category
  const cats = {};
  (state.directives || []).forEach(d => {
    if (!cats[d.category]) cats[d.category] = [];
    cats[d.category].push(d);
  });

  // Directives section
  nav.appendChild(makeSection('Directives', Object.entries(cats).map(([cat, dirs]) =>
    makeSection(cat, dirs.map(d =>
      makeItem(d.name, () => navigate('directive', d.id),
        state.activeSection === 'directive' && state.selectedProfile === d.id)
    ), true)
  ), false));

  // Profiles section
  nav.appendChild(makeSection('Profiles',
    (state.profiles || []).map(p =>
      makeItem(p.name || p.id, () => navigate('profile', p.id),
        state.activeSection === 'profile' && state.selectedProfile === p.id)
    )
  ));

  // History
  nav.appendChild(makeSection('History', [
    makeItem('Generations', () => navigate('generations'),
      state.activeSection === 'generations'),
  ]));
}

function makeSection(title, children, nested = false) {
  const sec = document.createElement('div');
  sec.className = 'nav-section';

  const hdr = document.createElement('div');
  hdr.className = 'nav-section-header';
  hdr.innerHTML = `<span>${title}</span><span class="chevron">▾</span>`;

  const body = document.createElement('div');
  body.className = 'nav-section-body';
  children.forEach(c => c && body.appendChild(c));

  hdr.addEventListener('click', () => {
    hdr.classList.toggle('collapsed');
    body.classList.toggle('hidden');
  });

  sec.appendChild(hdr);
  sec.appendChild(body);
  return sec;
}

function makeItem(label, onClick, active = false) {
  const el = document.createElement('div');
  el.className = 'nav-item' + (active ? ' active' : '');
  el.textContent = label;
  el.addEventListener('click', onClick);
  return el;
}
