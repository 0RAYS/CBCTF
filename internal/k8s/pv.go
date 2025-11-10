package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"context"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreatePVOptions struct {
	Name    string
	Labels  map[string]string
	Server  string
	Path    string
	Storage string
}

func CreatePV(ctx context.Context, options CreatePVOptions) (*corev1.PersistentVolume, bool, string) {
	var (
		pv      *corev1.PersistentVolume
		storage resource.Quantity
		err     error
	)
	storage, err = resource.ParseQuantity(options.Storage)
	if err != nil {
		log.Logger.Warningf("Failed to parse storage resource: %s", err)
		return nil, false, i18n.UnknownError
	}
	pv = &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:   options.Name,
			Labels: options.Labels,
		},
		Spec: corev1.PersistentVolumeSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
			Capacity: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceStorage: storage,
			},
			StorageClassName:              "",
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				NFS: &corev1.NFSVolumeSource{
					Server: options.Server,
					Path:   options.Path,
				},
			},
		},
	}
	pv, err = kubeClient.CoreV1().PersistentVolumes().Create(ctx, pv, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create pv: %s", err)
		return nil, false, i18n.CreatePVError
	}
	return pv, true, i18n.Success
}

func GetPV(ctx context.Context, name string) (*corev1.PersistentVolume, bool, string) {
	pv, err := kubeClient.CoreV1().PersistentVolumes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.PVNotFound
		}
		log.Logger.Warningf("Failed to get pv: %s", err)
		return nil, false, i18n.GetPVError
	}
	return pv, true, i18n.Success
}
