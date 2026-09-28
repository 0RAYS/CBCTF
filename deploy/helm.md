> For AI agents: the complete documentation index is available at /llms.txt, the full documentation bundle is available at /llms-full.txt.

# Helm 部署

CBCTF Chart 位于仓库根目录的 `chart/`。默认会创建应用 Deployment、Service、Ingress、ServiceAccount、ClusterRole、共享 PVC，以及内置
PostgreSQL 和 Redis。

## 前置要求

- Kubernetes 集群可用
- Helm 可用
- 集群可以拉取 `ghcr.io/0rays/cbctf`、PostgreSQL、Redis 以及题目镜像
- 如启用持久化，集群需要可用 StorageClass
- 动态附件输入与抓包文件需要平台和运行时 Pod 共享同一 PVC；多节点场景使用 `ReadWriteMany`
- VPC 靶机需要提前安装 Kube-OVN 和 Multus CNI
- KubeVirt VM 靶机需要提前安装 KubeVirt，并确认节点支持虚拟化

## 安装

添加 Helm Repo

```bash
helm repo add cbctf https://cbctf.0rays.club
helm repo update
```

**须使用自定义 values**

Helm values 会渲染为容器内 `/app/config.yaml`。应用 Deployment 固定为单副本、`Recreate` 策略；ConfigMap 变化触发 Pod 替换，会有短暂中断，并非滚动更新。

```bash
helm show values cbctf/cbctf > values.yaml
helm upgrade --install cbctf cbctf/cbctf -n cbctf --create-namespace -f values.yaml --wait --timeout 10m
```

:::info
本文与仓库 Chart `0.0.29` 对应。使用仓库代码时可将命令中的 `cbctf/cbctf` 换成 `./chart`；新配置能力需要匹配版本的后端镜像。

运行时设置首次写入数据库后，以数据库为准。修改 values 中的 `host`、JWT、注册、上传限制或 Kubernetes 运行参数再升级，并不会覆盖已经保存的设置；请到「系统管理」修改并按需重启。参见[配置来源与生效时机](/deploy/settings.md)。
:::

## 升级和卸载

```bash
helm upgrade cbctf cbctf/cbctf -n cbctf -f values.yaml
helm uninstall cbctf -n cbctf
```

共享 PVC 默认带有保留策略，卸载不会删除 `/app/data` 中的数据。PostgreSQL 和 Redis 的 PVC 也应在确认备份后再手动清理。

## 常用 Values

| 配置项                         | 说明                          | 示例                    |
| --------------------------- | --------------------------- | --------------------- |
| `image.repository`          | 应用镜像仓库                      | `ghcr.io/0rays/cbctf` |
| `image.tag`                 | 应用镜像标签                      | `latest`              |
| `imagePullSecrets`          | 私有镜像拉取 Secret               | `[{name: regcred}]`   |
| `imageCredentials.*`        | Chart 自动创建镜像仓库 Secret 的内联凭据 | `registry: ghcr.io`   |
| `timezone`                  | 容器时区                        | `Asia/Shanghai`       |
| `service.type`              | Service 类型                  | `ClusterIP`           |
| `service.port`              | Service 端口                  | `8000`                |
| `ingress.enabled`           | 是否启用 Ingress                | `true`                |
| `ingress.className`         | IngressClass                | `nginx`               |
| `ingress.hosts`             | 域名和路径                       | `ctf.example.com`     |
| `ingress.tls`               | TLS Secret 配置               | `cbctf-tls`           |
| `resources`                 | 应用 Pod 资源限制                 | `requests.cpu: 500m`  |
| `persistence.enabled`       | 是否创建共享 PVC                  | `true`                |
| `persistence.storageClass`  | 共享 PVC 的 StorageClass       | `nfs-client`          |
| `persistence.accessMode`    | 访问模式                        | `ReadWriteMany`       |
| `persistence.size`          | 共享 PVC 容量                   | `20Gi`                |
| `persistence.existingClaim` | 复用已有 PVC                    | `cbctf-data`          |

## 应用配置

上传大小限制已拆分为 `cbctf.gin.upload.picture`、`cbctf.gin.upload.challenge`、`cbctf.gin.upload.writeup`。旧的 `cbctf.gin.upload.max` 不再生效。

