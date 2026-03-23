package model

const (
	AdminGroupName     = "admin"
	OrganizerGroupName = "organizer"
	UserGroupName      = "user"
)

var DefaultGroups = []Group{
	{Name: AdminGroupName, Description: "系统管理员", Default: true},
	{Name: OrganizerGroupName, Description: "赛事主办方", Default: true},
	{Name: UserGroupName, Description: "选手", Default: true},
}

var DefaultGroupRoleMap = map[string]string{
	AdminGroupName:     AdminRoleName,
	OrganizerGroupName: OrganizerRoleName,
	UserGroupName:      UserRoleName,
}

// Group 用户组
// ManyToMany User
// BelongsTo Role
type Group struct {
	Users       []User `gorm:"many2many:user_groups;" json:"-"`
	RoleID      uint   `gorm:"default:null;index" json:"role_id"`
	Role        Role   `json:"-"`
	Name        string `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
	BaseModel
}
