package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetTrafficResp(traffic model.Traffic) gin.H {
	return gin.H{
		"id":       traffic.ID,
		"pod_id":   traffic.PodID,
		"src_ip":   traffic.SrcIP,
		"dst_ip":   traffic.DstIP,
		"src_port": traffic.SrcPort,
		"dst_port": traffic.DstPort,
		"payload":  traffic.Payload,
		"time":     traffic.Time,
		"type":     traffic.Type,
	}
}
