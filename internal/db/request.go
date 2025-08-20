package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RequestRepo struct {
	BasicRepo[model.Request]
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
	}
}

type UpdateRequestOptions struct{}

func (u UpdateRequestOptions) Convert2Map() map[string]any {
	return make(map[string]any)
}

func InitRequestRepo(tx *gorm.DB) *RequestRepo {
	return &RequestRepo{
		BasicRepo: BasicRepo[model.Request]{
			DB: tx,
		},
	}
}

func (r *RequestRepo) GetIPByMagic(magic string) ([]string, bool, string) {
	var ipL []string
	res := r.DB.Model(&model.Request{}).Where("magic = ?", magic).Distinct("ip").Find(&ipL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Requests: %s", res.Error)
		return nil, false, i18n.GetRequestError
	}
	if res.RowsAffected == 0 {
		return nil, false, i18n.RequestNotFound
	}
	return ipL, true, i18n.Success
}

func (r *RequestRepo) CountIP() (int64, bool, string) {
	var count int64
	res := r.DB.Model(&model.Request{}).Distinct("ip").Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Reuqest: %s", res.Error)
		return 0, false, i18n.GetRequestError
	}
	return count, true, i18n.Success
}
