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
	Name      string
	ContestID uint
	Desc      string
	Captcha   string
	Avatar    model.AvatarURL
	Banned    bool
	Hidden    bool
	CaptainID uint
	Last      time.Time
}

func (c CreateTeamOptions) Convert2Model() model.Model {
	return model.Team{
		Name:      c.Name,
		ContestID: c.ContestID,
		Desc:      c.Desc,
		Captcha:   c.Captcha,
		Avatar:    c.Avatar,
		Banned:    c.Banned,
		Hidden:    c.Hidden,
		CaptainID: c.CaptainID,
		Last:      c.Last,
	}
}

type UpdateTeamOptions struct {
	Name      *string
	Desc      *string
	Captcha   *string
	Avatar    *model.AvatarURL
	Banned    *bool
	Hidden    *bool
	CaptainID *uint
	Score     *float64
	Rank      *int
	Last      *time.Time
	UserCount *int64
}

func (u UpdateTeamOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Desc != nil {
		options["desc"] = *u.Desc
	}
	if u.Captcha != nil {
		options["captcha"] = *u.Captcha
	}
	if u.Avatar != nil {
		options["avatar"] = *u.Avatar
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

func (t *TeamRepo) GetBy2ID(userID, contestID uint, optionsL ...GetOptions) (model.Team, model.RetVal) {
	options := GetOptions{}
	if len(optionsL) > 0 {
		options = optionsL[0]
	}
	if options.Conditions == nil {
		options.Conditions = make(map[string]any)
	}
	options.Conditions["contest_id"] = contestID
	user, ret := InitUserRepo(t.DB).GetByID(userID, GetOptions{
		Selects:  []string{"id"},
		Preloads: map[string]GetOptions{"Teams": options},
	})
	if !ret.OK {
		return model.Team{}, ret
	}
	if len(user.Teams) == 0 {
		return model.Team{}, model.RetVal{Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": model.Team{}.GetModelName()}}
	}
	return user.Teams[0], model.SuccessRetVal()
}

func (t *TeamRepo) Delete(idL ...uint) model.RetVal {
	teamL, _, ret := t.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id", "name", "contest_id"},
		Preloads: map[string]GetOptions{
			"Users":       {Selects: []string{"id"}},
			"Submissions": {Selects: []string{"id", "team_id"}},
			"TeamFlags":   {Selects: []string{"id", "team_id"}},
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
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": model.Team{}.GetModelName(), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
