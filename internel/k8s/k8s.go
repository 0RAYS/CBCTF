package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"bytes"
	"context"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"os"
	"time"

	authorizationv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	kubeClient     *kubernetes.Clientset
	kubeConfig     *rest.Config
	adminAPIConfig *api.Config
	namespaceName  string
	err            error
)

func Run() {
	namespaceName = config.Env.K8S.Namespace
	if _, err = os.Stat(config.Env.K8S.Config.User); err != nil {
		log.Logger.Fatalf("Make sure the config.k8s.config.user configured correctly: %s", err)
	}
	kubeConfig, err = clientcmd.BuildConfigFromFlags("", config.Env.K8S.Config.User)
	if err != nil {
		log.Logger.Fatalf("Failed to load k8s user config: %s", err)
	}
	kubeConfig.QPS = 100
	kubeConfig.Burst = 200
	kubeClient, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init k8s client: %s", err)
	}
}

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if _, err = os.Stat(config.Env.K8S.Config.Admin); err != nil {
		log.Logger.Fatalf("Invalid config.k8s.config.admin: %s", err)
	}
	adminAPIConfig, err = clientcmd.LoadFromFile(config.Env.K8S.Config.Admin)
	if err != nil {
		log.Logger.Fatalf("Failed to load admin config: %s", err)
	}
	kubeConfig, err = clientcmd.BuildConfigFromFlags("", config.Env.K8S.Config.Admin)
	if err != nil {
		log.Logger.Fatalf("Failed to create client config: %s", err)
	}
	kubeClient, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init k8s client: %s", err)
	}
	log.Logger.Info("Admin config loaded")

	dynamicClient, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to create dynamic client: %s", err)
	}

	namespaceName = config.Env.K8S.Namespace
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(gen(namespaceName)), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			log.Logger.Fatalf("Failed to convert raw object to unstructured map: %s", err)
		}
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		gr, err := restmapper.GetAPIGroupResources(kubeClient.Discovery())
		if err != nil {
			log.Logger.Fatalf("Failed to get API group resources: %s", err)
		}
		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Logger.Fatalf("Failed to get REST mapping for GroupVersionKind: %s", err)
		}
		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace(namespaceName)
			}
			dri = dynamicClient.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dynamicClient.Resource(mapping.Resource)
		}
		if _, err = dri.Create(ctx, unstructuredObj, metav1.CreateOptions{}); err != nil {
			log.Logger.Warningf("Failed to create resource: %s", err)
			continue
		}
		log.Logger.Infof("Resource %s/%s created successfully", mapping.Resource.Resource, unstructuredObj.GetName())
	}
	if err != io.EOF {
		log.Logger.Fatalf("Unknown error: %s", err)
	}

	updateNodeIPs(ctx)
	if err := writeKubeConfig(ctx); err != nil {
		log.Logger.Fatalf("Failed to save kubeconfig to %s.conf: %s ", namespaceName, err)
	}
	config.Env.K8S.Config.User = fmt.Sprintf("%s.conf", namespaceName)
	tmp := config.Env.K8S.Config.Admin
	if err := config.Save(config.Env); err != nil {
		log.Logger.Fatalf("Failed to update config: %s", err)
	}
	log.Logger.Infof("Kubeconfig saved to %s.conf, please remove the %s and restart", namespaceName, tmp)
	os.Exit(0)
}

func updateNodeIPs(ctx context.Context) {
	nodes, err := kubeClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Fatalf("Failed to list nodes: %v", err)
	}
	var ips []string
	for _, node := range nodes.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeInternalIP || addr.Type == corev1.NodeExternalIP {
				ips = append(ips, addr.Address)
			}
		}
	}
	config.Env.K8S.Nodes = ips
}

// writeKubeConfig 写入一个低权限的 kubeconfig 文件
func writeKubeConfig(ctx context.Context) error {
	svcAccountName := fmt.Sprintf("%s-admin", namespaceName)
	svcAccountSecretName := fmt.Sprintf("%s-admin-secret", namespaceName)
	svcAccountToken, err := kubeClient.CoreV1().Secrets(namespaceName).Get(ctx, svcAccountSecretName, metav1.GetOptions{})
	if err != nil {
		log.Logger.Fatalf("Failed to get service account token secret: %v", err)
	}
	token := string(svcAccountToken.Data["token"])
	ca := svcAccountToken.Data["ca.crt"]
	host := kubeConfig.Host
	configCTX := adminAPIConfig.Contexts[adminAPIConfig.CurrentContext]
	return clientcmd.WriteToFile(api.Config{
		Clusters: map[string]*api.Cluster{
			configCTX.Cluster: {
				Server:                   host,
				CertificateAuthorityData: ca,
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			svcAccountName: {
				Token: token,
			},
		},
		Contexts: map[string]*api.Context{
			fmt.Sprintf("%s-admin@kubernetes-%s", namespaceName, svcAccountName): {
				Cluster:   configCTX.Cluster,
				AuthInfo:  svcAccountName,
				Namespace: namespaceName,
			},
		},
		CurrentContext: fmt.Sprintf("%s-admin@kubernetes-%s", namespaceName, svcAccountName),
	}, fmt.Sprintf("%s.conf", namespaceName))
}

// CheckPermission checks if the user has permission to access the resources
func CheckPermission() {
	if _, err = os.Stat(config.Env.K8S.Config.User); err != nil {
		log.Logger.Fatalf("Make sure the config.k8s.config.user configured correctly: %s", err)
	}
	kubeConfig, err = clientcmd.BuildConfigFromFlags("", config.Env.K8S.Config.User)
	if err != nil {
		log.Logger.Fatalf("Failed to load k8s user config: %s", err)
	}
	kubeClient, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init k8s client: %s", err)
	}
	log.Logger.Infof("Checking permission in namespace %s", namespaceName)
	groups := map[string][]string{
		"":                      {"pods", "services", "configmaps", "pods/exec"},
		"networking.k8s.io":     {"networkpolicies"},
		"crd.projectcalico.org": {"ippools"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	for group, resources := range groups {
		for _, resource := range resources {
			accessReview := &authorizationv1.SelfSubjectAccessReview{
				Spec: authorizationv1.SelfSubjectAccessReviewSpec{
					ResourceAttributes: &authorizationv1.ResourceAttributes{
						Namespace: namespaceName,
						Group:     group,
						Version:   "*",
						Resource:  resource,
						Verb:      "*",
					},
				},
			}
			res, err := kubeClient.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, accessReview, metav1.CreateOptions{})
			if err != nil {
				log.Logger.Warningf("Failed to check permissions: %v", err)
			}
			if !res.Status.Allowed {
				log.Logger.Warningf("User does NOT have permission to access %s-%s in namespace cbctf.", group, resource)
				log.Logger.Warningf("Reason: %s", res.Status.Reason)
				log.Logger.Warningf("EvaluationError: %s", res.Status.EvaluationError)
			}
		}
	}
	log.Logger.Infof("User has permission to access all needed resources in namespace %s", namespaceName)
}
