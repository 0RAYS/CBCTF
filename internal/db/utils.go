package db

import (
	"CBCTF/internal/model"
	"gorm.io/gorm"
	"regexp"
)

// isValidEmail 邮箱格式验证
func isValidEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	if regexp.MustCompile(pattern).MatchString(email) {
		return true
	}
	return false
}

// isUniqueEmail 邮箱不能重复
func isUniqueEmail(email string, v interface{}) bool {
	var res *gorm.DB
	switch v.(type) {
	case model.User:
		res = DB.Model(&model.User{}).Where("email = ?", email).Find(&model.User{}).Limit(1)
	case model.Admin:
		res = DB.Model(&model.Admin{}).Where("email = ?", email).Find(&model.Admin{}).Limit(1)
	default:
		return false
	}
	if res.RowsAffected > 0 {
		return false
	}
	return true
}

// isUniqueName 对象名不能重复, 但在此处不考虑Team
func isUniqueName(name string, v interface{}) bool {
	var res *gorm.DB
	switch v.(type) {
	case model.User:
		res = DB.Model(&model.User{}).Where("name = ?", name).Find(&model.User{}).Limit(1)
	case model.Admin:
		res = DB.Model(&model.Admin{}).Where("name = ?", name).Find(&model.Admin{}).Limit(1)
	case model.Contest:
		res = DB.Model(&model.Contest{}).Where("name = ?", name).Find(&model.Contest{}).Limit(1)
	default:
		return false
	}
	if res.RowsAffected > 0 {
		return false
	}
	return true
}

// isUniqueTeamName 在每个Contest中, 队伍名不能重复
func isUniqueTeamName(name string, id uint) bool {
	res := DB.Model(&model.Team{}).Where("name = ? AND contest_id = ?", name, id).Find(&model.Team{}).Limit(1)
	if res.RowsAffected > 0 {
		return false
	}
	return true
}

// isUniqueTeamMember model.User 不能在同一个 model.Contest 出现多次
func isUniqueTeamMember(contestID uint, userID uint) bool {
	var tmp []model.User
	err := DB.Model(&model.User{ID: userID}).Where("contest_id = ?", contestID).Association("Contests").Find(&tmp)
	if len(tmp) > 0 || err != nil {
		return false
	}
	return true
}

// isMemberInTeam model.User 是否在 model.Team 中
func isMemberInTeam(teamID uint, userID uint) bool {
	var tmp []model.Team
	err := DB.Model(&model.User{ID: userID}).Where("team_id = ?", teamID).Association("Teams").Find(&tmp)
	if len(tmp) > 0 || err != nil {
		return true
	}
	return false
}
