<script lang="ts">
	import { page } from '$app/state';
	import { PUBLIC_PB_ADDR } from '$env/static/public';
	import {
		pb,
		PB_COLLECTION_IMAGE_QUEUES,
		PB_COLLECTION_GENERATED_IMAGES
	} from '$lib/pb/backend-pb';
	import { GetDefaultSettings } from '$lib/pb/default-settings';
	import {
		Badge,
		Button,
		Card,
		Heading,
		Hr,
		Img,
		Input,
		Label,
		Modal,
		Spinner,
		Toggle
	} from 'flowbite-svelte';
	import { Section } from 'flowbite-svelte-blocks';
	import { RefreshOutline } from 'flowbite-svelte-icons';
	import { type RecordModel } from 'pocketbase';
	import { onMount, tick } from 'svelte';

	let loading = $state(false);
	let queueData = $state<RecordModel>();
	let defaultSettings = $state<RecordModel>();
	let fileToken = $state('');

	async function GetQueueData() {
		loading = true;
		let id = page.url.searchParams.get('id') ?? '';
		try {
			fileToken = await pb.files.getToken();
			queueData = await pb.collection(PB_COLLECTION_IMAGE_QUEUES).getOne(id, {
				expand: 'generated_images_via_queue'
			});
		} catch {}
		await tick();
		refreshUpscaleButtonDisabled();
		loading = false;
	}

	onMount(() => {
		GetQueueData();
		GetDefaultSettings().then((v) => (defaultSettings = v));
	});

	let showConfirmUpscaleModal = $state(false);
	let disableUpscaleSelectedButton = $state(true);
	let upscaleTimes = $state(1);
	function refreshUpscaleButtonDisabled() {
		const selectedImages = queueData?.expand?.generated_images_via_queue?.filter(
			(x: RecordModel) => x.selected
		);
		disableUpscaleSelectedButton = (selectedImages?.length ?? 0) == 0;
	}

	async function HandleConfirmUpscale(e: SubmitEvent) {
		e.preventDefault();
		loading = true;

		const batch = pb.createBatch();
		queueData?.expand?.generated_images_via_queue?.forEach((element: RecordModel) => {
			if (element.selected) {
				batch.collection(PB_COLLECTION_GENERATED_IMAGES).update(element.id, {
					selected: true
				});
			} else {
				batch.collection(PB_COLLECTION_GENERATED_IMAGES).delete(element.id);
			}
		});
		batch.collection(PB_COLLECTION_IMAGE_QUEUES).update(queueData!.id, {
			user_confirmed_upscale: true,
			upscale_times: upscaleTimes
		});
		await batch.send();

		loading = false;
		location.reload();
	}
</script>

<svelte:head>
	<title>View Queue | Taprium</title>
</svelte:head>

<Section>
	{#if loading}
		<Spinner />
	{:else if queueData != undefined}
		<div class="flex items-center justify-between">
			<Heading>Prompts</Heading>
			<div>
				<Button onclick={GetQueueData} color="alternative"><RefreshOutline /></Button>
			</div>
		</div>
		<br />
		<Card class="bg-green-500 p-4 " size="xl">
			<p class="leading-tight font-normal text-white dark:text-gray-400">
				{queueData.positive_prompt}
			</p>
			<div>
				<Badge color="green" border>Positive</Badge>
			</div>
		</Card>
		<br />
		<Card class="bg-red-500 p-4 " size="xl">
			<p class="leading-tight font-normal text-white dark:text-gray-400">
				{queueData.negative_prompt}
			</p>
			<div>
				<Badge color="red" border>Positive</Badge>
			</div>
		</Card>

		<br />
		Status: {queueData.status}
		<br />
		Upscale Confirmed: {queueData.user_confirmed_upscale}
		<br />
		Size: {queueData.width} W * {queueData.height} H
		<br />
		Queued Image Number: {queueData.number}
		<br />
		Upscale Confirmed: {queueData.user_confirmed_upscale}
		{#if queueData.user_confirmed_upscale}
			<br />
			Upscale Times: {queueData.upscale_times}
		{/if}

		<Hr />

		<div class="flex items-center justify-between">
			<Heading
				>Images [ {queueData.expand?.generated_images_via_queue?.length ?? 0}/{queueData.number} ]</Heading
			>
			<Button
				color="cyan"
				disabled={disableUpscaleSelectedButton}
				onclick={() => {
					upscaleTimes = defaultSettings?.upscale_times ?? 2;
					showConfirmUpscaleModal = true;
				}}
			>
				Upscale Selected
			</Button>
		</div>
		<br />
		<div class="grid grid-cols-1 gap-2 sm:grid-cols-2 md:grid-cols-4">
			{#each queueData.expand?.generated_images_via_queue as i}
				<Card class="p-2">
					<Img
						id={i.id}
						src={`${PUBLIC_PB_ADDR.replace(/\/+$/, '')}/api/files/${PB_COLLECTION_GENERATED_IMAGES}/${i.id}/${i.image}?token=${fileToken}`}
					/>

					{#if i.upscaled}
						<Badge>Upscaled</Badge>
					{/if}
					{#if !queueData.user_confirmed_upscale}
						<div class=" m-auto mt-auto pt-2">
							<Toggle
								bind:checked={i.selected}
								onchange={(v) => {
									refreshUpscaleButtonDisabled();
								}}
							>
								Upscale This
							</Toggle>
						</div>
					{/if}
				</Card>
			{/each}
		</div>
	{:else}
		<Heading>Queue Not Found</Heading>
	{/if}
</Section>

<Modal title="Confirm Upscale?" bind:open={showConfirmUpscaleModal}>
	<p>Unselected images will be deleted, the deletion was not recoverable.</p>
	<form class="flex flex-col space-y-6" onsubmit={HandleConfirmUpscale}>
		<div>
			<Label class="space-y-2">
				<span>Default Upscale Times</span>
				<Input type="number" min={1} bind:value={upscaleTimes} required />
			</Label>
		</div>
		<div>
			<Button type="submit">Yes</Button>
			<Button value="decline" color="alternative" onclick={() => (showConfirmUpscaleModal = false)}>
				Cancel
			</Button>
		</div>
	</form>
</Modal>
