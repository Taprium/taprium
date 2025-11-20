<script lang="ts">
	import { Section } from 'flowbite-svelte-blocks';
	import { Card, Hr } from 'flowbite-svelte';
</script>

<Section>
	<div class="grid grid-cols-1 gap-2 md:grid-cols-4">
		<Card href="/img_gen" class="p-4 sm:p-6 md:p-8">
			<h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
				Image Generation
			</h5>
			<p class="leading-tight font-normal text-gray-700 dark:text-gray-400">Total: 200</p>
			<p class="leading-tight font-normal text-gray-700 dark:text-gray-400">In Queue: 4</p>
			<p class="leading-tight font-normal text-gray-700 dark:text-gray-400">
				To Generate Images: 40
			</p>
		</Card>
		<Card class="p-4 sm:p-6 md:p-8">
			<h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">Upscale</h5>
			<p class="leading-tight font-normal text-gray-700 dark:text-gray-400">In Queue: 4</p>
			<p class="leading-tight font-normal text-gray-700 dark:text-gray-400">Processing: 4</p>
		</Card>
	</div>
	<Hr />
	<Card href="/upscale_runners" class="p-4 sm:p-6 md:p-8">
		<h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
			Upscale Runners
		</h5>
		<p class="leading-tight font-normal text-gray-700 dark:text-gray-400">Total: 10</p>
		<p class="leading-tight font-normal text-gray-700 dark:text-gray-400">Online: 5</p>
	</Card>
</Section>
