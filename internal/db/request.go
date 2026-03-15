package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type RequestRepo struct {
	BaseRepo[model.Request]
}

type CreateRequestOptions struct {
	IP        string
	Time      time.Time
	Method    string
	Path      string
	URL       string
	UserAgent string
	Status    int
	Referer   string
	Magic     string
	UserID    sql.Null[uint]
}

func (c CreateRequestOptions) Convert2Model() model.Model {
	return model.Request{
		IP:        c.IP,
		Time:      c.Time,
		Method:    c.Method,
		Path:      c.Path,
		URL:       c.URL,
		UserAgent: c.UserAgent,
		Status:    c.Status,
		Referer:   c.Referer,
		Magic:     c.Magic,
		UserID:    c.UserID,
	}
}

func InitRequestRepo(tx *gorm.DB) *RequestRepo {
	return &RequestRepo{
		BaseRepo: BaseRepo[model.Request]{
			DB: tx,
		},
	}
}

func (r *RequestRepo) Insert(requests ...model.Request) model.RetVal {
	if res := r.DB.CreateInBatches(requests, len(requests)); res.Error != nil {
		log.Logger.Warningf("Failed to create requests: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": model.Request{}.ModelName(), "Error": res.Error}}
	}
	return model.SuccessRetVal()
}

func (r *RequestRepo) CountIP() (int64, model.RetVal) {
	var count int64
	res := r.DB.Model(&model.Request{}).Distinct("ip").Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Request: %s", res.Error)
		return 0, model.RetVal{Msg: i18n.Model.Request.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return count, model.SuccessRetVal()
}

type UserIP struct {
	UserID    uint
	IP        string
	FirstTime time.Time
}

func (r *RequestRepo) GetUserIP(userIDL ...uint) ([]UserIP, model.RetVal) {
	if len(userIDL) == 0 {
		return nil, model.SuccessRetVal()
	}
	var userIPL []UserIP
	res := r.DB.Model(&model.Request{}).Select("user_id, ip, MIN(time) AS first_time").
		Where("user_id IN ?", userIDL).
		Group("user_id, ip").
		Scan(&userIPL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Request: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Request.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return userIPL, model.SuccessRetVal()
}
