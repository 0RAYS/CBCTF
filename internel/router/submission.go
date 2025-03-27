package router

import (
	f "CBCTF/internel/form"
	db "CBCTF/internel/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSubmissions(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	submissions, count, ok, msg := db.InitSubmissionRepo(db.DB.WithContext(ctx)).GetAll(form.Limit, form.Offset, false, 0)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"submissions": &submissions, "count": count}})
}
