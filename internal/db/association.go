package db

import (
	"CBCTF/internal/model"
	"gorm.io/gorm"
)

// AppendUserToTeam Many2Many
func AppendUserToTeam(tx *gorm.DB, user model.User, team model.Team) error {
	return tx.Model(&team).Association("Users").Append(&user)
}

// AppendUserToContest Many2Many
func AppendUserToContest(tx *gorm.DB, user model.User, contest model.Contest) error {
	return tx.Model(&contest).Association("Users").Append(&user)
}

// AppendTeamToContest HasMany
func AppendTeamToContest(tx *gorm.DB, team model.Team, contest model.Contest) error {
	return tx.Model(&contest).Association("Teams").Append(&team)
}

// DeleteUserFromTeam Many2Many
func DeleteUserFromTeam(tx *gorm.DB, user model.User, team model.Team) error {
	return tx.Model(&team).Association("Users").Delete(&user)
}

// DeleteUserFromContest Many2Many
func DeleteUserFromContest(tx *gorm.DB, user model.User, contest model.Contest) error {
	return tx.Model(&contest).Association("Users").Delete(&user)
}

// DeleteTeamFromContest HasMany
func DeleteTeamFromContest(tx *gorm.DB, team model.Team, contest model.Contest) error {
	return tx.Model(&contest).Association("Teams").Delete(&team)
}
