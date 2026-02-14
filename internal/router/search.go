package router

import (
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SearchIP(ctx *gin.Context) {
	var form dto.SearchIP
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ip, err := utils.SearchIP(form.IP)
	if err != nil {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	data := gin.H{"iso": ip.Country.ISOCode, "timezone": ip.Location.TimeZone}
	latitude, longitude := ip.Location.Latitude, ip.Location.Longitude
	if latitude != nil {
		data["latitude"] = *latitude
	}
	if longitude != nil {
		data["longitude"] = *longitude
	}
	if strings.Contains(strings.ToLower(i18n.DetectLanguage(ctx)), "zh-cn") {
		data["country"] = ip.Country.Names.SimplifiedChinese
		data["city"] = ip.City.Names.SimplifiedChinese
		data["subdivision"] = ""
		for _, sub := range ip.Subdivisions {
			data["subdivision"] = fmt.Sprintf("%s / %s", data["subdivision"], sub.Names.SimplifiedChinese)
		}
	} else {
		data["country"] = ip.Country.Names.English
		data["city"] = ip.City.Names.English
		data["subdivision"] = ""
		for _, sub := range ip.Subdivisions {
			data["subdivision"] = fmt.Sprintf("%s / %s", data["subdivision"], sub.Names.English)
		}
	}
	data["subdivision"] = strings.Trim(data["subdivision"].(string), " / ")
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}

func Search(ctx *gin.Context) {
	//var form dto.SearchForm
	//if ret := form.Bind(ctx); !ret.OK {
	//	ctx.JSON(http.StatusOK, ret)
	//	return
	//}
	//query := ctx.Request.URL.Query()
	//options := db.GetOptions{Search: make(map[string]any)}
	//switch form.Model {
	//case "user":
	//	allowedKeys := []string{"name", "email", "id"}
	//	for key, value := range query {
	//		if slices.Contains(allowedKeys, key) {
	//			if len(value) > 0 {
	//				options.Search[key] = value[0]
	//			}
	//		}
	//	}
	//	users, count, ret := db.InitUserRepo(db.DB).List(form.Limit, form.Offset, options)
	//	if !ok {
	//		ctx.JSON(http.StatusOK, ret)
	//		return
	//	}
	//	data := make([]gin.H, 0)
	//	for _, user := range users {
	//		data = append(data, resp.GetUserResp(user, true))
	//	}
	//	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"results": data, "count": count}})
	//	return
	//case "contest":
	//	allowedKeys := []string{"name", "id"}
	//	for key, value := range query {
	//		if slices.Contains(allowedKeys, key) {
	//			if len(value) > 0 {
	//				options.Search[key] = value[0]
	//			}
	//		}
	//	}
	//	options.Preloads = map[string]db.GetOptions{
	//		"Teams": {}, "Users": {}, "Notices": {},
	//	}
	//	contests, count, ret := db.InitContestRepo(db.DB).List(form.Limit, form.Offset, options)
	//	if !ok {
	//		ctx.JSON(http.StatusOK, ret)
	//		return
	//	}
	//	data := make([]gin.H, 0)
	//	for _, contest := range contests {
	//		data = append(data, resp.GetContestResp(contest, true))
	//	}
	//	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"results": data, "count": count}})
	//	return
	//case "team":
	//	allowedKeys := []string{"name", "id"}
	//	for key, value := range query {
	//		if slices.Contains(allowedKeys, key) {
	//			if len(value) > 0 {
	//				options.Search[key] = value[0]
	//			}
	//		}
	//	}
	//	options.Preloads = map[string]db.GetOptions{"Users": {}}
	//	teams, count, ret := db.InitTeamRepo(db.DB).List(form.Limit, form.Offset, options)
	//	if !ok {
	//		ctx.JSON(http.StatusOK, ret)
	//		return
	//	}
	//	data := make([]gin.H, 0)
	//	for _, team := range teams {
	//		data = append(data, resp.GetTeamResp(team))
	//	}
	//	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"results": data, "count": count}})
	//	return
	//case "challenge":
	//	allowedKeys := []string{"name", "id", "category", "type"}
	//	for key, value := range query {
	//		if slices.Contains(allowedKeys, key) {
	//			if key == "id" {
	//				key = "rand_id"
	//			}
	//			if len(value) > 0 {
	//				options.Search[key] = value[0]
	//			}
	//		}
	//	}
	//	options.Preloads = map[string]db.GetOptions{"ChallengeFlags": {}, "Dockers": {}}
	//	challenges, count, ret := db.InitChallengeRepo(db.DB).List(form.Limit, form.Offset, options)
	//	if !ok {
	//		ctx.JSON(http.StatusOK, ret)
	//		return
	//	}
	//	data := make([]gin.H, 0)
	//	for _, challenge := range challenges {
	//		data = append(data, resp.GetChallengeResp(challenge))
	//	}
	//	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"results": data, "count": count}})
	//	return
	//default:
	//	ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Request.BadRequest})
	//	return
	//}
}
