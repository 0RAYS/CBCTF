package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"slices"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListNodes(ctx context.Context) (*corev1.NodeList, model.RetVal) {
	nodes, err := kubeClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to list nodes: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "Nodes", "Error": err.Error()}}
	}
	return nodes, model.SuccessRetVal()
}

func ListSchedulableNodes(ctx context.Context) ([]*corev1.Node, model.RetVal) {
	allNodes, ret := ListNodes(ctx)
	if !ret.OK || allNodes == nil {
		return nil, ret
	}
	nodes := make([]*corev1.Node, 0)
	for _, node := range allNodes.Items {
		schedulable := true
		for _, taint := range node.Spec.Taints {
			if taint.Effect == corev1.TaintEffectNoSchedule {
				schedulable = false
				break
			}
		}
		if schedulable {
			nodes = append(nodes, &node)
		}
	}
	return nodes, model.SuccessRetVal()
}

func GetNodeImageList(ctx context.Context) (map[string][]string, model.RetVal) {
	nodes, ret := ListSchedulableNodes(ctx)
	if !ret.OK {
		return nil, ret
	}
	images := make(map[string][]string)
	for _, node := range nodes {
		images[node.Name] = make([]string, 0)
		for _, containerImage := range node.Status.Images {
			for _, name := range containerImage.Names {
				if strings.TrimSpace(name) == "" || strings.Contains(name, "@sha256:") {
					continue
				}
				if !slices.Contains(images[node.Name], name) {
					images[node.Name] = append(images[node.Name], name)
				}
			}
		}
	}
	return images, model.SuccessRetVal()
}
