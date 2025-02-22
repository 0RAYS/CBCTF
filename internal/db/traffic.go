package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/traffic"
	"CBCTF/internal/utils"
	"fmt"
	"gorm.io/gorm"
)

func SaveTraffic(tx *gorm.DB, docker model.Docker) (bool, string) {
	res := tx.Model(&model.Traffic{}).Where("path = ?", docker.TrafficPath()).Find(&model.Traffic{}).Limit(1)
	if res.RowsAffected > 0 || res.Error != nil {
		return true, "Success"
	}
	connections, ok, msg := traffic.ReadPcap(docker.TrafficPath())
	if !ok {
		if docker.DeletedAt.Valid && msg == "PcapNotFound" {
			msg = "HasNoTraffic"
		}
		return ok, msg
	}
	for _, conn := range connections {
		t := model.InitTraffic(conn, docker)
		res := tx.Model(&model.Traffic{}).Create(&t)
		if res.Error != nil {
			log.Logger.Warningf("Failed to save traffic: %s", res.Error)
			return false, "SaveTrafficError"
		}
	}
	return true, "Success"
}

func getTrafficByID(tx *gorm.DB, column string, id uint, limit, offset int) ([]model.Traffic, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	var traffics []model.Traffic
	var count int64
	res := tx.Model(&model.Traffic{}).Where(fmt.Sprintf("%s = ?", column), id)
	if err := res.Count(&count).Error; err != nil {
		log.Logger.Warningf("Failed to count traffic: %s", err)
		return make([]model.Traffic, 0), -1, false, "UnknownError"
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	res = res.Limit(limit).Offset(offset).Find(&traffics)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get traffic: %s", res.Error)
		return make([]model.Traffic, 0), -1, false, "GetTrafficError"
	}
	return traffics, count, true, ""
}

func GetTrafficByDocker(tx *gorm.DB, dockerID uint, limit, offset int) ([]model.Traffic, int64, bool, string) {
	return getTrafficByID(tx, "docker_id", dockerID, limit, offset)
}

func GetTrafficByTeam(tx *gorm.DB, teamID uint, limit, offset int) ([]model.Traffic, int64, bool, string) {
	return getTrafficByID(tx, "team_id", teamID, limit, offset)
}
