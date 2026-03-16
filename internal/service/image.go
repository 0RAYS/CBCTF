package service

import (
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

func GetNodeImageList() (map[string][]string, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return k8s.GetNodeImageList(ctx)
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

		if corev1.PullPolicy(form.PullPolicy) != corev1.PullAlways && slices.ContainsFunc(node.Status.Images, func(image corev1.ContainerImage) bool {
			for _, name := range image.Names {
				if name == imageName {
					return true
				}
			}
			return false
		}) {
			continue
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
		var chunks [][]string
		for i := 0; i < len(images); i += 5 {
			end := i + 5
			if end > len(images) {
				end = len(images)
			}
			chunks = append(chunks, images[i:end])
		}
		for _, chunk := range chunks {
			if _, ret = k8s.CreateJob(ctx, k8s.CreateJobOptions{
				Name:       fmt.Sprintf("image-puller-%s", utils.RandStr(5)),
				Images:     chunk,
				PullPolicy: form.PullPolicy,
				NodeSelector: map[string]string{
					"kubernetes.io/hostname": nodeName,
				},
			}); !ret.OK {
				return ret
			}
		}
	}
	return model.SuccessRetVal()
}
