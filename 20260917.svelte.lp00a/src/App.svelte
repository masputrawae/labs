<script>
  import { onMount } from "svelte";

  let todos = $state([]);
  let newTask = $state("");
  let editingId = $state(null);
  let editTask = $state("");
  let errorMessage = $state("");

  const jsonHeaders = {
    "Content-Type": "application/json",
  };

  async function getTodos() {
    try {
      errorMessage = "";

      const res = await fetch("/api/todos");

      if (!res.ok) {
        throw new Error("Gagal mengambil data todo");
      }

      const json = await res.json();
      todos = json.data;
    } catch (error) {
      errorMessage = error.message;
    }
  }

  async function addTodo() {
    const task = newTask.trim();

    if (!task) return;

    try {
      const res = await fetch("/api/todos", {
        method: "POST",
        headers: jsonHeaders,
        body: JSON.stringify({ task }),
      });

      if (!res.ok) {
        throw new Error("Gagal menambahkan todo");
      }

      const json = await res.json();

      todos.push(json.data);
      newTask = "";
    } catch (error) {
      errorMessage = error.message;
    }
  }

  function startEditing(todo) {
    editingId = todo.id;
    editTask = todo.task;
  }

  function cancelEditing() {
    editingId = null;
    editTask = "";
  }

  async function updateTodo(id) {
    const task = editTask.trim();

    if (!task) return;

    try {
      const res = await fetch(`/api/todos/${id}`, {
        method: "PATCH",
        headers: jsonHeaders,
        body: JSON.stringify({ task }),
      });

      if (!res.ok) {
        throw new Error("Gagal mengubah todo");
      }

      const json = await res.json();

      todos = todos.map((todo) => (todo.id === id ? json.data : todo));

      cancelEditing();
    } catch (error) {
      errorMessage = error.message;
    }
  }

  async function toggleTodo(todo) {
    try {
      const res = await fetch(`/api/todos/${todo.id}`, {
        method: "PATCH",
        headers: jsonHeaders,
        body: JSON.stringify({
          isDone: !todo.isDone,
        }),
      });

      if (!res.ok) {
        throw new Error("Gagal mengubah status todo");
      }

      const json = await res.json();

      todos = todos.map((item) => (item.id === todo.id ? json.data : item));
    } catch (error) {
      errorMessage = error.message;
    }
  }

  async function deleteTodo(id) {
    try {
      const res = await fetch(`/api/todos/${id}`, {
        method: "DELETE",
      });

      if (!res.ok) {
        throw new Error("Gagal menghapus todo");
      }

      todos = todos.filter((todo) => todo.id !== id);

      if (editingId === id) {
        cancelEditing();
      }
    } catch (error) {
      errorMessage = error.message;
    }
  }

  onMount(() => {
    getTodos();
  });
</script>

<form
  onsubmit={(event) => {
    event.preventDefault();
    addTodo();
  }}
>
  <label for="new-task">Task:</label>

  <input
    id="new-task"
    type="text"
    bind:value={newTask}
    placeholder="Masukkan task"
    required
  />

  <button type="submit">Save</button>
</form>

{#if errorMessage}
  <p role="alert">{errorMessage}</p>
{/if}

{#if todos.length === 0}
  <p>Belum ada todo.</p>
{:else}
  <ul>
    {#each todos as todo (todo.id)}
      <li id={`todo-id-${todo.id}`}>
        {#if editingId === todo.id}
          <label for={`edit-task-${todo.id}`}> Edit task: </label>

          <input
            id={`edit-task-${todo.id}`}
            type="text"
            bind:value={editTask}
            required
          />

          <button type="button" onclick={() => updateTodo(todo.id)}>
            Update
          </button>

          <button type="button" onclick={cancelEditing}> Cancel </button>
        {:else}
          <p>
            Task:
            {todo.task}
          </p>

          <button type="button" onclick={() => startEditing(todo)}>
            Edit
          </button>
        {/if}

        <p>
          Status:
          <button type="button" onclick={() => toggleTodo(todo)}>
            {todo.isDone ? "Done" : "Not done"}
          </button>
        </p>

        <p>Created At: {todo.createdAt}</p>
        <p>Updated At: {todo.updatedAt}</p>

        <button type="button" onclick={() => deleteTodo(todo.id)}>
          Delete
        </button>
      </li>
    {/each}
  </ul>
{/if}
