package task

import (
	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const GenAttachmentTaskType = "tasks:attachment"

type GenAttachmentPayload struct {
	UserID           uint
	ContestChallenge model.ContestChallenge
	Team             model.Team
	TeamFlags        []model.TeamFlag
}

func EnqueueGenAttachmentTask(userID uint, contestChallenge model.ContestChallenge, team model.Team, teamFlags []model.TeamFlag) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(GenAttachmentPayload{userID, contestChallenge, team, teamFlags})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(GenAttachmentTaskType, payload)
	return client.Enqueue(task, asynq.Timeout(2*time.Minute))
}

func HandleGenAttachmentTask(_ context.Context, t *asynq.Task) error {
	var payload GenAttachmentPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	if ok, _ := k8s.GenAttachment(payload.ContestChallenge, payload.Team, payload.TeamFlags); !ok {
		//websocket.Send(false, payload.UserID, wm.ErrorLevel, wm.GenerateAttachmentWSType, "Generate Attachment", "Failed")
	} else {
		//websocket.Send(false, payload.UserID, wm.SuccessLevel, wm.GenerateAttachmentWSType, "Generate Attachment", "Done")
	}
	return nil
}
