---
title: 认证与 OAuth
description: 配置 CBCTF 本地认证、邮箱验证、OAuth/OIDC 登录和用户组自动分配。
---

# 认证与 OAuth

## 本地认证

CBCTF 使用用户名、密码和图形验证码登录，登录态保存在 HttpOnly Cookie 中。当前登录表单不使用邮箱代替用户名。

### JWT 配置

| 配置项              | 说明               |
|------------------|------------------|
| `gin.jwt.secret` | JWT 签名密钥，必须替换默认值 |

### 邮箱验证

需要先配置 SMTP，并为用户分配 `self:activate` 权限。注册会发送验证邮件，也可从个人设置重新发送；邮件打开 `/platform/#/verify?token=...`，前端调用 `POST /verify` 完成验证。密码找回通过登录页操作，成功重置密码也会标记邮箱已验证。

## OAuth / OIDC

平台提供 `oauth2` 和 `cas` 两类协议。OAuth2 可以对接提供授权、Token 和 UserInfo 端点的 OIDC 服务，但界面不是自动 Discovery 配置，需要手动填写端点与 claim 提取表达式。公开入口为：

- `GET /oauth`
- `GET /oauth/{uri}`
- `GET /oauth/{uri}/callback`
- `GET /oauth/token`

后端完成第三方回调后，会跳转到前端 `#/oauth/callback` 页面继续登录流程。

### OAuth 配置字段

| 字段                  | 说明                       |
|---------------------|--------------------------|
| `provider`          | 提供商名称                    |
| `protocol`          | 创建/更新 API 的查询参数，`oauth2` 或 `cas`，不是正文中的开关 |
| `scopes`            | OAuth2 scope 字符串列表 |
| `auth_url`          | 授权端点                     |
| `token_url`         | Token 端点                 |
| `user_info_url`     | 用户信息端点                   |
| `callback_url`      | 第三方回调地址                  |
| `client_id`         | 客户端 ID                   |
| `client_secret`     | 客户端密钥                    |
| `uri`               | 平台路由标识，对应 `/oauth/{uri}` |
| `id_claim`          | 用户 ID 提取表达式              |
| `name_claim`        | 用户名提取表达式                 |
| `email_claim`       | 邮箱提取表达式                  |
| `picture_claim`     | 头像提取表达式                  |
| `description_claim` | 简介提取表达式                  |
| `groups_claim`      | 用户组字段                    |
| `admin_group`       | 命中后自动加入 `admin` 用户组的组名   |
| `default_group`     | 首次登录后自动加入的默认分组           |
| `on`                | 是否启用该提供商                 |

### GitHub 示例

GitHub OAuth App 回调地址示例：`https://your.domain.com/oauth/github/callback`

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
  "scopes": ["user:email"],
  "name_claim": "{login}",
  "email_claim": "{email}",
  "picture_claim": "{avatar_url}",
  "description_claim": "{bio}"
}
```

在后台「OAuth」编辑内置 GitHub 配置即可；填写 Client ID/Secret，核对 callback URL、`user:email` scope 与上述 claim，再启用。GitHub 邮箱处理会读取已验证的主邮箱。旧数据库中已保存的 `{picture_url}` 不会因代码修正自动迁移，需要改为 `{avatar_url}`。

新建 OAuth/CAS 提供商默认未启用，创建后再编辑启用。CAS 使用登录端点与验证端点，不使用 OAuth2 Token 交换字段。更改 `host` 不会自动更新已经保存的提供商 callback URL。

## 注册控制

| 配置                            | 说明               |
|-------------------------------|------------------|
| `registration.enabled: true`  | 开放公开注册           |
| `registration.enabled: false` | 禁止本地自助注册；已启用的第三方登录仍有自己的用户创建流程 |
| `registration.default_group`  | 注册用户自动加入的分组 ID   |

## 多提供商并存

平台可同时启用本地登录和多个 OAuth 提供商。登录页会展示所有启用的提供商入口。

提供商 `default_group` 是分组 ID；`admin_group` 按第三方用户组命中后加入平台管理员组。部署前用普通第三方账号检查实际授权，避免把平台默认分组与提供商分组混为一谈。跨域 Cookie 和回调页部署要求见[前后端分离](../../deploy/separation)。

组声明及默认分组在第三方账号首次创建时应用，后续登录目前更新头像、简介和原始用户信息，不自动持续同步组成员关系。已有账号的授权应在 RBAC 页面维护。
