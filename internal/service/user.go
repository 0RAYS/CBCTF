package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/oauth"
	"CBCTF/internal/utils"

	"gorm.io/gorm"
)

func CreateUser(tx *gorm.DB, form dto.RegisterForm) (model.User, model.RetVal) {
	return db.InitUserRepo(tx).Create(db.CreateUserOptions{
		Name:           form.Name,
		Password:       utils.HashPassword(form.Password),
		Email:          form.Email,
		Provider:       oauth.LocalProvider,
		ProviderUserID: utils.UUID(),
		OauthRaw:       "{}",
	})
}

func AdminCreateUser(tx *gorm.DB, form dto.CreateUserForm) (model.User, model.RetVal) {
	return db.InitUserRepo(tx).Create(db.CreateUserOptions{
		Name:           form.Name,
		Password:       utils.HashPassword(form.Password),
		Email:          form.Email,
		Description:    form.Description,
		Verified:       form.Verified,
		Banned:         form.Banned,
		Hidden:         form.Hidden,
		Provider:       oauth.LocalProvider,
		ProviderUserID: utils.UUID(),
		OauthRaw:       "{}",
	})
}

func VerifyUser(tx *gorm.DB, form dto.LoginForm) (model.User, model.RetVal) {
	repo := db.InitUserRepo(tx)
	user, ret := repo.GetByUniqueKey("name", form.Name)
	if !ret.OK {
		return model.User{}, model.RetVal{Msg: i18n.Model.User.NamePasswordWrong}
	}
	if !utils.CompareHashAndPassword(user.Password, form.Password) {
		return model.User{}, model.RetVal{Msg: i18n.Model.User.NamePasswordWrong}
	}
	return user, model.SuccessRetVal()
}

func ChangeUserPwd(tx *gorm.DB, user model.User, form dto.ChangePasswordForm) model.RetVal {
	repo := db.InitUserRepo(tx)
	if user.Password != model.NeverLoginPWD && !utils.CompareHashAndPassword(user.Password, form.OldPassword) {
		return model.RetVal{Msg: i18n.Model.User.PasswordWrong}
	}
	if utils.CheckPassword(form.NewPassword) < 2 {
		return model.RetVal{Msg: i18n.Model.User.WeakPassword}
	}
	return repo.Update(user.ID, db.UpdateUserOptions{Password: new(utils.HashPassword(form.NewPassword))})
}

func UpdateSelf(tx *gorm.DB, user model.User, form dto.UpdateSelfForm) model.RetVal {
	repo := db.InitUserRepo(tx)
	options := db.UpdateUserOptions{
		Name:        form.Name,
		Email:       form.Email,
		Description: form.Description,
	}
	if form.Email != nil && *form.Email != user.Email {
		options.Verified = new(false)
	}
	return repo.Update(user.ID, options)
}

func UpdateUser(tx *gorm.DB, user model.User, form dto.UpdateUserForm) model.RetVal {
	repo := db.InitUserRepo(tx)
	options := db.UpdateUserOptions{
		Name:        form.Name,
		Email:       form.Email,
		Description: form.Description,
		Hidden:      form.Hidden,
		Banned:      form.Banned,
		Verified:    form.Verified,
	}
	if form.Email != nil && *form.Email != user.Email {
		options.Verified = new(false)
	}
	if form.Password != nil {
		options.Password = new(utils.HashPassword(*form.Password))
	}
	return repo.Update(user.ID, options)
}

func DeleteSelf(tx *gorm.DB, user model.User, form dto.DeleteSelfForm) model.RetVal {
	if !utils.CompareHashAndPassword(user.Password, form.Password) {
		return model.RetVal{Msg: i18n.Model.User.PasswordWrong}
	}
	return DeleteUser(tx, user)
}

func DeleteUser(tx *gorm.DB, user model.User) model.RetVal {
	repo := db.InitUserRepo(tx)
	if repo.CountAssociation(user, "Contests") > 0 {
		return model.RetVal{Msg: i18n.Model.User.InContest}
	}
	return repo.Delete(user.ID)
}
