package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/email"
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSmtp(ctx *gin.Context) {
	smtp := middleware.GetSmtp(ctx)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetSmtpResp(smtp)))
}

func GetSmtps(ctx *gin.Context) {
	var form f.ListModelsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	smtps, count, ret := db.InitSmtpRepo(db.DB).List(form.Limit, form.Offset)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, smtp := range smtps {
		data = append(data, resp.GetSmtpResp(smtp))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"smtps": data, "count": count}))
}

func CreateSmtp(ctx *gin.Context) {
	var form f.CreateSmtpForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
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
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetSmtpResp(smtp)))
}

func UpdateSmtp(ctx *gin.Context) {
	var form f.UpdateSmtpForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
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
		ctx.JSON(http.StatusOK, ret)
		return
	}
	newSmtp, ret := db.InitSmtpRepo(db.DB).GetByID(smtp.ID)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	email.DelSenders(smtp)
	if newSmtp.On {
		email.AddSenders(smtp)
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, ret)
}

func DeleteSmtp(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteSmtpEventType)
	smtp := middleware.GetSmtp(ctx)
	if ret := db.InitSmtpRepo(db.DB).Delete(smtp.ID); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	email.DelSenders(smtp)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}
