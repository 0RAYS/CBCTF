package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/utils"
	"context"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type CreateJobOptions struct {
	Name         string
	Labels       map[string]string
	Images       []string
	PullPolicy   string
	NodeSelector map[string]string
}

func CreateJob(ctx context.Context, options CreateJobOptions) (*batchv1.Job, bool, string) {
	var (
		job *batchv1.Job
		err error
	)
	containers := make([]corev1.Container, 0)
	for _, image := range options.Images {
		containers = append(containers, corev1.Container{
			Name:            strings.ToLower(utils.RandStr(10)),
			ImagePullPolicy: corev1.PullPolicy(options.PullPolicy),
			Image:           image,
			Command:         []string{"echo", "Success"},
		})
	}
	job = &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: GlobalNamespace,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: utils.Ptr[int32](0),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("image-puller-%s", strings.ToLower(utils.RandStr(5))),
					Namespace: GlobalNamespace,
				},
				Spec: corev1.PodSpec{
					NodeSelector:  options.NodeSelector,
					Containers:    containers,
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
	job, err = kubeClient.BatchV1().Jobs(GlobalNamespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Job: %v", err)
		return nil, false, i18n.CreateJobError
	}
	return job, true, i18n.Success
}

func GetJob(ctx context.Context, name string) (*batchv1.Job, bool, string) {
	job, err := kubeClient.BatchV1().Jobs(GlobalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.JobNotFound
		}
		log.Logger.Warningf("Failed to get Job: %v", err)
		return nil, false, i18n.GetJobError
	}
	return job, true, i18n.Success
}

func ListJobs(ctx context.Context) (*batchv1.JobList, bool, string) {
	jobList, err := kubeClient.BatchV1().Jobs(GlobalNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to list Jobs: %v", err)
		return nil, false, i18n.GetJobError
	}
	return jobList, true, i18n.Success
}

func DeleteJob(ctx context.Context, name string) (bool, string) {
	err := kubeClient.BatchV1().Jobs(GlobalNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Job: %v", err)
		return false, i18n.GetJobError
	}
	return true, i18n.Success
}
