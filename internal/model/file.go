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
	ID          string         `json:"id" gorm:"primarykey"`
	Filename    string         `json:"filename"`
	Size        int64          `json:"size"`
	Path        string         `json:"-"`
	Uploader    uint           `json:"uploader"`
	Type        string         `json:"type"`
	Challenge   bool           `json:"challenge"`
	ChallengeID uint           `json:"challenge_id"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index" `
}

func InitFile(path string, uploader uint, fileHeader *multipart.FileHeader, challenge bool) File {
	return File{
		ID:        utils.RandomString(),
		Filename:  fileHeader.Filename,
		Size:      fileHeader.Size,
		Path:      path,
		Uploader:  uploader,
		Type:      strings.ToLower(p.Ext(fileHeader.Filename)),
		Challenge: false,
	}
}
