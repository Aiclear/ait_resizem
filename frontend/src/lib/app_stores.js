// Copyright (c) 2024 Barat Semet (https://github.com/barats)
// Resizem is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

import { writable } from '@macfja/svelte-persistent-store';
import {
	KEY_ASK_WHERE_TO_SAVE,
	KEY_CPU_VALUE,
	KEY_DOING,
	KEY_EXIF_ORIENTATION_VALUE,
	KEY_FILES_LIST,
	KEY_FILTER_VALUE,
	KEY_FORMAT_VALUE,
	KEY_GIF_COLORS_VALUE,
	KEY_HEIGHT_VALUE,
	KEY_JPEG_QUALITY_VALUE,
	KEY_PNG_COMPRESSION_VALUE,
	KEY_RESULT_LIST,
	KEY_TIFF_COMPRESSION_VALUE,
	KEY_WIDTH_VALUE,
	KEY_WATERMARK_TYPE,
	KEY_WATERMARK_TEXT,
	KEY_WATERMARK_IMAGE_PATH,
	KEY_WATERMARK_OPACITY,
	KEY_WATERMARK_POSITION,
	KEY_WATERMARK_OFFSET_X,
	KEY_WATERMARK_OFFSET_Y,
	KEY_WATERMARK_ROTATION,
	KEY_WATERMARK_SCALE,
	KEY_WATERMARK_FONT_SIZE,
	KEY_WATERMARK_FONT_COLOR,
	KEY_RENAME_ENABLED,
	KEY_RENAME_TEMPLATE,
	KEY_RENAME_START_INDEX,
	KEY_RENAME_INDEX_PADDING,
	KEY_RENAME_DATE_FORMAT
} from './app_consts';

export let doing = writable(KEY_DOING, false);
export let askWheretoSave = writable(KEY_ASK_WHERE_TO_SAVE, false);

export let filesList = writable(KEY_FILES_LIST, []);
export let resultList = writable(KEY_RESULT_LIST, []);

export let formatValue = writable(KEY_FORMAT_VALUE, 0); //Keep original format
export let filterValue = writable(KEY_FILTER_VALUE, 1); // Lanczos
export let widthValue = writable(KEY_WIDTH_VALUE);
export let heightValue = writable(KEY_HEIGHT_VALUE);

export let cpuUsageValue = writable(KEY_CPU_VALUE, 1);
export let jpegQualityValue = writable(KEY_JPEG_QUALITY_VALUE, 75);
export let gifColorsValue = writable(KEY_GIF_COLORS_VALUE, 256);
export let tiffCompressionValue = writable(KEY_TIFF_COMPRESSION_VALUE, 0);
export let pngCompressionValue = writable(KEY_PNG_COMPRESSION_VALUE, -1);
export let autoExifOrientation = writable(KEY_EXIF_ORIENTATION_VALUE, false);

export let watermarkType = writable(KEY_WATERMARK_TYPE, 0);
export let watermarkText = writable(KEY_WATERMARK_TEXT, '');
export let watermarkImagePath = writable(KEY_WATERMARK_IMAGE_PATH, '');
export let watermarkOpacity = writable(KEY_WATERMARK_OPACITY, 0.5);
export let watermarkPosition = writable(KEY_WATERMARK_POSITION, 8);
export let watermarkOffsetX = writable(KEY_WATERMARK_OFFSET_X, 10);
export let watermarkOffsetY = writable(KEY_WATERMARK_OFFSET_Y, 10);
export let watermarkRotation = writable(KEY_WATERMARK_ROTATION, 0);
export let watermarkScale = writable(KEY_WATERMARK_SCALE, 1.0);
export let watermarkFontSize = writable(KEY_WATERMARK_FONT_SIZE, 24);
export let watermarkFontColor = writable(KEY_WATERMARK_FONT_COLOR, '#FFFFFF');

export let renameEnabled = writable(KEY_RENAME_ENABLED, false);
export let renameTemplate = writable(KEY_RENAME_TEMPLATE, '{name}_{index:3}');
export let renameStartIndex = writable(KEY_RENAME_START_INDEX, 1);
export let renameIndexPadding = writable(KEY_RENAME_INDEX_PADDING, 3);
export let renameDateFormat = writable(KEY_RENAME_DATE_FORMAT, 'YYYYMMDD');
