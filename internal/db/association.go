package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func AppendUserToTeam(tx *gorm.DB, user model.User, team model.Team) model.RetVal {
	res := tx.Model(&model.UserTeam{}).Create(&model.UserTeam{UserID: user.ID, TeamID: team.ID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append User to Team: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": "UserTeam", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteUserFromTeam(tx *gorm.DB, user model.User, team model.Team) model.RetVal {
	res := tx.Model(&model.UserTeam{}).Where("user_id = ? AND team_id = ?", user.ID, team.ID).
		Delete(&model.UserTeam{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete User from Team: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": "UserTeam", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func AppendUserToContest(tx *gorm.DB, user model.User, contest model.Contest) model.RetVal {
	res := tx.Model(&model.UserContest{}).Create(&model.UserContest{UserID: user.ID, ContestID: contest.ID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append User to Contest: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": "UserContest", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteUserFromContest(tx *gorm.DB, user model.User, contest model.Contest) model.RetVal {
	res := tx.Model(&model.UserContest{}).Where("user_id = ? AND contest_id = ?", user.ID, contest.ID).
		Delete(&model.UserContest{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete User from Contest: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": "UserContest", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func AppendUserToGroup(tx *gorm.DB, user model.User, group model.Group) model.RetVal {
	res := tx.Model(&model.UserGroup{}).Create(&model.UserGroup{UserID: user.ID, GroupID: group.ID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append User to Group: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": "UserGroup", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteUserFromGroup(tx *gorm.DB, user model.User, group model.Group) model.RetVal {
	res := tx.Model(&model.UserGroup{}).Where("user_id = ? AND group_id = ?", user.ID, group.ID).
		Delete(&model.UserGroup{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete User from Group: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": "UserGroup", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func AssignPermissionToRole(tx *gorm.DB, permission model.Permission, role model.Role) model.RetVal {
	res := tx.Model(&model.RolePermission{}).Create(&model.RolePermission{RoleID: role.ID, PermissionID: permission.ID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to assign Role Permission: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": "RolePermission", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func RevokePermissionFromRole(tx *gorm.DB, permission model.Permission, role model.Role) model.RetVal {
	res := tx.Model(&model.RolePermission{}).Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).
		Delete(&model.RolePermission{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to revoke Role Permission: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": "RolePermission", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func GetGroupUsers(tx *gorm.DB, group model.Group, limit, offset int) ([]model.User, int64, model.RetVal) {
	var count int64
	if res := tx.Model(&model.UserGroup{}).Where("group_id = ?", group.ID).Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to count Group Users: %s", res.Error)
		return nil, 0, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": "UserGroup", "Error": res.Error.Error()}}
	}
	var users []model.User
	res := tx.Raw(`
		SELECT users.* FROM users
		INNER JOIN user_groups ON user_groups.user_id = users.id
		WHERE user_groups.group_id = ? AND users.deleted_at IS NULL
		ORDER BY users.id
		LIMIT ? OFFSET ?
	`, group.ID, limit, offset).Scan(&users)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Group Users: %s", res.Error)
		return nil, 0, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.User{}.ModelName(), "Error": res.Error.Error()}}
	}
	return users, count, model.SuccessRetVal()
}
