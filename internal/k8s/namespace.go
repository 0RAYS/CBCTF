package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateNamespaceOptions struct {
	Name   string
	Labels map[string]string
}

func CreateNamespace(ctx context.Context, options CreateNamespaceOptions) (*corev1.Namespace, model.RetVal) {
	var (
		namespace *corev1.Namespace
		err       error
	)
	namespace = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   options.Name,
			Labels: options.Labels,
		},
	}
	namespace, err = kubeClient.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Namespace: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "Namespace", "Error": err.Error()}}
	}
	return namespace, model.SuccessRetVal()
}
