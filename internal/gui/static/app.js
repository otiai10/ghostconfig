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
            loadOptions(),
            loadColors()
        ]);
        renderSections();
        renderOptions();
        setupEventListeners();
    } catch (error) {
        console.error('Failed to initialize:', error);
        document.getElementById('options').innerHTML =
            '<div class="no-results">Failed to load options. Please refresh the page.</div>';
    }
}

// API calls
async function loadOptions() {
    const response = await fetch('/api/options');
    if (!response.ok) throw new Error('Failed to load options');
    state.sections = await response.json();
}

async function loadColors() {
    const response = await fetch('/api/colors');
    if (!response.ok) throw new Error('Failed to load colors');
    state.colors = await response.json();
}

async function loadFonts() {
    if (state.fonts.length > 0) return state.fonts;
    const response = await fetch('/api/fonts');
    if (!response.ok) throw new Error('Failed to load fonts');
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
        throw new Error('Failed to save config');
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

    showStatus(`Saved: ${key} = ${value}`);
    renderOptions();
}

// Rendering functions
function renderSections() {
    const nav = document.getElementById('sections');
    const buttons = state.sections.map(section =>
        `<button class="section-btn" data-section="${section.name}">${section.name}</button>`
    );
    nav.innerHTML = `<button class="section-btn active" data-section="all">All</button>` + buttons.join('');
}

function renderOptions() {
    const container = document.getElementById('options');
    const filtered = filterOptions();

    if (filtered.length === 0) {
        container.innerHTML = '<div class="no-results">No options found</div>';
        return;
    }

    let html = '';
    let currentSection = '';

    for (const opt of filtered) {
        if (state.currentSection === 'all' && opt.section !== currentSection) {
            currentSection = opt.section;
            html += `<div class="section-header">${currentSection}</div>`;
        }

        const displayValue = opt.currentValue || opt.defaultValue || '(empty)';
        const isModified = opt.currentValue && opt.currentValue !== opt.defaultValue;
        const badge = isModified
            ? '<span class="modified-badge">modified</span>'
            : (!opt.currentValue ? '<span class="default-badge">default</span>' : '');

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
            <label>Custom:</label>
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
    container.innerHTML = '<div class="loading">Loading fonts...</div>';

    try {
        const fonts = await loadFonts();

        let html = `<input type="text" class="font-filter" id="font-filter" placeholder="Filter fonts..." value="">`;
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
        container.innerHTML = '<div class="no-results">Failed to load fonts</div>';
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
        showStatus('Failed to save: ' + error.message, true);
    }
}

// Event listeners
function setupEventListeners() {
    // Search
    document.getElementById('search').addEventListener('input', (e) => {
        state.searchQuery = e.target.value;
        renderOptions();
    });

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

    // Exit button
    document.getElementById('exit-btn').addEventListener('click', async () => {
        if (confirm('Exit Ghostty Config Editor?')) {
            try {
                await fetch('/api/exit', { method: 'POST' });
                window.close();
                // Fallback if window.close() is blocked by browser
                document.body.innerHTML = '<div style="display:flex;align-items:center;justify-content:center;height:100vh;color:#a6adc8;">Server stopped. You can close this tab.</div>';
            } catch (error) {
                showStatus('Failed to exit', true);
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
            e.preventDefault();
            document.getElementById('search').focus();
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
