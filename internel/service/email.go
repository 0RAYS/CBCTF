package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/redis"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

func SendEmail(user model.User) (bool, string) {
	id := utils.UUID()
	token, err := utils.Generate(user.ID, user.Name, "email", "email")
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		return false, i18n.UnknownError
	}
	ok, msg := redis.SetEmailVerifyToken(user.ID, id)
	if !ok {
		return false, msg
	}
	if err = utils.SendVerifyEmail(user.Email, token, id); err != nil {
		log.Logger.Warningf("Failed to send mail: %s", err)
		return false, i18n.SendEmailError
	}
	return true, i18n.Success
}

func VerifyEmail(tx *gorm.DB, form f.VerifyEmail) (bool, string) {
	claims, err := utils.Parse(form.Token)
	if err != nil {
		return false, i18n.InvalidEmailVerifyToken
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
		return true, i18n.Success
	}
	return false, i18n.InvalidEmailVerifyToken
}
