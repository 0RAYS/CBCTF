package service

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/compose-spec/compose-go/v2/types"
)

func Template2Yaml(template model.ChallengeTemplate, challengeFlags []model.ChallengeFlag) string {
	cfg := types.Project{
		Services: make(types.Services),
		Networks: make(types.Networks),
	}

	networks := make(map[string]model.NetworkDefinition)
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
						Published: expose.Published,
						Mode:      "ingress",
						Target:    uint32(expose.Port),
					})
				}
			}
			if len(container.Environment) > 0 {
				service.Environment = make(map[string]*string)
				for key, value := range container.Environment {
					service.Environment[key] = new(value)
				}
			}
			xVolumes := make(model.XVolumes, 0)
			for _, flag := range challengeFlags {
				if flag.Binding.PodKey != pod.Key || flag.Binding.ContainerKey != container.Key {
					continue
				}
				switch flag.Binding.Type {
				case model.EnvFlagBindingType:
					if service.Environment == nil {
						service.Environment = make(map[string]*string)
					}
					service.Environment[flag.Binding.Target] = new(flag.Value)
				case model.FileFlagBindingType:
					xVolumes = append(xVolumes, model.XVolume{
						Path:    flag.Binding.Target,
						Content: flag.Value,
					})
				}
			}
			if len(xVolumes) > 0 {
				if service.Extensions == nil {
					service.Extensions = make(types.Extensions)
				}
				service.Extensions[model.XVolumesExtension] = xVolumes
			}
			if len(pod.Networks) > 0 {
				service.Networks = make(map[string]*types.ServiceNetworkConfig)
				for _, network := range pod.Networks {
					service.Networks[network.Attachment.Name] = &types.ServiceNetworkConfig{
						Ipv4Address: network.Attachment.IP,
					}
					networks[network.Definition.Name] = network.Definition
				}
			}
			cfg.Services[service.Name] = service
		}
	}
	for name, network := range networks {
		cfg.Networks[name] = types.NetworkConfig{
			External: types.External(network.External),
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
