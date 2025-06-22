package service

import (
	"CBCTF/internel/email"
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
	"strings"
)

func CreateAdmin(tx *gorm.DB, form f.CreateAdminForm) (model.Admin, bool, string) {
	repo := db.InitAdminRepo(tx)
	if !email.IsValidEmail(form.Email) {
		return model.Admin{}, false, i18n.InvalidEmail
	}
	if !repo.IsUniqueEmail(form.Email) {
		return model.Admin{}, false, i18n.DuplicateEmail
	}
	if !repo.IsUniqueName(form.Name) {
		return model.Admin{}, false, i18n.DuplicateUserName
	}
	return repo.Create(db.CreateAdminOptions{
		Name:     form.Name,
		Password: utils.HashPassword(form.Password),
		Email:    form.Email,
		Avatar:   "",
		Verified: false,
	})
}

func VerifyAdmin(tx *gorm.DB, form f.LoginForm) (model.Admin, bool, string) {
	repo := db.InitAdminRepo(tx)
	admin, ok, msg := repo.GetByName(form.Name)
	if !ok {
		return model.Admin{}, false, msg
	}
	if !utils.CompareHashAndPassword(admin.Password, form.Password) {
		return model.Admin{}, false, i18n.NameOrPasswordError
	}
	return admin, true, i18n.Success
}

func UpdateUser(tx *gorm.DB, user model.User, form f.UpdateUserForm) (bool, string) {
	repo := db.InitUserRepo(tx)
	options := db.UpdateUserOptions{
		Desc:     form.Desc,
		Hidden:   form.Hidden,
		Banned:   form.Banned,
		Verified: form.Verified,
	}
	if form.Country != nil && *form.Country != user.Country {
		options.Country = utils.Ptr(strings.ToUpper(*form.Country))
	}
	if form.Email != nil && *form.Email != user.Email {
		if !email.IsValidEmail(*form.Email) {
			return false, i18n.InvalidEmail
		}
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
	if form.Password != nil {
		options.Password = utils.Ptr(utils.HashPassword(*form.Password))
	}
	return repo.Update(user.ID, options)
}

func UpdateAdmin(tx *gorm.DB, admin model.Admin, form f.UpdateAdminForm) (bool, string) {
	repo := db.InitAdminRepo(tx)
	options := db.UpdateAdminOptions{}
	if form.Email != nil && *form.Email != admin.Email {
		if !repo.IsUniqueEmail(*form.Email) {
			return false, i18n.DuplicateEmail
		}
		options.Email = form.Email
		options.Verified = utils.Ptr(false)
	}
	if form.Name != nil && *form.Name != admin.Name {
		if !repo.IsUniqueName(*form.Name) {
			return false, i18n.DuplicateUserName
		}
		options.Name = form.Name
	}
	return repo.Update(admin.ID, options)
}

func ChangeAdminPassword(tx *gorm.DB, admin model.Admin, form f.ChangePasswordForm) (bool, string) {
	if !utils.CompareHashAndPassword(admin.Password, form.OldPassword) {
		return false, i18n.PasswordError
	}
	if utils.CompareHashAndPassword(admin.Password, form.NewPassword) {
		return false, i18n.PasswordSame
	}
	hash := utils.HashPassword(form.NewPassword)
	repo := db.InitAdminRepo(tx)
	return repo.Update(admin.ID, db.UpdateAdminOptions{Password: &hash})
}
