package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreatePVCOptions struct {
	Name    string
	Labels  map[string]string
	Storage string
}

func GetPVC(ctx context.Context, name string) (*corev1.PersistentVolumeClaim, model.RetVal) {
	pvc, err := kubeClient.CoreV1().PersistentVolumeClaims(globalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "PVC"}}
		}
		log.Logger.Warningf("Failed to get pvc: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "PVC", "Error": err.Error()}}
	}
	return pvc, model.SuccessRetVal()
}
