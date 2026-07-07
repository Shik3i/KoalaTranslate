<script lang="ts">
  import TranslationSlot from '../../../components/TranslationSlot.svelte';
</script>

<div class="popup-wrapper">
  <!-- Simulated Extension Header -->
  <div class="popup-header">
    <div class="header-row">
      <h1 class="logo-title">
        <span class="logo-emoji">🐨</span>
        <TranslationSlot keyName="APP_TITLE" />
      </h1>
      <div class="popup-version">v2.5.3</div>
    </div>

    <div class="tabs">
      <div class="tab-btn"><TranslationSlot keyName="TAB_ROOM" /></div>
      <div class="tab-btn"><TranslationSlot keyName="TAB_SYNC" /></div>
      <div class="tab-btn active"><TranslationSlot keyName="TAB_SETTINGS" /></div>
      <div class="tab-btn"><TranslationSlot keyName="TAB_STATUS" /></div>
    </div>
  </div>

  <!-- Content Area -->
  <div class="popup-content">
    <!-- Username -->
    <div class="settings-row">
      <label><TranslationSlot keyName="LABEL_USERNAME" /></label>
      <input type="text" value="Anonymous Koala" class="width-fixed" readonly />
    </div>

    <!-- Privacy Details (Collapsed indicator) -->
    <div class="settings-panel">
      <div class="flex justify-between items-center pointer">
        <span class="bold-text"><TranslationSlot keyName="LABEL_PRIVACY_SETTINGS" /></span>
        <span>▼</span>
      </div>
    </div>

    <!-- Hide Clutter -->
    <div class="settings-row-toggle">
      <label class="margin-zero"><TranslationSlot keyName="LABEL_HIDE_CLUTTER" /></label>
      <label class="toggle-switch">
        <input type="checkbox" checked disabled />
        <span class="slider"></span>
      </label>
    </div>

    <!-- Auto Sync -->
    <div class="settings-row-toggle">
      <label class="margin-zero"><TranslationSlot keyName="LABEL_AUTO_SYNC_NEXT" /></label>
      <label class="toggle-switch">
        <input type="checkbox" disabled />
        <span class="slider"></span>
      </label>
    </div>

    <!-- Auto Copy -->
    <div class="settings-row-toggle">
      <label class="margin-zero"><TranslationSlot keyName="LABEL_AUTO_COPY_INVITE" /></label>
      <label class="toggle-switch">
        <input type="checkbox" checked disabled />
        <span class="slider"></span>
      </label>
    </div>

    <!-- Browser Notifications -->
    <div class="settings-row-toggle">
      <label class="margin-zero"><TranslationSlot keyName="LABEL_NOTIFICATIONS" /></label>
      <label class="toggle-switch">
        <input type="checkbox" disabled />
        <span class="slider"></span>
      </label>
    </div>

    <!-- Audio Settings -->
    <div class="settings-row-toggle bg-card border-card">
      <label class="margin-zero font-size-13"><TranslationSlot keyName="LABEL_AUDIO_PROCESSING" /></label>
      <button class="open-audio-settings-btn" disabled>
        <TranslationSlot keyName="AUDIO_OPEN_SETTINGS" />
      </button>
    </div>

    <!-- App Language -->
    <div class="settings-row bg-card border-card flex-col-align">
      <label><TranslationSlot keyName="LABEL_LANGUAGE" /></label>
      <select disabled class="lang-select">
        <option>🇬🇧 English</option>
        <option>🇩🇪 Deutsch</option>
        <option>🇫🇷 Français</option>
        <option>🇪🇸 Español</option>
      </select>
    </div>

    <!-- Troubleshooting buttons -->
    <div class="troubleshooting-card">
      <button class="secondary-action-btn"><TranslationSlot keyName="BTN_RESTART_TOUR" /></button>
      <button class="secondary-action-btn"><TranslationSlot keyName="BTN_REGEN_ID" /></button>
    </div>
  </div>
</div>

<style>
  .popup-wrapper {
    width: 320px;
    background: #0f172a;
    color: #f8fafc;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    border-radius: 12px;
    padding: 16px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
    border: 1px solid #1e293b;
    user-select: none;
    text-align: left;
    max-height: 480px;
    overflow-y: auto;
  }
  .header-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
  }
  .logo-title {
    font-size: 18px;
    margin: 0;
    color: #6366f1;
    letter-spacing: 1px;
    text-transform: uppercase;
    display: inline-flex;
    align-items: center;
    gap: 8px;
  }
  .logo-emoji {
    font-size: 24px;
  }
  .popup-version {
    font-size: 10px;
    color: #94a3b8;
    opacity: 0.5;
  }
  .tabs {
    display: flex;
    gap: 4px;
    background: #1e293b;
    padding: 4px;
    border-radius: 12px;
    margin-bottom: 20px;
  }
  .tab-btn {
    flex: 1;
    padding: 8px;
    text-align: center;
    color: #94a3b8;
    font-weight: 600;
    font-size: 12px;
    border-radius: 8px;
    cursor: pointer;
  }
  .tab-btn.active {
    background: #6366f1;
    color: #ffffff;
  }
  .settings-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: #1e293b;
    padding: 10px;
    border-radius: 8px;
    margin-bottom: 12px;
    border: 1px solid #334155;
  }
  .settings-row-toggle {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: #1e293b;
    padding: 10px;
    border-radius: 8px;
    margin-bottom: 12px;
    border: 1px solid #334155;
  }
  .width-fixed {
    width: 130px !important;
    padding: 6px !important;
    font-size: 12px !important;
    background: #0f172a !important;
    border: 1px solid #334155 !important;
    color: white !important;
  }
  .settings-panel {
    background: #1e293b;
    padding: 10px;
    border-radius: 8px;
    margin-bottom: 12px;
    border: 1px solid #334155;
    font-size: 12px;
  }
  .bold-text {
    font-weight: 700;
  }
  .margin-zero {
    margin-bottom: 0;
    font-size: 12px;
    text-transform: none;
    font-weight: 600;
  }
  .toggle-switch {
    position: relative;
    display: inline-block;
    width: 36px;
    height: 20px;
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
    background-color: #334155;
    border-radius: 20px;
    transition: 0.3s;
  }
  .slider:before {
    position: absolute;
    content: "";
    height: 14px;
    width: 14px;
    left: 3px;
    bottom: 3px;
    background-color: #94a3b8;
    border-radius: 50%;
    transition: 0.3s;
  }
  input:checked + .slider {
    background-color: #6366f1;
  }
  input:checked + .slider:before {
    transform: translateX(16px);
    background-color: white;
  }
  .font-size-13 {
    font-size: 13px !important;
  }
  .open-audio-settings-btn {
    font-size: 11px;
    color: #6366f1;
    background: rgba(99, 102, 241, 0.12);
    padding: 5px 12px;
    border-radius: 6px;
    font-weight: 600;
    border: none;
  }
  .flex-col-align {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  .lang-select {
    width: 100% !important;
    padding: 6px !important;
    background: #0f172a !important;
    border: 1px solid #334155 !important;
    color: white !important;
    border-radius: 8px !important;
  }
  .troubleshooting-card {
    margin-top: 15px;
    padding: 8px 0;
    border-top: 1px solid #1e293b;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .secondary-action-btn {
    width: 100%;
    padding: 8px;
    background: #334155;
    color: white;
    border-radius: 6px;
    font-size: 11px;
    border: none;
    font-weight: 600;
    text-align: center;
  }
  .flex { display: flex; }
  .items-center { align-items: center; }
  .justify-between { justify-content: space-between; }
  .pointer { cursor: pointer; }
</style>
