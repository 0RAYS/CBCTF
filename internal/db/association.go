package db

import (
	"CBCTF/internal/model"
	"gorm.io/gorm"
)

// AppendUserToTeam Many2Many
func AppendUserToTeam(tx *gorm.DB, user model.User, team model.Team) error {
	return tx.Model(&model.UserTeam{}).Create(&model.UserTeam{UserID: user.ID, TeamID: team.ID}).Error
}

// AppendUserToContest Many2Many
func AppendUserToContest(tx *gorm.DB, user model.User, contest model.Contest) error {
	return tx.Model(&model.UserContest{}).Create(&model.UserContest{UserID: user.ID, ContestID: contest.ID}).Error
}

// DeleteUserFromTeam Many2Many
func DeleteUserFromTeam(tx *gorm.DB, user model.User, team model.Team) error {
	return tx.Model(&model.UserTeam{}).Where("user_id = ? AND team_id = ?", user.ID, team.ID).
		Delete(&model.UserTeam{}).Error
}

// DeleteUserFromContest Many2Many
func DeleteUserFromContest(tx *gorm.DB, user model.User, contest model.Contest) error {
	return tx.Model(&model.UserContest{}).Where("user_id = ? AND contest_id = ?", user.ID, contest.ID).
		Delete(&model.UserContest{}).Error
}
