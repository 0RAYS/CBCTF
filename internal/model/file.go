package model

import (
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"mime/multipart"
	p "path"
	"strings"
	"time"
)

const (
	Avatar  = "avatar"
	WriteUP = "writeup"
)

type File struct {
	ID        string         `json:"id" gorm:"primarykey"`
	Filename  string         `json:"filename"`
	Size      int64          `json:"size"`
	Path      string         `json:"-"`
	Uploader  uint           `json:"uploader"`
	Suffix    string         `json:"suffix"`
	Hash      string         `json:"hash"`
	Type      string         `json:"type"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" `
}

func InitFile(path string, uploader uint, file *multipart.FileHeader, hash string, t string) File {
	tmp := File{
		ID:       utils.UUID(),
		Filename: file.Filename,
		Size:     file.Size,
		Path:     path,
		Uploader: uploader,
		Hash:     hash,
		Suffix:   strings.ToLower(p.Ext(file.Filename)),
		Type:     t,
	}
	return tmp
}
