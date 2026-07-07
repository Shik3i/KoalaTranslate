<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import './global.css';

  let { children } = $props();

  let loading = $state(true);
  let authenticated = $state(false);
  let username = $state('');

  // Svelte 5 runes context for active user
  import { setContext } from 'svelte';
  
  class UserState {
    authenticated = $state(false);
    username = $state('');
  }
  
  const userState = new UserState();
  setContext('userState', userState);

  async function checkAuth() {
    try {
      const res = await fetch('/api/auth/me');
      if (res.ok) {
        const data = await res.json();
        userState.authenticated = true;
        userState.username = data.username;
        authenticated = true;
        username = data.username;

        if (page.url.pathname === '/login') {
          await goto('/');
        }
      } else {
        userState.authenticated = false;
        userState.username = '';
        authenticated = false;
        username = '';

        if (page.url.pathname !== '/login') {
          await goto('/login');
        }
      }
    } catch (err) {
      console.error('Auth verification failed:', err);
      if (page.url.pathname !== '/login') {
        await goto('/login');
      }
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    checkAuth();
  });
</script>

<svelte:head>
  <title>KoalaTranslate · In-context translation portal</title>
</svelte:head>

{#if loading}
  <div class="loading-screen">
    <div class="spinner"></div>
    <p>Authenticating session...</p>
  </div>
{:else}
  {@render children()}
{/if}

<style>
  .loading-screen {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background: #090d16;
    color: #94a3b8;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  }
  .spinner {
    width: 40px;
    height: 40px;
    border: 3px solid rgba(99, 102, 241, 0.1);
    border-top-color: #6366f1;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 16px;
  }
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>
