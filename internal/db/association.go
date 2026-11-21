package db

import (
	"fmt"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

// appendAssociation inserts one join table row and reports a failure with the
// shared association context.
func appendAssociation[M any](tx *gorm.DB, row M, name, target, action, errKey string) model.RetVal {
	res := tx.Model(new(M)).Create(&row)
	if res.Error != nil {
		log.Logger.Warningf("Failed to %s: %s (target=%s)", action, res.Error, target)
		return model.FailRetVal(errKey, name, target, res.Error)
	}
	return model.SuccessRetVal()
}

// deleteAssociation removes the join table rows matching the condition.
func deleteAssociation[M any](tx *gorm.DB, name, where string, args []any, target, action, errKey string) model.RetVal {
	res := tx.Model(new(M)).Where(where, args...).Delete(new(M))
	if res.Error != nil {
		log.Logger.Warningf("Failed to %s: %s (target=%s)", action, res.Error, target)
		return model.FailRetVal(errKey, name, target, res.Error)
	}
	return model.SuccessRetVal()
}

// deleteAssociationByField removes the join table rows selected by run and
// reports the failure with the ids in the target.
func deleteAssociationByField(name, field string, ids []uint, errKey string, run func() *gorm.DB) model.RetVal {
	if len(ids) == 0 {
		return model.SuccessRetVal()
	}
	res := run()
	if res.Error != nil {
		target := fmt.Sprintf("%s_id=%v", field, ids)
		log.Logger.Warningf("Failed to delete %s by %s IDs %v: %s (target=%s)", name, field, ids, res.Error, target)
		return model.FailRetVal(errKey, name, target, res.Error)
	}
	return model.SuccessRetVal()
}

// deleteAssociationIDs removes the join table rows of the given ids.
func deleteAssociationIDs[M any](tx *gorm.DB, name, field string, ids []uint, errKey string) model.RetVal {
	return deleteAssociationByField(name, field, ids, errKey, func() *gorm.DB {
		return tx.Where(field+" IN ?", ids).Delete(new(M))
	})
}

func AppendUserToTeam(tx *gorm.DB, user model.User, team model.Team) model.RetVal {
	return appendAssociation(tx, model.UserTeam{UserID: user.ID, TeamID: team.ID},
		"UserTeam", fmt.Sprintf("team_id=%d,user_id=%d", team.ID, user.ID), "append User to Team", i18n.Model.UserTeam.CreateError)
}

func DeleteUserFromTeam(tx *gorm.DB, user model.User, team model.Team) model.RetVal {
	return deleteAssociation[model.UserTeam](tx, "UserTeam", "user_id = ? AND team_id = ?", []any{user.ID, team.ID},
		fmt.Sprintf("team_id=%d,user_id=%d", team.ID, user.ID), "delete User from Team", i18n.Model.UserTeam.DeleteError)
}

func DeleteUserTeamByUserID(tx *gorm.DB, userIDL ...uint) model.RetVal {
	return deleteAssociationIDs[model.UserTeam](tx, "UserTeam", "user", userIDL, i18n.Model.UserTeam.DeleteError)
}

func DeleteUserTeamByTeamID(tx *gorm.DB, teamIDL ...uint) model.RetVal {
	return deleteAssociationIDs[model.UserTeam](tx, "UserTeam", "team", teamIDL, i18n.Model.UserTeam.DeleteError)
}

func AppendUserToContest(tx *gorm.DB, user model.User, contest model.Contest) model.RetVal {
	return appendAssociation(tx, model.UserContest{UserID: user.ID, ContestID: contest.ID},
		"UserContest", fmt.Sprintf("contest_id=%d,user_id=%d", contest.ID, user.ID), "append User to Contest", i18n.Model.UserContest.CreateError)
}

func DeleteUserFromContest(tx *gorm.DB, user model.User, contest model.Contest) model.RetVal {
	return deleteAssociation[model.UserContest](tx, "UserContest", "user_id = ? AND contest_id = ?", []any{user.ID, contest.ID},
		fmt.Sprintf("contest_id=%d,user_id=%d", contest.ID, user.ID), "delete User from Contest", i18n.Model.UserContest.DeleteError)
}

func DeleteUserContestByUserID(tx *gorm.DB, userIDL ...uint) model.RetVal {
	return deleteAssociationIDs[model.UserContest](tx, "UserContest", "user", userIDL, i18n.Model.UserContest.DeleteError)
}

// DeleteUserContestByTeamID 按队伍删除参赛关系，队伍与比赛的关系由队伍自身决定。
func DeleteUserContestByTeamID(tx *gorm.DB, teamIDL ...uint) model.RetVal {
	return deleteAssociationByField("UserContest", "team", teamIDL, i18n.Model.UserContest.DeleteError, func() *gorm.DB {
		return tx.Exec(`
			DELETE FROM user_contests
			USING user_teams, teams
			WHERE user_contests.user_id = user_teams.user_id
				AND user_contests.contest_id = teams.contest_id
				AND user_teams.team_id = teams.id
				AND teams.id IN ?
		`, teamIDL)
	})
}

func AppendUserToGroup(tx *gorm.DB, user model.User, group model.Group) model.RetVal {
	return appendAssociation(tx, model.UserGroup{UserID: user.ID, GroupID: group.ID},
		"UserGroup", fmt.Sprintf("group_id=%d,user_id=%d", group.ID, user.ID), "append User to Group", i18n.Model.UserGroup.CreateError)
}

func DeleteUserFromGroup(tx *gorm.DB, user model.User, group model.Group) model.RetVal {
	return deleteAssociation[model.UserGroup](tx, "UserGroup", "user_id = ? AND group_id = ?", []any{user.ID, group.ID},
		fmt.Sprintf("group_id=%d,user_id=%d", group.ID, user.ID), "delete User from Group", i18n.Model.UserGroup.DeleteError)
}

func DeleteUserGroupByUserID(tx *gorm.DB, userIDL ...uint) model.RetVal {
	return deleteAssociationIDs[model.UserGroup](tx, "UserGroup", "user", userIDL, i18n.Model.UserGroup.DeleteError)
}

func DeleteUserGroupByGroupID(tx *gorm.DB, groupIDL ...uint) model.RetVal {
	return deleteAssociationIDs[model.UserGroup](tx, "UserGroup", "group", groupIDL, i18n.Model.UserGroup.DeleteError)
}

func AssignPermissionToRole(tx *gorm.DB, permission model.Permission, role model.Role) model.RetVal {
	return appendAssociation(tx, model.RolePermission{RoleID: role.ID, PermissionID: permission.ID},
		"RolePermission", fmt.Sprintf("permission_id=%d,role_id=%d", permission.ID, role.ID), "assign Role Permission", i18n.Model.RolePermission.CreateError)
}

func RevokePermissionFromRole(tx *gorm.DB, permission model.Permission, role model.Role) model.RetVal {
	return deleteAssociation[model.RolePermission](tx, "RolePermission", "role_id = ? AND permission_id = ?", []any{role.ID, permission.ID},
		fmt.Sprintf("permission_id=%d,role_id=%d", permission.ID, role.ID), "revoke Role Permission", i18n.Model.RolePermission.DeleteError)
}

func DeleteRolePermissionByRoleID(tx *gorm.DB, roleIDL ...uint) model.RetVal {
	return deleteAssociationIDs[model.RolePermission](tx, "RolePermission", "role", roleIDL, i18n.Model.RolePermission.DeleteError)
}
