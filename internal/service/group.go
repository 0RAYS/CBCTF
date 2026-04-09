package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/view"

	"gorm.io/gorm"
)

func BuildGroupView(tx *gorm.DB, group model.Group) view.GroupView {
	count, _ := db.InitGroupRepo(tx).CountUsers(group.ID)
	return view.GroupView{
		Group:     group,
		UserCount: count,
	}
}

func BuildGroupViews(tx *gorm.DB, groups []model.Group) []view.GroupView {
	views := make([]view.GroupView, 0, len(groups))
	for _, group := range groups {
		views = append(views, BuildGroupView(tx, group))
	}
	return views
}

func GetGroupView(tx *gorm.DB, group model.Group) view.GroupView {
	return BuildGroupView(tx, group)
}

func ListGroups(tx *gorm.DB, form dto.ListModelsForm) ([]view.GroupView, int64, model.RetVal) {
	groups, count, ret := db.InitGroupRepo(tx).List(form.Limit, form.Offset)
	if !ret.OK {
		return nil, 0, ret
	}
	return BuildGroupViews(tx, groups), count, model.SuccessRetVal()
}

func ListGroupUsers(tx *gorm.DB, group model.Group, form dto.ListModelsForm) ([]view.UserView, int64, model.RetVal) {
	users, count, ret := db.InitUserRepo(tx).GetByGroupID(group.ID, form.Limit, form.Offset)
	if !ret.OK {
		return nil, 0, ret
	}
	return BuildUserViews(tx, users, true), count, model.SuccessRetVal()
}

func CreateGroup(tx *gorm.DB, form dto.CreateGroupForm) (model.Group, model.RetVal) {
	return db.InitGroupRepo(tx).Create(db.CreateGroupOptions{
		RoleID:      form.RoleID,
		Name:        form.Name,
		Description: form.Description,
	})
}

func UpdateGroup(tx *gorm.DB, group model.Group, form dto.UpdateGroupForm) model.RetVal {
	if group.Default && form.Name != nil {
		return model.RetVal{Msg: i18n.Model.Group.CannotUpdateDefault}
	}
	return db.InitGroupRepo(tx).Update(group.ID, db.UpdateGroupOptions{
		RoleID:      form.RoleID,
		Name:        form.Name,
		Description: form.Description,
	})
}

func DeleteGroup(tx *gorm.DB, group model.Group) model.RetVal {
	if group.Default {
		return model.RetVal{Msg: i18n.Model.Group.CannotDeleteDefault}
	}
	return db.InitGroupRepo(tx).Delete(group.ID)
}

func AssignUserToGroup(tx *gorm.DB, group model.Group, form dto.AssignUserGroupForm) (model.User, model.RetVal) {
	user, ret := db.InitUserRepo(tx).GetByID(form.UserID)
	if !ret.OK {
		return model.User{}, ret
	}
	return user, db.AppendUserToGroup(tx, user, group)
}

func RemoveUserFromGroup(tx *gorm.DB, group model.Group, form dto.AssignUserGroupForm) (model.User, model.RetVal) {
	user, ret := db.InitUserRepo(tx).GetByID(form.UserID)
	if !ret.OK {
		return model.User{}, ret
	}
	return user, db.DeleteUserFromGroup(tx, user, group)
}
