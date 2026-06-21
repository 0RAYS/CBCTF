package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// ResizePicture only support for png jpg jpeg gif
func ResizePicture(path string, width, height int) error {
	if strings.EqualFold(filepath.Ext(path), ".gif") {
		return resizeGIF(path, width, height)
	}
	return resizeImage(path, width, height)
}

func resizeImage(path string, width, height int) error {
	img, err := openImage(path)
	if err != nil {
		return err
	}
	oldWidth, oldHeight := img.Bounds().Dx(), img.Bounds().Dy()
	if oldWidth < oldHeight {
		height = int(float64(oldHeight) * (float64(width) / float64(oldWidth)))
	} else {
		width = int(float64(oldWidth) * (float64(height) / float64(oldHeight)))
	}
	if err = os.Rename(path, fmt.Sprintf("%s.bak", path)); err != nil {
		return err
	}
	return saveImage(resizeLanczos(img, width, height), path)
}

func resizeGIF(path string, width, height int) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	img, err := gif.DecodeAll(file)
	if err != nil {
		return err
	}
	if len(img.Image) == 0 {
		return fmt.Errorf("gif has no frames")
	}

	oldWidth, oldHeight := img.Config.Width, img.Config.Height
	if oldWidth == 0 || oldHeight == 0 {
		bounds := img.Image[0].Bounds()
		oldWidth, oldHeight = bounds.Dx(), bounds.Dy()
	}
	if oldWidth < oldHeight {
		height = int(float64(oldHeight) * (float64(width) / float64(oldWidth)))
	} else {
		width = int(float64(oldWidth) * (float64(height) / float64(oldHeight)))
	}

	canvasBounds := image.Rect(0, 0, oldWidth, oldHeight)
	canvas := image.NewRGBA(canvasBounds)
	resized := &gif.GIF{
		Image:           make([]*image.Paletted, 0, len(img.Image)),
		Delay:           append([]int(nil), img.Delay...),
		LoopCount:       img.LoopCount,
		Disposal:        append([]byte(nil), img.Disposal...),
		BackgroundIndex: img.BackgroundIndex,
		Config: image.Config{
			ColorModel: img.Config.ColorModel,
			Width:      width,
			Height:     height,
		},
	}

	for i, frame := range img.Image {
		var previous *image.RGBA
		if i < len(img.Disposal) && img.Disposal[i] == gif.DisposalPrevious {
			previous = image.NewRGBA(canvas.Bounds())
			draw.Draw(previous, canvas.Bounds(), canvas, canvas.Bounds().Min, draw.Src)
		}

		draw.Draw(canvas, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)
		resized.Image = append(resized.Image, func(src image.Image) *image.Paletted {
			colors := make([]color.Color, 0, 256)
			colors = append(colors, color.Transparent)
			colors = append(colors, palette.Plan9[:255]...)

			dst := image.NewPaletted(src.Bounds(), colors)
			draw.FloydSteinberg.Draw(dst, src.Bounds(), src, src.Bounds().Min)
			return dst
		}(resizeLanczos(canvas, width, height)))

		if i >= len(img.Disposal) {
			continue
		}
		switch img.Disposal[i] {
		case gif.DisposalBackground:
			draw.Draw(canvas, frame.Bounds(), image.Transparent, image.Point{}, draw.Src)
		case gif.DisposalPrevious:
			if previous != nil {
				canvas = previous
			}
		}
	}

	if err = os.Rename(path, fmt.Sprintf("%s.bak", path)); err != nil {
		return err
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	return gif.EncodeAll(out, resized)
}

// Portions of the resize implementation are adapted from github.com/disintegration/imaging.
//
// The MIT License (MIT)
//
// Copyright (c) 2012 Grigory Dryapak
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

type resizeWeight struct {
	index  int
	weight float64
}

type resizeFilter struct {
	support float64
	kernel  func(float64) float64
}

var lanczosFilter = resizeFilter{
	support: 3,
	kernel: func(x float64) float64 {
		x = math.Abs(x)
		if x < 3 {
			return sinc(x) * sinc(x/3)
		}
		return 0
	},
}

func openImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	img, _, err := image.Decode(file)
	return img, err
}

