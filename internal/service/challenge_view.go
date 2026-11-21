package service

import (
	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/view"
)

// challengeFilter 把请求参数转换成题库聚合视图的过滤条件。
func challengeFilter(contestID *uint, t model.ChallengeType, category, name, description string) db.ChallengeFilter {
	filter := db.ChallengeFilter{
		Type:        t,
		Category:    category,
		Name:        name,
		Description: description,
	}
	if contestID != nil {
		filter.ContestID = *contestID
		filter.Unbound = *contestID == 0
	}
	return filter
}

func BuildChallengeView(tx *gorm.DB, challenge model.Challenge) view.ChallengeView {
	result := view.ChallengeView{
		Challenge: challenge,
		Flags:     make([]view.ChallengeFlagView, 0),
	}
	if challenge.ContestID.Valid {
		if contest, ret := db.InitContestRepo(tx).GetByID(challenge.ContestID.V); ret.OK {
			result.ContestName = contest.Name
		}
	}
	if challenge.Type != model.PodsChallengeType {
		for _, flag := range challenge.ChallengeFlags {
			result.Flags = append(result.Flags, view.ChallengeFlagView{
				ID:    flag.ID,
				Value: flag.Value,
			})
		}
	} else {
		result.DockerCompose = Template2Yaml(challenge.Template, challenge.ChallengeFlags)
	}

	file, ret := sourceFileRecord(tx, challenge)
	result.FileName = file.Filename
	if !ret.OK && ret.Msg != i18n.Model.NotFound {
		result.Unavailable = []string{"file"}
	}
	return result
}

// ListChallengeViews 返回跨比赛聚合的题目列表。
func ListChallengeViews(tx *gorm.DB, form dto.GetChallengesForm) ([]view.ChallengeView, int64, model.RetVal) {
	challenges, count, ret := db.InitChallengeRepo(tx).ListChallenges(
		form.Limit,
		form.Offset,
		challengeFilter(form.ContestID, form.Type, form.Category, form.Name, form.Description),
	)
	if !ret.OK {
		return nil, 0, ret
	}
	views := make([]view.ChallengeView, 0, len(challenges))
	for _, challenge := range challenges {
		views = append(views, BuildChallengeView(tx, challenge))
	}
	return views, count, model.SuccessRetVal()
}

func ListChallengeCategories(tx *gorm.DB, form dto.GetCategoriesForm) ([]string, model.RetVal) {
	return db.InitChallengeRepo(tx).ListCategories(challengeFilter(form.ContestID, form.Type, "", "", ""))
}

func GetChallengeWithFlags(tx *gorm.DB, challenge model.Challenge) (model.Challenge, model.RetVal) {
	return db.InitChallengeRepo(tx).GetByID(challenge.ID, db.GetOptions{
		Preloads: map[string]db.GetOptions{"ChallengeFlags": {}},
	})
}

func DeleteChallenge(tx *gorm.DB, challenge model.Challenge) model.RetVal {
	return db.WithTransactionDB(tx, func(tx *gorm.DB) model.RetVal {
		if ret := stopGeneratorResources(tx, db.GetOptions{Conditions: map[string]any{"challenge_id": challenge.ID}}); !ret.OK {
			return ret
		}
		return db.InitChallengeRepo(tx).Delete(challenge.RandID)
	})
}

// ListChallengesInContest 返回比赛对选手可见的题目。
func ListChallengesInContest(tx *gorm.DB, contest model.Contest, team model.Team, form dto.GetChallengesInContestForm) ([]view.ChallengeRuntimeView, int64, model.RetVal) {
	var (
		challenges []model.Challenge
		count      int64
		ret        model.RetVal
	)
	repo := db.InitChallengeRepo(tx)
	if form.Unsolved {
		ids, unsolvedCount, listRet := repo.ListUnsolvedID(team.ID, contest.ID, form.Category, form.Limit, form.Offset)
		if !listRet.OK {
			return nil, 0, listRet
		}
		count = unsolvedCount
		challenges, _, ret = repo.List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"id": ids},
			Preloads:   map[string]db.GetOptions{"ChallengeFlags": {}},
		})
	} else {
		options := db.GetOptions{
			Conditions: map[string]any{"contest_id": contest.ID, "hidden": false},
			Preloads:   map[string]db.GetOptions{"ChallengeFlags": {}},
		}
		if form.Category != "" {
			options.Conditions["category"] = form.Category
		}
		challenges, count, ret = repo.List(form.Limit, form.Offset, options)
	}
	if !ret.OK {
		return nil, 0, ret
	}
	views := make([]view.ChallengeRuntimeView, 0, len(challenges))
	for _, challenge := range challenges {
		item, ret := BuildChallengeRuntimeView(tx, team, challenge)
		if !ret.OK {
			return nil, count, ret
		}
		views = append(views, item)
	}
	return views, count, model.SuccessRetVal()
}

