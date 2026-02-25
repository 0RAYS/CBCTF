package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/utils"
	"encoding/json"
	"fmt"
	"slices"

	"gorm.io/gorm"
)

func OauthLogin(tx *gorm.DB, provider model.Oauth, response map[string]any) (model.User, model.RetVal) {
	id, ok := utils.GetClaimValue[string](response, provider.IDClaim)
	if !ok {
		log.Logger.Warningf("Failed to get user_id by provider %s: %s", provider.Provider, response)
		return model.User{}, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": "Get value failed"}}
	}
	name, ok := utils.GetClaimValue[string](response, provider.NameClaim)
	if !ok {
		name = fmt.Sprintf("%s_%s", provider.Provider, utils.RandStr(10))
	}
	email, ok := utils.GetClaimValue[string](response, provider.EmailClaim)
	if !ok {
		email = fmt.Sprintf("%s_%s@example.com", provider.Provider, utils.RandStr(10))
	}
	picture, _ := utils.GetClaimValue[string](response, provider.PictureClaim)
	description, _ := utils.GetClaimValue[string](response, provider.DescriptionClaim)
	raw, _ := json.Marshal(response)
	userRepo := db.InitUserRepo(tx)
	user, ret := userRepo.Get(db.GetOptions{Conditions: map[string]any{"provider": provider.Provider, "provider_user_id": id}})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return model.User{}, ret
		}
		// 获取用户失败的时创建新用户
		user, ret = userRepo.Create(db.CreateUserOptions{
			Name:           name,
			Password:       model.NeverLoginPWD,
			Email:          email,
			Picture:        model.FileURL(picture),
			Description:    description,
			Verified:       true,
			Provider:       provider.Provider,
			ProviderUserID: id,
			OauthRaw:       string(raw),
		})
		if !ret.OK {
			return model.User{}, ret
		}
		if provider.GroupsClaim != "" {
			groupRepo := db.InitGroupRepo(tx)
			if groups, ok := utils.GetClaimValue[[]string](response, provider.GroupsClaim); ok {
				// 同步所有组
				for _, groupName := range groups {
					group, ret := groupRepo.GetByUniqueKey("name", groupName)
					if !ret.OK {
						continue
					}
					if !userRepo.IsInGroup(user.ID, group.Name) {
						db.AppendUserToGroup(tx, user, group)
					}
				}
				// 尝试添加到管理员组
				if slices.Contains(groups, provider.AdminGroup) {
					if !userRepo.IsInGroup(user.ID, model.AdminGroupName) {
						adminGroup, ret := db.InitGroupRepo(tx).GetByUniqueKey("name", model.AdminGroupName)
						if ret.OK {
							db.AppendUserToGroup(tx, user, adminGroup)
						}
					}
				}
			}
		}
		// 获取组声明或加组失败后尝试加入默认组
		if provider.DefaultGroup != 0 {
			defaultGroup, ret := db.InitGroupRepo(tx).GetByID(provider.DefaultGroup)
			if ret.OK {
				// 最终都无法获取到组则放弃加组
				if !userRepo.IsInGroup(user.ID, defaultGroup.Name) {
					db.AppendUserToGroup(tx, user, defaultGroup)
				}
			}
		}
		prometheus.RecordUserRegister(provider.Provider)
	} else {
		// 获取用户成功的时更新用户信息
		ret = userRepo.Update(user.ID, db.UpdateUserOptions{
			Name:        &name,
			Email:       &email,
			Description: &description,
			Picture:     new(model.FileURL(picture)),
			OauthRaw:    new(string(raw)),
		})
		if !ret.OK {
			return model.User{}, ret
		}
		prometheus.RecordUserLogin(provider.Provider)
	}
	return userRepo.Get(db.GetOptions{Conditions: map[string]any{"provider": provider.Provider, "provider_user_id": id}})
}
