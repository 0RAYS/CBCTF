package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/model"
)

func DeleteIPCollection(ctx context.Context, labels map[string]string) model.RetVal {
	return deleteCollection("IP", "IPs", labels, func(options metav1.ListOptions) error {
		return ovnClient.KubeovnV1().IPs().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	})
}
