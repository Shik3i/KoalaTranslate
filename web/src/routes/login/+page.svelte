<script lang="ts">
  import { getContext } from 'svelte';
  import { goto } from '$app/navigation';

  // Get shared userState context
  const userState = getContext<any>('userState');

  let username = $state('');
  let password = $state('');
  let errorMsg = $state('');
  let submitting = $state(false);

  async function handleLogin(e: Event) {
    e.preventDefault();
    if (!username || !password) {
      errorMsg = 'Please enter both username and password.';
      return;
    }

    submitting = true;
    errorMsg = '';

    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      });

      if (res.ok) {
        const data = await res.json();
        userState.authenticated = true;
        userState.username = data.username;
        await goto('/');
      } else {
        const data = await res.json();
        errorMsg = data.error || 'Invalid credentials.';
      }
    } catch (err) {
      errorMsg = 'Server connection error. Please try again.';
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>Login · KoalaTranslate</title>
</svelte:head>

<main class="login-container">
  <div class="login-card glass-panel">
    <div class="logo-area">
      <div class="logo-icon">🐨</div>
      <h1>KoalaTranslate</h1>
      <p class="subtitle text-muted">Translation portal for KoalaSync</p>
    </div>

    {#if errorMsg}
      <div class="error-banner">
        <span>⚠️</span>
        <p>{errorMsg}</p>
      </div>
    {/if}

    <form onsubmit={handleLogin}>
      <div class="form-group">
        <label for="username">Username</label>
        <input 
          id="username"
          type="text" 
          bind:value={username} 
          placeholder="Enter username" 
          disabled={submitting} 
          required 
        />
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input 
          id="password" 
          type="password" 
          bind:value={password} 
          placeholder="Enter password" 
          disabled={submitting} 
          required 
        />
      </div>

      <button type="submit" class="btn-primary w-full" disabled={submitting}>
        {#if submitting}
          Logging in...
        {:else}
          Sign In
        {/if}
      </button>
    </form>

    <div class="footer text-muted">
      Admin-only access portal. Registration disabled.
    </div>
  </div>
</main>

<style>
  .login-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    padding: 20px;
    background: radial-gradient(circle at 10% 20%, rgba(99, 102, 241, 0.08) 0%, transparent 40%),
                radial-gradient(circle at 90% 80%, rgba(168, 85, 247, 0.08) 0%, transparent 40%),
                #090d16;
  }
  .login-card {
    width: 100%;
    max-width: 400px;
    padding: 40px 32px;
  }
  .logo-area {
    text-align: center;
    margin-bottom: 30px;
  }
  .logo-icon {
    font-size: 48px;
    margin-bottom: 12px;
  }
  .logo-area h1 {
    font-size: 24px;
    font-weight: 800;
    margin: 0 0 6px 0;
    letter-spacing: -0.5px;
    background: linear-gradient(135deg, #a855f7, #6366f1);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
  .subtitle {
    font-size: 13px;
    margin: 0;
  }
  .form-group {
    margin-bottom: 20px;
  }
  .error-banner {
    display: flex;
    gap: 12px;
    align-items: center;
    background: rgba(239, 68, 68, 0.15);
    border: 1px solid rgba(239, 68, 68, 0.3);
    padding: 12px 16px;
    border-radius: var(--radius-md);
    margin-bottom: 24px;
    font-size: 13px;
  }
  .error-banner p {
    margin: 0;
    color: #fca5a5;
  }
  .footer {
    text-align: center;
    font-size: 11px;
    margin-top: 30px;
    border-top: 1px solid rgba(255, 255, 255, 0.05);
    padding-top: 20px;
  }
</style>
