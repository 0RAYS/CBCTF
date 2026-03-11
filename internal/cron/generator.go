package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"slices"
	"time"

	"github.com/robfig/cron/v3"
)

// stopUnCtrlGenerator 关闭不受控的以 `gen` 为命名前缀的 pod
func stopUnCtrlGenerator(c *cron.Cron) {
	function := exec("StopUnCtrlGenerator", func() model.RetVal {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		pods, ret := k8s.GetPodList(ctx, map[string]string{k8s.GeneratorPodTag: k8s.GeneratorPodTag})
		cancel()
		if !ret.OK {
			log.Logger.Warningf("Failed to get generators %v", ret)
			return ret
		}
		generators, _, ret := db.InitGeneratorRepo(db.DB).List(-1, -1)
		if !ret.OK {
			log.Logger.Warningf("Failed to get generators %v", ret)
			return ret
		}
		for _, pod := range pods.Items {
			if !slices.ContainsFunc(generators, func(generator model.Generator) bool {
				return generator.Name == pod.Name
			}) {
				ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
				k8s.DeletePod(ctx, pod.Name)
				cancel()
			}
		}
		return model.SuccessRetVal()
	})
	function()
	c.Schedule(cron.Every(time.Hour), cron.FuncJob(function))
}
