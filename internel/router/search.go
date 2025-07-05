package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	db "CBCTF/internel/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Search(ctx *gin.Context) {
	var form f.SearchForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	var (
		data  any
		count int64
		msg   string
	)
	switch form.Model {
	case "user":
		data, count, _, msg = db.InitUserRepo(db.DB.WithContext(ctx)).FuzzSearch(form.Limit, form.Offset, form.Key, form.Value)
	case "contest":
		data, count, _, msg = db.InitContestRepo(db.DB.WithContext(ctx)).FuzzSearch(form.Limit, form.Offset, form.Key, form.Value)
	case "team":
		data, count, _, msg = db.InitTeamRepo(db.DB.WithContext(ctx)).FuzzSearch(form.Limit, form.Offset, form.Key, form.Value)
	case "challenge":
		data, count, _, msg = db.InitChallengeRepo(db.DB.WithContext(ctx)).FuzzSearch(form.Limit, form.Offset, form.Key, form.Value)
	default:
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "data": data}})
}
