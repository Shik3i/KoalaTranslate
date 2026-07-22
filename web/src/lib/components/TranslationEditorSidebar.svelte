<script lang="ts">
  import { getContext } from 'svelte';

  // Retrieve shared state context
  const translationState = getContext<any>('translationState');
  const projectID = 'koalasync';

  let currentText = $state('');
  let item = $derived(translationState?.dict[translationState?.activeKey]);
  let originalText = $derived(item?.targetValue || '');
  let isDirty = $derived(currentText !== originalText);

  // Sync state when active key changes
  $effect(() => {
    if (item) {
      currentText = item.targetValue || '';
    } else {
      currentText = '';
    }
  });

  // Warnings checking
  let placeholderWarnings = $derived(
    item ? checkPlaceholders(item.sourceValue, currentText) : []
  );

  let lengthWarning = $derived(
    item && item.maxChars && currentText.length > item.maxChars
      ? `Length exceeds maximum character limit of ${item.maxChars}.`
      : ''
  );

  function checkPlaceholders(source: string, target: string): string[] {
    const patterns = [
      /\{\{[a-zA-Z0-9_-]+\}\}/g,
      /\{[a-zA-Z0-9_-]+\}/g,
      /%[s|d]/g,
      /\$[1-9]/g,
      /\$[a-zA-Z0-9_-]+\$/g
    ];

    const extracted: string[] = [];
    for (const re of patterns) {
      const matches = source.match(re);
      if (matches) {
        for (const m of matches) {
          if (!extracted.includes(m)) {
            extracted.push(m);
          }
        }
      }
    }

    const warnings: string[] = [];
    const targetLower = target.toLowerCase();
    for (const p of extracted) {
      if (!targetLower.includes(p.toLowerCase())) {
        warnings.push(`Missing placeholder: ${p}`);
      }
    }
    return warnings;
  }

  async function saveTranslation() {
    if (!item) return;
    try {
      const res = await fetch(`/api/projects/${projectID}/translations/${item.key}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          locale: translationState.locale,
          value: currentText
        })
      });

      if (res.ok) {
        // Update local state reactively
        item.targetValue = currentText;
        item.isDirty = true; // Mark as locally changed in UI
        
        // Re-evaluate warnings
        item.warnings = [...placeholderWarnings];
        if (lengthWarning) {
          item.warnings.push(lengthWarning);
        }
      } else {
        alert('Failed to save translation.');
      }
    } catch (err) {
      console.error(err);
      alert('Network error while saving.');
    }
  }

  function discardChanges() {
    currentText = originalText;
  }
</script>

<div class="sidebar-editor glass-panel">
  {#if !translationState?.activeKey}
    <div class="empty-state text-muted">
      <span>👁️</span>
      <p>Select any text block in the simulated preview screen to edit its translation.</p>
    </div>
  {:else if item}
    <div class="editor-header">
      <div class="flex justify-between items-center">
        <h3>Translation Editor</h3>
        <button class="close-btn" onclick={() => { translationState.activeKey = '' }}>×</button>
      </div>
      <div class="key-badge">{item.key}</div>
      {#if item.description}
        <p class="description text-muted">{item.description}</p>
      {/if}
    </div>

    <div class="editor-body">
      <!-- Reference Source -->
      <div class="field-group">
        <span class="field-label">Source Locale (English)</span>
        <div class="source-bubble">{item.sourceValue || '(Empty)'}</div>
      </div>

      <!-- Target Input -->
      <div class="field-group">
        <label for="translation-value">Target Translation ({translationState.locale.toUpperCase()})</label>
        <textarea 
          id="translation-value"
          bind:value={currentText}
          rows="4"
          placeholder="Type translation here..."
        ></textarea>
      </div>

      <!-- Character limit warning -->
      {#if lengthWarning}
        <div class="warning-banner border-error text-error">
          <p>{lengthWarning}</p>
        </div>
      {/if}

      <!-- Placeholder warning -->
      {#if placeholderWarnings.length > 0}
        <div class="warning-banner border-warning text-warning">
          {#each placeholderWarnings as warning}
            <p>⚠️ {warning}</p>
          {/each}
        </div>
      {/if}

      <!-- Action Footer -->
      <div class="editor-actions flex gap-2">
        <button 
          class="btn-primary flex-1" 
          disabled={!isDirty}
          onclick={saveTranslation}
        >
          Save Changes
        </button>
        <button 
          class="btn-secondary" 
          disabled={!isDirty}
          onclick={discardChanges}
        >
          Reset
        </button>
      </div>
    </div>
  {/if}
</div>

<style>
  .sidebar-editor {
    width: 350px;
    height: 100%;
    padding: 24px;
    display: flex;
    flex-direction: column;
    border-radius: 0 !important;
    border-left: 1px solid var(--border-color);
    border-right: none;
    border-top: none;
    border-bottom: none;
  }
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    text-align: center;
    font-size: 13px;
  }
  .empty-state span {
    font-size: 32px;
    margin-bottom: 12px;
  }
  .editor-header h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 700;
  }
  .close-btn {
    background: transparent;
    border: none;
    color: var(--color-text-muted);
    font-size: 20px;
    padding: 0;
    cursor: pointer;
  }
  .key-badge {
    display: inline-block;
    font-family: var(--font-mono);
    font-size: 11px;
    background: var(--bg-active);
    color: var(--color-accent-hover);
    padding: 3px 8px;
    border-radius: 4px;
    margin-top: 10px;
    font-weight: 600;
  }
  .description {
    font-size: 12px;
    margin: 8px 0 0 0;
    line-height: 1.4;
  }
  .editor-body {
    display: flex;
    flex-direction: column;
    gap: 20px;
    margin-top: 24px;
    overflow-y: auto;
    flex: 1;
  }
  .field-group .field-label,
  .field-group label {
    margin-bottom: 6px;
  }
  .source-bubble {
    background: rgba(15, 23, 42, 0.4);
    border: 1px solid var(--border-color);
    padding: 12px;
    border-radius: var(--radius-md);
    font-size: 13px;
    line-height: 1.4;
    white-space: pre-wrap;
  }
  .warning-banner {
    background: rgba(15, 23, 42, 0.4);
    border-left: 3px solid;
    padding: 8px 12px;
    border-radius: var(--radius-sm);
    font-size: 12px;
  }
  .warning-banner p {
    margin: 0;
    line-height: 1.4;
  }
  .border-error { border-left-color: var(--color-error); }
  .border-warning { border-left-color: var(--color-warning); }
  .editor-actions {
    margin-top: auto;
    padding-top: 20px;
  }
</style>
