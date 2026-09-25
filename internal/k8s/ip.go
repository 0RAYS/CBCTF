package k8s

import (
	"context"

	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

func DeleteIPCollection(ctx context.Context, labels map[string]string) model.RetVal {
	options, ret := deleteCollectionOptions("IP", labels)
	if !ret.OK {
		return ret
	}
	err := ovnClient.KubeovnV1().IPs().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete IP list: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "IP", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
