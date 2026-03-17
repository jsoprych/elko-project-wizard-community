// elko-project-wizard — main app orchestration
import { renderSidebar } from './sidebar.js';
import { renderWorkspace } from './workspace.js';

let state = {
  directives: [],
  profiles: [],
  selectedProfile: null,
  activeSection: 'directives',
};

async function loadData() {
  const [dirs, profs] = await Promise.all([
    fetch('/api/directives').then(r => r.json()),
    fetch('/api/profiles').then(r => r.json()),
  ]);
  state.directives = dirs || [];
  state.profiles   = profs || [];
}

export function getState() { return state; }

export async function navigate(section, id) {
  state.activeSection = section;
  state.selectedProfile = id || null;
  renderSidebar(state);
  renderWorkspace(state);
}

export async function refresh() {
  await loadData();
  renderSidebar(state);
  renderWorkspace(state);
}

// Theme
const savedTheme = localStorage.getItem('elko-theme') || 'dark';
document.documentElement.setAttribute('data-theme', savedTheme);

document.body.insertAdjacentHTML('beforeend',
  `<button class="theme-toggle" id="theme-toggle">☀ / ☾</button>`);
document.getElementById('theme-toggle').addEventListener('click', () => {
  const t = document.documentElement.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
  document.documentElement.setAttribute('data-theme', t);
  localStorage.setItem('elko-theme', t);
});

document.getElementById('btn-new-profile').addEventListener('click', () => {
  navigate('new-profile');
});

// Boot
(async () => {
  await loadData();
  renderSidebar(state);
  renderWorkspace(state);
})();
