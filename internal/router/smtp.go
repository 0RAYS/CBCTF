package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/email"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSmtp(ctx *gin.Context) {
	smtp := middleware.GetSmtp(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetSmtpResp(smtp)})
}

func GetSmtps(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	smtps, count, ok, msg := db.InitSmtpRepo(db.DB).List(form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, smtp := range smtps {
		data = append(data, resp.GetSmtpResp(smtp))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"smtps": data, "count": count}})
}

func CreateSmtp(ctx *gin.Context) {
	var form f.CreateSmtpForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateSmtpEventType)
	tx := db.DB.Begin()
	smtp, ok, msg := service.CreateSmtp(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.GetSmtpResp(smtp)})
}

func UpdateSmtp(ctx *gin.Context) {
	var form f.UpdateSmtpForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateSmtpEventType)
	smtp := middleware.GetSmtp(ctx)
	tx := db.DB.Begin()
	if ok, msg := service.UpdateSmtp(tx, smtp, form); !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	newSmtp, ok, msg := db.InitSmtpRepo(db.DB).GetByID(smtp.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	email.DelSenders(smtp)
	if newSmtp.On {
		email.AddSenders(smtp)
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteSmtp(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteSmtpEventType)
	smtp := middleware.GetSmtp(ctx)
	tx := db.DB.Begin()
	if ok, msg := db.InitSmtpRepo(tx).Delete(smtp.ID); !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	email.DelSenders(smtp)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}
