---
title: 配置说明
description: 说明 CBCTF 配置来源、Helm values、配置文件和数据库运行时设置。
---

# 配置说明

CBCTF 从内置默认值和 `config.yaml` 读取配置。Helm 会把 values 渲染成 `/app/config.yaml`，优先级高于内置默认值。

启动参数 `-c <文件>` 用于选择配置文件，默认是 `config.yaml`。当前实现没有绑定 `CBCTF_*` 环境变量，也不自动生成缺失的配置文件。不要把 Kubernetes 注入的环境变量当作应用配置来源。

以下配置只来自部署配置，不进入 `settings` 表：

- PostgreSQL/GORM 配置：`gorm.*`
- Redis 连接信息：`redis.*`
- 数据存储目录：`path`
- Gin 监听地址和端口：`gin.host`、`gin.port`
- 运行时共享 PVC：`k8s.shared_volume_claim`（空值回退到 `{namespace}-shared-volume`）

`internal/db/setting.go` 中列出的运行时设置首次启动时写入 `settings`；之后以数据库值为准。系统配置不会写回 `config.yaml`。所以升级 Helm values 不会覆盖数据库中的同名设置。

旧的 `log.level` / `log.save` 已不在应用配置结构中；当前日志输出由 Logrus 和 Redis 日志 Hook 处理，不能用这两个键控制级别或文件轮转。`gorm.log.level` 与 `asynq.log.level` 仍有效。

## 保存与生效

「系统管理」保存后会读取数据库设置到**当前进程**，但并非所有已初始化组件都会随之重建：

| 设置 | 生效时机 |
| --- | --- |
| 注册开关、默认分组、后续请求读取的公开地址、Webhook/作弊白名单 | 当前进程后续操作使用新值 |
| 抓包开关、worker 镜像、生成器池容量、FRP 镜像等 | 后续新建工作负载使用；已有实例不自动重建 |
| Gin 模式、可信代理、CORS、上传限制、全局限流；Asynq 并发；Kubernetes 命名空间 | 需要后台「重启」重建路由、客户端或 worker |
| 固定部署配置 | 修改文件/values 后重新创建 Pod；后台重启不会重新读取配置文件 |

存在多个平台进程时，保存不会广播到其他进程。不要只修改 `k8s.namespace` 就把它当作迁移：目标命名空间、RoleBinding、PVC 和历史工作负载也需要处理，Helm 部署通常保持 Release namespace。

上传大小限制已拆分为 `gin.upload.picture`、`gin.upload.challenge`、`gin.upload.writeup`。旧的 `gin.upload.max` 不再生效。

## 常用配置项

| 配置项                          | 说明                                         | 示例                          |
|------------------------------|--------------------------------------------|-----------------------------|
| `host`                       | 平台对外访问地址，用于根路径跳转、邮件链接、OAuth 回调等，不要带尾部 `/`  | `https://ctf.example.com`   |
| `path`                       | 数据根目录，存放附件、动态附件、流量文件、GeoIP 数据库等            | `./data`                    |
| `gin.mode`                   | Gin 模式                                     | `release`                   |
| `gin.host`                   | HTTP 监听地址                                  | `0.0.0.0`                   |
| `gin.port`                   | HTTP 监听端口                                  | `8000`                      |
| `gin.proxies`                | 可信代理 IP 或 CIDR                             | `10.244.0.0/16`             |
| `gin.upload.picture`         | 图片上传大小限制，单位 MiB                            | `8`                         |
| `gin.upload.challenge`       | 题目附件上传大小限制，单位 MiB                          | `8`                         |
| `gin.upload.writeup`         | 题解上传大小限制，单位 MiB                            | `8`                         |
| `gin.ratelimit.global`       | 全局每分钟限流                                    | `100`                       |
| `gin.ratelimit.whitelist`    | 跳过限流的 IP 或 CIDR                            | `127.0.0.1`                 |
| `gin.origins`                | 允许的浏览器请求 Origin，用于跨域请求与认证 cookie 判断         | `https://ctf.example.com`   |
| `gin.jwt.secret`             | JWT 签名密钥，生产环境必须替换默认值                       | `change-me-long-random`     |
| `gin.metrics.whitelist`      | 允许访问 `/metrics` 的 IP 或 CIDR                | `10.0.0.0/8`                |
| `gin.pprof.whitelist`        | 允许访问 `/debug/pprof/*` 的 IP 或 CIDR           | `127.0.0.1`                 |
| `gorm.postgres.host`         | PostgreSQL 地址                              | `postgres.cbctf.svc`        |
| `gorm.postgres.port`         | PostgreSQL 端口                              | `5432`                      |
| `gorm.postgres.user`         | PostgreSQL 用户名                             | `cbctf`                     |
| `gorm.postgres.pwd`          | PostgreSQL 密码                              | `example-postgres-password` |
| `gorm.postgres.db`           | PostgreSQL 数据库名                            | `cbctf`                     |
| `gorm.postgres.sslmode`      | 布尔值，`false` 为 `disable`，`true` 为 `require` | `false`                     |
| `redis.host`                 | Redis 地址                                   | `redis.cbctf.svc`           |
| `redis.port`                 | Redis 端口                                   | `6379`                      |
| `redis.pwd`                  | Redis 密码                                   | `example-redis-password`    |
| `registration.enabled`       | 是否允许公开注册                                   | `true`                      |
| `registration.default_group` | 新用户默认分组 ID，`0` 表示不指定                       | `0`                         |

