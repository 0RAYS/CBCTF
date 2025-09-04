package db

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

type DiffUpdateDeviceOptions struct {
	Count int64
}

func (d DiffUpdateDeviceOptions) Convert2Expr() map[string]any {
	options := make(map[string]any)
	if d.Count != 0 {
		options["count"] = gorm.Expr("count + ?", d.Count)
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

func (d *DeviceRepo) RecordDevice(options CreateDeviceOptions) (bool, string) {
	device := options.Convert2Model().(model.Device)
	res := d.DB.Model(&model.Device{}).FirstOrCreate(&device, device)
	if res.Error != nil {
		return false, i18n.GetDeviceError
	}
	return d.DiffUpdate(device.ID, DiffUpdateDeviceOptions{Count: 1})
}

func (d *DeviceRepo) GetByMagic(magic string) ([]model.Device, bool, string) {
	var devices []model.Device
	res := d.DB.Model(&model.Device{}).Where("magic = ?", magic).Order("id").Find(&devices)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Devices: %s", res.Error)
		return nil, false, i18n.GetDeviceError
	}
	if res.RowsAffected == 0 {
		return nil, false, i18n.DeviceNotFound
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
