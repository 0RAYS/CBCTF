package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

func GetPVC(ctx context.Context, name string) (*corev1.PersistentVolumeClaim, model.RetVal) {
	pvc, err := kubeClient.CoreV1().PersistentVolumeClaims(globalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.FailRetVal(i18n.K8S.NotFound, "PVC", name, nil)
		}
		log.Logger.Warningf("Failed to get PVC %s: %s", name, err)
		return nil, model.FailRetVal(i18n.K8S.GetError, "PVC", name, err)
	}
	return pvc, model.SuccessRetVal()
}
