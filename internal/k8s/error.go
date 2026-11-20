package k8s

import (
	"errors"

	apierror "k8s.io/apimachinery/pkg/api/errors"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

// errMissingPodIdentity marks a deletion that was refused because the Pod
// could not be identified safely.
var errMissingPodIdentity = errors.New("refusing deletion without a Pod name and observed UID")

// errMissingNamespace marks a deletion that was refused because the target
// namespace was not given.
var errMissingNamespace = errors.New("refusing deletion without a namespace")

// errMissingConfigMap marks a creation that reported success without
// returning the created ConfigMap.
var errMissingConfigMap = errors.New("config map creation returned no object")

// errMissingPodObject marks a creation that reported success without
// returning the created Pod.
var errMissingPodObject = errors.New("pod creation returned no object")

// errMissingPodLookup marks a lookup that reported success without returning
// the Pod.
var errMissingPodLookup = errors.New("pod lookup returned no object")

// errMissingVMObject marks a creation that reported success without
// returning the created VirtualMachine.
var errMissingVMObject = errors.New("virtual machine creation returned no object")

// errMissingServiceObject marks a creation that reported success without
// returning the created Service.
var errMissingServiceObject = errors.New("service creation returned no object")

// errMissingVPCObject marks a creation that reported success without
// returning the created VPC.
var errMissingVPCObject = errors.New("vpc creation returned no object")

// getFailed reports a failed resource lookup: a missing resource is reported
// as not found, every other failure as a lookup error.
func getFailed(kind, target string, err error) model.RetVal {
	if apierror.IsNotFound(err) {
		return model.FailRetVal(i18n.K8S.NotFound, kind, target, nil)
	}
	log.Logger.Warningf("Failed to get %s %s: %s", kind, target, err)
	return model.FailRetVal(i18n.K8S.GetError, kind, target, err)
}
