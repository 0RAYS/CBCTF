package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func GetContestIDByUserID(tx *gorm.DB, userID uint) ([]uint, model.RetVal) {
	var ucs []model.UserContest
	res := tx.Model(&model.UserContest{}).Where("user_id = ?", userID).Find(&ucs)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get contest: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.Contest{}.ModelName(), "Error": res.Error.Error()}}
	}
	var idL []uint
	for _, uc := range ucs {
		idL = append(idL, uc.ContestID)
	}
	return idL, model.SuccessRetVal()
}

// AppendUserToTeam Many2Many
func AppendUserToTeam(tx *gorm.DB, user model.User, team model.Team) model.RetVal {
	res := tx.Model(&model.UserTeam{}).Create(&model.UserTeam{UserID: user.ID, TeamID: team.ID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append User to Team: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": "UserTeam", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

// AppendUserToContest Many2Many
func AppendUserToContest(tx *gorm.DB, user model.User, contest model.Contest) model.RetVal {
	res := tx.Model(&model.UserContest{}).Create(&model.UserContest{UserID: user.ID, ContestID: contest.ID})
	if res.Error != nil {
		log.Logger.Warningf("Failed to append User to Contest: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": "UserContest", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

// DeleteUserFromTeam Many2Many
func DeleteUserFromTeam(tx *gorm.DB, user model.User, team model.Team) model.RetVal {
	res := tx.Model(&model.UserTeam{}).Where("user_id = ? AND team_id = ?", user.ID, team.ID).
		Delete(&model.UserTeam{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete User from Team: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": "UserTeam", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

// DeleteUserFromContest Many2Many
func DeleteUserFromContest(tx *gorm.DB, user model.User, contest model.Contest) model.RetVal {
	res := tx.Model(&model.UserContest{}).Where("user_id = ? AND contest_id = ?", user.ID, contest.ID).
		Delete(&model.UserContest{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete User from Contest: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": "UserContest", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
