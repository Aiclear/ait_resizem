<!-- Copyright (c) 2024 Barat Semet (https://github.com/barats)
Resizem is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details. -->

<script>
	import { Modal } from 'flowbite-svelte';
	import { selectedFile, showImageDetail, imageDetailData } from '$lib/app_stores';
	import { _ } from 'svelte-i18n';
	import { FileImageOutline, CheckCircleOutline, InfoCircleSolid } from 'flowbite-svelte-icons';

	let loading = false;
	let error = null;

	async function loadImageDetail() {
		if (!$selectedFile) return;

		loading = true;
		error = null;

		try {
			const getImageDetail = window['go']?.['rmanager']?.['FileManager']?.['GetImageDetail'];
			if (getImageDetail) {
				const detail = await getImageDetail($selectedFile);
				$imageDetailData = detail;
			} else {
				error = $_('image_detail.load_error');
			}
		} catch (err) {
			console.error('Failed to load image detail:', err);
			error = err.message || $_('image_detail.load_error');
		} finally {
			loading = false;
		}
	}

	function closeModal() {
		$showImageDetail = false;
		$selectedFile = null;
		$imageDetailData = null;
		error = null;
	}

	$: if ($showImageDetail && $selectedFile) {
		loadImageDetail();
	}

	$: if (!$showImageDetail) {
		closeModal();
	}
</script>

{#if $showImageDetail}
	<Modal bind:open={$showImageDetail} size="lg" on:close={closeModal}>
		<div class="flex items-center gap-2 border-b border-gray-200 pb-4">
			<InfoCircleSolid class="h-6 w-6 text-blue-500" />
			<h3 class="text-lg font-semibold text-gray-900">
				{$_('image_detail.title')}
			</h3>
		</div>

		<div class="max-h-[60vh] overflow-y-auto py-4">
			{#if loading}
				<div class="flex items-center justify-center py-8">
					<div class="h-8 w-8 animate-spin rounded-full border-4 border-blue-500 border-t-transparent" />
				</div>
			{:else if error}
				<div class="rounded-lg bg-red-50 p-4 text-center text-red-600">
					{error}
				</div>
			{:else if $imageDetailData}
				<div class="space-y-6">
					<div class="flex items-center gap-3 rounded-lg bg-gray-50 p-4">
						<FileImageOutline class="h-10 w-10 text-blue-500" />
						<div class="min-w-0 flex-1">
							<p class="font-medium text-gray-900 truncate" title={$imageDetailData.file_name}>
								{$imageDetailData.file_name}
							</p>
							<p class="text-sm text-gray-500 truncate" title={$imageDetailData.directory}>
								📁 {$imageDetailData.directory}
							</p>
						</div>
					</div>

					<div class="grid grid-cols-2 gap-4">
						<div class="rounded-lg border border-gray-200 p-4">
							<div class="flex items-center gap-2 text-gray-500">
								<span class="text-sm">📐 {$_('image_detail.dimensions')}</span>
							</div>
							<p class="mt-1 text-lg font-semibold text-gray-900">
								{$imageDetailData.width} × {$imageDetailData.height} px
							</p>
						</div>

						<div class="rounded-lg border border-gray-200 p-4">
							<div class="flex items-center gap-2 text-gray-500">
								<span class="text-sm">📄 {$_('image_detail.file_size')}</span>
							</div>
							<p class="mt-1 text-lg font-semibold text-gray-900">
								{$imageDetailData.file_size_str}
							</p>
						</div>

						<div class="rounded-lg border border-gray-200 p-4">
							<div class="flex items-center gap-2 text-gray-500">
								<span class="text-sm">🖼️ {$_('image_detail.format')}</span>
							</div>
							<p class="mt-1 text-lg font-semibold text-gray-900 uppercase">
								{$imageDetailData.format || '-'}
							</p>
						</div>

						<div class="rounded-lg border border-gray-200 p-4">
							<div class="flex items-center gap-2 text-gray-500">
								<span class="text-sm">🎨 {$_('image_detail.color_mode')}</span>
							</div>
							<p class="mt-1 text-lg font-semibold text-gray-900">
								{$imageDetailData.color_mode || '-'}
							</p>
						</div>
					</div>

					<div class="space-y-3 rounded-lg border border-gray-200 p-4">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2 text-gray-500">
								<span class="text-sm">📅 {$_('image_detail.created_time')}</span>
							</div>
							<span class="font-medium text-gray-900">
								{$imageDetailData.created_time || '-'}
							</span>
						</div>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2 text-gray-500">
								<span class="text-sm">📅 {$_('image_detail.modified_time')}</span>
							</div>
							<span class="font-medium text-gray-900">
								{$imageDetailData.modified_time || '-'}
							</span>
						</div>
					</div>

					<div class="rounded-lg border border-gray-200 p-4">
						<p class="text-sm text-gray-500">{$_('image_detail.full_path')}</p>
						<p class="mt-1 break-all text-sm font-mono text-gray-700">
							{$imageDetailData.file_path}
						</p>
					</div>
				</div>
			{/if}
		</div>

		<div class="flex justify-end border-t border-gray-200 pt-4">
			<button
				class="rounded-lg bg-blue-500 px-6 py-2 text-sm font-medium text-white hover:bg-blue-600 focus:outline-none focus:ring-4 focus:ring-blue-300"
				on:click={closeModal}
			>
				{$_('image_detail.close')}
			</button>
		</div>
	</Modal>
{/if}
