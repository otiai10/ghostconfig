// i18n state
let i18nData = {
    currentLang: 'en',
    languages: ['en', 'ja'],
    messages: {}
};

// Language display names
const langNames = {
    en: 'EN',
    ja: 'JA'
};

// Translate function
function t(key) {
    const msgs = i18nData.messages[i18nData.currentLang] || {};
    return msgs[key] || key;
}

// Application state
const state = {
    sections: [],
    currentSection: 'all',
    searchQuery: '',
    colors: [],
    fonts: [],
    currentOption: null
};

// Initialize the application
async function init() {
    try {
        await Promise.all([
            loadI18n(),
            loadOptions(),
            loadColors()
        ]);
        renderLangSwitcher();
        applyI18n();
        renderSections();
        renderOptions();
        setupEventListeners();
    } catch (error) {
        console.error(t('gui.error.init'), error);
        document.getElementById('options').innerHTML =
            `<div class="no-results">${t('gui.error.load_options')}</div>`;
    }
}

// Load i18n messages
async function loadI18n() {
    const response = await fetch('/api/i18n');
    if (!response.ok) return;
    const data = await response.json();
    i18nData.languages = data.languages || ['en', 'ja'];
    i18nData.messages = data.messages || {};

    // Check localStorage for saved language preference
    const savedLang = localStorage.getItem('ghostconfig-lang');
    if (savedLang && i18nData.languages.includes(savedLang)) {
        i18nData.currentLang = savedLang;
    } else {
        i18nData.currentLang = data.defaultLang || 'en';
    }
}

// Render language switcher select
function renderLangSwitcher() {
    const container = document.getElementById('lang-switcher');
    const options = i18nData.languages.map(lang => {
        const selected = lang === i18nData.currentLang ? 'selected' : '';
        return `<option value="${lang}" ${selected}>${langNames[lang] || lang.toUpperCase()}</option>`;
    });
    container.innerHTML = `<select id="lang-select">${options.join('')}</select>`;
}

// Switch language
function switchLanguage(lang) {
    if (!i18nData.languages.includes(lang)) return;

    i18nData.currentLang = lang;
    localStorage.setItem('ghostconfig-lang', lang);

    // Update UI
    renderLangSwitcher();
    applyI18n();
    renderSections();
    renderOptions();
}

// Apply i18n to static elements
function applyI18n() {
    document.title = t('app.title');
    const h1 = document.querySelector('h1');
    if (h1) h1.textContent = t('app.title');
    const search = document.getElementById('search');
    if (search) search.placeholder = t('gui.search_placeholder');
    const exitBtn = document.getElementById('exit-btn');
    if (exitBtn) exitBtn.textContent = t('gui.exit');
    const modalCancel = document.getElementById('modal-cancel');
    if (modalCancel) modalCancel.textContent = t('gui.cancel');
    const modalSave = document.getElementById('modal-save');
    if (modalSave) modalSave.textContent = t('gui.save');
    const loading = document.querySelector('.loading');
    if (loading) loading.textContent = t('gui.loading');
}

// API calls
async function loadOptions() {
    const response = await fetch('/api/options');
    if (!response.ok) throw new Error(t('gui.error.load_options_api'));
    state.sections = await response.json();
}

async function loadColors() {
    const response = await fetch('/api/colors');
    if (!response.ok) throw new Error(t('gui.error.load_colors'));
    state.colors = await response.json();
}

async function loadFonts() {
    if (state.fonts.length > 0) return state.fonts;
    const response = await fetch('/api/fonts');
    if (!response.ok) throw new Error(t('gui.error.load_fonts'));
    state.fonts = await response.json();
    return state.fonts;
}

async function saveConfig(key, value) {
    const response = await fetch('/api/config', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ key, value })
    });

    if (!response.ok) {
        throw new Error(t('gui.error.save'));
    }

    // Update local state
    for (const section of state.sections) {
        for (const opt of section.options) {
            if (opt.key === key) {
                opt.currentValue = value;
                break;
            }
        }
    }

    showStatus(t('msg.saved').replace('%s', key).replace('%s', value));
    renderOptions();
}

// Rendering functions
function renderSections() {
    const nav = document.getElementById('sections');
    const buttons = state.sections.map(section =>
        `<button class="section-btn" data-section="${section.name}">${section.name}</button>`
    );
    nav.innerHTML = `<button class="section-btn active" data-section="all">${t('gui.all')}</button>` + buttons.join('');
}

