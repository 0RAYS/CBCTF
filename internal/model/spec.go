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

// DeepCopy 返回 ChallengeTemplate 的深拷贝。
// 题目加入比赛后需要持有一份独立副本，模板内含多层 slice / map，
// 浅拷贝会让题库与比赛共享底层数组，改动一侧会污染另一侧。
func (t ChallengeTemplate) DeepCopy() ChallengeTemplate {
	if t.Pods == nil {
		return ChallengeTemplate{}
	}
	pods := make([]ChallengePodTemplate, 0, len(t.Pods))
	for _, pod := range t.Pods {
		pods = append(pods, ChallengePodTemplate{
			Key:          pod.Key,
			Name:         pod.Name,
			ServicePorts: pod.ServicePorts.DeepCopy(),
			Networks:     pod.Networks.DeepCopy(),
			Containers:   copyChallengeContainers(pod.Containers),
		})
	}
	return ChallengeTemplate{Pods: pods}
}

func copyChallengeContainers(containers []ChallengeContainerTemplate) []ChallengeContainerTemplate {
	if containers == nil {
		return nil
	}
	result := make([]ChallengeContainerTemplate, 0, len(containers))
	for _, container := range containers {
		result = append(result, ChallengeContainerTemplate{
			Key:          container.Key,
			Name:         container.Name,
			Image:        container.Image,
			CPU:          container.CPU,
			Memory:       container.Memory,
			WorkingDir:   container.WorkingDir,
			Command:      append([]string(nil), container.Command...),
			Environment:  copyStringMap(container.Environment),
			KubeVirt:     container.KubeVirt,
			Bootloader:   container.Bootloader,
			SecureBoot:   container.SecureBoot,
			UserData:     container.UserData.DeepCopy(),
			VolumeMounts: append([]XVolume(nil), container.VolumeMounts...),
			Exposes:      container.Exposes.DeepCopy(),
		})
	}
	return result
}

func (c CloudInitConfig) DeepCopy() CloudInitConfig {
	if c.Users == nil && c.Groups == nil && c.WriteFiles == nil {
		return CloudInitConfig{}
	}
	result := CloudInitConfig{WriteFiles: append([]CloudInitWriteFile(nil), c.WriteFiles...)}
	for _, user := range c.Users {
		user.SSHAuthorizedKeys = append([]string(nil), user.SSHAuthorizedKeys...)
		result.Users = append(result.Users, user)
	}
	for _, group := range c.Groups {
		group.Members = append([]string(nil), group.Members...)
		result.Groups = append(result.Groups, group)
	}
	return result
}

func copyStringMap(source map[string]string) map[string]string {
	if source == nil {
		return nil
	}
	result := make(map[string]string, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}
