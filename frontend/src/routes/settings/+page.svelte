<script lang="ts">
	import { pb, PB_COLLECTION_SETTINGS } from '$lib/pb/backend-pb';
	import { GetDefaultSettings } from '$lib/pb/default-settings';
	import {
		Button,
		ButtonGroup,
		Dropdown,
		DropdownItem,
		Heading,
		Input,
		InputAddon,
		Label,
		Spinner
	} from 'flowbite-svelte';
	import { Section } from 'flowbite-svelte-blocks';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	import type { RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';

	let loading = $state(true);
	let settingsRecord = $state<RecordModel>();

	function getSettings() {
		GetDefaultSettings()
			.then((v) => {
				settingsRecord = v;
			})
			.finally(() => {
				loading = false;
			});
	}

	let backends: string[] = $state([]);

	async function getSettingsOptions() {
		const result = await pb.send('/select-options/settings', {});
		console.log(result);
		backends = result.image_backend;
	}

	onMount(() => {
		getSettings();
		getSettingsOptions();
	});

	async function handleSave(e: SubmitEvent) {
		e.preventDefault();
		if (settingsRecord == undefined) return;
		await pb.collection(PB_COLLECTION_SETTINGS).update(settingsRecord!.id, settingsRecord);
		loading = true;
		getSettings();
	}
</script>

<svelte:head>
	<title>Settings - Taprium</title>
</svelte:head>

<Section>
	{#if loading}
		<Spinner />
	{:else if settingsRecord != undefined}
		<form class="flex flex-col space-y-6" onsubmit={handleSave}>
			<div class="flex items-center justify-between">
				<Heading class="mb-5">Settings</Heading>
				<Button type="submit">Save</Button>
			</div>
			<Label class="space-y-2">
				<span>Default Image Size</span>
				<br />
				<ButtonGroup>
					<InputAddon>Width</InputAddon>
					<Input type="number" min={512} bind:value={settingsRecord.img_width} />
				</ButtonGroup>
				<span class="mx-2">*</span>
				<ButtonGroup>
					<InputAddon>Height</InputAddon>
					<Input type="number" min={512} bind:value={settingsRecord.img_height} />
				</ButtonGroup>
			</Label>
			<Label class="space-y-2">
				<span>Default Queue Image Count</span>
				<Input type="number" min={1} bind:value={settingsRecord.default_queue_count} />
			</Label>
			<Label class="space-y-2">
				<span>Default Upscale Times</span>
				<Input type="number" min={1} bind:value={settingsRecord.upscale_times} />
			</Label>
			<Label class="space-y-2">
				<span>Upscale Timeout Seconds</span>
				<Input type="number" min={1} bind:value={settingsRecord.upscale_timeout_in_second} />
			</Label>

			<Label class="space-y-2">
				<span>Image generation backend</span>
				<br />
				<Button>
					{settingsRecord.image_backend.length == 0
						? 'cloudflare-worker-ai'
						: settingsRecord.image_backend}
					<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" />
				</Button>
				<Dropdown simple>
					{#each backends as iBackend}
						<DropdownItem
							onclick={(v: any) => {
								settingsRecord!.image_backend = iBackend;
							}}>{iBackend}</DropdownItem
						>
					{/each}
					<!-- <DropdownItem>Dashboard</DropdownItem>
				<DropdownItem>Settings</DropdownItem>
				<DropdownItem>Earnings</DropdownItem>
				<DropdownItem>Sign out</DropdownItem> -->
				</Dropdown>
			</Label>
		</form>
	{:else}{/if}
</Section>
