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
      <div class="tab-btn active"><TranslationSlot keyName="TAB_ROOM" /></div>
      <div class="tab-btn"><TranslationSlot keyName="TAB_SYNC" /></div>
      <div class="tab-btn"><TranslationSlot keyName="TAB_SETTINGS" /></div>
      <div class="tab-btn"><TranslationSlot keyName="TAB_STATUS" /></div>
    </div>
  </div>

  <!-- Content Area focusing on Relay Server Settings -->
  <div class="popup-content">
    <div class="highlight-container">
      <div class="server-picker-card">
        <label><TranslationSlot keyName="LABEL_SERVER" /></label>
        <div class="server-buttons">
          <button class="server-btn"><TranslationSlot keyName="BTN_SERVER_OFFICIAL" /></button>
          <button class="server-btn active"><TranslationSlot keyName="BTN_SERVER_CUSTOM" /></button>
        </div>
        
        <div class="input-highlight">
          <TranslationSlot keyName="PLACEHOLDER_SERVER_URL" />
          <input 
            type="text" 
            class="server-url-input" 
            value="wss://custom-relay-server.koalastuff.net" 
            placeholder="wss://your-server:3000"
            readonly
          />
        </div>
      </div>
    </div>

    <!-- Background join button to give context -->
    <button class="join-btn" disabled>
      <TranslationSlot keyName="BTN_JOIN_ROOM" />
    </button>
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
  .highlight-container {
    padding: 2px;
    border: 2px solid #a855f7;
    border-radius: 10px;
    margin-bottom: 14px;
    box-shadow: 0 0 12px rgba(168, 85, 247, 0.3);
  }
  .server-picker-card {
    background: #1e293b;
    border-radius: 8px;
    padding: 10px;
  }
  .server-buttons {
    display: flex;
    gap: 4px;
    margin-bottom: 8px;
  }
  .server-btn {
    flex: 1;
    padding: 6px;
    font-size: 11px;
    background: transparent;
    color: #94a3b8;
    border: none;
    border-radius: 6px;
    font-weight: 600;
  }
  .server-btn.active {
    background: #6366f1;
    color: white;
  }
  .input-highlight {
    background: rgba(99, 102, 241, 0.05);
    border-radius: 6px;
  }
  .server-url-input {
    width: 100%;
    background: #0f172a;
    border: 1px solid #6366f1;
    color: white;
    padding: 8px;
    border-radius: 6px;
    font-size: 11px;
    outline: none;
  }
  .join-btn {
    width: 100%;
    padding: 12px;
    background: #334155;
    border: none;
    color: #94a3b8;
    font-weight: 700;
    border-radius: 8px;
    cursor: not-allowed;
    opacity: 0.6;
  }
</style>
