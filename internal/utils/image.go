package utils

import (
	"fmt"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func ResizeImage(path string, width, height int) error {
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
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)
	if err = imaging.Save(img, fmt.Sprintf("%s.bak%s", path, filepath.Ext(path))); err != nil {
		return err
	}
	return imaging.Save(resizedImg, path)
}
