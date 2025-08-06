package service

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"CBCTF/internal/utils"
	"fmt"
	"slices"
	"strings"

	"gorm.io/gorm"
)

func GetChallenges(tx *gorm.DB, form f.GetChallengesForm) ([]model.Challenge, int64, bool, string) {
	options := db.GetOptions{
		Conditions: make(map[string]any),
		Preloads:   map[string]db.GetOptions{"Dockers": {}, "ChallengeFlags": {}},
	}
	if form.Type != "" {
		options.Conditions["type"] = form.Type
	}
	if form.Category != "" {
		options.Conditions["category"] = form.Category
	}
	return db.InitChallengeRepo(tx).List(form.Limit, form.Offset, options)
}

func CreateChallenge(tx *gorm.DB, form f.CreateChallengeForm) (model.Challenge, bool, string) {
	challengeRepo, challengeFlagRepo := db.InitChallengeRepo(tx), db.InitChallengeFlagRepo(tx)
	challenge, ok, msg := challengeRepo.Create(db.CreateChallengeOptions{
		RandID:          utils.UUID(),
		Name:            form.Name,
		Desc:            form.Desc,
		Type:            form.Type,
		Category:        form.Category,
		GeneratorImage:  form.GeneratorImage,
		Options:         form.Options,
		NetworkPolicies: form.NetworkPolicies,
	})
	if !ok {
		return model.Challenge{}, false, msg
	}
	switch form.Type {
	case model.StaticChallengeType:
		for _, flag := range form.Flags {
			if _, ok, msg = challengeFlagRepo.Create(db.CreateChallengeFlagOptions{
				ChallengeID: challenge.ID,
				Value:       flag,
			}); !ok {
				return model.Challenge{}, false, msg
			}
		}
	case model.QuestionChallengeType:
		answer := ""
		for _, option := range form.Options {
			if option.Correct {
				answer += fmt.Sprintf("%s,", option.RandID)
			}
		}
		answer = strings.TrimSuffix(answer, ",")
		if _, ok, msg = challengeFlagRepo.Create(db.CreateChallengeFlagOptions{
			ChallengeID: challenge.ID,
			Value:       answer,
		}); !ok {
			return model.Challenge{}, false, msg
		}
	case model.DynamicChallengeType:
		for _, flag := range form.Flags {
			if _, ok, msg = challengeFlagRepo.Create(db.CreateChallengeFlagOptions{
				ChallengeID: challenge.ID,
				Value:       flag,
			}); !ok {
				return model.Challenge{}, false, msg
			}
		}
	case model.PodsChallengeType:
		config, ok, msg := utils.LoadDockerComposeYaml(form.DockerCompose)
		if !ok {
			return model.Challenge{}, false, msg
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
			tmp := model.Network{External: network.External.External, Name: network.Name}
			if len(network.Ipam.Config) > 0 {
				tmp.CIDR = network.Ipam.Config[0].Subnet
				tmp.Gateway = network.Ipam.Config[0].Gateway
			}
			networksMap[network.Name] = tmp
		}
		flagOptions := make([]db.CreateChallengeFlagOptions, 0)
		dockerRepo := db.InitDockerRepo(tx)
		for _, app := range config.Services {
			name := app.Name
			if app.ContainerName != "" {
				name = app.ContainerName
			}
			if app.Image == "" {
				return model.Challenge{}, false, i18n.InvalidDockerImage
			}
			environment := make(model.StringMap)
			for k, v := range app.Environment {
				if !strings.HasPrefix(k, model.EnvFlagPrefix) {
					environment[k] = *v
				}
			}
			ports := make(model.Exposes, 0)
			tmp := make([]string, 0)
			for _, port := range app.Ports {
				target := fmt.Sprintf("%d/%s", port.Target, port.Protocol)
				if !slices.Contains(tmp, target) {
					ports = append(ports, model.Expose{Port: int32(port.Target), Protocol: port.Protocol})
					tmp = append(tmp, target)
				}
			}
			networks := make(model.Networks, 0)
			for key, value := range app.Networks {
				if key == "default" {
					continue
				}
				network, ok := networksMap[key]
				if !ok {
					log.Logger.Warningf("Network %s not found in networks", key)
					return model.Challenge{}, false, i18n.UnknownError
				}
				network.IP = value.Ipv4Address
				networks = append(networks, network)
			}
			docker, ok, msg := dockerRepo.Create(db.CreateDockerOptions{
				ChallengeID: challenge.ID,
				Name:        name,
				Image:       app.Image,
				CPU:         app.CPUS,
				Memory:      int64(app.MemLimit),
				WorkingDir:  app.WorkingDir,
				Command:     model.StringList(app.Command),
				Exposes:     ports,
				Environment: environment,
				Networks:    networks,
			})
			if !ok {
				return model.Challenge{}, false, msg
			}
			for k, v := range app.Environment {
				if strings.HasPrefix(k, model.EnvFlagPrefix) {
					flagOptions = append(flagOptions, db.CreateChallengeFlagOptions{
						ChallengeID: challenge.ID,
						DockerID:    &docker.ID,
						Name:        k,
						Value:       *v,
						InjectType:  model.EnvInjectType,
					})
				}
			}
			for _, volume := range app.Volumes {
				if value, ok := volumeFlag[volume.Source]; ok {
					flagOptions = append(flagOptions, db.CreateChallengeFlagOptions{
						ChallengeID: challenge.ID,
						DockerID:    &docker.ID,
						Name:        volume.Source,
						Value:       value,
						InjectType:  model.VolumeInjectType,
						Path:        volume.Target,
					})
				}
			}
		}
		for _, options := range flagOptions {
			if _, ok, msg = challengeFlagRepo.Create(options); !ok {
				return model.Challenge{}, false, msg
			}
		}
	default:
		return model.Challenge{}, false, i18n.InvalidChallengeType
	}
	return challengeRepo.GetByID(challenge.ID, db.GetOptions{
		Preloads: map[string]db.GetOptions{"Dockers": {}, "ChallengeFlags": {}},
	})
}

