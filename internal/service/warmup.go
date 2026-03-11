package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
)

func GetNodeImageList() (map[string][]string, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return k8s.GetNodeImageList(ctx)
}

func WarmUpContestChallengeImage(form dto.WarmUpImageForm) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if form.PullPolicy == string(corev1.PullNever) {
		return model.SuccessRetVal()
	}
	nodes, ret := k8s.ListSchedulableNodes(ctx)
	if !ret.OK {
		return ret
	}
	for _, node := range nodes {
		images := form.Images
		if corev1.PullPolicy(form.PullPolicy) != corev1.PullAlways {
			images = slices.DeleteFunc(images, func(image string) bool {
				if strings.TrimSpace(image) == "" {
					return true
				}
				for _, containerImage := range node.Status.Images {
					for _, name := range containerImage.Names {
						if name == image {
							return true
						}
					}
				}
				return false
			})
		}
		if len(images) > 0 {
			var chunks [][]string
			for i := 0; i < len(images); i += 5 {
				end := i + 5
				if end > len(images) {
					end = len(images)
				}
				chunks = append(chunks, images[i:end])
			}
			for _, chunk := range chunks {
				if _, ret = k8s.CreateJob(ctx, k8s.CreateJobOptions{
					Name:       fmt.Sprintf("image-puller-%s", utils.RandStr(5)),
					Images:     chunk,
					PullPolicy: form.PullPolicy,
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": node.Name,
					},
				}); !ret.OK {
					return ret
				}
			}
		}
	}
	return model.SuccessRetVal()
}

func GetContestVictims(tx *gorm.DB, contest model.Contest, form dto.GetContestVictimsForm) ([]model.Victim, int64, model.RetVal) {
	var challengeID uint
	if form.ChallengeID != "" {
		challenge, ret := db.InitChallengeRepo(tx).GetByRandID(form.ChallengeID)
		if !ret.OK || challenge.Type != model.PodsChallengeType {
			return nil, 0, ret
		}
		challengeID = challenge.ID
	}
	options := db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads: map[string]db.GetOptions{
			"Pods":             {},
			"User":             {},
			"Team":             {},
			"ContestChallenge": {},
		},
	}
	if challengeID != 0 {
		options.Conditions["challenge_id"] = challengeID
	}
	if form.TeamID != 0 {
		options.Conditions["team_id"] = form.TeamID
	}
	if form.UserID != 0 {
		options.Conditions["user_id"] = form.UserID
	}
	victims, count, ret := db.InitVictimRepo(tx).List(form.Limit, form.Offset, options)
	return slices.DeleteFunc(victims, func(victim model.Victim) bool {
		if !victim.TeamID.Valid || !victim.ContestChallengeID.Valid || !victim.ContestID.Valid {
			count--
			return true
		}
		return false
	}), count, ret
}

func StartContestVictims(tx *gorm.DB, contest model.Contest, form dto.StartContestVictimsForm) model.RetVal {
	if len(form.Challenges) == 0 || len(form.Teams) == 0 {
		return model.SuccessRetVal()
	}
	challenges, _, ret := db.InitChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"type": model.PodsChallengeType},
	})
	if !ret.OK {
		return ret
	}
	challengeIDL := make([]uint, 0)
	for _, challenge := range challenges {
		challengeIDL = append(challengeIDL, challenge.ID)
	}
	teams, _, ret := db.InitTeamRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "id": form.Teams},
	})
	if !ret.OK {
		return ret
	}
	if len(challengeIDL) == 0 || len(teams) == 0 {
		return model.SuccessRetVal()
	}
	contestChallenges, _, ret := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "challenge_id": challengeIDL},
		Preloads:   map[string]db.GetOptions{"ContestFlags": {}},
	})
	if !ret.OK {
		return ret
	}
	if len(contestChallenges) == 0 {
		return model.SuccessRetVal()
	}
	victimRepo := db.InitVictimRepo(tx)
	for _, contestChallenge := range contestChallenges {
		for _, team := range teams {
			if CheckIfSolved(tx, team, contestChallenge.ContestFlags) {
				continue
			}
			if !CheckIfGenerated(tx, team, contestChallenge.ContestFlags) {
				if _, ret = CreateTeamFlag(tx, team, contest, contestChallenge); !ret.OK {
					continue
				}
			}
			_, ret = StartVictim(tx, team.CaptainID, team.ID, contest.ID, contestChallenge.ID, contestChallenge.ChallengeID)
			if !ret.OK {
				victim, ret := victimRepo.HasAliveVictim(team.ID, contestChallenge.ChallengeID)
				if !ret.OK {
					continue
				}
				StopVictim(tx, victim)
			}
		}
	}
	return model.SuccessRetVal()
}

// StopContestVictims tx 无需开启事务
func StopContestVictims(tx *gorm.DB, form dto.StopContestVictimsForm) model.RetVal {
	if len(form.Victims) == 0 {
		return model.SuccessRetVal()
	}
	victims, _, ret := db.InitVictimRepo(tx).List(-1, -1, db.GetOptions{Conditions: map[string]any{"id": form.Victims}})
	if !ret.OK {
		return ret
	}
	for _, victim := range victims {
		if ret = StopVictim(tx, victim); !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal()
}
