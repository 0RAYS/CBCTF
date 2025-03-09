package form

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Strings []string

func (s Strings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Strings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Strings value")
	}
	return json.Unmarshal(bytes, s)
}

// GetModelsForm for get models list
type GetModelsForm struct {
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit"`
}

// ChangePasswordForm for user or admin change password
type ChangePasswordForm struct {
	OldPassword string `form:"oldPassword" json:"oldPassword" binding:"required"`
	NewPassword string `form:"newPassword" json:"newPassword" binding:"required"`
}
