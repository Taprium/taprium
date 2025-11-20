<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import CSRPagination from '$lib/components/CSRPagination.svelte';
	import { pb, PB_COLLECTION_GENERATE_QUEUES } from '$lib/pb/backend-pb';
	import {
		Badge,
		Button,
		Heading,
		Input,
		Label,
		Modal,
		Spinner,
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell
	} from 'flowbite-svelte';
	import { Section } from 'flowbite-svelte-blocks';
	import type { ListResult, RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';

	let loading = $state(false);
	let showAddModal = $state(false);
	let queues = $state<ListResult<RecordModel>>();
	let pageIndex = $state(1);

	async function GetQueues() {
		loading = true;

		pageIndex = parseInt(page.url.searchParams.get('page') ?? '1');
		queues = await pb.collection(PB_COLLECTION_GENERATE_QUEUES).getList(pageIndex, 20, {
			filter: '',
			expand: 'generated_images_via_queue',
			sort: '-updated'
		});

		loading = false;
	}

	onMount(() => {
		GetQueues();
	});

	let addQueuePositive = $state('');
	let addQueueNegative = $state('');
	let addQueueNumber = $state(1);
	let addQueueLoading = $state(false);

	async function HandleAddQueue(e: SubmitEvent) {
		e.preventDefault();
		addQueueLoading = true;
		let record = await pb.collection(PB_COLLECTION_GENERATE_QUEUES).create({
			positive_prompt: addQueuePositive,
			negative_prompt: addQueueNegative,
			number: addQueueNumber,
			status: 'queue'
		});
		goto(`/img_gen/edit?id=${record.id}`);
	}
</script>

<svelte:head>
	<title>Image Generation Queues | AI Shared</title>
</svelte:head>

<Section>
	<div class="flex items-center justify-between">
		<Heading class="mb-5 text-purple-600">Image Generation Queues</Heading>

		<div>
			<Button
				onclick={() => {
					addQueueLoading = false;
					addQueuePositive = '';
					addQueueNegative = '';
					addQueueNumber = 1;
					showAddModal = true;
				}}>Add</Button
			>
		</div>
	</div>
	{#if loading}
		<Spinner />
	{:else if queues != undefined}
		<CSRPagination
			data={queues}
			loadListCallBack={() => {
				GetQueues();
			}}
		/>
		<Table>
			<TableHead>
				<TableHeadCell>prompt</TableHeadCell>
				<TableHeadCell>number</TableHeadCell>
				<TableHeadCell>status</TableHeadCell>
				<TableHeadCell>queue at</TableHeadCell>
				<TableHeadCell>operation</TableHeadCell>
			</TableHead>
			<TableBody>
				{#each queues.items as q}
					<TableBodyRow>
						<TableBodyCell>
							<div class="flex flex-col gap-2">
								<div class="text-wrap">
									<Badge color="green">{q.positive_prompt}</Badge>
								</div>
								<div class="text-wrap">
									<Badge color="red">{q.negative_prompt}</Badge>
								</div>
							</div>
						</TableBodyCell>
						<TableBodyCell>
							Generated: {q.expand?.generated_images_via_queue?.length ?? 0}
							<br />
							Requested: {q.number}
						</TableBodyCell>
						<TableBodyCell>{q.status}</TableBodyCell>
						<TableBodyCell>
							{q.created.split(' ')[0]}
						</TableBodyCell>
						<TableBodyCell>
							<Button
								onclick={() => {
									addQueueLoading = false;
									addQueuePositive = q.positive_prompt;
									addQueueNegative = q.negative_prompt;
									addQueueNumber = 1;
									showAddModal = true;
								}}>Queue More</Button
							>
							<Button href="/img_gen/view?id={q.id}" color="alternative">View</Button>
						</TableBodyCell>
					</TableBodyRow>
				{/each}
			</TableBody>
		</Table>
		<CSRPagination
			data={queues}
			loadListCallBack={() => {
				GetQueues();
			}}
		/>
	{/if}
</Section>

<Modal title="Add Image Generation Queue" bind:open={showAddModal}>
	{#if addQueueLoading}
		<Spinner />
	{:else}
		<form onsubmit={HandleAddQueue}>
			<div class="flex flex-col space-y-6">
				<Label class="space-y-2">
					<span>Positive Prompt</span>
					<Input type="text" required bind:value={addQueuePositive} />
				</Label>
				<Label class="space-y-2">
					<span>Negative Prompt</span>
					<Input type="text" bind:value={addQueueNegative} />
				</Label>
				<Label class="space-y-2">
					<span>Generate Image Count</span>
					<Input type="number" min={1} required bind:value={addQueueNumber} />
				</Label>
			</div>
			<Button type="submit">Submit</Button>
			<Button
				color="alternative"
				onclick={() => {
					showAddModal = false;
				}}
			>
				Cancel
			</Button>
		</form>
	{/if}
</Modal>
