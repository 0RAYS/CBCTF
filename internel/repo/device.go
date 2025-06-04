package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type DeviceRepo struct {
	Basic[model.Device]
}

type CreateDeviceOptions struct {
	UserID uint
	Magic  string
	Count  int
}

func (c CreateDeviceOptions) Convert2Model() model.Model {
	return model.Device{
		UserID: c.UserID,
		Magic:  c.Magic,
		Count:  c.Count,
	}
}

type UpdateDeviceOptions struct {
	Count *int
}

func (u UpdateDeviceOptions) Convert2Map() map[string]any {
	data := make(map[string]any)
	if u.Count != nil {
		data["count"] = *u.Count
	}
	return data
}

func InitDeviceRepo(tx *gorm.DB) *DeviceRepo {
	return &DeviceRepo{
		Basic: Basic[model.Device]{
			DB: tx,
		},
	}
}

func (d *DeviceRepo) GetBy2ID(userID uint, magic string) (model.Device, bool, string) {
	var device model.Device
	res := d.DB.Model(&model.Device{}).Where("user_id = ? AND magic = ?", userID, magic).Limit(1).Find(&device)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Device: %s", res.Error)
		return model.Device{}, false, model.Device{}.GetErrorString()
	}
	if res.RowsAffected == 0 {
		return model.Device{}, false, model.Device{}.NotFoundErrorString()
	}
	return device, true, i18n.Success
}
