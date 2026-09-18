<script>
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { isAuth } from '$lib/api/auth';
	import { onMount } from 'svelte';

	let loading = $state(true);
	onMount(async () => {
		try {
			const ok = await isAuth();
			if (!ok) {
				await goto(resolve('/login'));
			}
		} catch (error) {
			console.log(error);
		} finally {
			loading = false;
		}
	});
</script>

{#if loading}
	<p>Loading</p>
{:else}
	<h1>Hay Welcome</h1>
{/if}
