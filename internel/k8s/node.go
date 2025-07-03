package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"slices"
)

func ListNodes(ctx context.Context) (*corev1.NodeList, bool, string) {
	nodes, err := kubeClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to list nodes: %v", err)
		return nodes, false, i18n.GetNodeListError
	}
	return nodes, true, i18n.Success
}

func ListSchedulableNodes(ctx context.Context) ([]*corev1.Node, bool, string) {
	allNodes, ok, msg := ListNodes(ctx)
	if !ok {
		return nil, ok, msg
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
	return nodes, true, i18n.Success
}

func GetNodeIPList(ctx context.Context) ([]string, bool, string) {
	nodes, ok, msg := ListNodes(ctx)
	if !ok {
		return make([]string, 0), false, msg
	}
	ips := make([]string, 0)
	for _, node := range nodes.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeInternalIP || addr.Type == corev1.NodeExternalIP {
				ips = append(ips, addr.Address)
			}
		}
	}
	return ips, true, i18n.Success
}

func GetNodeImageList(ctx context.Context) (map[string][]string, bool, string) {
	nodes, ok, msg := ListSchedulableNodes(ctx)
	if !ok {
		return make(map[string][]string), false, msg
	}
	images := make(map[string][]string)
	for _, node := range nodes {
		images[node.Name] = make([]string, 0)
		for _, containerImage := range node.Status.Images {
			for _, name := range containerImage.Names {
				if !slices.Contains(images[node.Name], name) {
					images[node.Name] = append(images[node.Name], name)
				}
			}
		}
	}
	return images, true, i18n.Success
}
