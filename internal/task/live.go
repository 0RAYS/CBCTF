package task

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"slices"
	"sort"
	"time"

	"github.com/hibiken/asynq"
)

var inspector *asynq.Inspector

func ListLiveTasks(status, queue string, limit, offset int) ([]*asynq.TaskInfo, int64, []string, []string, model.RetVal) {
	if inspector == nil {
		return nil, 0, nil, nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": "task inspector unavailable"}}
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
		if end > len(allTasks) {
			end = len(allTasks)
		}
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
	switch state {
	case "pending":
		return listPendingTasks(queue, info.Pending)
	case "active":
		return listActiveTasks(queue, info.Active)
	case "scheduled":
		return listScheduledTasks(queue, info.Scheduled)
	case "retry":
		return listRetryTasks(queue, info.Retry)
	case "archived":
		return listArchivedTasks(queue, info.Archived)
	case "completed":
		return listCompletedTasks(queue, info.Completed)
	default:
		return listActiveTasks(queue, info.Active)
	}
}

func listPendingTasks(queue string, size int) ([]*asynq.TaskInfo, error) {
	if size <= 0 {
		return []*asynq.TaskInfo{}, nil
	}
	return inspector.ListPendingTasks(queue, asynq.Page(1), asynq.PageSize(size))
}

func listActiveTasks(queue string, size int) ([]*asynq.TaskInfo, error) {
	if size <= 0 {
		return []*asynq.TaskInfo{}, nil
	}
	return inspector.ListActiveTasks(queue, asynq.Page(1), asynq.PageSize(size))
}

func listScheduledTasks(queue string, size int) ([]*asynq.TaskInfo, error) {
	if size <= 0 {
		return []*asynq.TaskInfo{}, nil
	}
	return inspector.ListScheduledTasks(queue, asynq.Page(1), asynq.PageSize(size))
}

func listRetryTasks(queue string, size int) ([]*asynq.TaskInfo, error) {
	if size <= 0 {
		return []*asynq.TaskInfo{}, nil
	}
	return inspector.ListRetryTasks(queue, asynq.Page(1), asynq.PageSize(size))
}

func listArchivedTasks(queue string, size int) ([]*asynq.TaskInfo, error) {
	if size <= 0 {
		return []*asynq.TaskInfo{}, nil
	}
	return inspector.ListArchivedTasks(queue, asynq.Page(1), asynq.PageSize(size))
}

func listCompletedTasks(queue string, size int) ([]*asynq.TaskInfo, error) {
	if size <= 0 {
		return []*asynq.TaskInfo{}, nil
	}
	return inspector.ListCompletedTasks(queue, asynq.Page(1), asynq.PageSize(size))
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
