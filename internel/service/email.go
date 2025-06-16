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
	token, err := utils.Generate(user.ID, user.Name, false, "email")
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		return false, i18n.UnknownError
	}
	if err = redis.SetEmailVerifyToken(user.ID, id); err != nil {
		log.Logger.Warningf("Failed to set email verify token: %s", err)
		return false, i18n.SetEmailVerifyTokenError
	}
	go func() {
		if err = utils.SendVerifyEmail(user.Email, token, id); err != nil {
			log.Logger.Warningf("Failed to send mail: %s", err)
		}
	}()
	return true, i18n.Success
}

func VerifyEmail(tx *gorm.DB, form f.VerifyEmail) (bool, string) {
	claims, err := utils.Parse(form.Token)
	if err != nil {
		return false, i18n.InvalidEmailVerifyToken
	}
	id, err := redis.GetEmailVerifyToken(claims.UserID)
	if err != nil {
		return false, id
	}
	if form.ID == id {
		verified := true
		repo := db.InitUserRepo(tx)
		ok, msg := repo.Update(claims.UserID, db.UpdateUserOptions{Verified: &verified})
		if !ok {
			return false, msg
		}
		if err = redis.DelEmailVerifyToken(claims.UserID); err != nil {
			log.Logger.Warningf("Failed to delete email verify token: %s", err)
			return false, i18n.DelEmailVerifyTokenError
		}
		return true, i18n.Success
	}
	return false, i18n.InvalidEmailVerifyToken
}
