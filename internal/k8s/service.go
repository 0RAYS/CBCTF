package k8s

import (
	"CBCTF/internal/log"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// DeleteService 删除 Service, 目前主要是靶机的端口映射
func DeleteService(name string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	err := Client.CoreV1().Services(NamespaceName).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to delete service: %v", err)
		return false, "DeleteServiceError"
	}
	return true, "Success"
}
