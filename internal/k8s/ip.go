package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"context"
	"fmt"
	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

type CreateIPOptions struct {
	Name    string
	Labels  map[string]string
	Subnet  string
	PodName string
	IP      string
}

func GetIPList(ctx context.Context, labels ...map[string]string) (*kubeovnv1.IPList, bool, string) {
	var options metav1.ListOptions
	if len(labels) > 0 {
		var selector string
		for k, v := range labels[0] {
			selector += fmt.Sprintf("%s=%s,", k, v)
		}
		options = metav1.ListOptions{
			LabelSelector: strings.TrimSuffix(selector, ","),
		}
	}
	ipList, err := kubeOVNClient.KubeovnV1().IPs().List(ctx, options)
	if err != nil {
		log.Logger.Warningf("Failed to get IP list: %v", err)
		return nil, false, i18n.GetIPError
	}
	return ipList, true, i18n.Success
}

func DeleteIPList(ctx context.Context, labels ...map[string]string) (bool, string) {
	var options metav1.ListOptions
	if len(labels) > 0 {
		var selector string
		for k, v := range labels[0] {
			selector += fmt.Sprintf("%s=%s,", k, v)
		}
		options = metav1.ListOptions{
			LabelSelector: strings.TrimSuffix(selector, ","),
		}
	}
	err := kubeOVNClient.KubeovnV1().IPs().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete IP list: %v", err)
		return false, i18n.DeleteIPError
	}
	return true, i18n.Success
}

func DeleteIPForce(ctx context.Context, name string) (bool, string) {
	_, err := kubeOVNClient.KubeovnV1().IPs().Patch(ctx, name, types.MergePatchType, forceDelete, metav1.PatchOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete IP: %v", err)
		return false, i18n.DeleteIPError
	}
	return true, i18n.Success
}

func DeleteIPListForce(ctx context.Context, labels ...map[string]string) (bool, string) {
	ipL, ok, msg := GetIPList(ctx, labels...)
	if !ok {
		return false, msg
	}
	for _, ip := range ipL.Items {
		if ok, msg = DeleteIPForce(ctx, ip.Name); !ok {
			return false, msg
		}
	}
	return true, i18n.Success
}
