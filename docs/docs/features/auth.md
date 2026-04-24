---
sidebar_position: 7
---

# 认证与 OAuth

## 本地认证

CBCTF 默认使用用户名 + 密码登录

### JWT 配置

| 配置项 | 说明 |
|--------|------|
| `gin.jwt.secret` | JWT 签名密钥, 必须替换默认值 |

当前代码中 **不存在** `gin.jwt.static` 配置项。

### 设备绑定

前端会使用 FingerprintJS 生成设备标识, 并通过请求头 `X-M` 发送给后端。JWT 中也会绑定该标识; 若同一账号从不同设备使用, 可能触发 `token_magic` 作弊检测。

### 邮箱验证

需要先配置 SMTP, 并为用户分配 `self:activate` 权限。用户调用 `POST /me/activate` 后, 平台发送验证邮件; 用户访问邮件中的 `/verify` 链接后, 账号会被标记为已验证。

## OAuth / OIDC

平台支持配置多个 OAuth / OIDC 提供商。公开入口为: 

- `GET /oauth`
- `GET /oauth/{uri}`
- `GET /oauth/{uri}/callback`
- `GET /oauth/token`

后端完成第三方回调后, 会跳转到前端 `#/oauth/callback` 页面继续登录流程。

### OAuth 配置字段

| 字段 | 说明 |
|------|------|
| `provider` | 提供商名称 |
| `auth_url` | 授权端点 |
| `token_url` | Token 端点 |
| `user_info_url` | 用户信息端点 |
| `callback_url` | 第三方回调地址 |
| `client_id` | 客户端 ID |
| `client_secret` | 客户端密钥 |
| `uri` | 平台路由标识, 对应 `/oauth/{uri}` |
| `id_claim` | 用户 ID 提取表达式 |
| `name_claim` | 用户名提取表达式 |
| `email_claim` | 邮箱提取表达式 |
| `picture_claim` | 头像提取表达式 |
| `description_claim` | 简介提取表达式 |
| `groups_claim` | 用户组字段 |
| `admin_group` | 命中后自动加入 `admin` 用户组的组名 |
| `default_group` | 首次登录后自动加入的默认分组 |
| `on` | 是否启用该提供商 |

### GitHub 示例

GitHub OAuth App 回调地址示例: 

```text
https://your.domain.com/oauth/github/callback
```

```json
{
  "provider": "Github",
  "auth_url": "https://github.com/login/oauth/authorize",
  "token_url": "https://github.com/login/oauth/access_token",
  "user_info_url": "https://api.github.com/user",
  "callback_url": "https://your.domain.com/oauth/github/callback",
  "client_id": "your-client-id",
  "client_secret": "your-client-secret",
  "uri": "github",
  "id_claim": "{id}",
  "name_claim": "{name}",
  "email_claim": "{email}",
  "picture_claim": "{picture_url}",
  "description_claim": "{bio}"
}
```

## 注册控制

| 配置 | 说明 |
|------|------|
| `registration.enabled: true` | 开放公开注册 |
| `registration.enabled: false` | 禁止自助注册, 仅管理员可创建用户 |
| `registration.default_group` | 注册用户自动加入的分组 ID |

## 多提供商并存

平台可同时启用本地登录和多个 OAuth 提供商。登录页会展示所有启用的提供商入口。
