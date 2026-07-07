<script lang="ts">
  import TranslationSlot from '../../../components/TranslationSlot.svelte';
</script>

<div class="audio-options-wrapper">
  <!-- Header Shell -->
  <header class="page-header">
    <div>
      <span class="back-link"><TranslationSlot keyName="AUDIO_BACK" /></span>
      <p class="eyebrow">KoalaSync</p>
      <h1 class="page-title"><TranslationSlot keyName="AUDIO_PAGE_TITLE" /></h1>
    </div>
    <div class="master-toggle">
      <span class="toggle-label"><TranslationSlot keyName="AUDIO_MASTER_TOGGLE" /></span>
      <label class="toggle-switch">
        <input type="checkbox" checked disabled />
        <span class="slider"></span>
      </label>
    </div>
  </header>

  <!-- Compressor Settings Section -->
  <section class="panel">
    <div class="section-heading flex justify-between items-center">
      <h2 class="section-title"><TranslationSlot keyName="AUDIO_COMPRESSOR" /></h2>
      <div class="inline-toggle flex items-center">
        <span class="toggle-label-sm"><TranslationSlot keyName="AUDIO_COMPRESSOR_ENABLE" /></span>
        <label class="toggle-switch-sm">
          <input type="checkbox" checked disabled />
          <span class="slider"></span>
        </label>
      </div>
    </div>

    <!-- Presets -->
    <div class="preset-group">
      <div class="preset-card active">
        <input type="radio" checked disabled />
        <span><TranslationSlot keyName="AUDIO_PRESET_RECOMMENDED" /></span>
      </div>
      <div class="preset-card">
        <input type="radio" disabled />
        <span><TranslationSlot keyName="AUDIO_PRESET_DYNAMIC_RANGE" /></span>
      </div>
      <div class="preset-card">
        <input type="radio" disabled />
        <span><TranslationSlot keyName="AUDIO_PRESET_VOCAL_ENHANCEMENT" /></span>
      </div>
      <div class="preset-card">
        <input type="radio" disabled />
        <span><TranslationSlot keyName="AUDIO_PRESET_SMOOTH" /></span>
      </div>
      <div class="preset-card">
        <input type="radio" disabled />
        <span><TranslationSlot keyName="AUDIO_PRESET_CUSTOM" /></span>
      </div>
    </div>

    <!-- Sliders Grid -->
    <div class="custom-grid">
      <div class="control-row">
        <label class="param-label"><TranslationSlot keyName="AUDIO_PARAM_THRESHOLD" /></label>
        <div class="slider-wrapper">
          <input type="range" min="-60" max="0" value="-24" disabled />
          <span class="value-display">-24</span>
          <span class="unit">dB</span>
        </div>
      </div>
      <div class="control-row">
        <label class="param-label"><TranslationSlot keyName="AUDIO_PARAM_KNEE" /></label>
        <div class="slider-wrapper">
          <input type="range" min="0" max="40" value="30" disabled />
          <span class="value-display">30</span>
          <span class="unit">dB</span>
        </div>
      </div>
      <div class="control-row">
        <label class="param-label"><TranslationSlot keyName="AUDIO_PARAM_RATIO" /></label>
        <div class="slider-wrapper">
          <input type="range" min="1" max="20" value="12" disabled />
          <span class="value-display">12</span>
          <span class="unit">:1</span>
        </div>
      </div>
      <div class="control-row">
        <label class="param-label"><TranslationSlot keyName="AUDIO_PARAM_ATTACK" /></label>
        <div class="slider-wrapper">
          <input type="range" min="0" max="1000" value="3" disabled />
          <span class="value-display">3</span>
          <span class="unit">ms</span>
        </div>
      </div>
      <div class="control-row">
        <label class="param-label"><TranslationSlot keyName="AUDIO_PARAM_RELEASE" /></label>
        <div class="slider-wrapper">
          <input type="range" min="0" max="1000" value="250" disabled />
          <span class="value-display">250</span>
          <span class="unit">ms</span>
        </div>
      </div>
    </div>
  </section>

  <!-- Equalizer Section -->
  <section class="panel muted-panel">
    <div class="section-heading">
      <h2 class="section-title"><TranslationSlot keyName="AUDIO_EQUALIZER" /></h2>
    </div>
    <p class="coming-soon-text"><TranslationSlot keyName="AUDIO_COMING_SOON" /></p>
  </section>