func saveImage(img image.Image, path string) error {
	format := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	if format != "jpg" && format != "jpeg" && format != "png" {
		return fmt.Errorf("unsupported image format: %s", format)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	if format == "png" {
		err = png.Encode(file, img)
	} else {
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
	}

	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	return err
}

func resizeLanczos(img image.Image, width, height int) *image.NRGBA {
	return resizeImageWithFilter(img, width, height, lanczosFilter)
}

func resizeImageWithFilter(img image.Image, width, height int, filter resizeFilter) *image.NRGBA {
	dstW, dstH := width, height
	if dstW < 0 || dstH < 0 || (dstW == 0 && dstH == 0) {
		return &image.NRGBA{}
	}

	srcW, srcH := img.Bounds().Dx(), img.Bounds().Dy()
	if srcW <= 0 || srcH <= 0 {
		return &image.NRGBA{}
	}

	if dstW == 0 {
		dstW = int(math.Max(1, math.Floor(float64(dstH)*float64(srcW)/float64(srcH)+0.5)))
	}
	if dstH == 0 {
		dstH = int(math.Max(1, math.Floor(float64(dstW)*float64(srcH)/float64(srcW)+0.5)))
	}

	if srcW != dstW && srcH != dstH {
		return resizeVertical(resizeHorizontal(img, dstW, filter), dstH, filter)
	}
	if srcW != dstW {
		return resizeHorizontal(img, dstW, filter)
	}
	if srcH != dstH {
		return resizeVertical(img, dstH, filter)
	}
	return cloneNRGBA(img)
}

func resizeHorizontal(img image.Image, width int, filter resizeFilter) *image.NRGBA {
	src := newImageScanner(img)
	dst := image.NewNRGBA(image.Rect(0, 0, width, src.height))
	weights := precomputeResizeWeights(width, src.width, filter)

	parallelRange(0, src.height, func(ys <-chan int) {
		scanLine := make([]uint8, src.width*4)
		for y := range ys {
			src.scan(0, y, src.width, y+1, scanLine)
			dstOffset := y * dst.Stride
			for x := range weights {
				var r, g, b, a float64
				for _, weight := range weights[x] {
					srcOffset := weight.index * 4
					s := scanLine[srcOffset : srcOffset+4 : srcOffset+4]
					aw := float64(s[3]) * weight.weight
					r += float64(s[0]) * aw
					g += float64(s[1]) * aw
					b += float64(s[2]) * aw
					a += aw
				}
				if a != 0 {
					dstPixel := dst.Pix[dstOffset+x*4 : dstOffset+x*4+4 : dstOffset+x*4+4]
					aInv := 1 / a
					dstPixel[0] = clampUint8(r * aInv)
					dstPixel[1] = clampUint8(g * aInv)
					dstPixel[2] = clampUint8(b * aInv)
					dstPixel[3] = clampUint8(a)
				}
			}
		}
	})

	return dst
}

func resizeVertical(img image.Image, height int, filter resizeFilter) *image.NRGBA {
	src := newImageScanner(img)
	dst := image.NewNRGBA(image.Rect(0, 0, src.width, height))
	weights := precomputeResizeWeights(height, src.height, filter)

	parallelRange(0, src.width, func(xs <-chan int) {
		scanLine := make([]uint8, src.height*4)
		for x := range xs {
			src.scan(x, 0, x+1, src.height, scanLine)
			for y := range weights {
				var r, g, b, a float64
				for _, weight := range weights[y] {
					srcOffset := weight.index * 4
					s := scanLine[srcOffset : srcOffset+4 : srcOffset+4]
					aw := float64(s[3]) * weight.weight
					r += float64(s[0]) * aw
					g += float64(s[1]) * aw
					b += float64(s[2]) * aw
					a += aw
				}
				if a != 0 {
					dstOffset := y*dst.Stride + x*4
					dstPixel := dst.Pix[dstOffset : dstOffset+4 : dstOffset+4]
					aInv := 1 / a
					dstPixel[0] = clampUint8(r * aInv)
					dstPixel[1] = clampUint8(g * aInv)
					dstPixel[2] = clampUint8(b * aInv)
					dstPixel[3] = clampUint8(a)
				}
			}
		}
	})

	return dst
}

func precomputeResizeWeights(dstSize, srcSize int, filter resizeFilter) [][]resizeWeight {
	scale := float64(srcSize) / float64(dstSize)
	filterScale := scale
	if filterScale < 1 {
		filterScale = 1
	}
	radius := math.Ceil(filterScale * filter.support)

	weights := make([][]resizeWeight, dstSize)
	tmp := make([]resizeWeight, 0, dstSize*int(radius+2)*2)
	for dst := 0; dst < dstSize; dst++ {
		center := (float64(dst)+0.5)*scale - 0.5
		begin := int(math.Ceil(center - radius))
		if begin < 0 {
			begin = 0
		}
		end := int(math.Floor(center + radius))
		if end > srcSize-1 {
			end = srcSize - 1
		}

		var sum float64
		for src := begin; src <= end; src++ {
			weight := filter.kernel((float64(src) - center) / filterScale)
			if weight != 0 {
				sum += weight
				tmp = append(tmp, resizeWeight{index: src, weight: weight})
			}
		}
		if sum != 0 {
			for i := range tmp {
				tmp[i].weight /= sum
			}
		}

		weights[dst] = tmp
		tmp = tmp[len(tmp):]
	}

	return weights
}

type imageScanner struct {
	image         image.Image
	width, height int
}

func newImageScanner(img image.Image) *imageScanner {
	return &imageScanner{
		image:  img,
		width:  img.Bounds().Dx(),
		height: img.Bounds().Dy(),
	}
}

func (s *imageScanner) scan(x1, y1, x2, y2 int, dst []uint8) {
	bounds := s.image.Bounds()
	idx := 0
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			c := color.NRGBAModel.Convert(s.image.At(bounds.Min.X+x, bounds.Min.Y+y)).(color.NRGBA)
			dst[idx+0] = c.R
			dst[idx+1] = c.G
			dst[idx+2] = c.B
			dst[idx+3] = c.A
			idx += 4
		}
	}
}

func cloneNRGBA(img image.Image) *image.NRGBA {
	src := newImageScanner(img)
	dst := image.NewNRGBA(image.Rect(0, 0, src.width, src.height))
	rowSize := src.width * 4
	parallelRange(0, src.height, func(ys <-chan int) {
		for y := range ys {
			offset := y * dst.Stride
			src.scan(0, y, src.width, y+1, dst.Pix[offset:offset+rowSize])
		}
	})
	return dst
}

func parallelRange(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}

	procs := runtime.GOMAXPROCS(0)
	if procs > count {
		procs = count
	}

	items := make(chan int, count)
	for i := start; i < stop; i++ {
		items <- i
	}
	close(items)

	var wg sync.WaitGroup
	for i := 0; i < procs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(items)
		}()
	}
	wg.Wait()
}

func clampUint8(x float64) uint8 {
	value := int64(x + 0.5)
	if value > 255 {
		return 255
	}
	if value > 0 {
		return uint8(value)
	}
	return 0
}

func sinc(x float64) float64 {
	if x == 0 {
		return 1
	}
	return math.Sin(math.Pi*x) / (math.Pi * x)
}
