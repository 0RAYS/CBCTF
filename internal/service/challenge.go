package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"net"
	"net/netip"
	"regexp"
	"slices"
	"strings"

	"github.com/compose-spec/compose-go/v2/types"
	"gorm.io/gorm"
)

var (
	composePortNameRegex = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	composeEnvKeyRegex   = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
)

func invalidComposeYamlRetVal(err string) model.RetVal {
	return model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": err}}
}

func serviceKubeVirt(service types.ServiceConfig) bool {
	kubeVirt, _ := service.Extensions[model.XKubeVirtExtension].(bool)
	return kubeVirt
}

func validateChallengeCompose(config *types.Project) model.RetVal {
	if len(config.Services) == 0 {
		return invalidComposeYamlRetVal("At least one service is required")
	}

	hasVpcNetworks := false
	networkSubnets := make(map[string]struct{})
	networkGateways := make(map[string]struct{})
	for _, network := range config.Networks {
		name := strings.TrimSpace(strings.TrimPrefix(network.Name, config.Name+"_"))
		if name == "" || name == "default" {
			continue
		}
		hasVpcNetworks = true
		if len(network.Ipam.Config) == 0 {
			return invalidComposeYamlRetVal("Empty IPAM")
		}
		subnet := strings.TrimSpace(network.Ipam.Config[0].Subnet)
		gateway := strings.TrimSpace(network.Ipam.Config[0].Gateway)
		if subnet == "" {
			return invalidComposeYamlRetVal(fmt.Sprintf("network %s subnet is required", name))
		}
		if _, ok := networkSubnets[subnet]; ok {
			return invalidComposeYamlRetVal(fmt.Sprintf("network %s subnet must be unique", name))
		}
		networkSubnets[subnet] = struct{}{}
		if gateway == "" {
			return invalidComposeYamlRetVal(fmt.Sprintf("network %s gateway is required", name))
		}
		if _, ok := networkGateways[gateway]; ok {
			return invalidComposeYamlRetVal(fmt.Sprintf("network %s gateway must be unique", name))
		}
		networkGateways[gateway] = struct{}{}
	}

	serviceNames := make(map[string]struct{})
	containerNames := make(map[string]struct{})
	nonKubeVirtServiceCount := 0
	openPortCount := 0
	assignedIPs := make(map[string]map[string]struct{})

	for _, service := range config.Services {
		name := strings.TrimSpace(service.Name)
		if name == "" {
			return invalidComposeYamlRetVal("service name is required")
		}
		if _, ok := serviceNames[name]; ok {
			return invalidComposeYamlRetVal(fmt.Sprintf("service %s name must be unique", name))
		}
		serviceNames[name] = struct{}{}
		label := name

		containerName := strings.TrimSpace(service.ContainerName)
		if containerName != "" {
			if _, ok := containerNames[containerName]; ok {
				return invalidComposeYamlRetVal(fmt.Sprintf("%s container_name must be unique", label))
			}
			containerNames[containerName] = struct{}{}
		}

		if strings.TrimSpace(service.Image) == "" {
			return model.RetVal{Msg: i18n.Model.Challenge.EmptyImage}
		}
		kubeVirt := serviceKubeVirt(service)
		if kubeVirt && !hasVpcNetworks {
			return invalidComposeYamlRetVal(fmt.Sprintf("%s x-kubevirt requires VPC networks", label))
		}
		if boot, ok := service.Extensions[model.XBootExtension].(model.XBoot); ok {
			bootloader := strings.TrimSpace(boot.Bootloader)
			if kubeVirt && bootloader != "" && bootloader != "bios" && bootloader != "efi" {
				return invalidComposeYamlRetVal(fmt.Sprintf("%s bootloader should be bios or efi", label))
			}
		}

		if !kubeVirt {
			nonKubeVirtServiceCount++
			servicePortTargets := make(map[uint32]struct{})
			for _, port := range service.Ports {
				if port.Target == 0 {
					return invalidComposeYamlRetVal(fmt.Sprintf("%s port target is required", label))
				}
				if _, ok := servicePortTargets[port.Target]; ok {
					return invalidComposeYamlRetVal(fmt.Sprintf("%s port target must be unique within the same service", label))
				}
				servicePortTargets[port.Target] = struct{}{}
				if strings.TrimSpace(port.Published) != "" && !composePortNameRegex.MatchString(port.Published) {
					return invalidComposeYamlRetVal(fmt.Sprintf("%s port name can only contain letters, numbers, underscores, and hyphens", label))
				}
				protocol := strings.TrimSpace(port.Protocol)
				if protocol != "" && protocol != "tcp" && protocol != "udp" {
					return invalidComposeYamlRetVal(fmt.Sprintf("%s port protocol should be tcp or udp", label))
				}
				openPortCount++
			}

			for key := range service.Environment {
				if strings.TrimSpace(key) == "" {
					return invalidComposeYamlRetVal(fmt.Sprintf("%s environment variable name is required", label))
				}
				if !composeEnvKeyRegex.MatchString(key) {
					return invalidComposeYamlRetVal(fmt.Sprintf("%s environment variable name format is invalid", label))
				}
			}

			volumeTargets := make(map[string]struct{})
			if volumes, ok := service.Extensions[model.XVolumesExtension].(model.XVolumes); ok {
				for _, volume := range volumes {
					target := strings.TrimSpace(volume.Path)
					if target == "" {
						return invalidComposeYamlRetVal(fmt.Sprintf("%s file Flag mount path is required", label))
					}
					if _, ok := volumeTargets[target]; ok {
						return invalidComposeYamlRetVal(fmt.Sprintf("%s file Flag target must be unique within the same service", label))
					}
					volumeTargets[target] = struct{}{}
				}
			}
		}

		if hasVpcNetworks && len(service.Networks) == 0 {
			return invalidComposeYamlRetVal(fmt.Sprintf("%s choose a network and fill in IP after configuring custom networks", label))
		}
		serviceNetworkNames := make(map[string]struct{})
		for networkName, network := range service.Networks {
			if network == nil {
				return invalidComposeYamlRetVal(fmt.Sprintf("%s empty network config", label))
			}
			name := strings.TrimSpace(networkName)
			if name == "" {
				return invalidComposeYamlRetVal(fmt.Sprintf("%s network name is required", label))
			}
			if _, ok := serviceNetworkNames[name]; ok {
				return invalidComposeYamlRetVal(fmt.Sprintf("%s network cannot be selected more than once", label))
			}
			serviceNetworkNames[name] = struct{}{}
			ipv4Address := strings.TrimSpace(network.Ipv4Address)
			if ipv4Address == "" {
				return invalidComposeYamlRetVal(fmt.Sprintf("%s %s ipv4_address is required", label, name))
			}
			ip, err := netip.ParseAddr(ipv4Address)
			if err != nil || !ip.Is4() {
				return invalidComposeYamlRetVal(fmt.Sprintf("%s %s ipv4_address is invalid", label, name))
			}
			if _, ok := assignedIPs[name]; !ok {
				assignedIPs[name] = make(map[string]struct{})
			}
			if _, ok := assignedIPs[name][ipv4Address]; ok {
				return invalidComposeYamlRetVal(fmt.Sprintf("%s %s ipv4_address must be unique within the same network", label, name))
			}
			assignedIPs[name][ipv4Address] = struct{}{}

			macAddress := strings.TrimSpace(network.MacAddress)
			if kubeVirt && macAddress == "" {
				return invalidComposeYamlRetVal(fmt.Sprintf("%s %s mac_address is required when x-kubevirt is true", label, name))
			}
			if macAddress != "" {
				if _, err := net.ParseMAC(macAddress); err != nil {
					return invalidComposeYamlRetVal(fmt.Sprintf("%s %s mac_address is invalid", label, name))
				}
			}
		}
	}

	if nonKubeVirtServiceCount > 0 && openPortCount == 0 {
		return invalidComposeYamlRetVal("Container challenges must expose at least one port")
	}
	return model.SuccessRetVal()
}

