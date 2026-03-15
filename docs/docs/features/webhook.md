---
sidebar_position: 8
---

# Webhook

Webhook 允许平台在特定事件发生时，向外部系统（Slack、Discord、自动化脚本等）发送 HTTP 通知。

## Webhook 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | Webhook 名称（仅用于标识） |
| `url` | string | 目标 URL |
| `method` | string | HTTP 方法（`POST`、`GET` 等） |
| `headers` | object | 自定义 HTTP 请求头，如 `{"Authorization": "Bearer token"}` |
| `timeout` | int | 请求超时（毫秒） |
| `retry` | int | 失败后重试次数 |
| `on` | bool | 是否启用此 Webhook |
| `events` | []string | 订阅的事件类型列表 |

## 触发事件

Webhook 支持订阅平台的多种事件，主要包括：

**用户事件**
- 用户注册、登录、修改信息

**比赛事件**
- 比赛创建、开始、结束

**队伍事件**
- 队伍创建、加入、退出

**题目事件**
- 题目初始化、重置、靶机启停

**Flag 事件**
- Flag 提交成功（`flag_correct`）
- Flag 提交错误（`flag_wrong`）
- 首血、二血、三血事件

## 投递历史

每个 Webhook 都维护投递历史记录，包含：

- `success`：累计成功投递次数
- `failure`：累计失败投递次数
- 最后成功/失败时间
- 详细历史日志（包含完整请求和响应内容）

通过 `GET /admin/webhook/:webhookID/history` 查看。

## 安全配置

在 `config.yaml` 中配置 Webhook 目标 URL 白名单：

```yaml
webhook:
  whitelist:
    - "https://hooks.slack.com/"
    - "https://discord.com/api/webhooks/"
```

配置后，平台仅允许向匹配前缀的 URL 发送 Webhook 请求，防止 SSRF 攻击。留空 `[]` 表示不限制。

## 自定义认证

通过 `headers` 字段添加认证信息：

```json
{
  "Authorization": "Bearer your-token",
  "X-Secret": "your-secret"
}
```

## 重试机制

Webhook 投递通过 Asynq 异步队列处理。若请求失败（超时、非 2xx 响应），按 `retry` 配置的次数自动重试，重试间隔指数递增。

## 示例：Slack 通知

```json
{
  "name": "Slack Flag 通知",
  "url": "https://hooks.slack.com/services/YOUR/WEBHOOK/URL",
  "method": "POST",
  "headers": {"Content-Type": "application/json"},
  "timeout": 5000,
  "retry": 3,
  "on": true,
  "events": ["flag_correct", "first_blood", "second_blood", "third_blood"]
}
```
