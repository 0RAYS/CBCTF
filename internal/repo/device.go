package repo

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type DeviceRepo struct {
	BasicRepo[model.Device]
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
	options := make(map[string]any)
	if u.Count != nil {
		options["count"] = *u.Count
	}
	return options
}

func InitDeviceRepo(tx *gorm.DB) *DeviceRepo {
	return &DeviceRepo{
		BasicRepo: BasicRepo[model.Device]{
			DB: tx,
		},
	}
}

func (d *DeviceRepo) GetByMagic(magic string) ([]model.Device, bool, string) {
	var devices []model.Device
	res := d.DB.Model(&model.Device{}).Where("magic = ?", magic).Order("id").Find(&devices)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Devices: %s", res.Error)
		return make([]model.Device, 0), false, i18n.GetDeviceError
	}
	if res.RowsAffected == 0 {
		return make([]model.Device, 0), false, i18n.DeviceNotFound
	}
	return devices, true, i18n.Success
}

func (d *DeviceRepo) GetBy2ID(userID uint, magic string) (model.Device, bool, string) {
	var device model.Device
	res := d.DB.Model(&model.Device{}).Where("user_id = ? AND magic = ?", userID, magic).Limit(1).Find(&device)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Device: %s", res.Error)
		return model.Device{}, false, i18n.GetDeviceError
	}
	if res.RowsAffected == 0 {
		return model.Device{}, false, i18n.DeviceNotFound
	}
	return device, true, i18n.Success
}
