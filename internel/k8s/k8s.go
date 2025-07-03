package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"context"
	"fmt"
	projectcalicov3 "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
	rbacv1 "k8s.io/api/rbac/v1"
	"os"
	"time"

	"github.com/projectcalico/api/pkg/client/clientset_generated/clientset"
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
	kubeClient                  *kubernetes.Clientset
	calicoClient                *clientset.Clientset
	svcAccountToken             *corev1.Secret
	kubeConfig                  *rest.Config
	adminAPIConfig              *api.Config
	namespaceName               string
	svcAccountName              string
	svcAccountSecretName        string
	ipPoolName                  string
	adminRoleName               string
	adminRoleBindingName        string
	adminClusterRoleName        string
	adminClusterRoleBindingName string
)

func Init(run bool) {
	var err error
	namespaceName = config.Env.K8S.Namespace
	svcAccountName = fmt.Sprintf("%s-admin", namespaceName)
	svcAccountSecretName = fmt.Sprintf("%s-admin-secret", namespaceName)
	ipPoolName = fmt.Sprintf("%s-ip-pool", namespaceName)
	adminRoleName = fmt.Sprintf("%s-admin-role", namespaceName)
	adminRoleBindingName = fmt.Sprintf("%s-admin-role-binding", namespaceName)
	adminClusterRoleName = fmt.Sprintf("%s-admin-cluster-role", namespaceName)
	adminClusterRoleBindingName = fmt.Sprintf("%s-admin-cluster-role-binding", namespaceName)
	if run {
		if _, err := os.Stat(config.Env.K8S.Config.User); err != nil {
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
}

func InitResources() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	loadKubeConfig()
	initClients()

	log.Logger.Debugf("Checking resources in namespace %s", namespaceName)
	updateNodeIPs(ctx)
	ensureNamespace(ctx)
	ensureServiceAccount(ctx)
	ensureSecret(ctx)
	ensureIPPool(ctx)
	ensureRole(ctx)
	ensureRoleBinding(ctx)
	ensureClusterRole(ctx)
	ensureClusterRoleBinding(ctx)

	if err := writeKubeConfig(); err != nil {
		log.Logger.Fatalf("Failed to save kubeconfig to %s.conf: %s ", namespaceName, err)
	}
	config.Env.K8S.Config.User = fmt.Sprintf("./%s.conf", namespaceName)
	tmp := config.Env.K8S.Config.Admin
	if err := config.Save(config.Env); err != nil {
		log.Logger.Fatalf("Failed to update config: %s", err)
	}
	log.Logger.Infof("Kubeconfig saved to %s.conf, please remove the %s and restart", namespaceName, tmp)
	os.Exit(0)
}

// CheckPermission checks if the user has permission to access the resources
func CheckPermission() {
	var err error
	if _, err := os.Stat(config.Env.K8S.Config.User); err != nil {
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
	groups := map[string]map[string][]string{
		"":                      {"pods": {"*"}, "services": {"*"}, "configmaps": {"*"}, "pods/exec": {"*"}, "nodes": {"get", "list", "watch"}},
		"batch":                 {"jobs": {"*"}},
		"networking.k8s.io":     {"networkpolicies": {"*"}},
		"crd.projectcalico.org": {"ippools": {"*"}},
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
	if _, err := os.Stat(config.Env.K8S.Config.Admin); err != nil {
		log.Logger.Fatalf("Invalid config.k8s.config.admin: %s", err)
	}
	var err error
	adminAPIConfig, err = clientcmd.LoadFromFile(config.Env.K8S.Config.Admin)
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
	calicoClient, err = clientset.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init Calico client: %s", err)
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

func ensureServiceAccount(ctx context.Context) {
	if _, err := kubeClient.CoreV1().ServiceAccounts(namespaceName).Get(ctx, svcAccountName, metav1.GetOptions{}); err != nil {
		log.Logger.Infof("Creating ServiceAccount %s...", svcAccountName)
		_, err = kubeClient.CoreV1().ServiceAccounts(namespaceName).Create(ctx, &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      svcAccountName,
				Namespace: namespaceName,
			},
		}, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Failed to create ServiceAccount: %v", err)
		}
	}
}

func ensureSecret(ctx context.Context) {
	var err error
	if svcAccountToken, err = kubeClient.CoreV1().Secrets(namespaceName).Get(ctx, svcAccountSecretName, metav1.GetOptions{}); err != nil {
		log.Logger.Infof("Creating secret %s...", svcAccountSecretName)
		svcAccountToken, err = kubeClient.CoreV1().Secrets(namespaceName).Create(ctx, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      svcAccountSecretName,
				Namespace: namespaceName,
				Annotations: map[string]string{
					"kubernetes.io/service-account.name": svcAccountName,
				},
			},
			Type: corev1.SecretTypeServiceAccountToken,
		}, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Failed to create secret: %v", err)
		}
	}
}

