package k8s

import (
	"context"
	"slices"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
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
		schedulable := !node.Spec.Unschedulable
		ready := false
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
				ready = true
			}
		}
		schedulable = schedulable && ready
		for _, taint := range node.Spec.Taints {
			if taint.Effect == corev1.TaintEffectNoSchedule || taint.Effect == corev1.TaintEffectNoExecute {
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

func ListNodeImages(ctx context.Context) (map[string][]string, model.RetVal) {
	nodes, ret := ListSchedulableNodes(ctx)
	if !ret.OK {
		return nil, ret
	}
	images := make(map[string][]string)
	for _, node := range nodes {
		images[node.Name] = make([]string, 0)
		for name := range nodeImageSet(node) {
			images[node.Name] = append(images[node.Name], name)
		}
		slices.Sort(images[node.Name])
	}
	return images, model.SuccessRetVal()
}

// Use the same canonical inventory for display and pre-pull decisions.
func nodeImageSet(node *corev1.Node) map[string]struct{} {
	images := make(map[string]struct{})
	for _, entry := range node.Status.Images {
		for _, name := range entry.Names {
			if name = NormalizeImage(name); name != "" {
				images[name] = struct{}{}
			}
		}
	}
	return images
}