| 配置项                                | 说明                                | 示例                        |
| ---------------------------------- | --------------------------------- | ------------------------- |
| `cbctf.host`                       | 平台公开访问地址，不要带尾部 `/`                | `https://ctf.example.com` |
| `cbctf.gin.mode`                   | Gin 运行模式                          | `release`                 |
| `cbctf.gin.host`                   | 容器内监听地址                           | `0.0.0.0`                 |
| `cbctf.gin.port`                   | 容器内监听端口                           | `8000`                    |
| `cbctf.gin.upload.picture`         | 图片上传大小限制，单位 MiB                   | `8`                       |
| `cbctf.gin.upload.challenge`       | 题目附件上传大小限制，单位 MiB                 | `8`                       |
| `cbctf.gin.upload.writeup`         | 题解上传大小限制，单位 MiB                   | `8`                       |
| `cbctf.gin.proxies`                | 可信代理 IP 或 CIDR                    | `10.244.0.0/16`           |
| `cbctf.gin.origins`                | 允许的浏览器请求 Origin                   | `https://ctf.example.com` |
| `cbctf.gin.ratelimit.global`       | 全局限流                              | `100`                     |
| `cbctf.gin.jwt.secret`             | JWT 签名密钥                          | `change-me-long-random`   |
| `cbctf.gin.metrics.whitelist`      | 允许访问 `/metrics` 的 IP 或 CIDR       | `10.0.0.0/8`              |
| `cbctf.gin.pprof.whitelist`        | 允许访问 `/debug/pprof/*` 的 IP 或 CIDR | `127.0.0.1`               |
| `cbctf.asynq.queues.traffic`       | 靶机流量解析任务并发                        | `2`                       |
| `cbctf.registration.enabled`       | 是否允许公开注册                          | `true`                    |
| `cbctf.registration.default_group` | 新用户默认分组 ID，`0` 表示不指定              | `0`                       |
| `cbctf.cheat.ip.whitelist`         | 作弊检测 IP 白名单                       | `10.0.0.0/8`              |
| `cbctf.webhook.whitelist`          | Webhook 目标白名单                     | `example.com`             |

JWT、PostgreSQL 和 Redis 密钥会写入 `/app/config.yaml`。

管理后台不能修改 PostgreSQL/GORM、Redis、数据目录、Gin 监听地址和监听端口。

## PostgreSQL 和 Redis

| 配置项                            | 说明                        | 示例                          |
| ------------------------------ | ------------------------- | --------------------------- |
| `postgres.enabled`             | 是否部署内置 PostgreSQL         | `true`                      |
| `postgres.auth.database`       | 数据库名                      | `cbctf`                     |
| `postgres.auth.username`       | 用户名                       | `cbctf`                     |
| `postgres.auth.password`       | PostgreSQL 密码             | `example-postgres-password` |
| `postgres.persistence.enabled` | PostgreSQL 数据持久化          | `true`                      |
| `postgres.persistence.size`    | PostgreSQL PVC 容量         | `5Gi`                       |
| `postgres.extraConfig`         | 追加到 `postgresql.conf` 的配置 | `max_connections = 500`     |
| `redis.enabled`                | 是否部署内置 Redis              | `true`                      |
| `redis.auth.password`          | Redis 密码                  | `example-redis-password`    |
| `redis.persistence.enabled`    | Redis 数据持久化               | `true`                      |
| `redis.persistence.size`       | Redis PVC 容量              | `1Gi`                       |

关闭内置组件时必须设置 `externalHost`，地址按原值使用，不追加命名空间。认证和端口仍使用对应的 `auth` 与 `service.port`：

```yaml
postgres:
  enabled: false
  externalHost: postgres.example.com
  service:
    port: 5432
  auth:
    database: cbctf
    username: cbctf
    password: replace-with-database-password
redis:
  enabled: false
  externalHost: redis.example.com
  service:
    port: 6379
  auth:
    password: replace-with-redis-password
cbctf:
  gorm:
    postgres:
      sslmode: true
```

PostgreSQL 使用 `pg_trgm`；初始化会尝试创建扩展。平台有 HTTP、任务查询和 advisory-lock 连接池，外部 PostgreSQL 应提供足够连接数，并使用直连或 session pooling，不能使用 transaction pooling 承载会话级 advisory lock。Redis 配置目前只有 host、port、密码，不提供 Sentinel、Cluster 或 TLS 参数。

