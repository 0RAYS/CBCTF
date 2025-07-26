package service

import (
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"fmt"
	"gorm.io/gorm"
	"sync"
)

var UserDeviceMutex sync.Map

func RecordDevice(tx *gorm.DB, userID uint, magic string, ip string) (model.Device, bool, string) {
	deviceRepo := db.InitDeviceRepo(tx)
	if devices, ok, _ := deviceRepo.GetByMagic(magic); ok {
		cheatRepo := db.InitCheatRepo(tx)
		for _, device := range devices {
			if userID != device.UserID {
				cheatRepo.Create(db.CreateCheatOptions{
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
	if device, ok, msg := deviceRepo.GetBy2ID(userID, magic); ok {
		mu, _ := UserDeviceMutex.LoadOrStore(userID, &sync.Mutex{})
		mu.(*sync.Mutex).Lock()
		defer mu.(*sync.Mutex).Unlock()
		count := device.Count + 1
		ok, msg = deviceRepo.Update(device.ID, db.UpdateDeviceOptions{Count: &count})
		if !ok {
			return model.Device{}, false, msg
		}
		return deviceRepo.GetByID(device.ID)
	}
	return deviceRepo.Create(db.CreateDeviceOptions{UserID: userID, Magic: magic, Count: 1})
}
