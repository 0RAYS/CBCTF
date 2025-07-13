package resp

import (
	"CBCTF/internel/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func Dockers2Yaml(dockers []model.Docker, challengeFlags []model.ChallengeFlag) string {
	baseYaml := `
services:
%s
`
	var volumeStr string
	volumeFlags := make(map[uint]map[string]string)
	envFlags := make(map[uint]map[string]string)
	for i, flag := range challengeFlags {
		if flag.DockerID == nil {
			continue
		}
		switch flag.InjectType {
		case model.VolumeInjectType:
			if volumeFlags[*flag.DockerID] == nil {
				volumeFlags[*flag.DockerID] = make(map[string]string)
			}
			volumeFlags[*flag.DockerID] = make(map[string]string)
			name := fmt.Sprintf("%s_%d", model.VolumeFlagPrefix, i)
			volumeFlags[*flag.DockerID][name] = flag.Path
			volumeStr += fmt.Sprintf("  %s:\n", name)
			volumeStr += fmt.Sprintf("    labels:\n")
			volumeStr += fmt.Sprintf("      - %s=%s\n", model.VolumeFlagLabelKey, flag.Value)
		case model.EnvInjectType:
			if envFlags[*flag.DockerID] == nil {
				envFlags[*flag.DockerID] = make(map[string]string)
			}
			name := fmt.Sprintf("%s_%d", model.EnvFlagPrefix, i)
			envFlags[*flag.DockerID][name] = flag.Value
		default:
			continue
		}
	}
	volumeStr = strings.Trim(volumeStr, "\n")

	var (
		serviceStr string
		networks   = make(map[string]model.Network)
	)
	for _, docker := range dockers {
		serviceStr += fmt.Sprintf("  %s:\n", docker.Name)
		serviceStr += fmt.Sprintf("    image: %s\n", docker.Image)
		if docker.CPU > 0 {
			serviceStr += fmt.Sprintf("    cpus: %f\n", docker.CPU)
		}
		if docker.Memory > 0 {
			serviceStr += fmt.Sprintf("    mem_limit: %d\n", docker.Memory)
		}
		if docker.WorkingDir != "" {
			serviceStr += fmt.Sprintf("    working_dir: %s\n", docker.WorkingDir)
		}
		if docker.Command != nil && len(docker.Command) > 0 {
			commandStr := "["
			for _, cmd := range docker.Command {
				commandStr += fmt.Sprintf("\"%s\", ", cmd)
			}
			commandStr = commandStr[:len(commandStr)-2] + "]"
			serviceStr += fmt.Sprintf("    command: %s\n", commandStr)
		}
		if docker.Expose != nil && len(docker.Expose) > 0 {
			serviceStr += "    expose:\n"
			for _, port := range docker.Expose {
				serviceStr += fmt.Sprintf("      - \"%s\"\n", port)
			}
		}
		if docker.Environment != nil || len(envFlags[docker.ID]) > 0 {
			serviceStr += "    environment:\n"
			if docker.Environment != nil && len(docker.Environment) > 0 {
				for key, value := range docker.Environment {
					serviceStr += fmt.Sprintf("      - %s=%s\n", key, value)
				}
			}
			if flags, ok := envFlags[docker.ID]; ok {
				for key, value := range flags {
					serviceStr += fmt.Sprintf("      - %s=%s\n", key, value)
				}
			}
		}
		if flags, ok := volumeFlags[docker.ID]; ok {
			serviceStr += "    volumes:\n"
			for key, path := range flags {
				serviceStr += fmt.Sprintf("      - %s:%s\n", key, path)
			}
		}
		if docker.Networks != nil && len(docker.Networks) > 0 {
			serviceStr += "    networks:\n"
			for _, network := range docker.Networks {
				networkName := strings.ReplaceAll(network.Subnet, ".", "_")
				networkName = fmt.Sprintf("network_%s", strings.ReplaceAll(networkName, "/", "_"))
				serviceStr += fmt.Sprintf("      %s:\n        ipv4_address: %s\n", networkName, network.IP)
				networks[networkName] = network
			}
		}
		serviceStr += "\n"
	}
	serviceStr = strings.Trim(serviceStr, "\n")

	var networkStr string
	for name, network := range networks {
		networkStr += fmt.Sprintf("  %s:\n", name)
		networkStr += fmt.Sprintf("    external: %v\n", network.External)
		networkStr += fmt.Sprintf("    ipam:\n      config:\n        - subnet: %s\n          gateway: %s\n", network.Subnet, network.Gateway)
		networkStr += "\n"
	}
	networkStr = strings.Trim(networkStr, "\n")
	baseYaml = strings.Trim(fmt.Sprintf(baseYaml, serviceStr), "\n")
	if volumeStr != "" {
		baseYaml += fmt.Sprintf("\n\nvolumes:\n%s", volumeStr)
	}
	if networkStr != "" {
		baseYaml += fmt.Sprintf("\n\nnetworks:\n%s", networkStr)
	}
	return strings.Trim(baseYaml, "\n")
}

// GetChallengeResp 需要预加载 ChallengeFlags, Dockers
func GetChallengeResp(challenge model.Challenge) gin.H {
	flags := make([]gin.H, 0)
	if challenge.Type != model.PodsChallengeType {
		for _, flag := range challenge.ChallengeFlags {
			flags = append(flags, gin.H{"id": flag.ID, "value": flag.Value})
		}
	}
	return gin.H{
		"id":               challenge.RandID,
		"name":             challenge.Name,
		"desc":             challenge.Desc,
		"category":         challenge.Category,
		"type":             challenge.Type,
		"generator_image":  challenge.GeneratorImage,
		"flags":            flags,
		"docker_compose":   Dockers2Yaml(challenge.Dockers, challenge.ChallengeFlags),
		"network_policies": challenge.NetworkPolicies,
	}
}
