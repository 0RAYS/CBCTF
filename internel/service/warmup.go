package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/k8s"
	"context"
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
	_, ok, msg := k8s.CreateDaemonSet(ctx, k8s.CreateDaemonSetOptions{
		Images:     form.Images,
		PullPolicy: form.PullPolicy,
	})
	return ok, msg
}
