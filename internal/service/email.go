package service

import (
	"CBCTF/internal/config"
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
	token, err := utils.GenerateToken(user.ID, user.Name, id, config.Env.Gin.JWT.Secret)
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
	claims, err := utils.ParseToken(form.Token, config.Env.Gin.JWT.Secret)
	if err != nil || !utils.CompareMagic(form.ID, claims.X) {
		return model.RetVal{Msg: i18n.Model.Email.InvalidVerifyToken}
	}
	if _, ret := redis.GetEmailVerifyToken(claims.UserID); !ret.OK {
		return ret
	}
	if ret := redis.DelEmailVerifyToken(claims.UserID); !ret.OK {
		return ret
	}
	return db.InitUserRepo(tx).Update(claims.UserID, db.UpdateUserOptions{Verified: new(true)})
}

// SendPasswordResetEmail 向用户邮箱发送密码重置链接
// 即使邮箱不存在也返回成功，防止用户枚举
func SendPasswordResetEmail(tx *gorm.DB, form dto.ForgotPasswordForm) model.RetVal {
	user, ret := db.InitUserRepo(tx).GetByUniqueField("email", form.Email)
	if !ret.OK {
		return ret
	}
	id := utils.UUID()
	token, err := utils.GenerateToken(user.ID, user.Name, id, config.Env.Gin.JWT.Secret)
	if err != nil {
		log.Logger.Warningf("Failed to generate password reset token: %s", err)
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	if ret := redis.SetPasswordResetToken(user.ID, id); !ret.OK {
		return ret
	}
	if _, err = task.EnqueueSendResetPasswordEmailTask(user.Email, token, id); err != nil {
		log.Logger.Warningf("Failed to enqueue send reset password email task: %s", err)
		return model.RetVal{Msg: i18n.Task.EnqueueError, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

// ResetUserPassword 验证重置 token 并更新密码，同时将邮箱设为已验证
func ResetUserPassword(tx *gorm.DB, form dto.ResetPasswordForm) model.RetVal {
	claims, err := utils.ParseToken(form.Token, config.Env.Gin.JWT.Secret)
	if err != nil || !utils.CompareMagic(form.ID, claims.X) {
		return model.RetVal{Msg: i18n.Model.User.InvalidResetToken}
	}
	if _, ret := redis.GetPasswordResetToken(claims.UserID); !ret.OK {
		return model.RetVal{Msg: i18n.Model.User.InvalidResetToken}
	}
	if ret := redis.DelPasswordResetToken(claims.UserID); !ret.OK {
		return ret
	}
	hashedPwd := utils.HashPassword(form.Password)
	verified := true
	return db.InitUserRepo(tx).Update(claims.UserID, db.UpdateUserOptions{
		Password: &hashedPwd,
		Verified: &verified,
	})
}

func ListEmails(tx *gorm.DB, smtp model.Smtp, form dto.ListModelsForm) ([]model.Email, int64, model.RetVal) {
	options := db.GetOptions{Sort: []string{"created_at DESC"}}
	if smtp.ID > 0 {
		options.Conditions = map[string]any{"smtp_id": smtp.ID}
	}
	return db.InitEmailRepo(tx).List(form.Limit, form.Offset, options)
}
