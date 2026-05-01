// Copyright (c) 2024 Barat Semet (https://github.com/barats)
// Resizem is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package rmanager

import (
	"context"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"resizem/rimage"
	"strings"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Listening for DnD events
func FileDropEventListener(ctx context.Context) {
	runtime.OnFileDrop(ctx, func(x, y int, paths []string) {
		for _, p := range paths {
			processPath(ctx, p)
		}
	})
}

// processPath processes a single path - if it's a directory, recursively scan for images
func processPath(ctx context.Context, path string) {
	info, err := os.Stat(path)
	if err != nil {
		log.Println("Error accessing path:", path, err)
		return
	}

	if info.IsDir() {
		// Recursively scan directory for image files
		filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && hasImageExt(filePath) {
				SendFilePathEvent(ctx, filePath)
			}
			return nil
		})
	} else {
		// It's a file, check if it's an image
		if hasImageExt(path) {
			SendFilePathEvent(ctx, path)
		}
	}
}

///////////////////////// FileManager /////////////////////////

type FileManager struct {
	ctx context.Context
}

func NewFileManager() *FileManager {
	return &FileManager{}
}

func (fm *FileManager) OnStartup(ctx context.Context) {
	fm.ctx = ctx
}

// Accepts request from frontend to cancel running jobs
func (fm *FileManager) CancelHandleFiles() {
	CancelJobs()
}

// Accepts files from fronted and start to handle them
func (fm *FileManager) StartHandleFiles(files []string, opts rimage.ImageOptions) {
	if len(files) > 0 {
		StartJobs(fm.ctx, files, opts)
	}
}

func (fm *FileManager) OpenFilesDialog() {
	filter := runtime.FileFilter{
		DisplayName: "Image Files",
		Pattern:     "*.jpg;*.jpeg;*.png;*.webp;*.tif;*.tiff;*.bmp;*.gif",
	}

	home, _ := os.UserHomeDir()
	if strings.EqualFold(runtime.Environment(fm.ctx).Platform, "windows") {
		// Known issue
		// https://github.com/wailsapp/wails/issues/1381
		home = ""
	}

	files, _ := runtime.OpenMultipleFilesDialog(fm.ctx, runtime.OpenDialogOptions{
		Title:            Resizem.Name,
		DefaultDirectory: home,
		Filters:          []runtime.FileFilter{filter},
	})

	for _, p := range files {
		SendFilePathEvent(fm.ctx, p)
	} //end of for
}

func (fm *FileManager) OpenDirectoryDialog() string {
	home, _ := os.UserHomeDir()
	path, _ := runtime.OpenDirectoryDialog(fm.ctx, runtime.OpenDialogOptions{
		Title:            Resizem.Name,
		DefaultDirectory: home,
	})

	return path
}

// ImageDetail contains detailed information about an image file
type ImageDetail struct {
	FilePath     string `json:"file_path"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	FileSizeStr  string `json:"file_size_str"`
	Format       string `json:"format"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	ColorMode    string `json:"color_mode"`
	CreatedTime  string `json:"created_time"`
	ModifiedTime string `json:"modified_time"`
	Directory    string `json:"directory"`
}

// GetImageDetail returns detailed information about an image file
func (fm *FileManager) GetImageDetail(filePath string) (*ImageDetail, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory: %s", filePath)
	}

	detail := &ImageDetail{
		FilePath:     filePath,
		FileName:     filepath.Base(filePath),
		FileSize:     info.Size(),
		FileSizeStr:  formatFileSize(info.Size()),
		Directory:    filepath.Dir(filePath),
		ModifiedTime: info.ModTime().Format("2006-01-02 15:04:05"),
	}

	// Get created time (platform specific)
	detail.CreatedTime = getCreatedTime(info)

	// Get image format and dimensions
	ext := strings.ToLower(filepath.Ext(filePath))
	detail.Format = strings.TrimPrefix(ext, ".")

	// Open image to get dimensions and color mode
	file, err := os.Open(filePath)
	if err != nil {
		return detail, nil
	}
	defer file.Close()

	config, format, err := image.DecodeConfig(file)
	if err == nil {
		detail.Width = config.Width
		detail.Height = config.Height
		if format != "" {
			detail.Format = format
		}
	}

	// Seek back to read image for color mode
	file.Seek(0, 0)
	img, _, err := image.Decode(file)
	if err == nil {
		detail.ColorMode = detectColorMode(img)
	}

	return detail, nil
}

