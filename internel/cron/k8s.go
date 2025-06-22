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
					if k == "victim" || k == "generator" {
						if _, ok, _ = k8s.GetPod(ctx, v); !ok {
							k8s.DeleteConfigMapListByPodName(ctx, k, v)
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
					if k == "victim" || k == "generator" {
						if _, ok, _ = k8s.GetPod(ctx, v); !ok {
							k8s.DeleteServiceListByPodName(ctx, k, v)
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
					if k == "victim" || k == "generator" {
						if _, ok, _ = k8s.GetPod(ctx, v); !ok {
							k8s.DeleteNetworkPolicyListByPodName(ctx, k, v)
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
