package model

import "CBCTF/internel/i18n"

// Device
// BelongsTo User
// HasMany Request
type Device struct {
	UserID uint   `json:"user_id"`
	User   User   `json:"-"`
	Magic  string `json:"magic"`
	Count  int    `json:"count"`
	BasicModel
}

func (d Device) GetModelName() string {
	return "Device"
}

func (d Device) GetVersion() uint {
	return d.Version
}

func (d Device) CreateErrorString() string {
	return i18n.CreateDeviceError
}

func (d Device) DeleteErrorString() string {
	return i18n.DeleteDeviceError
}

func (d Device) GetErrorString() string {
	return i18n.GetDeviceError
}

func (d Device) NotFoundErrorString() string {
	return i18n.DeviceNotFound
}

func (d Device) UpdateErrorString() string {
	return i18n.UpdateDeviceError
}

func (d Device) GetUniqueKey() []string {
	return []string{"id"}
}

func (d Device) GetForeignKeys() []string {
	return []string{"id", "user_id"}
}
