package model

type Oauth struct {
	AuthURL              string  `json:"auth_url"`
	TokenURL             string  `json:"token_url"`
	UserInfoURL          string  `json:"user_info_url"`
	CallbackURL          string  `json:"callback_url"`
	ClientID             string  `json:"client_id"`
	ClientSecret         string  `json:"client_secret"`
	Provider             string  `gorm:"type:varchar(255);uniqueIndex;not null" json:"provider"`
	Uri                  string  `json:"uri"`
	RespIDField          string  `json:"id_field"`
	RespNameField        string  `json:"name_field"`
	RespEmailField       string  `json:"email_field"`
	RespPictureField     string  `json:"picture_field"`
	RespDescriptionField string  `json:"description_field"`
	On                   bool    `json:"on"`
	Picture              FileURL `json:"picture"`
	BaseModel
}

func (o Oauth) ModelName() string {
	return "Oauth"
}

func (o Oauth) GetBaseModel() BaseModel {
	return o.BaseModel
}

func (o Oauth) UniqueFields() []string {
	return []string{"id", "provider"}
}

func (o Oauth) QueryFields() []string {
	return []string{"id", "provider", "on"}
}
