package service

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"

	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/task"
)

func challengeSidecarImages(challengeType model.ChallengeType) []string {
	var sidecars []string
	if challengeType == model.DynamicChallengeType {
		sidecars = append(sidecars, config.Env.K8S.WorkerImage)
	}
	if challengeType == model.PodsChallengeType {
		if config.Env.K8S.CaptureEnabled {
			sidecars = append(sidecars, config.Env.K8S.CaptureImage)
		}
		if config.Env.K8S.Frp.On {
			sidecars = append(sidecars, config.Env.K8S.Frp.FrpcImage, config.Env.K8S.Frp.NginxImage)
		}
	}
	return sidecars
}

func ChallengeImages(challenge model.Challenge) []string {
	return k8s.ChallengeImages(challenge, challengeSidecarImages(challenge.Type)...)
}

// ListChallengeImages covers the whole library, including challenges not yet
// referenced by a contest, plus the platform's enabled runtime dependencies.
func ListChallengeImages(tx *gorm.DB) ([]string, model.RetVal) {
	challenges, ret := db.InitChallengeRepo(tx.Select("type", "generator_image", "template")).FindAll()
	if !ret.OK {
		return nil, ret
	}
	sidecars := append(challengeSidecarImages(model.DynamicChallengeType), challengeSidecarImages(model.PodsChallengeType)...)
	images := k8s.ChallengeImages(model.Challenge{}, sidecars...)
	for _, challenge := range challenges {
		images = append(images, ChallengeImages(challenge)...)
	}
	slices.Sort(images)
	return slices.Compact(images), model.SuccessRetVal()
}

func warmChallengeImages(challenge model.Challenge) {
	if err := task.EnqueuePrepullTask(ChallengeImages(challenge)); err != nil {
		log.Logger.Warningf("Failed to enqueue image warmup: challenge_id=%d error=%v", challenge.ID, err)
	}
}

func ListNodeImages(ctx context.Context) (map[string][]string, model.RetVal) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	return k8s.ListNodeImages(ctx)
}

func PullContestChallengeImage(ctx context.Context, form dto.PullImageForm) model.RetVal {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
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
		imageName = k8s.NormalizeImage(imageName)
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
