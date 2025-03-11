package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateSubmission is a function to create a new submission
func CreateSubmission(tx *gorm.DB, contest model.Contest, team model.Team, user model.User, usage model.Usage, value string) (model.Submission, bool, string) {
	if usage.Attempt != 0 && usage.Attempt <= CountAttempts(tx, contest.ID, team.ID, usage.ChallengeID) {
		return model.Submission{}, false, "NotAllowSubmit"
	}
	if _, ok, msg := GetFlagBy3ID(tx, contest.ID, team.ID, usage.ChallengeID); !ok {
		return model.Submission{}, false, msg
	}
	solved := VerifyFlag(tx, contest.ID, team.ID, usage.ChallengeID, value)
	if solved {
		if ok, msg := Solve(tx, usage.ID, team.ID, contest.Blood); !ok {
			return model.Submission{}, false, msg
		}
	}
	team, ok, msg := GetTeamByID(tx, team.ID)
	if !ok {
		return model.Submission{}, false, msg
	}
	submission := model.InitSubmission(usage.ID, contest.ID, usage.ChallengeID, team.ID, user.ID, value, solved, team.Score)
	if err := tx.Model(model.Submission{}).Create(&submission).Error; err != nil {
		return model.Submission{}, false, "CreateSubmissionError"
	}
	return submission, true, "Success"
}

func IsSolved(tx *gorm.DB, contestID, teamID uint, challengeID string) bool {
	var submission model.Submission
	res := tx.Model(model.Submission{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ? AND solved = ?", contestID, teamID, challengeID, true).Find(&submission).Limit(1)
	if res.RowsAffected < 1 {
		return false
	}
	return true
}

func CountAttempts(tx *gorm.DB, contestID, teamID uint, challengeID string) int64 {
	var count int64
	res := tx.Model(model.Submission{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count attempts: %v", res.Error)
		return 0
	}
	return count
}

// GetSubmissions is a function to get submissions
func GetSubmissions(tx *gorm.DB, limit, offset int, teamIDL ...uint) ([]model.Submission, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	var submissions []model.Submission
	var count int64
	res := tx.Model(model.Submission{})
	if len(teamIDL) > 0 {
		res = res.Where("team_id IN ?", teamIDL)
	}
	if res.Count(&count).Error != nil {
		log.Logger.Warningf("Failed to count submissions: %v", res.Error)
		return make([]model.Submission, 0), 0, false, "UnknownError"
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	if res = res.Order("created_at desc").Limit(limit).Offset(offset).Find(&submissions); res.Error != nil {
		log.Logger.Warningf("Failed to get submissions: %v", res.Error)
		return nil, 0, false, "GetSubmissionError"
	}
	return submissions, count, true, "Success"
}

func GetTeamSolved(tx *gorm.DB, teamID uint) ([]model.Submission, bool, string) {
	var submissions []model.Submission
	res := tx.Model(model.Submission{}).Order("created_at asc").
		Where("team_id = ? AND solved = ?", teamID, true).Find(&submissions)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get submissions: %v", res.Error)
		return make([]model.Submission, 0), false, "GetSubmissionError"
	}
	return submissions, true, "Success"
}

func CalcTeamScore(tx *gorm.DB, contestID, teamID uint) (float64, bool, string) {
	solved, ok, msg := GetTeamSolved(tx, teamID)
	if !ok {
		return 0, false, msg
	}
	var score float64
	for _, submission := range solved {
		usage, ok, msg := GetUsageBy2ID(tx, contestID, submission.ChallengeID)
		if !ok {
			log.Logger.Warningf("Failed to get usage: %s", msg)
			continue
		}
		rate := 0.0
		for {
			if usage.First == teamID {
				rate = 0.05
				break
			}
			if usage.Second == teamID {
				rate = 0.03
				break
			}
			if usage.Third == teamID {
				rate = 0.01
				break
			}
			break
		}
		if rate > 0 {
			score += usage.Score * (1 + rate)
		} else {
			score += usage.CurrentScore
		}
	}
	return score, true, "Success"
}

func GetTeamSolvedState(tx *gorm.DB, team model.Team) ([]gin.H, bool, string) {
	solved, ok, msg := GetTeamSolved(tx, team.ID)
	if !ok {
		return make([]gin.H, 0), false, msg
	}
	usages, ok, msg := GetUsageByContestID(tx, team.ContestID, false)
	if !ok {
		return make([]gin.H, 0), false, msg
	}
	categories := make(map[string]string)
	for _, v := range usages {
		categories[v.ChallengeID] = v.Category
	}
	allCount := make(map[string]int64)
	for _, usage := range usages {
		allCount[usage.Category] += 1
	}
	solvedCount := make(map[string]int64)
	for _, submission := range solved {
		solvedCount[categories[submission.ChallengeID]] += 1
	}
	var tmp []gin.H
	for k, v := range allCount {
		if _, ok := solvedCount[k]; !ok {
			solvedCount[k] = 0
		}
		tmp = append(tmp, gin.H{"category": k, "solved": solvedCount[k], "all": v})
	}
	return tmp, true, "Success"
}

func GetContestSolved(tx *gorm.DB, contestID uint) ([]model.Submission, bool, string) {
	var submissions []model.Submission
	res := tx.Model(model.Submission{}).
		Where("contest_id = ? AND solved = ?", contestID, true).Find(&submissions)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get submissions: %v", res.Error)
		return make([]model.Submission, 0), false, "GetSubmissionError"
	}
	return submissions, true, "Success"
}
