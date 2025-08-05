package service

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"CBCTF/internal/utils"
	"fmt"
	"gorm.io/gorm"
	"os"
	"strings"
)

// LoadTraffic 简单记录涉及到的 IP 地址
func LoadTraffic(tx *gorm.DB, victim model.Victim) (bool, string) {
	trafficRepo := db.InitTrafficRepo(tx)
	optionsL := make(map[string]db.CreateTrafficOptions)
	count, _, _ := trafficRepo.Count(db.CountOptions{Conditions: map[string]any{"victim_id": victim.ID}})
	if count > 0 {
		return true, i18n.Success
	}
	go func(victim model.Victim) {
		err := utils.Zip(victim.TrafficBasePath(), victim.TrafficZipPath())
		if err != nil {
			log.Logger.Warningf("Failed to zip .pcap files: %v", err)
		}
	}(victim)
	dir, err := os.ReadDir(victim.TrafficBasePath())
	if err != nil {
		log.Logger.Warningf("Failed to read dir: %v", err)
		return false, ""
	}
	for _, file := range dir {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".pcap") || !strings.HasSuffix(file.Name(), ".pcapng") {
			continue
		}
		packet, err := utils.ReadPcap(fmt.Sprintf("%s/%s", victim.TrafficBasePath(), file.Name()))
		if err != nil {
			if os.IsNotExist(err) {
				return false, i18n.PcapNotFound
			}
			log.Logger.Warningf("Failed to read pcap file %s: %s", file.Name(), err)
			return false, i18n.UnknownError
		}
		for _, conn := range packet {
			connID := fmt.Sprintf("%s:%d-%s:%d-%s", conn.SrcIP, conn.SrcPort, conn.DstIP, conn.DstPort, conn.Type)
			if options, exists := optionsL[connID]; exists {
				options.Count += 1
				optionsL[connID] = options
			} else {
				optionsL[connID] = db.CreateTrafficOptions{
					VictimID: victim.ID,
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
