---
title: Webhook
description: 按实际事件列表、秒级超时与主机白名单配置 Webhook，检查投递历史和重试。
---

# Webhook

入口：「管理后台 → Webhook」。平台将匹配事件放入 Asynq 队列，以 JSON 请求体投递到你的接收服务。

## 先配置目标白名单

在「系统管理」配置 `webhook.whitelist`。**空列表拒绝创建目标，不是不限制**。白名单匹配主机名、主机名加端口、IP 或 IP CIDR，不接受 URL 路径前缀：

```yaml
webhook:
  whitelist:
    - hooks.example.com
    - notify.example.com:8443
    - 192.0.2.10
    - 192.0.2.0/24
```

Helm 初次部署写入 `cbctf.webhook.whitelist`；已有数据库通过系统管理更新。白名单在创建/修改 URL 时校验，不是对已有投递、DNS 解析结果和 HTTP 重定向的网络防火墙。

## 创建与启用

填写名称、目标 URL、方法、超时、重试次数、请求头，选择订阅事件并勾选启用。创建和编辑均支持启用开关。

| 字段 | 含义 |
| --- | --- |
| `name` / `url` | 名称与目标完整 URL |
| `method` | 仅 `POST` 或 `GET`，两者都携带 JSON 请求体 |
| `headers` | 请求头对象，默认发送 `Content-Type: application/json` |
| `timeout` | **秒**，`0` 使用 HTTP 客户端默认 30 秒 |
| `retry` | 失败后的最大重试次数，`0` 不重试 |
| `events` | 事件名称字符串列表；空列表订阅所有事件 |
| `on` | 是否启用 |

## 事件名称和负载

编辑器从 `GET /admin/webhook/events` 加载实际事件；以该列表为准。常见值有：

- 用户：`register`、`login`、`oauth_login`、`update_user`
- 队伍：`create_team`、`join_team`、`leave_team`
- 题目：`init_challenge`、`reset_challenge`、`start_victim`、`stop_victim`
- 提交：`submit_flag`
- 公告：`create_notice`、`update_notice`、`delete_notice`

没有 `flag_correct`、`flag_wrong`、`first_blood` 等独立 Webhook 事件，也没有比赛时间到达后自动产生的开始/结束通知。

实际负载只包含事件类型、来源 IP 和相关对象 ID：

```json
{
  "type": "submit_flag",
  "ip": "192.0.2.20",
  "models": {"Self": 12, "Contest": 3, "Team": 8}
}
```

`models` 的键依事件上下文变化。负载不包含 Flag 内容、解题正误、完整对象或事件成功状态，不应仅凭 `submit_flag` 就宣布解题成功。

## 可用配置示例

```json
{
  "name": "赛事事件接收器",
  "url": "https://hooks.example.com/cbctf",
  "method": "POST",
  "headers": {"Authorization": "Bearer your-token"},
  "timeout": 10,
  "retry": 3,
  "on": true,
  "events": ["submit_flag", "create_notice"]
}
```

Slack、Discord 等通常要求自己的消息格式，需要接收服务将上述负载转换后再发送，不能直接把 CBCTF 的 JSON 当作聊天消息模板。

## 历史与重试

页面提供全局和单个 Webhook 的历史，展示成功状态、HTTP 响应码、耗时、错误及关联事件，不保存完整请求/响应正文。

网络错误、超时和非 2xx 响应会向任务队列返回失败，按 `retry` 重试；每次尝试分别记录历史。接收器应考虑重复投递。配置修改影响后续入队任务，已入队任务携带的是当时的目标配置。
