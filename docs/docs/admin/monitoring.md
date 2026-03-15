---
sidebar_position: 8
---

# 监控与日志

## 系统状态

通过 `GET /admin/system/status`（`admin:system:status`）查看各服务的健康状态，包括 MySQL、Redis、Asynq 等组件的连接状态。

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

- HTTP 请求数和延迟（按路由、方法、状态码分类）
- Asynq 任务队列深度和处理速率
- 作弊检测计数
- 活跃靶机数量

## 日志系统

平台使用 Logrus 记录日志，支持文件轮转。

### 日志级别

通过 `log.level` 配置，支持 `DEBUG` `INFO` `WARNING` `ERROR`。

### 日志持久化

设置 `log.save: true` 后，日志持久化到 `{path}/logs/` 目录下，自动按日期轮转。

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
- 选手 Writeup（PDF/ZIP）

支持在线下载（`admin:file:read`）和删除（`admin:file:delete`）。

## 系统配置在线更新

无需重启服务器即可更新配置：

```bash
# 读取当前配置
GET /admin/system/config    # 需 admin:system:read

# 更新配置
PUT /admin/system/config    # 需 admin:system:update

# 应用配置（热重启）
POST /admin/system/restart  # 需 admin:system:restart
```

:::caution
热重启会短暂中断服务（通常 1-2 秒）。正式比赛期间谨慎使用。
:::
