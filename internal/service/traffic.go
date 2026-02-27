package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	r "CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GetTraffic(victim model.Victim, form dto.GetTrafficForm) ([]utils.Connection, []string, int64, model.RetVal) {
	connections, ret := r.GetTraffic(victim)
	if !ret.OK {
		return nil, nil, 0, ret
	}
	if len(connections) < 1 {
		ret = r.UpdateTraffics(victim)
		if !ret.OK {
			return nil, nil, 0, ret
		}
		connections, ret = r.GetTraffic(victim)
		if !ret.OK {
			return nil, nil, 0, ret
		}
		if len(connections) < 1 {
			return make([]utils.Connection, 0), nil, 0, model.SuccessRetVal()
		}
	}
	totalDuration := int64(connections[len(connections)-1].Time.Sub(connections[0].Time))/1e9 + 1
	ip := make(map[string]bool)
	startIndex := 0
	endIndex := len(connections) - 1
	for i, connection := range connections {
		if _, ok := ip[connection.SrcIP]; !ok {
			ip[connection.SrcIP] = true
		}
		if _, ok := ip[connection.DstIP]; !ok {
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
	return connections[startIndex:endIndex], ipL, totalDuration, model.SuccessRetVal()
}

// LoadTraffic 简单记录涉及到的 IP 地址
func LoadTraffic(tx *gorm.DB, victim model.Victim) model.RetVal {
	trafficRepo := db.InitTrafficRepo(tx)
	optionsL := make(map[string]db.CreateTrafficOptions)
	count, _ := trafficRepo.Count(db.CountOptions{Conditions: map[string]any{"victim_id": victim.ID}})
	if count > 0 {
		return model.SuccessRetVal()
	}
	go func(victim model.Victim) {
		if err := utils.Zip(victim.TrafficBasePath(), victim.TrafficZipPath()); err != nil {
			log.Logger.Warningf("Failed to zip .pcap files: %s", err)
			return
		}
		size, hash, err := utils.GetFileInfoByPath(victim.TrafficZipPath())
		if err != nil {
			log.Logger.Warningf("Failed to get file info %s: %s", victim.TrafficZipPath(), err)
			return
		}
		db.InitFileRepo(db.DB).Create(db.CreateFileOptions{
			RandID:   utils.UUID(),
			Filename: "traffics.zip",
			Size:     size,
			Path:     model.FilePath(victim.TrafficZipPath()),
			Model:    victim.ModelName(),
			ModelID:  victim.ID,
			Suffix:   ".zip",
			Hash:     hash,
			Type:     model.TrafficFileType,
		})
	}(victim)
	connections, err := utils.ReadPcapDir(victim.TrafficBasePath())
	if err != nil {
		log.Logger.Warningf("Failed to read pcap: %s", err)
		return model.RetVal{Msg: i18n.Model.File.ReadPcapError, Attr: map[string]any{"Error": err.Error()}}
	}
	for _, conn := range connections {
		connID := fmt.Sprintf("%s-%s-%s-%s", conn.SrcIP, conn.DstIP, conn.Type, conn.Subtype)
		if options, ok := optionsL[connID]; ok {
			options.Count += 1
			options.Size += conn.Size
			optionsL[connID] = options
		} else {
			optionsL[connID] = db.CreateTrafficOptions{
				VictimID: victim.ID,
				SrcIP:    conn.SrcIP,
				DstIP:    conn.DstIP,
				Type:     conn.Type,
				Subtype:  conn.Subtype,
				Size:     conn.Size,
				Count:    1,
			}
		}
	}
	for _, options := range optionsL {
		_, ret := trafficRepo.Create(options)
		if !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal()
}
