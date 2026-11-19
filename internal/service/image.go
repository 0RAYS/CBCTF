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
	"CBCTF/internal/redis"
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

func ContestChallengeImages(contestChallenge model.ContestChallenge) []string {
	return k8s.ContestChallengeImages(contestChallenge, challengeSidecarImages(contestChallenge.Type)...)
}

func warmContestChallengeImages(contestChallenge model.ContestChallenge) {
	if err := task.EnqueuePrepullTask(ContestChallengeImages(contestChallenge)); err != nil {
		log.Logger.Warningf("Failed to enqueue image warmup: contest_challenge_id=%d error=%v", contestChallenge.ID, err)
	}
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

// ListNodeNames returns the schedulable nodes the detection and the warmup target.
func ListNodeNames(ctx context.Context) ([]string, model.RetVal) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	nodes, ret := k8s.ListSchedulableNodes(ctx)
	if !ret.OK {
		return nil, ret
	}
	names := make([]string, 0, len(nodes))
	for _, node := range nodes {
		names = append(names, node.Name)
	}
	return names, model.SuccessRetVal()
}

// GlobalImageScope identifies one detection scope, so the last result of a scope
// never mixes with another one.
func GlobalImageScope() string { return "global" }

func ContestImageScope(contestID uint) string { return fmt.Sprintf("contest:%d", contestID) }

// DetectImages starts a detection run for one scope. The run probes the nodes
// with imagePullPolicy: Never and keeps only its last result.
func DetectImages(ctx context.Context, scope string, images []string) model.RetVal {
	batch := model.NewBatch(1)
	if len(images) == 0 {
		batch.Skip(scope, "no images")
		return batch.Result(ctx)
	}
	if err := task.EnqueueImageDetect(scope, images); err != nil {
		log.Logger.Warningf("Failed to enqueue image detection: scope=%s error=%v", scope, err)
		batch.Fail(scope, "enqueue", model.FailRetVal(i18n.Task.EnqueueError, "ImageDetect", scope, err))
		return batch.Result(ctx)
	}
	batch.Success(scope, "queued")
	return batch.Result(ctx)
}

// LoadImageDetection returns the last detection result of a scope, if any.
func LoadImageDetection(ctx context.Context, scope string) *model.ImageDetection {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	detection, err := redis.LoadImageDetection(ctx, scope)
	if err != nil {
		log.Logger.Warningf("Failed to load the image detection result: scope=%s error=%v", scope, err)
		return nil
	}
	return detection
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
			return model.FailRetVal(i18n.Response.BadRequest, "Node", nodeName, fmt.Errorf("unknown node: %s", nodeName))
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

	names := make([]string, 0, len(targetImages))
	for name := range targetImages {
		names = append(names, name)
	}
	slices.Sort(names)
	batch := model.NewBatch(len(names))
	for _, nodeName := range names {
		if ctx.Err() != nil {
			return batch.Result(ctx)
		}
		images := targetImages[nodeName]
		if err := task.EnqueuePrepullTargets(images, []string{nodeName}, form.PullPolicy); err != nil {
			log.Logger.Warningf("Warmup enqueue failed: node=%s error=%v", nodeName, err)
			batch.Fail(nodeName, "enqueue", model.FailRetVal(i18n.Task.EnqueueError, "Node", nodeName, err))
			return batch.Result(ctx)
		}
		batch.Success(nodeName, "queued")
	}
	return batch.Result(ctx)
}
