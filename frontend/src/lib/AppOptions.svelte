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
	import { NumberInput, Label, Select, Helper } from 'flowbite-svelte';
	import { Range, Toggle } from 'flowbite-svelte';
	import { OutputImagesTypes, ResampleFilterTypes } from '$lib/wailsjs/go/rmanager/TypeManager.js';
	import { onMount } from 'svelte';
	import { _ } from 'svelte-i18n';
	import {
		filterValue,
		formatValue,
		heightValue,
		widthValue,
		watermarkType,
		watermarkText,
		watermarkImagePath,
		watermarkOpacity,
		watermarkPosition,
		watermarkOffsetX,
		watermarkOffsetY,
		watermarkRotation,
		watermarkScale,
		watermarkFontSize,
		watermarkFontColor,
		renameEnabled,
		renameTemplate,
		renameStartIndex,
		renameIndexPadding,
		renameDateFormat
	} from './app_stores';

	let showWidthHelper = false,
		showHeightHelper = false;

	let allFormats, allFilters;
	let showWatermarkOptions = false;
	let showRenameOptions = false;

	const watermarkTypes = [
		{ name: '无水印', value: 0 },
		{ name: '文字水印', value: 1 },
		{ name: '图片水印', value: 2 }
	];

	const watermarkPositions = [
		{ name: '左上角', value: 0 },
		{ name: '顶部居中', value: 1 },
		{ name: '右上角', value: 2 },
		{ name: '左侧居中', value: 3 },
		{ name: '居中', value: 4 },
		{ name: '右侧居中', value: 5 },
		{ name: '左下角', value: 6 },
		{ name: '底部居中', value: 7 },
		{ name: '右下角', value: 8 }
	];

	onMount(() => {
		OutputImagesTypes().then((data) => {
			var keep = { name: $_('home.options.keep'), value: 0 };
			data.push(keep);
			allFormats = data;
			console.log('remote call OutputImagesTypes() done');
		});

		ResampleFilterTypes().then((data) => {
			allFilters = data;
			console.log('remote call ResampleFilterTypes() done');
		});
	});

	$: if ($widthValue < 0 || $widthValue === null) {
		showWidthHelper = true;
	} else {
		showWidthHelper = false;
	}

	$: if ($heightValue < 0 || $heightValue === null) {
		showHeightHelper = true;
	} else {
		showHeightHelper = false;
	}

	$: showWatermarkOptions = $watermarkType !== 0;
</script>