function renderOptions() {
    const container = document.getElementById('options');
    const filtered = filterOptions();

    if (filtered.length === 0) {
        container.innerHTML = `<div class="no-results">${t('gui.no_options')}</div>`;
        return;
    }

    let html = '';

    for (const opt of filtered) {
        const displayValue = opt.currentValue || opt.defaultValue || '(empty)';
        const isModified = opt.currentValue && opt.currentValue !== opt.defaultValue;
        const badge = isModified
            ? `<span class="modified-badge">${t('gui.modified')}</span>`
            : (!opt.currentValue ? `<span class="default-badge">${t('gui.default')}</span>` : '');

        let valueHtml = escapeHtml(displayValue);
        if (opt.type === 'color' && displayValue && displayValue !== '(empty)') {
            const colorVal = displayValue.startsWith('#') ? displayValue : `#${displayValue}`;
            valueHtml = `<span class="color-preview" style="background: ${colorVal}"></span> ${escapeHtml(displayValue)}`;
        }

        html += `
            <div class="option-card" data-key="${opt.key}">
                <div class="option-header">
                    <span class="option-key">${escapeHtml(opt.key)}</span>
                    ${badge}
                </div>
                <div class="option-value">${valueHtml}</div>
                <div class="option-description">${escapeHtml(opt.description || '')}</div>
            </div>
        `;
    }

    container.innerHTML = html;
}

function filterOptions() {
    let options = [];

    for (const section of state.sections) {
        for (const opt of section.options) {
            const matchesSection = state.currentSection === 'all' || section.name === state.currentSection;
            const matchesSearch = !state.searchQuery ||
                opt.key.toLowerCase().includes(state.searchQuery.toLowerCase()) ||
                (opt.description && opt.description.toLowerCase().includes(state.searchQuery.toLowerCase()));

            if (matchesSection && matchesSearch) {
                options.push({ ...opt, section: section.name });
            }
        }
    }

    return options;
}

function findOption(key) {
    for (const section of state.sections) {
        for (const opt of section.options) {
            if (opt.key === key) {
                return opt;
            }
        }
    }
    return null;
}

// Modal and editor functions
function openEditor(option) {
    state.currentOption = option;
    const modal = document.getElementById('modal');
    const title = document.getElementById('modal-title');
    const description = document.getElementById('modal-description');
    const container = document.getElementById('modal-input-container');

    title.textContent = option.key;
    description.textContent = option.description || '';

    const currentValue = option.currentValue || option.defaultValue || '';

    switch (option.type) {
        case 'color':
            renderColorPicker(container, currentValue);
            break;
        case 'font':
            renderFontPicker(container, currentValue);
            break;
        default:
            container.innerHTML = `<input type="text" id="edit-value" value="${escapeHtml(currentValue)}" placeholder="${escapeHtml(option.defaultValue || '')}">`;
    }

    modal.classList.remove('hidden');

    const input = container.querySelector('input');
    if (input) {
        input.focus();
        input.select();
    }
}

function closeModal() {
    document.getElementById('modal').classList.add('hidden');
    state.currentOption = null;
}

function renderColorPicker(container, currentValue) {
    const normalizedCurrent = currentValue.replace('#', '').toLowerCase();

    let html = '<div class="color-grid">';
    for (const color of state.colors) {
        const isSelected = color.value.toLowerCase() === normalizedCurrent;
        html += `
            <button class="color-option ${isSelected ? 'selected' : ''}"
                    data-value="${color.value}"
                    style="background: #${color.value}"
                    title="${color.name}">
            </button>
        `;
    }
    html += '</div>';

    html += `
        <div class="custom-color">
            <label>${t('gui.custom_label')}</label>
            <input type="color" id="custom-color-picker" value="#${normalizedCurrent || 'ffffff'}">
            <input type="text" id="custom-color-hex" value="${normalizedCurrent}" placeholder="ffffff">
        </div>
    `;

    container.innerHTML = html;

    // Event listeners for color picker
    container.querySelectorAll('.color-option').forEach(btn => {
        btn.addEventListener('click', () => {
            container.querySelectorAll('.color-option').forEach(b => b.classList.remove('selected'));
            btn.classList.add('selected');
            document.getElementById('custom-color-hex').value = btn.dataset.value;
            document.getElementById('custom-color-picker').value = '#' + btn.dataset.value;
        });
    });

    document.getElementById('custom-color-picker').addEventListener('input', (e) => {
        document.getElementById('custom-color-hex').value = e.target.value.replace('#', '');
        container.querySelectorAll('.color-option').forEach(b => b.classList.remove('selected'));
    });

    document.getElementById('custom-color-hex').addEventListener('input', (e) => {
        const hex = e.target.value.replace('#', '');
        if (/^[0-9a-fA-F]{6}$/.test(hex)) {
            document.getElementById('custom-color-picker').value = '#' + hex;
        }
        container.querySelectorAll('.color-option').forEach(b => b.classList.remove('selected'));
    });
}

