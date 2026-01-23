package model

const NeverLoginPWD = "never_login_pwd"

// User
// ManyToMany Teams
// ManyToMany Contests
// HasMany Devices
// HasMany Submissions
type User struct {
	Teams          []*Team      `gorm:"many2many:user_teams;" json:"-"`
	Contests       []*Contest   `gorm:"many2many:user_contests;" json:"-"`
	Devices        []Device     `json:"-"`
	Submissions    []Submission `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name           string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Password       string       `gorm:"not null" json:"-"`
	Email          string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Country        string       `gorm:"default:'CN'" json:"country"`
	Avatar         AvatarURL    `json:"avatar"`
	Desc           string       `json:"desc"`
	Verified       bool         `gorm:"default:false" json:"verified"`
	Hidden         bool         `gorm:"default:false" json:"hidden"`
	Banned         bool         `gorm:"default:false" json:"banned"`
	Score          float64      `gorm:"default:0" json:"score"`
	Solved         int64        `gorm:"default:0" json:"solved"`
	Provider       string       `gorm:"type:varchar(255);index:idx_provider_id,unique;not null" json:"provider"`
	ProviderUserID string       `gorm:"type:varchar(255);index:idx_provider_id,unique;not null" json:"provider_user_id"`
	OauthRaw       string       `json:"oauth_raw"`
	BaseModel
}

func (u User) GetModelName() string {
	return "User"
}

func (u User) GetBaseModel() BaseModel {
	return u.BaseModel
}

func (u User) GetUniqueKey() []string {
	return []string{"id", "name", "email"}
}

func (u User) GetAllowedQueryFields() []string {
	return []string{"id", "name", "email", "country", "desc", "verified", "banned", "hidden"}
}
