package db

import (
	"CBCTF/internal/model"
	"regexp"
)

// isUniqueUserName 用户名不能重复
func isUniqueUserName(name string) bool {
	res := DB.Model(&model.User{}).Where("name = ?", name).Find(&model.User{})
	if res.RowsAffected > 0 {
		return false
	}
	return true
}

// isUniqueTeamName 通常比赛中队伍名不能重复
func isUniqueTeamName(contestID uint, name string) bool {
	var contest model.Contest
	DB.Model(&model.Contest{}).Preload("Teams").Where("id = ?", contestID).Find(&contest).Limit(1)
	for _, team := range contest.Teams {
		if team.Name == name {
			return false
		}
	}
	return true
}

// isNotRepeatPlayer 同个用户不能在同场比赛中加入两个队伍
func isNotRepeatPlayer(contestID uint, userID uint) bool {
	var contest model.Contest
	DB.Model(&model.Contest{}).Preload("Teams.Users").Where("id = ?", contestID).Find(&contest).Limit(1)
	for _, team := range contest.Teams {
		for _, user := range team.Users {
			if user.ID == userID {
				return false
			}
		}
	}
	return true
}

// isUniqueContestName 赛事名称不能重复
func isUniqueContestName(name string) bool {
	res := DB.Model(&model.Contest{}).Where("name = ?", name).Find(&model.Contest{})
	if res.RowsAffected > 0 {
		return false
	}
	return true
}

// isUniqueEmail 注册邮箱不能重复
func isUniqueEmail(email string) bool {
	res := DB.Model(&model.User{}).Where("email = ?", email).Find(&model.User{})
	if res.RowsAffected > 0 {
		return false
	}
	return true
}

// isValidEmail 校验邮箱格式
func isValidEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	if regexp.MustCompile(pattern).MatchString(email) {
		return true
	}
	return false
}

// isTeamUser User 是 Team 成员，team 需要预加载
func isTeamUser(user model.User, team model.Team) bool {
	for _, teamUser := range team.Users {
		if user.ID == teamUser.ID {
			return true
		}
	}
	return false
}

// isContestTeam Team 是 Contest 队伍，contest 需要预加载
func isContestTeam(team model.Team, contest model.Contest) bool {
	for _, contestTeam := range contest.Teams {
		if team.ID == contestTeam.ID {
			return true
		}
	}
	return false
}