func ensureIPPool(ctx context.Context) {
	if _, err := calicoClient.ProjectcalicoV3().IPPools().Get(ctx, ipPoolName, metav1.GetOptions{}); err == nil {
		_ = calicoClient.ProjectcalicoV3().IPPools().Delete(ctx, ipPoolName, metav1.DeleteOptions{})
	}
	_, err := calicoClient.ProjectcalicoV3().IPPools().Create(ctx, &projectcalicov3.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: ipPoolName},
		Spec: projectcalicov3.IPPoolSpec{
			CIDR:        config.Env.K8S.IPPool.CIDR,
			IPIPMode:    projectcalicov3.IPIPModeNever,
			NATOutgoing: true,
			BlockSize:   config.Env.K8S.IPPool.BlockSize,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Fatalf("Failed to create IPPool: %v", err)
	}
}

func ensureRole(ctx context.Context) {
	if _, err := kubeClient.RbacV1().Roles(namespaceName).Get(ctx, adminRoleName, metav1.GetOptions{}); err == nil {
		_ = kubeClient.RbacV1().Roles(namespaceName).Delete(ctx, adminRoleName, metav1.DeleteOptions{})
	}
	_, err := kubeClient.RbacV1().Roles(namespaceName).Create(ctx, &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      adminRoleName,
			Namespace: namespaceName,
		},
		Rules: []rbacv1.PolicyRule{
			{APIGroups: []string{""}, Resources: []string{"pods", "services", "configmaps", "pods/exec"}, Verbs: []string{"*"}},
			{APIGroups: []string{"batch"}, Resources: []string{"jobs"}, Verbs: []string{"*"}},
			{APIGroups: []string{"networking.k8s.io"}, Resources: []string{"networkpolicies"}, Verbs: []string{"*"}},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Fatalf("Failed to create Role: %v", err)
	}
}

func ensureRoleBinding(ctx context.Context) {
	if _, err := kubeClient.RbacV1().RoleBindings(namespaceName).Get(ctx, adminRoleBindingName, metav1.GetOptions{}); err == nil {
		_ = kubeClient.RbacV1().RoleBindings(namespaceName).Delete(ctx, adminRoleBindingName, metav1.DeleteOptions{})
	}
	_, err := kubeClient.RbacV1().RoleBindings(namespaceName).Create(ctx, &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      adminRoleBindingName,
			Namespace: namespaceName,
		},
		Subjects: []rbacv1.Subject{{
			Kind: "ServiceAccount", Name: svcAccountName, Namespace: namespaceName,
		}},
		RoleRef: rbacv1.RoleRef{
			Kind: "Role", Name: adminRoleName, APIGroup: "rbac.authorization.k8s.io",
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Fatalf("Failed to create RoleBinding: %v", err)
	}
}

func ensureClusterRole(ctx context.Context) {
	if _, err := kubeClient.RbacV1().ClusterRoles().Get(ctx, adminClusterRoleName, metav1.GetOptions{}); err == nil {
		_ = kubeClient.RbacV1().ClusterRoles().Delete(ctx, adminClusterRoleName, metav1.DeleteOptions{})
	}
	_, err := kubeClient.RbacV1().ClusterRoles().Create(ctx, &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{Name: adminClusterRoleName},
		Rules: []rbacv1.PolicyRule{
			{APIGroups: []string{""}, Resources: []string{"nodes"}, Verbs: []string{"get", "list", "watch"}},
			{APIGroups: []string{"crd.projectcalico.org"}, Resources: []string{"ippools"}, Verbs: []string{"*"}},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Fatalf("Failed to create ClusterRole: %v", err)
	}
}

func ensureClusterRoleBinding(ctx context.Context) {
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Get(ctx, adminClusterRoleBindingName, metav1.GetOptions{}); err == nil {
		_ = kubeClient.RbacV1().ClusterRoleBindings().Delete(ctx, adminClusterRoleBindingName, metav1.DeleteOptions{})
	}
	_, err := kubeClient.RbacV1().ClusterRoleBindings().Create(ctx, &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{Name: adminClusterRoleBindingName},
		Subjects: []rbacv1.Subject{{
			Kind: "ServiceAccount", Name: svcAccountName, Namespace: namespaceName,
		}},
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole", Name: adminClusterRoleName, APIGroup: "rbac.authorization.k8s.io",
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Fatalf("Failed to create ClusterRoleBinding: %v", err)
	}
}

// writeKubeConfig 写入一个低权限的 kubeconfig 文件
func writeKubeConfig() error {
	token := string(svcAccountToken.Data["token"])
	ca := svcAccountToken.Data["ca.crt"]
	host := kubeConfig.Host
	ctx := adminAPIConfig.Contexts[adminAPIConfig.CurrentContext]
	return clientcmd.WriteToFile(api.Config{
		Clusters: map[string]*api.Cluster{
			ctx.Cluster: {
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
				Cluster:   ctx.Cluster,
				AuthInfo:  svcAccountName,
				Namespace: namespaceName,
			},
		},
		CurrentContext: fmt.Sprintf("%s-admin@kubernetes-%s", namespaceName, svcAccountName),
	}, fmt.Sprintf("%s.conf", namespaceName))
}