// detectColorMode detects the color mode of an image
func detectColorMode(img image.Image) string {
	bounds := img.Bounds()
	if bounds.Empty() {
		return "Unknown"
	}

	// Check a few sample pixels to determine color mode
	// Sample pixels from corners and center
	samplePoints := []image.Point{
		bounds.Min,
		{bounds.Max.X - 1, bounds.Min.Y},
		{bounds.Min.X, bounds.Max.Y - 1},
		{bounds.Max.X - 1, bounds.Max.Y - 1},
		{(bounds.Min.X + bounds.Max.X) / 2, (bounds.Min.Y + bounds.Max.Y) / 2},
	}

	hasAlpha := false
	hasColor := false
	isGrayscale := true

	for _, p := range samplePoints {
		if !p.In(bounds) {
			continue
		}
		c := img.At(p.X, p.Y)
		r, g, b, a := c.RGBA()

		// Check alpha
		if a < 0xffff {
			hasAlpha = true
		}

		// Check if it's grayscale (r == g == b)
		if r != g || g != b {
			isGrayscale = false
			hasColor = true
		}

		// Check if there's any color
		if r > 0 || g > 0 || b > 0 {
			if r != g || g != b {
				hasColor = true
			}
		}
	}

	// Also check the image type for more accurate detection
	switch img.(type) {
	case *image.Gray:
		return "Grayscale"
	case *image.Gray16:
		return "Grayscale 16-bit"
	case *image.RGBA:
		if hasAlpha {
			return "RGBA"
		}
		return "RGB"
	case *image.RGBA64:
		if hasAlpha {
			return "RGBA 64-bit"
		}
		return "RGB 64-bit"
	case *image.NRGBA:
		return "NRGBA"
	case *image.NRGBA64:
		return "NRGBA 64-bit"
	case *image.Alpha:
		return "Alpha"
	case *image.Alpha16:
		return "Alpha 16-bit"
	case *image.CMYK:
		return "CMYK"
	case *image.Paletted:
		if pal, ok := img.(*image.Paletted); ok {
			// Check palette for alpha
			for _, c := range pal.Palette {
				if _, _, _, a := c.RGBA(); a < 0xffff {
					hasAlpha = true
					break
				}
			}
			if hasAlpha {
				return "Paletted (with transparency)"
			}
			return "Paletted"
		}
		return "Paletted"
	}

	// Fallback to pixel-based detection
	if isGrayscale {
		if hasAlpha {
			return "Grayscale with Alpha"
		}
		return "Grayscale"
	}
	if hasAlpha {
		return "RGBA"
	}
	return "RGB"
}

// getCreatedTime returns the creation time of a file (platform specific)
func getCreatedTime(info os.FileInfo) string {
	if stat, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
		createdTime := time.Unix(0, stat.CreationTime.Nanoseconds())
		return createdTime.Format("2006-01-02 15:04:05")
	}
	// Fallback to modification time if creation time not available
	return info.ModTime().Format("2006-01-02 15:04:05")
}

// formatFileSize converts file size to human-readable string
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// Simply, check image format by its extension
// We will do the REAL type-check when worker goroutine starts the job
func hasImageExt(path string) bool {
	ext := filepath.Ext(path)
	log.Println("ext --> " + ext)
	if len(ext) >= 3 { // file extension should be long enough, otherwise there's no point to check for it
		ext = strings.ToLower(ext)
		if strings.EqualFold(".bmp", ext) || strings.EqualFold(".gif", ext) ||
			strings.EqualFold(".jpg", ext) || strings.EqualFold(".jpeg", ext) ||
			strings.EqualFold(".png", ext) ||
			strings.EqualFold(".tif", ext) || strings.EqualFold(".tiff", ext) ||
			strings.EqualFold(".webp", ext) {
			return true
		}
		return false
	}
	return false
}
