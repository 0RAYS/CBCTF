package model

import (
	"database/sql"
)

const (
	ChallengeFileType = "file"
	PictureFileType   = "picture"
	WriteupFileType   = "writeup"
	TrafficFileType   = "traffic"
)

// File
// BelongsTo Admin
// BelongsTo User
// BelongsTo Team
// BelongsTo Contest
type File struct {
	AdminID     sql.Null[uint] `gorm:"default:null" json:"admin_id"`
	UserID      sql.Null[uint] `gorm:"default:null" json:"user_id"`
	TeamID      sql.Null[uint] `gorm:"default:null" json:"team_id"`
	ContestID   sql.Null[uint] `gorm:"default:null" json:"contest_id"`
	OauthID     sql.Null[uint] `gorm:"default:null" json:"oauth_id"`
	ChallengeID sql.Null[uint] `gorm:"default:null" json:"challenge_id"`
	RandID      string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"rand_id"`
	Filename    string         `json:"filename"`
	Size        int64          `json:"size"`
	Path        string         `json:"-"`
	Suffix      string         `json:"suffix"`
	Hash        string         `json:"hash"`
	Type        string         `json:"type"`
	BaseModel
}

func (f File) GetModelName() string {
	return "File"
}

func (f File) GetBaseModel() BaseModel {
	return f.BaseModel
}

func (f File) GetUniqueField() []string {
	return []string{"id", "rand_id"}
}

func (f File) GetAllowedQueryFields() []string {
	return []string{
		"id", "rand_id", "model_name", "model_id", "filename", "size", "suffix", "hash", "type",
	}
}
