package service

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/task"
	"CBCTF/internal/utils"

	"gorm.io/gorm"
)

func SendEmail(user model.User) (bool, string) {
	id := utils.UUID()
	token, err := utils.GenerateToken(user.ID, user.Name, false, id)
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		return false, i18n.UnknownError
	}
	if err = redis.SetEmailVerifyToken(user.ID, id); err != nil {
		log.Logger.Warningf("Failed to set email verify token: %s", err)
		return false, i18n.SetEmailVerifyTokenError
	}
	if _, err = task.EnqueueSendEmailTask(user.Email, token, id); err != nil {
		log.Logger.Warningf("Failed to enqueue send email task: %s", err)
		return false, i18n.SendEmailError
	}
	return true, i18n.Success
}

func VerifyEmail(tx *gorm.DB, form f.VerifyEmail) (bool, string) {
	claims, err := utils.ParseToken(form.Token)
	if err != nil || !utils.CompareMagic(form.ID, claims.X) {
		return false, i18n.InvalidEmailVerifyToken
	}
	if _, err = redis.GetEmailVerifyToken(claims.UserID); err != nil {
		return false, i18n.GetEmailVerifyTokenError
	}
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
