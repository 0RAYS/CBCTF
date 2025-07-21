package cron

import (
	"CBCTF/internel/k8s"
	"context"
	"github.com/robfig/cron/v3"
	"time"
)

func ClearUnCtrlResource(c *cron.Cron) {
	function := exec("ClearUnCtrlResource", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		if configmaps, ok, _ := k8s.GetConfigMapList(ctx); ok {
			for _, cm := range configmaps.Items {
				for k, v := range cm.Labels {
					if k == k8s.VictimPodTag || k == k8s.GeneratorPodTag {
						if _, ok, _ = k8s.GetPod(ctx, v); !ok {
							k8s.DeleteConfigMapList(ctx, map[string]string{k: v})
						}
					}
				}
			}
		}
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
		if services, ok, _ := k8s.GetServiceList(ctx); ok {
			for _, cm := range services.Items {
				for k, v := range cm.Labels {
					if k == k8s.VictimPodTag || k == k8s.GeneratorPodTag {
						if _, ok, _ = k8s.GetPod(ctx, v); !ok {
							k8s.DeleteServiceList(ctx, map[string]string{k: v})
						}
					}
				}
			}
		}
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
		if policies, ok, _ := k8s.GetNetworkPolicyList(ctx); ok {
			for _, cm := range policies.Items {
				for k, v := range cm.Labels {
					if k == k8s.VictimPodTag || k == k8s.GeneratorPodTag {
						if _, ok, _ = k8s.GetPod(ctx, v); !ok {
							k8s.DeleteNetworkPolicyList(ctx, map[string]string{k: v})
						}
					}
				}
			}
		}
		cancel()
	})
	function()
	c.Schedule(cron.Every(time.Hour), cron.FuncJob(function))
}
