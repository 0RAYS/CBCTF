package repo

import (
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
	URL       string
	UserAgent string
	Status    int
	Referer   string
	Magic     string
}

func InitRequestRepo(tx *gorm.DB) *RequestRepo {
	return &RequestRepo{Repo: Repo[model.Request]{DB: tx, Model: "Request"}}
}

//func (r *RequestRepo) Init(tx *gorm.DB) RequestRepo {
//	r.DB = tx
//	return *r
//}

//func (r *RequestRepo) Create(options CreateRequestOptions) (model.Request, bool, string) {
//	ip, err := utils.S2S[model.Request](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.Request: %s", err)
//		return model.Request{}, false, "Options2ModelError"
//	}
//	if res := r.DB.Model(&model.Request{}).Create(&ip); res.Error != nil {
//		log.Logger.Warningf("Failed to create Request: %s", res.Error)
//		return model.Request{}, false, "CreateIPError"
//	}
//	return ip, true, "Success"
//}

//func (r *RequestRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Request, bool, string) {
//	switch key {
//	case "id":
//		value = value.(uint)
//	default:
//		return model.Request{}, false, "UnsupportedKey"
//	}
//	var request model.Request
//	res := r.DB.Model(&model.Request{}).Where(key+" = ?", value)
//	res = model.GetPreload(res, model.Request{}, preload, depth).Limit(1).Find(&request)
//	if res.RowsAffected == 0 {
//		return model.Request{}, false, "RequestNotFound"
//	}
//	return request, true, "Success"
//}

//func (r *RequestRepo) GetByID(id uint, preload bool, depth int) (model.Request, bool, string) {
//	return r.getByUniqueKey("id", id, preload, depth)
//}

//func (r *RequestRepo) Count() (int64, bool, string) {
//	var count int64
//	res := r.DB.Model(&model.Request{}).Count(&count)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to count Requests: %s", res.Error)
//		return 0, false, "CountModelError"
//	}
//	return count, true, "Success"
//}

func (r *RequestRepo) CountIP() (int64, bool, string) {
	var count int64
	res := r.DB.Model(&model.Request{}).Distinct("ip").Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Requests: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

//func (r *RequestRepo) GetAll(limit, offset int, preload bool, depth int) ([]model.Request, int64, bool, string) {
//	var (
//		requests       = make([]model.Request, 0)
//		count, ok, msg = r.Count()
//	)
//	if !ok {
//		return requests, count, false, msg
//	}
//	res := r.DB.Model(&model.Request{})
//	res =  model.GetPreload(res, model.Request{}, preload, depth).Limit(limit).Offset(offset).Find(&requests)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to get all Requests: %s", res.Error)
//		return requests, count, false, "GetRequestError"
//	}
//	return requests, count, true, "Success"
//}

//func (r *RequestRepo) Delete(idL ...uint) (bool, string) {
//	res := r.DB.Model(&model.Request{}).Where("id IN ?", idL).Delete(&model.Request{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Request: %s", res.Error)
//		return false, "DeleteRequestError"
//	}
//	return true, "Success"
//}
