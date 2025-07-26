package service

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"CBCTF/internal/traffic"
	"CBCTF/internal/utils"
	"fmt"
	"gorm.io/gorm"
)

// LoadTraffic victim 需要预加载 Pods
func LoadTraffic(tx *gorm.DB, victim model.Victim) (bool, string) {
	trafficRepo := db.InitTrafficRepo(tx)
	optionsL := make(map[string]db.CreateTrafficOptions)
	go func(victim model.Victim) {
		err := utils.Zip(victim.TrafficPaths(), victim.TrafficZipPath())
		if err != nil {
			log.Logger.Warningf("Failed to zip .pcap files: %v", err)
		}
	}(victim)
	for _, pod := range victim.Pods {
		_, count, _, _ := trafficRepo.List(1, 0, db.GetOptions{
			Conditions: map[string]any{
				"pod_id": pod.ID,
			},
		})
		if count > 0 {
			return true, i18n.Success
		}
		packet, ok, msg := traffic.ReadPcap(pod.TrafficPcapPath())
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
	return true, i18n.Success
}
