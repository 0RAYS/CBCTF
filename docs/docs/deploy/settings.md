---
sidebar_position: 5
---

# 配置说明

CBCTF 的核心配置文件为程序根目录下的 `config.yaml`，使用 YAML 格式。

## 完整默认配置

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
    static: false
  metrics:
    whitelist:
      - 127.0.0.1
      - ::1
gorm:
  mysql:
    host: 127.0.0.1
    port: 3306
    user: cbctf
    pwd: password
    db: cbctf
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
          - from: 40000
            to: 60000
            exclude:
              - 50000
  generator_worker: 2
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

---

## 环境变量覆盖

所有配置项均可通过以 `CBCTF_` 为前缀的环境变量覆盖，使用 `_` 作为层级分隔符：

```bash
CBCTF_GIN_PORT=9000              # 覆盖 gin.port
CBCTF_GORM_MYSQL_PWD=newpassword # 覆盖 gorm.mysql.pwd
CBCTF_LOG_LEVEL=DEBUG            # 覆盖 log.level
```

环境变量优先级高于 `config.yaml`。

---

## 字段说明

### `host`

- **类型**：string
- **默认值**：`http://127.0.0.1:8000`
- **说明**：平台对外公开访问 URL，用于 OAuth 回调地址、邮件中的链接、前端跨域校验等。**不含末尾斜线。**

### `path`

- **类型**：string
- **默认值**：`./data`
- **说明**：程序数据存储根目录，需具备完整读写权限。存储内容包括上传的附件、GeoIP 数据库等。

### `log.level`

- **类型**：string
- **默认值**：`info`
- **说明**：日志输出等级，支持 `DEBUG` `INFO` `WARNING` `ERROR`。

### `log.save`

- **类型**：bool
- **默认值**：`true`
- **说明**：是否将日志持久化到 `./logs/` 目录，支持文件轮转。

### `asynq.concurrency`

- **类型**：int
- **默认值**：`50`
- **说明**：Asynq 异步任务队列最大并发 worker 数。过低会导致任务积压，过高会占用过多资源。

### `asynq.log.level`

- **类型**：string
- **默认值**：`warning`
- **说明**：Asynq 内部日志等级，支持 `DEBUG` `INFO` `WARNING` `ERROR`。

### `gin.mode`

- **类型**：string
- **默认值**：`release`
- **说明**：Gin 框架运行模式。`release` 为生产模式，`debug` 输出详细路由日志。

### `gin.host`

- **类型**：string
- **默认值**：`127.0.0.1`
- **说明**：HTTP 服务监听地址。有反向代理时可仅监听 `127.0.0.1`；直接对外暴露时设为 `0.0.0.0`。

### `gin.port`

- **类型**：int
- **默认值**：`8000`
- **说明**：HTTP 服务监听端口。

### `gin.upload.max`

- **类型**：int
- **默认值**：`8`
- **说明**：单次文件上传大小限制，单位 MiB。

### `gin.proxies`

- **类型**：[]string
- **默认值**：`["127.0.0.1"]`
- **说明**：可信反向代理服务器 IP 列表，支持 CIDR 格式。配置后平台才能从 `X-Forwarded-For` 正确获取真实客户端 IP（影响作弊检测、速率限制、IP 日志等）。

### `gin.cors`

- **类型**：[]string
- **默认值**：`["http://127.0.0.1:8000"]`
- **说明**：允许跨域访问的前端地址列表。前后端同域部署时填平台地址，分离部署时填前端域名。

### `gin.ratelimit.global`

- **类型**：int
- **默认值**：`100`
- **说明**：全局请求速率限制（每秒请求数），超出返回 `429 Too Many Requests`。部分敏感接口有独立更严格的限制。

### `gin.ratelimit.whitelist`

- **类型**：[]string
- **默认值**：`["::1", "127.0.0.1"]`
- **说明**：不受速率限制约束的 IP 列表，支持 CIDR 格式。

### `gin.log.whitelist`

- **类型**：[]string
- **默认值**：`["/metrics", "/platform/*filepath"]`
- **说明**：不记录访问日志的请求路径列表。支持 Gin 风格的路由通配符（`*filepath`）。

### `gin.jwt.secret`

- **类型**：string
- **默认值**：`0rays-jbnrz`（示例，**必须修改**）
- **说明**：JWT 签名密钥。泄露将导致认证绕过。建议使用 32 位以上随机字符串。`static: false` 时，平台会通过定时任务定期轮转密钥。

### `gin.jwt.static`

- **类型**：bool
- **默认值**：`false`
- **说明**：设为 `true` 时禁止动态轮转 JWT 密钥（`secret` 值永久有效）。不推荐在生产环境使用。

