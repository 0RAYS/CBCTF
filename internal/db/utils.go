package db

import (
	"CBCTF/internal/model"
	"gorm.io/gorm"
	"regexp"
)

// IsValidEmail 邮箱格式验证
func IsValidEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	if regexp.MustCompile(pattern).MatchString(email) {
		return true
	}
	return false
}

// IsUniqueEmail 邮箱不能重复
func IsUniqueEmail(tx *gorm.DB, email string) bool {
	var res *gorm.DB
	res = tx.Model(&model.User{}).Where("email = ?", email).Find(&model.User{}).Limit(1)
	if res.RowsAffected > 0 {
		return false
	}
	res = tx.Model(&model.Admin{}).Where("email = ?", email).Find(&model.Admin{}).Limit(1)
	if res.RowsAffected > 0 {
		return false
	}
	return true
}

// IsUniqueName 对象名不能重复, 但在此处不考虑Team
func IsUniqueName(tx *gorm.DB, name string, v interface{}) bool {
	var res *gorm.DB
	switch v.(type) {
	case model.User:
		res = tx.Model(&model.User{}).Where("name = ?", name).Find(&model.User{}).Limit(1)
	case model.Admin:
		res = tx.Model(&model.Admin{}).Where("name = ?", name).Find(&model.Admin{}).Limit(1)
	case model.Contest:
		res = tx.Model(&model.Contest{}).Where("name = ?", name).Find(&model.Contest{}).Limit(1)
	default:
		return false
	}
	if res.RowsAffected > 0 {
		return false
	}
	return true
}

// IsUniqueTeamName 在每个Contest中, 队伍名不能重复, 锁定 Team 表
func IsUniqueTeamName(tx *gorm.DB, name string, id uint) bool {
	res := tx.Model(&model.Team{}).Where("name = ? AND contest_id = ?", name, id).Find(&model.Team{}).Limit(1)
	if res.RowsAffected > 0 {
		return false
	}
	return true
}

// IsUniqueTeamMember model.User 不能在同一个 model.Contest 出现多次, 锁定关联表
func IsUniqueTeamMember(tx *gorm.DB, contestID uint, userID uint) bool {
	var tmp []model.User
	err := tx.Model(&model.User{ID: userID}).Where("contest_id = ?", contestID).Association("Contests").Find(&tmp)
	if len(tmp) > 0 || err != nil {
		return false
	}
	return true
}

// IsMemberInTeam model.User 是否在 model.Team 中, 锁定关联表
func IsMemberInTeam(tx *gorm.DB, teamID uint, userID uint) bool {
	var tmp []model.Team
	err := tx.Model(&model.User{ID: userID}).Where("team_id = ?", teamID).Association("Teams").Find(&tmp)
	if err != nil || len(tmp) <= 0 {
		return false
	}
	return true
}

// IsValidChallengeType 题目类型验证
func IsValidChallengeType(t int) bool {
	if t >= 0 && t < 3 {
		return true
	}
	return false
}
