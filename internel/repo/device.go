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

func (d *DeviceRepo) GetBy2ID(userID uint, magic string) (model.Device, bool, string) {
	var device model.Device
	res := d.DB.Model(&model.Device{}).Where("user_id = ? AND magic = ?", userID, magic).Limit(1).Find(&device)
	if res.RowsAffected == 0 {
		return model.Device{}, false, "DeviceNotFound"
	}
	return device, true, "Success"
}

func (d *DeviceRepo) Update(id uint, options UpdateDeviceOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Device: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		device, ok, msg := d.GetByID(id, false)
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
