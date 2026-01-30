<script lang="ts">
	import { Section } from 'flowbite-svelte-blocks';
	import { Card, Hr } from 'flowbite-svelte';
	import { onDestroy, onMount } from 'svelte';
	import {
		pb,
		PB_COLLECTION_IMAGE_QUEUES,
		PB_COLLECTION_GENERATED_IMAGES,
		PB_COLLECTION_UPSCALE_RUNNERS
	} from '$lib/pb/backend-pb';
	import type { RecordModel } from 'pocketbase';

	let toGenerateImages = $state<RecordModel[]>();
	let upscaleData = $state<RecordModel[]>();
	let upscaleRunnerData = $state<RecordModel[]>();

	async function GetImageGenerationData() {
		toGenerateImages = await pb.collection(PB_COLLECTION_IMAGE_QUEUES).getFullList({
			filter: 'status="queue" || status="processing"'
		});
	}
	async function GetUpscaleData() {
		upscaleData = await pb.collection(PB_COLLECTION_GENERATED_IMAGES).getFullList({
			filter: 'selected=true && upscaled=false'
		});
	}

	let upscaleRunnerValidTime = new Date();
	async function GetUpscaleRunnerData() {
		upscaleRunnerData = await pb.collection(PB_COLLECTION_UPSCALE_RUNNERS).getFullList({});
		upscaleRunnerValidTime = new Date();
		upscaleRunnerValidTime.setMinutes(upscaleRunnerValidTime.getMinutes() - 1);
		// console.log(upscaleRunnerValidTime);
	}

	let interval: NodeJS.Timeout;
	async function getDashboardData() {
		interval = setInterval(async () => {
			await GetImageGenerationData();
			await GetUpscaleData();
			await GetUpscaleRunnerData();
		}, 500);
	}

	onDestroy(() => {
		clearInterval(interval);
	});

	onMount(() => {
		getDashboardData();
	});
</script>

<svelte:head>
	<title>Taprium</title>
</svelte:head>

<Section>
	<div class="grid grid-cols-1 gap-2 md:grid-cols-4">
		<Card href="/img_gen" class="p-4 sm:p-6 md:p-8">
			<h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
				Image Generation
			</h5>
			<p class="leading-tight font-normal text-gray-700 dark:text-gray-400">
				To Generate Images: {toGenerateImages?.reduce((acc, x) => acc + x.number, 0) ?? '0'}
			</p>
		</Card>
		<Card class="p-4 sm:p-6 md:p-8">
			<h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">Upscale</h5>
			<p class="leading-tight font-normal text-gray-700 dark:text-gray-400">
				To Upscale: {upscaleData?.length ?? '0'}
			</p>
		</Card>
	</div>
	<Hr />
	<Card href="/upscale_runners" class="p-4 sm:p-6 md:p-8">
		<h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
			Upscale Runners
		</h5>
		<p class="leading-tight font-normal text-gray-700 dark:text-gray-400">
			Total: {upscaleRunnerData?.length ?? '0'}
		</p>
		<p class="leading-tight font-normal text-gray-700 dark:text-gray-400">
			Online: {upscaleRunnerData?.filter((x) => {
				console.log(new Date(x.pinged_at));
				return x.pinged_at != '' && new Date(x.pinged_at) > upscaleRunnerValidTime;
			}).length ?? '0'}
		</p>
	</Card>
</Section>