func GetChallenges(tx *gorm.DB, form dto.GetChallengesForm) ([]model.Challenge, int64, model.RetVal) {
	options := db.GetOptions{
		Conditions: make(map[string]any),
		Search:     make(map[string]string),
		Preloads:   map[string]db.GetOptions{"ChallengeFlags": {}},
	}
	if form.Type != "" {
		options.Conditions["type"] = form.Type
	}
	if form.Category != "" {
		options.Conditions["category"] = form.Category
	}
	if form.Name != "" {
		options.Search["name"] = form.Name
	}
	if form.Description != "" {
		options.Search["description"] = form.Description
	}
	return db.InitChallengeRepo(tx).List(form.Limit, form.Offset, options)
}

func buildChallengeTemplate(dockerCompose string) (model.ChallengeTemplate, []db.CreateChallengeFlagOptions, model.RetVal) {
	prefix := utils.RandHexStr(10)
	config, err := utils.LoadDockerComposeYaml(dockerCompose, prefix, map[string]any{
		model.XVolumesExtension:   model.XVolumes{},
		model.XKubeVirtExtension:  false,
		model.XBootExtension:      model.XBoot{},
		model.XCloudInitExtension: model.XCloudInit{},
	})
	if err != nil {
		log.Logger.Warningf("Failed to load DockerCompose: %v", err)
		return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	if ret := validateChallengeCompose(config); !ret.OK {
		return model.ChallengeTemplate{}, nil, ret
	}
	prefix = fmt.Sprintf("%s_", prefix)
	networksMap := make(map[string]model.NetworkDefinition)
	for _, network := range config.Networks {
		network.Name = strings.TrimPrefix(network.Name, prefix)
		if network.Name == "default" {
			continue
		}
		if len(network.Ipam.Config) == 0 {
			return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": "Empty IPAM"}}
		}
		subnet, err := netip.ParsePrefix(network.Ipam.Config[0].Subnet)
		if err != nil {
			return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": err.Error()}}
		}
		gateway, err := netip.ParseAddr(network.Ipam.Config[0].Gateway)
		if err != nil || !subnet.Contains(gateway) {
			return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": "Invalid gateway"}}
		}
		networksMap[network.Name] = model.NetworkDefinition{
			External: bool(network.External),
			Name:     network.Name,
			CIDR:     network.Ipam.Config[0].Subnet,
			Gateway:  network.Ipam.Config[0].Gateway,
		}
	}

	template := model.ChallengeTemplate{
		Pods: make([]model.ChallengePodTemplate, 0, len(config.Services)),
	}
	flagOptions := make([]db.CreateChallengeFlagOptions, 0)
	for _, app := range config.Services {
		name := app.Name
		if app.ContainerName != "" {
			name = app.ContainerName
		}
		if app.Image == "" {
			return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Challenge.EmptyImage}
		}
		containerKey := strings.ToLower(name)
		podKey := containerKey
		environment := make(model.StringMap)
		for k, v := range app.Environment {
			if !strings.HasPrefix(k, model.EnvFlagPrefix) {
				environment[k] = *v
			}
		}
		ports := make(model.Exposes, 0)
		seenPorts := make([]string, 0)
		for _, port := range app.Ports {
			target := fmt.Sprintf("%d/%s", port.Target, port.Protocol)
			if !slices.Contains(seenPorts, target) {
				ports = append(ports, model.Expose{
					Port:      int32(port.Target),
					Protocol:  port.Protocol,
					Published: port.Published,
				})
				seenPorts = append(seenPorts, target)
			}
		}
		networks := make(model.Networks, 0)
		for key, value := range app.Networks {
			if key == "default" {
				continue
			}
			if value == nil {
				return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": "Empty network name"}}
			}
			ip, err := netip.ParseAddr(value.Ipv4Address)
			if err != nil {
				return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": "Invalid ip"}}
			}
			networkDefinition, ok := networksMap[key]
			if !ok {
				log.Logger.Warningf("Network %s not found in networks", key)
				return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": "Invalid network"}}
			}
			subnet, err := netip.ParsePrefix(networkDefinition.CIDR)
			if err != nil {
				return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": err.Error()}}
			}
			if !subnet.Contains(ip) {
				return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": "Invalid subnet"}}
			}
			networks = append(networks, model.Network{
				Definition: networkDefinition,
				Attachment: model.NetworkAttachment{
					Name: key,
					IP:   value.Ipv4Address,
					MAC:  value.MacAddress,
				},
			})
		}
		if len(networksMap) > 0 && len(networks) == 0 {
			return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": "Invalid networks"}}
		}
		containerTemplate := model.ChallengeContainerTemplate{
			Key:         containerKey,
			Name:        name,
			Image:       app.Image,
			CPU:         app.CPUS,
			Memory:      int64(app.MemLimit),
			WorkingDir:  app.WorkingDir,
			Command:     model.StringList(app.Command),
			Environment: environment,
			Exposes:     append(model.Exposes(nil), ports...),
		}
		if boot, ok := app.Extensions[model.XBootExtension].(model.XBoot); ok {
			containerTemplate.Bootloader = boot.Bootloader
			containerTemplate.SecureBoot = boot.SecureBoot
		}
		if kubeVirt, ok := app.Extensions[model.XKubeVirtExtension].(bool); ok {
			containerTemplate.KubeVirt = kubeVirt
		}
		if cloudInit, ok := app.Extensions[model.XCloudInitExtension].(model.XCloudInit); ok {
			containerTemplate.UserData = cloudInit.UserData
		}
		template.Pods = append(template.Pods, model.ChallengePodTemplate{
			Key:          podKey,
			Name:         name,
			ServicePorts: append(model.Exposes(nil), ports...),
			Networks:     networks,
			Containers:   []model.ChallengeContainerTemplate{containerTemplate},
		})
		for k, v := range app.Environment {
			if strings.HasPrefix(k, model.EnvFlagPrefix) {
				flagOptions = append(flagOptions, db.CreateChallengeFlagOptions{
					Value: *v,
					Binding: model.FlagBinding{
						PodKey:       podKey,
						ContainerKey: containerKey,
						Type:         model.EnvFlagBindingType,
						Target:       k,
					},
				})
			}
		}
		if volumes, ok := app.Extensions[model.XVolumesExtension].(model.XVolumes); ok {
			for _, volume := range volumes {
				if volume.Path == "" || volume.Content == "" {
					continue
				}
				flagOptions = append(flagOptions, db.CreateChallengeFlagOptions{
					Value: volume.Content,
					Binding: model.FlagBinding{
						PodKey:       podKey,
						ContainerKey: containerKey,
						Type:         model.FileFlagBindingType,
						Target:       volume.Path,
					},
				})
			}
		}
	}
	return template, flagOptions, model.SuccessRetVal()
}

