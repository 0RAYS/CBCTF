package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func GetTeamIDByUserID(tx *gorm.DB, userID uint) ([]uint, bool, string) {
	var idL []uint
	res := tx.Model(&model.UserTeam{}).Where("user_id = ?", userID).Find(&idL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get team: %s", res.Error)
		return nil, false, i18n.GetTeamError
	}
	return idL, true, i18n.Success
}

func GetContestIDByUserID(tx *gorm.DB, userID uint) ([]uint, bool, string) {
	var idL []uint
	res := tx.Model(&model.UserContest{}).Where("user_id = ?", userID).Find(&idL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get contest: %s", res.Error)
		return nil, false, i18n.GetContestError
	}
	return idL, true, i18n.Success
}

func GetUserIDByTeamID(tx *gorm.DB, teamID uint) ([]uint, bool, string) {
	var idL []uint
	res := tx.Model(&model.UserTeam{}).Where("team_id = ?", teamID).Find(&idL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get user: %s", res.Error)
		return nil, false, i18n.GetUserError
	}
	return idL, true, i18n.Success
}

func GetUserIDByContestID(tx *gorm.DB, contestID uint) ([]uint, bool, string) {
	var idL []uint
	res := tx.Model(&model.UserContest{}).Where("contest_id = ?", contestID).Find(&idL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get user: %s", res.Error)
		return nil, false, i18n.GetUserError
	}
	return idL, true, i18n.Success
}

// AppendUserToTeam Many2Many
func AppendUserToTeam(tx *gorm.DB, userID, teamID uint) (bool, string) {
	res := tx.Model(&model.UserTeam{}).Create(&model.UserTeam{UserID: userID, TeamID: teamID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append User to Team: %s", res.Error)
		return false, i18n.AppendUserToTeamError
	}
	if ok, msg := InitTeamRepo(tx).Update(teamID, UpdateTeamOptions{DiffUserCount: 1}); !ok {
		return false, msg
	}
	return InitUserRepo(tx).Update(userID, UpdateUserOptions{DiffTeamCount: 1})
}

// AppendUserToContest Many2Many
func AppendUserToContest(tx *gorm.DB, userID, contestID uint) (bool, string) {
	res := tx.Model(&model.UserContest{}).Create(&model.UserContest{UserID: userID, ContestID: contestID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append User to Contest: %s", res.Error)
		return false, i18n.AppendUserToContestError
	}
	if ok, msg := InitContestRepo(tx).Update(contestID, UpdateContestOptions{DiffUserCount: 1}); !ok {
		return false, msg
	}
	return InitUserRepo(tx).Update(userID, UpdateUserOptions{DiffContestCount: 1})
}

// DeleteUserFromTeam Many2Many
func DeleteUserFromTeam(tx *gorm.DB, userID, teamID uint) (bool, string) {
	res := tx.Model(&model.UserTeam{}).Where("user_id = ? AND team_id = ?", userID, teamID).
		Delete(&model.UserTeam{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete User from Team: %s", res.Error)
		return false, i18n.DeleteUserFromTeamError
	}
	if ok, msg := InitTeamRepo(tx).Update(teamID, UpdateTeamOptions{DiffUserCount: -1}); !ok {
		return false, msg
	}
	return InitUserRepo(tx).Update(userID, UpdateUserOptions{DiffTeamCount: -1})
}

// DeleteUserFromContest Many2Many
func DeleteUserFromContest(tx *gorm.DB, userID, contestID uint) (bool, string) {
	res := tx.Model(&model.UserContest{}).Where("user_id = ? AND contest_id = ?", userID, contestID).
		Delete(&model.UserContest{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete User from Contest: %s", res.Error)
		return false, i18n.DeleteUserFromContestError
	}
	if ok, msg := InitContestRepo(tx).Update(contestID, UpdateContestOptions{DiffUserCount: -1}); !ok {
		return false, msg
	}
	return InitUserRepo(tx).Update(userID, UpdateUserOptions{DiffContestCount: -1})
}
