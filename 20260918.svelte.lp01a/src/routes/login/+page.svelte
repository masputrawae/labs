<script lang="ts">
	import z from 'zod';
	import InputField from '$lib/components/InputField.svelte';
	import { onMount } from 'svelte';
	import { isAuth } from '$lib/api/auth';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	const inputSchema = z.object({
		username: z.string().min(3).max(20),
		password: z.string().min(8)
	});

	type Input = z.infer<typeof inputSchema>;
	type Errors = Partial<Record<keyof Input, string>>;

	let loading = $state(true);
	let input = $state<Input>({
		username: '',
		password: ''
	});

	let errors = $state<Errors>({});

	const post = (event: SubmitEvent) => {
		event.preventDefault();
		const result = inputSchema.safeParse(input);
		if (!result.success) {
			const flattened = z.flattenError(result.error);
			errors = Object.fromEntries(
				Object.entries(flattened.fieldErrors).map(([key, messages]) => [key, messages?.[0] ?? ''])
			) as Errors;
			return;
		}

		try {
			errors = {};
			console.log(input);
		} catch (error) {
			console.log(error);
		} finally {
			loading = false;
		}
	};

	onMount(async () => {
		try {
			const ok = await isAuth();
			if (ok) {
				await goto(resolve('/'));
			}
		} catch (error) {
			console.log(error);
		} finally {
			loading = false;
		}
	});
</script>

{#if loading}
	<p>loading</p>
{:else}
	<form onsubmit={post}>
		<header>
			<h1>Login</h1>
			<p>Hi, welcome back. To verify your identity, please enter your username and password.</p>
		</header>

		<main>
			<InputField
				label="Username"
				id="registerUsername"
				name="username"
				type="text"
				bind:value={input.username}
				error={errors.username}
			/>

			<InputField
				label="Password"
				id="registerPassword"
				name="password"
				type="password"
				bind:value={input.password}
				error={errors.password}
			/>
		</main>

		<footer>
			<button type="submit">Login</button>
			<p>
				Don't have an account yet? Sign up
				<a href={resolve('/register')}>here</a>.
			</p>
		</footer>
	</form>
{/if}