func UpdateChallenge(tx *gorm.DB, challenge model.Challenge, form f.UpdateChallengeForm) (bool, string) {
	switch challenge.Type {
	case model.StaticChallengeType, model.DynamicChallengeType:
		oldChallengeFlagID := make([]uint, 0)
		for _, flag := range challenge.ChallengeFlags {
			oldChallengeFlagID = append(oldChallengeFlagID, flag.ID)
		}
		challengeFlagRepo := db.InitChallengeFlagRepo(tx)
		for _, flag := range form.Flags {
			if slices.Contains(oldChallengeFlagID, flag.ID) {
				if ok, msg := challengeFlagRepo.Update(flag.ID, db.UpdateChallengeFlagOptions{
					Value: &flag.Value,
				}); !ok {
					return false, msg
				}
				oldChallengeFlagID = slices.DeleteFunc(oldChallengeFlagID, func(id uint) bool {
					return id == flag.ID
				})
			} else {
				if _, ok, msg := challengeFlagRepo.Create(db.CreateChallengeFlagOptions{
					ChallengeID: challenge.ID,
					Value:       flag.Value,
				}); !ok {
					return false, msg
				}
			}
		}
		if ok, msg := challengeFlagRepo.Delete(oldChallengeFlagID...); !ok {
			return false, msg
		}
		return db.InitChallengeRepo(tx).Update(challenge.ID, db.UpdateChallengeOptions{
			Name:           form.Name,
			Desc:           form.Desc,
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
				if ok, msg := repo.Update(challenge.ChallengeFlags[0].ID, db.UpdateChallengeFlagOptions{
					Value: &answer,
				}); !ok {
					return false, msg
				}
			} else {
				if _, ok, msg := repo.Create(db.CreateChallengeFlagOptions{
					ChallengeID: challenge.ID,
					Value:       answer,
				}); !ok {
					return false, msg
				}
			}
		}
		return db.InitChallengeRepo(tx).Update(challenge.ID, db.UpdateChallengeOptions{
			Name:     form.Name,
			Desc:     form.Desc,
			Category: form.Category,
			Options:  form.Options,
		})
	case model.PodsChallengeType:
		return db.InitChallengeRepo(tx).Update(challenge.ID, db.UpdateChallengeOptions{
			Name:            form.Name,
			Desc:            form.Desc,
			Category:        form.Category,
			NetworkPolicies: form.NetworkPolicies,
		})
	default:
		return false, i18n.InvalidChallengeType
	}
}
