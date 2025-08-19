package cron

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"slices"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	corev1 "k8s.io/api/core/v1"
)

// prepareGenerator 预启动生成器 Pod, 同时关闭长时运行 Pod, 重置资源
func prepareGenerator(c *cron.Cron) {
	function := exec("PrepareGenerator", func() {
		contests, _, ok, _ := db.InitContestRepo(db.DB).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"hidden": false},
			Selects:    []string{"id", "start", "duration"},
		})
		if !ok {
			return
		}
		contestChallengeRepo := db.InitContestChallengeRepo(db.DB)
		for _, contest := range contests {
			if contest.IsOver() {
				continue
			}
			contestChallengeL, _, ok, _ := contestChallengeRepo.List(-1, -1, db.GetOptions{
				Conditions: map[string]any{"contest_id": contest.ID, "type": model.DynamicChallengeType},
				Selects:    []string{"id", "name", "challenge_id"},
				Preloads: map[string]db.GetOptions{
					"Challenge": {Selects: []string{"id", "rand_id", "generator_image"}},
				},
			})
			if !ok {
				continue
			}
			for _, contestChallenge := range contestChallengeL {
				timeoutL := make([]*corev1.Pod, 0)
				k8s.GeneratorMapMutex.RLock()
				for _, generator := range k8s.GeneratorMap[contestChallenge.ID] {
					// TODO 此处需要注意平台所在时区与K8S节点需要相同
					if generator.Status.Phase != corev1.PodRunning || generator.CreationTimestamp.Add(time.Hour).Before(time.Now()) {
						timeoutL = append(timeoutL, generator)
					}
				}
				k8s.GeneratorMapMutex.RUnlock()
				for _, generator := range timeoutL {
					k8s.StopGenerator(contestChallenge, generator)
				}
				k8s.GeneratorMapMutex.RLock()
				length := len(k8s.GeneratorMap[contestChallenge.ID])
				k8s.GeneratorMapMutex.RUnlock()
				for i := 0; i < len(config.Env.K8S.Nodes)*config.Env.K8S.GeneratorWorker-length; i++ {
					go k8s.StartGenerator(contestChallenge)
				}
			}
		}
	})
	function()
	c.Schedule(cron.Every(2*time.Minute), cron.FuncJob(function))
}

// stopUnCtrlGenerator 关闭不受控的 (k8s.GeneratorMap) 以 `gen` 为命名前缀的 pod
func stopUnCtrlGenerator(c *cron.Cron) {
	function := exec("StopUnCtrlGenerator", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		pods, ok, msg := k8s.GetPodList(ctx)
		cancel()
		if !ok {
			log.Logger.Warningf("Failed to get generators %s", msg)
			return
		}
		generators := make(map[string]*corev1.Pod)
		names := make([]string, 0)
		k8s.GeneratorMapMutex.RLock()
		for _, v := range k8s.GeneratorMap {
			for _, generator := range v {
				generators[generator.Name] = generator
				names = append(names, generator.Name)
			}
		}
		k8s.GeneratorMapMutex.RUnlock()
		for _, pod := range pods.Items {
			if strings.HasPrefix(pod.Name, "gen") && !slices.Contains(names, pod.Name) {
				ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
				_, _ = k8s.DeletePod(ctx, pod.Name)
				cancel()
			}
		}
	})
	function()
	c.Schedule(cron.Every(time.Hour), cron.FuncJob(function))
}
