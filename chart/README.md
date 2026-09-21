# CBCTF

Helm Charts

The chart renders `.Values.cbctf` into `/app/config.yaml`; sensitive settings are stored in the ConfigMap. Change default passwords before production use.

These fields are deployment-only:

- PostgreSQL connection: `gorm.postgres.*`
- Redis connection: `redis.*`
- Data path: `path`
- Gin listen address and port: `gin.host`, `gin.port`

`cbctf.gin.pprof.whitelist` controls the IP/CIDR sources allowed to access `/debug/pprof/*`. Keep it limited to loopback or trusted operator networks in production.

`cbctf.gorm.postgres.sslmode` is a boolean value.

- `false` maps to PostgreSQL DSN `sslmode=disable`
- `true` maps to PostgreSQL DSN `sslmode=require`

## Kubernetes RBAC

The application runs a startup SelfSubjectAccessReview check and exits if required Kubernetes permissions are missing. Namespaced runtime permissions are granted by a **Role + RoleBinding only in the
release namespace**. The ClusterRole/ClusterRoleBinding retain only namespace GET
(restricted to the release namespace), node discovery, self-access reviews and
Kube-OVN cluster-scoped resources. No wildcard permissions are used. The backend's
self-check uses separate resource/subresource fields for `pods/exec` and `pods/log`.

The combined runtime permissions are:

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

## Scheduling rollout notes (0.0.24)

- Upgrade the chart RBAC and backend together. Older binaries self-check subresources
  incorrectly and cannot use the namespace-name-restricted GET rule reliably.
- Cluster-scoped Kube-OVN deletecollection cannot be confined by label in RBAC. Use a
  dedicated platform namespace/cluster and keep the application SA out of challenge
  Pods; the backend also validates deletion selectors.
- `startupProbe` defaults to a five-minute initialization allowance. Override it for
  large DB migrations; liveness/readiness are gated until startup succeeds.
- `terminationGracePeriodSeconds` defaults to 120. HTTP drains for at most 30 seconds,
  and task workers have a 30-second shutdown allowance. Running cron/DB operations
  can still extend shutdown; this is not a hard bound on the entire process.
- Start/stop workers use a separate PostgreSQL advisory-lock pool. Budget additional
  connections (up to the same configured maximum as the task query pool). Direct or
  session-pooled PostgreSQL is required; transaction-mode PgBouncer is unsupported.
- A process crash can still leave pending resources; inspect the admin diagnostics and
  reconcile explicitly. This release does not implement a durable DB/queue outbox.
