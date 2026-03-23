package resp

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"strconv"

	"github.com/compose-spec/compose-go/types"
	"github.com/gin-gonic/gin"
)

func Dockers2Yaml(dockers []model.Docker, challengeFlags []model.ChallengeFlag) string {
	cfg := types.Project{
		Services: make(types.Services, 0),
		Networks: make(types.Networks),
		Volumes:  make(types.Volumes),
	}

	volumeFlags := make(map[uint]map[string]string)
	envFlags := make(map[uint]map[string]string)
	for _, flag := range challengeFlags {
		if !flag.DockerID.Valid {
			continue
		}
		switch flag.InjectType {
		case model.VolumeFlagInjectType:
			if volumeFlags[flag.DockerID.V] == nil {
				volumeFlags[flag.DockerID.V] = make(map[string]string)
			}
			volumeFlags[flag.DockerID.V] = make(map[string]string)
			volumeFlags[flag.DockerID.V][flag.Name] = flag.Path
			cfg.Volumes[flag.Name] = types.VolumeConfig{
				Labels: map[string]string{
					model.VolumeFlagLabelKey: flag.Value,
				},
			}
		case model.EnvFlagInjectType:
			if envFlags[flag.DockerID.V] == nil {
				envFlags[flag.DockerID.V] = make(map[string]string)
			}
			envFlags[flag.DockerID.V][flag.Name] = flag.Value
		default:
			continue
		}
	}

	var networks = make(map[string]model.Network)
	for _, docker := range dockers {
		service := types.ServiceConfig{
			Name:       docker.Name,
			Image:      docker.Image,
			CPUS:       docker.CPU,
			MemLimit:   types.UnitBytes(docker.Memory),
			WorkingDir: docker.WorkingDir,
			Command:    types.ShellCommand(docker.Command),
		}
		if docker.Command != nil && len(docker.Command) > 0 {
			service.Command = types.ShellCommand(docker.Command)
		}
		if docker.Exposes != nil && len(docker.Exposes) > 0 {
			service.Ports = make([]types.ServicePortConfig, 0)
			for _, expose := range docker.Exposes {
				service.Ports = append(service.Ports, types.ServicePortConfig{
					Protocol:  expose.Protocol,
					Published: strconv.Itoa(int(expose.Port)),
					Mode:      "ingress",
					Target:    uint32(expose.Port),
				})
			}
		}
		if docker.Environment != nil || len(envFlags[docker.ID]) > 0 {
			service.Environment = make(map[string]*string)
			if docker.Environment != nil && len(docker.Environment) > 0 {
				for key, value := range docker.Environment {
					service.Environment[key] = &value
				}
			}
			if flags, ok := envFlags[docker.ID]; ok {
				for key, value := range flags {
					service.Environment[key] = &value
				}
			}
		}
		if flags, ok := volumeFlags[docker.ID]; ok {
			service.Volumes = make([]types.ServiceVolumeConfig, 0)
			for key, path := range flags {
				service.Volumes = append(service.Volumes, types.ServiceVolumeConfig{
					Type:   "volume",
					Source: key,
					Target: path,
				})
			}
		}
		if docker.Networks != nil && len(docker.Networks) > 0 {
			service.Networks = make(map[string]*types.ServiceNetworkConfig)
			for _, network := range docker.Networks {
				service.Networks[network.Name] = &types.ServiceNetworkConfig{
					Ipv4Address: network.IP,
				}
				networks[network.Name] = network
			}
		}
		cfg.Services = append(cfg.Services, service)
	}
	for name, network := range networks {
		cfg.Networks[name] = types.NetworkConfig{
			External: types.External{
				External: network.External,
			},
			Ipam: types.IPAMConfig{
				Config: []*types.IPAMPool{
					{
						Subnet:  network.CIDR,
						Gateway: network.Gateway,
					},
				},
			},
		}
	}
	res, err := cfg.MarshalYAML()
	if err != nil {
		log.Logger.Warningf("Failed to convert dockers to YAML: %s", err)
		return ""
	}
	return string(res)
}

// GetChallengeResp model.Challenge Preload model.Docker model.ChallengeFlag
func GetChallengeResp(challenge model.Challenge) gin.H {
	flags := make([]gin.H, 0)
	if challenge.Type != model.PodsChallengeType {
		for _, flag := range challenge.ChallengeFlags {
			flags = append(flags, gin.H{"id": flag.ID, "value": flag.Value})
		}
	}
	dockerCompose := ""
	if challenge.Type == model.PodsChallengeType {
		dockerCompose = Dockers2Yaml(challenge.Dockers, challenge.ChallengeFlags)
	}
	file, _ := db.InitFileRepo(db.DB).Get(db.GetOptions{
		Conditions: map[string]any{"model": model.ModelName(challenge), "model_id": challenge.ID, "type": model.ChallengeFileType},
	})
	return gin.H{
		"id":               challenge.RandID,
		"name":             challenge.Name,
		"description":      challenge.Description,
		"category":         challenge.Category,
		"type":             challenge.Type,
		"generator_image":  challenge.GeneratorImage,
		"flags":            flags,
		"docker_compose":   dockerCompose,
		"options":          challenge.Options,
		"network_policies": challenge.NetworkPolicies,
		"file":             file.Filename,
	}
}

func GetSimpleChallengeResp(challenge model.Challenge) gin.H {
	return gin.H{
		"id":               challenge.RandID,
		"name":             challenge.Name,
		"description":      challenge.Description,
		"category":         challenge.Category,
		"type":             challenge.Type,
		"generator_image":  challenge.GeneratorImage,
		"options":          challenge.Options,
		"network_policies": challenge.NetworkPolicies,
	}
}
