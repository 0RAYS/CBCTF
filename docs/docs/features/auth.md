---
sidebar_position: 7
---

# 认证与 OAuth

## 本地认证

CBCTF 默认使用邮箱 + 密码认证。登录成功后颁发 JWT，存储于浏览器 `localStorage`。

### JWT 配置

| 配置项 | 说明 |
|--------|------|
| `gin.jwt.secret` | JWT 签名密钥，留空时自动生成 32 位随机串 |
| `gin.jwt.static` | `false`（默认）：平台定期轮转 JWT 密钥；`true`：密钥固定不变 |

### 设备绑定

平台使用 [FingerprintJS](https://github.com/fingerprintjs/fingerprintjs) 生成设备指纹（通过请求头 `X-M` 传递），并将其绑定到 JWT 中。若同一账号从不同设备登录，会触发作弊检测（`token_magic` 类型）。

### 邮箱验证

需要配置 SMTP 并为用户角色分配 `self:activate` 权限。用户调用 `POST /me/activate` 发送激活邮件，点击邮件中的链接完成验证。

## OAuth / OIDC

CBCTF 支持任意标准 OIDC 提供商，可同时配置多个 OAuth 提供商。

### OAuth 配置字段

| 字段 | 说明 |
|------|------|
| `provider` | 提供商标识（如 `github`、`oidc`、`google`） |
| `auth_url` | OAuth 授权端点 URL |
| `token_url` | OAuth token 端点 URL |
| `user_info_url` | 用户信息端点 URL（OIDC UserInfo） |
| `client_id` | OAuth 客户端 ID |
| `client_secret` | OAuth 客户端密钥 |
| `uri` | URL 中的标识符，对应路由 `/oauth/{uri}/redirect` |
| `id_claim` | 用户 ID 字段名（如 `sub`） |
| `name_claim` | 用户名字段名（如 `name`、`login`） |
| `email_claim` | 邮箱字段名（如 `email`） |
| `picture_claim` | 头像 URL 字段名（如 `picture`、`avatar_url`） |
| `description_claim` | 简介字段名（可选） |
| `groups_claim` | 用户所属群组字段名（可选，用于管理员自动提权） |
| `default_group` | OAuth 登录用户默认加入的分组 ID |

### 管理员组

若配置了 `groups_claim`，当用户的群组列表中包含指定的管理员组名称时，平台会自动为该用户分配管理员角色。

### GitHub OAuth 示例

在 GitHub 创建 OAuth App（Settings → Developer settings → OAuth Apps），回调 URL 设为：
`https://your.domain.com/api/oauth/github/callback`

```json
{
  "provider": "github",
  "auth_url": "https://github.com/login/oauth/authorize",
  "token_url": "https://github.com/login/oauth/access_token",
  "user_info_url": "https://api.github.com/user",
  "client_id": "your-client-id",
  "client_secret": "your-client-secret",
  "uri": "github",
  "id_claim": "id",
  "name_claim": "login",
  "email_claim": "email",
  "picture_claim": "avatar_url"
}
```

## 注册控制

| 配置 | 说明 |
|------|------|
| `registration.enabled: true` | 开放公开注册（默认） |
| `registration.enabled: false` | 禁止自助注册，管理员仍可手动创建用户 |
| `registration.default_group` | 注册用户自动加入的分组 ID |

## 多提供商并存

可以同时配置本地认证 + 多个 OAuth 提供商。登录页面会显示所有可用的登录方式，用户可选择任意方式登录。不同提供商登录的用户通过邮箱关联为同一账号。
