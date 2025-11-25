<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import {
		Avatar,
		Button,
		Dropdown,
		DropdownGroup,
		DropdownHeader,
		DropdownItem,
		Navbar,
		NavBrand,
		NavHamburger,
		NavLi,
		NavUl
	} from 'flowbite-svelte';
	import { pb } from '$lib/pb/backend-pb';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	let { children } = $props();

	async function SignOut() {
		pb.authStore.clear();
		location.reload();
	}

	$effect(() => {
		if (!pb.authStore.isValid && page.url.pathname != '/login') {
			goto('/login');
		}
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Navbar>
	<NavBrand href="/">
		<!-- <img src="/images/flowbite-svelte-icon-logo.svg" class="me-3 h-6 sm:h-9" alt="Flowbite Logo" /> -->
		<span class="self-center text-xl font-semibold whitespace-nowrap dark:text-white">
			Taprium
		</span>
	</NavBrand>
	<div class="flex items-center md:order-2">
		{#if pb.authStore.isValid}
			<Avatar id="avatar-menu" />
			<NavHamburger />
		{:else}
			<Button href="/login">Login</Button>
		{/if}
	</div>
	{#if pb.authStore.isValid}
		<Dropdown placement="bottom" triggeredBy="#avatar-menu">
			<!-- <DropdownHeader>
			<span class="block text-sm">Bonnie Green</span>
			<span class="block truncate text-sm font-medium">name@flowbite.com</span>
		</DropdownHeader> -->
			<DropdownGroup>
				<!-- <DropdownItem>Dashboard</DropdownItem> -->
				<DropdownItem href="/settings">Settings</DropdownItem>
				<!-- <DropdownItem>Earnings</DropdownItem> -->
			</DropdownGroup>
			<DropdownItem onclick={SignOut}>Sign out</DropdownItem>
		</Dropdown>
	{/if}

	<NavUl>
		<NavLi href="/img_gen">Image Gen</NavLi>
		<NavLi href="/upscale_runners">Upscale Runners</NavLi>
	</NavUl>
</Navbar>

{@render children()}
