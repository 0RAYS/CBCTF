package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// DeleteCollection must never silently become a namespace/cluster-wide delete.
// Require a concrete workload owner (or a subnet owner for Kube-OVN IP cleanup),
// validate selector syntax and reject ambiguous multiple variadic arguments.
func deleteCollectionOptions(resource string, filters ...map[string]string) (metav1.ListOptions, model.RetVal) {
	invalid := func() (metav1.ListOptions, model.RetVal) {
		return metav1.ListOptions{}, model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": resource, "Error": "refusing deletion without a valid, non-empty workload ownership selector"}}
	}
	if len(filters) != 1 || len(filters[0]) == 0 {
		return invalid()
	}
	owned := false
	for key, value := range filters[0] {
		if key == "" || value == "" {
			return invalid()
		}
		if key == "victim_id" || key == "generator_id" {
			id, err := strconv.ParseUint(value, 10, 64)
			if err != nil || id == 0 {
				return invalid()
			}
			owned = true
		}
		if resource == "IP" && key == "ovn.kubernetes.io/subnet" {
			owned = true
		}
	}
	if !owned {
		return invalid()
	}
	selector, err := labels.ValidatedSelectorFromSet(labels.Set(filters[0]))
	if err != nil || selector.Empty() {
		return invalid()
	}
	return metav1.ListOptions{LabelSelector: selector.String()}, model.SuccessRetVal()
}
