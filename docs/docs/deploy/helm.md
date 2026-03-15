---
sidebar_position: 3
---

# Helm 部署

CBCTF 提供官方 Helm Chart，可一键在 Kubernetes 集群上部署完整平台，包含 MySQL、Redis 及所有必要的 RBAC 配置。

## 前置要求

- Helm 3.10+
- 已搭建 Kubernetes 集群（参考 [集群搭建](/docs/deploy/cluster)）
- 支持 `ReadWriteMany` 的 StorageClass（NFS 等）

## 基本操作

### 添加 Chart 仓库

```bash
helm repo add 0rays https://0rays.github.io/CBCTF-charts
helm repo update
```

### 安装

```bash
# 使用默认值快速安装
helm install cbctf 0rays/cbctf \
  --namespace cbctf \
  --create-namespace \
  --set cbctf.host=https://ctf.example.com

# 推荐：下载并自定义 values 文件后安装
helm show values 0rays/cbctf > my-values.yaml
# 编辑 my-values.yaml ...
helm install cbctf 0rays/cbctf \
  --namespace cbctf \
  --create-namespace \
  -f my-values.yaml
```

### 升级

```bash
helm upgrade cbctf 0rays/cbctf --namespace cbctf -f my-values.yaml
```

### 卸载

```bash
helm uninstall cbctf --namespace cbctf
```

:::caution
卸载不会自动删除 PVC。如需清理所有数据：
```bash
kubectl delete pvc --all -n cbctf
```
:::

### 查看初始管理员密码

```bash
kubectl logs -n cbctf deploy/cbctf | grep "Init Admin"
```

---

## Values 完整参数说明

### `image.*` — 镜像配置

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `image.repository` | `ghcr.io/0rays/cbctf` | 镜像仓库地址 |
| `image.tag` | `latest` | 镜像标签，生产环境建议固定版本如 `20250101` |
| `image.pullPolicy` | `IfNotPresent` | 镜像拉取策略 |

### `imageCredentials.*` — 私有仓库凭证

Chart 会自动创建 `docker-registry` 类型的 Secret：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `imageCredentials.registry` | `""` | 仓库地址，如 `ghcr.io` |
| `imageCredentials.username` | `""` | 用户名 |
| `imageCredentials.password` | `""` | 密码或 Token |

### `timezone` — 时区

```yaml
timezone: "Asia/Shanghai"   # 留空使用 UTC
```

### `service.*` — 服务配置

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `service.type` | `ClusterIP` | 服务类型：`ClusterIP` / `NodePort` / `LoadBalancer` |
| `service.port` | `8000` | 服务端口 |

### `ingress.*` — Ingress 配置

```yaml
ingress:
  enabled: false
  className: "nginx"
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
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

### `cbctf.host` — 平台公网地址

平台对外公开 URL，用于邮件链接、OAuth 回调等。**必须配置，不含末尾斜线。**

```yaml
cbctf:
  host: "https://ctf.example.com"
