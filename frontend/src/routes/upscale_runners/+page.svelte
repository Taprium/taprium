<script lang="ts">
	import { pb, PB_COLLECTION_UPSCALE_RUNNERS } from '$lib/pb/backend-pb';
	import {
		Badge,
		Button,
		Heading,
		Input,
		Label,
		Modal,
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Toggle
	} from 'flowbite-svelte';
	import { Section } from 'flowbite-svelte-blocks';
	import { EyeOutline, EyeSlashOutline, RefreshOutline } from 'flowbite-svelte-icons';
	import type { RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';

	let loading = $state(false);
	let runners = $state<Array<RecordModel>>();

	async function GetRunners() {
		loading = true;

		runners = await pb.collection(PB_COLLECTION_UPSCALE_RUNNERS).getFullList();

		loading = false;
	}

	onMount(() => {
		GetRunners();
	});

	async function handleEnabled(id: string, value: boolean) {
		await pb.collection(PB_COLLECTION_UPSCALE_RUNNERS).update(id, {
			enabled: value
		});
		GetRunners();
	}

	async function handleDelete(id: string) {
		await pb.collection(PB_COLLECTION_UPSCALE_RUNNERS).delete(id);
		GetRunners();
	}

	let showChangePasswordModal = $state(false);
	let changePasswordRunner = $state<RecordModel>();
	let changePasswordError = $state('');

	let changePasswordValue = $state('');
	let changePasswordConfirmValue = $state('');
	let changePasswordShowValue = $state(false);
	async function handleChangePassword(e: SubmitEvent) {
		e.preventDefault();
		if (changePasswordValue != changePasswordConfirmValue) {
			changePasswordError = 'Password value does not match.';
			return;
		} else if (changePasswordValue == '') {
			changePasswordError = 'Password can not be empty.';
			return;
		} else {
			changePasswordError = '';
		}

		await pb.collection(PB_COLLECTION_UPSCALE_RUNNERS).update(changePasswordRunner!.id, {
			password: changePasswordValue,
			passwordConfirm: changePasswordConfirmValue
		});
		showChangePasswordModal = false;
	}

	let showAddRunnerModal = $state(false);
	let addRunnerNameValue = $state('');
	let addRunnerPasswordValue = $state('');
	let addRunnerPasswordConfirmValue = $state('');
	let addRunnerEnabled = $state(false);
	let addRunnerShowPassword = $state(false);
	let addRunnerError = $state('');
	async function handleAddRummer(e: SubmitEvent) {
		e.preventDefault();
		await pb.collection(PB_COLLECTION_UPSCALE_RUNNERS).create({
			verified: true,
			name: addRunnerNameValue,
			enabled: addRunnerEnabled,
			password: addRunnerPasswordValue,
			passwordConfirm: addRunnerPasswordConfirmValue
		});
		showAddRunnerModal = false;
		GetRunners();
	}
</script>

<svelte:head>
	<title>Upscale Runners - Taprium</title>
</svelte:head>

<Section>
	<div class="flex items-center justify-between">
		<Heading class="mb-5 text-green-600">Upscale Runners</Heading>
		<div class="flex items-end gap-2">
			<Button
				color="alternative"
				onclick={() => {
					GetRunners();
				}}
			>
				<RefreshOutline />
			</Button>
			<Button
				onclick={() => {
					addRunnerNameValue = '';
					addRunnerPasswordValue = '';
					addRunnerPasswordConfirmValue = '';
					addRunnerError = '';
					addRunnerShowPassword = false;
					showAddRunnerModal = true;
				}}
			>
				Add
			</Button>
		</div>
	</div>

	<Table>
		<TableHead>
			<TableHeadCell>Name</TableHeadCell>
			<TableHeadCell>Enabled</TableHeadCell>
			<TableHeadCell>Last Seen</TableHeadCell>
			<TableHeadCell>Operation</TableHeadCell>
		</TableHead>
		<TableBody>
			{#each runners as r}
				<TableBodyRow>
					<TableBodyCell>{r.name}</TableBodyCell>
					<TableBodyCell>
						<Badge color={r.enabled ? 'green' : 'red'}>{r.enabled}</Badge>
					</TableBodyCell>
					<TableBodyCell>{new Date(r.pinged_at)}</TableBodyCell>
					<TableBodyCell>
						{#if r.enabled}
							<Button
								onclick={() => {
									handleEnabled(r.id, false);
								}}
							>
								Disable
							</Button>
						{:else}
							<Button
								onclick={() => {
									handleEnabled(r.id, true);
								}}
							>
								Enable
							</Button>
						{/if}

						<Button
							onclick={() => {
								changePasswordConfirmValue = '';
								changePasswordValue = '';
								changePasswordShowValue = false;
								changePasswordRunner = r;
								showChangePasswordModal = true;
							}}>Change Password</Button
						>

						<Button
							color="alternative"
							onclick={() => {
								handleDelete(r.id);
							}}
						>
							Delete
						</Button>
					</TableBodyCell>
				</TableBodyRow>
			{/each}
		</TableBody>
	</Table>
</Section>

<Modal
	title="Change password for [{changePasswordRunner?.name}]"
	bind:open={showChangePasswordModal}
>
	<form class="flex flex-col space-y-6" onsubmit={handleChangePassword}>
		{#if changePasswordError}
			<Label color="red">{changePasswordError}</Label>
		{/if}
		<Label class="space-y-2">
			<span>Password - at least 8 characters</span>
			<Input
				type={changePasswordShowValue ? 'text' : 'password'}
				name="password"
				minlength={8}
				required
				bind:value={changePasswordValue}
				onchange={() => {
					if (changePasswordValue != changePasswordConfirmValue) {
						changePasswordError = 'Password value does not match.';
					} else if (changePasswordValue.length < 8) {
						changePasswordError = 'Password must be at least 8 characters.';
					} else {
						changePasswordError = '';
					}
				}}
			>
				{#snippet right()}
					<button
						class="pointer-events-auto"
						onclick={() => {
							changePasswordShowValue = !changePasswordShowValue;
						}}
					>
						{#if changePasswordShowValue}
							<EyeSlashOutline class="h-6 w-6" />
						{:else}
							<EyeOutline class="h-6 w-6" />
						{/if}
					</button>
				{/snippet}</Input
			>
		</Label>
		<Label class="space-y-2">
			<span>Confirm Password - at least 8 characters</span>
			<Input
				type={changePasswordShowValue ? 'text' : 'password'}
				name="password"
				minlength={8}
				required
				bind:value={changePasswordConfirmValue}
				onchange={() => {
					if (changePasswordValue != changePasswordConfirmValue) {
						changePasswordError = 'Password value does not match.';
					} else {
						changePasswordError = '';
					}
				}}
			/>
		</Label>
		<div>
			<Button type="submit">Submit</Button>
			<Button color="alternative" onclick={() => (showChangePasswordModal = false)}>Cancel</Button>
		</div>
	</form>
</Modal>

<Modal title="Add upscale runners" bind:open={showAddRunnerModal}>
	<form class="flex flex-col space-y-6" onsubmit={handleAddRummer}>
		{#if addRunnerError}
			<Label color="red">{addRunnerError}</Label>
		{/if}

		<Label class="space-y-2">
			<span>Runner Name</span>
			<Input type="text" required bind:value={addRunnerNameValue} />
		</Label>
		<Label class="space-y-2">
			<span>Password - at least 8 characters</span>
			<Input
				type={addRunnerShowPassword ? 'text' : 'password'}
				name="password"
				minlength={8}
				required
				bind:value={addRunnerPasswordValue}
				onchange={() => {
					if (addRunnerPasswordValue != addRunnerPasswordConfirmValue) {
						addRunnerError = 'Password value does not match.';
					} else if (changePasswordValue.length < 8) {
						addRunnerError = 'Password must be at least 8 characters.';
					} else {
						addRunnerError = '';
					}
				}}
			>
				{#snippet right()}
					<button
						class="pointer-events-auto"
						onclick={() => {
							addRunnerShowPassword = !addRunnerShowPassword;
						}}
					>
						{#if addRunnerShowPassword}
							<EyeSlashOutline class="h-6 w-6" />
						{:else}
							<EyeOutline class="h-6 w-6" />
						{/if}
					</button>
				{/snippet}</Input
			>
		</Label>
		<Label class="space-y-2">
			<span>Confirm Password - at least 8 characters</span>
			<Input
				type={addRunnerShowPassword ? 'text' : 'password'}
				name="password"
				minlength={8}
				required
				bind:value={addRunnerPasswordConfirmValue}
				onchange={() => {
					if (addRunnerPasswordValue != addRunnerPasswordConfirmValue) {
						addRunnerError = 'Password value does not match.';
					} else if (addRunnerPasswordValue.length < 8) {
						addRunnerError = 'Password must be at least 8 characters.';
					} else {
						addRunnerError = '';
					}
				}}
			/>
		</Label>
		<Toggle bind:checked={addRunnerEnabled}>Enabled</Toggle>
		<div>
			<Button type="submit">Submit</Button>
			<Button color="alternative" onclick={() => (showAddRunnerModal = false)}>Cancel</Button>
		</div>
	</form>
</Modal>
