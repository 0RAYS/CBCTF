package model

import (
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"mime/multipart"
	p "path"
	"strings"
	"time"
)

type File struct {
	ID        string         `json:"id" gorm:"primarykey"`
	Filename  string         `json:"filename"`
	Size      int64          `json:"size"`
	Path      string         `json:"-"`
	Uploader  uint           `json:"uploader"`
	Admin     bool           `json:"admin"`
	Type      string         `json:"type"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" `
}

func InitFile(path string, uploader uint, admin bool, fileHeader *multipart.FileHeader) File {
	return File{
		ID:       utils.RandomString(),
		Filename: fileHeader.Filename,
		Size:     fileHeader.Size,
		Path:     path,
		Uploader: uploader,
		Admin:    admin,
		Type:     strings.ToLower(p.Ext(fileHeader.Filename)),
	}
}
