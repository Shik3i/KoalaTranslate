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
      <div class="popup-version">
        v2.5.3
      </div>
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
      <button class="create-room-btn">
        <TranslationSlot keyName="BTN_CREATE_ROOM" />
      </button>

      <div class="advanced-details">
        <div class="details-summary">
          <TranslationSlot keyName="MANUAL_CONNECT_HEADER" />
          <span class="chevron">▼</span>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  /* Localized styles matching KoalaSync popup */
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
    box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
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
  }
  .chevron {
    font-size: 8px;
  }
</style>
