package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/netip"
	"slices"
)

type NetworkDefinition struct {
	Name     string `json:"name"`
	CIDR     string `json:"cidr"`
	Gateway  string `json:"gateway"`
	External bool   `json:"external"`
}

type NetworkAttachment struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

// Network is the persisted legacy shape used in challenge and pod specs.
// Keep the JSON flat for existing data, but expose definition/attachment views
// so callers do not have to mix network-level and pod-level fields manually.
type Network struct {
	Name     string `json:"name"`
	CIDR     string `json:"cidr"`
	Gateway  string `json:"gateway"`
	IP       string `json:"ip"`
	External bool   `json:"external"`
}

func (n Network) Definition() NetworkDefinition {
	return NetworkDefinition{
		Name:     n.Name,
		CIDR:     n.CIDR,
		Gateway:  n.Gateway,
		External: n.External,
	}
}

func (n Network) Attachment() NetworkAttachment {
	return NetworkAttachment{
		Name: n.Name,
		IP:   n.IP,
	}
}

func NewNetwork(definition NetworkDefinition, attachment NetworkAttachment) Network {
	return Network{
		Name:     definition.Name,
		CIDR:     definition.CIDR,
		Gateway:  definition.Gateway,
		IP:       attachment.IP,
		External: definition.External,
	}
}

type Networks []Network

func (n Networks) Value() (driver.Value, error) {
	n = slices.DeleteFunc(n, func(n Network) bool {
		if n.CIDR == "" || n.Gateway == "" || n.IP == "" {
			return true
		}
		cidr, err := netip.ParsePrefix(n.CIDR)
		if err != nil {
			return true
		}
		gateway, err := netip.ParseAddr(n.Gateway)
		if err != nil || !cidr.Contains(gateway) {
			return true
		}
		ip, err := netip.ParseAddr(n.IP)
		if err != nil {
			return true
		}
		return !cidr.Contains(ip)
	})
	return json.Marshal(n)
}

func (n *Networks) Scan(value any) error {
	if err := scanJSON(value, n); err != nil {
		return fmt.Errorf("failed to scan Networks value")
	}
	return nil
}

type Expose struct {
	Published string `json:"publish"`
	Port      int32  `json:"port"`
	Protocol  string `json:"protocol"`
}

type Exposes []Expose

func (e Exposes) Value() (driver.Value, error) {
	e = slices.DeleteFunc(e, func(n Expose) bool {
		return n.Port < 0 || n.Port > 65535
	})
	return json.Marshal(e)
}

func (e *Exposes) Scan(value interface{}) error {
	if err := scanJSON(value, e); err != nil {
		return fmt.Errorf("failed to scan Exposes value")
	}
	return nil
}
