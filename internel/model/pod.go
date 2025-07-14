package model

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"slices"
)

type Pod struct {
	VictimID     uint        `json:"victim_id"`
	Victim       Victim      `json:"-"`
	Containers   []Container `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Traffics     []Traffic   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name         string      `json:"name"`
	ExposedIP    string      `json:"exposed_ip"`
	PodPorts     Exposes     `gorm:"type:json" json:"pod_ports"`
	ExposedPorts Int32List   `gorm:"type:json" json:"exposed_ports"`
	IPs          IPs         `gorm:"default:null;type:json" json:"ips"`
	BasicModel
}

func (p Pod) GetModelName() string {
	return "Pod"
}

func (p Pod) GetVersion() uint {
	return p.Version
}

func (p Pod) CreateErrorString() string {
	return i18n.CreatePodError
}

func (p Pod) DeleteErrorString() string {
	return i18n.DeletePodError
}

func (p Pod) GetErrorString() string {
	return i18n.GetPodError
}

func (p Pod) NotFoundErrorString() string {
	return i18n.PodNotFound
}

func (p Pod) UpdateErrorString() string {
	return i18n.UpdatePodError
}

func (p Pod) GetUniqueKey() []string {
	return []string{"id"}
}

func (p Pod) GetForeignKeys() []string {
	return []string{"id", "victim_id"}
}

func (p Pod) TrafficPath() string {
	return fmt.Sprintf("%s/traffics/victim-%d/pod-%d-%s.pcap", config.Env.Path, p.VictimID, p.ID, p.Name)
}

func (p Pod) RemoteAddr() []string {
	data := make([]string, 0)
	for _, port := range p.ExposedPorts {
		data = append(data, fmt.Sprintf("%s:%d", p.ExposedIP, port))
	}
	return data
}

type IP struct {
	Name    string
	Subnet  string
	PodName string
	IP      string
}

type IPs []IP

func (i IPs) Value() (driver.Value, error) {
	i = slices.DeleteFunc(i, func(i IP) bool {
		if net.ParseIP(i.IP) == nil {
			return true
		}
		return false
	})
	return json.Marshal(i)
}

func (i *IPs) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan IPs value")
	}
	return json.Unmarshal(b, i)
}
