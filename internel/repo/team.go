package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type TeamRepo struct {
	BasicRepo[model.Team]
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
}

func (u UpdateTeamOptions) Convert2Map() map[string]any {
	data := make(map[string]any)
	if u.Name != nil {
		data["name"] = *u.Name
	}
	if u.Desc != nil {
		data["desc"] = *u.Desc
	}
	if u.Captcha != nil {
		data["captcha"] = *u.Captcha
	}
	if u.Avatar != nil {
		data["avatar"] = *u.Avatar
	}
	if u.Banned != nil {
		data["banned"] = *u.Banned
	}
	if u.Hidden != nil {
		data["hidden"] = *u.Hidden
	}
	if u.CaptainID != nil {
		data["captain_id"] = *u.CaptainID
	}
	if u.Score != nil {
		data["score"] = *u.Score
	}
	if u.Rank != nil {
		data["rank"] = *u.Rank
	}
	if u.Last != nil {
		data["last"] = *u.Last
	}
	return data
}

func InitTeamRepo(tx *gorm.DB) *TeamRepo {
	return &TeamRepo{
		BasicRepo: BasicRepo[model.Team]{
			DB: tx,
		},
	}
}

func (t *TeamRepo) IsUniqueName(contestID uint, name string) bool {
	_, ok, _ := t.Get(GetOptions{
		Conditions: map[string]any{
			"contest_id": contestID,
			"name":       name,
		},
		Selects: []string{"id"},
	})
	return !ok
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

func (t *TeamRepo) GetByName(contestID uint, name string, optionsL ...GetOptions) (model.Team, bool, string) {
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

func (t *TeamRepo) GetBy2ID(userID, contestID uint, optionsL ...GetOptions) (model.Team, bool, string) {
	options := GetOptions{}
	if len(optionsL) > 0 {
		options = optionsL[0]
	}
	if options.Conditions == nil {
		options.Conditions = make(map[string]any)
	}
	options.Conditions["contest_id"] = contestID
	user, ok, msg := InitUserRepo(t.DB).GetByID(userID, GetOptions{
		Selects: []string{"id"},
		Preloads: map[string]GetOptions{
			"Teams": options,
		},
	})
	if !ok {
		return model.Team{}, false, msg
	}
	if len(user.Teams) == 0 {
		return model.Team{}, false, i18n.TeamNotFound
	}
	return *user.Teams[0], true, i18n.Success
}

func (t *TeamRepo) Delete(idL ...uint) (bool, string) {
	teamL, _, ok, msg := t.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id", "name", "contest_id"},
		Preloads: map[string]GetOptions{
			"Users":       {Selects: []string{"id"}},
			"Submissions": {Selects: []string{"id"}},
			"TeamFlags":   {Selects: []string{"id"}},
		},
	})
	if !ok && msg != i18n.TeamNotFound {
		return false, msg
	}
	submissionIDL, teamFlagIDL := make([]uint, 0), make([]uint, 0)
	for _, team := range teamL {
		deletedName := fmt.Sprintf("%s_deleted_%s", team.Name, utils.RandStr(6))
		if ok, msg = t.Update(team.ID, UpdateTeamOptions{
			Name: &deletedName,
		}); !ok {
			return false, msg
		}
		for _, user := range team.Users {
			if ok, msg = DeleteUserFromContest(t.DB, user.ID, team.ContestID); !ok {
				return false, msg
			}
			if ok, msg = DeleteUserFromTeam(t.DB, user.ID, team.ID); !ok {
				return false, msg
			}
		}
		for _, submission := range team.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
		for _, teamFlag := range team.TeamFlags {
			teamFlagIDL = append(teamFlagIDL, teamFlag.ID)
		}
	}
	if ok, msg = InitSubmissionRepo(t.DB).Delete(submissionIDL...); !ok {
		return false, msg
	}
	if ok, msg = InitTeamFlagRepo(t.DB).Delete(teamFlagIDL...); !ok {
		return false, msg
	}
	if res := t.DB.Model(&model.Team{}).Where("id IN ?", idL).Delete(&model.Team{}); res.Error != nil {
		log.Logger.Errorf("Failed to delete Team: %s", res.Error)
		return false, i18n.DeleteTeamError
	}
	return true, i18n.Success
}
