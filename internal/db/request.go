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

func (r *RequestRepo) CountIP() (int64, model.RetVal) {
	var count int64
	res := r.DB.Model(&model.Request{}).Distinct("ip").Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Request: %s", res.Error)
		return 0, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.Request{}.ModelName(), "Error": res.Error.Error()}}
	}
	return count, model.SuccessRetVal()
}

func (r *RequestRepo) GetUserIP(userID uint) ([]string, model.RetVal) {
	var ipL []string
	res := r.DB.Model(&model.Request{}).Where("user_id = ?", userID).Distinct("ip").Find(&ipL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Request: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.Request{}.ModelName(), "Error": res.Error.Error()}}
	}
	return ipL, model.SuccessRetVal()
}
