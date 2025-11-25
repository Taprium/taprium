<script lang="ts">
	import { Section, Register } from 'flowbite-svelte-blocks';
	import { Button, Checkbox, Label, Input } from 'flowbite-svelte';
	import { onMount } from 'svelte';
	import { pb, PB_COLLECTION_USERS } from '$lib/pb/backend-pb';
	import { goto } from '$app/navigation';

	let email: string, password: string;

	onMount(() => {
		if (pb.authStore.isValid) {
			goto('/');
		}
	});

	async function SignIn() {
		try {
			const authData = await pb.collection(PB_COLLECTION_USERS).authWithPassword(email, password);
			location.reload();
		} catch (e: any) {}
	}
</script>

<Section name="login">
	<Register href="/">
		{#snippet top()}
			Taprium
		{/snippet}
		<div class="space-y-4 p-6 sm:p-8 md:space-y-6">
			<form class="flex flex-col space-y-6" on:submit|preventDefault={SignIn}>
				<h3 class="p-0 text-xl font-medium text-gray-900 dark:text-white">Login</h3>
				<Label class="space-y-2">
					<span>Your email</span>
					<Input
						type="email"
						name="email"
						placeholder="name@company.com"
						required
						bind:value={email}
					/>
				</Label>
				<Label class="space-y-2">
					<span>Your password</span>
					<Input
						type="password"
						name="password"
						placeholder="•••••"
						required
						bind:value={password}
					/>
				</Label>
				<!-- <div class="flex items-start">
					<Checkbox>Remember me</Checkbox>
					<a href="/" class="ml-auto text-sm text-blue-700 hover:underline dark:text-blue-500"
						>Forgot password?</a
					>
				</div> -->
				<Button type="submit" class="w-full1">Sign in</Button>
				<!-- <p class="text-sm font-light text-gray-500 dark:text-gray-400">
					Don’t have an account yet? <a
						href="/"
						class="font-medium text-primary-600 hover:underline dark:text-primary-500">Sign up</a
					>
				</p> -->
			</form>
		</div>
	</Register>
</Section>
