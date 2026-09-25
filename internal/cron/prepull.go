package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"CBCTF/internal/task"
)

func warmChallengeImagesTask() model.RetVal {
	challenges, _, ret := db.InitChallengeRepo(db.CronDB).List(-1, -1)
	if !ret.OK {
		return ret
	}
	for _, challenge := range challenges {
		if err := task.EnqueuePrepullTask(service.ChallengeImages(challenge)); err != nil {
			return model.RetVal{Msg: i18n.Task.EnqueueError, Attr: map[string]any{"Error": err.Error()}}
		}
	}
	return model.SuccessRetVal()
}
