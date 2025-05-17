package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateConfigMap(ctx context.Context, configMapName string, data map[string]string) (*corev1.ConfigMap, bool, string) {
	if _, ok, _ := GetConfigMap(ctx, configMapName); ok {
		DeleteConfigMap(ctx, configMapName)
	}
	var err error
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: NamespaceName,
		},
		Data: data,
	}
	configMap, err = client.CoreV1().ConfigMaps(NamespaceName).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil && !apierror.IsAlreadyExists(err) {
		log.Logger.Warningf("Failed to create ConfigMap: %v", err)
		return nil, false, i18n.CreateConfigMapError
	}
	return configMap, true, i18n.Success
}

func GetConfigMap(ctx context.Context, configMapName string) (*corev1.ConfigMap, bool, string) {
	configMap, err := client.CoreV1().ConfigMaps(NamespaceName).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.ConfigMapNotFound
		}
		log.Logger.Warningf("Failed to get ConfigMap: %v", err)
		return nil, false, i18n.GetConfigMapError
	}
	return configMap, true, i18n.Success
}

func DeleteConfigMap(ctx context.Context, configMapName string) (bool, string) {
	err := client.CoreV1().ConfigMaps(NamespaceName).Delete(ctx, configMapName, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete ConfigMap: %v", err)
		return false, i18n.DeleteConfigMapError
	}
	return true, i18n.Success
}
