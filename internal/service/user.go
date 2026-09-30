package service

import (
	"gorm.io/gorm"

	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/oa"
	"CBCTF/internal/utils"
	"CBCTF/internal/view"
)

func BuildUserView(tx *gorm.DB, user model.User, includeCounts bool) view.UserView {
	result := view.UserView{
		User:           user,
		HasAdminAccess: db.InitPermissionRepo(tx).HasAdminAccess(user.ID),
	}
	if includeCounts {
		var ret model.RetVal
		result.TeamCount, ret = db.InitUserRepo(tx).CountTeams(user.ID)
		if !ret.OK {
			result.Unavailable = append(result.Unavailable, "teams")
		}
		result.ContestCount, ret = db.InitUserRepo(tx).CountContests(user.ID)
		if !ret.OK {
			result.Unavailable = append(result.Unavailable, "contests")
		}
	}
	return result
}

func BuildUserViews(tx *gorm.DB, users []model.User, includeCounts bool) []view.UserView {
	views := make([]view.UserView, 0, len(users))
	for _, user := range users {
		views = append(views, BuildUserView(tx, user, includeCounts))
	}
	return views
}

func RegisterUser(tx *gorm.DB, form dto.RegisterForm) (model.User, model.RetVal) {
	var result model.User
	ret := db.WithTransactionDB(tx, func(tx *gorm.DB) model.RetVal {
		user, ret := db.InitUserRepo(tx).Create(model.User{
			Name:           form.Name,
			Password:       utils.HashPassword(form.Password),
			Email:          form.Email,
			Provider:       oa.LocalProvider,
			ProviderUserID: utils.UUID(),
			OauthRaw:       "{}",
		})
		if !ret.OK {
			return ret
		}
		if config.Env.Registration.DefaultGroup != 0 {
			defaultGroup, groupRet := db.InitGroupRepo(tx).GetByID(config.Env.Registration.DefaultGroup)
			if groupRet.OK {
				if ret = db.AppendUserToGroup(tx, user, defaultGroup); !ret.OK {
					return ret
				}
			}
		}
		result = user
		return model.SuccessRetVal()
	})
	return result, ret
}

func AdminCreateUser(tx *gorm.DB, form dto.CreateUserForm) (model.User, model.RetVal) {
	return db.InitUserRepo(tx).Create(model.User{
		Name:           form.Name,
		Password:       utils.HashPassword(form.Password),
		Email:          form.Email,
		Description:    form.Description,
		Verified:       form.Verified,
		Banned:         form.Banned,
		Hidden:         form.Hidden,
		Provider:       oa.LocalProvider,
		ProviderUserID: utils.UUID(),
		OauthRaw:       "{}",
	})
}

func VerifyUser(tx *gorm.DB, form dto.LoginForm) (model.User, model.RetVal) {
	repo := db.InitUserRepo(tx)
	user, ret := repo.GetByUniqueField("name", form.Name)
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
	return db.WithTransactionDB(tx, func(tx *gorm.DB) model.RetVal {
		repo := db.InitUserRepo(tx)
		count, ret := repo.CountContests(user.ID)
		if !ret.OK {
			return ret
		}
		if count > 0 {
			return model.RetVal{Msg: i18n.Model.User.InContest}
		}
		return repo.Delete(user.ID)
	})
}

func ListUsers(tx *gorm.DB, form dto.ListUsersForm) ([]view.UserView, int64, model.RetVal) {
	options := db.GetOptions{
		Search: userSearchConditions(form),
		Sort:   []string{"id DESC"},
	}
	users, count, ret := db.InitUserRepo(tx).List(form.Limit, form.Offset, options)
	if !ret.OK {
		return nil, 0, ret
	}
	return BuildUserViews(tx, users, true), count, model.SuccessRetVal()
}

func ListUsersNotInGroup(tx *gorm.DB, group model.Group, form dto.ListUsersForm) ([]view.UserView, int64, model.RetVal) {
	options := db.GetOptions{Search: userSearchConditions(form)}
	users, count, ret := db.InitUserRepo(tx).GetNotInGroupID(group.ID, form.Limit, form.Offset, options)
	if !ret.OK {
		return nil, 0, ret
	}
	return BuildUserViews(tx, users, true), count, model.SuccessRetVal()
}

func userSearchConditions(form dto.ListUsersForm) map[string]string {
	conditions := make(map[string]string)
	if form.Name != "" {
		conditions["name"] = form.Name
	}
	if form.Email != "" {
		conditions["email"] = form.Email
	}
	if form.Description != "" {
		conditions["description"] = form.Description
	}
	return conditions
}

func GetAccessibleRoutes(tx *gorm.DB, userID uint) ([]string, model.RetVal) {
	permNames, ret := db.InitPermissionRepo(tx).GetUserPermissions(userID)
	if !ret.OK {
		return nil, ret
	}

	permSet := make(map[string]struct{}, len(permNames))
	for _, name := range permNames {
		permSet[name] = struct{}{}
	}

	routes := make([]string, 0)
	for route, perm := range model.RoutePermissions {
		if _, ok := permSet[perm]; ok {
			routes = append(routes, route)
		}
	}
	return routes, model.SuccessRetVal()
}
