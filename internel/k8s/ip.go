package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateIPOptions struct {
	Name    string
	Labels  map[string]string
	Subnet  string
	PodName string
	IP      string
}

func CreateIP(ctx context.Context, options CreateIPOptions) (*kubeovnv1.IP, bool, string) {
	var (
		ip  *kubeovnv1.IP
		err error
	)
	ip = &kubeovnv1.IP{
		ObjectMeta: metav1.ObjectMeta{
			Name: options.Name,
		},
		Spec: kubeovnv1.IPSpec{
			Subnet:      options.Subnet,
			PodType:     "",
			Namespace:   GlobalNamespace,
			PodName:     options.PodName,
			V4IPAddress: options.IP,
		},
	}
	ip, err = kubeOVNClient.KubeovnV1().IPs().Create(ctx, ip, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create IP: %s", err)
		return nil, false, i18n.CreateIPError
	}
	return ip, true, i18n.Success
}

func GetIP(ctx context.Context, name string) (*kubeovnv1.IP, bool, string) {
	ip, err := kubeOVNClient.KubeovnV1().IPs().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.IPNotFound
		}
		log.Logger.Warningf("Failed to get IP: %s", err)
		return nil, false, i18n.GetIPError
	}
	return ip, true, i18n.Success
}

func GetIPList(ctx context.Context) (*kubeovnv1.IPList, bool, string) {
	ipList, err := kubeOVNClient.KubeovnV1().IPs().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to get IP list: %s", err)
		return nil, false, i18n.GetIPError
	}
	return ipList, true, i18n.Success
}

func DeleteIP(ctx context.Context, name string) (bool, string) {
	err := kubeOVNClient.KubeovnV1().IPs().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete IP: %s", err)
		return false, i18n.DeleteIPError
	}
	return true, i18n.Success
}
