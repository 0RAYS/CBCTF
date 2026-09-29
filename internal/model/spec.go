package model

import (
	"strings"
	"time"

	"go.yaml.in/yaml/v4"
)

type FlagBindingType string

const EnvFlagPrefix = "FLAG"

const (
	XVolumesExtension   = "x-volumes"
	XKubeVirtExtension  = "x-kubevirt"
	XBootExtension      = "x-boot"
	XCloudInitExtension = "x-cloudinit"
)

type XVolume struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content"`
}

type XBoot struct {
	Bootloader string `json:"bootloader" yaml:"bootloader"`
	SecureBoot bool   `json:"secure_boot,omitempty" yaml:"secure_boot,omitempty"`
}

type XCloudInit struct {
	Users             []CloudInitUser      `yaml:"users,omitempty" json:"users"`
	Groups            []CloudInitGroup     `yaml:"groups,omitempty" json:"groups"`
	WriteFiles        []CloudInitWriteFile `yaml:"write_files,omitempty" json:"write_files"`
	SSHAuthorizedKeys []string             `yaml:"ssh_authorized_keys,omitempty" json:"ssh_authorized_keys"`
}

type CloudInitConfig struct {
	Users             []CloudInitUser      `yaml:"users,omitempty" json:"users"`
	Groups            []CloudInitGroup     `yaml:"groups,omitempty" json:"groups"`
	WriteFiles        []CloudInitWriteFile `yaml:"write_files,omitempty" json:"write_files"`
	SSHAuthorizedKeys []string             `yaml:"ssh_authorized_keys,omitempty" json:"ssh_authorized_keys"`
}

type CloudInitUser struct {
	Name              string   `yaml:"name" json:"name"`
	Gecos             string   `yaml:"gecos" json:"gecos"`
	Groups            []string `yaml:"groups,omitempty" json:"groups"`
	Sudo              []string `yaml:"sudo,omitempty" json:"sudo"`
	Shell             string   `yaml:"shell" json:"shell"`
	HomeDir           string   `yaml:"homedir" json:"homedir"`
	LockPasswd        bool     `yaml:"lock_passwd" json:"lock_passwd"`
	Passwd            string   `yaml:"passwd" json:"passwd"`
	PlainTextPasswd   string   `yaml:"plain_text_passwd" json:"plain_text_passwd"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty" json:"ssh_authorized_keys"`
	NoCreateHome      bool     `yaml:"no_create_home" json:"no_create_home"`
	System            bool     `yaml:"system" json:"system"`
}

type CloudInitGroup struct {
	Name    string   `yaml:"name" json:"name"`
	Members []string `yaml:"members,omitempty" json:"members"`
}

func (x XCloudInit) CloudConfig() CloudInitConfig {
	return CloudInitConfig{
		Users:             x.Users,
		Groups:            x.Groups,
		WriteFiles:        x.WriteFiles,
		SSHAuthorizedKeys: x.SSHAuthorizedKeys,
	}
}

type CloudInitWriteFile struct {
	Path        string `yaml:"path" json:"path"`
	Content     string `yaml:"content" json:"content"`
	Owner       string `yaml:"owner" json:"owner"`
	Permissions string `yaml:"permissions" json:"permissions"`
	Encoding    string `yaml:"encoding" json:"encoding"`
	Append      bool   `yaml:"append" json:"append"`
	Defer       bool   `yaml:"defer" json:"defer"`
}

func (c CloudInitConfig) Empty() bool {
	return len(c.Users) == 0 && len(c.Groups) == 0 && len(c.WriteFiles) == 0 && len(c.SSHAuthorizedKeys) == 0
}

func (c CloudInitConfig) String() (string, error) {
	if c.Empty() {
		return "", nil
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	return "#cloud-config\n" + strings.TrimSpace(string(data)), nil
}

const (
	EnvFlagBindingType           FlagBindingType = "env"
	FileFlagBindingType          FlagBindingType = "file"
	CloudInitFileFlagBindingType FlagBindingType = "cloudinit_file"
)

type FlagBinding struct {
	PodKey       string          `json:"pod_key"`
	ContainerKey string          `json:"container_key"`
	Type         FlagBindingType `json:"type"`
	Target       string          `json:"target"`
}

type ChallengeContainerTemplate struct {
	Key          string            `json:"key"`
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	CPU          float32           `json:"cpu"`
	Memory       int64             `json:"memory"`
	WorkingDir   string            `json:"working_dir"`
	Command      []string          `json:"command"`
	Environment  map[string]string `json:"environment"`
	KubeVirt     bool              `json:"kubevirt"`
	Bootloader   string            `json:"bootloader"`
	SecureBoot   bool              `json:"secure_boot"`
	UserData     CloudInitConfig   `json:"user_data"`
	VolumeMounts []XVolume         `json:"volume_mounts"`
	Exposes      Exposes           `json:"exposes"`
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

type VictimContainerSpec struct {
	Key          string            `json:"key"`
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	Resources    ResourceSpec      `json:"resources"`
	WorkingDir   string            `json:"working_dir"`
	Command      []string          `json:"command"`
	Environment  map[string]string `json:"environment"`
	KubeVirt     bool              `json:"kubevirt"`
	Bootloader   string            `json:"bootloader"`
	SecureBoot   bool              `json:"secure_boot"`
	UserData     CloudInitConfig   `json:"user_data"`
	VolumeMounts []XVolume         `json:"volume_mounts"`
	Exposes      Exposes           `json:"exposes"`
}

type ResourceSpec struct {
	CPUMillis   int64 `json:"cpu_millis"`
	MemoryBytes int64 `json:"memory_bytes"`
}

type PodSpec struct {
	Key          string                `json:"key"`
	ServicePorts Exposes               `json:"service_ports"`
	Networks     Networks              `json:"networks"`
	Containers   []VictimContainerSpec `json:"containers"`
}

type VictimSpec struct {
	Pods            []PodSpec       `json:"pods"`
	NetworkPlan     VPC             `json:"network_plan"`
	NetworkPolicies NetworkPolicies `json:"network_policies"`
	FrpEnabled      bool            `json:"frp_enabled"`
}

type VictimResources struct {
	Submitted     bool               `json:"submitted"`
	ReadyDeadline time.Time          `json:"ready_deadline"`
	UIDs          map[string]string  `json:"uids"`
	NodePorts     []NodePortEndpoint `json:"node_ports"`
	NetworkPlan   VPC                `json:"network_plan"`
	PodNames      []string           `json:"pod_names"`
	FrpcPodNames  []string           `json:"frpc_pod_names"`
}

type NodePortEndpoint struct {
	PodName  string   `json:"pod_name"`
	Endpoint Endpoint `json:"endpoint"`
}
