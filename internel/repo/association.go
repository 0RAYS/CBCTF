package repo

import (
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type AssociationRepo struct {
	DB *gorm.DB
}

func (a *AssociationRepo) IsUniqueMember(contestID, userID uint) bool {
	res := a.DB.Model(&model.UserContest{}).
		Where("contest_id = ? AND user_id = ?", contestID, userID).Find(&model.UserContest{}).Limit(1)
	return res.RowsAffected == 0
}

func (a *AssociationRepo) AppendUserToTeam(userID, teamID uint) bool {
	return a.DB.Model(&model.UserTeam{}).Create(&model.UserTeam{UserID: userID, TeamID: teamID}).Error == nil
}

func (a *AssociationRepo) AppendUserToContest(userID, contestID uint) bool {
	return a.DB.Model(&model.UserContest{}).Create(&model.UserContest{UserID: userID, ContestID: contestID}).Error == nil
}

func (a *AssociationRepo) DeleteUserFromTeam(userID, teamID uint) bool {
	return a.DB.Model(&model.UserTeam{}).Where("user_id = ? AND team_id = ?", userID, teamID).
		Delete(&model.UserTeam{}).Error == nil
}

func (a *AssociationRepo) DeleteUserFromContest(userID, contestID uint) bool {
	return a.DB.Model(&model.UserContest{}).Where("user_id = ? AND contest_id = ?", userID, contestID).
		Delete(&model.UserContest{}).Error == nil
}
