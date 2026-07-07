<script lang="ts">
  import { getContext } from 'svelte';

  interface Props {
    keyName: string;
    maxChars?: number;
    children?: any;
  }

  let { keyName, maxChars, children }: Props = $props();

  // Retrieve translationState context
  const translationState = getContext<any>('translationState');

  // Reactively compute translations info
  let item = $derived(translationState?.dict[keyName]);
  let text = $derived(item?.targetValue || item?.sourceValue || keyName);
  let isDirty = $derived(item?.isDirty || false);
  let hasWarnings = $derived(item?.warnings && item?.warnings.length > 0);
  let isSelected = $derived(translationState?.activeKey === keyName);

  function selectKey(e: MouseEvent) {
    e.stopPropagation();
    if (translationState) {
      translationState.activeKey = keyName;
    }
  }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<span 
  class="translation-slot" 
  class:selected={isSelected}
  class:dirty={isDirty}
  class:warning={hasWarnings}
  onclick={selectKey}
  title={`Key: ${keyName}${item?.description ? '\nDesc: ' + item.description : ''}`}
>
  {text}
  {#if isDirty}
    <span class="dirty-indicator" title="Unsaved changes local to SQLite"></span>
  {/if}
</span>

<style>
  .translation-slot {
    position: relative;
    cursor: pointer;
    border: 1px dashed transparent;
    border-radius: 4px;
    padding: 1px 4px;
    transition: all 0.2s ease;
    display: inline-flex;
    align-items: center;
  }

  .translation-slot:hover {
    border-color: var(--color-accent-hover);
    background: rgba(99, 102, 241, 0.1);
  }

  .translation-slot.selected {
    border-color: var(--color-accent);
    background: rgba(99, 102, 241, 0.2);
    box-shadow: 0 0 8px rgba(99, 102, 241, 0.25);
  }

  .translation-slot.dirty {
    border-bottom: 2px solid var(--color-warning);
  }

  .translation-slot.warning {
    border-bottom: 2px solid var(--color-error);
  }

  .dirty-indicator {
    position: absolute;
    top: -2px;
    right: -2px;
    width: 6px;
    height: 6px;
    background-color: var(--color-warning);
    border-radius: 50%;
    box-shadow: 0 0 4px var(--color-warning);
  }
</style>
