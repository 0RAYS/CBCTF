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
		BasicRepo: BasicRepo[model.Request]{
			DB: tx,
		},
	}
}

func (r *RequestRepo) GetByMagic(magic string, optionsL ...GetOptions) ([]model.Request, int64, bool, string) {
	options := GetOptions{}
	if len(optionsL) > 0 {
		options = optionsL[0]
	}
	if len(options.Conditions) == 0 {
		options.Conditions = make(map[string]any)
	}
	options.Conditions["magic"] = magic
	return r.List(-1, -1, options)
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

func (r *RequestRepo) GetUserIP(userID uint) ([]string, bool, string) {
	var ipL []string
	res := r.DB.Model(&model.Request{}).Distinct("ip").Where("user_id = ?", userID).Find("ip", &ipL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Reuqest: %s", res.Error)
		return nil, false, i18n.GetRequestError
	}
	return ipL, true, i18n.Success
}
