package service

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"strings"

	"gorm.io/gorm"
)

func CreateAdmin(tx *gorm.DB, form f.CreateAdminForm) (model.Admin, model.RetVal) {
	return db.InitAdminRepo(tx).Create(db.CreateAdminOptions{
		Name:     form.Name,
		Password: utils.HashPassword(form.Password),
		Email:    form.Email,
		Avatar:   "",
		Verified: false,
	})
}

func VerifyAdmin(tx *gorm.DB, form f.LoginForm) (model.Admin, model.RetVal) {
	repo := db.InitAdminRepo(tx)
	admin, ret := repo.GetByUniqueKey("name", form.Name)
	if !ret.OK {
		return model.Admin{}, ret
	}
	if !utils.CompareHashAndPassword(admin.Password, form.Password) {
		return model.Admin{}, model.RetVal{Msg: i18n.Model.User.NamePasswordWrong}
	}
	return admin, model.SuccessRetVal()
}

func UpdateUser(tx *gorm.DB, user model.User, form f.UpdateUserForm) model.RetVal {
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
		options.Verified = utils.Ptr(false)
	}
	if form.Password != nil {
		options.Password = utils.Ptr(utils.HashPassword(*form.Password))
	}
	return repo.Update(user.ID, options)
}

func UpdateAdmin(tx *gorm.DB, admin model.Admin, form f.UpdateAdminForm) model.RetVal {
	repo := db.InitAdminRepo(tx)
	options := db.UpdateAdminOptions{}
	if form.Email != nil && *form.Email != admin.Email {
		options.Verified = utils.Ptr(false)
	}
	return repo.Update(admin.ID, options)
}

func ChangeAdminPassword(tx *gorm.DB, admin model.Admin, form f.ChangePasswordForm) model.RetVal {
	if !utils.CompareHashAndPassword(admin.Password, form.OldPassword) {
		return model.RetVal{Msg: i18n.Model.User.PasswordWrong}
	}
	hash := utils.HashPassword(form.NewPassword)
	return db.InitAdminRepo(tx).Update(admin.ID, db.UpdateAdminOptions{Password: &hash})
}
