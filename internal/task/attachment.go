package task

import (
	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"CBCTF/internal/websocket"
	wm "CBCTF/internal/websocket/model"
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const GenAttachmentTaskType = "tasks:attachment"

type GenAttachmentPayload struct {
	UserID    uint
	Challenge model.Challenge
	TeamID    uint
	Flags     []string
}

func EnqueueGenAttachmentTask(userID uint, challenge model.Challenge, team model.Team, teamFlags []model.TeamFlag) (*asynq.TaskInfo, error) {
	var flags []string
	for _, flag := range teamFlags {
		flags = append(flags, flag.Value)
	}
	payload, err := msgpack.Marshal(GenAttachmentPayload{userID, challenge, team.ID, flags})
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if ret := k8s.GenAttachment(ctx, payload.Challenge, payload.TeamID, payload.Flags); !ret.OK {
		websocket.Send(false, payload.UserID, wm.ErrorLevel, wm.GenerateAttachmentWSType, "Generate Attachment", "Failed")
	} else {
		websocket.Send(false, payload.UserID, wm.SuccessLevel, wm.GenerateAttachmentWSType, "Generate Attachment", "Done")
	}
	return nil
}
