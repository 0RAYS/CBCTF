package model

type Permission struct {
	Roles       []Role `gorm:"many2many:role_permissions" json:"-"`
	Name        string `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Resource    string `gorm:"type:varchar(255);index;not null" json:"resource"`
	Operation   string `gorm:"type:varchar(255);not null" json:"operation"`
	Description string `json:"description"`
	BaseModel
}

func (p Permission) TableName() string {
	return "permissions"
}

func (p Permission) ModelName() string {
	return "Permission"
}

func (p Permission) GetBaseModel() BaseModel {
	return p.BaseModel
}

func (p Permission) UniqueFields() []string {
	return []string{"id", "name"}
}

func (p Permission) QueryFields() []string {
	return []string{"id", "name", "resource", "operation", "description"}
}
