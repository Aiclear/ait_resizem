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
	"image"
	"image/color"
	"image/draw"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func AddWatermark(baseImg image.Image, opts WatermarkOptions) (image.Image, error) {
	if opts.Type == WatermarkTypeNone {
		return baseImg, nil
	}

	result := imaging.Clone(baseImg)

	if opts.Type == WatermarkTypeImage {
		return addImageWatermark(result, opts)
	}

	if opts.Type == WatermarkTypeText {
		return addTextWatermark(result, opts)
	}

	return result, nil
}

func addImageWatermark(baseImg *image.NRGBA, opts WatermarkOptions) (image.Image, error) {
	if opts.ImagePath == "" {
		return baseImg, nil
	}

	file, err := os.Open(opts.ImagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	watermarkImg, err := imaging.Decode(file)
	if err != nil {
		return nil, err
	}

	if opts.Scale > 0 && opts.Scale != 1.0 {
		newWidth := int(float64(watermarkImg.Bounds().Dx()) * opts.Scale)
		newHeight := int(float64(watermarkImg.Bounds().Dy()) * opts.Scale)
		watermarkImg = imaging.Resize(watermarkImg, newWidth, newHeight, imaging.Lanczos)
	}

	if opts.Rotation != 0 {
		watermarkImg = imaging.Rotate(watermarkImg, opts.Rotation, color.Transparent)
	}

	if opts.Opacity < 1.0 {
		watermarkImg = applyOpacity(watermarkImg, opts.Opacity)
	}

	pos := calculateWatermarkPosition(
		baseImg.Bounds().Dx(), baseImg.Bounds().Dy(),
		watermarkImg.Bounds().Dx(), watermarkImg.Bounds().Dy(),
		opts.Position, opts.OffsetX, opts.OffsetY,
	)

	draw.Draw(baseImg, watermarkImg.Bounds().Add(pos), watermarkImg, watermarkImg.Bounds().Min, draw.Over)

	return baseImg, nil
}

func addTextWatermark(baseImg *image.NRGBA, opts WatermarkOptions) (image.Image, error) {
	if opts.Text == "" {
		return baseImg, nil
	}

	fontSize := opts.FontSize
	if fontSize <= 0 {
		fontSize = 24
	}

	face := basicfont.Face7x13

	textColor := parseColor(opts.FontColor)
	if textColor == nil {
		textColor = color.White
	}

	lines := strings.Split(opts.Text, "\n")
	lineHeight := int(fontSize) * 2

	maxWidth := 0
	for _, line := range lines {
		width := len(line) * 7
		if width > maxWidth {
			maxWidth = width
		}
	}

	textImgWidth := maxWidth + 20
	textImgHeight := len(lines)*lineHeight + 20

	textImg := image.NewNRGBA(image.Rect(0, 0, textImgWidth, textImgHeight))

	drawer := &font.Drawer{
		Dst:  textImg,
		Src:  image.NewUniform(textColor),
		Face: face,
	}

	for i, line := range lines {
		drawer.Dot = fixed.P(10, 10+(i+1)*lineHeight)
		drawer.DrawString(line)
	}

	var watermarkImg image.Image = textImg

	if opts.Scale > 0 && opts.Scale != 1.0 {
		newWidth := int(float64(watermarkImg.Bounds().Dx()) * opts.Scale)
		newHeight := int(float64(watermarkImg.Bounds().Dy()) * opts.Scale)
		watermarkImg = imaging.Resize(watermarkImg, newWidth, newHeight, imaging.Lanczos)
	}

	if opts.Rotation != 0 {
		watermarkImg = imaging.Rotate(watermarkImg, opts.Rotation, color.Transparent)
	}

	if opts.Opacity < 1.0 {
		watermarkImg = applyOpacity(watermarkImg, opts.Opacity)
	}

	pos := calculateWatermarkPosition(
		baseImg.Bounds().Dx(), baseImg.Bounds().Dy(),
		watermarkImg.Bounds().Dx(), watermarkImg.Bounds().Dy(),
		opts.Position, opts.OffsetX, opts.OffsetY,
	)

	draw.Draw(baseImg, watermarkImg.Bounds().Add(pos), watermarkImg, watermarkImg.Bounds().Min, draw.Over)

	return baseImg, nil
}

func applyOpacity(img image.Image, opacity float64) image.Image {
	if opacity >= 1.0 {
		return img
	}

	bounds := img.Bounds()
	result := image.NewNRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			c.A = uint8(float64(c.A) * opacity)
			result.Set(x, y, c)
		}
	}

	return result
}

