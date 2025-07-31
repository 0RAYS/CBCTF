package service

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"gorm.io/gorm"
)

func CreateOauthProvider(tx *gorm.DB, form f.CreateOauthProviderForm) (model.Oauth, bool, string) {
	return db.InitOauthRepo(tx).Create(db.CreateOauthOptions{
		AuthURL:         form.AuthURL,
		TokenURL:        form.TokenURL,
		UserInfoURL:     form.UserInfoURL,
		ClientID:        form.ClientID,
		ClientSecret:    form.ClientSecret,
		Provider:        form.Provider,
		URI:             form.URI,
		RespIDField:     form.RespIDField,
		RespNameField:   form.RespNameField,
		RespEmailField:  form.RespEmailField,
		RespAvatarField: form.RespAvatarField,
		RespDescField:   form.RespDescField,
		On:              false,
	})
}

func UpdateOauthProvider(tx *gorm.DB, oauth model.Oauth, form f.UpdateOauthProviderForm) (bool, string) {
	return db.InitOauthRepo(tx).Update(oauth.ID, db.UpdateOauthOptions{
		AuthURL:         form.AuthURL,
		TokenURL:        form.TokenURL,
		UserInfoURL:     form.UserInfoURL,
		ClientID:        form.ClientID,
		ClientSecret:    form.ClientSecret,
		Provider:        form.Provider,
		URI:             form.URI,
		RespIDField:     form.RespIDField,
		RespNameField:   form.RespNameField,
		RespEmailField:  form.RespEmailField,
		RespAvatarField: form.RespAvatarField,
		RespDescField:   form.RespDescField,
		On:              form.On,
	})
}
