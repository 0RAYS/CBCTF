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
		if flag.DockerID == nil {
			continue
		}
		switch flag.InjectType {
		case model.VolumeInjectType:
			volumeFlags[*flag.DockerID] = make(map[string]string)
			name := fmt.Sprintf("%s_%d", model.VolumeFlagPrefix, i)
			volumeFlags[*flag.DockerID][name] = flag.Path
			volumeStr += fmt.Sprintf("\t%s:\n", name)
			volumeStr += fmt.Sprintf("\t\tlabels:\n")
			volumeStr += fmt.Sprintf("\t\t\t- %s=%s\n", model.VolumeFlagLabelKey, flag.Value)
		case model.EnvInjectType:
			envFlags[*flag.DockerID] = make(map[string]string)
			name := fmt.Sprintf("%s_%d", model.EnvFlagPrefix, i)
			envFlags[*flag.DockerID][name] = flag.Value
		default:
			continue
		}
	}
	volumeStr = strings.Trim(volumeStr, "\n")

	var serviceStr string
	for _, docker := range dockers {
		serviceStr += fmt.Sprintf("\t%s:\n", docker.Name)
		serviceStr += fmt.Sprintf("\t\timage: %s\n", docker.Image)
		if docker.PullPolicy != nil && *docker.PullPolicy != "" {
			serviceStr += fmt.Sprintf("\t\tpull_policy: %s\n", *docker.PullPolicy)
		}
		if docker.WorkingDir != nil && *docker.WorkingDir != "" {
			serviceStr += fmt.Sprintf("\t\tworking_dir: %s\n", *docker.WorkingDir)
		}
		if docker.Command != nil && len(docker.Command) > 0 {
			commandStr := "["
			for _, cmd := range docker.Command {
				commandStr += fmt.Sprintf("\"%s\", ", cmd)
			}
			commandStr = commandStr[:len(commandStr)-2] + "]"
			serviceStr += fmt.Sprintf("\t\tcommand: %s\n", commandStr)
		}
		if docker.Expose != nil && len(docker.Expose) > 0 {
			serviceStr += "\t\texpose:\n"
			for _, port := range docker.Expose {
				serviceStr += fmt.Sprintf("\t\t\t- \"%s\"\n", port)
			}
		}
		if docker.Environment != nil || len(envFlags[docker.ID]) > 0 {
			serviceStr += "\t\tenvironment:\n"
			if docker.Environment != nil && len(docker.Environment) > 0 {
				for key, value := range docker.Environment {
					serviceStr += fmt.Sprintf("\t\t\t- %s=%s\n", key, value)
				}
			}
			if flags, ok := envFlags[docker.ID]; ok {
				for key, value := range flags {
					serviceStr += fmt.Sprintf("\t\t\t- %s=%s\n", key, value)
				}
			}
		}
		if flags, ok := volumeFlags[docker.ID]; ok {
			serviceStr += "\t\tvolumes:\n"
			for key, path := range flags {
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
