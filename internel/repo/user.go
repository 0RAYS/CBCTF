package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"fmt"
	"gorm.io/gorm"
)

type UserRepo struct {
	Basic[model.User]
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
		Basic: Basic[model.User]{
			DB: tx,
		},
	}
}

func (u *UserRepo) IsUniqueName(name string) bool {
	_, ok, _ := u.getUniqueByKey("name", name)
	return !ok
}

func (u *UserRepo) IsUniqueEmail(email string) bool {
	_, ok, _ := u.getUniqueByKey("email", email)
	return !ok
}

func (u *UserRepo) GetByName(name string, preloadL ...string) (model.User, bool, string) {
	return u.getUniqueByKey("name", name, preloadL...)
}

func (u *UserRepo) Delete(idL ...uint) (bool, string) {
	submissionIDL := make([]uint, 0)
	for _, id := range idL {
		user, ok, msg := u.GetByID(id, "Teams", "Submissions")
		if !ok {
			return false, msg
		}
		deletedName := fmt.Sprintf("%s_deleted_%s", user.Name, utils.RandStr(6))
		deletedEmail := fmt.Sprintf("%s_deleted_%s", user.Email, utils.RandStr(6))
		if ok, msg = u.Update(id, UpdateUserOptions{
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
	if ok, msg := InitSubmissionRepo(u.DB).Delete(submissionIDL...); !ok {
		return false, msg
	}
	if res := u.DB.Model(&model.User{}).Where("id IN ?", idL).Delete(&model.User{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete User: %s", res.Error)
		return false, model.User{}.DeleteErrorString()
	}
	return true, i18n.Success
}
