package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/netip"
	"slices"
)

type NetworkDefinition struct {
	Name    string `json:"name"`
	CIDR    string `json:"cidr"`
	Gateway string `json:"gateway"`
}

type NetworkAttachment struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	MAC  string `json:"mac"`
}

type Network struct {
	Definition NetworkDefinition `json:"definition"`
	Attachment NetworkAttachment `json:"attachment"`
}

type Networks []Network

func (n Networks) Value() (driver.Value, error) {
	return json.Marshal(slices.DeleteFunc(n, func(n Network) bool {
		if n.Definition.CIDR == "" || n.Definition.Gateway == "" || n.Attachment.IP == "" {
			return true
		}
		cidr, err := netip.ParsePrefix(n.Definition.CIDR)
		if err != nil {
			return true
		}
		gateway, err := netip.ParseAddr(n.Definition.Gateway)
		if err != nil || !cidr.Contains(gateway) {
			return true
		}
		ip, err := netip.ParseAddr(n.Attachment.IP)
		if err != nil {
			return true
		}
		return !cidr.Contains(ip)
	}))
}

func (n *Networks) Scan(value any) error {
	if err := scanJSON(value, n); err != nil {
		return fmt.Errorf("failed to scan Networks value")
	}
	return nil
}

type Expose struct {
	Published string `json:"publish"`
	Protocol  string `json:"protocol"`
	Port      int32  `json:"port"`
}

type Exposes []Expose

func (e Exposes) Value() (driver.Value, error) {
	return json.Marshal(slices.DeleteFunc(e, func(n Expose) bool {
		return n.Port < 0 || n.Port > 65535
	}))
}

func (e *Exposes) Scan(value interface{}) error {
	if err := scanJSON(value, e); err != nil {
		return fmt.Errorf("failed to scan Exposes value")
	}
	return nil
}
