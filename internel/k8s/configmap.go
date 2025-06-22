package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/utils"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type CreateConfigMapOptions struct {
	PodName string
	Data    map[string]string
}

func CreateConfigMap(ctx context.Context, options CreateConfigMapOptions) (*corev1.ConfigMap, bool, string) {
	var (
		configMap *corev1.ConfigMap
		err       error
	)
	configMap = &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("cm-%s", strings.ToLower(utils.RandStr(10))),
			Namespace: namespaceName,
			Labels: map[string]string{
				VictimPodTag: options.PodName,
			},
		},
		Data: options.Data,
	}
	configMap, err = kubeClient.CoreV1().ConfigMaps(namespaceName).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create ConfigMap: %v", err)
		return nil, false, i18n.CreateConfigMapError
	}
	return configMap, true, i18n.Success
}

func GetConfigMapList(ctx context.Context) (*corev1.ConfigMapList, bool, string) {
	configMapList, err := kubeClient.CoreV1().ConfigMaps(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.ConfigMapNotFound
		}
		log.Logger.Warningf("Failed to list ConfigMap: %v", err)
		return nil, false, i18n.GetConfigMapError
	}
	return configMapList, true, i18n.Success
}

func GetConfigMapListByPodName(ctx context.Context, key, podName string) (*corev1.ConfigMapList, bool, string) {
	configMapList, err := kubeClient.CoreV1().ConfigMaps(namespaceName).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", key, podName),
	})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.ConfigMapNotFound
		}
		log.Logger.Warningf("Failed to list Pod %s ConfigMap: %v", podName, err)
		return nil, false, i18n.GetConfigMapError
	}
	return configMapList, true, i18n.Success
}

func DeleteConfigMap(ctx context.Context, configMapName string) (bool, string) {
	err := kubeClient.CoreV1().ConfigMaps(namespaceName).Delete(ctx, configMapName, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete ConfigMap: %v", err)
		return false, i18n.DeleteConfigMapError
	}
	return true, i18n.Success
}

// DeleteConfigMapListByPodName TODO: 有可能删不干净
func DeleteConfigMapListByPodName(ctx context.Context, key, podName string) (bool, string) {
	configMapList, ok, msg := GetConfigMapListByPodName(ctx, key, podName)
	if !ok {
		if msg != i18n.ConfigMapNotFound {
			return false, msg
		}
		return true, i18n.Success
	}
	for _, cm := range configMapList.Items {
		if ok, msg = DeleteConfigMap(ctx, cm.Name); !ok {
			return false, msg
		}
	}
	return true, i18n.Success
}
