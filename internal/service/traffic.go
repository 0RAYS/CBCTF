package service

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	r "CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GetTraffic(victim model.Victim, form f.GetTrafficForm) ([]utils.Connection, []string, int64, bool, string) {
	connections, ok, msg := r.GetTraffic(victim)
	if !ok {
		return nil, nil, 0, false, msg
	}
	if len(connections) < 1 {
		ok, msg = r.UpdateTraffics(victim)
		if !ok {
			return nil, nil, 0, false, msg
		}
		connections, ok, msg = r.GetTraffic(victim)
		if !ok {
			return nil, nil, 0, false, msg
		}
		if len(connections) < 1 {
			return make([]utils.Connection, 0), nil, 0, true, i18n.Success
		}
	}
	totalDuration := int64(connections[len(connections)-1].Time.Sub(connections[0].Time))/1e9 + 1
	ip := make(map[string]bool)
	startIndex := 0
	endIndex := len(connections) - 1
	for i, connection := range connections {
		if _, exists := ip[connection.SrcIP]; !exists {
			ip[connection.SrcIP] = true
		}
		if _, exists := ip[connection.DstIP]; !exists {
			ip[connection.DstIP] = true
		}
		if connection.TimeShift < time.Duration(form.TimeShift*1e9) {
			startIndex = i
		}
		if endIndex == len(connections)-1 && connection.TimeShift > time.Duration((form.TimeShift+form.Duration)*1e9) {
			endIndex = i
		}
	}
	ipL := make([]string, 0)
	for k := range ip {
		ipL = append(ipL, k)
	}
	return connections[startIndex:endIndex], ipL, totalDuration, true, i18n.Success
}

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
			log.Logger.Warningf("Failed to zip .pcap files: %s", err)
		}
	}(victim)
	connections, err := utils.ReadPcapDir(victim.TrafficBasePath())
	if err != nil {
		log.Logger.Warningf("Failed to read pcap: %s", err)
		return false, i18n.ReadPcapError
	}
	for _, conn := range connections {
		connID := fmt.Sprintf("%s:%d-%s:%d-%s-%s", conn.SrcIP, conn.SrcPort, conn.DstIP, conn.DstPort, conn.Type, conn.Subtype)
		if options, exists := optionsL[connID]; exists {
			options.Count += 1
			options.Size += conn.Size
			optionsL[connID] = options
		} else {
			optionsL[connID] = db.CreateTrafficOptions{
				VictimID: victim.ID,
				SrcIP:    conn.SrcIP,
				DstIP:    conn.DstIP,
				SrcPort:  conn.SrcPort,
				DstPort:  conn.DstPort,
				Type:     conn.Type,
				Subtype:  conn.Subtype,
				Size:     conn.Size,
				Count:    1,
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
