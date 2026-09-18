<script lang="ts">
	import z from 'zod';
	import InputField from '$lib/components/InputField.svelte';
	import { onMount } from 'svelte';
	import { isAuth } from '$lib/api/auth';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	const inputSchema = z
		.object({
			email: z.email(),
			username: z.string().min(3).max(20),
			password: z.string().min(8),
			confirmPassword: z.string()
		})
		.refine(({ password, confirmPassword }) => password === confirmPassword, {
			message: 'incorrect password',
			path: ['confirmPassword']
		});

	type Input = z.infer<typeof inputSchema>;
	type Errors = Partial<Record<keyof Input, string>>;

	let loading = $state(true);
	let input = $state<Input>({
		email: '',
		username: '',
		password: '',
		confirmPassword: ''
	});

	let errors = $state<Errors>({});

	const post = (event: SubmitEvent) => {
		event.preventDefault();
		loading = true;

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
			<h1>Create Account</h1>
			<p>Hello and welcome. To proceed with the registration, please fill in the details below.</p>
		</header>

		<main>
			<InputField
				label="Email"
				id="registerEmail"
				name="email"
				type="email"
				bind:value={input.email}
				error={errors.email}
			/>

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

			<InputField
				label="Confirm Password"
				id="registerConfirmPassword"
				name="confirmPassword"
				type="password"
				bind:value={input.confirmPassword}
				error={errors.confirmPassword}
			/>
		</main>

		<footer>
			<button type="submit">Register</button>
			<p>
				Do you already have an account? Log in <a href={resolve('/login')}>here</a>.
			</p>
		</footer>
	</form>
{/if}
