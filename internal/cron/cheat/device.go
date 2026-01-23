package cheat

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func CheckSameDevice(contest model.Contest) {
	userIDL, ret := db.GetUserIDByContestID(db.DB, contest.ID)
	if !ret.OK {
		return
	}
	type tmp struct {
		Time   time.Time
		UserID uint
	}
	deviceUserMap := make(map[string][]tmp)
	for _, userID := range userIDL {
		devices, _, ret := db.InitDeviceRepo(db.DB).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"user_id": userID},
		})
		if !ret.OK {
			continue
		}
		for _, device := range devices {
			deviceUserMap[device.Magic] = append(deviceUserMap[device.Magic], tmp{device.CreatedAt, userID})
		}
	}
	repo := db.InitCheatRepo(db.DB)
	for magic, users := range deviceUserMap {
		if len(users) > 1 {
			var str []string
			for _, user := range users {
				str = append(str, strconv.Itoa(int(user.UserID)))
			}
			for _, user := range users {
				repo.Create(db.CreateCheatOptions{
					UserID:    sql.Null[uint]{V: user.UserID, Valid: true},
					ContestID: sql.Null[uint]{V: contest.ID, Valid: true},
					Magic:     magic,
					Reason:    fmt.Sprintf(model.SameDeviceMagic, fmt.Sprintf("User %s", strings.Join(str, ","))),
					Type:      model.Suspicious,
					Time:      user.Time,
				})
			}
		}
	}
}
