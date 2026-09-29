package view

import "CBCTF/internal/model"

type GroupView struct {
	Unavailable []string
	Group       model.Group
	UserCount   int64
}
