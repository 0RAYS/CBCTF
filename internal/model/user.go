package model

const NeverLoginPWD = "never_login_pwd"

// User
// ManyToMany Teams
// ManyToMany Contests
// ManyToMany Groups
// HasMany Submissions
type User struct {
	Teams          []Team       `gorm:"many2many:user_teams;" json:"-"`
	Contests       []Contest    `gorm:"many2many:user_contests;" json:"-"`
	Submissions    []Submission `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Groups         []Group      `gorm:"many2many:user_groups;" json:"-"`
	Name           string       `gorm:"type:varchar(255);uniqueIndex:idx_users_name_active,where:deleted_at IS NULL;not null" json:"name"`
	Password       string       `gorm:"not null" json:"-"`
	Email          string       `gorm:"type:varchar(255);uniqueIndex:idx_users_email_active,where:deleted_at IS NULL;not null" json:"email"`
	Picture        FileURL      `json:"picture"`
	Description    string       `json:"description"`
	Verified       bool         `gorm:"default:false" json:"verified"`
	Hidden         bool         `gorm:"default:false" json:"hidden"`
	Banned         bool         `gorm:"default:false" json:"banned"`
	Score          float64      `gorm:"default:0" json:"score"`
	Solved         int64        `gorm:"default:0" json:"solved"`
	Provider       string       `gorm:"type:varchar(255);uniqueIndex:idx_users_provider_id_active,where:deleted_at IS NULL;not null" json:"provider"`
	ProviderUserID string       `gorm:"type:varchar(255);uniqueIndex:idx_users_provider_id_active,where:deleted_at IS NULL;not null" json:"provider_user_id"`
	OauthRaw       string       `json:"oauth_raw"`
	BaseModel
}