async function renderFontPicker(container, currentValue) {
    container.innerHTML = `<div class="loading">${t('gui.loading_fonts')}</div>`;

    try {
        const fonts = await loadFonts();

        let html = `<input type="text" class="font-filter" id="font-filter" placeholder="${t('gui.filter_fonts')}" value="">`;
        html += '<div class="font-list" id="font-list">';

        for (const font of fonts) {
            const isSelected = font === currentValue;
            html += `
                <div class="font-option ${isSelected ? 'selected' : ''}" data-font="${escapeHtml(font)}">
                    <span class="font-name">${escapeHtml(font)}</span>
                    <span class="font-sample" style="font-family: '${escapeHtml(font)}', monospace;">abcABC0123</span>
                </div>`;
        }

        html += '</div>';
        container.innerHTML = html;

        // Font filter
        document.getElementById('font-filter').addEventListener('input', (e) => {
            const filter = e.target.value.toLowerCase();
            document.querySelectorAll('.font-option').forEach(opt => {
                const match = opt.dataset.font.toLowerCase().includes(filter);
                opt.style.display = match ? '' : 'none';
            });
        });

        // Font selection
        document.getElementById('font-list').addEventListener('click', (e) => {
            if (e.target.classList.contains('font-option')) {
                document.querySelectorAll('.font-option').forEach(o => o.classList.remove('selected'));
                e.target.classList.add('selected');
            }
        });
    } catch (error) {
        container.innerHTML = `<div class="no-results">${t('gui.error.load_fonts')}</div>`;
    }
}

async function saveCurrentEdit() {
    if (!state.currentOption) return;

    let value = '';
    const option = state.currentOption;
    const container = document.getElementById('modal-input-container');

    switch (option.type) {
        case 'color':
            value = document.getElementById('custom-color-hex').value;
            break;
        case 'font':
            const selected = container.querySelector('.font-option.selected');
            value = selected ? selected.dataset.font : '';
            break;
        default:
            value = document.getElementById('edit-value').value;
    }

    try {
        await saveConfig(option.key, value);
        closeModal();
    } catch (error) {
        showStatus(t('gui.error.save_prefix') + error.message, true);
    }
}

// Event listeners
function setupEventListeners() {
    // Search
    const search = document.getElementById('search');
    if (search) {
        search.addEventListener('input', (e) => {
            state.searchQuery = e.target.value;
            renderOptions();
        });
    }

    // Section navigation
    document.getElementById('sections').addEventListener('click', (e) => {
        if (e.target.classList.contains('section-btn')) {
            document.querySelectorAll('.section-btn').forEach(b => b.classList.remove('active'));
            e.target.classList.add('active');
            state.currentSection = e.target.dataset.section;
            renderOptions();
        }
    });

    // Option click
    document.getElementById('options').addEventListener('click', (e) => {
        const card = e.target.closest('.option-card');
        if (card) {
            const key = card.dataset.key;
            const option = findOption(key);
            if (option) {
                openEditor(option);
            }
        }
    });

    // Modal events
    document.getElementById('modal-cancel').addEventListener('click', closeModal);
    document.getElementById('modal-save').addEventListener('click', saveCurrentEdit);
    document.querySelector('.modal-backdrop').addEventListener('click', closeModal);

    // Language switcher
    document.getElementById('lang-switcher').addEventListener('change', (e) => {
        if (e.target.id === 'lang-select') {
            switchLanguage(e.target.value);
        }
    });

    // Exit button
    document.getElementById('exit-btn').addEventListener('click', async () => {
        if (confirm(t('gui.exit_confirm'))) {
            try {
                await fetch('/api/exit', { method: 'POST' });
                window.close();
                // Fallback if window.close() is blocked by browser
                document.body.innerHTML = `<div style="display:flex;align-items:center;justify-content:center;height:100vh;color:#a6adc8;">${t('gui.server_stopped')}</div>`;
            } catch (error) {
                showStatus(t('gui.error.exit'), true);
            }
        }
    });

    // Keyboard shortcuts
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            closeModal();
        }
        if (e.key === 'Enter' && !document.getElementById('modal').classList.contains('hidden')) {
            saveCurrentEdit();
        }
        if (e.key === '/' && document.activeElement.tagName !== 'INPUT') {
            const search = document.getElementById('search');
            if (search) {
                e.preventDefault();
                search.focus();
            }
        }
    });
}

// Utility functions
function escapeHtml(str) {
    if (!str) return '';
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

function showStatus(message, isError = false) {
    const status = document.getElementById('status');
    status.textContent = message;
    status.style.color = isError ? 'var(--warning)' : 'var(--success)';

    setTimeout(() => {
        status.textContent = '';
    }, 3000);
}

// Start the application
document.addEventListener('DOMContentLoaded', init);
