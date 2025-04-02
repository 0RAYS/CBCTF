package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type DeviceRepo struct {
	Repo[model.Device]
}

type CreateDeviceOptions struct {
	UserID uint
	Magic  string
	Count  int
}

type UpdateDeviceOptions struct {
	Count *int `json:"count"`
}

func InitDeviceRepo(tx *gorm.DB) *DeviceRepo {
	return &DeviceRepo{Repo: Repo[model.Device]{DB: tx, Model: "Device"}}
}

//func (d *DeviceRepo) Create(options CreateDeviceOptions) (model.Device, bool, string) {
//	device, err := utils.S2S[model.Device](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.Device: %s", err)
//		return model.Device{}, false, "Options2ModelError"
//	}
//	if res := d.DB.Model(&model.Device{}).Create(&device); res.Error != nil {
//		log.Logger.Warningf("Failed to create Device: %s", res.Error)
//		return model.Device{}, false, "CreateDeviceError"
//	}
//	return device, true, "Success"
//}

//func (d *DeviceRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Device, bool, string) {
//	switch key {
//	case "id":
//		value = value.(uint)
//	default:
//		return model.Device{}, false, "UnsupportedKey"
//	}
//	var device model.Device
//	res := d.DB.Model(&model.Device{}).Where(key+" = ?", value)
//	res = model.GetPreload(res, model.Device{}, preload, depth).Find(&device).Limit(1)
//	if res.RowsAffected == 0 {
//		return model.Device{}, false, "DeviceNotFound"
//	}
//	return device, true, "Success"
//}

//func (d *DeviceRepo) GetByID(id uint, preload bool, depth int) (model.Device, bool, string) {
//	return d.getByUniqueKey("id", id, preload, depth)
//}

func (d *DeviceRepo) GetBy2ID(userID uint, magic string) (model.Device, bool, string) {
	var device model.Device
	res := d.DB.Model(&model.Device{}).Where("user_id = ? AND magic = ?", userID, magic).Find(&device).Limit(1)
	if res.RowsAffected == 0 {
		return model.Device{}, false, "DeviceNotFound"
	}
	return device, true, "Success"
}

//func (d *DeviceRepo) Count() (int64, bool, string) {
//	var count int64
//	res := d.DB.Model(&model.Device{}).Count(&count)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to count Devices: %s", res.Error)
//		return 0, false, "CountModelError"
//	}
//	return count, true, "Success"
//}

//func (d *DeviceRepo) GetAll(limit, offset int, preload bool, depth int) ([]model.Device, int64, bool, string) {
//	var (
//		devices        = make([]model.Device, 0)
//		count, ok, msg = d.Count()
//	)
//	if !ok {
//		return devices, count, false, msg
//	}
//	res := d.DB.Model(&model.Device{})
//	res = model.GetPreload(res, model.Device{}, preload, depth).Find(&devices).Limit(limit).Offset(offset)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to get all Devices: %s", res.Error)
//		return devices, count, false, "GetDeviceError"
//	}
//	return devices, count, true, "Success"
//}

func (d *DeviceRepo) Update(id uint, options UpdateDeviceOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Device: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		device, ok, msg := d.GetByID(id, false, 0)
		if !ok {
			return false, msg
		}
		data["version"] = device.Version + 1
		res := d.DB.Model(&model.Device{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, device.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Device: %s", res.Error)
			return false, "UpdateDeviceError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}

//func (d *DeviceRepo) Delete(idL ...uint) (bool, string) {
//	res := d.DB.Model(&model.Device{}).Where("id IN ?", idL).Delete(&model.Device{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Device: %s", res.Error)
//		return false, "DeleteDeviceError"
//	}
//	return true, "Success"
//}
