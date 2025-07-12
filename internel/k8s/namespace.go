package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
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
		log.Logger.Warningf("Failed to create Namespace: %v", err)
		return nil, false, i18n.CreateNamespaceError
	}
	return namespace, true, i18n.Success
}

func GetNamespace(ctx context.Context, name string) (*corev1.Namespace, bool, string) {
	namespace, err := kubeClient.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.NamespaceNotFound
		}
		log.Logger.Warningf("Failed to get Namespace: %v", err)
		return nil, false, i18n.GetNamespaceError
	}
	return namespace, true, i18n.Success
}

func GetNamespaceList(ctx context.Context) (*corev1.NamespaceList, bool, string) {
	namespaceList, err := kubeClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to list Namespace: %v", err)
		return nil, false, i18n.GetNamespaceError
	}
	return namespaceList, true, i18n.Success
}

func DeleteNamespace(ctx context.Context, name string) (bool, string) {
	err := kubeClient.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Namespace: %v", err)
		return false, i18n.DeleteNamespaceError
	}
	return true, i18n.Success
}
