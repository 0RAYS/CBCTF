package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/k8s"
	"context"
	corev1 "k8s.io/api/core/v1"
	"time"
)

func GetNodeImageList() (map[string][]string, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return k8s.GetNodeImageList(ctx)
}

func WarmUpContestChallengeImage(form f.WarmUpImageForm) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if form.PullPolicy == string(corev1.PullNever) {
		return true, i18n.Success
	}
	nodes, ok, msg := k8s.ListSchedulableNodes(ctx)
	if !ok {
		return false, msg
	}
	for _, node := range nodes {
		if _, ok, msg = k8s.CreateJob(ctx, k8s.CreateJobOptions{
			Images:     form.Images,
			PullPolicy: form.PullPolicy,
			NodeSelector: map[string]string{
				"kubernetes.io/hostname": node.Name,
			},
		}); !ok {
			return false, msg
		}
	}
	return true, i18n.Success
}
