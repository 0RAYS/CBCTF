package resp

import (
	"CBCTF/internal/utils"
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"
)

type Statistics struct {
	SrcIP   string
	DstIP   string
	Type    string
	Subtype string
	Count   int64
	Size    int64
}

func GetTrafficResp(connections []utils.Connection) gin.H {
	data := make(map[string]Statistics)
	ipL := make([]string, 0)
	for _, connection := range connections {
		key := fmt.Sprintf("%s-%s-%s-%s", connection.SrcIP, connection.DstIP, connection.Type, connection.Subtype)
		if stats, exists := data[key]; exists {
			stats.Count += 1
			stats.Size += int64(connection.Size)
			data[key] = stats
		} else {
			data[key] = Statistics{
				SrcIP:   connection.SrcIP,
				DstIP:   connection.DstIP,
				Type:    connection.Type,
				Subtype: connection.Subtype,
				Count:   1,
				Size:    int64(connection.Size),
			}
		}
		if !slices.Contains(ipL, connection.SrcIP) {
			ipL = append(ipL, connection.SrcIP)
		}
		if !slices.Contains(ipL, connection.DstIP) {
			ipL = append(ipL, connection.DstIP)
		}
	}
	conn := make([]gin.H, 0)
	for _, stats := range data {
		conn = append(conn, gin.H{
			"src_ip":  stats.SrcIP,
			"dst_ip":  stats.DstIP,
			"type":    stats.Type,
			"subtype": stats.Subtype,
			"count":   stats.Count,
			"size":    stats.Size,
		})
	}
	return gin.H{"connections": conn, "ip": ipL}
}
