package model

type Device struct {
	UserID uint   `gorm:"index:idx_user_id_magic,unique;" json:"user_id"`
	Magic  string `gorm:"index:idx_user_id_magic,unique;" json:"magic"`
}
