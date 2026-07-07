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
    <!-- Active Room Status -->
    <div class="info-card border-left-accent flex justify-between items-center">
      <div>
        <label><TranslationSlot keyName="LABEL_ACTIVE_ROOM" /></label>
        <div class="active-room-id">KS-7A2D</div>
      </div>
      <div class="active-server"><TranslationSlot keyName="ACTIVE_SERVER_OFFICIAL" /></div>
    </div>

    <!-- Invite Link -->
    <div class="info-card">
      <label><TranslationSlot keyName="LABEL_INVITE_LINK" /></label>
      <div class="invite-box">
        <input type="text" value="https://sync.koalastuff.net/#/room/KS-7A2D" readonly />
        <button class="copy-btn" title="Copy Invite Link">
          <TranslationSlot keyName="BTN_COPY_INVITE" />
        </button>
      </div>
    </div>

    <!-- Host Control -->
    <div class="info-card">
      <div class="flex justify-between items-center">
        <label><TranslationSlot keyName="LABEL_HOST_CONTROL" /></label>
        <span class="badge badge-host"><TranslationSlot keyName="BADGE_HOST" /></span>
      </div>
      <div class="toggle-row">
        <span><TranslationSlot keyName="LABEL_HOST_ONLY_TOGGLE" /></span>
        <label class="toggle-switch">
          <input type="checkbox" checked disabled />
          <span class="slider"></span>
        </label>
      </div>
    </div>

    <!-- Peers List -->
    <div class="peers-section">
      <label><TranslationSlot keyName="LABEL_PEERS_IN_ROOM" /></label>
      <div class="info-card peer-list">
        <div class="peer-item flex justify-between items-center">
          <div class="peer-info">
            <span class="status-dot online"></span>
            <span>Anonymous Koala (<TranslationSlot keyName="LABEL_YOU" />)</span>
          </div>
          <span class="badge"><TranslationSlot keyName="BADGE_HOST" /></span>
        </div>
        <div class="peer-item flex justify-between items-center">
          <div class="peer-info">
            <span class="status-dot online"></span>
            <span>Friendly Panda</span>
          </div>
          <span class="badge"><TranslationSlot keyName="BADGE_GUEST" /></span>
        </div>
      </div>
    </div>

    <!-- Leave Room -->
    <button class="leave-btn">
      <TranslationSlot keyName="BTN_LEAVE_ROOM" />
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
  .info-card {
    background: #1e293b;
    padding: 12px;
    border-radius: 8px;
    margin-bottom: 12px;
    border: 1px solid #334155;
  }
  .border-left-accent {
    border-left: 4px solid #6366f1;
  }
  .active-room-id {
    font-weight: 700;
    color: #6366f1;
    font-size: 16px;
    letter-spacing: 1px;
  }
  .active-server {
    font-size: 10px;
    color: #94a3b8;
  }
  .invite-box {
    display: flex;
    gap: 8px;
    margin-top: 8px;
  }
  .invite-box input {
    flex: 1;
    font-size: 11px;
    background: #0f172a;
    border: 1px solid #334155;
    color: white;
    padding: 8px;
    border-radius: 6px;
  }
  .copy-btn {
    padding: 8px;
    background: #334155;
    border: none;
    color: white;
    border-radius: 6px;
    font-size: 11px;
  }
  .badge {
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
    padding: 2px 8px;
    border-radius: 8px;
    background: #334155;
    color: #f8fafc;
  }
  .badge-host {
    background: #6366f1;
  }
  .toggle-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 12px;
    font-size: 12px;
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
    background-color: #6366f1;
    border-radius: 20px;
    transition: 0.3s;
  }
  .slider:before {
    position: absolute;
    content: "";
    height: 14px;
    width: 14px;
    right: 3px;
    bottom: 3px;
    background-color: white;
    border-radius: 50%;
    transition: 0.3s;
  }
  .peer-list {
    padding: 4px 12px;
    max-height: 120px;
    overflow-y: auto;
  }
  .peer-item {
    padding: 8px 0;
    border-bottom: 1px solid #334155;
    font-size: 12px;
  }
  .peer-item:last-child {
    border-bottom: none;
  }
  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    display: inline-block;
    margin-right: 6px;
  }
  .status-dot.online {
    background: #22c55e;
    box-shadow: 0 0 8px #22c55e;
  }
  .peer-info {
    display: flex;
    align-items: center;
  }
  .leave-btn {
    width: 100%;
    padding: 10px;
    background: #1e293b;
    border: 1px solid #ef4444;
    color: #ef4444;
    font-weight: 600;
    border-radius: 8px;
    cursor: pointer;
    margin-top: 8px;
  }
  .flex { display: flex; }
  .justify-between { justify-content: space-between; }
  .items-center { align-items: center; }
</style>