### `gin.metrics.whitelist`

- **类型**：[]string
- **默认值**：`["127.0.0.1", "::1"]`
- **说明**：允许访问 `/metrics`（Prometheus 指标）端点的 IP 列表，支持 CIDR。未在白名单内的请求返回 `403`。

### `gorm.mysql.host`

- **类型**：string
- **默认值**：`127.0.0.1`

### `gorm.mysql.port`

- **类型**：int
- **默认值**：`3306`

### `gorm.mysql.user`

- **类型**：string
- **默认值**：`cbctf`

### `gorm.mysql.pwd`

- **类型**：string
- **默认值**：`password`（**必须修改**）

### `gorm.mysql.db`

- **类型**：string
- **默认值**：`cbctf`

### `gorm.mysql.mxopen`

- **类型**：int
- **默认值**：`100`
- **说明**：数据库连接池最大连接数。

### `gorm.mysql.mxidle`

- **类型**：int
- **默认值**：`10`
- **说明**：数据库连接池最大空闲连接数。

### `gorm.log.level`

- **类型**：string
- **默认值**：`silent`
- **说明**：GORM SQL 日志等级，支持 `INFO` `WARNING` `ERROR` `SILENT`。

### `redis.host`

- **类型**：string
- **默认值**：`127.0.0.1`

### `redis.port`

- **类型**：int
- **默认值**：`6379`

### `redis.pwd`

- **类型**：string
- **默认值**：`password`
- **说明**：Redis 密码。无密码时留空。

### `k8s.config`

- **类型**：string
- **默认值**：`./admin.yaml`
- **说明**：Kubernetes kubeconfig 文件路径，用于访问 K8s API。

### `k8s.namespace`

- **类型**：string
- **默认值**：`cbctf`
- **说明**：动态题目容器部署的 Kubernetes 命名空间，同时作为相关 K8s 资源名称前缀。

### `k8s.tcpdump`

- **类型**：string
- **默认值**：`nicolaka/netshoot:latest`
- **说明**：流量捕获功能使用的 sidecar 容器镜像。

### `k8s.frp.on`

- **类型**：bool
- **默认值**：`false`
- **说明**：是否启用 FRP 端口转发功能，用于将题目容器端口暴露给选手。

### `k8s.frp.frpc`

- **类型**：string
- **默认值**：`snowdreamtech/frpc:latest`
- **说明**：frpc 客户端容器镜像。

### `k8s.frp.nginx`

- **类型**：string
- **默认值**：`nginx:latest`
- **说明**：Nginx 反向代理容器镜像，用于 HTTP 类题目的端口转发。

### `k8s.frp.frps`

- **类型**：[]object
- **默认值**：`[]`
- **说明**：frps 服务器列表，多个时随机选择。每个条目包含 `host`、`port`、`token` 和 `allowed`（端口范围列表）。

### `k8s.generator_worker`

- **类型**：int
- **默认值**：`2`
- **说明**：每个动态附件题目的生成器 Pod 数量倍率，实际 Pod 数 = `节点数 × generator_worker`。

### `cheat.ip.whitelist`

- **类型**：[]string
- **默认值**：私有地址段
- **说明**：不参与作弊检测（IP 共享检测）的 IP 列表，支持 CIDR。默认包含所有 RFC 1918 私有地址段，防止内网 IP 误报。

### `webhook.whitelist`

- **类型**：[]string
- **默认值**：`[]`
- **说明**：允许的 Webhook 目标 URL 前缀列表。空列表时不限制目标地址；配置后仅允许向匹配前缀的 URL 发送请求。

### `registration.enabled`

- **类型**：bool
- **默认值**：`true`
- **说明**：是否开放公开注册。设为 `false` 后禁止新用户自助注册，管理员仍可手动创建用户。

### `registration.default_group`

- **类型**：int
- **默认值**：`0`
- **说明**：新注册用户自动加入的分组 ID，`0` 表示不分配分组。

### `geocity_db`

- **类型**：string
- **默认值**：`./data/GeoLite2-City.mmdb`
- **说明**：MaxMind GeoLite2-City 数据库文件路径，用于 IP 地理位置解析（事件日志、用户 IP 信息显示）。文件需从 [MaxMind](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) 自行下载。缺少此文件时其他功能不受影响。

---

:::warning 安全建议
- **`gin.jwt.secret`**：必须修改为随机字符串，泄露将导致认证绕过
- **`gin.metrics.whitelist`**：限制为 Prometheus 服务 IP，不要对公网开放 `/metrics`
- **`gorm.mysql.pwd` / `redis.pwd`**：使用强密码，避免使用默认值
:::
