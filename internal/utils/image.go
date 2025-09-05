package utils

import (
	"fmt"
	"os"

	"github.com/disintegration/imaging"
)

func ResizeImage(path string, width, height int) error {
	img, err := imaging.Open(path)
	if err != nil {
		return err
	}
	if err = os.Rename(path, fmt.Sprintf("%s.bak", path)); err != nil {
		return err
	}
	oldWidth, oldHeight := img.Bounds().Dx(), img.Bounds().Dy()
	if oldWidth < oldHeight {
		height = int(float64(oldHeight) * (float64(width) / float64(oldWidth)))
	} else {
		width = int(float64(oldWidth) * (float64(height) / float64(oldHeight)))
	}
	return imaging.Save(imaging.Resize(img, width, height, imaging.Lanczos), path)
}
