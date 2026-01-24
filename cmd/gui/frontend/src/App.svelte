<script>
  import { onMount } from 'svelte';
  import { GetEntries, GetTodos, UpdateTodoStatus } from './wailsjs/go/main/App';

  let entries = [];
  let todos = [];
  let view = 'entries';
  let loading = true;

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    loading = true;
    try {
      entries = (await GetEntries()) || [];
      todos = (await GetTodos()) || [];
    } catch (err) {
      console.error('Failed to load data:', err);
    }
    loading = false;
  }

  async function cycleTodoStatus(todo) {
    const cycle = { open: 'next', next: 'done', done: 'open' };
    const newStatus = cycle[todo.status] || 'open';
    await UpdateTodoStatus(todo.id, newStatus);
    await loadData();
  }

  function formatDate(timestamp) {
    if (!timestamp) return '';
    return new Date(timestamp).toLocaleDateString();
  }

  function statusIndicator(status) {
    switch (status) {
      case 'next': return '[>]';
      case 'done': return '[x]';
      default: return '[ ]';
    }
  }
</script>

<main>
  <header>
    <h1>amos</h1>
    <nav>
      <button class:active={view === 'entries'} on:click={() => view = 'entries'}>
        entries ({entries.length})
      </button>
      <button class:active={view === 'todos'} on:click={() => view = 'todos'}>
        todos ({todos.filter(t => t.status !== 'done').length})
      </button>
    </nav>
  </header>

  {#if loading}
    <p>loading...</p>
  {:else if view === 'entries'}
    <section class="entries">
      {#if entries.length === 0}
        <p class="empty">no entries yet</p>
      {:else}
        <ul>
          {#each entries as entry}
            <li>
              <span class="date">{formatDate(entry.created_at)}</span>
              <strong>{entry.title}</strong>
              {#if entry.tags && entry.tags.length > 0}
                <span class="tags">{entry.tags.map(t => '@' + t).join(' ')}</span>
              {/if}
            </li>
          {/each}
        </ul>
      {/if}
    </section>
  {:else}
    <section class="todos">
      {#if todos.length === 0}
        <p class="empty">no todos yet</p>
      {:else}
        <ul>
          {#each todos as todo}
            <li class={todo.status}>
              <button class="status-btn" on:click={() => cycleTodoStatus(todo)}>
                {statusIndicator(todo.status)}
              </button>
              <span class="todo-title">{todo.title}</span>
              {#if todo.tags && todo.tags.length > 0}
                <span class="tags">{todo.tags.map(t => '@' + t).join(' ')}</span>
              {/if}
            </li>
          {/each}
        </ul>
      {/if}
    </section>
  {/if}
</main>

<style>
  :global(body) {
    font-family: monospace;
    background: #000;
    color: #fff;
    margin: 0;
    padding: 0;
  }

  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
  }

  header {
    border-bottom: 1px solid #fff;
    padding-bottom: 10px;
    margin-bottom: 20px;
  }

  h1 {
    margin: 0 0 10px 0;
    font-size: 24px;
  }

  nav button {
    background: none;
    border: 1px solid #fff;
    color: #fff;
    font-family: monospace;
    padding: 5px 10px;
    margin-right: 10px;
    cursor: pointer;
  }

  nav button.active {
    background: #fff;
    color: #000;
  }

  ul {
    list-style: none;
    padding: 0;
  }

  li {
    padding: 8px 0;
    border-bottom: 1px solid #333;
  }

  li.done {
    opacity: 0.5;
  }

  .date {
    color: #888;
    margin-right: 10px;
  }

  .tags {
    color: #888;
    margin-left: 10px;
  }

  .empty {
    color: #555;
    text-align: center;
    padding: 40px 0;
  }

  .status-btn {
    background: none;
    border: none;
    color: #fff;
    font-family: monospace;
    cursor: pointer;
    padding: 0;
    margin-right: 8px;
  }

  li.next .status-btn {
    color: #0f0;
  }

  li.done .status-btn {
    color: #555;
  }
</style>
