> For AI agents: the complete documentation index is available at /llms.txt, the full documentation bundle is available at /llms-full.txt.

# 监控与日志

## 系统状态

「仪表盘」展示系统资源、网络收发、平均请求耗时和用户、比赛、提交等统计，查看需要 `admin:system:status` 权限。部分统计暂不可用时，可稍后刷新；持续异常则检查应用日志和对应服务的状态。

## 任务队列

在「任务队列」切换历史记录和实时队列，查看排队、正在执行、重试及已结束的任务，需要 `admin:task:read` 权限。批量操作部分失败时，可结合逐项结果中的原因和任务日志排查。

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

需要进一步定位 CPU 或内存异常时，运维人员可通过 `/debug/pprof/*` 获取性能诊断数据。请只允许可信运维来源访问。

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

管理后台「日志」可在线查看平台近期运行记录；也可通过 Kubernetes 查看应用容器日志。

### 日志级别

后台日志页可按级别筛选。排障时先查看错误和警告，再结合相邻时间的详细记录。数据库日志级别通过 `gorm.log.level` 设置，任务日志级别通过 `asynq.log.level` 设置。

### 日志持久化

后台只保留有限的近期日志。需要赛后审计或长期检索时，请接入集群日志采集和归档服务。

### 在线查看日志

查看后台日志需要 `admin:log:read` 权限，无需登录服务器。加载失败时可使用刷新重新读取。

## 访问日志

访问日志包含请求方法、路径、状态码和耗时，可用于定位失败请求和响应较慢的操作。

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

在「文件」页面查看和管理平台文件，需要 `admin:file:list` 权限，包括：

- 题目附件（`attachment.zip`、`generator.zip`）
- 流量捕获文件（`.pcap`）
- 选手 Writeup（PDF/DOC/DOCX）

支持在线下载（`admin:file:read`）和删除（`admin:file:delete`）。

## 系统运行时设置

在「系统管理」修改并保存运行配置。查看、修改和后台重启分别需要 `admin:system:read`、`admin:system:update`、`admin:system:restart` 权限。

PostgreSQL/GORM、Redis、数据目录、Gin 监听地址和监听端口在页面中只读，只能通过部署配置修改。

共享 PVC `k8s.shared_volume_claim` 同样是部署参数，修改后需要重新创建平台 Pod。后台「重启」用于重新加载在线设置，不会重新读取磁盘配置文件。各项设置的生效方式见[配置说明](/deploy/settings.md)。

:::warning
重启会短暂中断服务。正式比赛期间谨慎使用。
:::
