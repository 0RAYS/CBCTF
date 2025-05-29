package model

type Container struct {
	PodID    uint    `json:"pod_id"`
	Pod      Pod     `json:"-"`
	Name     string  `json:"name"`
	Image    string  `json:"image"`
	Hostname string  `json:"hostname"`
	Flags    Strings `gorm:"type:json" json:"flags"`
	PodPorts Ports   `gorm:"type:json" json:"pod_ports"`
	BaseModel
}
