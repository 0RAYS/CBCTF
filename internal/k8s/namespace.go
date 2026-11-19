package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
)

func GetNamespace(ctx context.Context, name string) (*corev1.Namespace, model.RetVal) {
	ns, err := kubeClient.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, model.FailRetVal(i18n.K8S.GetError, "Namespace", name, err)
	}
	return ns, model.SuccessRetVal()
}
