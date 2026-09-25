package k8s

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
)

type CreateJobOptions struct {
	Name         string
	Labels       map[string]string
	Images       []string
	PullPolicy   string
	SelectedNode string
}

func CreateJob(ctx context.Context, options CreateJobOptions) (*batchv1.Job, model.RetVal) {
	var (
		job *batchv1.Job
		err error
	)
	job = &batchv1.Job{
		Name:      options.Name,
		Namespace: globalNamespace,
		Labels:    options.Labels,
		Spec: batchv1.JobSpec{
			BackoffLimit:            new(int32(0)),
			TTLSecondsAfterFinished: new(int32(3600)),
			ActiveDeadlineSeconds:   new(int64(600)),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    options.Labels,
					Name:      fmt.Sprintf("image-puller-%s", utils.RandHexStr(5)),
					Namespace: globalNamespace,
				},
				Spec: corev1.PodSpec{
					AutomountServiceAccountToken: new(false),
					EnableServiceLinks:           new(false),
					Affinity: func() *corev1.Affinity {
						if options.SelectedNode != "" {
							return &corev1.Affinity{
								NodeAffinity: &corev1.NodeAffinity{
									RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
										NodeSelectorTerms: []corev1.NodeSelectorTerm{
											{
												MatchFields: []corev1.NodeSelectorRequirement{
													{
														Key:      "metadata.name",
														Operator: corev1.NodeSelectorOpIn,
														Values:   []string{options.SelectedNode},
													},
												},
											},
										},
									},
								},
							}
						}
						return nil
					}(),
					Containers: func() []corev1.Container {
						containers := make([]corev1.Container, 0, len(options.Images))
						for _, image := range options.Images {
							containers = append(containers, corev1.Container{
								Name:            utils.RandHexStr(10),
								ImagePullPolicy: corev1.PullPolicy(options.PullPolicy),
								Image:           image,
								Command:         []string{"echo", "Success"},
								RestartPolicy:   new(corev1.ContainerRestartPolicyNever),
							})
						}
						return containers
					}(),
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
	job, err = kubeClient.BatchV1().Jobs(globalNamespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Job: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "Job", "Error": err.Error()}}
	}
	return job, model.SuccessRetVal()
}
