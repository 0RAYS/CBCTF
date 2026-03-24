---
sidebar_position: 3
---

# Helm 部署

CBCTF 提供官方 Helm Chart，可在 Kubernetes 中部署应用、共享存储接入以及可选的内置 PostgreSQL/Redis。

## 前置要求

- Helm 3.10+
- 可用的 Kubernetes 集群
- 支持 `ReadWriteMany` 的存储类，用于附件与动态附件共享
- 如需 VPC 网络，需提前安装 Kube-OVN 与 Multus

## 基本操作

### 添加仓库

```bash
helm repo add 0rays https://cbctf.0rays.club/CBCTF
helm repo update
```

### 安装

```bash
helm install cbctf 0rays/cbctf \
  --namespace cbctf \
  --create-namespace \
  --set cbctf.host=https://ctf.example.com
```

推荐先导出默认值再修改：

```bash
helm show values 0rays/cbctf > my-values.yaml
helm install cbctf 0rays/cbctf \
  --namespace cbctf \
  --create-namespace \
  -f my-values.yaml
```

### 升级与卸载

```bash
helm upgrade cbctf 0rays/cbctf -n cbctf -f my-values.yaml
helm uninstall cbctf -n cbctf
```

:::caution
卸载不会自动清理 PVC；若需清空数据，请手动删除对应 PVC。
:::

### 查看初始管理员密码

```bash
kubectl logs -n cbctf deployment/cbctf | grep "Init Admin"
```

## 核心 Values

### 镜像与运行环境

| 参数 | 说明 |
|------|------|
| `image.repository` / `image.tag` | CBCTF 应用镜像 |
| `image.pullPolicy` | 镜像拉取策略 |
| `timezone` | 容器时区 |
| `imagePullSecrets` | 额外拉取凭据 |
| `imageCredentials.*` | 自动生成 registry Secret 的内联凭据 |

### 服务暴露

| 参数 | 说明 |
|------|------|
| `service.type` | `ClusterIP` / `NodePort` / `LoadBalancer` |
| `service.port` | 应用服务端口 |
| `ingress.*` | Ingress 主机、路径和 TLS |

### `cbctf.*`

#### `cbctf.host`

平台公开访问地址。必须与最终对外地址一致，并且不要带尾部 `/`。

#### `cbctf.log.*`

- `cbctf.log.level`
- `cbctf.log.save`

#### `cbctf.gin.*`

- `mode`
- `host`
- `port`
- `upload.max`
- `proxies`
- `cors`
- `ratelimit.global`
- `ratelimit.whitelist`
- `log.whitelist`
- `jwt.secret`
- `metrics.whitelist`

Chart 会将 `cbctf.gin.jwt.secret` 写入 Secret，并通过环境变量注入容器。

#### `cbctf.asynq.*`

- `concurrency`
- `log.level`
- `queues.victim`
- `queues.generator`
- `queues.attachment`
- `queues.email`
- `queues.webhook`
- `queues.image`

`cbctf.asynq.concurrency` 保留为全局默认值；当前应用会按任务类型分别启动独立的 Asynq worker 池，实际生效的并发主要由 `cbctf.asynq.queues.*` 控制：

- `queues.victim`：容器靶机启停任务
- `queues.generator`：动态附件生成器启停任务
- `queues.attachment`：附件生成任务
- `queues.email`：邮件任务
- `queues.webhook`：Webhook 任务
- `queues.image`：图片处理任务

#### `cbctf.gorm.log.level`

控制 GORM 日志等级。PostgreSQL 连接信息由 Chart 根据 `postgres.*` 自动拼装。

#### `cbctf.gorm.postgres.sslmode`

PostgreSQL SSL 开关，类型为布尔值：

- `false` -> DSN `sslmode=disable`
- `true` -> DSN `sslmode=require`

#### `cbctf.k8s.*`

- `tcpdump`
- `frp.on`
- `frp.frpc`
- `frp.nginx`
- `frp.frps`
- `externalNetwork.enabled`
- `externalNetwork.cidr`
- `externalNetwork.gateway`
- `externalNetwork.interface`

当前 Chart 和应用代码中 **不存在** `cbctf.k8s.generator_worker`、`cbctf.k8s.kubeovnRBAC`、`cbctf.k8s.multusRBAC` 这些可配置项。

#### 其他配置

- `cbctf.cheat.ip.whitelist`
- `cbctf.registration.enabled`
- `cbctf.registration.default_group`
- `cbctf.webhook.whitelist`

### 持久化

| 参数 | 说明 |
|------|------|
| `persistence.enabled` | 是否启用共享存储 |
| `persistence.storageClass` | RWX 存储类 |
| `persistence.accessMode` | 默认 `ReadWriteMany` |
| `persistence.size` | 共享卷容量 |
| `persistence.existingClaim` | 复用已有 PVC |

动态附件依赖共享卷；若使用不支持 RWX 的存储类，动态附件会失败。

### PostgreSQL / Redis

| 参数 | 说明 |
|------|------|
| `postgres.enabled` | 是否部署内置 PostgreSQL |
| `postgres.externalHost` | 关闭内置 PostgreSQL 时必须提供 |
| `postgres.auth.*` | 用户名、密码、数据库名 |
| `postgres.persistence.*` | PostgreSQL 数据卷 |
| `postgres.extraConfig` | 额外 PostgreSQL 配置 |
| `redis.enabled` | 是否部署内置 Redis |
| `redis.externalHost` | 关闭内置 Redis 时必须提供 |
| `redis.auth.password` | Redis 密码 |
| `redis.persistence.*` | Redis 数据卷 |

## 示例

### 最小部署

```yaml
cbctf:
  host: "http://localhost:8000"
  gin:
    cors:
      - "http://localhost:8000"
    jwt:
      secret: "my-dev-secret"

postgres:
  auth:
    password: "postgres-password"

redis:
  auth:
    password: "redis-password"
```

```bash
helm install cbctf 0rays/cbctf -n cbctf --create-namespace -f values.yaml
kubectl port-forward svc/cbctf 8000:8000 -n cbctf
```

### Ingress + TLS + 外部网络

```yaml
cbctf:
  host: "https://ctf.example.com"
  gin:
    cors:
      - "https://ctf.example.com"
    jwt:
      secret: "your-very-long-random-secret"
    proxies:
      - "10.244.0.0/16"
  k8s:
    externalNetwork:
      enabled: true
      cidr: "192.168.0.0/24"
      gateway: "192.168.0.1"
      interface: "eth0"

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

persistence:
  storageClass: nfs-client
  size: 50Gi
```

## Chart 生成的资源

默认会创建：

- Deployment / Service / Ingress
- ServiceAccount / ClusterRole / ClusterRoleBinding
- 共享 PVC
- 可选的 PostgreSQL 与 Redis StatefulSet
- 可选的外部网络 Subnet 与 NAD

Chart 已内置 Kube-OVN 与 Multus 所需 ClusterRole 规则，无需额外启用开关。
