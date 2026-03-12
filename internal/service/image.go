package service

import (
	"CBCTF/internal/dto"
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
	for _, node := range nodes {
		images := form.Images
		if corev1.PullPolicy(form.PullPolicy) != corev1.PullAlways {
			images = slices.DeleteFunc(images, func(image string) bool {
				if strings.TrimSpace(image) == "" {
					return true
				}
				for _, containerImage := range node.Status.Images {
					for _, name := range containerImage.Names {
						if name == image {
							return true
						}
					}
				}
				return false
			})
		}
		if len(images) > 0 {
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
						"kubernetes.io/hostname": node.Name,
					},
				}); !ret.OK {
					return ret
				}
			}
		}
	}
	return model.SuccessRetVal()
}
