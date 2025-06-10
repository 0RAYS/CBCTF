package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
	"sync"
)

var UserDeviceMutex sync.Map

func RecordDevice(tx *gorm.DB, userID uint, magic string) (model.Device, bool, string) {
	var (
		repo            = db.InitDeviceRepo(tx)
		device, ok, msg = repo.GetBy2ID(userID, magic)
		count           int
	)
	if ok {
		mu, _ := UserDeviceMutex.LoadOrStore(userID, &sync.Mutex{})
		mu.(*sync.Mutex).Lock()
		defer mu.(*sync.Mutex).Unlock()
		count = device.Count + 1
		ok, msg = repo.Update(device.ID, db.UpdateDeviceOptions{Count: &count})
		if !ok {
			return model.Device{}, false, msg
		}
		return repo.GetByID(device.ID)
	}
	return repo.Create(db.CreateDeviceOptions{UserID: userID, Magic: magic, Count: 1})
}
