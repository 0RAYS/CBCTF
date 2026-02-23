package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

var models = []model.Model{
	model.Challenge{}, model.ChallengeFlag{}, model.Cheat{}, model.Container{}, model.Contest{},
	model.ContestChallenge{}, model.ContestFlag{}, model.Device{}, model.Docker{}, model.Email{}, model.Event{},
	model.File{}, model.Group{}, model.Notice{}, model.Oauth{}, model.Permission{}, model.Pod{}, model.Request{},
	model.Role{}, model.Setting{}, model.Smtp{}, model.Submission{}, model.Team{}, model.TeamFlag{}, model.Traffic{},
	model.User{}, model.Victim{}, model.Webhook{}, model.WebhookHistory{},
}

func GetAllowQueryModels(ctx *gin.Context) {
	data := gin.H{}
	for _, m := range models {
		if len(m.QueryFields()) > 0 {
			data[m.ModelName()] = m.QueryFields()
		}
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}

func Search(ctx *gin.Context) {
	var form dto.SearchModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	var m model.Model
	var fields []string
	var found bool
	for _, m = range models {
		if fields = m.QueryFields(); len(fields) > 0 && m.ModelName() == form.Model {
			found = true
			break
		}
	}
	if !found {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": "Not allowed model"}})
		return
	}
	options := db.GetOptions{Search: make(map[string]string), Sort: make([]string, 0)}
	for key, value := range form.Search {
		if !slices.Contains(fields, key) {
			continue
		}
		options.Search[key] = value
	}
	for key, value := range form.Sort {
		if !slices.Contains(fields, key) {
			continue
		}
		switch strings.ToLower(value) {
		case "asc", "desc":
		default:
			continue
		}
		options.Sort = append(options.Sort, fmt.Sprintf("%s %s", key, strings.ToLower(value)))
	}
	ms, count, ret := db.Search(m, form.Limit, form.Offset, options)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"count": count, "models": ms}))
}
