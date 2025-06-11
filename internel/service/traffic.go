package service

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/traffic"
	"CBCTF/internel/utils"
	"fmt"
	"gorm.io/gorm"
)

// LoadTraffic victim 需要预加载 Pods
func LoadTraffic(tx *gorm.DB, victim model.Victim) (bool, string) {
	trafficRepo := db.InitTrafficRepo(tx)
	optionsL := make(map[string]db.CreateTrafficOptions)
	for _, pod := range victim.Pods {
		_, _, ok, _ := trafficRepo.ListWithConditions(1, 0, db.GetOptions{
			{Key: "pod_id", Value: pod.ID, Op: "and"},
		}, false)
		if ok {
			return true, i18n.Success
		}
		packet, ok, msg := traffic.ReadPcap(pod.TrafficPath())
		if !ok {
			if pod.DeletedAt.Valid && msg == i18n.TrafficNotFound {
				msg = i18n.HasNoTraffic
			}
			return ok, msg
		}
		for _, conn := range packet {
			connID := fmt.Sprintf("%s:%d-%s:%d-%s", conn.SrcIP, conn.SrcPort, conn.DstIP, conn.DstPort, conn.Type)
			if options, exists := optionsL[connID]; exists {
				options.Count += 1
				optionsL[connID] = options
			} else {
				optionsL[connID] = db.CreateTrafficOptions{
					VictimID: victim.ID,
					PodID:    pod.ID,
					SrcIP:    conn.SrcIP,
					DstIP:    conn.DstIP,
					SrcPort:  conn.SrcPort,
					DstPort:  conn.DstPort,
					Type:     conn.Type,
					Count:    1,
				}
			}
		}
	}
	for _, options := range optionsL {
		_, ok, msg := trafficRepo.Create(options)
		if !ok {
			return false, msg
		}
	}
	err := utils.Zip(victim.TrafficPaths(), victim.TrafficZipPath())
	if err != nil {
		return false, i18n.ZipError
	}
	return true, i18n.Success
}
