package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func GetTeamIDByUserID(tx *gorm.DB, userID uint) ([]uint, bool, string) {
	var uts []model.UserTeam
	res := tx.Model(&model.UserTeam{}).Where("user_id = ?", userID).Find(&uts)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get team: %s", res.Error)
		return nil, false, i18n.GetTeamError
	}
	var idL []uint
	for _, ut := range uts {
		idL = append(idL, ut.TeamID)
	}
	return idL, true, i18n.Success
}

func GetContestIDByUserID(tx *gorm.DB, userID uint) ([]uint, bool, string) {
	var ucs []model.UserContest
	res := tx.Model(&model.UserContest{}).Where("user_id = ?", userID).Find(&ucs)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get contest: %s", res.Error)
		return nil, false, i18n.GetContestError
	}
	var idL []uint
	for _, uc := range ucs {
		idL = append(idL, uc.ContestID)
	}
	return idL, true, i18n.Success
}

func GetUserIDByTeamID(tx *gorm.DB, teamID uint) ([]uint, bool, string) {
	var uts []model.UserTeam
	res := tx.Model(&model.UserTeam{}).Where("team_id = ?", teamID).Find(&uts)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get user: %s", res.Error)
		return nil, false, i18n.GetUserError
	}
	var idL []uint
	for _, ut := range uts {
		idL = append(idL, ut.UserID)
	}
	return idL, true, i18n.Success
}

func GetUserIDByContestID(tx *gorm.DB, contestID uint) ([]uint, bool, string) {
	var ucs []model.UserContest
	res := tx.Model(&model.UserContest{}).Where("contest_id = ?", contestID).Find(&ucs)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get user: %s", res.Error)
		return nil, false, i18n.GetUserError
	}
	var idL []uint
	for _, uc := range ucs {
		idL = append(idL, uc.UserID)
	}
	return idL, true, i18n.Success
}

// AppendUserToTeam Many2Many
func AppendUserToTeam(tx *gorm.DB, user model.User, team model.Team) (bool, string) {
	res := tx.Model(&model.UserTeam{}).Create(&model.UserTeam{UserID: user.ID, TeamID: team.ID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append User to Team: %s", res.Error)
		return false, i18n.AppendUserToTeamError
	}
	return true, i18n.Success
}

// AppendUserToContest Many2Many
func AppendUserToContest(tx *gorm.DB, user model.User, contest model.Contest) (bool, string) {
	res := tx.Model(&model.UserContest{}).Create(&model.UserContest{UserID: user.ID, ContestID: contest.ID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append User to Contest: %s", res.Error)
		return false, i18n.AppendUserToContestError
	}
	return true, i18n.Success
}

// DeleteUserFromTeam Many2Many
func DeleteUserFromTeam(tx *gorm.DB, user model.User, team model.Team) (bool, string) {
	res := tx.Model(&model.UserTeam{}).Where("user_id = ? AND team_id = ?", user.ID, team.ID).
		Delete(&model.UserTeam{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete User from Team: %s", res.Error)
		return false, i18n.DeleteUserFromTeamError
	}
	return true, i18n.Success
}

// DeleteUserFromContest Many2Many
func DeleteUserFromContest(tx *gorm.DB, user model.User, contest model.Contest) (bool, string) {
	res := tx.Model(&model.UserContest{}).Where("user_id = ? AND contest_id = ?", user.ID, contest.ID).
		Delete(&model.UserContest{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete User from Contest: %s", res.Error)
		return false, i18n.DeleteUserFromContestError
	}
	return true, i18n.Success
}
