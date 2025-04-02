package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
	"time"
)

func CreateContest(tx *gorm.DB, form f.CreateContestForm) (model.Contest, bool, string) {
	repo := db.InitContestRepo(tx)
	return repo.Create(db.CreateContestOptions{
		Name:      form.Name,
		Desc:      form.Desc,
		Captcha:   form.Captcha,
		Avatar:    "",
		Prefix:    form.Prefix,
		Size:      form.Size,
		Start:     form.Start,
		Duration:  time.Duration(form.Duration) * time.Second,
		Blood:     form.Blood,
		Hidden:    form.Hidden,
		Rules:     form.Rules,
		Prizes:    form.Prizes,
		Timelines: form.Timelines,
	})
}

func UpdateContest(tx *gorm.DB, contest model.Contest, form f.UpdateContestForm) (bool, string) {
	repo := db.InitContestRepo(tx)
	if form.Name != nil && *form.Name != contest.Name {
		if !repo.IsUniqueName(*form.Name) {
			return false, "DuplicateContestName"
		}
	}
	duration := time.Duration(*form.Duration) * time.Second
	return repo.Update(contest.ID, db.UpdateContestOptions{
		Name:      form.Name,
		Desc:      form.Desc,
		Captcha:   form.Captcha,
		Prefix:    form.Prefix,
		Size:      form.Size,
		Start:     form.Start,
		Duration:  &duration,
		Blood:     form.Blood,
		Hidden:    form.Hidden,
		Rules:     form.Rules,
		Prizes:    form.Prizes,
		Timelines: form.Timelines,
	})
}
