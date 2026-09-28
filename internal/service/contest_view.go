package service

import (
	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"CBCTF/internal/view"
)

func BuildContestView(tx *gorm.DB, contest model.Contest) view.ContestView {
	contestRepo := db.InitContestRepo(tx)
	teamCount, _ := contestRepo.CountTeams(contest.ID)
	userCount, _ := contestRepo.CountUsers(contest.ID)
	noticeCount, _ := contestRepo.CountNotices(contest.ID)
	result := view.ContestView{
		Contest:     contest,
		TeamCount:   teamCount,
		UserCount:   userCount,
		NoticeCount: noticeCount,
		StatsReady:  true,
	}
	champion, _, _ := GetTeamRanking(tx, contest, 1, 0)
	if len(champion) > 0 {
		result.Highest = champion[0].Score
	}
	result.SolvedCount, _ = db.InitSubmissionRepo(tx).Count(db.CountOptions{
		Conditions: map[string]any{"solved": true, "contest_id": contest.ID},
	})
	return result
}

func ListContests(tx *gorm.DB, form dto.ListModelsForm, admin bool) ([]view.ContestView, int64, model.RetVal) {
	options := db.GetOptions{Sort: []string{"id DESC"}}
	if !admin {
		options.Conditions = map[string]any{"hidden": false}
	}
	contests, count, ret := db.InitContestRepo(tx).List(form.Limit, form.Offset, options)
	if !ret.OK {
		return nil, 0, ret
	}
	views := make([]view.ContestView, 0, len(contests))
	if len(contests) == 0 {
		return views, count, model.SuccessRetVal()
	}

	contestIDs := make([]uint, 0, len(contests))
	for _, contest := range contests {
		contestIDs = append(contestIDs, contest.ID)
	}

	contestRepo := db.InitContestRepo(tx)
	teamCountMap, _ := contestRepo.CountTeamsMap(contestIDs...)
	userCountMap, _ := contestRepo.CountUsersMap(contestIDs...)
	noticeCountMap := make(map[uint]int64, len(contests))
	for _, contest := range contests {
		noticeCount, _ := contestRepo.CountNotices(contest.ID)
		noticeCountMap[contest.ID] = noticeCount
	}

	for _, contest := range contests {
		views = append(views, view.ContestView{
			Contest:     contest,
			TeamCount:   teamCountMap[contest.ID],
			UserCount:   userCountMap[contest.ID],
			NoticeCount: noticeCountMap[contest.ID],
		})
	}
	return views, count, model.SuccessRetVal()
}

func DeleteContest(tx *gorm.DB, contest model.Contest) model.RetVal {
	return db.WithTransactionDB(tx, func(tx *gorm.DB) model.RetVal {
		if ret := stopGeneratorResources(tx, db.GetOptions{Conditions: map[string]any{"contest_id": contest.ID}}); !ret.OK {
			return ret
		}
		return db.InitContestRepo(tx).Delete(contest.ID)
	})
}
