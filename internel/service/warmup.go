package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/k8s"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"context"
	"fmt"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	"strings"
	"time"
)

func GetNodeImageList() (map[string][]string, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return k8s.GetNodeImageList(ctx)
}

func WarmUpContestChallengeImage(form f.WarmUpImageForm) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if form.PullPolicy == string(corev1.PullNever) {
		return true, i18n.Success
	}
	nodes, ok, msg := k8s.ListSchedulableNodes(ctx)
	if !ok {
		return false, msg
	}
	for _, node := range nodes {
		if _, ok, msg = k8s.CreateJob(ctx, k8s.CreateJobOptions{
			Name:       fmt.Sprintf("image-puller-%s", utils.RandStr(5)),
			Images:     form.Images,
			PullPolicy: form.PullPolicy,
			NodeSelector: map[string]string{
				"kubernetes.io/hostname": node.Name,
			},
		}); !ok {
			return false, msg
		}
	}
	return true, i18n.Success
}

func GetContestVictims(tx *gorm.DB, contest model.Contest, form f.GetContestVictimsForm) ([]model.Victim, int64, bool, string) {
	var (
		victims            = make([]model.Victim, 0)
		contestChallengeID uint
	)
	if form.ChallengeID != "" {
		challenge, ok, msg := db.InitChallengeRepo(tx).GetByRandID(form.ChallengeID, db.GetOptions{
			Selects: []string{"id", "type"},
		})
		if !ok || challenge.Type != model.PodsChallengeType {
			return victims, 0, false, msg
		}
		contestChallenge, ok, msg := db.InitContestChallengeRepo(tx).Get(db.GetOptions{
			Conditions: map[string]any{
				"contest_id":   contest.ID,
				"challenge_id": challenge.ID,
			},
			Selects: []string{"id"},
		})
		if !ok {
			return victims, 0, false, msg
		}
		contestChallengeID = contestChallenge.ID
	}
	options := db.GetOptions{
		Conditions: make(map[string]any),
		Preloads: map[string]db.GetOptions{
			"Pods":             {},
			"User":             {Selects: []string{"id", "name"}},
			"Team":             {Selects: []string{"id", "name"}},
			"ContestChallenge": {Selects: []string{"id", "name"}},
		},
	}
	if contestChallengeID != 0 {
		options.Conditions["contest_challenge_id"] = contestChallengeID
	}
	if form.TeamID != 0 {
		options.Conditions["team_id"] = form.TeamID
	}
	if form.UserID != 0 {
		options.Conditions["user_id"] = form.UserID
	}
	return db.InitVictimRepo(tx).List(form.Limit, form.Offset, options)
}

func StartContestVictims(tx *gorm.DB, contest model.Contest, form f.StartContestVictimsForm) (bool, string) {
	if len(form.Challenges) == 0 || len(form.Teams) == 0 {
		return true, i18n.Success
	}
	challengeIDL := make([]uint, 0)
	for _, randID := range form.Challenges {
		challenge, ok, _ := db.InitChallengeRepo(tx).GetByRandID(randID, db.GetOptions{
			Conditions: map[string]any{"type": model.PodsChallengeType},
		})
		if !ok {
			continue
		}
		challengeIDL = append(challengeIDL, challenge.ID)
	}
	teams := make([]model.Team, 0)
	for _, id := range form.Teams {
		team, ok, _ := db.InitTeamRepo(tx).GetByID(id, db.GetOptions{
			Conditions: map[string]any{"contest_id": contest.ID},
		})
		if !ok {
			continue
		}
		teams = append(teams, team)
	}
	if len(challengeIDL) == 0 || len(teams) == 0 {
		return true, i18n.Success
	}
	contestChallenges, _, ok, msg := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "challenge_id": challengeIDL},
		Preloads: map[string]db.GetOptions{
			"ContestFlags": {},
		},
	})
	if !ok {
		return false, msg
	}
	if len(contestChallenges) == 0 {
		return true, i18n.Success
	}
	go func(teams []model.Team, contestChallenges []model.ContestChallenge) {
		for _, contestChallenge := range contestChallenges {
			for _, team := range teams {
				if !CheckIfGenerated(db.DB, team, contestChallenge) {
					tx2 := db.DB.Begin()
					if _, ok, msg = CreateTeamFlag(tx2, team, contestChallenge); !ok {
						tx2.Rollback()
						continue
					}
					tx2.Commit()
				}
				tx2 := db.DB.Begin()
				_, ok, msg = StartTeamVictim(tx, model.User{BasicModel: model.BasicModel{ID: team.CaptainID}}, team, contestChallenge)
				if !ok {
					tx2.Rollback()
					continue
				}
				tx2.Commit()
			}
		}
	}(teams, contestChallenges)
	return true, i18n.Success
}

func StopContestVictims(tx *gorm.DB, form f.StopContestVictimsForm) (bool, string) {
	if len(form.Victims) == 0 {
		return true, i18n.Success
	}
	victimRepo := db.InitVictimRepo(tx)
	victims, _, ok, msg := victimRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{
			"id": form.Victims,
		},
		Preloads: map[string]db.GetOptions{
			"Pods": {},
		},
	})
	if !ok {
		return false, msg
	}
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
