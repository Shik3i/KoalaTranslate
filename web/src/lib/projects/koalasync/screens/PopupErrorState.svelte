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
      <div class="tab-btn"><TranslationSlot keyName="TAB_SETTINGS" /></div>
      <div class="tab-btn active"><TranslationSlot keyName="TAB_STATUS" /></div>
    </div>
  </div>

  <!-- Content Area -->
  <div class="popup-content">
    <!-- Connection Status Error -->
    <label><TranslationSlot keyName="LABEL_CONN_STATUS" /></label>
    <div class="info-card flex justify-between items-center">
      <div class="flex items-center">
        <span class="status-dot offline"></span>
        <span class="status-text text-error"><TranslationSlot keyName="CONN_STATUS_DISCONNECTED" /></span>
      </div>
      <button class="retry-btn">
        <TranslationSlot keyName="BTN_RETRY" />
      </button>
    </div>

    <!-- Video Debug Info -->
    <label><TranslationSlot keyName="LABEL_VIDEO_DEBUG" /></label>
    <div class="info-card debug-card text-muted">
      <TranslationSlot keyName="VIDEO_DEBUG_EMPTY" />
    </div>

    <!-- Technical Logs -->
    <div class="flex justify-between items-center header-margin">
      <label><TranslationSlot keyName="LABEL_LOGS" /></label>
      <span class="clear-btn"><TranslationSlot keyName="BTN_CLEAR" /></span>
    </div>
    <div class="log-list">
      <div class="log-entry log-error">[13:02:44] WebSocket connection failed to wss://sync.koalastuff.net:3000</div>
      <div class="log-entry log-error">[13:02:45] ERR_CONN_TIMEOUT: Connection timed out. Please try again.</div>
      <div class="log-entry">[13:02:45] Reconnection loop starting (retry 1 of 5)...</div>
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
  .info-card {
    background: #1e293b;
    padding: 12px;
    border-radius: 8px;
    margin-bottom: 14px;
    border: 1px solid #334155;
  }
  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    display: inline-block;
    margin-right: 6px;
  }
  .status-dot.offline {
    background: #ef4444;
    box-shadow: 0 0 8px #ef4444;
  }
  .status-text {
    font-weight: bold;
    font-size: 13px;
  }
  .text-error {
    color: #ef4444;
  }
  .retry-btn {
    padding: 4px 10px;
    font-size: 11px;
    background: #334155;
    border: none;
    color: white;
    font-weight: bold;
    border-radius: 6px;
    cursor: pointer;
  }
  .debug-card {
    font-size: 10px;
    font-family: monospace;
    min-height: 48px;
    line-height: 1.4;
  }
  .header-margin {
    margin-top: 10px;
    margin-bottom: 6px;
  }
  .clear-btn {
    color: #6366f1;
    font-size: 10px;
    cursor: pointer;
    font-weight: bold;
  }
  .log-list {
    background: #000000;
    font-family: monospace;
    font-size: 10px;
    padding: 8px;
    border-radius: 8px;
    max-height: 100px;
    overflow-y: auto;
    border: 1px solid #334155;
  }
  .log-entry {
    margin-bottom: 4px;
    color: #cccccc;
  }
  .log-error {
    color: #ef4444;
  }
  .flex { display: flex; }
  .items-center { align-items: center; }
  .justify-between { justify-content: space-between; }
</style>
