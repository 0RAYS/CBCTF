package resp

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"strconv"

	"github.com/compose-spec/compose-go/types"
	"github.com/gin-gonic/gin"
)

func Template2Yaml(template model.ChallengeTemplate, challengeFlags []model.ChallengeFlag) string {
	cfg := types.Project{
		Services: make(types.Services, 0),
		Networks: make(types.Networks),
		Volumes:  make(types.Volumes),
	}

	networks := make(map[string]model.Network)
	for _, pod := range template.Pods {
		for _, container := range pod.Containers {
			service := types.ServiceConfig{
				Name:       container.Name,
				Image:      container.Image,
				CPUS:       container.CPU,
				MemLimit:   types.UnitBytes(container.Memory),
				WorkingDir: container.WorkingDir,
				Command:    types.ShellCommand(container.Command),
			}
			if len(container.Exposes) > 0 {
				service.Ports = make([]types.ServicePortConfig, 0, len(container.Exposes))
				for _, expose := range container.Exposes {
					service.Ports = append(service.Ports, types.ServicePortConfig{
						Protocol:  expose.Protocol,
						Published: strconv.Itoa(int(expose.Port)),
						Mode:      "ingress",
						Target:    uint32(expose.Port),
					})
				}
			}
			if len(container.Environment) > 0 {
				service.Environment = make(map[string]*string)
				for key, value := range container.Environment {
					v := value
					service.Environment[key] = &v
				}
			}
			for _, flag := range challengeFlags {
				if flag.Binding.PodKey != pod.Key || flag.Binding.ContainerKey != container.Key {
					continue
				}
				switch flag.Binding.Type {
				case model.EnvFlagBindingType:
					if service.Environment == nil {
						service.Environment = make(map[string]*string)
					}
					v := flag.Value
					service.Environment[flag.Binding.Target] = &v
				case model.FileFlagBindingType:
					service.Volumes = append(service.Volumes, types.ServiceVolumeConfig{
						Type:   "volume",
						Source: flag.Name,
						Target: flag.Binding.Target,
					})
					cfg.Volumes[flag.Name] = types.VolumeConfig{
						Labels: map[string]string{
							model.VolumeFlagLabelKey: flag.Value,
						},
					}
				}
			}
			if len(pod.Networks) > 0 {
				service.Networks = make(map[string]*types.ServiceNetworkConfig)
				for _, network := range pod.Networks {
					service.Networks[network.Name] = &types.ServiceNetworkConfig{
						Ipv4Address: network.IP,
					}
					networks[network.Name] = network
				}
			}
			cfg.Services = append(cfg.Services, service)
		}
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
		log.Logger.Warningf("Failed to convert template to YAML: %s", err)
		return ""
	}
	return string(res)
}

func GetChallengeResp(challenge model.Challenge) gin.H {
	flags := make([]gin.H, 0)
	if challenge.Type != model.PodsChallengeType {
		for _, flag := range challenge.ChallengeFlags {
			flags = append(flags, gin.H{"id": flag.ID, "value": flag.Value})
		}
	}
	dockerCompose := ""
	if challenge.Type == model.PodsChallengeType {
		dockerCompose = Template2Yaml(challenge.Template, challenge.ChallengeFlags)
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
