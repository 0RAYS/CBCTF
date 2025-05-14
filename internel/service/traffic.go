package service

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/traffic"
	"CBCTF/internel/utils"
	"encoding/hex"
	"gorm.io/gorm"
)

// LoadTraffic model.Victim 需要预加载 model.Pod
func LoadTraffic(tx *gorm.DB, victim model.Victim) (bool, string) {
	repo := db.InitTrafficRepo(tx)
	for _, pod := range victim.Pods {
		_, _, ok, _ := repo.GetByKey("pod_id", pod.ID, 1, 0)
		if ok {
			return true, i18n.Success
		}
		packet, ok, msg := traffic.ReadPcap(pod.TrafficPath())
		if !ok {
			if pod.DeletedAt.Valid && msg == "PcapNotFound" {
				msg = "HasNoTraffic"
			}
			return ok, msg
		}
		for _, conn := range packet {
			_, t := conn.ParsePayload()
			_, ok, msg := repo.Create(db.CreateTrafficOptions{
				VictimID: victim.ID,
				PodID:    pod.ID,
				SrcIP:    conn.SrcIP,
				DstIP:    conn.DstIP,
				SrcPort:  conn.SrcPort,
				DstPort:  conn.DstPort,
				Payload:  hex.EncodeToString(conn.Payload),
				Time:     conn.Time,
				Type:     t,
				Path:     pod.TrafficPath(),
			})
			if !ok {
				return false, msg
			}
		}
	}
	err := utils.Zip(victim.TrafficPaths(), victim.TrafficZipPath())
	if err != nil {
		return false, i18n.ZipError
	}
	return true, i18n.Success
}
