package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type RequestRepo struct {
	BaseRepo[model.Request]
}

func InitRequestRepo(tx *gorm.DB) *RequestRepo {
	return &RequestRepo{
		BaseRepo: BaseRepo[model.Request]{
			DB: tx,
		},
	}
}

func (r *RequestRepo) Create(requests ...model.Request) model.RetVal {
	if len(requests) == 0 {
		return model.SuccessRetVal()
	}
	if res := r.DB.Model(&model.Request{}).CreateInBatches(requests, 200); res.Error != nil {
		log.Logger.Warningf("Failed to create requests: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": model.Name(model.Request{}), "Error": res.Error}}
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
	FirstTime time.Time
	IP        string
	UserID    uint
}

func (r *RequestRepo) ListSharedContestUserIPs(contestID uint, start, end time.Time) ([]UserIP, model.RetVal) {
	if contestID == 0 {
		return nil, model.SuccessRetVal()
	}

	sharedIPs := r.DB.Table("requests").
		Select("requests.ip").
		Joins("INNER JOIN user_contests ON user_contests.user_id = requests.user_id").
		Joins("INNER JOIN users ON users.id = requests.user_id AND users.deleted_at IS NULL").
		Where("user_contests.contest_id = ? AND requests.deleted_at IS NULL", contestID).
		Where("requests.time >= ? AND requests.time <= ?", start, end).
		Group("requests.ip").
		Having("COUNT(DISTINCT requests.user_id) > 1")

	var userIPL []UserIP
	res := r.DB.Table("requests").
		Select("requests.user_id, requests.ip, MIN(requests.time) AS first_time").
		Joins("INNER JOIN user_contests ON user_contests.user_id = requests.user_id").
		Joins("INNER JOIN users ON users.id = requests.user_id AND users.deleted_at IS NULL").
		Where("user_contests.contest_id = ? AND requests.deleted_at IS NULL AND requests.ip IN (?)", contestID, sharedIPs).
		Where("requests.time >= ? AND requests.time <= ?", start, end).
		Group("requests.user_id, requests.ip").
		Order("requests.ip ASC, first_time ASC, requests.user_id ASC").
		Scan(&userIPL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to list shared request IPs: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Request.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return userIPL, model.SuccessRetVal()
}
