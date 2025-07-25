package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreatePVCOptions struct {
	Name    string
	Labels  map[string]string
	Storage string
}

func CreatePVC(ctx context.Context, options CreatePVCOptions) (*corev1.PersistentVolumeClaim, bool, string) {
	var (
		pvc     *corev1.PersistentVolumeClaim
		storage resource.Quantity
		err     error
	)
	storage, err = resource.ParseQuantity(options.Storage)
	if err != nil {
		log.Logger.Warningf("Failed to parse storage resource: %s", err)
		return nil, false, i18n.UnknownError
	}
	pvc = &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   options.Name,
			Labels: options.Labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
			Resources: corev1.VolumeResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: storage,
				},
			},
		},
	}
	pvc, err = kubeClient.CoreV1().PersistentVolumeClaims(GlobalNamespace).Create(ctx, pvc, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create pvc: %s", err)
		return nil, false, i18n.CreatePVCError
	}
	return pvc, true, i18n.Success
}

func GetPVC(ctx context.Context, name string) (*corev1.PersistentVolumeClaim, bool, string) {
	pvc, err := kubeClient.CoreV1().PersistentVolumeClaims(GlobalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.PVCNotFound
		}
		log.Logger.Warningf("Failed to get pvc: %s", err)
		return nil, false, i18n.GetPVCError
	}
	return pvc, true, i18n.Success
}

func DeletePVC(ctx context.Context, name string) (bool, string) {
	err := kubeClient.CoreV1().PersistentVolumeClaims(GlobalNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete pvc: %s", err)
		return false, i18n.DeletePVCError
	}
	return true, i18n.Success
}
