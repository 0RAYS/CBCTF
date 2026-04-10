package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func SearchIP(ctx *gin.Context) {
	var form dto.SearchIP
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ip, err := utils.SearchIP(form.IP, config.Env.GeoCityDB)
	if err != nil {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
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
	resp.JSON(ctx, model.SuccessRetVal(data))
}