</div>

<style>
  .audio-options-wrapper {
    width: 100%;
    max-width: 580px;
    background: #0f172a;
    color: #f8fafc;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    border-radius: 12px;
    padding: 24px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
    border: 1px solid #1e293b;
    text-align: left;
    user-select: none;
  }
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    border-bottom: 1px solid #1e293b;
    padding-bottom: 16px;
    margin-bottom: 20px;
  }
  .back-link {
    font-size: 13px;
    color: #6366f1;
    font-weight: 600;
    cursor: pointer;
  }
  .eyebrow {
    font-size: 10px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: #94a3b8;
    margin: 4px 0 2px 0;
  }
  .page-title {
    font-size: 24px;
    margin: 0;
    font-weight: 800;
  }
  .master-toggle {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 6px;
  }
  .toggle-label {
    font-size: 11px;
    text-transform: uppercase;
    font-weight: bold;
    color: #94a3b8;
  }
  .toggle-switch {
    position: relative;
    display: inline-block;
    width: 44px;
    height: 24px;
  }
  .toggle-switch input {
    opacity: 0;
    width: 0;
    height: 0;
  }
  .slider {
    position: absolute;
    cursor: pointer;
    inset: 0;
    background-color: #6366f1;
    border-radius: 24px;
    transition: 0.3s;
  }
  .slider:before {
    position: absolute;
    content: "";
    height: 18px;
    width: 18px;
    right: 3px;
    bottom: 3px;
    background-color: white;
    border-radius: 50%;
    transition: 0.3s;
  }
  .panel {
    background: #1e293b;
    border: 1px solid #334155;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;
  }
  .muted-panel {
    border-style: dashed;
    background: rgba(30, 41, 59, 0.4);
    opacity: 0.7;
  }
  .section-title {
    font-size: 16px;
    margin: 0;
    font-weight: 700;
  }
  .toggle-label-sm {
    font-size: 12px;
    color: #94a3b8;
    margin-right: 8px;
  }
  .toggle-switch-sm {
    position: relative;
    display: inline-block;
    width: 36px;
    height: 20px;
  }
  .toggle-switch-sm input {
    opacity: 0;
    width: 0;
    height: 0;
  }
  .toggle-switch-sm .slider:before {
    height: 14px;
    width: 14px;
    right: 3px;
    bottom: 3px;
  }
  .preset-group {
    display: flex;
    gap: 8px;
    margin: 16px 0;
    flex-wrap: wrap;
  }
  .preset-card {
    background: #0f172a;
    border: 1px solid #334155;
    padding: 8px 12px;
    border-radius: 6px;
    font-size: 11px;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 6px;
    color: #94a3b8;
  }
  .preset-card.active {
    border-color: #6366f1;
    background: rgba(99, 102, 241, 0.05);
    color: white;
  }
  .custom-grid {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-top: 12px;
    border-top: 1px solid #334155;
    padding-top: 16px;
  }
  .control-row {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
  .param-label {
    font-size: 11px;
    font-weight: 700;
    text-transform: uppercase;
    color: #94a3b8;
  }
  .slider-wrapper {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .slider-wrapper input[type="range"] {
    flex: 1;
    accent-color: #6366f1;
    cursor: pointer;
  }
  .value-display {
    width: 40px;
    text-align: right;
    font-family: monospace;
    font-size: 13px;
    color: white;
  }
  .unit {
    width: 24px;
    font-size: 11px;
    color: #94a3b8;
  }
  .coming-soon-text {
    font-size: 13px;
    color: #94a3b8;
    margin: 8px 0 0 0;
  }
  .flex { display: flex; }
  .items-center { align-items: center; }
  .justify-between { justify-content: space-between; }
</style>
