package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"net/netip"
	"slices"
	"strings"

	"gorm.io/gorm"
)

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
	config, ret := utils.LoadDockerComposeYaml(dockerCompose)
	if !ret.OK || config == nil {
		return model.ChallengeTemplate{}, nil, ret
	}
	volumeFlag := make(map[string]string)
	for _, volume := range config.Volumes {
		volumeName := strings.TrimPrefix(volume.Name, "_")
		if strings.HasPrefix(volumeName, model.VolumeFlagPrefix) {
			for k, v := range volume.Labels {
				if k == model.VolumeFlagLabelKey {
					volumeFlag[volumeName] = v
				}
			}
		}
	}
	networksMap := make(map[string]model.Network)
	for _, network := range config.Networks {
		network.Name = strings.TrimPrefix(network.Name, "_")
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
		networksMap[network.Name] = model.Network{
			External: network.External.External,
			Name:     network.Name,
			CIDR:     network.Ipam.Config[0].Subnet,
			Gateway:  network.Ipam.Config[0].Gateway,
		}
	}

	template := model.ChallengeTemplate{
		Version: 1,
		Pods:    make([]model.ChallengePodTemplate, 0, len(config.Services)),
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
				ports = append(ports, model.Expose{Port: int32(port.Target), Protocol: port.Protocol})
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
			network, ok := networksMap[key]
			if !ok {
				log.Logger.Warningf("Network %s not found in networks", key)
				return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": "Invalid network"}}
			}
			subnet, err := netip.ParsePrefix(network.CIDR)
			if err != nil {
				return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": err.Error()}}
			}
			if !subnet.Contains(ip) {
				return model.ChallengeTemplate{}, nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": "Invalid subnet"}}
			}
			network.IP = value.Ipv4Address
			networks = append(networks, network)
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
					Name:  k,
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
		for _, volume := range app.Volumes {
			if value, ok := volumeFlag[volume.Source]; ok {
				flagOptions = append(flagOptions, db.CreateChallengeFlagOptions{
					Name:  volume.Source,
					Value: value,
					Binding: model.FlagBinding{
						PodKey:       podKey,
						ContainerKey: containerKey,
						Type:         model.FileFlagBindingType,
						Target:       volume.Target,
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
		Options:         form.Options,
		NetworkPolicies: form.NetworkPolicies,
		TemplateVersion: 1,
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
	case model.QuestionChallengeType:
		answer := make([]string, 0)
		for _, option := range form.Options {
			if option.Correct {
				answer = append(answer, option.RandID)
			}
		}
		if _, ret = challengeFlagRepo.Create(db.CreateChallengeFlagOptions{
			ChallengeID: challenge.ID,
			Value:       strings.Join(answer, ","),
		}); !ret.OK {
			return model.Challenge{}, ret
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
			Options:        form.Options,
			GeneratorImage: form.GeneratorImage,
		})
	case model.QuestionChallengeType:
		if form.Options != nil {
			answer := ""
			for _, option := range *form.Options {
				if option.Correct {
					answer += fmt.Sprintf("%s,", option.RandID)
				}
			}
			answer = strings.TrimSuffix(answer, ",")
			repo := db.InitChallengeFlagRepo(tx)
			if len(challenge.ChallengeFlags) > 0 {
				if ret := repo.Update(challenge.ChallengeFlags[0].ID, db.UpdateChallengeFlagOptions{
					Value: &answer,
				}); !ret.OK {
					return ret
				}
			} else {
				if _, ret := repo.Create(db.CreateChallengeFlagOptions{
					ChallengeID: challenge.ID,
					Value:       answer,
				}); !ret.OK {
					return ret
				}
			}
		}
		return db.InitChallengeRepo(tx).Update(challenge.ID, db.UpdateChallengeOptions{
			Name:        form.Name,
			Description: form.Description,
			Category:    form.Category,
			Options:     form.Options,
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
			version := challenge.TemplateVersion + 1
			if ret := db.InitChallengeRepo(tx).Update(challenge.ID, db.UpdateChallengeOptions{
				Name:            form.Name,
				Description:     form.Description,
				Category:        form.Category,
				NetworkPolicies: form.NetworkPolicies,
				Template:        &template,
				TemplateVersion: &version,
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
