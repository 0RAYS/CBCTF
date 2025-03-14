package k8s

import (
	"CBCTF/internal/log"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// GetNode 获取 Node 信息
func GetNode(name string) (*corev1.Node, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	node, err := Client.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to get node: %v", err)
		return &corev1.Node{}, false, "GetNodeError"
	}
	return node, true, "Success"
}
