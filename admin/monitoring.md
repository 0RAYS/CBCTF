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

### 指标内容

平台自身指标统一使用 `cbctf_` 前缀，PostgreSQL 与连接池指标沿用 `postgres_`、`gorm_` 前缀。仓库内的 `grafana/CBCTF-Dashboard.json` 是配套的 Grafana 仪表盘，可直接导入使用。

#### 服务状态

| 指标                                                                                                                                     | 说明                                                               |
| -------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------- |
| `cbctf_build_info{version,commit,go_version}`                                                                                          | 构建版本信息，`version`/`commit` 由镜像构建时注入                               |
| `cbctf_uptime_seconds`                                                                                                                 | 进程运行时长                                                           |
| `cbctf_dependency_ready{dependency}`                                                                                                   | 依赖就绪状态，`dependency` 取 `postgres`、`valkey`、`kubernetes`，值为 1 表示可用 |
| `cbctf_ctf_collection_success{section}`、`postgres_collection_success`、`cbctf_redis_collection_success`、`cbctf_task_collection_success` | 各采集器最近一次采集是否成功，值为 0 时对应面板数据已过期                                   |

#### HTTP 请求

| 指标                                                              | 说明                                                                                                                   |
| --------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| `cbctf_http_requests_total{method,path,code,status_class}`      | 请求计数。`code` 是平台业务响应码（业务码与 HTTP 语义对齐），但响应行始终是 HTTP 200，不能当作 HTTP 状态码使用；`status_class` 是由 `code` 派生的 `2xx`/`4xx`/`5xx` |
| `cbctf_http_request_duration_seconds{method,path,status_class}` | 请求处理耗时分布                                                                                                             |
| `cbctf_http_response_size_bytes{method,path,status_class}`      | 响应体大小分布                                                                                                              |
| `cbctf_http_request_size_bytes{method,path}`                    | 请求体大小分布                                                                                                              |
| `cbctf_http_in_flight_requests`                                 | 当前正在处理的请求数                                                                                                           |
| `cbctf_http_panics_total{path}`                                 | 被中间件捕获的 panic 次数，任何增长都需要排查                                                                                           |

#### CTF 业务

| 指标                                                                         | 说明                                            |
| -------------------------------------------------------------------------- | --------------------------------------------- |
| `cbctf_flag_submissions_total{contest_id,challenge_type,status}`           | Flag 提交计数                                     |
| `cbctf_flag_solve_duration_seconds{contest_id,challenge_type}`             | 从比赛开始到解出该题的耗时，反映赛题难度                          |
| `cbctf_blood_total{contest_id,blood_order}`                                | 一血/二血/三血事件计数                                  |
| `cbctf_contest_teams{contest_id}`、`cbctf_contest_participants{contest_id}` | 各比赛的战队数与参赛人数                                  |
| `cbctf_contests{status}`                                                   | 按状态（`coming`/`running`/`paused`/`over`）统计比赛数量 |
| `cbctf_challenges{contest_id,challenge_type}`                              | 各比赛的题目数量                                      |
| `cbctf_victims{status}`、`cbctf_victims_active`                             | 靶机容器状态分布与活跃数量                                 |
| `cbctf_generators{status}`                                                 | 生成器池状态分布                                      |
| `cbctf_cheat_detections_total{contest_id,reason_type}`                     | 作弊检测事件计数                                      |
| `cbctf_cheat_reviews_total{contest_id,action}`                             | 管理员复核作弊线索的动作计数                                |
| `cbctf_cheats_pending`                                                     | 待复核的作弊线索数量                                    |

#### 靶机生命周期与 Kubernetes

| 指标                                                | 说明                                                  |
| ------------------------------------------------- | --------------------------------------------------- |
| `cbctf_victim_transitions_total{action,result}`   | 靶机状态跃迁计数，`action` 取 `start`/`stop`/`expire`/`ready` |
| `cbctf_victim_lifecycle_duration_seconds{result}` | 从创建到容器就绪的耗时，即调度与镜像拉取的真实延迟                           |
| `cbctf_victim_readiness_failures_total{reason}`   | 靶机始终未能就绪的数量与原因                                      |
| `cbctf_k8s_operation_errors_total{operation}`     | 失败的 Kubernetes 操作计数                                 |

#### 用户、文件、邮件与安全

