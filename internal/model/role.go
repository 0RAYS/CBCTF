package model

type Role struct {
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"-"`
	Name        string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description string       `json:"description"`
	BaseModel
}

func (r Role) TableName() string {
	return "roles"
}

func (r Role) ModelName() string {
	return "Role"
}

func (r Role) GetBaseModel() BaseModel {
	return r.BaseModel
}

func (r Role) UniqueFields() []string {
	return []string{"id", "name"}
}

func (r Role) QueryFields() []string {
	return []string{"id", "name", "description"}
}
