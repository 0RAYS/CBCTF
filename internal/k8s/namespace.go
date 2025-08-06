package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateNamespaceOptions struct {
	Name   string
	Labels map[string]string
}

func CreateNamespace(ctx context.Context, options CreateNamespaceOptions) (*corev1.Namespace, bool, string) {
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
		return nil, false, i18n.CreateNamespaceError
	}
	return namespace, true, i18n.Success
}
