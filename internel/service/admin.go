package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

func CreateAdmin(tx *gorm.DB, form f.CreateAdminForm) (model.Admin, bool, string) {
	repo := db.InitAdminRepo(tx)
	if !utils.IsValidEmail(form.Email) {
		return model.Admin{}, false, "InvalidEmail"
	}
	if !repo.IsUniqueEmail(form.Email) {
		return model.Admin{}, false, "DuplicateEmail"
	}
	if !repo.IsUniqueName(form.Name) {
		return model.Admin{}, false, "DuplicateUsername"
	}
	return repo.Create(db.CreateAdminOptions{
		Name:     form.Name,
		Password: form.Email,
		Email:    form.Email,
		Avatar:   "",
		Verified: false,
	})
}

func VerifyAdmin(tx *gorm.DB, form f.LoginForm) (model.Admin, bool, string) {
	repo := db.InitAdminRepo(tx)
	admin, ok, msg := repo.GetByName(form.Name, true, 0)
	if !ok {
		return model.Admin{}, false, msg
	}
	if !utils.CompareHashAndPassword(admin.Password, form.Password) {
		return model.Admin{}, false, "NameOrPasswordError"
	}
	return admin, true, "Success"
}

func UpdateUser(tx *gorm.DB, user model.User, form f.UpdateUserForm) (bool, string) {
	repo := db.InitUserRepo(tx)
	options := db.UpdateUserOptions{
		Desc:     form.Desc,
		Country:  form.Country,
		Hidden:   form.Hidden,
		Banned:   form.Banned,
		Verified: form.Verified,
	}
	if *form.Email != "" && *form.Email != user.Email {
		if !utils.IsValidEmail(*form.Email) {
			return false, "InvalidEmail"
		}
		if !repo.IsUniqueEmail(*form.Email) {
			return false, "DuplicateEmail"
		}
		verified := false
		options.Email = form.Email
		options.Verified = &verified
	}
	if *form.Name != "" && *form.Name != user.Name {
		if !repo.IsUniqueName(*form.Name) {
			return false, "DuplicateUsername"
		}
		options.Name = form.Name
	}
	if *form.Password != "" {
		password := utils.HashPassword(*form.Password)
		options.Password = &password
	}
	return repo.Update(user.ID, options)
}

func DeleteUser(tx *gorm.DB, user model.User) (bool, string) {
	repo := db.InitUserRepo(tx)
	return repo.Delete(user.ID)
}

func UpdateAdmin(tx *gorm.DB, admin model.Admin, form f.UpdateAdminForm) (bool, string) {
	repo := db.InitAdminRepo(tx)
	options := db.UpdateAdminOptions{}
	if form.Email != nil && *form.Email != admin.Email {
		if repo.IsUniqueEmail(*form.Email) {
			return false, "DuplicateEmail"
		}
		options.Email = form.Email
		verified := false
		options.Verified = &verified
	}
	if form.Name != nil && *form.Name != admin.Name {
		if !repo.IsUniqueName(*form.Name) {
			return false, "DuplicateUsername"
		}
		options.Name = form.Name
	}
	return repo.Update(admin.ID, options)
}

func ChangeAdminPassword(tx *gorm.DB, admin model.Admin, form f.ChangePasswordForm) (bool, string) {
	if !utils.CompareHashAndPassword(admin.Password, form.OldPassword) {
		return false, "PasswordError"
	}
	if utils.CompareHashAndPassword(admin.Password, form.NewPassword) {
		return false, "PasswordSame"
	}
	hash := utils.HashPassword(form.NewPassword)
	repo := db.InitAdminRepo(tx)
	return repo.Update(admin.ID, db.UpdateAdminOptions{Password: &hash})
}
