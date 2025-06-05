package model

import (
	"CBCTF/internel/config"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type AvatarURL string

func (a AvatarURL) Value() (driver.Value, error) {
	if a == "" {
		return nil, nil
	}
	return strings.TrimPrefix(string(a), strings.Trim(config.Env.Backend, "/")), nil
}

func (a *AvatarURL) Scan(value any) error {
	if value == nil || value.(string) == "" {
		*a = ""
		return nil
	}
	path, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan AvatarURL: %v", value)
	}
	*a = AvatarURL(strings.Trim(config.Env.Backend, "/") + path)
	return nil
}

type StringList []string

func (s StringList) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringList) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringList value")
	}
	return json.Unmarshal(bytes, s)
}

type UintList []uint

func (u UintList) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *UintList) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan UintList value")
	}
	return json.Unmarshal(bytes, u)
}

type StringMap map[string]string

func (s StringMap) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringMap) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringMap value")
	}
	return json.Unmarshal(bytes, s)
}

type Target struct {
	Hostname string   `json:"hostname"`
	CIDR     string   `json:"cidr"`
	Except   []string `json:"except"`
}

func (t Target) isValidIPBlock() bool {
	if t.Hostname != "" {
		return true
	}
	if t.Hostname == "" && t.CIDR == "" {
		return false
	}
	_, ipNet, err := net.ParseCIDR(t.CIDR)
	if err != nil {
		return false
	}
	for _, ex := range t.Except {
		_, exNet, err := net.ParseCIDR(ex)
		if err != nil {
			return false
		}
		if !ipNet.Contains(exNet.IP) {
			return false
		}
	}
	return true
}

type NetworkPolicy struct {
	From []Target `json:"from"`
	To   []Target `json:"to"`
}

var DefaultNetworkPolicy = NetworkPolicy{
	To: []Target{
		{
			CIDR: "0.0.0.0/0",
			Except: []string{
				"10.0.0.0/8",
				"172.16.0.0/12",
				"192.168.0.0/16",
				"100.64.0.0/10",
			},
		},
	},
}

type NetworkPolicies []NetworkPolicy

func (n NetworkPolicies) Value() (driver.Value, error) {
	for _, p := range n {
		for i, ipBlock := range p.From {
			if !ipBlock.isValidIPBlock() {
				p.From = append(p.From[:i], p.From[i+1:]...)
			}
		}
		for i, ipBlock := range p.To {
			if !ipBlock.isValidIPBlock() {
				p.To = append(p.To[:i], p.To[i+1:]...)
			}
		}
	}
	return json.Marshal(n)
}

func (n *NetworkPolicies) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan NetworkPolicy value")
	}
	return json.Unmarshal(bytes, n)
}
