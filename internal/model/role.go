package model

const (
	AdminRoleName     = "admin"
	OrganizerRoleName = "organizer"
	UserRoleName      = "user"
)

var DefaultRoles = []Role{
	{Name: AdminRoleName, Description: "系统管理员, 拥有全部权限", Default: true},
	{Name: OrganizerRoleName, Description: "赛事主办方, 拥有赛事相关管理权限", Default: true},
	{Name: UserRoleName, Description: "参赛选手, 拥有参赛相关权限", Default: true},
}

var DefaultRolePermissionMap = map[string][]string{
	AdminRoleName: {
		PermUserCreate, PermUserRead, PermUserUpdate, PermUserDelete, PermUserList,
		PermRoleCreate, PermRoleRead, PermRoleUpdate, PermRoleDelete, PermRoleList, PermRoleAssign, PermRoleRevoke,
		PermPermissionRead, PermPermissionList, PermPermissionAssign, PermPermissionRevoke,
		PermGroupCreate, PermGroupRead, PermGroupUpdate, PermGroupDelete, PermGroupList,
	},
	OrganizerRoleName: {
		PermUserRead, PermUserDelete, PermUserList,
	},
	UserRoleName: {
		PermUserRead, PermUserDelete,
	},
}

type Role struct {
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"-"`
	Name        string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description string       `json:"description"`
	Default     bool         `json:"default"`
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
	return []string{"id", "name", "description", "default"}
}
