package model

type FileType string

const (
	ChallengeFileType FileType = "file"
	PictureFileType   FileType = "picture"
	WriteupFileType   FileType = "writeup"
	TrafficFileType   FileType = "traffic"
)

// File
// BelongsTo Admin
// BelongsTo User
// BelongsTo Team
// BelongsTo Contest
type File struct {
	Model    string   `gorm:"not null" json:"model"`
	ModelID  uint     `gorm:"not null" json:"model_id"`
	RandID   string   `gorm:"type:varchar(36);uniqueIndex;not null" json:"rand_id"`
	Filename string   `json:"filename"`
	Size     int64    `json:"size"`
	Path     string   `json:"-"`
	Suffix   string   `json:"suffix"`
	Hash     string   `json:"hash"`
	Type     FileType `json:"type"`
	BaseModel
}

func (f File) TableName() string {
	return "files"
}

func (f File) ModelName() string {
	return "File"
}

func (f File) GetBaseModel() BaseModel {
	return f.BaseModel
}

func (f File) UniqueFields() []string {
	return []string{"id", "rand_id"}
}

func (f File) QueryFields() []string {
	return []string{
		"id", "rand_id", "model", "model_id", "filename", "size", "suffix", "hash", "type",
	}
}
