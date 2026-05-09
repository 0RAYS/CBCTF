package task

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// LoadTraffic 简单记录涉及到的 IP 地址
func LoadTraffic(tx *gorm.DB, victim model.Victim) model.RetVal {
	trafficRepo := db.InitTrafficRepo(tx)
	optionsL := make(map[string]db.CreateTrafficOptions)
	count, _ := trafficRepo.Count(db.CountOptions{Conditions: map[string]any{"victim_id": victim.ID}})
	if count > 0 {
		log.Logger.Debugf("Traffic already loaded: victim_id=%d records=%d", victim.ID, count)
		return model.SuccessRetVal()
	}
	func(victim model.Victim) {
		start := time.Now()
		log.Logger.Debugf("Enrich pcap with process info: victim_id=%d path=%s", victim.ID, victim.TrafficBasePath())
		if errs := utils.EnrichPcapDir(victim.TrafficBasePath()); len(errs) > 0 {
			for _, err := range errs {
				log.Logger.Warningf("Enrich pcap error: %v", err)
			}
		}
		log.Logger.Debugf("Archiving victim traffic pcaps: victim_id=%d path=%s", victim.ID, victim.TrafficBasePath())
		if err := utils.Zip(victim.TrafficBasePath(), victim.TrafficZipPath()); err != nil {
			log.Logger.Warningf("Failed to archive victim traffic pcaps: victim_id=%d error=%s", victim.ID, err)
			return
		}
		size, hash, err := utils.GetFileInfoByPath(victim.TrafficZipPath())
		if err != nil {
			log.Logger.Warningf("Failed to get traffic archive info: victim_id=%d path=%s error=%s", victim.ID, victim.TrafficZipPath(), err)
			return
		}
		if _, ret := db.InitFileRepo(tx).Create(db.CreateFileOptions{
			RandID:   utils.UUID(),
			Filename: "traffics.zip",
			Size:     size,
			Path:     model.FilePath(victim.TrafficZipPath()),
			Model:    model.Name(victim),
			ModelID:  victim.ID,
			Suffix:   ".zip",
			Hash:     hash,
			Type:     model.TrafficFileType,
		}); !ret.OK {
			log.Logger.Warningf("Failed to create traffic archive record: victim_id=%d reason=%s", victim.ID, ret.Msg)
			return
		}
		log.Logger.Debugf("Archived victim traffic pcaps: victim_id=%d size=%d duration=%s", victim.ID, size, time.Since(start))
	}(victim)
	start := time.Now()
	connections, err := utils.ReadPcapDir(victim.TrafficBasePath())
	if err != nil {
		log.Logger.Warningf("Failed to read victim pcaps: victim_id=%d path=%s error=%s", victim.ID, victim.TrafficBasePath(), err)
		return model.RetVal{Msg: i18n.Model.File.ReadPcapError, Attr: map[string]any{"Error": err.Error()}}
	}
	log.Logger.Debugf("Read victim pcaps: victim_id=%d packets=%d duration=%s", victim.ID, len(connections), time.Since(start))
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
	log.Logger.Infof("Traffic loaded: victim_id=%d packets=%d records=%d", victim.ID, len(connections), len(optionsL))
	return model.SuccessRetVal()
}