// ListAllChallengesInContest 返回比赛下的全部题目。
func ListAllChallengesInContest(tx *gorm.DB, contest model.Contest, form dto.GetAllChallengesInContestForm) ([]model.Challenge, int64, model.RetVal) {
	options := db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"ChallengeFlags": {}},
		Search:     make(map[string]string),
	}
	if form.Name != "" {
		options.Search["name"] = form.Name
	}
	if form.Category != "" {
		options.Conditions["category"] = form.Category
	}
	if form.Type != "" {
		options.Conditions["type"] = form.Type
	}
	return db.InitChallengeRepo(tx).List(form.Limit, form.Offset, options)
}

// ListChallengeCategoriesInContest 返回比赛下的题目分类。
func ListChallengeCategoriesInContest(tx *gorm.DB, contest model.Contest, form dto.GetCategoriesForm) ([]string, model.RetVal) {
	return db.InitChallengeRepo(tx).ListCategories(db.ChallengeFilter{ContestID: contest.ID, Type: form.Type})
}

func BuildChallengeRuntimeView(tx *gorm.DB, team model.Team, challenge model.Challenge) (view.ChallengeRuntimeView, model.RetVal) {
	attempts, ret := CountAttempts(tx, team, challenge)
	if !ret.OK {
		return view.ChallengeRuntimeView{}, ret
	}
	generated, ret := CheckIfGenerated(tx, team, challenge.ChallengeFlags)
	if !ret.OK {
		return view.ChallengeRuntimeView{}, ret
	}
	solved, ret := CheckIfSolved(tx, team, challenge.ChallengeFlags)
	if !ret.OK {
		return view.ChallengeRuntimeView{}, ret
	}
	return view.ChallengeRuntimeView{
		Challenge: challenge,
		Attempts:  attempts,
		Init:      generated,
		Solved:    solved,
		Remote:    GetVictimStatus(tx, team.ID, challenge.ID, challenge.Type),
		FileName:  buildSourceFileName(tx, challenge, team.ID),
	}, model.SuccessRetVal()
}

func GetChallengeStatus(tx *gorm.DB, team model.Team, challenge model.Challenge) (view.ChallengeStatusView, model.RetVal) {
	flags, _, ret := db.InitChallengeFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID},
	})
	if !ret.OK {
		return view.ChallengeStatusView{}, ret
	}
	challenge.ChallengeFlags = flags
	state, ret := BuildChallengeRuntimeView(tx, team, challenge)
	if !ret.OK {
		return view.ChallengeStatusView{}, ret
	}
	return view.ChallengeStatusView{
		Attempts: state.Attempts,
		Init:     state.Init,
		Solved:   state.Solved,
		Remote:   state.Remote,
		FileName: state.FileName,
	}, model.SuccessRetVal()
}

// CopyChallengesToContest 把已有题目复制到某场比赛，副本与源题目互不关联。
func CopyChallengesToContest(tx *gorm.DB, contest model.Contest, form dto.CopyChallengesToContestForm) ([]model.Challenge, []string, model.RetVal) {
	copied := make([]model.Challenge, 0)
	failed := make([]string, 0)
	batch := model.NewBatch(len(form.ChallengeIDs))
	repo := db.InitChallengeRepo(tx)
	for _, randID := range form.ChallengeIDs {
		if tx.Statement.Context.Err() != nil {
			return copied, failed, batch.Result(tx.Statement.Context)
		}
		source, ret := repo.GetByRandID(randID, db.GetOptions{
			Preloads: map[string]db.GetOptions{"ChallengeFlags": {}},
		})
		if !ret.OK {
			failed = append(failed, randID)
			batch.Fail(randID, "lookup", ret)
			if model.BatchDependencyFailed(ret) {
				return copied, failed, batch.Result(tx.Statement.Context)
			}
			continue
		}
		var created model.Challenge
		ret = db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
			challenge, ret := CopyChallenge(tx2, source, contest.ID)
			if !ret.OK {
				return ret
			}
			created, ret = db.InitChallengeRepo(tx2).GetByID(challenge.ID, db.GetOptions{
				Preloads: map[string]db.GetOptions{"ChallengeFlags": {}},
			})
			return ret
		})
		if !ret.OK {
			failed = append(failed, randID)
			batch.Fail(randID, "transaction", ret)
			if model.BatchDependencyFailed(ret) {
				return copied, failed, batch.Result(tx.Statement.Context)
			}
			continue
		}
		copied = append(copied, created)
		batch.Success(randID, "committed")
		warmChallengeImages(created)
		warmGeneratorPool(model.NewChallengeSource(created))
	}
	return copied, failed, batch.Result(tx.Statement.Context)
}

func ListChallengeFlagSolvers(tx *gorm.DB, flag model.ChallengeFlag) ([]view.FlagSolverView, model.RetVal) {
	rows, ret := db.InitSubmissionRepo(tx).ListFlagSolvers(flag.ID)
	if !ret.OK {
		return nil, ret
	}
	solvers := make([]view.FlagSolverView, 0, len(rows))
	for _, row := range rows {
		solvers = append(solvers, view.FlagSolverView{
			UserID:   row.UserID,
			UserName: row.UserName,
			TeamID:   row.TeamID,
			TeamName: row.TeamName,
			Score:    row.Score,
			SolvedAt: row.SolvedAt,
		})
	}
	return solvers, model.SuccessRetVal()
}
