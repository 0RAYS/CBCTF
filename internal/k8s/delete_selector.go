package k8s

import (
	"errors"
	"strconv"

	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

// errMissingOwnershipSelector marks a deletion that was refused because the
// workload ownership selector was missing or invalid.
var errMissingOwnershipSelector = errors.New("refusing deletion without a valid, non-empty workload ownership selector")

// deleteCollection removes every resource of a collection that matches the
// workload ownership selector. The selector is built and validated here, so a
// deletion can never widen into an unscoped one.
func deleteCollection(kind, plural string, filter map[string]string, del func(metav1.ListOptions) error) model.RetVal {
	options, ret := deleteCollectionOptions(kind, filter)
	if !ret.OK {
		return ret
	}
	err := del(options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete %s (%s): %s", plural, options.LabelSelector, err)
		return model.FailRetVal(i18n.K8S.DeleteError, kind, options.LabelSelector, err)
	}
	return model.SuccessRetVal()
}

// deleteCollectionOptions builds the selector of a collection deletion. It
// rejects everything that is not a concrete workload owner (or a subnet owner
// for Kube-OVN IP cleanup), so a delete can never widen into a namespace-wide
// or cluster-wide one.
func deleteCollectionOptions(resource string, filter map[string]string) (metav1.ListOptions, model.RetVal) {
	invalid := func() (metav1.ListOptions, model.RetVal) {
		return metav1.ListOptions{}, model.FailRetVal(i18n.K8S.DeleteError, resource, "", errMissingOwnershipSelector)
	}
	if len(filter) == 0 {
		return invalid()
	}
	owned := false
	for key, value := range filter {
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
	selector, err := labels.ValidatedSelectorFromSet(filter)
	if err != nil || selector.Empty() {
		return invalid()
	}
	return metav1.ListOptions{LabelSelector: selector.String()}, model.SuccessRetVal()
}
