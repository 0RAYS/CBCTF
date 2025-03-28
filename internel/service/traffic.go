package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/traffic"
	"encoding/hex"
	"gorm.io/gorm"
)

// SaveTraffic 从 .pcap 文件中保存流量至数据库
func SaveTraffic(tx *gorm.DB, container model.Container) (bool, string) {
	repo := db.InitTrafficRepo(tx)
	_, ok, _ := repo.GetByID(container.ID, false, 0)
	if ok {
		return true, "Success"
	}
	connections, ok, msg := traffic.ReadPcap(container.TrafficPath())
	if !ok {
		if container.DeletedAt.Valid && msg == "PcapNotFound" {
			msg = "HasNoTraffic"
		}
		return ok, msg
	}
	for _, conn := range connections {
		_, t := conn.ParsePayload()
		_, ok, msg := repo.Create(db.CreateTrafficOptions{
			SrcIP:       conn.SrcIP,
			DstIP:       conn.DstIP,
			SrcPort:     conn.SrcPort,
			DstPort:     conn.DstPort,
			Payload:     hex.EncodeToString(conn.Payload),
			Type:        t,
			ContainerID: container.ID,
		})
		if !ok {
			return false, msg
		}
	}
	return true, "Success"
}
