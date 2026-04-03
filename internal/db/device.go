package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type DeviceRepo struct {
	BaseRepo[model.Device]
}

type CreateDeviceOptions struct {
	UserID uint
	Magic  string
	Count  int64
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
		BaseRepo: BaseRepo[model.Device]{
			DB: tx,
		},
	}
}

type ContestDeviceUser struct {
	Magic     string
	UserID    uint
	FirstTime time.Time
}

func (d *DeviceRepo) ListSharedContestDevices(contestID uint, start, end time.Time) ([]ContestDeviceUser, model.RetVal) {
	if contestID == 0 {
		return nil, model.SuccessRetVal()
	}

	sharedMagics := d.DB.Table("devices").
		Select("devices.magic").
		Joins("INNER JOIN user_contests ON user_contests.user_id = devices.user_id").
		Joins("INNER JOIN users ON users.id = devices.user_id AND users.deleted_at IS NULL").
		Where("user_contests.contest_id = ? AND devices.deleted_at IS NULL", contestID).
		Where("devices.created_at <= ? AND devices.updated_at >= ?", end, start).
		Group("devices.magic").
		Having("COUNT(DISTINCT devices.user_id) > 1")

	var rows []ContestDeviceUser
	res := d.DB.Table("devices").
		Select("devices.magic, devices.user_id, MIN(devices.created_at) AS first_time").
		Joins("INNER JOIN user_contests ON user_contests.user_id = devices.user_id").
		Joins("INNER JOIN users ON users.id = devices.user_id AND users.deleted_at IS NULL").
		Where("user_contests.contest_id = ? AND devices.deleted_at IS NULL AND devices.magic IN (?)", contestID, sharedMagics).
		Where("devices.created_at <= ? AND devices.updated_at >= ?", end, start).
		Group("devices.magic, devices.user_id").
		Order("devices.magic ASC, first_time ASC, devices.user_id ASC").
		Scan(&rows)
	if res.Error != nil {
		log.Logger.Warningf("Failed to list shared contest devices: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.ModelName(model.Device{}), "Error": res.Error.Error()}}
	}
	return rows, model.SuccessRetVal()
}

func (d *DeviceRepo) RecordDevice(options CreateDeviceOptions, count int64) model.RetVal {
	if count == 0 {
		return model.SuccessRetVal()
	}

	device := options.Convert2Model().(model.Device)
	res := d.DB.Model(&model.Device{}).FirstOrCreate(&device, device)
	if res.Error != nil {
		return model.RetVal{Msg: i18n.Model.Device.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return d.DiffUpdate(device.ID, DiffUpdateDeviceOptions{Count: count})
}
