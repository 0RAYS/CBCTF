package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"context"
	"fmt"
	"os"
	"time"

	kubeovnclient "github.com/JBNRZ/kubeovn-api/pkg/client/clientset"
	authorizationv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	VictimPodTag    = "victim"
	GeneratorPodTag = "generator"
)

var (
	kubeClient     *kubernetes.Clientset
	kubeOVNClient  *kubeovnclient.Clientset
	kubeConfig     *rest.Config
	adminAPIConfig *api.Config
	namespaceName  string
	ipPoolName     string
)

func Init(run bool) {
	var err error
	namespaceName = config.Env.K8S.Namespace
	ipPoolName = fmt.Sprintf("%s-ip-pool", namespaceName)
	if run {
		if _, err = os.Stat(config.Env.K8S.Config); err != nil {
			log.Logger.Fatalf("Make sure the config.k8s.config.user configured correctly: %s", err)
		}
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", config.Env.K8S.Config)
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
}

func InitResources() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	loadKubeConfig()
	initClients()

	log.Logger.Debugf("Checking resources in namespace %s", namespaceName)
	updateNodeIPs(ctx)
	ensureNamespace(ctx)

	tmp := config.Env.K8S.Config
	if err := config.Save(config.Env); err != nil {
		log.Logger.Fatalf("Failed to update config: %s", err)
	}
	log.Logger.Infof("Kubeconfig saved to %s.conf, please remove the %s and restart", namespaceName, tmp)
	os.Exit(0)
}

// CheckPermission checks if the user has permission to access the resources
func CheckPermission() {
	var err error
	if _, err := os.Stat(config.Env.K8S.Config); err != nil {
		log.Logger.Fatalf("Make sure the config.k8s.config.user configured correctly: %s", err)
	}
	kubeConfig, err = clientcmd.BuildConfigFromFlags("", config.Env.K8S.Config)
	if err != nil {
		log.Logger.Fatalf("Failed to load k8s user config: %s", err)
	}
	kubeClient, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init k8s client: %s", err)
	}
	log.Logger.Infof("Checking permission in namespace %s", namespaceName)
	groups := map[string]map[string][]string{
		"":                  {"pods": {"*"}, "services": {"*"}, "configmaps": {"*"}, "pods/exec": {"*"}, "nodes": {"get", "list", "watch"}},
		"batch":             {"jobs": {"*"}},
		"networking.k8s.io": {"networkpolicies": {"*"}},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	for group, resources := range groups {
		for resource, verbL := range resources {
			for _, verb := range verbL {
				accessReview := &authorizationv1.SelfSubjectAccessReview{
					Spec: authorizationv1.SelfSubjectAccessReviewSpec{
						ResourceAttributes: &authorizationv1.ResourceAttributes{
							Namespace: namespaceName,
							Group:     group,
							Version:   "*",
							Resource:  resource,
							Verb:      verb,
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
					os.Exit(-1)
				}
			}
		}
	}
	log.Logger.Infof("User has permission to access all needed resources in namespace %s", namespaceName)
}

func loadKubeConfig() {
	if _, err := os.Stat(config.Env.K8S.Config); err != nil {
		log.Logger.Fatalf("Invalid config.k8s.config.admin: %s", err)
	}
	var err error
	adminAPIConfig, err = clientcmd.LoadFromFile(config.Env.K8S.Config)
	if err != nil {
		log.Logger.Fatalf("Failed to load admin config: %s", err)
	}
	kubeConfig, err = clientcmd.NewNonInteractiveClientConfig(*adminAPIConfig, adminAPIConfig.CurrentContext, &clientcmd.ConfigOverrides{}, nil).ClientConfig()
	if err != nil {
		log.Logger.Fatalf("Failed to create client config: %s", err)
	}
	log.Logger.Info("Admin config loaded")
}

func initClients() {
	var err error
	kubeClient, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init k8s client: %s", err)
	}
	kubeOVNClient, err = kubeovnclient.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init KubeOVN client: %s", err)
	}
}

func updateNodeIPs(ctx context.Context) {
	ips, ok, _ := GetNodeIPList(ctx)
	if !ok {
		os.Exit(-1)
	}
	config.Env.K8S.Nodes = ips
}

func ensureNamespace(ctx context.Context) {
	if _, err := kubeClient.CoreV1().Namespaces().Get(ctx, namespaceName, metav1.GetOptions{}); err != nil {
		log.Logger.Infof("Creating namespace %s...", namespaceName)
		_, err = kubeClient.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: namespaceName},
		}, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Failed to create namespace: %v", err)
		}
	}
}
