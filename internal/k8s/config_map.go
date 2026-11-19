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
	options, ret := deleteCollectionOptions("ConfigMap", labels)
	if !ret.OK {
		return ret
	}
	err := kubeClient.CoreV1().ConfigMaps(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete ConfigMaps (%s): %s", options.LabelSelector, err)
		return model.FailRetVal(i18n.K8S.DeleteError, "ConfigMap", options.LabelSelector, err)
	}
	return model.SuccessRetVal()
}