内置 PostgreSQL 的 `auth` 变量仅初始化空数据目录；对已有 PVC 修改 values 中的密码不会执行数据库 `ALTER ROLE`，应先协调修改数据库凭据。关闭任一数据库的 `persistence.enabled` 后使用 `emptyDir`，Pod 重建会丢失对应数据。

## 共享数据卷与镜像凭据

- 平台挂载 `/app/data`，默认 PVC 为 `{namespace}-shared-volume`。`persistence.existingClaim` 可以指定同命名空间已有 PVC，Chart 同时写入后端 `k8s.shared_volume_claim`，使生成器和抓包容器使用同一卷。
- `persistence.enabled: false` 仅使平台数据使用临时卷，不会为运行时 Pod 创建共享 PVC；这不适用于需要动态附件或抓包的部署。
- `imagePullSecrets` 与 `imageCredentials` 用于 Chart 管理的平台、PostgreSQL 和 Redis Pod。后端创建的题目、生成器、FRPC 和预热 Job 不会继承应用 Pod 的凭据；这些 Pod 使用运行命名空间的 `default` ServiceAccount，可将所需 Secret 配置到该 ServiceAccount。VM 的 containerDisk 私有镜像还需按 KubeVirt 的拉取凭据机制单独验证。
- Chart 的 `nodeSelector`、`tolerations`、`affinity`、`resources` 作用于平台 Pod，不等于题目调度和资源配置；题目资源来自 Compose。

## Kubernetes 靶机配置

| 配置项                           | 说明                                 | 示例                                  |
| ----------------------------- | ---------------------------------- | ----------------------------------- |
| `serviceAccount.create`       | 是否创建应用 ServiceAccount              | `true`                              |
| `cbctf.k8s.capture`           | 流量捕获镜像                             | `ghcr.io/domcyrus/rustnet:latest`   |
| `cbctf.k8s.captureEnabled`    | 是否开启抓包容器                           | `true`                              |
| `cbctf.k8s.priorityClassName` | 普通 Pod 的已有 PriorityClass，当前不应用到 VM | `""`                                |
| `cbctf.k8s.workerImage`       | 独立 worker 镜像地址，与主程序版本分开配置          | `ghcr.io/0rays/cbctf-worker:latest` |
| `cbctf.k8s.generatorPoolSize` | 每比赛、每动态题的生成器池容量                    | `2`                                 |
| `cbctf.k8s.frp.on`            | 是否启用 FRP 端口暴露                      | `false`                             |
| `cbctf.k8s.frp.frpc`          | FRP client 镜像                      | `ghcr.io/fatedier/frpc:v0.69.0`     |
| `cbctf.k8s.frp.nginx`         | FRP 转发辅助 Nginx 镜像                  | `nginx:latest`                      |
| `cbctf.k8s.frp.frps`          | FRPS 服务端、token 和端口池                | `host: frps.example.com`            |

Chart 使用命名空间 Role 管理 Pod、Service、Job、NetworkPolicy、ConfigMap、Multus NAD 和 VM，使用 ClusterRole 管理节点读取及 Kube-OVN 集群级资源。Chart 不会安装 KubeVirt、Kube-OVN 或 Multus，需要时请先在集群层面安装这些组件。