```

### `cbctf.log.*` — 日志

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `cbctf.log.level` | `info` | 日志等级：`DEBUG` / `INFO` / `WARNING` / `ERROR` |
| `cbctf.log.save` | `true` | 是否持久化日志文件 |

### `cbctf.gin.*` — HTTP 服务

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `cbctf.gin.mode` | `release` | Gin 运行模式：`debug` / `release` / `test` |
| `cbctf.gin.host` | `0.0.0.0` | 监听地址 |
| `cbctf.gin.port` | `8000` | 监听端口 |
| `cbctf.gin.upload.max` | `8` | 单文件上传大小限制（MiB） |
| `cbctf.gin.proxies` | `["127.0.0.1", "::1"]` | 可信反向代理 IP 列表（CIDR 支持） |
| `cbctf.gin.cors` | `["http://cbctf.local"]` | 允许跨域的前端地址列表 |
| `cbctf.gin.ratelimit.global` | `100` | 全局速率限制（req/s） |
| `cbctf.gin.ratelimit.whitelist` | `["::1", "127.0.0.1"]` | 不受限速约束的 IP 列表 |
| `cbctf.gin.log.whitelist` | `["/metrics", "/platform/*filepath"]` | 不记录访问日志的路径 |
| `cbctf.gin.jwt.secret` | `""` | JWT 密钥，留空自动生成 |
| `cbctf.gin.jwt.static` | `false` | 是否禁止动态更新 JWT 密钥 |
| `cbctf.gin.metrics.whitelist` | `["127.0.0.1", "::1", ...]` | 允许访问 `/metrics` 的 IP 列表 |

:::info JWT 密钥自动生成
`cbctf.gin.jwt.secret` 留空时，Helm 在首次安装时自动生成随机密钥并存入 Kubernetes Secret，`helm upgrade` 时保持不变。
:::

### `cbctf.asynq.*` — 异步任务队列

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `cbctf.asynq.concurrency` | `50` | 最大并发 worker 数 |
| `cbctf.asynq.log.level` | `warning` | Asynq 内部日志等级 |

### `cbctf.gorm.log.level` — 数据库日志

支持 `INFO` / `WARNING` / `ERROR` / `SILENT`，默认 `silent`。MySQL 连接信息由 Chart 从 `mysql.auth.*` 自动注入。

### `cbctf.k8s.*` — Kubernetes 集成

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `cbctf.k8s.tcpdump` | `nicolaka/netshoot:latest` | 流量捕获辅助容器镜像 |
| `cbctf.k8s.generator_worker` | `2` | 动态附件生成器 Pod 数量倍率（每节点） |
| `cbctf.k8s.frp.on` | `false` | 是否启用 FRP 端口转发 |
| `cbctf.k8s.frp.frpc` | `snowdreamtech/frpc:latest` | frpc 容器镜像 |
| `cbctf.k8s.frp.nginx` | `nginx:latest` | Nginx 反向代理容器镜像 |
| `cbctf.k8s.frp.frps` | `[]` | frps 服务器配置列表 |
| `cbctf.k8s.kubeovnRBAC` | `false` | 是否创建 KubeOVN RBAC 规则（VPC 模式必须设 true） |
| `cbctf.k8s.multusRBAC` | `false` | 是否创建 Multus RBAC 规则（VPC 模式必须设 true） |
| `cbctf.k8s.externalNetwork.enabled` | `false` | 是否配置靶机外部网络访问 |
| `cbctf.k8s.externalNetwork.cidr` | `""` | 宿主机网段，如 `192.168.100.0/24` |
| `cbctf.k8s.externalNetwork.gateway` | `""` | 宿主机网关 |
| `cbctf.k8s.externalNetwork.interface` | `""` | 宿主机网卡名（所有节点须一致） |

### `cbctf.cheat.ip.whitelist` — 作弊检测 IP 白名单

白名单内的 IP 不参与作弊检测，默认包含所有私有地址段。

### `cbctf.registration.*` — 注册配置

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `cbctf.registration.enabled` | `true` | 是否开放公开注册 |
| `cbctf.registration.default_group` | `0` | 新用户默认分组 ID（0 表示不分配） |

### `cbctf.webhook.whitelist` — Webhook 目标白名单

允许的 Webhook 目标 URL 前缀列表，空列表表示不限制。

### `persistence.*` — 持久化存储

CBCTF 需要 `ReadWriteMany` 类型的 PVC 用于附件文件共享。

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `persistence.enabled` | `true` | 是否启用持久化 |
| `persistence.storageClass` | `""` | StorageClass 名称，留空使用集群默认 |
| `persistence.accessMode` | `ReadWriteMany` | 访问模式，**必须为 RWX** |
| `persistence.size` | `20Gi` | 存储容量 |

:::warning StorageClass 要求
K3S 默认的 `local-path` **不支持** `ReadWriteMany`。部署前需配置支持 RWX 的 StorageClass（如 NFS Provisioner）。
:::

### `mysql.*` — MySQL

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `mysql.enabled` | `true` | 是否使用内置 MySQL |
| `mysql.image.repository` | `mysql` | 镜像 |
| `mysql.image.tag` | `8.0` | 版本 |
| `mysql.auth.rootPassword` | `""` | root 密码，留空自动生成 |
| `mysql.auth.database` | `cbctf` | 数据库名 |
| `mysql.auth.username` | `cbctf` | 用户名 |
| `mysql.auth.password` | `""` | 密码，留空自动生成 |
| `mysql.persistence.size` | `5Gi` | 数据卷容量 |
| `mysql.extraConfig` | `utf8mb4 + max_connections=500` | 追加到 `my.cnf` 的 MySQL 配置 |

### `redis.*` — Redis

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `redis.enabled` | `true` | 是否使用内置 Redis |
| `redis.image.tag` | `7-alpine` | 版本 |
| `redis.auth.password` | `""` | 密码，留空自动生成 |
| `redis.persistence.size` | `1Gi` | 数据卷容量 |

### `resources` / `nodeSelector` / `tolerations` / `affinity`

标准 Kubernetes Pod 调度配置，按需配置。

---

## 常见场景示例

### 最小化部署（ClusterIP + port-forward）

```yaml
cbctf:
  host: "http://localhost:8000"
  gin:
    cors:
      - "http://localhost:8000"
    jwt:
      secret: "my-dev-secret"

mysql:
  auth:
    password: "mysql-password"

redis:
  auth:
    password: "redis-password"
```

```bash
helm install cbctf 0rays/cbctf -f values.yaml -n cbctf --create-namespace
kubectl port-forward svc/cbctf 8000:8000 -n cbctf
```

### 生产部署（Ingress + TLS + VPC 网络）

```yaml
cbctf:
  host: "https://ctf.example.com"
  gin:
    cors:
      - "https://ctf.example.com"
    jwt:
      secret: "your-very-long-random-secret-here"
    proxies:
      - "10.244.0.0/16"    # Ingress Controller Pod CIDR
  k8s:
    kubeovnRBAC: true
    multusRBAC: true
    externalNetwork:
      enabled: true
      cidr: "192.168.0.0/24"
      gateway: "192.168.0.1"
      interface: "eth0"

ingress:
  enabled: true
  className: "nginx"
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
  hosts:
    - host: ctf.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: cbctf-tls
      hosts:
        - ctf.example.com

persistence:
  storageClass: "nfs-client"
  size: 50Gi

mysql:
  auth:
    rootPassword: "strong-root-password"
    password: "strong-cbctf-password"
  persistence:
    size: 20Gi

redis:
  auth:
    password: "strong-redis-password"
```
