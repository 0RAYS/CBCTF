package model

const (
	PermUserCreate = "user:create"
	PermUserRead   = "user:read"
	PermUserUpdate = "user:update"
	PermUserDelete = "user:delete"
	PermUserList   = "user:list"

	PermRoleCreate = "role:create"
	PermRoleRead   = "role:read"
	PermRoleUpdate = "role:update"
	PermRoleDelete = "role:delete"
	PermRoleList   = "role:list"
	PermRoleAssign = "role:assign"
	PermRoleRevoke = "role:revoke"

	PermPermissionRead   = "permission:read"
	PermPermissionList   = "permission:list"
	PermPermissionAssign = "permission:assign"
	PermPermissionRevoke = "permission:revoke"

	PermGroupCreate = "group:create"
	PermGroupRead   = "group:read"
	PermGroupUpdate = "group:update"
	PermGroupDelete = "group:delete"
	PermGroupList   = "group:list"
)

var Permissions = []Permission{
	{Name: PermUserCreate, Resource: "user", Operation: "create", Description: "创建用户"},
	{Name: PermUserRead, Resource: "user", Operation: "read", Description: "查看用户详情"},
	{Name: PermUserUpdate, Resource: "user", Operation: "update", Description: "更新用户"},
	{Name: PermUserDelete, Resource: "user", Operation: "delete", Description: "删除用户"},
	{Name: PermUserList, Resource: "user", Operation: "list", Description: "查看用户列表"},

	{Name: PermRoleCreate, Resource: "role", Operation: "create", Description: "创建角色"},
	{Name: PermRoleRead, Resource: "role", Operation: "read", Description: "查看角色详情"},
	{Name: PermRoleUpdate, Resource: "role", Operation: "update", Description: "更新角色"},
	{Name: PermRoleDelete, Resource: "role", Operation: "delete", Description: "删除角色"},
	{Name: PermRoleList, Resource: "role", Operation: "list", Description: "查看角色列表"},
	{Name: PermRoleAssign, Resource: "role", Operation: "assign", Description: "分配角色"},
	{Name: PermRoleRevoke, Resource: "role", Operation: "revoke", Description: "移除角色"},

	{Name: PermPermissionRead, Resource: "permission", Operation: "read", Description: "查看权限详情"},
	{Name: PermPermissionList, Resource: "permission", Operation: "list", Description: "查看权限列表"},
	{Name: PermPermissionAssign, Resource: "permission", Operation: "assign", Description: "分配权限"},
	{Name: PermPermissionRevoke, Resource: "permission", Operation: "revoke", Description: "移除权限"},

	{Name: PermGroupCreate, Resource: "group", Operation: "create", Description: "创建权限分组"},
	{Name: PermGroupRead, Resource: "group", Operation: "read", Description: "查看权限分组详情"},
	{Name: PermGroupUpdate, Resource: "group", Operation: "update", Description: "更新权限分组"},
	{Name: PermGroupDelete, Resource: "group", Operation: "delete", Description: "删除权限分组"},
	{Name: PermGroupList, Resource: "group", Operation: "list", Description: "查看权限分组列表"},
}

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
