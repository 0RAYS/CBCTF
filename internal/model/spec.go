package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type FlagBindingType string

const (
	EnvFlagPrefix      = "FLAG"
	VolumeFlagPrefix   = "FLAG"
	VolumeFlagLabelKey = "value"
)

const (
	EnvFlagBindingType  FlagBindingType = "env"
	FileFlagBindingType FlagBindingType = "file"
)

type FlagBinding struct {
	PodKey       string          `json:"pod_key"`
	ContainerKey string          `json:"container_key"`
	Type         FlagBindingType `json:"type"`
	Target       string          `json:"target"`
}

func (f FlagBinding) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *FlagBinding) Scan(value any) error {
	if err := scanJSON(value, f); err != nil {
		return fmt.Errorf("failed to scan FlagBinding value")
	}
	return nil
}

type FlagBindings []FlagBinding

func (f FlagBindings) Value() (driver.Value, error) {
	if len(f) == 0 {
		return nil, nil
	}
	return json.Marshal(f)
}

func (f *FlagBindings) Scan(value any) error {
	if err := scanJSON(value, f); err != nil {
		return fmt.Errorf("failed to scan FlagBindings value")
	}
	return nil
}

type ChallengeContainerTemplate struct {
	Key         string     `json:"key"`
	Name        string     `json:"name"`
	Image       string     `json:"image"`
	CPU         float32    `json:"cpu"`
	Memory      int64      `json:"memory"`
	WorkingDir  string     `json:"working_dir"`
	Command     StringList `json:"command"`
	Environment StringMap  `json:"environment"`
	Exposes     Exposes    `json:"exposes"`
}

type ChallengePodTemplate struct {
	Key          string                       `json:"key"`
	Name         string                       `json:"name"`
	ServicePorts Exposes                      `json:"service_ports"`
	Networks     Networks                     `json:"networks"`
	Containers   []ChallengeContainerTemplate `json:"containers"`
}

type ChallengeTemplate struct {
	Pods []ChallengePodTemplate `json:"pods"`
}

func (c ChallengeTemplate) Value() (driver.Value, error) {
	if len(c.Pods) == 0 {
		return nil, nil
	}
	return json.Marshal(c)
}

func (c *ChallengeTemplate) Scan(value any) error {
	if err := scanJSON(value, c); err != nil {
		return fmt.Errorf("failed to scan ChallengeTemplate value")
	}
	return nil
}

type VictimContainerSpec struct {
	Key         string     `json:"key"`
	Name        string     `json:"name"`
	Image       string     `json:"image"`
	CPU         float32    `json:"cpu"`
	Memory      int64      `json:"memory"`
	WorkingDir  string     `json:"working_dir"`
	Command     StringList `json:"command"`
	Environment StringMap  `json:"environment"`
	Files       StringMap  `json:"files"`
	Exposes     Exposes    `json:"exposes"`
}

type PodSpec struct {
	Key          string                `json:"key"`
	ServicePorts Exposes               `json:"service_ports"`
	Networks     Networks              `json:"networks"`
	Containers   []VictimContainerSpec `json:"containers"`
}

func (p PodSpec) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PodSpec) Scan(value any) error {
	if err := scanJSON(value, p); err != nil {
		return fmt.Errorf("failed to scan PodSpec value")
	}
	return nil
}

type VictimSpec struct {
	Pods            []PodSpec       `json:"pods"`
	NetworkPlan     VPC             `json:"network_plan"`
	NetworkPolicies NetworkPolicies `json:"network_policies"`
	FrpEnabled      bool            `json:"frp_enabled"`
}

func (v VictimSpec) Value() (driver.Value, error) {
	if len(v.Pods) == 0 && v.NetworkPlan.Name == "" && len(v.NetworkPolicies) == 0 && !v.FrpEnabled {
		return nil, nil
	}
	return json.Marshal(v)
}

func (v *VictimSpec) Scan(value any) error {
	if err := scanJSON(value, v); err != nil {
		return fmt.Errorf("failed to scan VictimSpec value")
	}
	return nil
}

type VictimResources struct {
	NetworkPlan  VPC        `json:"network_plan"`
	PodNames     StringList `json:"pod_names"`
	FrpcPodNames StringList `json:"frpc_pod_names"`
}

func (v VictimResources) Value() (driver.Value, error) {
	if v.NetworkPlan.Name == "" && len(v.PodNames) == 0 && len(v.FrpcPodNames) == 0 {
		return nil, nil
	}
	return json.Marshal(v)
}

func (v *VictimResources) Scan(value any) error {
	if err := scanJSON(value, v); err != nil {
		return fmt.Errorf("failed to scan VictimResources value")
	}
	return nil
}
