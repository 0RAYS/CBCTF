package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/oauth"
	"CBCTF/internal/utils"
	"fmt"

	"gorm.io/gorm"
)

type UserRepo struct {
	BaseRepo[model.User]
}

type CreateUserOptions struct {
	Name           string
	Password       string
	Email          string
	Picture        model.FileURL
	Description    string
	Verified       bool
	Hidden         bool
	Banned         bool
	Provider       string
	ProviderUserID string
	OauthRaw       string
}

func (c CreateUserOptions) Convert2Model() model.Model {
	return model.User{
		Name:           c.Name,
		Password:       c.Password,
		Email:          c.Email,
		Picture:        c.Picture,
		Description:    c.Description,
		Verified:       c.Verified,
		Hidden:         c.Hidden,
		Banned:         c.Banned,
		Provider:       c.Provider,
		ProviderUserID: c.ProviderUserID,
		OauthRaw:       c.OauthRaw,
	}
}

type UpdateUserOptions struct {
	Name           *string
	Password       *string
	Email          *string
	Description    *string
	Picture        *model.FileURL
	Verified       *bool
	Hidden         *bool
	Banned         *bool
	Score          *float64
	Solved         *int64
	ContestCount   *int64
	TeamCount      *int64
	ProviderUserID *string
	OauthRaw       *string
}

func (u UpdateUserOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Password != nil {
		options["password"] = *u.Password
	}
	if u.Email != nil {
		options["email"] = *u.Email
	}
	if u.Description != nil {
		options["description"] = *u.Description
	}
	if u.Picture != nil {
		options["picture"] = *u.Picture
	}
	if u.Verified != nil {
		options["verified"] = *u.Verified
	}
	if u.Hidden != nil {
		options["hidden"] = *u.Hidden
	}
	if u.Banned != nil {
		options["banned"] = *u.Banned
	}
	if u.Score != nil {
		options["score"] = *u.Score
	}
	if u.Solved != nil {
		options["solved"] = *u.Solved
	}
	if u.ContestCount != nil {
		options["contest_count"] = *u.ContestCount
	}
	if u.TeamCount != nil {
		options["team_count"] = *u.TeamCount
	}
	if u.ProviderUserID != nil {
		options["provider_user_id"] = *u.ProviderUserID
	}
	if u.OauthRaw != nil {
		options["oauth_raw"] = *u.OauthRaw
	}
	return options
}

func InitUserRepo(tx *gorm.DB) *UserRepo {
	return &UserRepo{
		BaseRepo: BaseRepo[model.User]{
			DB: tx,
		},
	}
}

func (u *UserRepo) InitAdmin() model.RetVal {
	count, ret := u.CountGroupUser(model.AdminGroupName)
	if !ret.OK {
		return ret
	}
	if count == 0 {
		pwd := utils.UUID()
		admin, ret := u.Insert(model.User{
			Name:           "admin",
			Password:       utils.HashPassword(pwd),
			Email:          "admin@0rays.club",
			Description:    "default administrator",
			Verified:       true,
			Hidden:         true,
			Banned:         false,
			Provider:       oauth.LocalProvider,
			ProviderUserID: utils.UUID(),
			OauthRaw:       "{}",
		})
		if !ret.OK {
			return ret
		}
		group, ret := InitGroupRepo(u.DB).GetByUniqueField("name", model.AdminGroupName)
		if !ret.OK {
			return ret
		}
		if ret = AppendUserToGroup(u.DB, admin, group); !ret.OK {
			return ret
		}
		log.Logger.Infof("Init Admin: Admin{ name: admin, password: %s, email: admin@0rays.club}", pwd)
	}
	return model.SuccessRetVal()
}

func (u *UserRepo) IsInGroup(userID uint, groupName string) bool {
	var count int64
	res := u.DB.Table("user_groups").
		Joins("INNER JOIN `groups` ON user_groups.group_id = groups.id").
		Where("user_groups.user_id = ? AND groups.name = ? AND groups.deleted_at IS NULL", userID, groupName).
		Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to check user group membership: %v", res.Error)
		return false
	}
	return count > 0
}

