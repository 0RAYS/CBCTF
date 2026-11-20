package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/model"
)

func DeleteEndpointCollection(ctx context.Context, labels map[string]string) model.RetVal {
	return deleteCollection("EndpointSlice", "EndpointSlices", labels, func(options metav1.ListOptions) error {
		return kubeClient.DiscoveryV1().EndpointSlices(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	})
}
