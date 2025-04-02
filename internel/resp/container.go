package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetContainerResp(container model.Container) gin.H {
	return gin.H{
		"id":       container.ID,
		"pod":      container.PodName,
		"start":    container.Start,
		"duration": container.Duration.Seconds(),
		"ip":       container.IP,
		"flags":    container.Flags,
	}
}
