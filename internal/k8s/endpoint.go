package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"

	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteEndpointCollection(ctx context.Context, labels ...map[string]string) model.RetVal {
	options, ret := deleteCollectionOptions("EndpointSlice", labels...)
	if !ret.OK {
		return ret
	}
	err := kubeClient.DiscoveryV1().EndpointSlices(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete EndpointSlice: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "EndpointSlice", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
