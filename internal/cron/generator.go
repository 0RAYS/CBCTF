package cron

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"errors"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	corev1 "k8s.io/api/core/v1"
)

// prepareGenerator 预启动生成器 Pod, 同时关闭长时运行 Pod, 重置资源
func prepareGenerator(c *cron.Cron) {
	function := exec("PrepareGenerator", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		nodes, ret := k8s.ListSchedulableNodes(ctx)
		cancel()
		if !ret.OK {
			log.Logger.Warningf("Failed to count nodes")
			return errors.New(ret.Msg)
		}
		contests, _, ret := db.InitContestRepo(db.DB).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"hidden": false},
		})
		if !ret.OK {
			return errors.New(ret.Msg)
		}
		contestChallengeRepo := db.InitContestChallengeRepo(db.DB)
		for _, contest := range contests {
			if contest.IsOver() {
				continue
			}
			contestChallengeL, _, ret := contestChallengeRepo.List(-1, -1, db.GetOptions{
				Conditions: map[string]any{"contest_id": contest.ID, "type": model.DynamicChallengeType},
				Preloads:   map[string]db.GetOptions{"Challenge": {}},
			})
			if !ret.OK {
				continue
			}
			for _, contestChallenge := range contestChallengeL {
				challenge := contestChallenge.Challenge
				timeoutL := make([]*corev1.Pod, 0)
				k8s.GeneratorMapMutex.RLock()
				for _, generator := range k8s.GeneratorMap[challenge.ID] {
					if generator.Pod.Status.Phase != corev1.PodRunning || generator.Start.Add(6*time.Hour).Before(time.Now()) {
						timeoutL = append(timeoutL, generator.Pod)
					}
				}
				k8s.GeneratorMapMutex.RUnlock()
				for _, generator := range timeoutL {
					ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
					k8s.StopGenerator(ctx, challenge, generator)
					cancel()
				}
				k8s.GeneratorMapMutex.RLock()
				length := len(k8s.GeneratorMap[challenge.ID])
				k8s.GeneratorMapMutex.RUnlock()
				var wg sync.WaitGroup
				for i := 0; i < len(nodes)*config.Env.K8S.GeneratorWorker-length; i++ {
					wg.Go(func() {
						gCtx, gCancel := context.WithTimeout(context.Background(), time.Minute)
						k8s.StartGenerator(gCtx, challenge)
						gCancel()
					})
				}
				wg.Wait()
			}
		}
		return nil
	})
	function()
	c.Schedule(cron.Every(2*time.Minute), cron.FuncJob(function))
}

// stopUnCtrlGenerator 关闭不受控的 (k8s.GeneratorMap) 以 `gen` 为命名前缀的 pod
func stopUnCtrlGenerator(c *cron.Cron) {
	function := exec("StopUnCtrlGenerator", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		pods, ret := k8s.GetPodList(ctx)
		cancel()
		if !ret.OK {
			log.Logger.Warningf("Failed to get generators %v", ret)
			return errors.New(ret.Msg)
		}
		names := make([]string, 0)
		k8s.GeneratorMapMutex.RLock()
		for _, v := range k8s.GeneratorMap {
			for _, generator := range v {
				names = append(names, generator.Pod.Name)
			}
		}
		k8s.GeneratorMapMutex.RUnlock()
		for _, pod := range pods.Items {
			if strings.HasPrefix(pod.Name, "gen") && !slices.Contains(names, pod.Name) {
				ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
				k8s.DeletePod(ctx, pod.Name)
				cancel()
			}
		}
		return nil
	})
	function()
	c.Schedule(cron.Every(time.Hour), cron.FuncJob(function))
}
