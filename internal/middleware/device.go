package middleware

import "sync"

type RequestDevice struct {
	UserID uint
	Magic  string
	Count  int64
}

type requestDeviceKey struct {
	UserID uint
	Magic  string
}

var (
	RequestDevicePool  = make(map[requestDeviceKey]int64)
	RequestDeviceMutex sync.Mutex
)

func RecordRequestDevice(userID uint, magic string, count int64) {
	if count == 0 {
		return
	}

	RequestDeviceMutex.Lock()
	RequestDevicePool[requestDeviceKey{UserID: userID, Magic: magic}] += count
	RequestDeviceMutex.Unlock()
}

func DrainRequestDevicePool() []RequestDevice {
	RequestDeviceMutex.Lock()
	if len(RequestDevicePool) == 0 {
		RequestDeviceMutex.Unlock()
		return nil
	}

	pool := RequestDevicePool
	RequestDevicePool = make(map[requestDeviceKey]int64)
	RequestDeviceMutex.Unlock()

	devices := make([]RequestDevice, 0, len(pool))
	for key, count := range pool {
		devices = append(devices, RequestDevice{
			UserID: key.UserID,
			Magic:  key.Magic,
			Count:  count,
		})
	}
	return devices
}
