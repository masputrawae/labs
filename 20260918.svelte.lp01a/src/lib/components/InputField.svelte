<script lang="ts">
	import { Eye, EyeOff } from '@lucide/svelte';

	type Props = {
		label: string;
		id: string;
		name: string;
		type?: 'text' | 'email' | 'password';
		error?: string;
	};

	let {
		label,
		id,
		name,
		type = 'text',
		error = '',
		value = $bindable('')
	}: Props & { value?: string } = $props();

	let showPassword = $state(false);

	let isPassword = $derived(type === 'password');
	let inputType = $derived(isPassword && !showPassword ? 'password' : 'text');
</script>

<div>
	<label for={id}>{label}</label>

	<div>
		<input
			type={inputType}
			name={name}
			id={id}
			bind:value
			aria-invalid={Boolean(error)}
			aria-describedby={error ? `${id}-error` : undefined}
		/>

		{#if isPassword}
			<button
				type="button"
				aria-label={showPassword ? 'Hide password' : 'Show password'}
				onclick={() => (showPassword = !showPassword)}
			>
				{#if showPassword}
					<EyeOff />
				{:else}
					<Eye />
				{/if}
			</button>
		{/if}
	</div>

	{#if error}
		<p id={`${id}-error`}>{error}</p>
	{/if}
</div>
