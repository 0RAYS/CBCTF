package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"

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
