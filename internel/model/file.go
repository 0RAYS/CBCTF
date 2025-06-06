package model

import "CBCTF/internel/i18n"

const (
	ChallengeFile = "file"
	AvatarFile    = "avatar"
	WriteUPFile   = "writeup"
)

type File struct {
	RandID    string `gorm:"uniqueIndex;not null" json:"rand_id"`
	Filename  string `json:"filename"`
	Size      int64  `json:"size"`
	Path      string `json:"-"`
	AdminID   uint   `json:"admin_id"`
	UserID    uint   `json:"user_id"`
	TeamID    uint   `json:"team_id"`
	ContestID uint   `json:"contest_id"`
	Suffix    string `json:"suffix"`
	Hash      string `json:"hash"`
	Type      string `json:"type"`
	Basic
}

func (f File) GetModelName() string {
	return "File"
}

func (f File) GetID() uint {
	return f.ID
}

func (f File) GetVersion() uint {
	return f.Version
}

func (f File) CreateErrorString() string {
	return i18n.CreateFileError
}

func (f File) DeleteErrorString() string {
	return i18n.DeleteFileError
}

func (f File) GetErrorString() string {
	return i18n.GetFileError
}

func (f File) NotFoundErrorString() string {
	return i18n.FileNotFound
}

func (f File) UpdateErrorString() string {
	return i18n.UpdateFileError
}

func (f File) GetUniqueKey() []string {
	// 虽然hash并不唯一, 但不影响功能
	return []string{"id", "rand_id", "hash"}
}
