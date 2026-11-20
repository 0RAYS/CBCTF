package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type CreateConfigMapOptions struct {
	Name   string
	Labels map[string]string
	Data   map[string]string
}

func CreateConfigMap(ctx context.Context, options CreateConfigMapOptions) (*corev1.ConfigMap, model.RetVal) {
	var (
		configMap *corev1.ConfigMap
		err       error
	)
	configMap = &corev1.ConfigMap{
		OwnerReferences: resourceOwners(ctx),
		Name:            options.Name,
		Namespace:       globalNamespace,
		Labels:          options.Labels,
		Data:            options.Data,
	}
	configMap, err = kubeClient.CoreV1().ConfigMaps(globalNamespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create ConfigMap %s: %s", options.Name, err)
		return nil, model.FailRetVal(i18n.K8S.CreateError, "ConfigMap", options.Name, err)
	}
	return configMap, model.SuccessRetVal()
}

func DeleteConfigMapCollection(ctx context.Context, labels map[string]string) model.RetVal {
	return deleteCollection("ConfigMap", "ConfigMaps", labels, func(options metav1.ListOptions) error {
		return kubeClient.CoreV1().ConfigMaps(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	})
}
