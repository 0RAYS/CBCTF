package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/task"
	"CBCTF/internal/utils"

	"gorm.io/gorm"
)

func SendEmail(user model.User) model.RetVal {
	id := utils.UUID()
	token, err := utils.GenerateToken(user.ID, user.Name, id)
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	if ret := redis.SetEmailVerifyToken(user.ID, id); !ret.OK {
		return ret
	}
	if _, err = task.EnqueueSendEmailTask(user.Email, token, id); err != nil {
		log.Logger.Warningf("Failed to enqueue send email task: %s", err)
		return model.RetVal{Msg: i18n.Task.EnqueueError, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func VerifyEmail(tx *gorm.DB, form dto.VerifyEmail) model.RetVal {
	claims, err := utils.ParseToken(form.Token)
	if err != nil || !utils.CompareMagic(form.ID, claims.X) {
		return model.RetVal{Msg: i18n.Model.Email.InvalidVerifyToken}
	}
	if _, ret := redis.GetEmailVerifyToken(claims.UserID); !ret.OK {
		return ret
	}
	repo := db.InitUserRepo(tx)
	ret := repo.Update(claims.UserID, db.UpdateUserOptions{Verified: new(true)})
	if !ret.OK {
		return ret
	}
	if ret = redis.DelEmailVerifyToken(claims.UserID); !ret.OK {
		return ret
	}
	return model.SuccessRetVal()
}

func ListEmails(tx *gorm.DB, smtp model.Smtp, form dto.ListModelsForm) ([]model.Email, int64, model.RetVal) {
	options := db.GetOptions{}
	if smtp.ID > 0 {
		options.Conditions = map[string]any{"smtp_id": smtp.ID}
	}
	return db.InitEmailRepo(tx).List(form.Limit, form.Offset, options)
}
