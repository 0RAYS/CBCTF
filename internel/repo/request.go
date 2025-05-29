package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
	"time"
)

type RequestRepo struct {
	Repo[model.Request]
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

func InitRequestRepo(tx *gorm.DB) *RequestRepo {
	return &RequestRepo{
		Repo: Repo[model.Request]{
			DB: tx, Model: "Request",
			CreateError:   i18n.CreateRequestError,
			DeleteError:   i18n.DeleteRequestError,
			GetError:      i18n.GetRequestError,
			NotFoundError: i18n.RequestNotFound,
		},
	}
}

func (r *RequestRepo) CountIP() (int64, bool, string) {
	var count int64
	res := r.DB.Model(&model.Request{}).Distinct("ip").Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Requests: %s", res.Error)
		return 0, false, i18n.CountModelError
	}
	return count, true, i18n.Success
}
