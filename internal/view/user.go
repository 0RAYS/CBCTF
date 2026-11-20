package view

import "CBCTF/internal/model"

type UserView struct {
	Unavailable    []string
	User           model.User
	HasAdminAccess bool
	TeamCount      int64
	ContestCount   int64
}

// UserRankingView 用户总分榜中的一行，Rank 是参与排名的用户中的名次。
type UserRankingView struct {
	User model.User
	Rank int
}
