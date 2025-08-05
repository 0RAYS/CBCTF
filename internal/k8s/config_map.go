package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type CreateConfigMapOptions struct {
	Name   string
	Labels map[string]string
	Data   map[string]string
}

func CreateConfigMap(ctx context.Context, options CreateConfigMapOptions) (*corev1.ConfigMap, bool, string) {
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
		log.Logger.Warningf("Failed to create ConfigMap: %v", err)
		return nil, false, i18n.CreateConfigMapError
	}
	return configMap, true, i18n.Success
}

func GetConfigMap(ctx context.Context, name string) (*corev1.ConfigMap, bool, string) {
	configMap, err := kubeClient.CoreV1().ConfigMaps(globalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.ConfigMapNotFound
		}
		log.Logger.Warningf("Failed to get ConfigMap: %v", err)
		return nil, false, i18n.GetConfigMapError
	}
	return configMap, true, i18n.Success
}

func GetConfigMapList(ctx context.Context, labels ...map[string]string) (*corev1.ConfigMapList, bool, string) {
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
	configMapList, err := kubeClient.CoreV1().ConfigMaps(globalNamespace).List(ctx, options)
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.ConfigMapNotFound
		}
		log.Logger.Warningf("Failed to list ConfigMap: %v", err)
		return nil, false, i18n.GetConfigMapError
	}
	return configMapList, true, i18n.Success
}

func DeleteConfigMap(ctx context.Context, configMapName string) (bool, string) {
	err := kubeClient.CoreV1().ConfigMaps(globalNamespace).Delete(ctx, configMapName, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete ConfigMap: %v", err)
		return false, i18n.DeleteConfigMapError
	}
	return true, i18n.Success
}

func DeleteConfigMapList(ctx context.Context, labels ...map[string]string) (bool, string) {
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
		log.Logger.Warningf("Failed to delete ConfigMap: %v", err)
		return false, i18n.DeleteConfigMapError
	}
	return true, i18n.Success
}
