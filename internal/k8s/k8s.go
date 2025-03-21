package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"fmt"
	"os"
	"time"

	authorizationv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	Client                 *kubernetes.Clientset
	Namespace              *corev1.Namespace
	SvcAccount             *corev1.ServiceAccount
	Secret                 *corev1.Secret
	Role                   *rbacv1.Role
	RoleBinding            *rbacv1.RoleBinding
	ClusterRole            *rbacv1.ClusterRole
	ClusterRoleBinding     *rbacv1.ClusterRoleBinding
	Config                 *rest.Config
	APIConfig              *api.Config
	NamespaceName          string
	SvcAccountName         string
	SecretName             string
	RoleName               string
	RoleBindingName        string
	ClusterRoleName        string
	ClusterRoleBindingName string
)

func Init() {
	var err error
	NamespaceName = config.Env.K8S.Namespace
	SvcAccountName = fmt.Sprintf("%s-admin", NamespaceName)
	SecretName = fmt.Sprintf("%s-admin-secret", NamespaceName)
	RoleName = fmt.Sprintf("%s-admin-role", NamespaceName)
	RoleBindingName = fmt.Sprintf("%s-admin-role-binding", NamespaceName)
	ClusterRoleName = fmt.Sprintf("%s-admin-cluster-role", NamespaceName)
	ClusterRoleBindingName = fmt.Sprintf("%s-admin-cluster-role-binding", NamespaceName)
	if _, err = os.Stat(config.Env.K8S.Config.User); config.Env.K8S.Config.User != "" && err == nil {
		APIConfig, err = clientcmd.LoadFromFile(config.Env.K8S.Config.User)
		if err != nil {
			log.Logger.Fatalf("Failed to load k8s user config: %s", err)
		}
	} else {
		APIConfig, err = clientcmd.LoadFromFile(config.Env.K8S.Config.Admin)
		if err != nil {
			log.Logger.Fatalf("Failed to load k8s admin config: %s", err)
		}
	}
	Config, err = clientcmd.NewNonInteractiveClientConfig(*APIConfig, APIConfig.CurrentContext, &clientcmd.ConfigOverrides{}, nil).ClientConfig()
	Config.QPS = 20
	Config.Burst = 40
	log.Logger.Info("K8S config loaded, initiating client...")
	Client, err = kubernetes.NewForConfig(Config)
	if err != nil {
		log.Logger.Fatalf("Failed to init k8s client: %s", err)
	}
	if !checkPermission() {
		os.Exit(-1)
	}
	initResources()
	if _, err = os.Stat(config.Env.K8S.Config.User); config.Env.K8S.Config.User == "" && err != nil {
		if writeKubeConfig() != nil {
			log.Logger.Fatalf("Failed to save kubeconfig to %s.conf: %s ", NamespaceName, err)
		}
		config.Env.K8S.Config.User = fmt.Sprintf("%s.conf", NamespaceName)
		tmp := config.Env.K8S.Config.Admin
		if err := config.Save(config.Env); err != nil {
			log.Logger.Fatalf("Failed to update config: %s", err)
		}
		log.Logger.Infof("Kubeconfig saved to %s.conf, please restart and remove the %s", NamespaceName, tmp)
		os.Exit(0)
	}
	log.Logger.Info("K8S client initialized")
}

