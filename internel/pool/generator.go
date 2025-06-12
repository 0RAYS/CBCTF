package pool

import (
	"CBCTF/internel/k8s"
	"CBCTF/internel/model"
	"os"
	"os/signal"
	"time"
)

type GenTask struct {
	Team             model.Team
	ContestChallenge model.ContestChallenge
	TeamFlagL        []model.TeamFlag
}

var (
	GenAttachmentPool = make(chan GenTask, 1000)
)

func PrepareAttachment() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	for {
		select {
		case <-stop:
			return
		case task := <-GenAttachmentPool:
			k8s.GenerateAttachment(task.ContestChallenge, task.Team, task.TeamFlagL)
		case <-time.After(1 * time.Second):

		}
	}
}