func (u *UserRepo) CountGroupUser(group string) (int64, model.RetVal) {
	var count int64
	res := u.DB.Table("users").
		Joins("INNER JOIN user_groups ON users.id = user_groups.user_id").
		Joins("INNER JOIN `groups` ON user_groups.group_id = groups.id").
		Where("groups.name = ? AND users.deleted_at IS NULL AND groups.deleted_at IS NULL", group).
		Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count role users: %v", res.Error)
		return 0, model.RetVal{Msg: i18n.Model.User.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return count, model.SuccessRetVal()
}

func (u *UserRepo) GetByName(name string, optionsL ...GetOptions) (model.User, model.RetVal) {
	return u.GetByUniqueField("name", name, optionsL...)
}

func (u *UserRepo) GetByTeamID(teamID uint, limit, offset int) ([]model.User, model.RetVal) {
	var users []model.User
	res := u.DB.Table("users").Select("users.*").
		Joins("INNER JOIN user_teams ON user_teams.user_id = users.id").
		Where("user_teams.team_id = ? AND users.deleted_at IS NULL", teamID).
		Limit(limit).Offset(offset).Scan(&users)
	if res.Error != nil {
		log.Logger.Fatalf("Failed to get Users: %v", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.User.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return users, model.SuccessRetVal()
}

func (u *UserRepo) GetIDByTeamID(teamID uint, limit, offset int) ([]uint, model.RetVal) {
	users, ret := u.GetByTeamID(teamID, limit, offset)
	if !ret.OK {
		return nil, ret
	}
	var userIDL []uint
	for _, user := range users {
		userIDL = append(userIDL, user.ID)
	}
	return userIDL, ret
}

func (u *UserRepo) GetByContestID(contestID uint, limit, offset int) ([]model.User, model.RetVal) {
	var users []model.User
	res := u.DB.Table("users").Select("users.*").
		Joins("INNER JOIN user_contests ON user_contests.user_id = users.id").
		Where("user_contests.contest_id = ? AND users.deleted_at IS NULL", contestID).
		Limit(limit).Offset(offset).Scan(&users)
	if res.Error != nil {
		log.Logger.Fatalf("Failed to get Users: %v", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.User.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return users, model.SuccessRetVal()
}

func (u *UserRepo) GetIDByContestID(contestID uint, limit, offset int) ([]uint, model.RetVal) {
	users, ret := u.GetByContestID(contestID, limit, offset)
	if !ret.OK {
		return nil, ret
	}
	var userIDL []uint
	for _, user := range users {
		userIDL = append(userIDL, user.ID)
	}
	return userIDL, ret
}

func (u *UserRepo) GetByGroupID(groupID uint, limit, offset int) ([]model.User, int64, model.RetVal) {
	var count int64
	if res := u.DB.Model(&model.UserGroup{}).Where("group_id = ?", groupID).Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to count Group Users: %s", res.Error)
		return nil, 0, model.RetVal{Msg: i18n.Model.UserGroup.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	var users []model.User
	res := u.DB.Table("users").Select("users.*").
		Joins("INNER JOIN user_groups ON user_groups.user_id = users.id").
		Where("user_groups.group_id = ? AND users.deleted_at IS NULL", groupID).
		Limit(limit).Offset(offset).Scan(&users)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Group Users: %s", res.Error)
		return nil, 0, model.RetVal{Msg: i18n.Model.User.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return users, count, model.SuccessRetVal()
}

func (u *UserRepo) Delete(idL ...uint) model.RetVal {
	userL, _, ret := u.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Preloads:   map[string]GetOptions{"Teams": {}, "Submissions": {}},
	})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	submissionIDL := make([]uint, 0)
	for _, user := range userL {
		if ret = u.Update(user.ID, UpdateUserOptions{
			Name:           new(fmt.Sprintf("%s_deleted_%s", user.Name, utils.RandStr(6))),
			Email:          new(fmt.Sprintf("%s_deleted_%s", user.Email, utils.RandStr(6))),
			ProviderUserID: new(fmt.Sprintf("%s_deleted_%s", user.ProviderUserID, utils.RandStr(6))),
		}); !ret.OK {
			return ret
		}
		for _, team := range user.Teams {
			if ret = DeleteUserFromContest(u.DB, user, model.Contest{BaseModel: model.BaseModel{ID: team.ContestID}}); !ret.OK {
				return ret
			}
			if ret = DeleteUserFromTeam(u.DB, user, team); !ret.OK {
				return ret
			}
		}
		for _, submission := range user.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
	}
	if ret = InitSubmissionRepo(u.DB).Delete(submissionIDL...); !ret.OK {
		return ret
	}
	if res := u.DB.Model(&model.User{}).Where("id IN ?", idL).Delete(&model.User{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete User: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.User.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
