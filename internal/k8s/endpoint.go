package k8s

import (
	"context"

	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

func DeleteEndpointCollection(ctx context.Context, labels map[string]string) model.RetVal {
	options, ret := deleteCollectionOptions("EndpointSlice", labels)
	if !ret.OK {
		return ret
	}
	err := kubeClient.DiscoveryV1().EndpointSlices(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete EndpointSlices (%s): %s", options.LabelSelector, err)
		return model.FailRetVal(i18n.K8S.DeleteError, "EndpointSlice", options.LabelSelector, err)
	}
	return model.SuccessRetVal()
}
