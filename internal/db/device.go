package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"fmt"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

type DiffUpdateDeviceOptions struct {
	Count int
}

func (d DiffUpdateDeviceOptions) Convert2Expr() map[string]clause.Expr {
	options := make(map[string]clause.Expr)
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

var UserDeviceMutex sync.Map

func (d *DeviceRepo) RecordDevice(userID uint, magic, ip string) (model.Device, bool, string) {
	if devices, ok, _ := d.GetByMagic(magic); ok {
		cheatRepo := InitCheatRepo(d.DB)
		for _, device := range devices {
			if userID != device.UserID {
				cheatRepo.Create(CreateCheatOptions{
					UserID:  &userID,
					Magic:   magic,
					IP:      ip,
					Reason:  fmt.Sprintf(model.SameDeviceMagic, userID, device.UserID),
					Type:    model.Suspicious,
					Checked: false,
				})
			}
		}
	}
	if device, ok, msg := d.GetBy2ID(userID, magic); ok {
		mu, _ := UserDeviceMutex.LoadOrStore(userID, &sync.Mutex{})
		mu.(*sync.Mutex).Lock()
		defer mu.(*sync.Mutex).Unlock()
		count := device.Count + 1
		ok, msg = d.Update(device.ID, UpdateDeviceOptions{Count: &count})
		if !ok {
			return model.Device{}, false, msg
		}
		return d.GetByID(device.ID)
	}
	return d.Create(CreateDeviceOptions{UserID: userID, Magic: magic, Count: 1})
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