func CreateChallenge(tx *gorm.DB, form dto.CreateChallengeForm) (model.Challenge, model.RetVal) {
	challengeRepo, challengeFlagRepo := db.InitChallengeRepo(tx), db.InitChallengeFlagRepo(tx)
	options := db.CreateChallengeOptions{
		RandID:          utils.UUID(),
		Name:            form.Name,
		Description:     form.Description,
		Type:            form.Type,
		Category:        form.Category,
		GeneratorImage:  form.GeneratorImage,
		NetworkPolicies: form.NetworkPolicies,
	}
	var podFlagOptions []db.CreateChallengeFlagOptions
	if form.Type == model.PodsChallengeType {
		template, flags, ret := buildChallengeTemplate(form.DockerCompose)
		if !ret.OK {
			return model.Challenge{}, ret
		}
		options.Template = template
		podFlagOptions = flags
	}
	challenge, ret := challengeRepo.Create(options)
	if !ret.OK {
		return model.Challenge{}, ret
	}
	switch form.Type {
	case model.StaticChallengeType:
		for _, flag := range form.Flags {
			if _, ret = challengeFlagRepo.Create(db.CreateChallengeFlagOptions{
				ChallengeID: challenge.ID,
				Value:       flag,
			}); !ret.OK {
				return model.Challenge{}, ret
			}
		}
	case model.DynamicChallengeType:
		for _, flag := range form.Flags {
			if _, ret = challengeFlagRepo.Create(db.CreateChallengeFlagOptions{
				ChallengeID: challenge.ID,
				Value:       flag,
			}); !ret.OK {
				return model.Challenge{}, ret
			}
		}
	case model.PodsChallengeType:
		for _, flag := range podFlagOptions {
			flag.ChallengeID = challenge.ID
			if _, ret = challengeFlagRepo.Create(flag); !ret.OK {
				return model.Challenge{}, ret
			}
		}
	default:
		return model.Challenge{}, model.RetVal{Msg: i18n.Model.Challenge.InvalidType}
	}
	return challengeRepo.GetByID(challenge.ID, db.GetOptions{
		Preloads: map[string]db.GetOptions{"ChallengeFlags": {}},
	})
}

