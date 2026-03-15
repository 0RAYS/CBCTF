---
sidebar_position: 7
---

# 邮件配置

CBCTF 通过 SMTP 发送邮件（目前用于邮箱验证功能）。

## SMTP 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `address` | string | 发件人邮箱地址（如 `noreply@example.com`） |
| `host` | string | SMTP 服务器地址 |
| `port` | int | SMTP 端口（通常 25/465/587） |
| `pwd` | string | SMTP 账号密码或授权码 |
| `on` | bool | 是否启用该 SMTP 服务器 |

## 添加 SMTP 服务器

通过 `POST /admin/smtp`（`admin:smtp:create`）创建：

```json
{
  "address": "noreply@example.com",
  "host": "smtp.example.com",
  "port": 587,
  "pwd": "your-smtp-password",
  "on": true
}
```

## 多 SMTP 服务器

可以配置多个 SMTP 服务器，平台发送邮件时会从启用的服务器列表中随机选择一个。适用于负载均衡或备用 SMTP 场景。

## 邮件使用场景

目前 CBCTF 使用邮件的场景：

- **邮箱验证**：用户调用 `POST /me/activate` 后，平台发送含激活链接的验证邮件

## 启用邮箱验证的前置条件

1. 配置至少一个启用状态的 SMTP 服务器
2. 为 `user` 角色（或相关角色）分配 `self:activate` 权限，用户才能触发发送激活邮件

## 发送统计

通过 `GET /admin/smtp/:smtpID/email`（`admin:smtp:list`）查看该 SMTP 服务器的邮件发送历史，包含成功/失败次数和详细日志。
