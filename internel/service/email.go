package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/redis"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

func VerifyEmail(tx *gorm.DB, form f.VerifyEmail) (bool, string) {
	claims, err := utils.Parse(form.Token)
	if err != nil {
		return false, "InvalidEmailVerifyToken"
	}
	id, ok := redis.GetEmailVerifyToken(claims.UserID)
	if !ok {
		return false, id
	}
	if form.ID == id {
		verified := true
		repo := db.InitUserRepo(tx)
		ok, msg := repo.Update(claims.UserID, db.UpdateUserOptions{Verified: &verified})
		if !ok {
			return false, msg
		}
		redis.DelEmailVerifyToken(claims.UserID)
		return true, "Success"
	}
	return false, "InvalidEmailVerifyToken"
}
