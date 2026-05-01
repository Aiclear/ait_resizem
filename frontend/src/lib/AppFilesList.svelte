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
	import { Listgroup, ListgroupItem } from 'flowbite-svelte';
	import { filesList, resultList, selectedFile, showImageDetail, imageDetailData, doing } from '$lib/app_stores';
	import { FileImageOutline, CheckCircleOutline, CloseCircleOutline, InfoCircleSolid } from 'flowbite-svelte-icons';
	import { _ } from 'svelte-i18n';

	async function handleFileClick(filePath) {
		if ($doing) return;
		$selectedFile = filePath;
		$showImageDetail = true;
	}

	function isSelected(filePath) {
		return $selectedFile === filePath;
	}

	function getFileName(path) {
		return path.split(/[/\\]/).pop() || path;
	}
</script>

<Listgroup class="w-full border-0">
	{#if $filesList.length > 0}
		{#each $filesList as item}
			<ListgroupItem
				class="flex cursor-pointer items-center justify-between gap-2 text-sm transition-colors {isSelected(item) ? 'bg-blue-50 hover:bg-blue-100' : 'hover:bg-gray-50'}"
				on:click={() => handleFileClick(item)}
			>
				<div class="flex min-w-0 flex-1 items-center gap-2">
					<FileImageOutline class="h-5 w-5 flex-shrink-0 {isSelected(item) ? 'text-blue-500' : 'text-gray-500'}" />
					<span class="truncate" title={item}>{getFileName(item)}</span>
				</div>
				<InfoCircleSolid
					class="h-4 w-4 flex-shrink-0 {isSelected(item) ? 'text-blue-500' : 'text-gray-400'}"
					title={$_('file_list.view_details')}
				/>
			</ListgroupItem>
		{/each}
	{/if}

	{#if $resultList.length > 0}
		{#each $resultList as item}
			<ListgroupItem class="flex gap-2 text-sm">
				{#if item.status === 1}
					<CheckCircleOutline class="h-5 w-5 text-green-500" />
					<span class="text-green-700">{item.name}</span>
				{:else}
					<CloseCircleOutline class="h-5 w-5 text-red-500" />
					<span class="text-red-700">{item.name} - {item.message}</span>
				{/if}
			</ListgroupItem>
		{/each}
	{/if}
</Listgroup>
