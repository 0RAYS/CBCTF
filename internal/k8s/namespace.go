package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNamespace(ctx context.Context, name string) (*corev1.Namespace, model.RetVal) {
	ns, err := kubeClient.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "Namespace", "Error": err.Error()}}
	}
	return ns, model.SuccessRetVal()
}
