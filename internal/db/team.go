package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TeamRepo struct {
	BaseRepo[model.Team]
}

type CreateTeamOptions struct {
	Name        string
	ContestID   uint
	Description string
	Captcha     string
	Picture     model.FileURL
	Banned      bool
	Hidden      bool
	CaptainID   uint
	Last        time.Time
}

func (c CreateTeamOptions) Convert2Model() model.Model {
	return model.Team{
		Name:        c.Name,
		ContestID:   c.ContestID,
		Description: c.Description,
		Captcha:     c.Captcha,
		Picture:     c.Picture,
		Banned:      c.Banned,
		Hidden:      c.Hidden,
		CaptainID:   c.CaptainID,
		Last:        c.Last,
	}
}

type UpdateTeamOptions struct {
	Name        *string
	Description *string
	Captcha     *string
	Picture     *model.FileURL
	Banned      *bool
	Hidden      *bool
	CaptainID   *uint
	Score       *float64
	Rank        *int
	Last        *time.Time
	UserCount   *int64
}

func (u UpdateTeamOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Description != nil {
		options["description"] = *u.Description
	}
	if u.Captcha != nil {
		options["captcha"] = *u.Captcha
	}
	if u.Picture != nil {
		options["picture"] = *u.Picture
	}
	if u.Banned != nil {
		options["banned"] = *u.Banned
	}
	if u.Hidden != nil {
		options["hidden"] = *u.Hidden
	}
	if u.CaptainID != nil {
		options["captain_id"] = *u.CaptainID
	}
	if u.Score != nil {
		options["score"] = *u.Score
	}
	if u.Rank != nil {
		options["rank"] = *u.Rank
	}
	if u.Last != nil {
		options["last"] = *u.Last
	}
	if u.UserCount != nil {
		options["user_count"] = *u.UserCount
	}
	return options
}

func InitTeamRepo(tx *gorm.DB) *TeamRepo {
	return &TeamRepo{
		BaseRepo: BaseRepo[model.Team]{
			DB: tx,
		},
	}
}

func (t *TeamRepo) IsInTeam(teamID uint, userID uint) bool {
	res := t.DB.Model(&model.UserTeam{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).Limit(1).Find(&model.UserTeam{})
	return res.RowsAffected == 1
}

func (t *TeamRepo) IsInContest(contestID uint, userID uint) bool {
	res := t.DB.Model(&model.UserContest{}).
		Where("contest_id = ? AND user_id = ?", contestID, userID).Limit(1).Find(&model.UserContest{})
	return res.RowsAffected == 1
}

func (t *TeamRepo) GetByName(contestID uint, name string, optionsL ...GetOptions) (model.Team, model.RetVal) {
	var options GetOptions
	if len(optionsL) > 0 {
		options = optionsL[0]
	}
	if options.Conditions == nil {
		options.Conditions = make(map[string]any)
	}
	options.Conditions["contest_id"] = contestID
	options.Conditions["name"] = name
	return t.Get(options)
}

func (t *TeamRepo) GetBy2ID(userID, contestID uint) (model.Team, model.RetVal) {
	var team model.Team
	res := t.DB.Raw(`
		SELECT teams.*
		FROM teams
		INNER JOIN user_teams ON user_teams.team_id = teams.id
		WHERE user_teams.user_id = ? AND teams.contest_id = ?
		LIMIT 1
	`, userID, contestID).Scan(&team)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Team: %s", res.Error)
		return model.Team{}, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": team.ModelName(), "Error": res.Error.Error()}}
	}
	return team, model.SuccessRetVal()
}

func (t *TeamRepo) Delete(idL ...uint) model.RetVal {
	teamL, _, ret := t.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Preloads: map[string]GetOptions{
			"Users":       {},
			"Submissions": {},
			"TeamFlags":   {},
		},
	})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	submissionIDL, teamFlagIDL := make([]uint, 0), make([]uint, 0)
	for _, team := range teamL {
		deletedName := fmt.Sprintf("%s_deleted_%s", team.Name, utils.RandStr(6))
		if ret = t.Update(team.ID, UpdateTeamOptions{
			Name: &deletedName,
		}); !ret.OK {
			return ret
		}
		for _, user := range team.Users {
			if ret = DeleteUserFromContest(t.DB, user, model.Contest{BaseModel: model.BaseModel{ID: team.ContestID}}); !ret.OK {
				return ret
			}
			if ret = DeleteUserFromTeam(t.DB, user, team); !ret.OK {
				return ret
			}
		}
		for _, submission := range team.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
		for _, teamFlag := range team.TeamFlags {
			teamFlagIDL = append(teamFlagIDL, teamFlag.ID)
		}
	}
	if ret = InitSubmissionRepo(t.DB).Delete(submissionIDL...); !ret.OK {
		return ret
	}
	if ret = InitTeamFlagRepo(t.DB).Delete(teamFlagIDL...); !ret.OK {
		return ret
	}
	if res := t.DB.Model(&model.Team{}).Where("id IN ?", idL).Delete(&model.Team{}); res.Error != nil {
		log.Logger.Errorf("Failed to delete Team: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": model.Team{}.ModelName(), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
