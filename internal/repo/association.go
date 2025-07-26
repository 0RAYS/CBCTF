package repo

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"gorm.io/gorm"
)

// AppendUserToTeam Many2Many
func AppendUserToTeam(tx *gorm.DB, userID, teamID uint) (bool, string) {
	res := tx.Model(&model.UserTeam{}).Create(&model.UserTeam{UserID: userID, TeamID: teamID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append user to team: %v", res.Error)
		return false, i18n.AppendUserToTeamError
	}
	return true, i18n.Success
}

// AppendUserToContest Many2Many
func AppendUserToContest(tx *gorm.DB, userID, contestID uint) (bool, string) {
	res := tx.Model(&model.UserContest{}).Create(&model.UserContest{UserID: userID, ContestID: contestID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append user to contest: %v", res.Error)
		return false, i18n.AppendUserToContestError
	}
	return true, i18n.Success
}

// DeleteUserFromTeam Many2Many
func DeleteUserFromTeam(tx *gorm.DB, userID, teamID uint) (bool, string) {
	res := tx.Model(&model.UserTeam{}).Where("user_id = ? AND team_id = ?", userID, teamID).
		Delete(&model.UserTeam{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete user from team: %v", res.Error)
		return false, i18n.DeleteUserFromTeamError
	}
	return true, i18n.Success
}

// DeleteUserFromContest Many2Many
func DeleteUserFromContest(tx *gorm.DB, userID, contestID uint) (bool, string) {
	res := tx.Model(&model.UserContest{}).Where("user_id = ? AND contest_id = ?", userID, contestID).
		Delete(&model.UserContest{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete user from contest: %v", res.Error)
		return false, i18n.DeleteUserFromContestError
	}
	return true, i18n.Success
}
