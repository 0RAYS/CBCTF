package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
	"strings"
)

func GetChallenges(tx *gorm.DB, form f.GetChallengesForm) ([]model.Challenge, int64, bool, string) {
	var conditions db.GetOptions
	if form.Type != "" {
		conditions = append(conditions, db.GetOption{Key: "type", Value: form.Type, Op: "and"})
	}
	if form.Category != "" {
		conditions = append(conditions, db.GetOption{Key: "category", Value: utils.ToTitle(form.Category), Op: "and"})
	}
	return db.InitChallengeRepo(tx).ListWithConditions(
		form.Limit, form.Offset, conditions,
		"DockerGroups", "ChallengeFlags", "DockerGroups.Dockers",
	)
}

func CreateChallenge(tx *gorm.DB, form f.CreateChallengeForm) (model.Challenge, bool, string) {
	challengeRepo, challengeFlagRepo := db.InitChallengeRepo(tx), db.InitChallengeFlagRepo(tx)
	challenge, ok, msg := challengeRepo.Create(db.CreateChallengeOptions{
		RandID:   utils.UUID(),
		Name:     form.Name,
		Desc:     form.Desc,
		Type:     form.Type,
		Category: utils.ToTitle(form.Category),
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
	case model.DynamicChallengeType:
		if ok, msg = challengeRepo.Update(challenge.ID, db.UpdateChallengeOptions{
			GeneratorImage: &form.GeneratorImage,
		}); !ok {
			return model.Challenge{}, false, msg
		}
		for _, flag := range form.Flags {
			if _, ok, msg = challengeFlagRepo.Create(db.CreateChallengeFlagOptions{
				ChallengeID: challenge.ID,
				Value:       flag,
			}); !ok {
				return model.Challenge{}, false, msg
			}
		}
	case model.PodsChallengeType:
		dockerGroupRepo, dockerRepo := db.InitDockerGroupRepo(tx), db.InitDockerRepo(tx)
		for _, group := range form.DockerGroups {
			if len(group.NetworkPolicies) == 0 {
				group.NetworkPolicies = append(group.NetworkPolicies, model.DefaultNetworkPolicy)
			}
			dockerGroup, ok, msg := dockerGroupRepo.Create(db.CreateDockerGroupOptions{
				ChallengeID:     challenge.ID,
				NetworkPolicies: group.NetworkPolicies,
			})
			if !ok {
				return model.Challenge{}, false, msg
			}
			config, ok, msg := utils.LoadDockerComposeYaml(group.Yaml)
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
			flagOptions := make([]db.CreateChallengeFlagOptions, 0)
			for _, service := range config.Services {
				name := service.Name
				if service.ContainerName != "" {
					name = service.ContainerName
				}
				if service.Image == "" {
					return model.Challenge{}, false, i18n.InvalidDockerImage
				}
				environment := make(model.StringMap)
				for k, v := range service.Environment {
					if !strings.HasPrefix(k, model.EnvFlagPrefix) {
						environment[k] = *v
					}
				}
				docker, ok, msg := dockerRepo.Create(db.CreateDockerOptions{
					DockerGroupID: dockerGroup.ID,
					Name:          name,
					Image:         service.Image,
					PullPolicy:    &service.PullPolicy,
					WorkingDir:    &service.WorkingDir,
					Command:       (*model.StringList)(&service.Command),
					Expose:        (*model.StringList)(&service.Expose),
					Environment:   &environment,
				})
				if !ok {
					return model.Challenge{}, false, msg
				}
				for k, v := range service.Environment {
					if strings.HasPrefix(k, model.EnvFlagPrefix) {
						flagOptions = append(flagOptions, db.CreateChallengeFlagOptions{
							ChallengeID: challenge.ID,
							DockerID:    &docker.ID,
							Value:       *v,
							InjectType:  model.EnvInjectType,
						})
					}
				}
				for _, volume := range service.Volumes {
					if value, ok := volumeFlag[volume.Source]; ok {
						flagOptions = append(flagOptions, db.CreateChallengeFlagOptions{
							ChallengeID: challenge.ID,
							DockerID:    &docker.ID,
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
		}
	default:
		return model.Challenge{}, false, i18n.InvalidChallengeType
	}
	return challengeRepo.GetByID(challenge.ID, "DockerGroups", "ChallengeFlags", "DockerGroups.Dockers")
}

func UpdateChallenge(tx *gorm.DB, challenge model.Challenge, form f.UpdateChallengeForm) (bool, string) {
	switch challenge.Type {
	case model.StaticChallengeType, model.DynamicChallengeType:
		if form.Flags != nil {
			challengeFlagIDL := make([]uint, 0)
			for _, flag := range challenge.ChallengeFlags {
				challengeFlagIDL = append(challengeFlagIDL, flag.ID)
			}
			if ok, msg := db.InitChallengeFlagRepo(tx).Delete(challengeFlagIDL...); !ok {
				return false, msg
			}
		}
		return db.InitChallengeRepo(tx).Update(challenge.ID, db.UpdateChallengeOptions{
			Name:           form.Name,
			Desc:           form.Desc,
			Category:       form.Category,
			GeneratorImage: form.GeneratorImage,
		})
	case model.PodsChallengeType:
		return db.InitChallengeRepo(tx).Update(challenge.ID, db.UpdateChallengeOptions{
			Name:     form.Name,
			Desc:     form.Desc,
			Category: form.Category,
		})
	default:
		return false, i18n.InvalidChallengeType
	}
}
