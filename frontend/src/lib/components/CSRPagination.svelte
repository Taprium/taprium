<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Button, ButtonGroup, Label, PaginationNav } from 'flowbite-svelte';
	import { RefreshOutline } from 'flowbite-svelte-icons';

	export let data: any;

	export let loadListCallBack: () => void;
</script>

<div class="my-2">
	<div class="mt-5 flex w-full items-center justify-between">
		<div class="flex items-center justify-center gap-2">
			<Button
				size="xs"
				color="alternative"
				onclick={() => {
					loadListCallBack();
				}}
			>
				<RefreshOutline />
			</Button>
			<Label>
				Items: {data.totalItems}, Pages: {data.totalPages}, PerPage: {data.perPage}
			</Label>
			<ButtonGroup size="sm">
				<Button
					disabled={data.page == 0 || data.page == 1}
					size="xs"
					onclick={() => {
						page.url.searchParams.delete('page');
						goto(page.url);
						loadListCallBack();
					}}
				>
					First
				</Button>
				<Button
					disabled={data.page == data.totalPages}
					size="xs"
					onclick={() => {
						page.url.searchParams.set('page', data.totalPages);
						goto(page.url);
						loadListCallBack();
					}}
				>
					Last
				</Button>
			</ButtonGroup>
		</div>
		<PaginationNav
			currentPage={data.page}
			totalPages={data.totalPages}
			onPageChange={(v: number) => {
				page.url.searchParams.set('page', v.toString());
				goto(page.url);
				loadListCallBack();
			}}
		/>
	</div>
</div>
