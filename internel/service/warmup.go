package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/k8s"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"context"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
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
