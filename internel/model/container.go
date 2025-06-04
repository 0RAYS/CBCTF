package model

type Container struct {
	PodID    uint       `json:"pod_id"`
	Pod      Pod        `json:"-"`
	Name     string     `json:"name"`
	Image    string     `json:"image"`
	Hostname string     `json:"hostname"`
	Flags    StringList `gorm:"type:json" json:"flags"`
	PodPorts PortList   `gorm:"type:json" json:"pod_ports"`
	BaseModel
}