## 异步任务配置

Redis 同时用于缓存和 Asynq 任务队列。以下配置影响后台任务并发：

| 配置项                       | 说明            | 示例        |
|---------------------------|---------------|-----------|
| `asynq.log.level`         | Asynq 日志级别    | `warning` |
| `asynq.queues.victim`     | 靶机启动、停止各自的任务并发 | `8`       |
| `asynq.queues.traffic`    | 靶机流量解析任务并发    | `2`       |
| `asynq.queues.generator`  | 动态附件生成器启停任务并发 | `3`       |
| `asynq.queues.attachment` | 附件 worker 分发并发上限 | `32`      |
| `asynq.queues.email`      | 邮件任务并发        | `10`      |
| `asynq.queues.webhook`    | Webhook 任务并发  | `15`      |
| `asynq.queues.image`      | 图片处理任务并发      | `10`      |

当靶机启动、附件生成或邮件发送堆积时，先检查 Redis 状态和任务日志，再根据资源情况调整对应队列并发。

## Kubernetes 配置

| 配置项             | 说明                                          | 示例                                |
|-----------------|---------------------------------------------|-----------------------------------|
| `k8s.namespace` | 靶机、生成器等资源所在命名空间                             | `cbctf`                           |
| `k8s.shared_volume_claim` | 平台和运行时工作负载共用的 PVC；部署参数，后台不修改 | `cbctf-data` |
| `k8s.capture`   | 流量捕获 sidecar 镜像                             | `ghcr.io/domcyrus/rustnet:latest` |
| `k8s.capture_enabled` | 是否创建抓包容器及相关挂载 | `true` |
| `k8s.priority_class_name` | 普通 Pod 使用的已有 PriorityClass，当前不应用到 VM | `""` |
| `k8s.worker_image` | 含 `/app/worker` 的独立 worker 镜像 | `ghcr.io/0rays/cbctf-worker:latest` |
| `k8s.generator_pool_size` | 每比赛、每动态题的 generator 目标池容量 | `2` |
| `k8s.frp.on`    | 是否启用 FRP 暴露靶机端口                             | `false`                           |
| `k8s.frp.frpc`  | FRP client 镜像                               | `ghcr.io/fatedier/frpc:v0.69.0`   |
| `k8s.frp.nginx` | FRP 辅助 Nginx 镜像                             | `nginx:latest`                    |
| `k8s.frp.frps`  | FRPS 地址、端口、token、端口池                        | `host: frps.example.com`          |

抓包开关、Pod PriorityClass、worker 镜像和生成器池容量均可在「系统管理 → Kubernetes 配置」中查看和修改，保存到数据库后供后续新建实例及补池使用。已有实例不会自动重建或缩容，其他平台副本需重新加载配置。

预热、调度、缓存与回收机制见[工作负载调度](./workloads)。启用 FRP 时使用 ClusterIP，否则使用 NodePort。

## Helm 与配置文件的关系

Helm 会把以下值渲染进 ConfigMap：

| Values                       | 说明       |
|------------------------------|----------|
| `cbctf.gin.jwt.secret`       | JWT 签名密钥 |
| `postgres.auth.password`     | 数据库密码    |
| `redis.auth.password`        | Redis 密码 |

应用不从环境变量读取配置，也不为这些值生成 Secret。Chart 的镜像仓库凭据 Secret 与这些应用配置是不同对象。

## 首次启动行为

- 连接 PostgreSQL，创建 `pg_trgm` 扩展（失败只记录警告）
- 自动迁移数据表
- 初始化运行时策略设置、品牌配置、权限、默认角色、默认分组、Cron 任务和 OAuth 默认项
- 如果管理员组中没有用户，创建 `admin` 用户并将初始密码打印到日志

## 在线配置

管理后台会展示固定部署参数；PostgreSQL/GORM、Redis、数据目录、Gin 监听地址和监听端口只能通过部署配置修改。

GeoLite2-City 数据库可通过管理后台上传，保存到 `{path}/GeoLite2-City.mmdb`。

固定部署参数修改后需要重新创建 Pod。`registration.default_group: 0` 表示不指定注册用户组；应先到 RBAC 页面查看实际分组 ID，再选择普通用户分组，避免新用户没有参赛权限。

## 安全配置建议

:::warning

- 生产环境不要使用默认 `gin.jwt.secret`
- 数据库、Redis、SMTP、OAuth、FRP token 使用强密码或随机密钥
- `host`、Ingress 域名应与用户访问地址一致，`gin.origins` 应包含前端页面的实际 `Origin`
- 部署在反向代理后时正确设置 `gin.proxies`
- `/metrics` 只允许 Prometheus 或可信来源访问
- 只将必要的目标加入 `webhook.whitelist`
  :::
