package model

import (
	"CBCTF/internal/i18n"
	"database/sql"
)

const (
	ChallengeFile = "file"
	AvatarFile    = "avatar"
	WriteUPFile   = "writeup"
)

// File
// BelongsTo Admin
// BelongsTo User
// BelongsTo Team
// BelongsTo Contest
type File struct {
	AdminID   sql.Null[uint] `gorm:"default:null" json:"admin_id"`
	Admin     *Admin         `json:"-"`
	UserID    sql.Null[uint] `gorm:"default:null" json:"user_id"`
	User      *User          `json:"-"`
	TeamID    sql.Null[uint] `gorm:"default:null" json:"team_id"`
	Team      *Team          `json:"-"`
	ContestID sql.Null[uint] `gorm:"default:null" json:"contest_id"`
	Contest   *Contest       `json:"-"`
	OauthID   sql.Null[uint] `gorm:"default:null" json:"oauth_id"`
	Oauth     *Oauth         `json:"-"`
	RandID    string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"rand_id"`
	Filename  string         `json:"filename"`
	Size      int64          `json:"size"`
	Path      string         `json:"-"`
	Suffix    string         `json:"suffix"`
	Hash      string         `json:"hash"`
	Type      string         `json:"type"`
	BasicModel
}

func (f File) GetModelName() string {
	return "File"
}

func (f File) GetVersion() uint {
	return f.Version
}

func (f File) GetBasicModel() BasicModel {
	return f.BasicModel
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
