package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
	"strings"
)

func CreateUser(tx *gorm.DB, form f.RegisterForm) (model.User, bool, string) {
	if !utils.IsValidEmail(form.Email) {
		return model.User{}, false, i18n.InvalidEmail
	}
	repo := db.InitUserRepo(tx)
	if !repo.IsUniqueName(form.Name) {
		return model.User{}, false, i18n.DuplicateUsername
	}
	if !repo.IsUniqueEmail(form.Email) {
		return model.User{}, false, i18n.DuplicateEmail
	}
	if utils.CheckPassword(form.Password) <= 1 {
		return model.User{}, false, i18n.WeakPassword
	}
	return repo.Create(db.CreateUserOptions{
		Name:     form.Name,
		Password: utils.HashPassword(form.Password),
		Email:    form.Email,
	})
}

func AdminCreateUser(tx *gorm.DB, form f.CreateUserForm) (model.User, bool, string) {
	if !utils.IsValidEmail(form.Email) {
		return model.User{}, false, i18n.InvalidEmail
	}
	repo := db.InitUserRepo(tx)
	if !repo.IsUniqueName(form.Name) {
		return model.User{}, false, i18n.DuplicateUsername
	}
	if !repo.IsUniqueEmail(form.Email) {
		return model.User{}, false, i18n.DuplicateEmail
	}
	return repo.Create(db.CreateUserOptions{
		Name:     form.Name,
		Password: utils.HashPassword(form.Password),
		Email:    form.Email,
		Desc:     form.Desc,
		Country:  strings.ToUpper(form.Country),
		Verified: form.Verified,
		Banned:   form.Banned,
		Hidden:   form.Hidden,
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
	if !utils.CompareHashAndPassword(user.Password, form.OldPassword) {
		return false, i18n.PasswordError
	}
	if utils.CompareHashAndPassword(user.Password, form.NewPassword) {
		return false, i18n.PasswordSame
	}
	if utils.CheckPassword(form.NewPassword) <= 1 {
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
		country := strings.ToUpper(*form.Country)
		options.Country = &country
	}
	if form.Email != nil && *form.Email != user.Email {
		if !utils.IsValidEmail(*form.Email) {
			return false, i18n.InvalidEmail
		}
		if !repo.IsUniqueEmail(*form.Email) {
			return false, i18n.DuplicateEmail
		}
		verified := false
		options.Email = form.Email
		options.Verified = &verified
	}
	if form.Name != nil && *form.Name != user.Name {
		if !repo.IsUniqueName(*form.Name) {
			return false, i18n.DuplicateUsername
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