// initResources initializes resources in the namespace
func initResources() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	log.Logger.Debugf("Checking resources in namespace %s", NamespaceName)
	if Namespace, err = Client.CoreV1().Namespaces().Get(ctx, NamespaceName, metav1.GetOptions{}); err != nil {
		log.Logger.Infof("Namespace %s not found, creating...", NamespaceName)
		Namespace = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: NamespaceName,
			},
		}
		_, err = Client.CoreV1().Namespaces().Create(ctx, Namespace, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Error creating namespace: %v", err)
		}
	}
	if SvcAccount, err = Client.CoreV1().ServiceAccounts(NamespaceName).Get(ctx, SvcAccountName, metav1.GetOptions{}); err != nil {
		log.Logger.Infof("ServiceAccount %s not found in %s namespace, creating...", SvcAccountName, NamespaceName)
		SvcAccount = &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      SvcAccountName,
				Namespace: NamespaceName,
			},
		}
		SvcAccount, err = Client.CoreV1().ServiceAccounts(NamespaceName).Create(ctx, SvcAccount, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Error creating ServiceAccount: %v", err)
		}
	}
	if Secret, err = Client.CoreV1().Secrets(NamespaceName).Get(ctx, SecretName, metav1.GetOptions{}); err != nil {
		log.Logger.Infof("Secret %s not found in %s namespace, creating...", RoleName, NamespaceName)
		Secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      SecretName,
				Namespace: NamespaceName,
				Annotations: map[string]string{
					"kubernetes.io/service-account.name": SvcAccountName,
				},
			},
			Type: corev1.SecretTypeServiceAccountToken,
		}
		Secret, err = Client.CoreV1().Secrets(NamespaceName).Create(context.TODO(), Secret, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Error creating Secret: %v", err)
		}
	}
	if Role, err = Client.RbacV1().Roles(NamespaceName).Get(ctx, RoleName, metav1.GetOptions{}); err != nil {
		log.Logger.Infof("Role %s not found in %s namespace, creating...", RoleName, NamespaceName)
		Role = &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      RoleName,
				Namespace: NamespaceName,
			},
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{"*"},
					Resources: []string{"*"},
					Verbs:     []string{"*"},
				},
			},
		}
		Role, err = Client.RbacV1().Roles(NamespaceName).Create(ctx, Role, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Error creating Role: %v", err)
		}
	}
	if RoleBinding, err = Client.RbacV1().RoleBindings(NamespaceName).Get(ctx, RoleBindingName, metav1.GetOptions{}); err != nil {
		log.Logger.Infof("RoleBinding %s not found in %s namespace, creating...", RoleBindingName, NamespaceName)
		RoleBinding = &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      RoleBindingName,
				Namespace: NamespaceName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      SvcAccountName,
					Namespace: NamespaceName,
				},
			},
			RoleRef: rbacv1.RoleRef{
				Kind:     "Role",
				Name:     RoleName,
				APIGroup: "rbac.authorization.k8s.io",
			},
		}
		RoleBinding, err = Client.RbacV1().RoleBindings(NamespaceName).Create(ctx, RoleBinding, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Error creating RoleBinding: %v", err)
		}
	}
	if ClusterRole, err = Client.RbacV1().ClusterRoles().Get(ctx, ClusterRoleName, metav1.GetOptions{}); err != nil {
		log.Logger.Warningf("Failed to get ClusterRole: %s", err)
		log.Logger.Infof("ClusterRole %s not found, creating...", ClusterRoleName)
		ClusterRole = &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: ClusterRoleName,
			},
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{"crd.projectcalico.org"},
					Resources: []string{"ippools"},
					Verbs:     []string{"*"},
				},
			},
		}
		ClusterRole, err = Client.RbacV1().ClusterRoles().Create(ctx, ClusterRole, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Error creating ClusterRole: %v", err)
		}
	}
	if ClusterRoleBinding, err = Client.RbacV1().ClusterRoleBindings().Get(ctx, ClusterRoleBindingName, metav1.GetOptions{}); err != nil {
		log.Logger.Warningf("Failed to get ClusterRole: %s", err)
		log.Logger.Infof("ClusterRoleBindings %s not found, creating...", ClusterRoleBindingName)
		ClusterRoleBinding = &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: ClusterRoleBindingName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      SvcAccountName,
					Namespace: NamespaceName,
				},
			},
			RoleRef: rbacv1.RoleRef{
				Kind:     "ClusterRole",
				Name:     ClusterRoleName,
				APIGroup: "rbac.authorization.k8s.io",
			},
		}
		ClusterRoleBinding, err = Client.RbacV1().ClusterRoleBindings().Create(ctx, ClusterRoleBinding, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Error creating ClusterRoleBinding: %v", err)
		}
	}
}

// checkPermission checks if the user has permission to access the resources
func checkPermission() bool {
	log.Logger.Debugf("Checking permission in namespace %s", NamespaceName)
	groups := map[string][]string{
		"*":                     {"*"},
		"crd.projectcalico.org": {"ippools"},
	}
	verbs := []string{"get", "list", "create", "update", "delete"}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	for group, resources := range groups {
		for _, resource := range resources {
			for _, verb := range verbs {
				accessReview := &authorizationv1.SelfSubjectAccessReview{
					Spec: authorizationv1.SelfSubjectAccessReviewSpec{
						ResourceAttributes: &authorizationv1.ResourceAttributes{
							Namespace: NamespaceName,
							Group:     group,
							Version:   "*",
							Resource:  resource,
							Verb:      verb,
						},
					},
				}
				res, err := Client.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, accessReview, metav1.CreateOptions{})
				if err != nil {
					log.Logger.Warningf("Failed to check permissions for verb %s: %v", verb, err)
					return false
				}
				if !res.Status.Allowed {
					log.Logger.Warningf("User does NOT have permission to %s all resources in namespace cbctf.", verb)
					log.Logger.Warningf("Reason: %s", res.Status.Reason)
					log.Logger.Warningf("EvaluationError: %s", res.Status.EvaluationError)
					return false
				}
			}
		}
	}
	return true
}

func writeKubeConfig() error {
	token := string(Secret.Data["token"])
	ca := Secret.Data["ca.crt"]
	host := Config.Host
	ctx := APIConfig.Contexts[APIConfig.CurrentContext]
	kubeConfig := api.Config{
		Clusters: map[string]*api.Cluster{
			ctx.Cluster: {
				Server:                   host,
				CertificateAuthorityData: ca,
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			SvcAccountName: {
				Token: token,
			},
		},
		Contexts: map[string]*api.Context{
			fmt.Sprintf("kubernetes-admin@kubernetes-%s", SvcAccountName): {
				Cluster:   ctx.Cluster,
				AuthInfo:  SvcAccountName,
				Namespace: NamespaceName,
			},
		},
		CurrentContext: fmt.Sprintf("kubernetes-admin@kubernetes-%s", SvcAccountName),
	}
	return clientcmd.WriteToFile(kubeConfig, fmt.Sprintf("%s.conf", NamespaceName))
}
