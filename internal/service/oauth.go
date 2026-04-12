package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
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
	id, ok := utils.GetClaimStringValue(response, provider.IDClaim)
	if !ok {
		log.Logger.Warningf("Failed to get user_id by provider %s: %s", provider.Provider, response)
		return model.User{}, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": "Get value failed"}}
	}
	name, ok := utils.GetClaimStringValue(response, provider.NameClaim)
	if !ok {
		name = fmt.Sprintf("%s_%s", provider.Provider, utils.RandStr(10))
	}
	email, ok := utils.GetClaimStringValue(response, provider.EmailClaim)
	if !ok {
		email = fmt.Sprintf("%s_%s@example.com", provider.Provider, utils.RandStr(10))
	}
	picture, _ := utils.GetClaimStringValue(response, provider.PictureClaim)
	description, _ := utils.GetClaimStringValue(response, provider.DescriptionClaim)
	raw, _ := json.Marshal(response)
	userRepo := db.InitUserRepo(tx)
	user, ret := userRepo.Get(db.GetOptions{Conditions: map[string]any{"provider": provider.Provider, "provider_user_id": id}})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return model.User{}, ret
		}
		// 获取用户失败的时创建新用户
		if !userRepo.IsUniqueKeyValue(0, "name", name) {
			name = fmt.Sprintf("%s_%s", provider.Provider, utils.RandStr(10))
		}
		if !userRepo.IsUniqueKeyValue(0, "email", email) {
			email = fmt.Sprintf("%s_%s@example.com", provider.Provider, utils.RandStr(10))
		}
		user, ret = userRepo.Insert(model.User{
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
			if groups, groupsOK := utils.GetClaimRawValue[[]string](response, provider.GroupsClaim); groupsOK {
				// 同步所有组
				for _, groupName := range groups {
					group, groupRet := groupRepo.GetByUniqueField("name", groupName)
					if !groupRet.OK {
						continue
					}
					if !userRepo.IsInGroup(user.ID, group.Name) {
						db.AppendUserToGroup(tx, user, group)
					}
				}
				// 尝试添加到管理员组
				if slices.Contains(groups, provider.AdminGroup) {
					if !userRepo.IsInGroup(user.ID, model.AdminGroupName) {
						adminGroup, adminGroupRet := db.InitGroupRepo(tx).GetByUniqueField("name", model.AdminGroupName)
						if adminGroupRet.OK {
							db.AppendUserToGroup(tx, user, adminGroup)
						}
					}
				}
			}
		}
		// 获取组声明或加组失败后尝试加入默认组
		if provider.DefaultGroup != 0 {
			defaultGroup, defaultGroupRet := db.InitGroupRepo(tx).GetByID(provider.DefaultGroup)
			if defaultGroupRet.OK {
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
			Description: &description,
			Picture:     new(model.FileURL(picture)),
			OauthRaw:    new(string(raw)),
			Email:       &email,
		})
		if !ret.OK {
			return model.User{}, ret
		}
		prometheus.RecordUserLogin(provider.Provider)
	}
	return userRepo.Get(db.GetOptions{Conditions: map[string]any{"provider": provider.Provider, "provider_user_id": id}})
}

func OauthLoginWithTransaction(tx *gorm.DB, provider model.Oauth, response map[string]any) (model.User, model.RetVal) {
	var user model.User
	ret := db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		var loginRet model.RetVal
		user, loginRet = OauthLogin(tx2, provider, response)
		return loginRet
	})
	return user, ret
}

func ListEnabledOauthProviders(tx *gorm.DB) ([]model.Oauth, model.RetVal) {
	providers, _, ret := db.InitOauthRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"on": true},
	})
	return providers, ret
}

func ListOauthProviders(tx *gorm.DB, form dto.ListModelsForm) ([]model.Oauth, int64, model.RetVal) {
	return db.InitOauthRepo(tx).List(form.Limit, form.Offset)
}

func CreateOauthProvider(tx *gorm.DB, form dto.CreateOauthProviderForm) (model.Oauth, model.RetVal) {
	return db.InitOauthRepo(tx).Create(db.CreateOauthOptions{
		AuthURL:          form.AuthURL,
		TokenURL:         form.TokenURL,
		UserInfoURL:      form.UserInfoURL,
		CallbackURL:      form.CallbackURL,
		ClientID:         form.ClientID,
		ClientSecret:     form.ClientSecret,
		Provider:         form.Provider,
		Uri:              form.Uri,
		IDClaim:          form.IDClaim,
		NameClaim:        form.NameClaim,
		EmailClaim:       form.EmailClaim,
		PictureClaim:     form.PictureClaim,
		DescriptionClaim: form.DescriptionClaim,
		GroupsClaim:      form.GroupsClaim,
		AdminGroup:       form.AdminGroup,
		DefaultGroup:     form.DefaultGroup,
		On:               false,
	})
}

func UpdateOauthProvider(tx *gorm.DB, oldOauth model.Oauth, form dto.UpdateOauthProviderForm) (model.Oauth, model.RetVal) {
	if ret := db.InitOauthRepo(tx).Update(oldOauth.ID, db.UpdateOauthOptions{
		AuthURL:          form.AuthURL,
		TokenURL:         form.TokenURL,
		UserInfoURL:      form.UserInfoURL,
		CallbackURL:      form.CallbackURL,
		ClientID:         form.ClientID,
		ClientSecret:     form.ClientSecret,
		Provider:         form.Provider,
		Uri:              form.Uri,
		IDClaim:          form.IDClaim,
		NameClaim:        form.NameClaim,
		EmailClaim:       form.EmailClaim,
		PictureClaim:     form.PictureClaim,
		DescriptionClaim: form.DescriptionClaim,
		GroupsClaim:      form.GroupsClaim,
		AdminGroup:       form.AdminGroup,
		DefaultGroup:     form.DefaultGroup,
		On:               form.On,
		Picture:          form.Picture,
	}); !ret.OK {
		return model.Oauth{}, ret
	}
	return db.InitOauthRepo(tx).GetByID(oldOauth.ID)
}

func DeleteOauthProvider(tx *gorm.DB, oauth model.Oauth) model.RetVal {
	return db.InitOauthRepo(tx).Delete(oauth.ID)
}
