package cheat

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"fmt"
	"strconv"
	"strings"
)

func CheckSameDevice(contest model.Contest) {
	userIDL, ret := db.GetUserIDByContestID(db.DB, contest.ID)
	if !ret.OK {
		return
	}
	deviceUserMap := make(map[string][]uint)
	devices, _, ret := db.InitDeviceRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"user_id": userIDL},
	})
	if !ret.OK {
		return
	}
	for _, device := range devices {
		if !slicesContains(deviceUserMap[device.Magic], device.UserID) {
			deviceUserMap[device.Magic] = append(deviceUserMap[device.Magic], device.UserID)
		}
	}
	repo := db.InitCheatRepo(db.DB)
	for magic, userIDs := range deviceUserMap {
		if len(userIDs) > 1 {
			var str []string
			for _, uid := range userIDs {
				str = append(str, strconv.Itoa(int(uid)))
			}
			repo.Create(db.CreateCheatOptions{
				ContestID:  contest.ID,
				Model:      model.CheatRefModel{model.User{}.ModelName(): userIDs},
				Magic:      magic,
				Reason:     fmt.Sprintf(model.SameDeviceMagic, fmt.Sprintf("User %s", strings.Join(str, ","))),
				ReasonType: model.ReasonTypeSameDevice,
				Type:       model.Suspicious,
				Time:       devices[0].CreatedAt,
			})
		}
	}
}

func slicesContains(s []uint, v uint) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
