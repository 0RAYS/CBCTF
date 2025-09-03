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
	BasicRepo[model.User]
}

type CreateUserOptions struct {
	Name           string
	Password       string
	Email          string
	Country        string
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
		Country:        c.Country,
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
	Country        *string
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
	if u.Country != nil {
		options["country"] = *u.Country
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

type DiffUpdateUserOptions struct {
	TeamCount    int64
	ContestCount int64
}

func (d DiffUpdateUserOptions) Convert2Expr() map[string]any {
	options := make(map[string]any)
	if d.TeamCount != 0 {
		options["team_count"] = gorm.Expr("team_count + ?", d.TeamCount)
	}
	if d.ContestCount != 0 {
		options["contest_count"] = gorm.Expr("contest_count + ?", d.ContestCount)
	}
	return options
}

func InitUserRepo(tx *gorm.DB) *UserRepo {
	return &UserRepo{
		BasicRepo: BasicRepo[model.User]{
			DB: tx,
		},
	}
}

func (u *UserRepo) IsUniqueName(name string) bool {
	_, ok, _ := u.GetByUniqueKey("name", name, GetOptions{Selects: []string{"id"}})
	return !ok
}

func (u *UserRepo) IsUniqueEmail(email string) bool {
	_, ok, _ := u.GetByUniqueKey("email", email, GetOptions{Selects: []string{"id"}})
	return !ok
}

func (u *UserRepo) GetByName(name string, optionsL ...GetOptions) (model.User, bool, string) {
	return u.GetByUniqueKey("name", name, optionsL...)
}

func (u *UserRepo) Delete(idL ...uint) (bool, string) {
	userL, _, ok, msg := u.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id", "name", "email", "provider_user_id"},
		Preloads: map[string]GetOptions{
			"Teams":       {Selects: []string{"id", "contest_id"}},
			"Submissions": {Selects: []string{"id", "user_id"}},
		},
	})
	if !ok && msg != i18n.UserNotFound {
		return false, msg
	}
	submissionIDL := make([]uint, 0)
	for _, user := range userL {
		deletedName := fmt.Sprintf("%s_deleted_%s", user.Name, utils.RandStr(6))
		deletedEmail := fmt.Sprintf("%s_deleted_%s", user.Email, utils.RandStr(6))
		deleteProviderUserID := fmt.Sprintf("%s_deleted_%s", user.ProviderUserID, utils.RandStr(6))
		if ok, msg = u.Update(user.ID, UpdateUserOptions{
			Name:           &deletedName,
			Email:          &deletedEmail,
			ProviderUserID: &deleteProviderUserID,
		}); !ok {
			return false, msg
		}
		for _, team := range user.Teams {
			if ok, msg = DeleteUserFromContest(u.DB, user, model.Contest{BasicModel: model.BasicModel{ID: team.ContestID}}); !ok {
				return false, msg
			}
			if ok, msg = DeleteUserFromTeam(u.DB, user, *team); !ok {
				return false, msg
			}
		}
		for _, submission := range user.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
	}
	if ok, msg = InitSubmissionRepo(u.DB).Delete(submissionIDL...); !ok {
		return false, msg
	}
	if res := u.DB.Model(&model.User{}).Where("id IN ?", idL).Delete(&model.User{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete User: %s", res.Error)
		return false, i18n.DeleteUserError
	}
	return true, i18n.Success
}
