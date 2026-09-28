> For AI agents: the complete documentation index is available at /llms.txt, the full documentation bundle is available at /llms-full.txt.

# 监控与日志

## 系统状态

「仪表盘」通过 `GET /admin/system/status`（`admin:system:status`）展示系统资源指标、网络收发、平均请求耗时和用户/比赛/提交等统计。它不是逐项探测 PostgreSQL、Redis、Asynq 连通性的健康检查接口；连接问题还需查看应用日志和各组件状态。

## 任务队列

通过 `GET /admin/tasks` 和 `GET /admin/tasks/live`（`admin:task:read`）查看 Asynq 后台任务的执行状态，包括各队列积压量、正在处理的任务、历史执行记录等。

## 定时任务

在后台「定时任务」页面可查看任务的执行次数和最近成功、失败时间，调整运行间隔，以及设置是否在服务器启动时立即执行一次。启动执行设置仅在下一次启动时生效，保存配置不会立即运行任务。详见[定时任务](/guide/features/cronjobs.md)。

## Prometheus 指标

平台在 `/metrics` 端点暴露 Prometheus 格式的指标数据。

### 访问控制

通过 `gin.metrics.whitelist` 配置允许访问的 IP：

```yaml
gin:
  metrics:
    whitelist:
      - 127.0.0.1
      - 10.0.0.0/8    # Prometheus 所在网段
```

未在白名单内的 IP 访问 `/metrics` 会返回 `403 Forbidden`。

## PProf 诊断

平台注册 Go `net/http/pprof` 诊断端点，路径为 `/debug/pprof/*`。这些端点可能暴露运行时栈、堆和性能信息，生产环境应只允许可信运维来源访问。

通过 `gin.pprof.whitelist` 配置允许访问的 IP 或 CIDR：

```yaml
gin:
  pprof:
    whitelist:
      - 127.0.0.1
      - 10.0.0.0/8    # 运维或跳板机网段
```

未在白名单内的 IP 访问 `/debug/pprof/*` 会返回 `403 Forbidden`。

### 指标内容

- HTTP 请求数和延迟（按路由、方法、状态码分类）
- Asynq 任务队列深度和处理速率
- 作弊检测计数
- 活跃靶机数量

## 日志系统

平台使用 Logrus 输出进程日志，并通过 Redis 日志 Hook 提供后台查询。

### 日志级别

当前进程日志初始化为 Debug 级别，后台日志页可按级别筛选。数据库日志使用 `gorm.log.level`，任务日志使用 `asynq.log.level`；旧的顶层 `log.level` 已不生效。

### 日志持久化

当前没有 `log.save` 文件轮转实现。长期留存容器日志请接入集群日志采集；后台 Redis 日志是有限缓存，不应视为永久归档。

### 在线查看日志

通过 `GET /admin/logs`（`admin:log:read`）在管理后台查看日志内容，无需登录服务器。

## 访问日志

Gin 框架记录每个 HTTP 请求的访问日志，包含请求方法、路径、状态码、耗时等。

### 排除噪声路径

通过 `gin.log.whitelist` 排除不需要记录的路径：

```yaml
gin:
  log:
    whitelist:
      - /metrics
      - /platform/*filepath
```

## 文件管理

通过 `GET /admin/files`（`admin:file:list`）查看和管理平台上所有已上传的文件，包括：

- 题目附件（`attachment.zip`、`generator.zip`）
- 流量捕获文件（`.pcap`）
- 选手 Writeup（PDF/DOC/DOCX）

支持在线下载（`admin:file:read`）和删除（`admin:file:delete`）。

## 系统运行时设置

可通过接口更新系统配置，并在需要时重启重载：

```bash
# 读取当前配置
GET /admin/system/config    # 需 admin:system:read

# 更新配置
PUT /admin/system/config    # 需 admin:system:update

# 重载配置
POST /admin/system/restart  # 需 admin:system:restart
```

PostgreSQL/GORM、Redis、数据目录、Gin 监听地址和监听端口在页面中只读，只能通过部署配置修改。

共享 PVC `k8s.shared_volume_claim` 同样是部署参数。后台重启会停止并重建 HTTP、任务、Cron 和 Kubernetes 客户端，不重新读取磁盘配置文件；部署参数修改需要重建 Pod。上传大小、CORS、代理与并发等初始化参数的生效方式见[配置说明](/deploy/settings.md)。

:::warning
重启会短暂中断服务。正式比赛期间谨慎使用。
:::
