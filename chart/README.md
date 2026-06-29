# CBCTF

Helm Charts

The chart renders `.Values.cbctf` into `/app/config.yaml`; sensitive settings are stored in the ConfigMap. Change default passwords before production use.

These fields are deployment-only:

- PostgreSQL connection: `gorm.postgres.*`
- Redis connection: `redis.*`
- Data path: `path`
- Gin listen address and port: `gin.host`, `gin.port`

`cbctf.gorm.postgres.sslmode` is a boolean value.

- `false` maps to PostgreSQL DSN `sslmode=disable`
- `true` maps to PostgreSQL DSN `sslmode=require`

## Kubernetes RBAC

The application runs a startup SelfSubjectAccessReview check and exits if required Kubernetes permissions are missing. The chart-created ClusterRole grants the runtime permissions used by the backend:

| API group | Resources | Verbs |
| --- | --- | --- |
| core | `pods` | `create`, `get`, `list`, `watch`, `delete`, `deletecollection` |
| core | `pods/exec` | `create` |
| core | `pods/log` | `get` |
| core | `services` | `create`, `list`, `delete` |
| core | `configmaps` | `create`, `deletecollection` |
| core | `persistentvolumeclaims` | `get` |
| core | `namespaces` | `get` |
| core | `nodes` | `list` |
| `batch` | `jobs` | `create` |
| `networking.k8s.io` | `networkpolicies` | `create`, `deletecollection` |
| `discovery.k8s.io` | `endpointslices` | `deletecollection` |
| `authorization.k8s.io` | `selfsubjectaccessreviews` | `create` |
| `k8s.cni.cncf.io` | `network-attachment-definitions` | `create`, `get`, `deletecollection` |
| `kubevirt.io` | `virtualmachines` | `create`, `get`, `deletecollection` |
| `kubeovn.io` | `subnets` | `create`, `get`, `deletecollection` |
| `kubeovn.io` | `vpcs` | `create`, `deletecollection` |
| `kubeovn.io` | `ips` | `deletecollection` |
