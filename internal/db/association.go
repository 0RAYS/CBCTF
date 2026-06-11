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
		return model.RetVal{Msg: i18n.Model.UserTeam.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteUserFromTeam(tx *gorm.DB, user model.User, team model.Team) model.RetVal {
	res := tx.Model(&model.UserTeam{}).Where("user_id = ? AND team_id = ?", user.ID, team.ID).
		Delete(&model.UserTeam{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete User from Team: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.UserTeam.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteUserTeamByUserID(tx *gorm.DB, userIDL ...uint) model.RetVal {
	if len(userIDL) == 0 {
		return model.SuccessRetVal()
	}
	res := tx.Where("user_id IN ?", userIDL).Delete(&model.UserTeam{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete UserTeam by user IDs %v: %s", userIDL, res.Error)
		return model.RetVal{Msg: i18n.Model.UserTeam.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteUserTeamByTeamID(tx *gorm.DB, teamIDL ...uint) model.RetVal {
	if len(teamIDL) == 0 {
		return model.SuccessRetVal()
	}
	res := tx.Where("team_id IN ?", teamIDL).Delete(&model.UserTeam{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete UserTeam by team IDs %v: %s", teamIDL, res.Error)
		return model.RetVal{Msg: i18n.Model.UserTeam.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func AppendUserToContest(tx *gorm.DB, user model.User, contest model.Contest) model.RetVal {
	res := tx.Model(&model.UserContest{}).Create(&model.UserContest{UserID: user.ID, ContestID: contest.ID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append User to Contest: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.UserContest.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteUserFromContest(tx *gorm.DB, user model.User, contest model.Contest) model.RetVal {
	res := tx.Model(&model.UserContest{}).Where("user_id = ? AND contest_id = ?", user.ID, contest.ID).
		Delete(&model.UserContest{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete User from Contest: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.UserContest.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteUserContestByUserID(tx *gorm.DB, userIDL ...uint) model.RetVal {
	if len(userIDL) == 0 {
		return model.SuccessRetVal()
	}
	res := tx.Where("user_id IN ?", userIDL).Delete(&model.UserContest{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete UserContest by user IDs %v: %s", userIDL, res.Error)
		return model.RetVal{Msg: i18n.Model.UserContest.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteUserContestByTeamID(tx *gorm.DB, teamIDL ...uint) model.RetVal {
	if len(teamIDL) == 0 {
		return model.SuccessRetVal()
	}
	res := tx.Exec(`
		DELETE FROM user_contests
		USING user_teams, teams
		WHERE user_contests.user_id = user_teams.user_id
			AND user_contests.contest_id = teams.contest_id
			AND user_teams.team_id = teams.id
			AND teams.id IN ?
	`, teamIDL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete UserContest by team IDs %v: %s", teamIDL, res.Error)
		return model.RetVal{Msg: i18n.Model.UserContest.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func AppendUserToGroup(tx *gorm.DB, user model.User, group model.Group) model.RetVal {
	res := tx.Model(&model.UserGroup{}).Create(&model.UserGroup{UserID: user.ID, GroupID: group.ID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append User to Group: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.UserGroup.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteUserFromGroup(tx *gorm.DB, user model.User, group model.Group) model.RetVal {
	res := tx.Model(&model.UserGroup{}).Where("user_id = ? AND group_id = ?", user.ID, group.ID).
		Delete(&model.UserGroup{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete User from Group: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.UserGroup.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteUserGroupByUserID(tx *gorm.DB, userIDL ...uint) model.RetVal {
	if len(userIDL) == 0 {
		return model.SuccessRetVal()
	}
	res := tx.Where("user_id IN ?", userIDL).Delete(&model.UserGroup{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete UserGroup by user IDs %v: %s", userIDL, res.Error)
		return model.RetVal{Msg: i18n.Model.UserGroup.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteUserGroupByGroupID(tx *gorm.DB, groupIDL ...uint) model.RetVal {
	if len(groupIDL) == 0 {
		return model.SuccessRetVal()
	}
	res := tx.Where("group_id IN ?", groupIDL).Delete(&model.UserGroup{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete UserGroup by group IDs %v: %s", groupIDL, res.Error)
		return model.RetVal{Msg: i18n.Model.UserGroup.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func AssignPermissionToRole(tx *gorm.DB, permission model.Permission, role model.Role) model.RetVal {
	res := tx.Model(&model.RolePermission{}).Create(&model.RolePermission{RoleID: role.ID, PermissionID: permission.ID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to assign Role Permission: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.RolePermission.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func RevokePermissionFromRole(tx *gorm.DB, permission model.Permission, role model.Role) model.RetVal {
	res := tx.Model(&model.RolePermission{}).Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).
		Delete(&model.RolePermission{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to revoke Role Permission: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.RolePermission.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteRolePermissionByRoleID(tx *gorm.DB, roleIDL ...uint) model.RetVal {
	if len(roleIDL) == 0 {
		return model.SuccessRetVal()
	}
	res := tx.Where("role_id IN ?", roleIDL).Delete(&model.RolePermission{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete RolePermission by role IDs %v: %s", roleIDL, res.Error)
		return model.RetVal{Msg: i18n.Model.RolePermission.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
