package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/email"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"

	"github.com/gin-gonic/gin"
)

func GetSmtps(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	smtps, count, ret := db.InitSmtpRepo(db.DB).List(form.Limit, form.Offset)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, smtp := range smtps {
		data = append(data, resp.GetSmtpResp(smtp))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"smtps": data, "count": count}))
}

func CreateSmtp(ctx *gin.Context) {
	var form dto.CreateSmtpForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateSmtpEventType)
	smtp, ret := db.InitSmtpRepo(db.DB).Create(db.CreateSmtpOptions{
		Address: form.Address,
		Host:    form.Host,
		Port:    form.Port,
		Pwd:     form.Pwd,
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetSmtpResp(smtp)))
}

func UpdateSmtp(ctx *gin.Context) {
	var form dto.UpdateSmtpForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateSmtpEventType)
	smtp := middleware.GetSmtp(ctx)
	if ret := db.InitSmtpRepo(db.DB).Update(smtp.ID, db.UpdateSmtpOptions{
		Address: form.Address,
		Host:    form.Host,
		Port:    form.Port,
		Pwd:     form.Pwd,
		On:      form.On,
	}); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	newSmtp, ret := db.InitSmtpRepo(db.DB).GetByID(smtp.ID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	email.DelSenders(smtp)
	if newSmtp.On {
		email.AddSenders(smtp)
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}

func DeleteSmtp(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteSmtpEventType)
	smtp := middleware.GetSmtp(ctx)
	if ret := db.InitSmtpRepo(db.DB).Delete(smtp.ID); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	email.DelSenders(smtp)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
