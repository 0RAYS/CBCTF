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
	Hash        string         `json:"hash"`
	Admin       bool           `json:"isadmin"`
	Attachment  bool           `json:"isattachment"`
	ChallengeID uint           `json:"challenge_id"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index" `
}

func InitFile(path string, uploader uint, file *multipart.FileHeader, hash string, isAdmin bool, isAttachment bool, idL ...uint) File {
	tmp := File{
		ID:         utils.RandomString(),
		Filename:   file.Filename,
		Size:       file.Size,
		Path:       path,
		Uploader:   uploader,
		Hash:       hash,
		Type:       strings.ToLower(p.Ext(file.Filename)),
		Admin:      isAdmin,
		Attachment: isAttachment,
	}
	if len(idL) > 0 {
		tmp.ChallengeID = idL[0]
	}
	return tmp
}
