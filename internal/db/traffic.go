package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/traffic"
	"fmt"
	"gorm.io/gorm"
)

func SaveTraffic(tx *gorm.DB, docker model.Docker) (bool, string) {
	connections, ok, msg := traffic.ReadPcap(docker.TrafficPath())
	if !ok {
		return ok, msg
	}
	for _, conn := range connections {
		res := tx.Model(&model.Traffic{}).Create(model.InitTraffic(conn, docker))
		if res.Error != nil {
			log.Logger.Warningf("Failed to save traffic: %s", res.Error)
			return false, "SaveTrafficError"
		}
	}
	return true, "Success"
}

func getTrafficByID(tx *gorm.DB, column string, id uint, limit, offset int) ([]model.Traffic, int64, bool, string) {
	if limit < 0 {
		limit = -1
	}
	if offset < 0 {
		offset = -1
	}
	var traffics []model.Traffic
	var count int64
	res := tx.Model(&model.Traffic{}).Where(fmt.Sprintf("%s = ?", column), id)
	if err := res.Count(&count).Error; err != nil {
		log.Logger.Warningf("Failed to count traffic: %s", err)
		return []model.Traffic{}, -1, false, "UnknownError"
	}
	res = res.Limit(limit).Offset(offset).Find(&traffics)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get traffic: %s", res.Error)
		return []model.Traffic{}, -1, false, "GetTrafficError"
	}
	return traffics, count, true, ""
}

func GetTrafficByDocker(tx *gorm.DB, dockerID uint, limit, offset int) ([]model.Traffic, int64, bool, string) {
	return getTrafficByID(tx, "docker_id", dockerID, limit, offset)
}

func GetTrafficByTeam(tx *gorm.DB, teamID uint, limit, offset int) ([]model.Traffic, int64, bool, string) {
	return getTrafficByID(tx, "team_id", teamID, limit, offset)
}
