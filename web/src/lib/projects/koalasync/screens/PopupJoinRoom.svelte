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

  <!-- Content Area -->
  <div class="popup-content">
    <div class="section-join">
      <button class="create-room-btn secondary-btn" disabled>
        <TranslationSlot keyName="BTN_CREATE_ROOM" />
      </button>

      <div class="advanced-details expanded">
        <div class="details-summary">
          <TranslationSlot keyName="MANUAL_CONNECT_HEADER" />
          <span class="chevron">▲</span>
        </div>

        <div class="details-content">
          <!-- Server Picker -->
          <div class="server-picker-card">
            <label><TranslationSlot keyName="LABEL_SERVER" /></label>
            <div class="server-buttons">
              <button class="server-btn"><TranslationSlot keyName="BTN_SERVER_OFFICIAL" /></button>
              <button class="server-btn active"><TranslationSlot keyName="BTN_SERVER_CUSTOM" /></button>
            </div>
            <input 
              type="text" 
              class="server-url-input" 
              value="wss://your-server:3000" 
              placeholder="wss://your-server:3000"
              readonly
            />
          </div>

          <!-- Room ID -->
          <div class="form-group">
            <label><TranslationSlot keyName="LABEL_ROOM_ID" /></label>
            <input type="text" value="KS-1234" readonly />
          </div>

          <!-- Password -->
          <div class="form-group">
            <label><TranslationSlot keyName="LABEL_PASSWORD" /></label>
            <input type="password" value="secretpassword" readonly />
          </div>

          <!-- Join Button -->
          <button class="join-btn">
            <TranslationSlot keyName="BTN_JOIN_ROOM" />
          </button>

          <!-- Public Rooms -->
          <div class="public-rooms-header">
            <label><TranslationSlot keyName="LABEL_PUBLIC_ROOMS" /></label>
            <span class="refresh-link"><TranslationSlot keyName="BTN_REFRESH" /></span>
          </div>

          <div class="public-rooms-card">
            <div class="public-room-item">
              <span class="room-name">KS-9901</span>
              <span class="room-count">3 peers</span>
            </div>
            <div class="public-room-item">
              <span class="room-name">MovieNight</span>
              <span class="room-count">1 peer</span>
            </div>
          </div>
        </div>
      </div>
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
  .create-room-btn {
    width: 100%;
    padding: 14px;
    font-size: 15px;
    font-weight: 700;
    border-radius: 8px;
    border: none;
    background: linear-gradient(135deg, #6366f1, #a855f7);
    color: white;
    cursor: pointer;
    margin-bottom: 16px;
  }
  .create-room-btn.secondary-btn {
    background: #1e293b;
    border: 1px dashed #334155;
    color: #94a3b8;
    box-shadow: none;
    opacity: 0.6;
  }
  .advanced-details {
    background: #1e293b;
    border: 1px solid #334155;
    border-radius: 8px;
    padding: 8px;
  }
  .details-summary {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 11px;
    font-weight: 700;
    color: #94a3b8;
    text-transform: uppercase;
    cursor: pointer;
    margin-bottom: 12px;
  }
  .chevron {
    font-size: 8px;
  }
  .server-picker-card {
    background: #0f172a;
    border: 1px solid #334155;
    border-radius: 8px;
    padding: 10px;
    margin-bottom: 14px;
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
  .server-url-input {
    width: 100%;
    background: #1e293b;
    border: 1px solid #334155;
    color: white;
    padding: 8px;
    border-radius: 6px;
    font-size: 11px;
  }
  .form-group {
    margin-bottom: 12px;
  }
  .form-group input {
    width: 100%;
    background: #1e293b;
    border: 1px solid #334155;
    color: white;
    padding: 10px;
    border-radius: 8px;
    font-size: 12px;
  }
  .join-btn {
    width: 100%;
    padding: 12px;
    background: #6366f1;
    border: none;
    color: white;
    font-weight: 700;
    border-radius: 8px;
    cursor: pointer;
    margin-top: 8px;
  }
  .public-rooms-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 16px;
    margin-bottom: 8px;
  }
  .public-rooms-header label {
    margin: 0;
  }
  .refresh-link {
    color: #6366f1;
    font-size: 10px;
    font-weight: bold;
    cursor: pointer;
  }
  .public-rooms-card {
    background: #0f172a;
    border-radius: 6px;
    padding: 4px;
    max-height: 80px;
    overflow-y: auto;
    border: 1px solid #334155;
  }
  .public-room-item {
    display: flex;
    justify-content: space-between;
    padding: 6px 8px;
    font-size: 11px;
    border-bottom: 1px solid #1e293b;
  }
  .public-room-item:last-child {
    border-bottom: none;
  }
  .room-name {
    color: #6366f1;
    font-weight: bold;
  }
  .room-count {
    color: #94a3b8;
  }
</style>
