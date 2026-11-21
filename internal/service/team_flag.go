package service

import (
	"fmt"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/task"
	"CBCTF/internal/utils"
	"CBCTF/internal/view"
)

func CreateTeamFlags(tx *gorm.DB, team model.Team, contest model.Contest) model.RetVal {
	contestChallenges, _, ret := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"ContestFlags": {}},
	})
	if !ret.OK {
		return ret
	}
	for _, contestChallenge := range contestChallenges {
		ret = db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
			teamFlags, ret := CreateTeamFlag(tx2, team, contest, contestChallenge)
			if !ret.OK {
				return ret
			}
			if err := enqueueDynamicAttachment(team.CaptainID, team, contestChallenge, teamFlags); err != nil {
				return model.FailRetVal(
					i18n.Model.CreateError,
					model.Name(model.TeamFlag{}),
					fmt.Sprintf("team_id=%d contest_challenge_id=%d", team.ID, contestChallenge.ID),
					err,
				)
			}
			return model.SuccessRetVal()
		})
		if !ret.OK {
			return model.FailRetVal(
				i18n.Model.CreateError,
				model.Name(model.TeamFlag{}),
				fmt.Sprintf("team_id=%d contest_challenge_id=%d", team.ID, contestChallenge.ID),
				ret.Cause(),
			)
		}
	}
	return model.SuccessRetVal()
}

// enqueueDynamicAttachment queues the attachment generation of a dynamic
// challenge for the team. Static challenges need no attachment.
func enqueueDynamicAttachment(userID uint, team model.Team, contestChallenge model.ContestChallenge, teamFlags []model.TeamFlag) error {
	if contestChallenge.Type != model.DynamicChallengeType {
		return nil
	}
	if err := task.EnqueueGenAttachmentTask(userID, model.NewContestChallengeSource(contestChallenge), team.ID, teamFlags); err != nil {
		log.Logger.Warningf("Failed to enqueue gen attachment task: %s", err)
		return err
	}
	return nil
}

// applyTeamChallenge applies a team flag change to a contest challenge and
// queues the attachment of a dynamic challenge. The contest challenge must
// carry its flags so that the change can resolve them.
func applyTeamChallenge(
	tx *gorm.DB,
	user model.User,
	team model.Team,
	contestChallenge model.ContestChallenge,
	apply func(*gorm.DB) ([]model.TeamFlag, model.RetVal),
) model.RetVal {
	contestFlags, _, ret := db.InitContestFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		return ret
	}
	contestChallenge.ContestFlags = contestFlags
	return db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		teamFlags, applyRet := apply(tx2)
		if !applyRet.OK {
			return applyRet
		}
		if err := enqueueDynamicAttachment(user.ID, team, contestChallenge, teamFlags); err != nil {
			return model.FailRetVal(i18n.Task.EnqueueError, "TeamFlag", fmt.Sprintf("team_id=%d,contest_challenge_id=%d", team.ID, contestChallenge.ID), err)
		}
		return model.SuccessRetVal()
	})
}

func renderChallengeFlagValue(value string) string {
	if result := model.StaticFlagTmpl.FindAllStringSubmatch(value, 1); len(result) > 0 {
		return result[0][1]
	}
	if result := model.DynamicFlagTmpl.FindAllStringSubmatch(value, 1); len(result) > 0 {
		return utils.Leet(result[0][1])
	}
	if result := model.UUIDFlagTmpl.FindAllStringSubmatch(value, 1); len(result) > 0 {
		return utils.UUID()
	}
	return value
}

// CreateTeamFlag 需要 contestChallenge 预加载 ContestFlags
func CreateTeamFlag(tx *gorm.DB, team model.Team, contest model.Contest, contestChallenge model.ContestChallenge) ([]model.TeamFlag, model.RetVal) {
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	teamFlagL := make([]model.TeamFlag, 0)
	for _, contestFlag := range contestChallenge.ContestFlags {
		teamFlag, ret := teamFlagRepo.Get(db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID, "contest_flag_id": contestFlag.ID},
		})
		if ret.OK {
			teamFlagL = append(teamFlagL, teamFlag)
			continue
		}
		if ret.Msg != i18n.Model.NotFound {
			return nil, ret
		}
		options := model.TeamFlag{
			TeamID:          team.ID,
			ContestFlagID:   contestFlag.ID,
			ChallengeFlagID: contestFlag.ChallengeFlagID,
			Solved:          false,
			Value:           renderChallengeFlagValue(contestFlag.Value),
		}
		if prefix := contest.Prefix; prefix != "" {
			options.Value = fmt.Sprintf("%s{%s}", contest.Prefix, options.Value)
		}
		teamFlag, ret = teamFlagRepo.CreateIfAbsent(options)
		if !ret.OK {
			return nil, ret
		}
		teamFlagL = append(teamFlagL, teamFlag)
	}
	return teamFlagL, model.SuccessRetVal()
}

