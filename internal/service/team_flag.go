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
	challenges, _, ret := db.InitChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"ChallengeFlags": {}},
	})
	if !ret.OK {
		return ret
	}
	for _, challenge := range challenges {
		ret = db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
			teamFlags, ret := CreateTeamFlag(tx2, team, contest, challenge)
			if !ret.OK {
				return ret
			}
			if err := enqueueDynamicAttachment(team.CaptainID, team, challenge, teamFlags); err != nil {
				return model.FailRetVal(
					i18n.Model.CreateError,
					model.Name(model.TeamFlag{}),
					fmt.Sprintf("team_id=%d challenge_id=%d", team.ID, challenge.ID),
					err,
				)
			}
			return model.SuccessRetVal()
		})
		if !ret.OK {
			return model.FailRetVal(
				i18n.Model.CreateError,
				model.Name(model.TeamFlag{}),
				fmt.Sprintf("team_id=%d challenge_id=%d", team.ID, challenge.ID),
				ret.Cause(),
			)
		}
	}
	return model.SuccessRetVal()
}

// enqueueDynamicAttachment queues the attachment generation of a dynamic
// challenge for the team. Static challenges need no attachment.
func enqueueDynamicAttachment(userID uint, team model.Team, challenge model.Challenge, teamFlags []model.TeamFlag) error {
	if challenge.Type != model.DynamicChallengeType {
		return nil
	}
	if err := task.EnqueueGenAttachmentTask(userID, model.NewChallengeSource(challenge), team.ID, teamFlags); err != nil {
		log.Logger.Warningf("Failed to enqueue gen attachment task: %s", err)
		return err
	}
	return nil
}

// applyTeamChallenge applies a team flag change to a challenge and
// queues the attachment of a dynamic challenge. The challenge must
// carry its flags so that the change can resolve them.
func applyTeamChallenge(
	tx *gorm.DB,
	user model.User,
	team model.Team,
	challenge model.Challenge,
	apply func(*gorm.DB) ([]model.TeamFlag, model.RetVal),
) model.RetVal {
	flags, _, ret := db.InitChallengeFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID},
	})
	if !ret.OK {
		return ret
	}
	challenge.ChallengeFlags = flags
	return db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		teamFlags, applyRet := apply(tx2)
		if !applyRet.OK {
			return applyRet
		}
		if err := enqueueDynamicAttachment(user.ID, team, challenge, teamFlags); err != nil {
			return model.FailRetVal(i18n.Task.EnqueueError, "TeamFlag", fmt.Sprintf("team_id=%d,challenge_id=%d", team.ID, challenge.ID), err)
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

// CreateTeamFlag 需要 challenge 预加载 ChallengeFlags
func CreateTeamFlag(tx *gorm.DB, team model.Team, contest model.Contest, challenge model.Challenge) ([]model.TeamFlag, model.RetVal) {
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	teamFlagL := make([]model.TeamFlag, 0)
	for _, flag := range challenge.ChallengeFlags {
		teamFlag, ret := teamFlagRepo.Get(db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID, "challenge_flag_id": flag.ID},
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
			ChallengeFlagID: flag.ID,
			Solved:          false,
			Value:           renderChallengeFlagValue(flag.Value),
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

// UpdateTeamFlag 需要 challenge 预加载 ChallengeFlags
func UpdateTeamFlag(tx *gorm.DB, team model.Team, contest model.Contest, challenge model.Challenge) ([]model.TeamFlag, model.RetVal) {
	submissionRepo := db.InitSubmissionRepo(tx)
	submissions, _, ret := submissionRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID, "challenge_id": challenge.ID, "solved": true},
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
	for _, flag := range challenge.ChallengeFlags {
		teamFlag, ret := teamFlagRepo.Get(db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID, "challenge_flag_id": flag.ID},
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
	return CreateTeamFlag(tx, team, contest, challenge)
}

func CheckIfGenerated(tx *gorm.DB, team model.Team, flags []model.ChallengeFlag) (bool, model.RetVal) {
	if len(flags) == 0 {
		return true, model.SuccessRetVal()
	}
	challengeFlagIDL := make([]uint, 0, len(flags))
	for _, flag := range flags {
		challengeFlagIDL = append(challengeFlagIDL, flag.ID)
	}
	count, ret := db.InitTeamFlagRepo(tx).CountGenerated(team.ID, challengeFlagIDL...)
	if !ret.OK {
		return false, ret
	}
	return count == int64(len(flags)), model.SuccessRetVal()
}

func ListTeamFlagViews(tx *gorm.DB, team model.Team) ([]view.TeamFlagChallengeView, model.RetVal) {
	teamFlags, _, ret := db.InitTeamFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID},
		Preloads: map[string]db.GetOptions{"ChallengeFlag": {
			Preloads: map[string]db.GetOptions{"Challenge": {}},
		}},
	})
	if !ret.OK {
		return nil, ret
	}

	groupIndices := make(map[uint]int)
	result := make([]view.TeamFlagChallengeView, 0)
	for _, flag := range teamFlags {
		id := flag.ChallengeFlag.ChallengeID
		index, ok := groupIndices[id]
		if !ok {
			result = append(result, view.TeamFlagChallengeView{
				Name:     flag.ChallengeFlag.Challenge.Name,
				Type:     flag.ChallengeFlag.Challenge.Type,
				Category: flag.ChallengeFlag.Challenge.Category,
				Hidden:   flag.ChallengeFlag.Challenge.Hidden,
				Flags:    make([]view.TeamFlagInfoView, 0),
			})
			index = len(result) - 1
			groupIndices[id] = index
		}
		result[index].Flags = append(result[index].Flags, view.TeamFlagInfoView{
			Value:        flag.Value,
			Solved:       flag.Solved,
			Template:     flag.ChallengeFlag.Value,
			InitScore:    flag.ChallengeFlag.Score,
			CurrentScore: flag.ChallengeFlag.CurrentScore,
			Decay:        flag.ChallengeFlag.Decay,
			MinScore:     flag.ChallengeFlag.MinScore,
			Solvers:      flag.ChallengeFlag.Solvers,
		})
	}
	return result, model.SuccessRetVal()
}

func InitTeamChallenge(tx *gorm.DB, user model.User, team model.Team, contest model.Contest, challenge model.Challenge) model.RetVal {
	return applyTeamChallenge(tx, user, team, challenge, func(tx2 *gorm.DB) ([]model.TeamFlag, model.RetVal) {
		return CreateTeamFlag(tx2, team, contest, challenge)
	})
}

func ResetTeamChallenge(tx *gorm.DB, user model.User, team model.Team, contest model.Contest, challenge model.Challenge) model.RetVal {
	ret := applyTeamChallenge(tx, user, team, challenge, func(tx2 *gorm.DB) ([]model.TeamFlag, model.RetVal) {
		return UpdateTeamFlag(tx2, team, contest, challenge)
	})
	if !ret.OK {
		return ret
	}

	if challenge.Type == model.PodsChallengeType {
		queueTeamVictimStop(tx, team.ID, challenge.ID)
	}
	return model.SuccessRetVal()
}
