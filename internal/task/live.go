package task

import (
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/hibiken/asynq"

	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
)

var inspector *asynq.Inspector

func ListLiveTasks(status, queue, taskID string, limit, offset int) ([]*asynq.TaskInfo, int64, []string, []string, model.RetVal) {
	if inspector == nil {
		return nil, 0, nil, nil, model.RetVal{Msg: i18n.Task.InspectorUnavailable}
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	queues, err := inspector.Queues()
	if err != nil {
		return nil, 0, nil, nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	sort.Strings(queues)
	if queue != "" {
		if !slices.Contains(queues, queue) {
			return []*asynq.TaskInfo{}, 0, queues, []string{}, model.SuccessRetVal()
		}
		queues = []string{queue}
	}
	if len(queues) == 0 {
		return []*asynq.TaskInfo{}, 0, []string{}, []string{}, model.SuccessRetVal()
	}

	typeSet := make(map[string]struct{})
	allTasks := make([]*asynq.TaskInfo, 0)
	for _, q := range queues {
		tasks, err := listAllTaskState(status, q)
		if err != nil {
			return nil, 0, nil, nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
		}
		for _, item := range tasks {
			if taskID != "" && !strings.Contains(item.ID, taskID) {
				continue
			}
			allTasks = append(allTasks, item)
			typeSet[item.Type] = struct{}{}
		}
	}

	sort.SliceStable(allTasks, func(i, j int) bool {
		ti := liveTaskSortTime(allTasks[i], status)
		tj := liveTaskSortTime(allTasks[j], status)
		if ti.Equal(tj) {
			if allTasks[i].Queue == allTasks[j].Queue {
				return allTasks[i].ID > allTasks[j].ID
			}
			return allTasks[i].Queue < allTasks[j].Queue
		}
		if status == "scheduled" {
			return ti.Before(tj)
		}
		return ti.After(tj)
	})

	total := int64(len(allTasks))
	if offset >= len(allTasks) {
		allTasks = []*asynq.TaskInfo{}
	} else {
		end := offset + limit
		end = min(end, len(allTasks))
		allTasks = allTasks[offset:end]
	}

	types := make([]string, 0, len(typeSet))
	for key := range typeSet {
		types = append(types, key)
	}
	sort.Strings(types)
	return allTasks, total, queues, types, model.SuccessRetVal()
}

func listAllTaskState(state string, queue string) ([]*asynq.TaskInfo, error) {
	info, err := inspector.GetQueueInfo(queue)
	if err != nil {
		return nil, err
	}
	lists := []struct {
		state string
		size  int
		fn    func(string, ...asynq.ListOption) ([]*asynq.TaskInfo, error)
	}{
		{"active", info.Active, inspector.ListActiveTasks},
		{"pending", info.Pending, inspector.ListPendingTasks},
		{"scheduled", info.Scheduled, inspector.ListScheduledTasks},
		{"retry", info.Retry, inspector.ListRetryTasks},
		{"archived", info.Archived, inspector.ListArchivedTasks},
		{"completed", info.Completed, inspector.ListCompletedTasks},
	}
	for _, list := range lists {
		if state == list.state {
			if list.size <= 0 {
				return []*asynq.TaskInfo{}, nil
			}
			return list.fn(queue, asynq.Page(1), asynq.PageSize(list.size))
		}
	}
	allTasks := make([]*asynq.TaskInfo, 0, info.Pending+info.Active+info.Scheduled+info.Retry+info.Archived+info.Completed)
	for _, list := range lists {
		if list.size <= 0 {
			continue
		}
		tasks, err := list.fn(queue, asynq.Page(1), asynq.PageSize(list.size))
		if err != nil {
			return nil, err
		}
		allTasks = append(allTasks, tasks...)
	}
	return allTasks, nil
}

func liveTaskSortTime(task *asynq.TaskInfo, status string) time.Time {
	switch status {
	case "scheduled", "retry":
		if !task.NextProcessAt.IsZero() {
			return task.NextProcessAt
		}
	case "archived":
		if !task.LastFailedAt.IsZero() {
			return task.LastFailedAt
		}
	case "completed":
		if !task.CompletedAt.IsZero() {
			return task.CompletedAt
		}
	}
	if !task.LastFailedAt.IsZero() {
		return task.LastFailedAt
	}
	if !task.CompletedAt.IsZero() {
		return task.CompletedAt
	}
	if !task.NextProcessAt.IsZero() {
		return task.NextProcessAt
	}
	return time.Time{}
}
