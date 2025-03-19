package k8s

import (
	"CBCTF/internal/log"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNode 获取 Node 信息
func GetNode(ctx context.Context, name string) (*corev1.Node, bool, string) {
	node, err := Client.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to get Node %s: %v", name, err)
		return &corev1.Node{}, false, "GetNodeError"
	}
	return node, true, "Success"
}
