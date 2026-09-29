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
	teamCount, teamRet := contestRepo.CountTeams(contest.ID)
	userCount, userRet := contestRepo.CountUsers(contest.ID)
	noticeCount, noticeRet := contestRepo.CountNotices(contest.ID)
	result := view.ContestView{
		Contest:     contest,
		TeamCount:   teamCount,
		UserCount:   userCount,
		NoticeCount: noticeCount,
		StatsReady:  true,
	}
	for field, ret := range map[string]model.RetVal{"teams": teamRet, "users": userRet, "notices": noticeRet} {
		if !ret.OK {
			result.Unavailable = append(result.Unavailable, field)
		}
	}
	champion, _, rankRet := GetTeamRanking(tx, contest, 1, 0)
	if !rankRet.OK {
		result.Unavailable = append(result.Unavailable, "highest")
	}
	if len(champion) > 0 {
		result.Highest = champion[0].Score
	}
	var solveRet model.RetVal
	result.SolvedCount, solveRet = db.InitSubmissionRepo(tx).Count(db.CountOptions{
		Conditions: map[string]any{"solved": true, "contest_id": contest.ID},
	})
	if !solveRet.OK {
		result.Unavailable = append(result.Unavailable, "solved")
	}
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
	teamCountMap, teamRet := contestRepo.CountTeamsMap(contestIDs...)
	userCountMap, userRet := contestRepo.CountUsersMap(contestIDs...)
	noticeCountMap := make(map[uint]int64, len(contests))
	missingNotice := make(map[uint]bool)
	for _, contest := range contests {
		noticeCount, ret := contestRepo.CountNotices(contest.ID)
		missingNotice[contest.ID] = !ret.OK
		noticeCountMap[contest.ID] = noticeCount
	}

	for _, contest := range contests {
		missing := make([]string, 0)
		if !teamRet.OK {
			missing = append(missing, "teams")
		}
		if !userRet.OK {
			missing = append(missing, "users")
		}
		if missingNotice[contest.ID] {
			missing = append(missing, "notices")
		}
		views = append(views, view.ContestView{
			Unavailable: missing,
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
