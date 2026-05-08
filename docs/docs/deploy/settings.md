---
sidebar_position: 5
---

# 配置说明

CBCTF 使用 `config.yaml` 作为主配置文件。源码运行时，如果文件不存在，程序会从内置默认配置生成 `config.yaml` 并退出；修改后再次启动即可。Helm 部署时，Chart 会把 values 渲染成容器内 `/app/config.yaml`，并通过环境变量注入敏感配置。

环境变量支持 `CBCTF_` 前缀覆盖配置，点号转换为下划线。例如 `gorm.postgres.pwd` 对应 `CBCTF_GORM_POSTGRES_PWD`。

## 常用配置项

| 配置项 | 说明 | 示例 |
|---|---|---|
| `host` | 平台对外访问地址，用于根路径跳转、邮件链接、OAuth 回调等，不要带尾部 `/` | `https://ctf.example.com` |
| `path` | 数据根目录，存放附件、动态附件、流量文件、GeoIP 数据库等 | `./data` |
| `log.level` | 应用日志级别 | `info` |
| `log.save` | 是否将日志写入文件 | `false` |
| `gin.mode` | Gin 模式 | `release` |
| `gin.host` | HTTP 监听地址 | `0.0.0.0` |
| `gin.port` | HTTP 监听端口 | `8000` |
| `gin.proxies` | 可信代理 IP 或 CIDR | `10.244.0.0/16` |
| `gin.upload.max` | 上传大小限制，单位 MiB | `8` |
| `gin.ratelimit.global` | 全局每分钟限流 | `100` |
| `gin.ratelimit.whitelist` | 跳过限流的 IP 或 CIDR | `127.0.0.1` |
| `gin.cors` | 允许跨域来源 | `https://ctf.example.com` |
| `gin.jwt.secret` | JWT 签名密钥，生产环境必须替换默认值 | `change-me-long-random` |
| `gin.metrics.whitelist` | 允许访问 `/metrics` 的 IP 或 CIDR | `10.0.0.0/8` |
| `gorm.postgres.host` | PostgreSQL 地址 | `postgres.cbctf.svc` |
| `gorm.postgres.port` | PostgreSQL 端口 | `5432` |
| `gorm.postgres.user` | PostgreSQL 用户名 | `cbctf` |
| `gorm.postgres.pwd` | PostgreSQL 密码 | `example-postgres-password` |
| `gorm.postgres.db` | PostgreSQL 数据库名 | `cbctf` |
| `gorm.postgres.sslmode` | 布尔值，`false` 为 `disable`，`true` 为 `require` | `false` |
| `redis.host` | Redis 地址 | `redis.cbctf.svc` |
| `redis.port` | Redis 端口 | `6379` |
| `redis.pwd` | Redis 密码 | `example-redis-password` |
| `registration.enabled` | 是否允许公开注册 | `true` |
| `registration.default_group` | 新用户默认分组 ID，`0` 表示不指定 | `0` |
| `geocity_db` | GeoLite2-City 数据库路径 | `./data/GeoLite2-City.mmdb` |

## 异步任务配置

Redis 同时用于缓存和 Asynq 任务队列。以下配置影响后台任务并发：

| 配置项 | 说明 | 示例 |
|---|---|---|
| `asynq.log.level` | Asynq 日志级别 | `warning` |
| `asynq.concurrency` | 全局默认并发值 | `50` |
| `asynq.queues.victim` | 靶机启停任务并发 | `2` |
| `asynq.queues.generator` | 动态附件生成器启停任务并发 | `3` |
| `asynq.queues.attachment` | 附件生成任务并发 | `10` |
| `asynq.queues.email` | 邮件任务并发 | `10` |
| `asynq.queues.webhook` | Webhook 任务并发 | `15` |
| `asynq.queues.image` | 图片处理任务并发 | `10` |

当靶机启动、附件生成或邮件发送堆积时，先检查 Redis 状态和任务日志，再根据资源情况调整对应队列并发。

## Kubernetes 配置

| 配置项 | 说明 | 示例 |
|---|---|---|
| `k8s.config` | kubeconfig 路径；Helm 部署时为 `/admin/admin.yaml` | `./admin.yaml` |
| `k8s.namespace` | 靶机、生成器等资源所在命名空间 | `cbctf` |
| `k8s.tcpdump` | 流量捕获 sidecar 镜像 | `nicolaka/netshoot:latest` |
| `k8s.frp.on` | 是否启用 FRP 暴露靶机端口 | `false` |
| `k8s.frp.frpc` | FRP client 镜像 | `snowdreamtech/frpc:latest` |
| `k8s.frp.nginx` | FRP 辅助 Nginx 镜像 | `nginx:latest` |
| `k8s.frp.frps` | FRPS 地址、端口、token、端口池 | `host: frps.example.com` |
| `k8s.external_networks.enabled` | 是否初始化外部网络资源 | `true` |
| `k8s.external_networks.interfaces[].interface` | 外部网卡名，也是节点标签值 | `ens192` |
| `k8s.external_networks.interfaces[].cidr` | 外部网络 CIDR | `192.168.0.0/24` |
| `k8s.external_networks.interfaces[].gateway` | 外部网络网关 | `192.168.0.1` |

动态靶机和动态附件要求应用能够访问 Kubernetes API，并拥有创建、查询、删除相关资源的权限。Helm Chart 会创建默认 ClusterRole；如果自行部署，需要按 Chart 中的权限范围配置 RBAC。

## Helm 与配置文件的关系

Helm Chart 会把以下敏感值放入 Secret，并用环境变量覆盖配置文件中的空值：

| 环境变量 | 来源 | 说明 |
|---|---|---|
| `CBCTF_GIN_JWT_SECRET` | 应用 Secret | JWT 签名密钥 |
| `CBCTF_GORM_POSTGRES_PWD` | PostgreSQL Secret | 数据库密码 |
| `CBCTF_REDIS_PWD` | Redis Secret | Redis 密码 |

如果 `cbctf.gin.jwt.secret`、`postgres.auth.password` 或 `redis.auth.password` 留空，Chart 会在首次安装时生成随机值，并在升级时复用已有 Secret。

## 首次启动行为

- 连接 PostgreSQL，创建 `pg_trgm` 扩展（失败只记录警告）。
- 自动迁移数据表。
- 初始化系统设置、品牌配置、权限、默认角色、默认分组、Cron 任务和 OAuth 默认项。
- 如果管理员组中没有用户，创建 `admin` 用户并将初始密码打印到日志。

## 在线配置

管理后台提供系统配置页面，可查看和更新部分运行配置，也可以触发服务重启。生产环境仍建议将关键部署参数固化在 Helm values、Secret 或配置管理系统中，避免 Pod 重建后配置漂移。

## 安全配置建议

- 生产环境不要使用默认 `gin.jwt.secret`。
- 数据库、Redis、SMTP、OAuth、FRP token 使用强密码或随机密钥。
- `host`、`gin.cors`、Ingress 域名应保持一致。
- 部署在反向代理后时正确设置 `gin.proxies`。
- `/metrics` 只允许 Prometheus 或可信来源访问。
- 只将必要的目标加入 `webhook.whitelist`。
