package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

// AppendUserToTeam Many2Many
func AppendUserToTeam(tx *gorm.DB, userID, teamID uint) error {
	res := tx.Model(&model.UserTeam{}).Create(&model.UserTeam{UserID: userID, TeamID: teamID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append user to team: %v", res.Error)
	}
	return res.Error
}

// AppendUserToContest Many2Many
func AppendUserToContest(tx *gorm.DB, userID, contestID uint) error {
	res := tx.Model(&model.UserContest{}).Create(&model.UserContest{UserID: userID, ContestID: contestID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append user to contest: %v", res.Error)
	}
	return res.Error
}

// DeleteUserFromTeam Many2Many
func DeleteUserFromTeam(tx *gorm.DB, userID, teamID uint) error {
	res := tx.Model(&model.UserTeam{}).Where("user_id = ? AND team_id = ?", userID, teamID).
		Delete(&model.UserTeam{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete user from team: %v", res.Error)
	}
	return res.Error
}

// DeleteUserFromContest Many2Many
func DeleteUserFromContest(tx *gorm.DB, userID, contestID uint) error {
	res := tx.Model(&model.UserContest{}).Where("user_id = ? AND contest_id = ?", userID, contestID).
		Delete(&model.UserContest{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete user from contest: %v", res.Error)
	}
	return res.Error
}