| API group              | Resources                        | Verbs                                                          | 用途                           |
| ---------------------- | -------------------------------- | -------------------------------------------------------------- | ---------------------------- |
| core                   | `pods`                           | `create`, `get`, `list`, `watch`, `delete`, `deletecollection` | 创建靶机、生成器、FRPC Pod，并等待状态和清理   |
| core                   | `pods/exec`                      | `create`                                                       | Pod 终端操作；附件 worker 不使用 Exec  |
| core                   | `pods/log`                       | `get`                                                          | 读取 Pod 日志                    |
| core                   | `services`                       | `create`, `list`, `delete`                                     | 创建 ClusterIP / NodePort 暴露   |
| core                   | `configmaps`                     | `create`, `get`, `list`, `watch`, `delete`, `deletecollection` | 文件配置、共享根对象缓存与 foreground GC  |
| core                   | `persistentvolumeclaims`         | `get`                                                          | 启动时检查共享 PVC                  |
| core                   | `namespaces`                     | `get`                                                          | 启动时检查靶机命名空间                  |
| core                   | `nodes`                          | `list`                                                         | 枚举节点镜像和预拉取目标节点               |
| `batch`                | `jobs`                           | `create`                                                       | 创建镜像预拉取 Job                  |
| `networking.k8s.io`    | `networkpolicies`                | `create`, `deletecollection`                                   | 创建和清理靶机网络策略                  |
| `discovery.k8s.io`     | `endpointslices`                 | `deletecollection`                                             | 清理 Service 产生的 EndpointSlice |
| `authorization.k8s.io` | `selfsubjectaccessreviews`       | `create`                                                       | 启动时执行权限自检                    |
| `k8s.cni.cncf.io`      | `network-attachment-definitions` | `create`, `get`, `deletecollection`                            | VPC 模式下创建和清理 Multus NAD      |
| `kubevirt.io`          | `virtualmachines`                | `create`, `get`, `list`, `watch`, `deletecollection`           | VM 创建、共享就绪缓存及清理              |
| `kubeovn.io`           | `subnets`                        | `create`, `get`, `deletecollection`                            | VPC 模式下创建和清理 Kube-OVN 子网     |
| `kubeovn.io`           | `vpcs`                           | `create`, `get`, `list`, `watch`, `delete`, `deletecollection` | VPC 资源树及级联删除确认               |
| `kubeovn.io`           | `ips`                            | `deletecollection`                                             | 清理 Kube-OVN IP 分配            |

Chart 不会安装 KubeVirt、Kube-OVN 或 Multus，需要时请先在集群层面安装这些组件。

如果使用自定义 ServiceAccount 或外部 RBAC，可以用下面的命令提前检查关键权限：

```bash
kubectl auth can-i create pods -n cbctf --as=system:serviceaccount:cbctf:cbctf
kubectl auth can-i watch pods -n cbctf --as=system:serviceaccount:cbctf:cbctf
kubectl auth can-i create selfsubjectaccessreviews.authorization.k8s.io --as=system:serviceaccount:cbctf:cbctf
kubectl auth can-i create network-attachment-definitions.k8s.cni.cncf.io -n cbctf --as=system:serviceaccount:cbctf:cbctf
kubectl auth can-i create virtualmachines.kubevirt.io -n cbctf --as=system:serviceaccount:cbctf:cbctf
kubectl auth can-i create subnets.kubeovn.io --as=system:serviceaccount:cbctf:cbctf
```

## Ingress 示例

```yaml
cbctf:
  host: "https://ctf.example.com"
  gin:
    origins:
      - "https://ctf.example.com"
    proxies:
      - "10.244.0.0/16"
    pprof:
      whitelist:
        - "127.0.0.1"
        - "10.244.0.0/16"
    jwt:
      secret: "change-me-long-random-secret"

postgres:
  auth:
    password: "replace-with-database-password"
redis:
  auth:
    password: "replace-with-redis-password"
persistence:
  storageClass: nfs-client

ingress:
  enabled: true
  className: nginx
  annotations:
    nginx.ingress.kubernetes.io/proxy-body-size: "16m"
  hosts:
    - host: ctf.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: cbctf-tls
      hosts:
        - ctf.example.com
```

## 安装后检查

```bash
kubectl get pods -n cbctf
kubectl logs -n cbctf deployment/cbctf
kubectl get pvc -n cbctf
kubectl get ingress -n cbctf
```

检查初始管理员密码：

```bash
kubectl logs -n cbctf deployment/cbctf | grep "Init Admin"
```

如果 Pod 反复重启，优先检查日志中的数据库、Redis、RBAC、PVC、KubeVirt、Kube-OVN/Multus 相关错误。

## 启动时资源检查

Helm 安装后，应用启动时只检查以下资源，不自动创建命名空间或 PVC：

- 命名空间：`{namespace}`
- 共享存储 PVC：`k8s.shared_volume_claim`，未指定时为 `{namespace}-shared-volume`
- Kubernetes API 权限：上方 RBAC 表中的所有 verbs

:::warning
命名空间或必需 RBAC 缺失会使启动失败；PVC 缺失会记录警告，随后动态附件和启用抓包的靶机可能无法启动。KubeVirt 资源不会在启动时创建，只有启动包含 `x-kubevirt: true` 的 VM 靶机时才会创建对应
`VirtualMachine`。
:::
