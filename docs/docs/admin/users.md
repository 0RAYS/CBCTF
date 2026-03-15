---
sidebar_position: 3
---

# 用户与分组

## 用户字段

| 字段 | 说明 |
|------|------|
| `name` | 用户名 |
| `email` | 邮箱（用于登录和邮件通知） |
| `verified` | 邮箱是否已验证 |
| `picture` | 头像 URL |
| `description` | 个人简介 |

## 手动创建用户

当 `registration.enabled: false` 时，管理员可通过 `POST /admin/users`（需 `admin:user:create`）手动创建用户，绕过注册限制。

## 查找用户

- **列表**：`GET /admin/users`，支持分页和筛选
- **全局搜索**：`/admin/search`（`admin:models:search`），跨所有模型搜索

## 用户分组（Groups）

分组用于批量管理参赛者，常见用途：按班级/院校/队伍组织用户、控制默认分配、赛前批量导入用户。

### 创建分组

`POST /admin/groups`（需 `admin:group:create`）

### 自动分配

在 `registration.default_group` 中配置分组 ID，新注册用户自动加入该分组。也可在 OAuth 配置的 `default_group` 字段中为 OAuth 登录用户单独设置默认分组。

### 手动添加/移除用户

- 添加用户到分组：`POST /admin/groups/:groupID/users`（`admin:user:assign`）
- 移除用户：`DELETE /admin/groups/:groupID/users`（`admin:user:revoke`）

## GeoIP 地理信息

若配置了 `geocity_db`（MaxMind GeoLite2-City），用户列表中会显示最近登录 IP 的地理位置（国家/城市）。

## IP 登录历史

通过 `GET /admin/ip`（`admin:ip:search`）搜索特定 IP 的登录历史，包含登录时间、地理位置、关联用户等信息，用于作弊调查。

## 邮箱验证流程

1. 确保已配置 SMTP 服务器（见 [邮件配置](/docs/admin/smtp)）
2. 为 `user` 角色分配 `self:activate` 权限
3. 用户调用 `POST /me/activate` 触发发送激活邮件
4. 用户点击邮件中的链接，`verified` 字段变为 `true`

## 重置用户密码

管理员通过 `PUT /admin/users/:userID`（`admin:user:update`）可直接修改用户密码，无需原密码验证。
