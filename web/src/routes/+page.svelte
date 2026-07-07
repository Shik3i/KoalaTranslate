<script lang="ts">
  import { onMount, setContext } from 'svelte';
  import { goto } from '$app/navigation';

  // Import simulated screens
  import PopupDefault from '$lib/projects/koalasync/screens/PopupDefault.svelte';
  import PopupRoomActive from '$lib/projects/koalasync/screens/PopupRoomActive.svelte';
  import PopupJoinRoom from '$lib/projects/koalasync/screens/PopupJoinRoom.svelte';
  import PopupErrorState from '$lib/projects/koalasync/screens/PopupErrorState.svelte';
  import PopupSettings from '$lib/projects/koalasync/screens/PopupSettings.svelte';
  import PopupRelaySettings from '$lib/projects/koalasync/screens/PopupRelaySettings.svelte';
  import AudioCompressorSettings from '$lib/projects/koalasync/screens/AudioCompressorSettings.svelte';

  // Import editor components
  import TranslationEditorSidebar from '$lib/components/TranslationEditorSidebar.svelte';

  // Reactively track state
  class TranslationState {
    locale = $state('de');
    activeKey = $state('');
    items = $state<any[]>([]);
    
    get dict() {
      const d: Record<string, any> = {};
      for (const item of this.items) {
        d[item.key] = item;
      }
      return d;
    }
  }

  const translationState = new TranslationState();
  setContext('translationState', translationState);

  // UI state
  let locales = $state<{ code: string; name: string; isSource: boolean }[]>([]);
  let activeScreen = $state('popup-default');
  let showImportModal = $state(false);
  let showExportModal = $state(false);

  // Import form state
  let importType = $state('manual');
  let importLocale = $state('de');
  let importLocaleName = $state('German');
  let importIsSource = $state(false);
  let importFile = $state<File | null>(null);
  let importUrl = $state('');
  let importLoading = $state(false);
  let importError = $state('');

  // Export form state
  let exportPrOwner = $state('Shik3i');
  let exportPrRepo = $state('KoalaSync');
  let exportPrBranch = $state('main');
  let exportPrPath = $state('_locales/de/messages.json');
  let exportLoading = $state(false);
  let exportSuccessUrl = $state('');
  let exportError = $state('');

  // Fetch locales list
  async function fetchLocales() {
    try {
      const res = await fetch('/api/projects/koalasync/locales');
      if (res.ok) {
        locales = await res.json();
        // If target locale not in fetched list, default to first target locale or en
        if (locales.length > 0 && !locales.some(l => l.code === translationState.locale)) {
          const targets = locales.filter(l => !l.isSource);
          translationState.locale = targets.length > 0 ? targets[0].code : locales[0].code;
        }
      }
    } catch (err) {
      console.error(err);
    }
  }

  // Fetch translations for the active locale
  async function fetchTranslations() {
    if (!translationState.locale) return;
    try {
      const res = await fetch(`/api/projects/koalasync/translations?locale=${translationState.locale}`);
      if (res.ok) {
        translationState.items = await res.json();
      }
    } catch (err) {
      console.error(err);
    }
  }

  // Handle active locale change
  $effect(() => {
    if (translationState.locale) {
      fetchTranslations();
      // Update PR path locale
      exportPrPath = `_locales/${translationState.locale}/messages.json`;
    }
  });

  onMount(async () => {
    await fetchLocales();
    await fetchTranslations();
  });

  // Handle Logout
  async function handleLogout() {
    await fetch('/api/auth/logout', { method: 'POST' });
    await goto('/login');
  }

  // Handle Manual File Import
  async function handleImportSubmit(e: Event) {
    e.preventDefault();
    importLoading = true;
    importError = '';

    const formData = new FormData();
    formData.append('locale', importLocale);
    formData.append('localeName', importLocaleName);
    formData.append('isSource', String(importIsSource));

    if (importType === 'manual') {
      if (!importFile) {
        importError = 'Please select a messages.json file.';
        importLoading = false;
        return;
      }
      formData.append('file', importFile);

      try {
        const res = await fetch('/api/projects/koalasync/import/upload', {
          method: 'POST',
          body: formData
        });
        if (res.ok) {
          showImportModal = false;
          importFile = null;
          await fetchLocales();
          await fetchTranslations();
        } else {
          const data = await res.json();
          importError = data.error || 'Import failed.';
        }
      } catch (err) {
        importError = 'Connection failed.';
      } finally {
        importLoading = false;
      }
    } else {
      // GitHub import
      if (!importUrl) {
        importError = 'Please enter a valid GitHub URL.';
        importLoading = false;
        return;
      }

      try {
        const res = await fetch('/api/projects/koalasync/import/github', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            url: importUrl,
            locale: importLocale,
            localeName: importLocaleName,
            isSource: importIsSource
          })
        });
        if (res.ok) {
          showImportModal = false;
          importUrl = '';
          await fetchLocales();
          await fetchTranslations();
        } else {
          const data = await res.json();
          importError = data.error || 'GitHub import failed.';
        }
      } catch (err) {
        importError = 'Connection failed.';
      } finally {
        importLoading = false;
      }
    }
  }

  // Handle download export
  function downloadTranslations() {
    window.location.href = `/api/projects/koalasync/export?locale=${translationState.locale}`;
  }

  // Handle GitHub PR export
  async function submitGitHubPR() {
    exportLoading = true;
    exportError = '';
    exportSuccessUrl = '';

    try {
      const res = await fetch('/api/projects/koalasync/export/github-pr', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          locale: translationState.locale,
          upstreamOwner: exportPrOwner,
          upstreamRepo: exportPrRepo,
          upstreamBranch: exportPrBranch,
          filePath: exportPrPath
        })
      });

      if (res.ok) {
        const data = await res.json();
        exportSuccessUrl = data.prUrl;
        // Refresh translations state to reset dirty indicators in UI
        await fetchTranslations();
      } else {
        const data = await res.json();
        exportError = data.error || 'Failed to submit PR.';
      }
    } catch (err) {
      exportError = 'Connection failed.';
    } finally {
      exportLoading = false;
    }
  }

  function handleFileChange(e: Event) {
    const target = e.target as HTMLInputElement;
    if (target.files && target.files.length > 0) {
      importFile = target.files[0];
    }
  }

  // Count dirty/changed translation items
  let dirtyCount = $derived(translationState.items.filter(item => item.isDirty).length);
