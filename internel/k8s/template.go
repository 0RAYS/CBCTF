package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"bytes"
	"fmt"
	"text/template"
)

const TemplateYAML = `
# Namespace
apiVersion: v1
kind: Namespace
metadata:
    name: {{ .namespaceName }}
---
# ServiceAccount
apiVersion: v1
kind: ServiceAccount
metadata:
    name: {{ .svcAccountName }}
    namespace: {{ .namespaceName }}
---
# ServiceAccountToken Secret
apiVersion: v1
kind: Secret
metadata:
    name: {{ .svcAccountSecretName }}
    namespace: {{ .namespaceName }}
    annotations:
        kubernetes.io/service-account.name: {{ .svcAccountName }}
type: kubernetes.io/service-account-token
---
# Calico IPPool
apiVersion: crd.projectcalico.org/v1
kind: IPPool
metadata:
    name: {{ .ipPoolName }}
spec:
    cidr: {{ .ipPoolCIDR }}
    ipipMode: Never
    natOutgoing: true
    blockSize: {{ ipPoolBlockSize }}
---
# Role
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
    name: {{ .adminRoleName }}
    namespace: {{ .namespaceName }}
rules:
    - apiGroups: [""]
      resources: ["pods", "services", "configmaps", "pods/exec"]
      verbs: ["*"]
    - apiGroups: ["networking.k8s.io"]
      resources: ["networkpolicies"]
      verbs: ["*"]
---
# RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
    name: {{ .adminRoleBindingName }}
    namespace: {{ .namespaceName }}
subjects:
    - kind: ServiceAccount
      name: {{ .svcAccountName }}
      namespace: {{ .namespaceName }}
roleRef:
    kind: Role
    name: {{ .adminRoleName }}
    apiGroup: rbac.authorization.k8s.io
---
# ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    name: {{ .adminClusterRoleName }}
rules:
    - apiGroups: ["crd.projectcalico.org"]
      resources: ["ippools"]
      verbs: ["*"]
---
# ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
    name: {{ .adminClusterRoleBindingName }}
subjects:
    - kind: ServiceAccount
      name: {{ .svcAccountName }}
      namespace: {{ .namespaceName }}
roleRef:
    kind: ClusterRole
    name: {{ .adminClusterRoleName }}
    apiGroup: rbac.authorization.k8s.io
`

func gen(namespaceName string) []byte {
	tmpl := template.Must(template.New("example").Parse(TemplateYAML))
	data := map[string]string{
		"namespaceName":               namespaceName,
		"svcAccountName":              fmt.Sprintf("%s-admin", namespaceName),
		"svcAccountSecretName":        fmt.Sprintf("%s-admin-token", namespaceName),
		"ipPoolName":                  fmt.Sprintf("%s-ippool", namespaceName),
		"adminRoleName":               fmt.Sprintf("%s-admin-role", namespaceName),
		"adminRoleBindingName":        fmt.Sprintf("%s-admin-rolebinding", namespaceName),
		"adminClusterRoleName":        fmt.Sprintf("%s-admin-clusterrole", namespaceName),
		"adminClusterRoleBindingName": fmt.Sprintf("%s-admin-clusterrolebinding", namespaceName),
		"ipPoolCIDR":                  config.Env.K8S.IPPool.CIDR,
		"ipPoolBlockSize":             fmt.Sprintf("%d", config.Env.K8S.IPPool.BlockSize),
	}
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		log.Logger.Fatalf("Failed to generate template yaml: %s", err)
	}
	return buf.Bytes()
}
