package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"slices"
	"time"
)

// stopUnCtrlGenerator 关闭不受控的 model.Generator
func stopUnCtrlGeneratorTask() model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	pods, ret := k8s.GetPodList(ctx, map[string]string{k8s.GeneratorPodTag: k8s.GeneratorPodTag})
	cancel()
	if !ret.OK {
		return ret
	}
	generators, _, ret := db.InitGeneratorRepo(db.DB).List(-1, -1)
	if !ret.OK {
		return ret
	}
	for _, pod := range pods.Items {
		if !slices.ContainsFunc(generators, func(generator model.Generator) bool {
			return generator.Name == pod.Name
		}) {
			ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
			if ret = k8s.DeletePod(ctx, pod.Name); ret.OK {
				log.Logger.Infof("Deleted uncontrolled generator pod: pod=%s", pod.Name)
			} else {
				log.Logger.Warningf("Failed to delete uncontrolled generator pod: pod=%s reason=%s", pod.Name, ret.Msg)
			}
			cancel()
		}
	}
	return model.SuccessRetVal()
}
