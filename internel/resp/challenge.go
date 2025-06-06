package resp

import (
	"CBCTF/internel/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func Docker2Yaml(dockers []model.Docker, challengeFlags []model.ChallengeFlag) string {
	baseYaml := `
services:
%s

volumes:
%s
`
	var volumeStr string
	volumeFlags := make(map[uint]map[string]string)
	envFlags := make(map[uint]map[string]string)
	for i, flag := range challengeFlags {
		switch flag.InjectType {
		case model.VolumeInjectType:
			volumeFlags[flag.DockerID] = make(map[string]string)
			name := fmt.Sprintf("%s_%d", model.VolumeFlagPrefix, i)
			volumeFlags[flag.DockerID][name] = flag.Path
			volumeStr += fmt.Sprintf("\t%s:\n", name)
			volumeStr += fmt.Sprintf("\t\tlabels:\n")
			volumeStr += fmt.Sprintf("\t\t\t- %s=%s\n", model.VolumeFlagLabelKey, flag.Value)
		case model.EnvInjectType:
			envFlags[flag.DockerID] = make(map[string]string)
			name := fmt.Sprintf("%s_%d", model.EnvFlagPrefix, i)
			envFlags[flag.DockerID][name] = flag.Value
		default:
			continue
		}
	}
	volumeStr = strings.Trim(volumeStr, "\n")

	var serviceStr string
	for _, docker := range dockers {
		var (
			commandStr    string
			entrypointStr string
		)
		serviceStr += fmt.Sprintf("\t%s:\n", docker.Name)
		serviceStr += fmt.Sprintf("\t\timage: %s\n", docker.Image)
		if docker.PullPolicy != nil {
			serviceStr += fmt.Sprintf("\t\tpull_policy: %s\n", *docker.PullPolicy)
		}
		if docker.Hostname != nil {
			serviceStr += fmt.Sprintf("\t\thostname: %s\n", *docker.Hostname)
		}
		if docker.WorkingDir != nil {
			serviceStr += fmt.Sprintf("\t\tworking_dir: %s\n", *docker.WorkingDir)
		}
		if docker.User != nil {
			serviceStr += fmt.Sprintf("\t\tuser: %s\n", *docker.User)
		}
		if docker.CPUCount != nil {
			serviceStr += fmt.Sprintf("\t\tcpu_count: %d\n", *docker.CPUCount)
		}
		if docker.CPUPercent != nil {
			serviceStr += fmt.Sprintf("\t\tcpu_percent: %.2f\n", *docker.CPUPercent)
		}
		if docker.CPUPeriod != nil {
			serviceStr += fmt.Sprintf("\t\tcpu_period: %d\n", *docker.CPUPeriod)
		}
		if docker.CPUQuota != nil {
			serviceStr += fmt.Sprintf("\t\tcpu_quota: %d\n", *docker.CPUQuota)
		}
		if docker.CPURTPeriod != nil {
			serviceStr += fmt.Sprintf("\t\tcpu_rt_period: %d\n", *docker.CPURTPeriod)
		}
		if docker.CPURTRuntime != nil {
			serviceStr += fmt.Sprintf("\t\tcpu_rt_runtime: %d\n", *docker.CPURTRuntime)
		}
		if docker.CPUS != nil {
			serviceStr += fmt.Sprintf("\t\tcpus: %.2f\n", *docker.CPUS)
		}
		if docker.CPUSet != nil {
			serviceStr += fmt.Sprintf("\t\tcpu_set: %s\n", *docker.CPUSet)
		}
		if docker.CPUShares != nil {
			serviceStr += fmt.Sprintf("\t\tcpu_shares: %d\n", *docker.CPUShares)
		}
		if docker.MemLimit != nil {
			serviceStr += fmt.Sprintf("\t\tmem_limit: %d\n", *docker.MemLimit)
		}
		if docker.MemReservation != nil {
			serviceStr += fmt.Sprintf("\t\tmem_reservation: %d\n", *docker.MemReservation)
		}
		if docker.MemSwapLimit != nil {
			serviceStr += fmt.Sprintf("\t\tmem_swap_limit: %d\n", *docker.MemSwapLimit)
		}
		if docker.MemSwappiness != nil {
			serviceStr += fmt.Sprintf("\t\tmem_swappiness: %d\n", *docker.MemSwappiness)
		}
		if docker.Command != nil {
			commandStr = "["
			for _, cmd := range *docker.Command {
				commandStr += fmt.Sprintf("\"%s\", ", cmd)
			}
			commandStr = commandStr[:len(commandStr)-2] + "]"
			serviceStr += fmt.Sprintf("\t\tcommand: %s\n", commandStr)
		}
		if docker.Entrypoint != nil {
			entrypointStr = "["
			for _, entry := range *docker.Entrypoint {
				entrypointStr += fmt.Sprintf("\"%s\", ", entry)
			}
			entrypointStr = entrypointStr[:len(entrypointStr)-2] + "]"
			serviceStr += fmt.Sprintf("\t\tentrypoint: %s\n", entrypointStr)
		}
		if docker.Expose != nil {
			serviceStr += "\t\texpose:\n"
			for _, port := range *docker.Expose {
				serviceStr += fmt.Sprintf("\t\t\t- \"%s\"\n", port)
			}
		}
		if docker.Environment != nil || len(envFlags[docker.ID]) > 0 {
			serviceStr += "\t\tenvironment:\n"
			if docker.Environment != nil {
				for key, value := range *docker.Environment {
					serviceStr += fmt.Sprintf("\t\t\t- %s=%s\n", key, value)
				}
			}
			if flags, ok := envFlags[docker.ID]; ok {
				for key, value := range flags {
					serviceStr += fmt.Sprintf("\t\t\t- %s=%s\n", key, value)
				}
			}
		}
		if volumeFlags, ok := volumeFlags[docker.ID]; ok {
			serviceStr += "\t\tvolumes:\n"
			for key, path := range volumeFlags {
				serviceStr += fmt.Sprintf("\t\t\t- %s:%s\n", key, path)
			}
		}
	}
	serviceStr = strings.Trim(serviceStr, "\n")

	return fmt.Sprintf(baseYaml, serviceStr, volumeStr)
}

// GetChallengeResp 需要预加载 DockerGroups, ChallengeFlags, DockerGroups.Dockers
func GetChallengeResp(challenge model.Challenge) gin.H {
	flags := make([]string, 0)
	if challenge.Type != model.PodsChallengeType {
		for _, flag := range challenge.ChallengeFlags {
			flags = append(flags, flag.Value)
		}
	}
	dockerGroups := make([]gin.H, 0)
	for _, group := range challenge.DockerGroups {
		dockerGroups = append(dockerGroups, gin.H{
			"yaml":             Docker2Yaml(group.Dockers, challenge.ChallengeFlags),
			"network_policies": group.NetworkPolicies,
		})
	}
	return gin.H{
		"id":              challenge.RandID,
		"name":            challenge.Name,
		"desc":            challenge.Desc,
		"category":        challenge.Category,
		"type":            challenge.Type,
		"generator_image": challenge.GeneratorImage,
		"flags":           flags,
		"docker_groups":   dockerGroups,
	}
}
