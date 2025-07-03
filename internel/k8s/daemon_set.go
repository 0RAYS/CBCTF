package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/utils"
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateDaemonSetOptions struct {
	Images     []string
	PullPolicy string
}

func CreateDaemonSet(ctx context.Context, options CreateDaemonSetOptions) (*appsv1.DaemonSet, bool, string) {
	var (
		daemonSet *appsv1.DaemonSet
		err       error
	)
	containers := make([]corev1.Container, 0)
	for _, image := range options.Images {
		containers = append(containers, corev1.Container{
			Name:            fmt.Sprintf("%s-puller", image),
			ImagePullPolicy: corev1.PullPolicy(options.PullPolicy),
			Image:           image,
			Command:         []string{"echo", "Success"},
		})
	}
	name := utils.RandStr(10)
	daemonSet = &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("image-puller-%s", utils.RandStr(5)),
			Namespace: namespaceName,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"image-puller": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"image-puller": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers:    containers,
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
	daemonSet, err = kubeClient.AppsV1().DaemonSets(namespaceName).Create(ctx, daemonSet, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create DaemonSet: %v", err)
		return nil, false, i18n.CreateDaemonSetError
	}
	return daemonSet, true, i18n.Success
}

func GetDaemonSet(ctx context.Context, name string) (*appsv1.DaemonSet, bool, string) {
	daemonSet, err := kubeClient.AppsV1().DaemonSets(namespaceName).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to get DaemonSet: %v", err)
		return nil, false, i18n.GetDaemonSetError
	}
	return daemonSet, true, i18n.Success
}

func ListDaemonSets(ctx context.Context) (*appsv1.DaemonSetList, bool, string) {
	daemonSetList, err := kubeClient.AppsV1().DaemonSets(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to list DaemonSets: %v", err)
		return nil, false, i18n.GetDaemonSetError
	}
	return daemonSetList, true, i18n.Success
}

func DeleteDaemonSet(ctx context.Context, name string) (bool, string) {
	if err := kubeClient.AppsV1().DaemonSets(namespaceName).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		log.Logger.Warningf("Failed to delete DaemonSet: %v", err)
		return false, i18n.GetDaemonSetError
	}
	return true, i18n.Success
}
