package service

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/oauth"
	db "CBCTF/internal/repo"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"strings"
)

func CreateUser(tx *gorm.DB, form f.RegisterForm) (model.User, bool, string) {
	repo := db.InitUserRepo(tx)
	if !repo.IsUniqueName(form.Name) {
		return model.User{}, false, i18n.DuplicateUserName
	}
	if !repo.IsUniqueEmail(form.Email) {
		return model.User{}, false, i18n.DuplicateEmail
	}
	return repo.Create(db.CreateUserOptions{
		Name:           form.Name,
		Password:       utils.HashPassword(form.Password),
		Email:          form.Email,
		Provider:       oauth.LocalProvider,
		ProviderUserID: utils.UUID(),
		OauthRaw:       "{}",
	})
}

func AdminCreateUser(tx *gorm.DB, form f.CreateUserForm) (model.User, bool, string) {
	repo := db.InitUserRepo(tx)
	if !repo.IsUniqueName(form.Name) {
		return model.User{}, false, i18n.DuplicateUserName
	}
	if !repo.IsUniqueEmail(form.Email) {
		return model.User{}, false, i18n.DuplicateEmail
	}
	return repo.Create(db.CreateUserOptions{
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

func VerifyUser(tx *gorm.DB, form f.LoginForm) (model.User, bool, string) {
	repo := db.InitUserRepo(tx)
	user, ok, msg := repo.GetByName(form.Name)
	if !ok {
		return model.User{}, false, msg
	}
	if !utils.CompareHashAndPassword(user.Password, form.Password) {
		return model.User{}, false, i18n.NameOrPasswordError
	}
	return user, true, i18n.Success
}

func ChangeUserPwd(tx *gorm.DB, user model.User, form f.ChangePasswordForm) (bool, string) {
	repo := db.InitUserRepo(tx)
	if user.Password != model.NeverLoginPWD && !utils.CompareHashAndPassword(user.Password, form.OldPassword) {
		return false, i18n.PasswordError
	}
	if utils.CheckPassword(form.NewPassword) < 2 {
		return false, i18n.WeakPassword
	}
	password := utils.HashPassword(form.NewPassword)
	return repo.Update(user.ID, db.UpdateUserOptions{Password: &password})
}

func UpdateSelf(tx *gorm.DB, user model.User, form f.UpdateSelfForm) (bool, string) {
	repo := db.InitUserRepo(tx)
	options := db.UpdateUserOptions{
		Desc: form.Desc,
	}
	if form.Country != nil && *form.Country != user.Country {
		options.Country = utils.Ptr(strings.ToUpper(*form.Country))
	}
	if form.Email != nil && *form.Email != user.Email {
		if !repo.IsUniqueEmail(*form.Email) {
			return false, i18n.DuplicateEmail
		}
		options.Email = form.Email
		options.Verified = utils.Ptr(false)
	}
	if form.Name != nil && *form.Name != user.Name {
		if !repo.IsUniqueName(*form.Name) {
			return false, i18n.DuplicateUserName
		}
		options.Name = form.Name
	}
	return repo.Update(user.ID, options)
}

func DeleteSelf(tx *gorm.DB, user model.User, form f.DeleteSelfForm) (bool, string) {
	repo := db.InitUserRepo(tx)
	if !utils.CompareHashAndPassword(user.Password, form.Password) {
		return false, i18n.PasswordError
	}
	return repo.Delete(user.ID)
}
