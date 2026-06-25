# Helm 部署

CBCTF Chart 位于仓库根目录的 `chart/`。默认会创建应用 Deployment、Service、Ingress、ServiceAccount、ClusterRole、共享 PVC，以及内置
PostgreSQL 和 Redis。

## 前置要求

- Kubernetes 集群可用
- Helm 可用
- 集群可以拉取 `ghcr.io/0rays/cbctf`、PostgreSQL、Redis 以及题目镜像
- 如启用持久化，集群需要可用 StorageClass
- 动态附件建议使用支持 `ReadWriteMany` 的共享存储
- VPC 靶机需要提前安装 Kube-OVN 和 Multus CNI
- KubeVirt VM 靶机需要提前安装 KubeVirt，并确认节点支持虚拟化

## 安装

添加 Helm Repo

```bash
helm repo add cbctf https://cbctf.0rays.club
```

**须使用自定义 values**

Helm values 会渲染为容器内 `/app/config.yaml`。升级 values 后会滚动重启：

```bash
helm show values cbctf/cbctf > values.yaml
helm install cbctf cbctf/cbctf -n cbctf --create-namespace -f values.yaml
```

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

| 配置项                                | 说明                          | 示例                        |
| ---------------------------------- | --------------------------- | ------------------------- |
| `cbctf.host`                       | 平台公开访问地址，不要带尾部 `/`          | `https://ctf.example.com` |
| `cbctf.log.level`                  | 应用日志级别                      | `info`                    |
| `cbctf.log.save`                   | 是否持久化日志                     | `false`                   |
| `cbctf.gin.mode`                   | Gin 运行模式                    | `release`                 |
| `cbctf.gin.host`                   | 容器内监听地址                     | `0.0.0.0`                 |
| `cbctf.gin.port`                   | 容器内监听端口                     | `8000`                    |
| `cbctf.gin.upload.picture`         | 图片上传大小限制，单位 MiB             | `8`                       |
| `cbctf.gin.upload.challenge`       | 题目附件上传大小限制，单位 MiB           | `8`                       |
| `cbctf.gin.upload.writeup`         | 题解上传大小限制，单位 MiB             | `8`                       |
| `cbctf.gin.proxies`                | 可信代理 IP 或 CIDR              | `10.244.0.0/16`           |
| `cbctf.gin.cors`                   | CORS 允许来源                   | `https://ctf.example.com` |
| `cbctf.gin.ratelimit.global`       | 全局限流                        | `100`                     |
| `cbctf.gin.jwt.secret`             | JWT 签名密钥                    | `change-me-long-random`   |
| `cbctf.gin.metrics.whitelist`      | 允许访问 `/metrics` 的 IP 或 CIDR | `10.0.0.0/8`              |
| `cbctf.asynq.queues.traffic`       | 靶机流量解析任务并发                  | `2`                       |
| `cbctf.registration.enabled`       | 是否允许公开注册                    | `true`                    |
| `cbctf.registration.default_group` | 新用户默认分组 ID，`0` 表示不指定        | `0`                       |
| `cbctf.cheat.ip.whitelist`         | 作弊检测 IP 白名单                 | `10.0.0.0/8`              |
| `cbctf.webhook.whitelist`          | Webhook 目标白名单               | `example.com`             |

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

:::info
当前 Chart values 中没有外部 PostgreSQL 或外部 Redis 的 `externalHost` 配置项。如果需要使用外部数据库，需要同步调整 Chart
模板或用等价的 Service 名称接入。
:::

## Kubernetes 靶机配置

| 配置项                     | 说明                    | 示例                                |
| ----------------------- | --------------------- | --------------------------------- |
| `serviceAccount.create` | 是否创建应用 ServiceAccount | `true`                            |
| `cbctf.k8s.capture`     | 流量捕获镜像                | `ghcr.io/domcyrus/rustnet:latest` |
| `cbctf.k8s.frp.on`      | 是否启用 FRP 端口暴露         | `false`                           |
| `cbctf.k8s.frp.frpc`    | FRP client 镜像         | `ghcr.io/fatedier/frpc:v0.69.0`   |
| `cbctf.k8s.frp.nginx`   | FRP 转发辅助 Nginx 镜像     | `nginx:latest`                    |
| `cbctf.k8s.frp.frps`    | FRPS 服务端、token 和端口池   | `host: frps.example.com`          |

Chart 创建的 ClusterRole 包含 Pod、Service、Job、NetworkPolicy、EndpointSlice、Multus NAD、KubeVirt VirtualMachine、Kube-OVN
Subnet/VPC/IP 等资源权限。Chart 不会安装 KubeVirt、Kube-OVN 或 Multus，需要时请先在集群层面安装这些组件。

## Ingress 示例

```yaml
cbctf:
  host: "https://ctf.example.com"
  gin:
    cors:
      - "https://ctf.example.com"
    proxies:
      - "10.244.0.0/16"
    jwt:
      secret: "change-me-long-random-secret"

ingress:
  enabled: true
  className: nginx
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

## 启动时资源检查与创建

Helm 安装后，应用启动时会检查或创建以下资源：

- 命名空间：`{namespace}`
- 共享存储 PVC：`{namespace}-shared-volume`

:::warning
PVC 缺失会导致动态附件不可用。KubeVirt 资源不会在启动时创建，只有启动包含 `x-kubevirt: true` 的 VM 靶机时才会创建对应
`VirtualMachine`。
:::
