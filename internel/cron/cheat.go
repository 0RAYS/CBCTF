package cron

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"fmt"
	"github.com/robfig/cron/v3"
	"net"
	"strings"
	"time"
)

func CheckRequestIP(contest model.Contest) {
	_, podCIDR, err := net.ParseCIDR(config.Env.K8S.IPPool.CIDR)
	if err != nil {
		log.Logger.Warningf("Failed to parse Pod IPPool CIDR: %v", err)
		return
	}
	teams, _, ok, _ := db.InitTeamRepo(db.DB).ListWithConditions(-1, -1, db.GetOptions{
		{Key: "contest_id", Value: contest.ID, Op: "and"},
	}, false, "Users", "Users.Devices")
	if !ok {
		return
	}
	ipTeamIDMap := make(map[string][]uint)
	requestRepo := db.InitRequestRepo(db.DB)
	victimRepo := db.InitVictimRepo(db.DB)
	for _, team := range teams {
		for _, user := range team.Users {
			for _, device := range user.Devices {
				ipL, ok, _ := requestRepo.GetIPByMagic(device.Magic)
				if !ok {
					continue
				}
				for _, ip := range ipL {
					netIP := net.ParseIP(ip)
					if netIP == nil {
						continue
					}
					if netIP.IsLoopback() {
						continue
					}
					if podCIDR.Contains(netIP) {
						continue
					}
					if !utils.In(team.ID, ipTeamIDMap[ip]) {
						ipTeamIDMap[ip] = append(ipTeamIDMap[ip], team.ID)
					}
				}
			}
		}
		victims, _, ok, _ := victimRepo.ListWithConditions(-1, -1, db.GetOptions{
			{Key: "team_id", Value: team.ID, Op: "and"},
		}, true, "Traffics")
		if !ok {
			continue
		}
		for _, victim := range victims {
			for _, traffics := range victim.Traffics {
				netIP := net.ParseIP(traffics.SrcIP)
				if netIP == nil {
					continue
				}
				if netIP.IsLoopback() {
					continue
				}
				if podCIDR.Contains(netIP) {
					continue
				}
				if !utils.In(victim.TeamID, ipTeamIDMap[traffics.SrcIP]) {
					ipTeamIDMap[traffics.SrcIP] = append(ipTeamIDMap[traffics.SrcIP], victim.TeamID)
				}
			}
		}
	}
	cheatRepo := db.InitCheatRepo(db.DB)
	for ip, teamIDL := range ipTeamIDMap {
		if len(teamIDL) > 1 {
			var tmp strings.Builder
			for _, teamID := range teamIDL {
				tmp.WriteString(fmt.Sprintf("Team-%d, ", teamID))
			}
			for _, teamID := range teamIDL {
				cheatRepo.Create(db.CreateCheatOptions{
					TeamID:    &teamID,
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
