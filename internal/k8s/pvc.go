package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/model"
)

func GetPVC(ctx context.Context, name string) (*corev1.PersistentVolumeClaim, model.RetVal) {
	pvc, err := kubeClient.CoreV1().PersistentVolumeClaims(globalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, getFailed("PVC", name, err)
	}
	return pvc, model.SuccessRetVal()
}
