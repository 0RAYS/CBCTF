package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/resp"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func Search(ctx *gin.Context) {
	var form f.SearchForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	query := ctx.Request.URL.Query()
	options := db.GetOptions{SearchConditions: make(map[string]any)}
	switch form.Model {
	case "user":
		allowedKeys := []string{"name", "email", "id"}
		for key, value := range query {
			if slices.Contains(allowedKeys, key) {
				if len(value) > 0 {
					options.SearchConditions[key] = value[0]
				}
			}
		}
		users, count, ok, msg := db.InitUserRepo(db.DB.WithContext(ctx)).List(form.Limit, form.Offset, options)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		data := make([]gin.H, 0)
		for _, user := range users {
			data = append(data, resp.GetUserResp(user, true))
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"results": data, "count": count}})
		return
	case "contest":
		allowedKeys := []string{"name", "id"}
		for key, value := range query {
			if slices.Contains(allowedKeys, key) {
				if len(value) > 0 {
					options.SearchConditions[key] = value[0]
				}
			}
		}
		options.Preloads = map[string]db.GetOptions{
			"Teams": {Selects: []string{"id"}}, "Users": {Selects: []string{"id"}}, "Notices": {Selects: []string{"id"}},
		}
		contests, count, ok, msg := db.InitContestRepo(db.DB.WithContext(ctx)).List(form.Limit, form.Offset, options)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		data := make([]gin.H, 0)
		for _, contest := range contests {
			data = append(data, resp.GetContestResp(contest, true))
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"results": data, "count": count}})
		return
	case "team":
		allowedKeys := []string{"name", "id"}
		for key, value := range query {
			if slices.Contains(allowedKeys, key) {
				if len(value) > 0 {
					options.SearchConditions[key] = value[0]
				}
			}
		}
		options.Preloads = map[string]db.GetOptions{
			"Users": {Selects: []string{"id"}},
		}
		teams, count, ok, msg := db.InitTeamRepo(db.DB.WithContext(ctx)).List(form.Limit, form.Offset, options)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		data := make([]gin.H, 0)
		for _, team := range teams {
			data = append(data, resp.GetTeamResp(team))
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"results": data, "count": count}})
		return
	case "challenge":
		allowedKeys := []string{"name", "id", "category", "type"}
		for key, value := range query {
			if slices.Contains(allowedKeys, key) {
				if key == "id" {
					key = "rand_id"
				}
				if len(value) > 0 {
					options.SearchConditions[key] = value[0]
				}
			}
		}
		challenges, count, ok, msg := db.InitChallengeRepo(db.DB.WithContext(ctx)).List(form.Limit, form.Offset, options)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		data := make([]gin.H, 0)
		for _, challenge := range challenges {
			data = append(data, resp.GetChallengeResp(challenge))
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"results": data, "count": count}})
		return
	default:
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
}
