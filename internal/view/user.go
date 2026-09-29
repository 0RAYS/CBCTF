package view

import "CBCTF/internal/model"

type UserView struct {
	Unavailable    []string
	User           model.User
	HasAdminAccess bool
	TeamCount      int64
	ContestCount   int64
}
