package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetRemoteStatus model.Usage 需要预加载
func GetRemoteStatus(tx *gorm.DB, usage model.Usage) gin.H {
	data := gin.H{
		"target":    make([]string, 0),
		"remaining": 0,
		"status":    "Running",
	}
	if usage.Challenge.Type != model.DockerChallenge && usage.Challenge.Type != model.DockersChallenge {
		data["status"] = "NotDocker"
		return data
	}
	repo := db.InitContainerRepo(tx)
	var minTime float64
	for _, container := range usage.Containers {
		_, ok, _ := repo.GetByID(container.ID, false, 0)
		if !ok {
			data["status"] = "Down"
			continue
		}
		if minTime == 0 || minTime > container.Remaining().Seconds() {
			minTime = container.Remaining().Seconds()
		}
		data["target"] = append(data["target"].([]string), container.RemoteAddr()...)
	}
	if len(data["target"].([]string)) > 0 && data["status"] == "Down" {
		data["status"] = "PartDown"
	}
	data["remaining"] = minTime
	return data
}