// UpdateTeamFlag 需要 contestChallenge 预加载 ContestFlags
func UpdateTeamFlag(tx *gorm.DB, team model.Team, contest model.Contest, contestChallenge model.ContestChallenge) ([]model.TeamFlag, model.RetVal) {
	submissionRepo := db.InitSubmissionRepo(tx)
	submissions, _, ret := submissionRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID, "contest_challenge_id": contestChallenge.ID, "solved": true},
	})
	if !ret.OK {
		return nil, ret
	}
	submissionIDL := make([]uint, 0)
	for _, submission := range submissions {
		submissionIDL = append(submissionIDL, submission.ID)
	}
	teamFlagIDL := make([]uint, 0)
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	for _, contestFlag := range contestChallenge.ContestFlags {
		teamFlag, ret := teamFlagRepo.Get(db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID, "contest_flag_id": contestFlag.ID},
		})
		if !ret.OK {
			return nil, ret
		}
		teamFlagIDL = append(teamFlagIDL, teamFlag.ID)
	}
	if ret = submissionRepo.Delete(submissionIDL...); !ret.OK {
		return nil, ret
	}
	if ret = teamFlagRepo.Delete(teamFlagIDL...); !ret.OK {
		return nil, ret
	}
	return CreateTeamFlag(tx, team, contest, contestChallenge)
}

func CheckIfGenerated(tx *gorm.DB, team model.Team, contestFlags []model.ContestFlag) (bool, model.RetVal) {
	if len(contestFlags) == 0 {
		return true, model.SuccessRetVal()
	}
	contestFlagIDL := make([]uint, 0, len(contestFlags))
	for _, contestFlag := range contestFlags {
		contestFlagIDL = append(contestFlagIDL, contestFlag.ID)
	}
	count, ret := db.InitTeamFlagRepo(tx).CountGenerated(team.ID, contestFlagIDL...)
	if !ret.OK {
		return false, ret
	}
	return count == int64(len(contestFlags)), model.SuccessRetVal()
}

func ListTeamFlagViews(tx *gorm.DB, team model.Team) ([]view.TeamFlagChallengeView, model.RetVal) {
	teamFlags, _, ret := db.InitTeamFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID},
		Preloads: map[string]db.GetOptions{"ContestFlag": {
			Preloads: map[string]db.GetOptions{"ContestChallenge": {}},
		}},
	})
	if !ret.OK {
		return nil, ret
	}

	groupIndices := make(map[uint]int)
	result := make([]view.TeamFlagChallengeView, 0)
	for _, flag := range teamFlags {
		id := flag.ContestFlag.ContestChallengeID
		index, ok := groupIndices[id]
		if !ok {
			result = append(result, view.TeamFlagChallengeView{
				Name:     flag.ContestFlag.ContestChallenge.Name,
				Type:     flag.ContestFlag.ContestChallenge.Type,
				Category: flag.ContestFlag.ContestChallenge.Category,
				Hidden:   flag.ContestFlag.ContestChallenge.Hidden,
				Flags:    make([]view.TeamFlagInfoView, 0),
			})
			index = len(result) - 1
			groupIndices[id] = index
		}
		result[index].Flags = append(result[index].Flags, view.TeamFlagInfoView{
			Value:        flag.Value,
			Solved:       flag.Solved,
			Template:     flag.ContestFlag.Value,
			InitScore:    flag.ContestFlag.Score,
			CurrentScore: flag.ContestFlag.CurrentScore,
			Decay:        flag.ContestFlag.Decay,
			MinScore:     flag.ContestFlag.MinScore,
			Solvers:      flag.ContestFlag.Solvers,
		})
	}
	return result, model.SuccessRetVal()
}

func InitTeamChallenge(tx *gorm.DB, user model.User, team model.Team, contest model.Contest, contestChallenge model.ContestChallenge) model.RetVal {
	return applyTeamChallenge(tx, user, team, contestChallenge, func(tx2 *gorm.DB) ([]model.TeamFlag, model.RetVal) {
		return CreateTeamFlag(tx2, team, contest, contestChallenge)
	})
}

func ResetTeamChallenge(tx *gorm.DB, user model.User, team model.Team, contest model.Contest, contestChallenge model.ContestChallenge) model.RetVal {
	ret := applyTeamChallenge(tx, user, team, contestChallenge, func(tx2 *gorm.DB) ([]model.TeamFlag, model.RetVal) {
		return UpdateTeamFlag(tx2, team, contest, contestChallenge)
	})
	if !ret.OK {
		return ret
	}

	if contestChallenge.Type == model.PodsChallengeType {
		queueTeamVictimStop(tx, team.ID, contestChallenge.ID)
	}
	return model.SuccessRetVal()
}
