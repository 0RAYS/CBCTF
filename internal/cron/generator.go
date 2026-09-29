package cron

import (
	"context"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

// stopUnCtrlGenerator 关闭不受控的 model.Generator
func stopUnCtrlGeneratorTask() model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	pods, ret := k8s.ListPods(ctx, map[string]string{k8s.RoleLabel: k8s.GeneratorPodTag})
	if !ret.OK {
		return ret
	}
	if pods == nil {
		return model.SuccessRetVal()
	}
	generators, _, ret := db.InitGeneratorRepo(db.CronDB.WithContext(ctx)).List(-1, -1)
	if !ret.OK {
		return ret
	}
	known := make(map[string]bool, len(generators))
	for _, generator := range generators {
		known[generator.Name] = true
	}
	batch := model.NewBatch(len(pods.Items))
	for _, pod := range pods.Items {
		if ctx.Err() != nil {
			return batch.Result(ctx)
		}
		if !known[pod.Name] {
			itemCtx, itemCancel := context.WithTimeout(ctx, time.Minute)
			if ret = k8s.DeletePod(itemCtx, pod.Name, pod.UID); ret.OK {
				batch.Success(pod.Name, "deleted")
				log.Logger.Infof("Deleted uncontrolled generator pod: pod=%s", pod.Name)
			} else {
				batch.Fail(pod.Name, "delete", ret)
				log.Logger.Warningf("Failed to delete uncontrolled generator pod: pod=%s reason=%s", pod.Name, ret.Msg)
			}
			itemCancel()
		} else {
			batch.Skip(pod.Name, "database_owned")
		}
	}
	return batch.Result(ctx)
}