</script>

<div class="workspace-layout">
  <!-- Left Sidebar (Portal Navigation) -->
  <aside class="sidebar-navigation glass-panel">
    <div class="portal-brand">
      <span class="logo-emoji">🐨</span>
      <h2>KoalaTranslate</h2>
    </div>

    <!-- Active Project Details -->
    <div class="project-info-card">
      <div class="project-title-row">
        <span class="project-icon">⚡</span>
        <span class="project-name">KoalaSync</span>
      </div>
      <p class="project-desc text-muted">Chrome & Firefox video playback synchronizer.</p>
    </div>

    <!-- Locale Selector -->
    <div class="nav-section">
      <label>Target Locale</label>
      <select bind:value={translationState.locale} class="locale-select">
        {#each locales as loc}
          <option value={loc.code}>
            {loc.isSource ? '🇬🇧' : '🌐'} {loc.name} ({loc.code})
          </option>
        {/each}
      </select>
    </div>

    <!-- Screen Simulator Selector -->
    <div class="nav-section flex-1">
      <label>Visual Screens</label>
      <nav class="screen-links">
        <button 
          class="screen-link" 
          class:active={activeScreen === 'popup-default'} 
          onclick={() => { activeScreen = 'popup-default' }}
        >
          <span>🚪</span> Popup Default
        </button>
        <button 
          class="screen-link" 
          class:active={activeScreen === 'popup-active'} 
          onclick={() => { activeScreen = 'popup-active' }}
        >
          <span>⚡</span> Popup Room Created
        </button>
        <button 
          class="screen-link" 
          class:active={activeScreen === 'popup-join'} 
          onclick={() => { activeScreen = 'popup-join' }}
        >
          <span>🔑</span> Popup Join Room
        </button>
        <button 
          class="screen-link" 
          class:active={activeScreen === 'popup-error'} 
          onclick={() => { activeScreen = 'popup-error' }}
        >
          <span>⚠️</span> Popup Error State
        </button>
        <button 
          class="screen-link" 
          class:active={activeScreen === 'popup-settings'} 
          onclick={() => { activeScreen = 'popup-settings' }}
        >
          <span>⚙️</span> Options General Settings
        </button>
        <button 
          class="screen-link" 
          class:active={activeScreen === 'popup-relay'} 
          onclick={() => { activeScreen = 'popup-relay' }}
        >
          <span>📡</span> Options Relay Settings
        </button>
        <button 
          class="screen-link" 
          class:active={activeScreen === 'audio-compressor'} 
          onclick={() => { activeScreen = 'audio-compressor' }}
        >
          <span>🔊</span> Options Compressor
        </button>
      </nav>
    </div>

    <!-- Actions Panel at Bottom -->
    <div class="sidebar-footer">
      {#if dirtyCount > 0}
        <div class="unsaved-badge">
          ⚠️ {dirtyCount} unsaved translation changes
        </div>
      {/if}
      
      <div class="flex gap-2">
        <button class="btn-secondary flex-1" onclick={() => { showImportModal = true }}>Import</button>
        <button class="btn-primary flex-1" onclick={() => { showExportModal = true }}>Export</button>
      </div>

      <button class="logout-btn w-full" onclick={handleLogout}>
        Logout
      </button>
    </div>
  </aside>

  <!-- Center Viewport (Screen Simulator) -->
  <main class="simulator-viewport">
    <div class="browser-header">
      <div class="browser-dots">
        <span class="dot red"></span>
        <span class="dot yellow"></span>
        <span class="dot green"></span>
      </div>
      <div class="browser-address-bar">
        <span>🔒</span>
        <span class="url-text">
          {activeScreen === 'audio-compressor' 
            ? 'chrome-extension://koalasync/audio-options.html' 
            : 'chrome-extension://koalasync/popup.html'}
        </span>
      </div>
    </div>

    <!-- Simulator Frame -->
    <div class="simulator-frame">
      {#if activeScreen === 'popup-default'}
        <PopupDefault />
      {:else if activeScreen === 'popup-active'}
        <PopupRoomActive />
      {:else if activeScreen === 'popup-join'}
        <PopupJoinRoom />
      {:else if activeScreen === 'popup-error'}
        <PopupErrorState />
      {:else if activeScreen === 'popup-settings'}
        <PopupSettings />
      {:else if activeScreen === 'popup-relay'}
        <PopupRelaySettings />
      {:else if activeScreen === 'audio-compressor'}
        <AudioCompressorSettings />
      {/if}
    </div>
    
    <div class="simulator-caption text-muted">
      💡 Hover and click text in the simulated window above to edit translation key values directly in real UI context.
    </div>
  </main>

  <!-- Right Sidebar (Translation Editor) -->
  <TranslationEditorSidebar />

  <!-- IMPORT MODAL DRAWER -->
  {#if showImportModal}
    <div class="modal-overlay" onclick={() => { showImportModal = false }}>
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div class="modal-card glass-panel" onclick={e => e.stopPropagation()}>
        <h2>Import translations</h2>
        <p class="modal-subtitle text-muted">Import messages.json dictionary key-value entries to SQLite.</p>

        {#if importError}
          <div class="error-banner border-error text-error">
            <p>{importError}</p>
          </div>
        {/if}

        <form onsubmit={handleImportSubmit}>
          <div class="form-group">
            <label>Import Method</label>
            <div class="import-tabs">
              <button 
                type="button"
                class="import-tab-btn" 
                class:active={importType === 'manual'}
                onclick={() => { importType = 'manual' }}
              >
                Manual File Upload
              </button>
              <button 
                type="button"
                class="import-tab-btn" 
                class:active={importType === 'github'}
                onclick={() => { importType = 'github' }}
              >
                GitHub raw/blob URL
              </button>
            </div>
          </div>

          <!-- Shared locale configurations -->
          <div class="flex gap-2">
            <div class="form-group flex-1">
              <label for="import-locale">Locale Code</label>
              <input id="import-locale" type="text" bind:value={importLocale} placeholder="e.g. de, fr, es" required />
            </div>
            <div class="form-group flex-1">
              <label for="import-localename">Locale Name</label>
              <input id="import-localename" type="text" bind:value={importLocaleName} placeholder="e.g. German, French" required />
            </div>
          </div>

          <div class="form-group flex items-center gap-2">
            <input type="checkbox" id="isSource" bind:checked={importIsSource} />
            <label for="isSource" class="margin-zero pointer">Register as baseline Source Locale (e.g. English)</label>
          </div>

          {#if importType === 'manual'}
            <div class="form-group">
              <label for="import-file">Upload messages.json</label>
              <input id="import-file" type="file" accept=".json" onchange={handleFileChange} required />
            </div>
          {:else}
            <div class="form-group">
              <label for="import-url">GitHub File URL</label>
              <input 
                id="import-url"
                type="text" 
                bind:value={importUrl} 
                placeholder="https://github.com/Shik3i/KoalaSync/blob/main/extension/locales/en.json" 
                required 
              />
              <p class="input-hint text-muted">Supports github.com blob and raw.githubusercontent.com URLs only.</p>
            </div>
          {/if}

          <div class="modal-footer flex gap-2">
            <button type="submit" class="btn-primary flex-1" disabled={importLoading}>
              {#if importLoading}
                Importing...
              {:else}
                Execute Import
              {/if}
            </button>
            <button type="button" class="btn-secondary" onclick={() => { showImportModal = false }} disabled={importLoading}>
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  {/if}

  <!-- EXPORT MODAL DRAWER -->
  {#if showExportModal}
    <div class="modal-overlay" onclick={() => { showExportModal = false }}>
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div class="modal-card glass-panel" onclick={e => e.stopPropagation()}>
        <h2>Export Translations ({translationState.locale.toUpperCase()})</h2>
        <p class="modal-subtitle text-muted">Retrieve completed local SQLite translation values back to file.</p>

        {#if exportError}
          <div class="error-banner border-error text-error">
            <p>{exportError}</p>
          </div>
        {/if}

        {#if exportSuccessUrl}
          <div class="success-banner">
            <span>🎉</span>
            <div>
              <p class="bold-text">Pull Request Created Successfully!</p>
              <a href={exportSuccessUrl} target="_blank" class="pr-link">View Pull Request #{exportSuccessUrl.split('/').pop()}</a>
            </div>
          </div>
        {/if}

        <div class="export-options flex flex-col gap-4">
          <!-- Manual Local Download -->
          <div class="export-box flex justify-between items-center bg-card border-card">
            <div>
              <h4 class="margin-zero">Download messages.json</h4>
              <p class="text-muted font-size-11 margin-zero">Directly download target locale i18n file to your local computer.</p>
            </div>
            <button class="btn-secondary" onclick={downloadTranslations}>Download</button>
          </div>

          <!-- Create GitHub Pull Request -->
          <div class="export-box bg-card border-card flex flex-col gap-4">
            <div>
              <h4 class="margin-zero">Create GitHub Pull Request</h4>
              <p class="text-muted font-size-11 margin-zero">Open a Pull Request with updated translations from bot-fork repository.</p>
            </div>

            <div class="pr-fields-grid">
              <div class="form-group">
                <label for="exp-owner">Upstream Owner</label>
                <input id="exp-owner" type="text" bind:value={exportPrOwner} placeholder="Shik3i" />
              </div>
              <div class="form-group">
                <label for="exp-repo">Upstream Repo</label>
                <input id="exp-repo" type="text" bind:value={exportPrRepo} placeholder="KoalaSync" />
              </div>
              <div class="form-group">
                <label for="exp-branch">Upstream Branch</label>
                <input id="exp-branch" type="text" bind:value={exportPrBranch} placeholder="main" />
              </div>
              <div class="form-group">
                <label for="exp-filepath">PR File Target Path</label>
                <input id="exp-filepath" type="text" bind:value={exportPrPath} />
              </div>
            </div>

            <button 
              class="btn-primary w-full" 
              disabled={exportLoading}
              onclick={submitGitHubPR}
            >
              {#if exportLoading}
                Submitting PR...
              {:else}
                Submit Pull Request
              {/if}
            </button>
          </div>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn-secondary w-full" onclick={() => { showExportModal = false }}>
            Close Export Drawer
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .workspace-layout {
    display: flex;
    height: 100vh;
    width: 100vw;
    background: #090d16;
    overflow: hidden;
  }
  
  /* Navigation Sidebar */
  .sidebar-navigation {
    width: 280px;
    height: 100%;
    padding: 24px;
    display: flex;
    flex-direction: column;
    border-radius: 0 !important;
    border-right: 1px solid var(--border-color);
    border-left: none;
    border-top: none;
    border-bottom: none;
  }
  .portal-brand {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 24px;
  }
  .portal-brand h2 {
    font-size: 18px;
    font-weight: 800;
    margin: 0;
    letter-spacing: -0.5px;
    background: linear-gradient(135deg, #a855f7, #6366f1);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
  .project-info-card {
    background: rgba(15, 23, 42, 0.4);
    border: 1px solid var(--border-color);
    padding: 14px;
    border-radius: var(--radius-md);
    margin-bottom: 20px;
  }
  .project-title-row {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 700;
    font-size: 14px;
    color: var(--color-text);
  }
  .project-desc {
    font-size: 11px;
    margin: 6px 0 0 0;
    line-height: 1.4;
  }
  .nav-section {
    margin-bottom: 20px;
  }
  .locale-select {
    width: 100%;
  }
  .screen-links {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
  .screen-link {
    width: 100%;
    text-align: left;
    padding: 10px 14px;
    background: transparent;
    color: var(--color-text-muted);
    font-size: 13px;
    border-radius: var(--radius-sm);
    justify-content: flex-start;
    font-weight: 500;
  }
  .screen-link:hover {
    background: rgba(255, 255, 255, 0.05);
    color: var(--color-text);
  }
  .screen-link.active {
    background: var(--bg-active);
    color: var(--color-accent-hover);
    font-weight: 700;
  }
  .sidebar-footer {
    margin-top: auto;
    display: flex;
    flex-direction: column;
    gap: 12px;
    border-top: 1px solid rgba(255, 255, 255, 0.05);
    padding-top: 20px;
  }
  .unsaved-badge {
    background: rgba(245, 158, 11, 0.1);
    border: 1px solid rgba(245, 158, 11, 0.25);
    color: #fcd34d;
    padding: 6px 12px;
    border-radius: var(--radius-sm);
    font-size: 11px;
    font-weight: bold;
    text-align: center;
  }
  .logout-btn {
    background: transparent;
    color: var(--color-text-muted);
    border: 1px solid transparent;
    font-size: 12px;
    padding: 8px;
  }
  .logout-btn:hover {
    color: var(--color-error);
    background: rgba(239, 68, 68, 0.05);
    border-color: rgba(239, 68, 68, 0.2);
  }

  /* Center Viewport */
  .simulator-viewport {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: radial-gradient(circle at 50% 50%, rgba(99, 102, 241, 0.04) 0%, transparent 60%),
                #090d16;
    padding: 40px;
    overflow-y: auto;
  }
  .browser-header {
    width: 100%;
    max-width: 580px;
    background: #1e293b;
    border: 1px solid var(--border-color);
    border-bottom: none;
    border-top-left-radius: 12px;
    border-top-right-radius: 12px;
    padding: 10px 16px;
    display: flex;
    align-items: center;
    gap: 16px;
  }
  .browser-dots {
    display: flex;
    gap: 6px;
  }
  .dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
  }
  .dot.red { background: #ef4444; }
  .dot.yellow { background: #f59e0b; }
  .dot.green { background: #10b981; }
  .browser-address-bar {
    flex: 1;
    background: #0f172a;
    border-radius: 6px;
    padding: 6px 12px;
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--color-text-muted);
    display: flex;
    align-items: center;
    gap: 8px;
    border: 1px solid rgba(255, 255, 255, 0.05);
  }
  .url-text {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .simulator-frame {
    width: 100%;
    max-width: 580px;
    background: #090d16;
    border: 1px solid var(--border-color);
    border-bottom-left-radius: 12px;
    border-bottom-right-radius: 12px;
    padding: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 480px;
  }
  .simulator-caption {
    margin-top: 16px;
    font-size: 12px;
    max-width: 480px;
    text-align: center;
    line-height: 1.4;
  }

  /* Modals */
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(4px);
    z-index: 100;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
  }
  .modal-card {
    width: 100%;
    max-width: 500px;
    padding: 32px;
    animation: scaleUp 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }
  @keyframes scaleUp {
    from { transform: scale(0.95); opacity: 0; }
    to { transform: scale(1); opacity: 1; }
  }
  .modal-card h2 {
    margin: 0 0 6px 0;
    font-size: 20px;
    font-weight: 800;
  }
  .modal-subtitle {
    font-size: 13px;
    margin: 0 0 24px 0;
  }
  .import-tabs {
    display: flex;
    background: rgba(0, 0, 0, 0.2);
    padding: 4px;
    border-radius: var(--radius-md);
    border: 1px solid var(--border-color);
    margin-bottom: 4px;
  }
  .import-tab-btn {
    flex: 1;
    padding: 8px;
    background: transparent;
    color: var(--color-text-muted);
    font-size: 12px;
    border-radius: var(--radius-sm);
  }
  .import-tab-btn.active {
    background: var(--color-accent);
    color: white;
  }
  .input-hint {
    font-size: 10px;
    margin: 4px 0 0 0;
  }
  .modal-footer {
    margin-top: 24px;
  }
  .export-box {
    padding: 16px;
    border-radius: var(--radius-md);
  }
  .bg-card {
    background: rgba(15, 23, 42, 0.4);
  }
  .border-card {
    border: 1px solid var(--border-color);
  }
  .pr-fields-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
  }
  .pr-fields-grid .form-group {
    margin-bottom: 0;
  }
  .error-banner {
    background: rgba(239, 68, 68, 0.1);
    border-left: 3px solid var(--color-error);
    padding: 10px 14px;
    border-radius: var(--radius-sm);
    margin-bottom: 16px;
    font-size: 13px;
  }
  .error-banner p { margin: 0; }
  .success-banner {
    display: flex;
    gap: 12px;
    background: rgba(16, 185, 129, 0.1);
    border: 1px solid rgba(16, 185, 129, 0.3);
    padding: 16px;
    border-radius: var(--radius-md);
    margin-bottom: 20px;
    font-size: 13px;
  }
  .success-banner p { margin: 0; }
  .pr-link {
    color: var(--color-success);
    font-weight: bold;
    text-decoration: underline;
    display: inline-block;
    margin-top: 4px;
  }
  .font-size-11 { font-size: 11px; }
  .margin-zero { margin: 0; }
</style>