<div class="grid grid-cols-1 gap-4">
	<div class="border-b pb-4">
		<Label class="text-lg font-medium mb-2">{$_('home.options.tab.resize')}</Label>
		<div class="grid grid-cols-1 gap-3">
			<div>
				<Label for="format-select">{$_('home.options.format')}</Label>
				<Select
					id="format-select"
					placeholder={$_('home.options.choose')}
					items={allFormats}
					bind:value={$formatValue}
				/>
			</div>
			<div>
				<Label for="filter-select">{$_('home.options.filter')}</Label>
				<Select id="filter-select" items={allFilters} bind:value={$filterValue} />
			</div>
			<div class="grid grid-cols-2 gap-5">
				<Label
					>{$_('home.options.width.title')}
					<NumberInput id="width" bind:value={$widthValue} />
					{#if showWidthHelper}
						<Helper class="mt-2" color="red">
							<span class="font-medium">{$_('home.options.width.helper1')}</span>
							{$_('home.options.width.helper2')}
						</Helper>
					{/if}
				</Label>
				<Label
					>{$_('home.options.height.title')}
					<NumberInput id="height" bind:value={$heightValue} />
					{#if showHeightHelper}
						<Helper class="mt-2" color="red">
							<span class="font-medium">{$_('home.options.height.helper1')}</span>
							{$_('home.options.height.helper2')}
						</Helper>
					{/if}
				</Label>
			</div>
		</div>
	</div>

	<div class="border-b pb-4">
		<div class="flex items-center justify-between mb-2">
			<Label class="text-lg font-medium">{$_('home.options.tab.watermark')}</Label>
		</div>
		<div class="grid grid-cols-1 gap-3">
			<div>
				<Label for="watermark-type">{$_('home.options.watermark.type')}</Label>
				<Select id="watermark-type" items={watermarkTypes} bind:value={$watermarkType} />
			</div>

			{#if showWatermarkOptions}
				{#if $watermarkType === 1}
					<div>
						<Label for="watermark-text">{$_('home.options.watermark.text')}</Label>
						<textarea
							id="watermark-text"
							class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							placeholder={$_('home.options.watermark.text_placeholder')}
							bind:value={$watermarkText}
							rows={2}
						/>
					</div>
					<div class="grid grid-cols-2 gap-3">
						<div>
							<Label for="watermark-font-size">{$_('home.options.watermark.font_size')}</Label>
							<NumberInput id="watermark-font-size" bind:value={$watermarkFontSize} />
						</div>
						<div>
							<Label for="watermark-font-color">{$_('home.options.watermark.font_color')}</Label>
							<input
								type="text"
								id="watermark-font-color"
								class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
								bind:value={$watermarkFontColor}
								placeholder="#FFFFFF"
							/>
						</div>
					</div>
				{/if}

				{#if $watermarkType === 2}
					<div>
						<Label for="watermark-image">{$_('home.options.watermark.image')}</Label>
						<input
							type="text"
							id="watermark-image"
							class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							placeholder={$_('home.options.watermark.image_placeholder')}
							bind:value={$watermarkImagePath}
						/>
					</div>
				{/if}

				<div>
					<Label for="watermark-position">{$_('home.options.watermark.position')}</Label>
					<Select id="watermark-position" items={watermarkPositions} bind:value={$watermarkPosition} />
				</div>
				<div class="grid grid-cols-2 gap-3">
					<div>
						<Label>{$_('home.options.watermark.opacity')} ({$watermarkOpacity})</Label>
						<Range min="0" max="1" step="0.1" bind:value={$watermarkOpacity} />
					</div>
					<div>
						<Label>{$_('home.options.watermark.scale')} ({$watermarkScale})</Label>
						<Range min="0.1" max="3" step="0.1" bind:value={$watermarkScale} />
					</div>
				</div>
				<div class="grid grid-cols-2 gap-3">
					<div>
						<Label for="watermark-offset-x">{$_('home.options.watermark.offset_x')}</Label>
						<NumberInput id="watermark-offset-x" bind:value={$watermarkOffsetX} />
					</div>
					<div>
						<Label for="watermark-offset-y">{$_('home.options.watermark.offset_y')}</Label>
						<NumberInput id="watermark-offset-y" bind:value={$watermarkOffsetY} />
					</div>
				</div>
				<div>
					<Label>{$_('home.options.watermark.rotation')} ({$watermarkRotation}°)</Label>
					<Range min="0" max="360" step="1" bind:value={$watermarkRotation} />
				</div>
			{/if}
		</div>
	</div>

	<div>
		<div class="flex items-center justify-between mb-2">
			<Label class="text-lg font-medium">{$_('home.options.tab.rename')}</Label>
			<Toggle bind:checked={$renameEnabled} />
		</div>

		{#if $renameEnabled}
			<div class="grid grid-cols-1 gap-3 pl-4">
				<div>
					<Label for="rename-template">{$_('home.options.rename.template')}</Label>
					<input
						type="text"
						id="rename-template"
						class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						bind:value={$renameTemplate}
					/>
					<Helper class="mt-1">
						<span class="text-xs">
							{$_('home.options.rename.template_help')}
						</span>
					</Helper>
				</div>
				<div class="grid grid-cols-2 gap-3">
					<div>
						<Label for="rename-start-index">{$_('home.options.rename.start_index')}</Label>
						<NumberInput id="rename-start-index" bind:value={$renameStartIndex} />
					</div>
					<div>
						<Label for="rename-index-padding">{$_('home.options.rename.index_padding')}</Label>
						<NumberInput id="rename-index-padding" bind:value={$renameIndexPadding} />
					</div>
				</div>
				<div>
					<Label for="rename-date-format">{$_('home.options.rename.date_format')}</Label>
					<input
						type="text"
						id="rename-date-format"
						class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						bind:value={$renameDateFormat}
						placeholder="YYYYMMDD"
					/>
					<Helper class="mt-1">
						<span class="text-xs">
							{$_('home.options.rename.date_format_help')}
						</span>
					</Helper>
				</div>
			</div>
		{/if}
	</div>
</div>
