---
sidebar_position: 5
---

# 配置说明

CBCTF 使用根目录 `config.yaml` 作为主配置文件；若文件不存在，程序首次启动时会自动生成默认配置并退出。

## 默认配置

```yaml
host: http://127.0.0.1:8000
path: ./data
log:
  level: info
  save: true
asynq:
  log:
    level: warning
  concurrency: 50
  queues:
    victim: 2
    generator: 3
    attachment: 10
    email: 10
    webhook: 15
    image: 10
gin:
  mode: release
  host: 127.0.0.1
  port: 8000
  proxies:
    - 127.0.0.1
  upload:
    max: 8
  ratelimit:
    global: 100
    whitelist:
      - ::1
      - 127.0.0.1
  cors:
    - http://127.0.0.1:8000
  log:
    whitelist:
      - /metrics
      - /platform/*filepath
  jwt:
    secret: 0rays-jbnrz
  metrics:
    whitelist:
      - 127.0.0.1
      - ::1
gorm:
  postgres:
    host: 127.0.0.1
    port: 5432
    user: cbctf
    pwd: password
    db: cbctf
    sslmode: false
    mxopen: 100
    mxidle: 10
  log:
    level: silent
redis:
  host: 127.0.0.1
  port: 6379
  pwd: password
k8s:
  config: ./admin.yaml
  namespace: cbctf
  tcpdump: nicolaka/netshoot:latest
  frp:
    on: false
    frpc: snowdreamtech/frpc:latest
    nginx: nginx:latest
    frps:
      - host: example.com
        port: 7000
        token: token
        allowed:
          - from: 10000
            to: 30000
            exclude:
              - 20000
cheat:
  ip:
    whitelist:
      - 127.0.0.1
      - ::1
      - 10.0.0.0/8
      - 192.168.0.0/16
      - 172.16.0.0/12
      - 100.64.0.0/10
webhook:
  whitelist: []
registration:
  enabled: true
  default_group: 0
geocity_db: ./data/GeoLite2-City.mmdb
```

## 环境变量覆盖

所有配置项都可通过 `CBCTF_` 前缀环境变量覆盖，例如：

```bash
CBCTF_GIN_PORT=9000
CBCTF_GORM_POSTGRES_PWD=newpassword
CBCTF_GORM_POSTGRES_SSLMODE=true
CBCTF_LOG_LEVEL=DEBUG
```

`.` 会被 `_` 替换，因此 `gorm.postgres.pwd` 对应 `CBCTF_GORM_POSTGRES_PWD`。

## 关键字段

### `host`

平台对外访问地址，用于 OAuth 回调、邮件链接和静态资源跳转，不要带尾部 `/`。

### `path`

数据根目录。题目附件、动态附件、流量文件和 GeoIP 数据库都依赖这里的读写权限。

### `log.level` / `log.save`

- `level` 支持 `DEBUG`、`INFO`、`WARNING`、`ERROR`
- `save: true` 时日志写入 `./logs/%Y%m%d.log`

### `asynq.*`

异步任务队列配置，负责邮件、Webhook、动态附件生成等后台任务。

- `asynq.log.level`：Asynq worker 日志级别
- `asynq.concurrency`：全局默认并发值
- `asynq.queues.victim`：容器靶机启停任务并发
- `asynq.queues.generator`：动态附件生成器启停任务并发
- `asynq.queues.attachment`：附件生成任务并发
- `asynq.queues.email`：邮件任务并发
- `asynq.queues.webhook`：Webhook 任务并发
- `asynq.queues.image`：图片处理任务并发

当前版本会为不同任务类型启动独立的 Asynq worker 池，因此限流应优先调整 `asynq.queues.*`。`asynq.concurrency` 保留为默认值和兼容字段。

### `gin.*`

- `gin.host` / `gin.port`：HTTP 监听地址
- `gin.proxies`：可信代理 IP 或 CIDR
- `gin.upload.max`：上传大小限制，单位 MiB
- `gin.ratelimit.*`：全局限流与白名单
- `gin.cors`：允许跨域的前端地址
- `gin.log.whitelist`：不记录访问日志的路由
- `gin.jwt.secret`：JWT 签名密钥，必须替换默认值
- `gin.metrics.whitelist`：允许访问 `/metrics` 的 IP 白名单

当前代码中 **不存在** `gin.jwt.static` 配置项。

### `gorm.postgres.*`

PostgreSQL 连接信息与连接池参数：

- `host`、`port`、`user`、`pwd`、`db`
- `mxopen`：最大连接数
- `mxidle`：最大空闲连接数

#### `gorm.postgres.sslmode`

布尔开关：

- `false` -> DSN `sslmode=disable`
- `true` -> DSN `sslmode=require`

### `redis.*`

Redis 既用作缓存，也用作 Asynq 的任务队列后端。

### `k8s.*`

- `k8s.config`：kubeconfig 路径；容器内 Helm 部署时会写成 `/admin/admin.yaml`
- `k8s.namespace`：题目相关资源所在命名空间
- `k8s.tcpdump`：流量捕获 sidecar 镜像
- `k8s.frp.on`：是否启用 FRP 暴露题目端口
- `k8s.frp.frpc` / `k8s.frp.nginx`：FRP 相关镜像
- `k8s.frp.frps`：FRPS 服务端配置列表

当前代码中 **不存在** `k8s.generator_worker` 配置项。

### `cheat.ip.whitelist`

IP 白名单；命中的地址不会参与 IP 相关作弊检测。

### `webhook.whitelist`

Webhook 允许访问的目标白名单，支持 IP、CIDR、主机名和 `host:port`。

### `registration.*`

- `registration.enabled`：是否允许公开注册
- `registration.default_group`：注册后自动加入的分组 ID

### `geocity_db`

GeoLite2-City 数据库路径。配置后，可在后台查看 IP 的地理位置。

## 启动时的 K8s 资源检查

程序启动后会检查以下 K8s 资源：

- 命名空间 `{namespace}`
- PVC `{namespace}-shared-volume`
- Subnet `{namespace}-external-network`
- `kube-system/{namespace}-external-network` 对应的 NAD

含义如下：

- 命名空间缺失：启动失败
- PVC 缺失：动态附件不可用
- 外部网络资源缺失：VPC 网络模式不可用

## 在线配置更新

管理后台的系统配置页对应以下接口：

- `GET /admin/system/config`
- `PUT /admin/system/config`
- `POST /admin/system/restart`

`POST /admin/system/restart` 会向当前进程发送 `SIGUSR1`，触发热重启。
