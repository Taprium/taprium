<script lang="ts">
	import { page } from '$app/state';
	import { PUBLIC_PB_ADDR } from '$env/static/public';
	import {
		pb,
		PB_COLLECTION_GENERATE_QUEUES,
		PB_COLLECTION_GENERATED_IMAGES
	} from '$lib/pb/backend-pb';
	import { Badge, Button, Card, Heading, Hr, Img, Modal, Spinner, Toggle } from 'flowbite-svelte';
	import { Section } from 'flowbite-svelte-blocks';
	import { RefreshOutline } from 'flowbite-svelte-icons';
	import { type RecordModel } from 'pocketbase';
	import { onMount, tick } from 'svelte';

	let loading = $state(false);
	let queueData = $state<RecordModel>();

	async function GetQueueData() {
		loading = true;
		let id = page.url.searchParams.get('id') ?? '';
		try {
			queueData = await pb.collection(PB_COLLECTION_GENERATE_QUEUES).getOne(id, {
				expand: 'generated_images_via_queue'
			});
		} catch {}
		await tick();
		loading = false;
	}

	onMount(() => {
		GetQueueData();
	});

	let showConfirmUpscaleModal = $state(false);
	let disableUpscaleSelectedButton = $state(true);
	function refreshUpscaleButtonDisabled() {
		const selectedImages = queueData?.expand?.generated_images_via_queue?.filter(
			(x: RecordModel) => x.selected
		);
		disableUpscaleSelectedButton = selectedImages.length == 0;
	}

	async function HandleConfirmUpscale() {
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
		batch.collection(PB_COLLECTION_GENERATE_QUEUES).update(queueData!.id, {
			user_confirmed_upscale: true
		});
		await batch.send();

		loading = false;
		location.reload();
	}
</script>

<svelte:head>
	<title>View Queue | AI Shared</title>
</svelte:head>

<Section>
	{#if loading}
		<Spinner />
	{:else if queueData != undefined}
		<div class="flex items-center justify-between">
			<Heading>Prompts</Heading>
			<Button onclick={GetQueueData}><RefreshOutline /></Button>
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
		Upscale confirmed: {queueData.user_confirmed_upscale}

		<Hr />

		<div class="flex items-center justify-between">
			<Heading>Images</Heading>
			<Button
				color="cyan"
				disabled={disableUpscaleSelectedButton}
				onclick={() => {
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
						src={`${PUBLIC_PB_ADDR}/api/files/${PB_COLLECTION_GENERATED_IMAGES}/${i.id}/${i.image}`}
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
	{#snippet footer()}
		<Button onclick={HandleConfirmUpscale}>Yes</Button>
		<Button value="decline" color="alternative" onclick={() => (showConfirmUpscaleModal = false)}
			>Cancel</Button
		>
	{/snippet}
</Modal>