func calculateWatermarkPosition(baseWidth, baseHeight, watermarkWidth, watermarkHeight int,
	position WatermarkPosition, offsetX, offsetY int) image.Point {

	var x, y int

	switch position {
	case PositionTopLeft:
		x = offsetX
		y = offsetY
	case PositionTopCenter:
		x = (baseWidth - watermarkWidth) / 2
		y = offsetY
	case PositionTopRight:
		x = baseWidth - watermarkWidth - offsetX
		y = offsetY
	case PositionMiddleLeft:
		x = offsetX
		y = (baseHeight - watermarkHeight) / 2
	case PositionCenter:
		x = (baseWidth - watermarkWidth) / 2
		y = (baseHeight - watermarkHeight) / 2
	case PositionMiddleRight:
		x = baseWidth - watermarkWidth - offsetX
		y = (baseHeight - watermarkHeight) / 2
	case PositionBottomLeft:
		x = offsetX
		y = baseHeight - watermarkHeight - offsetY
	case PositionBottomCenter:
		x = (baseWidth - watermarkWidth) / 2
		y = baseHeight - watermarkHeight - offsetY
	case PositionBottomRight:
		x = baseWidth - watermarkWidth - offsetX
		y = baseHeight - watermarkHeight - offsetY
	}

	x = max(x, 0)
	y = max(y, 0)

	return image.Point{X: x, Y: y}
}

func parseColor(colorStr string) color.Color {
	colorStr = strings.TrimSpace(colorStr)
	if colorStr == "" {
		return color.White
	}

	if strings.HasPrefix(colorStr, "#") {
		return parseHexColor(colorStr)
	}

	colorMap := map[string]color.Color{
		"white":   color.White,
		"black":   color.Black,
		"red":     color.RGBA{255, 0, 0, 255},
		"green":   color.RGBA{0, 255, 0, 255},
		"blue":    color.RGBA{0, 0, 255, 255},
		"yellow":  color.RGBA{255, 255, 0, 255},
		"magenta": color.RGBA{255, 0, 255, 255},
		"cyan":    color.RGBA{0, 255, 255, 255},
		"gray":    color.RGBA{128, 128, 128, 255},
	}

	if c, ok := colorMap[strings.ToLower(colorStr)]; ok {
		return c
	}

	return color.White
}

func parseHexColor(hex string) color.Color {
	hex = strings.TrimPrefix(hex, "#")

	var r, g, b, a uint8 = 0, 0, 0, 255

	switch len(hex) {
	case 3:
		r = parseHexPair(string(hex[0]) + string(hex[0]))
		g = parseHexPair(string(hex[1]) + string(hex[1]))
		b = parseHexPair(string(hex[2]) + string(hex[2]))
	case 4:
		r = parseHexPair(string(hex[0]) + string(hex[0]))
		g = parseHexPair(string(hex[1]) + string(hex[1]))
		b = parseHexPair(string(hex[2]) + string(hex[2]))
		a = parseHexPair(string(hex[3]) + string(hex[3]))
	case 6:
		r = parseHexPair(hex[0:2])
		g = parseHexPair(hex[2:4])
		b = parseHexPair(hex[4:6])
	case 8:
		r = parseHexPair(hex[0:2])
		g = parseHexPair(hex[2:4])
		b = parseHexPair(hex[4:6])
		a = parseHexPair(hex[6:8])
	}

	return color.NRGBA{R: r, G: g, B: b, A: a}
}

func parseHexPair(pair string) uint8 {
	val, err := strconv.ParseUint(pair, 16, 8)
	if err != nil {
		return 0
	}
	return uint8(val)
}

func ApplyWatermarkToImage(img ResizemImage, opts WatermarkOptions) (image.Image, error) {
	return AddWatermark(img.Data(), opts)
}
