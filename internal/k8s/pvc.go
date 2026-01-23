package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
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

func CreatePVC(ctx context.Context, options CreatePVCOptions) (*corev1.PersistentVolumeClaim, model.RetVal) {
	var (
		pvc     *corev1.PersistentVolumeClaim
		storage resource.Quantity
		err     error
	)
	storage, err = resource.ParseQuantity(options.Storage)
	if err != nil {
		log.Logger.Warningf("Failed to parse storage resource: %s", err)
		return nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
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
			StorageClassName: utils.Ptr(""),
			VolumeName:       nfsVolumeName,
		},
	}
	pvc, err = kubeClient.CoreV1().PersistentVolumeClaims(globalNamespace).Create(ctx, pvc, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create pvc: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "PVC", "Error": err.Error()}}
	}
	return pvc, model.SuccessRetVal()
}

func GetPVC(ctx context.Context, name string) (*corev1.PersistentVolumeClaim, model.RetVal) {
	pvc, err := kubeClient.CoreV1().PersistentVolumeClaims(globalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "PVC"}}
		}
		log.Logger.Warningf("Failed to get pvc: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "PVC", "Error": err.Error()}}
	}
	return pvc, model.SuccessRetVal()
}
