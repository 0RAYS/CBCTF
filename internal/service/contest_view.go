package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"CBCTF/internal/view"

	"gorm.io/gorm"
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

func BuildContestViews(tx *gorm.DB, contests []model.Contest) []view.ContestView {
	views := make([]view.ContestView, 0, len(contests))
	if len(contests) == 0 {
		return views
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
	return views
}

func GetContestView(tx *gorm.DB, contest model.Contest) view.ContestView {
	return BuildContestView(tx, contest)
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
	return BuildContestViews(tx, contests), count, model.SuccessRetVal()
}

func DeleteContest(tx *gorm.DB, contest model.Contest) model.RetVal {
	return db.InitContestRepo(tx).Delete(contest.ID)
}

func DeleteContestWithTransaction(tx *gorm.DB, contest model.Contest) model.RetVal {
	return db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		return DeleteContest(tx2, contest)
	})
}
