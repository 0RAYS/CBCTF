package service

import (
	"CBCTF/internal/email"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	db "CBCTF/internal/repo"
	"CBCTF/internal/utils"

	"gorm.io/gorm"
)

func SendEmail(user model.User) (bool, string) {
	id := utils.UUID()
	token, err := utils.GenerateToken(user.ID, user.Name, false, "email")
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		return false, i18n.UnknownError
	}
	if err = redis.SetEmailVerifyToken(user.ID, id); err != nil {
		log.Logger.Warningf("Failed to set email verify token: %s", err)
		return false, i18n.SetEmailVerifyTokenError
	}
	go func() {
		if err = email.SendVerifyEmail(user.Email, token, id); err != nil {
			go prometheus.IncEmailSentMetrics(false)
			log.Logger.Warningf("Failed to send mail: %s", err)
		} else {
			go prometheus.IncEmailSentMetrics(true)
		}
	}()
	return true, i18n.Success
}

func VerifyEmail(tx *gorm.DB, form f.VerifyEmail) (bool, string) {
	claims, err := utils.ParseToken(form.Token)
	if err != nil {
		return false, i18n.InvalidEmailVerifyToken
	}
	id, err := redis.GetEmailVerifyToken(claims.UserID)
	if err != nil {
		return false, i18n.GetEmailVerifyTokenError
	}
	if form.ID == id {
		repo := db.InitUserRepo(tx)
		ok, msg := repo.Update(claims.UserID, db.UpdateUserOptions{Verified: utils.Ptr(true)})
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
