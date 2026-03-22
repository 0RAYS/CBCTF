package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: globalNamespace,
			Labels:    options.Labels,
		},
		Data: options.Data,
	}
	configMap, err = kubeClient.CoreV1().ConfigMaps(globalNamespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create ConfigMap: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "ConfigMap", "Error": err.Error()}}
	}
	return configMap, model.SuccessRetVal()
}

func DeleteConfigMapList(ctx context.Context, labels ...map[string]string) model.RetVal {
	var options metav1.ListOptions
	if len(labels) > 0 {
		var selector string
		for k, v := range labels[0] {
			selector += fmt.Sprintf("%s=%s,", k, v)
		}
		options = metav1.ListOptions{
			LabelSelector: strings.TrimSuffix(selector, ","),
		}
	}
	err := kubeClient.CoreV1().ConfigMaps(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete ConfigMap: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "ConfigMap", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
