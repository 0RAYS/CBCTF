package service

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/k8s"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"slices"
	"strings"
	"time"
)

func StartTeamVictim(tx *gorm.DB, user model.User, team model.Team, contestChallenge model.ContestChallenge) (model.Victim, bool, string) {
	challenge, ok, msg := db.InitChallengeRepo(tx).
		GetByID(contestChallenge.ChallengeID, db.GetOptions{
			Preloads: map[string]db.GetOptions{
				"Dockers": {
					Preloads: map[string]db.GetOptions{
						"ChallengeFlags": {},
					},
				},
			},
		})
	if !ok {
		return model.Victim{}, false, msg
	}
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	victimRepo := db.InitVictimRepo(tx)
	podRepo := db.InitPodRepo(tx)
	containerRepo := db.InitContainerRepo(tx)
	if victim, ok, _ := victimRepo.HasAliveVictim(team.ID, contestChallenge.ID); ok {
		return victim, true, i18n.Success
	}
	podPorts := make(model.Int32List, 0)
	for _, docker := range challenge.Dockers {
		for _, port := range docker.Exposes {
			if !slices.Contains(podPorts, port.Port) {
				podPorts = append(podPorts, port.Port)
			}
		}
	}
	vOptions := db.CreateVictimOptions{
		ContestChallengeID: contestChallenge.ID,
		TeamID:             team.ID,
		UserID:             user.ID,
		Start:              time.Now(),
		Duration:           time.Hour,
	}
	victim, ok, msg := victimRepo.Create(vOptions)
	if !ok {
		return model.Victim{}, false, msg
	}
	pOptions := db.CreatePodOptions{
		VictimID: victim.ID,
		Name:     victim.GenPodName(challenge.RandID),
		PodPorts: podPorts,
	}
	pod, ok, msg := podRepo.Create(pOptions)
	if !ok {
		return model.Victim{}, false, msg
	}
	for _, docker := range challenge.Dockers {
		envFlagL := make(model.StringList, 0)
		volumeFlagL := make(model.StringMap)
		for _, challengeFlag := range docker.ChallengeFlags {
			teamFlag, ok, msg := teamFlagRepo.Get(db.GetOptions{
				Conditions: map[string]any{
					"team_id":           team.ID,
					"challenge_flag_id": challengeFlag.ID,
				},
			})
			if !ok {
				return model.Victim{}, false, msg
			}
			switch challengeFlag.InjectType {
			case model.EnvInjectType:
				envFlagL = append(envFlagL, teamFlag.Value)
			case model.VolumeInjectType:
				volumeFlagL[challengeFlag.Path] = teamFlag.Value
			default:
				return model.Victim{}, false, i18n.InvalidChallengeFlagInjectType
			}
		}
		cOptions := db.CreateContainerOptions{
			PodID:       pod.ID,
			Name:        fmt.Sprintf("%s-%s", pod.Name, strings.ToLower(utils.RandStr(5))),
			Image:       docker.Image,
			Hostname:    docker.Name,
			WorkingDir:  docker.WorkingDir,
			Command:     docker.Command,
			Environment: docker.Environment,
			EnvFlags:    envFlagL,
			VolumeFlags: volumeFlagL,
			Exposes:     docker.Exposes,
		}
		container, ok, msg := containerRepo.Create(cOptions)
		if !ok {
			return model.Victim{}, false, msg
		}
		pod.Containers = append(pod.Containers, container)
	}
	victim.Pods = append(victim.Pods, pod)
	ip, ports, ok, msg := k8s.StartVictim(victim)
	if !ok {
		go k8s.StopVictim(victim)
		return model.Victim{}, false, msg
	}
	ok, msg = podRepo.Update(pod.ID, db.UpdatePodOptions{
		ExposedIP:    &ip,
		ExposedPorts: (*model.Int32List)(&ports),
	})
	if !ok {
		return model.Victim{}, false, msg
	}
	victim.Pods[0].ExposedIP = ip
	victim.Pods[0].ExposedPorts = ports
	return victim, true, i18n.Success
}

// GetTeamVictimStatus contestChallenge 需要预加载 model.ContestChallenge
func GetTeamVictimStatus(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) gin.H {
	data := gin.H{
		"target":    make([]string, 0),
		"remaining": 0,
		"status":    "Down",
	}
	if contestChallenge.Challenge.Type != model.PodsChallengeType {
		data["status"] = "NotDocker"
		return data
	}
	victims, _, ok, _ := db.InitVictimRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{
			"contest_challenge_id": contestChallenge.ID,
			"team_id":              team.ID,
		},
		Preloads: map[string]db.GetOptions{
			"Pods": {},
		},
	})
	if !ok {
		return data
	}
	if len(victims) == 0 {
		return data
	}
	if len(victims) > 1 {
		data["status"] = "Error"
		return data
	}
	for _, pod := range victims[0].Pods {
		data["target"] = append(data["target"].([]string), pod.RemoteAddr()...)
	}
	data["status"] = "Running"
	data["remaining"] = victims[0].Remaining().Seconds()
	return data
}

func StopTeamVictim(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) (bool, string) {
	victimRepo := db.InitVictimRepo(tx)
	victims, _, ok, msg := victimRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{
			"team_id":              team.ID,
			"contest_challenge_id": contestChallenge.ID,
		},
		Preloads: map[string]db.GetOptions{
			"Pods": {},
		},
	})
	if !ok {
		return false, msg
	}
	// 预期中, len(victims) == 1, 考虑意外情况
	victimIDL := make([]uint, 0)
	for _, victim := range victims {
		ok, msg = k8s.StopVictim(victim)
		if !ok {
			return false, msg
		}
		duration := time.Now().Sub(victim.Start)
		if ok, msg = victimRepo.Update(victim.ID, db.UpdateVictimOptions{
			Duration: &duration,
		}); !ok {
			return false, msg
		}
		victimIDL = append(victimIDL, victim.ID)
		LoadTraffic(tx, victim)
	}
	return victimRepo.Delete(victimIDL...)
}
