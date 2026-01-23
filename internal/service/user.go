package service

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/oauth"
	"CBCTF/internal/utils"
	"strings"

	"gorm.io/gorm"
)

func CreateUser(tx *gorm.DB, form f.RegisterForm) (model.User, model.RetVal) {
	return db.InitUserRepo(tx).Create(db.CreateUserOptions{
		Name:           form.Name,
		Password:       utils.HashPassword(form.Password),
		Email:          form.Email,
		Provider:       oauth.LocalProvider,
		ProviderUserID: utils.UUID(),
		OauthRaw:       "{}",
	})
}

func AdminCreateUser(tx *gorm.DB, form f.CreateUserForm) (model.User, model.RetVal) {
	return db.InitUserRepo(tx).Create(db.CreateUserOptions{
		Name:           form.Name,
		Password:       utils.HashPassword(form.Password),
		Email:          form.Email,
		Desc:           form.Desc,
		Country:        strings.ToUpper(form.Country),
		Verified:       form.Verified,
		Banned:         form.Banned,
		Hidden:         form.Hidden,
		Provider:       oauth.LocalProvider,
		ProviderUserID: utils.UUID(),
		OauthRaw:       "{}",
	})
}

func VerifyUser(tx *gorm.DB, form f.LoginForm) (model.User, model.RetVal) {
	repo := db.InitUserRepo(tx)
	user, ret := repo.GetByUniqueKey("name", form.Name)
	if !ret.OK {
		return model.User{}, ret
	}
	if !utils.CompareHashAndPassword(user.Password, form.Password) {
		return model.User{}, model.RetVal{Msg: i18n.Model.User.NamePasswordWrong}
	}
	return user, model.SuccessRetVal()
}

func ChangeUserPwd(tx *gorm.DB, user model.User, form f.ChangePasswordForm) model.RetVal {
	repo := db.InitUserRepo(tx)
	if user.Password != model.NeverLoginPWD && !utils.CompareHashAndPassword(user.Password, form.OldPassword) {
		return model.RetVal{Msg: i18n.Model.User.PasswordWrong}
	}
	if utils.CheckPassword(form.NewPassword) < 2 {
		return model.RetVal{Msg: i18n.Model.User.WeakPassword}
	}
	password := utils.HashPassword(form.NewPassword)
	return repo.Update(user.ID, db.UpdateUserOptions{Password: &password})
}

func UpdateSelf(tx *gorm.DB, user model.User, form f.UpdateSelfForm) model.RetVal {
	repo := db.InitUserRepo(tx)
	options := db.UpdateUserOptions{
		Desc: form.Desc,
	}
	if form.Country != nil && *form.Country != user.Country {
		options.Country = utils.Ptr(strings.ToUpper(*form.Country))
	}
	if form.Email != nil && *form.Email != user.Email {
		options.Verified = utils.Ptr(false)
	}
	return repo.Update(user.ID, options)
}

func DeleteSelf(tx *gorm.DB, user model.User, form f.DeleteSelfForm) model.RetVal {
	if !utils.CompareHashAndPassword(user.Password, form.Password) {
		return model.RetVal{Msg: i18n.Model.User.PasswordWrong}
	}
	contestIDL, ret := db.GetContestIDByUserID(tx, user.ID)
	if !ret.OK {
		return ret
	}
	repo := db.InitContestRepo(tx)
	for _, id := range contestIDL {
		contest, ret := repo.GetByID(id, db.GetOptions{Selects: []string{"id", "start", "duration"}})
		if !ret.OK {
			return ret
		}
		if contest.IsRunning() {
			return model.RetVal{Msg: i18n.Model.Contest.IsRunning}
		}
	}
	return db.InitUserRepo(tx).Delete(user.ID)
}
