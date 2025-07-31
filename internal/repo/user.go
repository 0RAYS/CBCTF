package repo

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
	Name     string
	Password string
	Email    string
	Country  string
	Avatar   model.AvatarURL
	Desc     string
	Verified bool
	Hidden   bool
	Banned   bool
	OauthRaw string
}

func (c CreateUserOptions) Convert2Model() model.Model {
	return model.User{
		Name:     c.Name,
		Password: c.Password,
		Email:    c.Email,
		Country:  c.Country,
		Avatar:   c.Avatar,
		Desc:     c.Desc,
		Verified: c.Verified,
		Hidden:   c.Hidden,
		Banned:   c.Banned,
		OauthRaw: c.OauthRaw,
	}
}

type UpdateUserOptions struct {
	Name     *string
	Password *string
	Email    *string
	Country  *string
	Desc     *string
	Avatar   *model.AvatarURL
	Verified *bool
	Hidden   *bool
	Banned   *bool
	Score    *float64
	Solved   *int64
}

func (u UpdateUserOptions) Convert2Map() map[string]any {
	data := make(map[string]any)
	if u.Name != nil {
		data["name"] = *u.Name
	}
	if u.Password != nil {
		data["password"] = *u.Password
	}
	if u.Email != nil {
		data["email"] = *u.Email
	}
	if u.Country != nil {
		data["country"] = *u.Country
	}
	if u.Desc != nil {
		data["desc"] = *u.Desc
	}
	if u.Avatar != nil {
		data["avatar"] = *u.Avatar
	}
	if u.Verified != nil {
		data["verified"] = *u.Verified
	}
	if u.Hidden != nil {
		data["hidden"] = *u.Hidden
	}
	if u.Banned != nil {
		data["banned"] = *u.Banned
	}
	if u.Score != nil {
		data["score"] = *u.Score
	}
	if u.Solved != nil {
		data["solved"] = *u.Solved
	}
	return data
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
		Selects:    []string{"id", "name", "email"},
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
		if ok, msg = u.Update(user.ID, UpdateUserOptions{
			Name:  &deletedName,
			Email: &deletedEmail,
		}); !ok {
			return false, msg
		}
		for _, team := range user.Teams {
			if ok, msg = DeleteUserFromContest(u.DB, user.ID, team.ContestID); !ok {
				return false, msg
			}
			if ok, msg = DeleteUserFromTeam(u.DB, user.ID, team.ID); !ok {
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
