package model

type Group struct {
	Users       []User `gorm:"many2many:user_groups;" json:"-"`
	RoleID      uint   `gorm:"default:null" json:"role_id"`
	Name        string `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description string `json:"description"`
	BaseModel
}

func (g Group) TableName() string {
	return "groups"
}

func (g Group) ModelName() string {
	return "Group"
}

func (g Group) GetBaseModel() BaseModel {
	return g.BaseModel
}

func (g Group) UniqueFields() []string {
	return []string{"id", "name"}
}

func (g Group) QueryFields() []string {
	return []string{"id", "role_id", "name", "description"}
}
