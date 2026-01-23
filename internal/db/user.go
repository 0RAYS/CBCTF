package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
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
	Avatar         model.AvatarURL
	Desc           string
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
		Avatar:         c.Avatar,
		Desc:           c.Desc,
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
	Desc           *string
	Avatar         *model.AvatarURL
	Verified       *bool
	Hidden         *bool
	Banned         *bool
	Score          *float64
	Solved         *int64
	ContestCount   *int64
	TeamCount      *int64
	ProviderUserID *string
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
	if u.Desc != nil {
		options["desc"] = *u.Desc
	}
	if u.Avatar != nil {
		options["avatar"] = *u.Avatar
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
	return options
}

func InitUserRepo(tx *gorm.DB) *UserRepo {
	return &UserRepo{
		BaseRepo: BaseRepo[model.User]{
			DB: tx,
		},
	}
}

func (u *UserRepo) GetByName(name string, optionsL ...GetOptions) (model.User, model.RetVal) {
	return u.GetByUniqueKey("name", name, optionsL...)
}

func (u *UserRepo) Delete(idL ...uint) model.RetVal {
	userL, _, ret := u.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id", "name", "email", "provider_user_id"},
		Preloads: map[string]GetOptions{
			"Teams":       {Selects: []string{"id", "contest_id"}},
			"Submissions": {Selects: []string{"id", "user_id"}},
		},
	})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	submissionIDL := make([]uint, 0)
	for _, user := range userL {
		deletedName := fmt.Sprintf("%s_deleted_%s", user.Name, utils.RandStr(6))
		deletedEmail := fmt.Sprintf("%s_deleted_%s", user.Email, utils.RandStr(6))
		deleteProviderUserID := fmt.Sprintf("%s_deleted_%s", user.ProviderUserID, utils.RandStr(6))
		if ret = u.Update(user.ID, UpdateUserOptions{
			Name:           &deletedName,
			Email:          &deletedEmail,
			ProviderUserID: &deleteProviderUserID,
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
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": model.User{}.GetModelName(), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
