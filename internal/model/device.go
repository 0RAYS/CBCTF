package model

type Device struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `gorm:"index:idx_user_id_magic,unique;" json:"user_id"`
	Magic  string `gorm:"index:idx_user_id_magic,unique;" json:"magic"`
	Count  int    `json:"count"`
}
