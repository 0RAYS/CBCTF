package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"fmt"

	authorizationv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	Client          *kubernetes.Clientset
	NamespaceName   string
	SvcAccountName  string
	RoleName        string
	RoleBindingName string
	Namespace       *corev1.Namespace
	SvcAccount      *corev1.ServiceAccount
	Role            *rbacv1.Role
	RoleBinding     *rbacv1.RoleBinding
	Config          *rest.Config
)

func Init() {
	var err error
	NamespaceName = config.Env.K8S.Namespace
	SvcAccountName = fmt.Sprintf("%s-admin", NamespaceName)
	RoleName = fmt.Sprintf("%s-admin-role", NamespaceName)
	RoleBindingName = fmt.Sprintf("%s-admin-role-binding", NamespaceName)
	Config, err = clientcmd.BuildConfigFromFlags("", config.Env.K8S.Config)
	if err != nil {
		log.Logger.Errorf("Failed to load k8s admin config")
	}
	Config.QPS = 100
	Config.Burst = 200
	Client, err = kubernetes.NewForConfig(Config)
	if err != nil {
		log.Logger.Errorf("Failed to init k8s client")
	}
	if !checkPermission() {
		log.Logger.Fatalf("Failed to check permission")
	}
	initResources()
}

func initResources() {
	var err error
	log.Logger.Infof("Checking resources in namespace %s", NamespaceName)
	if Namespace, err = Client.CoreV1().Namespaces().Get(context.TODO(), NamespaceName, metav1.GetOptions{}); err != nil {
		log.Logger.Infof("Namespace %s not found, creating...", NamespaceName)
		Namespace = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: NamespaceName,
			},
		}
		_, err = Client.CoreV1().Namespaces().Create(context.TODO(), Namespace, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Error creating namespace: %v", err)
		}
	}
	if SvcAccount, err = Client.CoreV1().ServiceAccounts(NamespaceName).Get(context.TODO(), SvcAccountName, metav1.GetOptions{}); err != nil {
		log.Logger.Infof("ServiceAccount %s not found in %s namespace, creating...", SvcAccountName, NamespaceName)
		SvcAccount = &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      SvcAccountName,
				Namespace: NamespaceName,
			},
		}
		SvcAccount, err = Client.CoreV1().ServiceAccounts(NamespaceName).Create(context.TODO(), SvcAccount, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Error creating ServiceAccount: %v", err)
		}
	}
	if Role, err = Client.RbacV1().Roles(NamespaceName).Get(context.TODO(), RoleName, metav1.GetOptions{}); err != nil {
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
		Role, err = Client.RbacV1().Roles(NamespaceName).Create(context.TODO(), Role, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Error creating Role: %v", err)
		}
	}
	if RoleBinding, err = Client.RbacV1().RoleBindings(NamespaceName).Get(context.TODO(), RoleBindingName, metav1.GetOptions{}); err != nil {
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
		RoleBinding, err = Client.RbacV1().RoleBindings(NamespaceName).Create(context.TODO(), RoleBinding, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Error creating RoleBinding: %v", err)
		}
	}
}

func checkPermission() bool {
	log.Logger.Infof("Checking permissions for k8s")
	verbs := []string{"get", "list", "create", "update", "delete"}
	resourceAttributes := &authorizationv1.ResourceAttributes{
		Namespace: NamespaceName,
		Group:     "*",
		Version:   "*",
		Resource:  "*",
	}

	for _, verb := range verbs {
		resourceAttributes.Verb = verb
		accessReview := &authorizationv1.SelfSubjectAccessReview{
			Spec: authorizationv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: resourceAttributes,
			},
		}
		res, err := Client.AuthorizationV1().SelfSubjectAccessReviews().Create(context.TODO(), accessReview, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Errorf("Failed to check permissions for verb %s: %v", verb, err)
		}
		if !res.Status.Allowed {
			log.Logger.Errorf("User does NOT have permission to %s all resources in namespace cbctf.", verb)
			log.Logger.Errorf("Reason: %s", res.Status.Reason)
			log.Logger.Errorf("EvaluationError: %s", res.Status.EvaluationError)
			return false
		}
	}
	return true
}
