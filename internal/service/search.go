package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/view"
	"fmt"
	"slices"

	"gorm.io/gorm"
)

var searchableModels = []model.Model{
	model.Challenge{}, model.Cheat{}, model.Contest{}, model.ContestChallenge{}, model.Device{}, model.Email{},
	model.Event{}, model.File{}, model.Group{}, model.Notice{}, model.Oauth{}, model.Permission{}, model.Request{},
	model.Role{}, model.Setting{}, model.Smtp{}, model.Submission{}, model.Team{}, model.TeamFlag{}, model.Traffic{},
	model.User{}, model.Victim{}, model.Webhook{}, model.WebhookHistory{},
}

func GetAllowQueryModels() map[string]view.SearchModelView {
	data := make(map[string]view.SearchModelView)
	for _, m := range searchableModels {
		if fields := model.QueryFields(m); len(fields) > 0 {
			data[model.ModelName(m)] = view.SearchModelView{
				Query:  fields,
				Search: model.SearchFields(m),
			}
		}
	}
	return data
}

func SearchModels(_ *gorm.DB, form dto.SearchModelsForm) (any, int64, model.RetVal) {
	var (
		m            model.Model
		queryFields  []string
		searchFields []string
		found        bool
	)
	for _, m = range searchableModels {
		if queryFields = model.QueryFields(m); len(queryFields) > 0 && model.ModelName(m) == form.Model {
			searchFields = model.SearchFields(m)
			found = true
			break
		}
	}
	if !found {
		return nil, 0, model.RetVal{Msg: i18n.Response.BadRequest, Attr: map[string]any{"Error": "Not allowed model"}}
	}
	options := db.GetOptions{Search: make(map[string]string), Sort: make([]string, 0)}
	for key, value := range form.Search {
		if !slices.Contains(searchFields, key) {
			continue
		}
		options.Search[key] = value
	}
	for key, value := range form.Sort {
		if !slices.Contains(queryFields, key) {
			continue
		}
		switch value {
		case "asc", "desc":
		default:
			continue
		}
		options.Sort = append(options.Sort, fmt.Sprintf("%s %s", key, value))
	}
	return db.Search(m, form.Limit, form.Offset, options)
}
