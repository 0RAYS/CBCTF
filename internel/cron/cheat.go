package cron

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"fmt"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

func CheckRequestIP(contest model.Contest) {
	teams, _, ok, _ := db.InitTeamRepo(db.DB).ListWithConditions(-1, -1, db.GetOptions{
		{Key: "contest_id", Value: contest.ID, Op: "and"},
	}, false, "Users", "Users.Devices")
	if !ok {
		return
	}
	ipTeamUserIDMap := make(map[string][][2]uint)
	requestRepo := db.InitRequestRepo(db.DB)
	trafficRepo := db.InitTrafficRepo(db.DB)
	for _, team := range teams {
		for _, user := range team.Users {
			for _, device := range user.Devices {
				ipL, ok, _ := requestRepo.GetIPByMagic(device.Magic)
				if !ok {
					continue
				}
				for _, ip := range ipL {
					if !utils.In([2]uint{team.ID, user.ID}, ipTeamUserIDMap[ip]) {
						ipTeamUserIDMap[ip] = append(ipTeamUserIDMap[ip], [2]uint{team.ID, user.ID})
					}
				}
			}
			traffics, _, ok, _ := trafficRepo.ListWithConditions(-1, -1, db.GetOptions{
				{Key: "contest_id", Value: contest.ID, Op: "and"},
				{Key: "team_id", Value: team.ID, Op: "and"},
				{Key: "user_id", Value: user.ID, Op: "and"},
			}, false)
			if !ok {
				continue
			}
			for _, traffic := range traffics {
				if !utils.In([2]uint{team.ID, user.ID}, ipTeamUserIDMap[traffic.SrcIP]) {
					ipTeamUserIDMap[traffic.SrcIP] = append(ipTeamUserIDMap[traffic.SrcIP], [2]uint{team.ID, user.ID})
				}
			}
		}
	}
	cheatRepo := db.InitCheatRepo(db.DB)
	for ip, teamUserIDL := range ipTeamUserIDMap {
		if len(teamUserIDL) > 1 {
			var tmp strings.Builder
			for _, teamUserID := range teamUserIDL {
				tmp.WriteString(fmt.Sprintf("Team%d-User%d, ", teamUserID[0], teamUserID[1]))
			}
			for _, teamUserID := range teamUserIDL {
				teamID, userID := teamUserID[0], teamUserID[1]
				cheatRepo.Create(db.CreateCheatOptions{
					UserID:    &teamID,
					TeamID:    &userID,
					ContestID: &contest.ID,
					IP:        ip,
					Reason:    fmt.Sprintf(model.SameRequestIP, strings.Trim(tmp.String(), ", "), ip),
					Type:      model.Suspicious,
					Checked:   false,
				})
			}
		}
	}
}

func CheckCheat() {
	function := exec("CheckCheat", func() {
		contests, _, ok, _ := db.InitContestRepo(db.DB).List(-1, -1)
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				continue
			}
			go CheckRequestIP(contest)
		}
	})
	function()
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(function))
}
