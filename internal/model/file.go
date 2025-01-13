package model

import (
	"gorm.io/gorm"
	"mime/multipart"
	p "path"
)

type File struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Path    string `json:"path"`
	OwnerID uint   `json:"owner"`
	Type    string `json:"type"`
	gorm.Model
}

func InitFile(ownerID uint, random string, path string, fileHeader *multipart.FileHeader) File {
	return File{
		ID:      random,
		Name:    fileHeader.Filename,
		Size:    fileHeader.Size,
		Path:    path,
		OwnerID: ownerID,
		Type:    p.Ext(path),
	}
}
