package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/task"
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

func ChallengeImages(challenge model.Challenge) []string {
	sidecars := []string{}
	if challenge.Type == model.DynamicChallengeType {
		sidecars = append(sidecars, config.Env.K8S.WorkerImage)
	}
	if challenge.Type == model.PodsChallengeType {
		if config.Env.K8S.CaptureEnabled {
			sidecars = append(sidecars, config.Env.K8S.CaptureImage)
		}
		if config.Env.K8S.Frp.On {
			sidecars = append(sidecars, config.Env.K8S.Frp.FrpcImage, config.Env.K8S.Frp.NginxImage)
		}
	}
	return k8s.ChallengeImages(challenge, sidecars...)
}

func warmChallengeImages(challenge model.Challenge) {
	if err := task.EnqueuePrepullTask(ChallengeImages(challenge)); err != nil {
		log.Logger.Warningf("Failed to enqueue image warmup: challenge_id=%d error=%v", challenge.ID, err)
	}
}

func ListNodeImages() (map[string][]string, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return k8s.ListNodeImages(ctx)
}

func PullContestChallengeImage(form dto.PullImageForm) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if form.PullPolicy == string(corev1.PullNever) {
		return model.SuccessRetVal()
	}
	nodes, ret := k8s.ListSchedulableNodes(ctx)
	if !ret.OK {
		return ret
	}

	nodeMap := make(map[string]*corev1.Node, len(nodes))
	for _, node := range nodes {
		nodeMap[node.Name] = node
	}

	targetImages := make(map[string][]string)
	seen := make(map[string]map[string]struct{})
	for _, target := range form.Targets {
		nodeName := strings.TrimSpace(target.Node)
		imageName := strings.TrimSpace(target.Image)
		if nodeName == "" || imageName == "" {
			continue
		}
		node, ok := nodeMap[nodeName]
		if !ok || node == nil {
			return model.RetVal{Msg: i18n.Response.BadRequest, Attr: map[string]any{"Error": fmt.Sprintf("Unknown node: %s", nodeName)}}
		}

		if _, ok = seen[nodeName]; !ok {
			seen[nodeName] = make(map[string]struct{})
		}
		if _, ok = seen[nodeName][imageName]; ok {
			continue
		}
		seen[nodeName][imageName] = struct{}{}
		targetImages[nodeName] = append(targetImages[nodeName], imageName)
	}

	for nodeName, images := range targetImages {
		if err := task.EnqueuePrepullTargets(images, []string{nodeName}, form.PullPolicy); err != nil {
			return model.RetVal{Msg: i18n.Task.EnqueueError, Attr: map[string]any{"Error": err.Error()}}
		}
	}
	return model.SuccessRetVal()
}
