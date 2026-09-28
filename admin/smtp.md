> For AI agents: the complete documentation index is available at /llms.txt, the full documentation bundle is available at /llms-full.txt.

# 邮件配置

平台通过 SMTP 发送邮箱验证、密码重置和管理员测试邮件。入口为「管理后台 → SMTP」。

## 配置步骤

1. 创建 SMTP，填写邮箱、服务器地址、端口及密码或授权码。
2. 如需参与自动发信，勾选「启用」。平台从启用且成功连接的账号池中随机选择发件账号。
3. 保存后使用「测试」，填写收件邮箱，确认实际收信。
4. 查看该配置或全局邮件历史，再用普通账号验证注册邮件与密码找回。

测试邮件使用指定 SMTP 直接发信，不经过随机账号池，也不要求该配置已启用。测试成功不等于其他启用账号都可用。

## 字段

| 字段        | 类型     | 说明                     |
| --------- | ------ | ---------------------- |
| `address` | string | 发件邮箱，同时用作 SMTP 登录账号    |
| `host`    | string | SMTP 服务器主机名            |
| `port`    | int    | 通常为 465 或 587，按服务商要求填写 |
| `pwd`     | string | 密码或应用授权码               |
| `on`      | bool   | 创建/编辑时控制是否启用           |

当前表单没有独立的 SMTP 用户名、TLS 模式或证书设置；使用 gomail 的端口/TLS 行为。编辑界面密码留空时保留原密码。

```json
{
  "address": "noreply@example.com",
  "host": "smtp.example.com",
  "port": 587,
  "pwd": "your-smtp-password",
  "on": true
}
```

## 验证与密码找回

- 本地注册自动提交验证邮件任务；个人设置可重新发送，需 `self:activate` 权限。
- 验证邮件打开 `{host}/platform/#/verify?token=...`，页面调用 `POST /verify` 完成验证。
- 登录页“忘记密码”发送重置链接，打开 `{host}/platform/#/reset-password?token=...`；成功重置同时标记邮箱已验证。
- `host` 必须是收件人可访问的地址。注册成功或任务入队不表示邮件已送达，需要检查任务和邮件历史。

## 权限与排查

创建使用 `admin:smtp:create`，测试使用 `admin:smtp:test`，邮件列表使用 `admin:smtp:list`。相关接口：

```text
POST /admin/smtp/:smtpID/test          请求 {"to":"recipient@example.com"}
GET  /admin/smtp/:smtpID/email         指定账号邮件历史
GET  /admin/email                     全局邮件历史
```

若测试成功但验证邮件失败，检查 Redis、Asynq 邮件任务、SMTP 启用状态，以及其他启用账号是否连接成功。若账号池为空，可修复配置后重新保存启用状态或重启服务以重新连接。
