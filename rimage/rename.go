// Copyright (c) 2024 Barat Semet (https://github.com/barats)
// Resizem is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package rimage

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RenameContext struct {
	Index       int
	Total       int
	OriginalName string
	OriginalExt  string
	OriginalPath string
	Width       int
	Height      int
}

func GenerateNewName(template string, ctx RenameContext, opts RenameOptions) string {
	if !opts.Enabled || template == "" {
		return ctx.OriginalName
	}

	result := template

	result = replaceIndexPlaceholder(result, ctx.Index, opts.StartIndex, opts.IndexPadding)
	result = replaceDatePlaceholder(result, opts.DateFormat)
	result = replaceNamePlaceholder(result, ctx.OriginalName)
	result = replaceExtPlaceholder(result, ctx.OriginalExt)
	result = replaceSizePlaceholder(result, ctx.Width, ctx.Height)
	result = replaceTotalPlaceholder(result, ctx.Total)

	return result
}

func replaceIndexPlaceholder(template string, index, startIndex, padding int) string {
	re := regexp.MustCompile(`\{index(:(\d+))?\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	for _, match := range matches {
		placeholder := match[0]
		actualIndex := index + startIndex
		paddingSize := padding
		if match[2] != "" {
			paddingSize, _ = strconv.Atoi(match[2])
		}
		formatted := fmt.Sprintf("%0*d", paddingSize, actualIndex)
		template = strings.Replace(template, placeholder, formatted, 1)
	}

	return template
}

func replaceDatePlaceholder(template string, dateFormat string) string {
	re := regexp.MustCompile(`\{date(:([^}]+))?\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	now := time.Now()

	for _, match := range matches {
		placeholder := match[0]
		format := dateFormat
		if match[2] != "" {
			format = match[2]
		}
		formatted := formatDate(now, format)
		template = strings.Replace(template, placeholder, formatted, 1)
	}

	return template
}

func formatDate(t time.Time, format string) string {
	formatMap := map[string]string{
		"YYYY": strconv.Itoa(t.Year()),
		"YY":   fmt.Sprintf("%02d", t.Year()%100),
		"MM":   fmt.Sprintf("%02d", t.Month()),
		"M":    strconv.Itoa(int(t.Month())),
		"DD":   fmt.Sprintf("%02d", t.Day()),
		"D":    strconv.Itoa(t.Day()),
		"HH":   fmt.Sprintf("%02d", t.Hour()),
		"H":    strconv.Itoa(t.Hour()),
		"mm":   fmt.Sprintf("%02d", t.Minute()),
		"m":    strconv.Itoa(t.Minute()),
		"ss":   fmt.Sprintf("%02d", t.Second()),
		"s":    strconv.Itoa(t.Second()),
	}

	result := format
	for key, value := range formatMap {
		result = strings.ReplaceAll(result, key, value)
	}

	return result
}

func replaceNamePlaceholder(template string, originalName string) string {
	template = strings.ReplaceAll(template, "{name}", originalName)

	re := regexp.MustCompile(`\{name:(\d+),(\d+)\}`)
	matches := re.FindAllStringSubmatch(template, -1)
	for _, match := range matches {
		placeholder := match[0]
		start, _ := strconv.Atoi(match[1])
		length, _ := strconv.Atoi(match[2])
		if start < 0 {
			start = 0
		}
		if start >= len(originalName) {
			template = strings.Replace(template, placeholder, "", 1)
			continue
		}
		end := start + length
		if end > len(originalName) {
			end = len(originalName)
		}
		extracted := originalName[start:end]
		template = strings.Replace(template, placeholder, extracted, 1)
	}

	re = regexp.MustCompile(`\{name:(\d+)\}`)
	matches = re.FindAllStringSubmatch(template, -1)
	for _, match := range matches {
		placeholder := match[0]
		length, _ := strconv.Atoi(match[1])
		if length < 0 {
			length = len(originalName) + length
			if length < 0 {
				length = 0
			}
			if length > len(originalName) {
				template = strings.Replace(template, placeholder, originalName, 1)
			} else {
				template = strings.Replace(template, placeholder, originalName[len(originalName)-length:], 1)
			}
		} else {
			if length > len(originalName) {
				length = len(originalName)
			}
			template = strings.Replace(template, placeholder, originalName[:length], 1)
		}
	}

	return template
}

func replaceExtPlaceholder(template string, originalExt string) string {
	ext := strings.TrimPrefix(originalExt, ".")
	template = strings.ReplaceAll(template, "{ext}", ext)
	template = strings.ReplaceAll(template, "{EXT}", strings.ToUpper(ext))
	return template
}

func replaceSizePlaceholder(template string, width, height int) string {
	template = strings.ReplaceAll(template, "{width}", strconv.Itoa(width))
	template = strings.ReplaceAll(template, "{height}", strconv.Itoa(height))
	template = strings.ReplaceAll(template, "{size}", fmt.Sprintf("%dx%d", width, height))
	return template
}

func replaceTotalPlaceholder(template string, total int) string {
	return strings.ReplaceAll(template, "{total}", strconv.Itoa(total))
}

func GetOriginalNameFromPath(path string) (string, string) {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	return name, ext
}
