package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

func ResizePicture(path string, width, height int) error {
	if strings.EqualFold(filepath.Ext(path), ".gif") {
		return resizeGIF(path, width, height)
	}
	return resizeImage(path, width, height)
}

func resizeImage(path string, width, height int) error {
	img, err := imaging.Open(path)
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
	return imaging.Save(imaging.Resize(img, width, height, imaging.Lanczos), path)
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
		}(imaging.Resize(canvas, width, height, imaging.Lanczos)))

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
