package model

const NeverLoginPWD = "never_login_pwd"

// User
// ManyToMany Teams
// ManyToMany Contests
// HasMany Devices
// HasMany Submissions
type User struct {
	Teams          []Team       `gorm:"many2many:user_teams;" json:"-"`
	Contests       []Contest    `gorm:"many2many:user_contests;" json:"-"`
	Submissions    []Submission `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Groups         []Group      `gorm:"many2many:user_groups;" json:"-"`
	Name           string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Password       string       `gorm:"not null" json:"-"`
	Email          string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Picture        FileURL      `json:"picture"`
	Description    string       `json:"description"`
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

func (u User) TableName() string {
	return "users"
}

func (u User) ModelName() string {
	return "User"
}

func (u User) GetBaseModel() BaseModel {
	return u.BaseModel
}

func (u User) UniqueFields() []string {
	return []string{"id", "name", "email"}
}

func (u User) QueryFields() []string {
	return []string{"id", "name", "email", "description", "verified", "banned", "hidden", "provider"}
}
