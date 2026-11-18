package k8s

import (
	"context"
	"slices"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

// ListImagePullSecrets returns the names of the image pull secrets in the
// platform namespace, so an administrator can pick an existing secret instead
// of creating a new one.
func ListImagePullSecrets(ctx context.Context) ([]string, model.RetVal) {
	secrets, err := kubeClient.CoreV1().Secrets(globalNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to list image pull secrets: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "Secret", "Error": err.Error()}}
	}
	names := make([]string, 0, len(secrets.Items))
	for _, secret := range secrets.Items {
		switch secret.Type {
		case corev1.SecretTypeDockerConfigJson, corev1.SecretTypeDockercfg:
			names = append(names, secret.Name)
		}
	}
	slices.Sort(names)
	return names, model.SuccessRetVal()
}

// HasImagePullSecret reports whether the named image pull secret exists in the
// platform namespace.
func HasImagePullSecret(ctx context.Context, name string) (bool, model.RetVal) {
	if name == "" {
		return false, model.SuccessRetVal()
	}
	names, ret := ListImagePullSecrets(ctx)
	if !ret.OK {
		return false, ret
	}
	return slices.Contains(names, name), model.SuccessRetVal()
}
