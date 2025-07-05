package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
)

func Search(ctx *gin.Context) {
	var form f.SearchForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	switch form.Model {
	case "user":
		allowedKeys := []string{"name", "email", "id"}
		if !slices.Contains(allowedKeys, form.Key) {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
			return
		}
		users, count, ok, msg := db.InitUserRepo(db.DB.WithContext(ctx)).FuzzSearch(form.Limit, form.Offset, form.Key, form.Value)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		data := make([]gin.H, 0)
		for _, user := range users {
			data = append(data, resp.GetUserResp(user, true))
		}
		ctx.JSON(http.StatusOK, gin.H{"results": data, "count": count})
		return
	case "contest":
		allowedKeys := []string{"name", "id"}
		if !slices.Contains(allowedKeys, form.Key) {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
			return
		}
		contests, count, ok, msg := db.InitContestRepo(db.DB.WithContext(ctx)).FuzzSearch(form.Limit, form.Offset, form.Key, form.Value)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		data := make([]gin.H, 0)
		for _, contest := range contests {
			data = append(data, resp.GetContestResp(contest, true))
		}
		ctx.JSON(http.StatusOK, gin.H{"results": data, "count": count})
		return
	case "team":
		allowedKeys := []string{"name", "id"}
		if !slices.Contains(allowedKeys, form.Key) {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
			return
		}
		teams, count, ok, msg := db.InitTeamRepo(db.DB.WithContext(ctx)).FuzzSearch(form.Limit, form.Offset, form.Key, form.Value)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		data := make([]gin.H, 0)
		for _, team := range teams {
			data = append(data, resp.GetTeamResp(team))
		}
		ctx.JSON(http.StatusOK, gin.H{"results": data, "count": count})
		return
	case "challenge":
		allowedKeys := []string{"name", "id"}
		if !slices.Contains(allowedKeys, form.Key) {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
			return
		}
		challenges, count, ok, msg := db.InitChallengeRepo(db.DB.WithContext(ctx)).FuzzSearch(form.Limit, form.Offset, form.Key, form.Value)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		data := make([]gin.H, 0)
		for _, challenge := range challenges {
			data = append(data, resp.GetChallengeResp(challenge))
		}
		ctx.JSON(http.StatusOK, gin.H{"results": data, "count": count})
		return
	default:
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
}