func UpdateChallenge(tx *gorm.DB, challenge model.Challenge, form dto.UpdateChallengeForm) model.RetVal {
	switch challenge.Type {
	case model.StaticChallengeType, model.DynamicChallengeType:
		oldChallengeFlagID := make([]uint, 0)
		for _, flag := range challenge.ChallengeFlags {
			oldChallengeFlagID = append(oldChallengeFlagID, flag.ID)
		}
		challengeFlagRepo := db.InitChallengeFlagRepo(tx)
		for _, flag := range form.Flags {
			if slices.Contains(oldChallengeFlagID, flag.ID) {
				if ret := challengeFlagRepo.Update(flag.ID, db.UpdateChallengeFlagOptions{
					Value: &flag.Value,
				}); !ret.OK {
					return ret
				}
				oldChallengeFlagID = slices.DeleteFunc(oldChallengeFlagID, func(id uint) bool {
					return id == flag.ID
				})
			} else {
				if _, ret := challengeFlagRepo.Create(db.CreateChallengeFlagOptions{
					ChallengeID: challenge.ID,
					Value:       flag.Value,
				}); !ret.OK {
					return ret
				}
			}
		}
		if ret := challengeFlagRepo.Delete(oldChallengeFlagID...); !ret.OK {
			return ret
		}
		return db.InitChallengeRepo(tx).Update(challenge.ID, db.UpdateChallengeOptions{
			Name:           form.Name,
			Description:    form.Description,
			Category:       form.Category,
			GeneratorImage: form.GeneratorImage,
		})
	case model.PodsChallengeType:
		if form.DockerCompose != nil {
			challengeFlagIDL := make([]uint, 0, len(challenge.ChallengeFlags))
			for _, flag := range challenge.ChallengeFlags {
				challengeFlagIDL = append(challengeFlagIDL, flag.ID)
			}
			if ret := db.InitChallengeFlagRepo(tx).Delete(challengeFlagIDL...); !ret.OK {
				return ret
			}
			template, flags, ret := buildChallengeTemplate(*form.DockerCompose)
			if !ret.OK {
				return ret
			}
			if ret := db.InitChallengeRepo(tx).Update(challenge.ID, db.UpdateChallengeOptions{
				Name:            form.Name,
				Description:     form.Description,
				Category:        form.Category,
				NetworkPolicies: form.NetworkPolicies,
				Template:        &template,
			}); !ret.OK {
				return ret
			}
			repo := db.InitChallengeFlagRepo(tx)
			for _, flag := range flags {
				flag.ChallengeID = challenge.ID
				if _, ret := repo.Create(flag); !ret.OK {
					return ret
				}
			}
			return model.SuccessRetVal()
		}
		return db.InitChallengeRepo(tx).Update(challenge.ID, db.UpdateChallengeOptions{
			Name:            form.Name,
			Description:     form.Description,
			Category:        form.Category,
			NetworkPolicies: form.NetworkPolicies,
		})
	default:
		return model.RetVal{Msg: i18n.Model.Challenge.InvalidType}
	}
}