| 指标                                                                                                               | 说明              |
| ---------------------------------------------------------------------------------------------------------------- | --------------- |
| `cbctf_user_registrations_total{oauth_provider}`                                                                 | 注册计数            |
| `cbctf_user_logins_total{oauth_provider,status}`                                                                 | 登录尝试计数          |
| `cbctf_password_resets_total{status}`                                                                            | 密码重置计数          |
| `cbctf_file_uploads_total{file_type}`、`cbctf_file_upload_size_bytes{file_type}`                                  | 文件上传计数与大小分布     |
| `cbctf_email_sent_total{email_kind,status}`                                                                      | 邮件发送计数          |
| `cbctf_webhook_deliveries_total{event_type,status}`、`cbctf_webhook_delivery_duration_seconds{event_type,status}` | Webhook 投递计数与耗时 |
| `cbctf_rate_limit_hits_total{endpoint}`                                                                          | 限流命中计数          |

#### 定时任务与异步任务

| 指标                                                                                           | 说明                        |
| -------------------------------------------------------------------------------------------- | ------------------------- |
| `cbctf_cron_job_runs_total{job_name,status}`、`cbctf_cron_job_duration_seconds{job_name}`     | 定时任务运行计数与耗时               |
| `cbctf_task_enqueued_total{task_type,status}`、`cbctf_task_processed_total{task_type,status}` | 异步任务入队与处理计数               |
| `cbctf_task_processing_duration_seconds{task_type,status}`                                   | 异步任务处理耗时                  |
| `cbctf_task_retries_total{task_type}`                                                        | 异步任务重试计数                  |
| `cbctf_task_queue_depth{queue,state}`                                                        | 队列瞬时深度，按任务状态拆分            |
| `cbctf_task_queue_latency_seconds{queue}`                                                    | 队列中最老待处理任务的等待时间，积压告警的主要依据 |
| `cbctf_task_queue_memory_bytes{queue}`                                                       | 队列占用的内存                   |

#### 依赖组件

| 指标                                                                                          | 说明                                              |
| ------------------------------------------------------------------------------------------- | ----------------------------------------------- |
| `gorm_dbstats_*`                                                                            | GORM 连接池的连接数、等待次数与等待耗时                          |
| `cbctf_redis_pool_*`                                                                        | Valkey/Redis 连接池连接数、空闲数、待处理请求、命中与未命中、等待时长       |
| `cbctf_redis_command_duration_seconds{command}`、`cbctf_redis_command_errors_total{command}` | Redis 命令耗时与错误计数                                 |
| `postgres_stat_database_*`                                                                  | 事务提交/回滚、死锁与冲突、块命中/读取及块读取耗时、临时文件数与写入量等库级统计       |
| `postgres_bgwriter_*`                                                                       | Checkpoint 次数与读写耗时，checkpoint、后台写入器与后端直接写出的缓冲区数 |
| `postgres_server_info{version}`                                                             | PostgreSQL 版本信息                                 |
| `postgres_stat_activity_backends{datname,state}`                                            | 各数据库的连接状态分布                                     |
| `postgres_n_live_tup`、`postgres_n_dead_tup`、`postgres_n_mod_since_last_analyze`             | 表的存活行、死元组与未 analyze 的变更量                        |
| `postgres_*_scan_total`、`postgres_*_blks_*_total`                                           | 表的扫描次数与块命中/读取计数                                 |
| `postgres_last_*_timestamp_seconds`                                                         | 表最近一次 vacuum / analyze 时间                       |

PostgreSQL 17 从 `pg_stat_bgwriter` 移除了 checkpoint 与后端写入列，采集器改从 `pg_stat_checkpointer` 和 `pg_stat_io` 读取，指标名与语义不变。

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

## 日志系统

管理后台「日志」可在线查看平台近期运行记录；也可通过 Kubernetes 查看应用容器日志。

### 日志级别

后台日志页按所选级别精确筛选，只显示该级别的日志，选择「全部」可查看所有级别。应用日志固定按 `DEBUG` 级别记录，因此页面上各级别都能查到；数据库日志级别通过 `gorm.log.level` 设置，任务日志级别通过 `asynq.log.level` 设置。

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

PostgreSQL/GORM、Valkey、数据目录、Gin 监听地址和监听端口在页面中只读，只能通过部署配置修改。

共享 PVC `k8s.shared_volume_claim` 同样是部署参数，修改后需要重新创建平台 Pod。后台「重启」用于重新加载在线设置，不会重新读取磁盘配置文件。各项设置的生效方式见[配置说明](/deploy/settings.md)。

:::warning
重启会短暂中断服务。正式比赛期间谨慎使用。
:::
