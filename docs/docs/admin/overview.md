---
sidebar_position: 1
---

# 管理后台概览

## 访问入口

管理后台位于 `/platform/#/admin`，需要具有 `admin` 角色或拥有相应管理权限的用户才能访问。

## 初始管理员

平台首次启动时，若数据库中不存在管理员账号，自动创建初始管理员并将凭据打印到日志：

```
Init Admin: Admin{ name: admin, password: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx, email: admin@0rays.club}
```

:::warning
首次登录后请立即修改初始密码。
:::

## admin 与 organizer 的区别

| 角色 | 特点 |
|------|------|
| `admin` | 拥有所有权限（`SetFullAccess`），可管理整个平台 |
| `organizer` | 拥有部分管理权限，通常用于授权赛事组织者管理特定比赛 |

organizer 角色由管理员通过 RBAC 系统自定义权限集合，详见 [权限与角色](/docs/admin/rbac)。

## 后台导航结构

| 功能区 | 说明 |
|--------|------|
| 系统状态 | 查看各服务健康状态、Prometheus 指标 |
| 用户管理 | 用户列表、创建用户、分配角色、管理分组 |
| 角色权限 | 创建角色、配置权限、分配给用户 |
| OAuth | 配置 OAuth/OIDC 提供商 |
| SMTP | 配置邮件服务 |
| Webhook | 配置出站 Webhook |
| 题目管理 | 全局题库，创建/编辑题目，管理附件 |
| 比赛管理 | 创建比赛，管理队伍、题目、公告、作弊记录 |
| 文件管理 | 查看/下载/删除附件、流量 pcap、Writeup |
| 日志 | 查看访问日志和系统日志 |

## 全局搜索

访问 `/admin/search`（需 `admin:models:search` 权限），支持跨模型全局搜索，快速定位用户、队伍、题目、提交记录等数据。

## 系统配置在线更新

通过 `PUT /admin/system/config` 可在线更新配置，通过 `POST /admin/system/restart` 热重启服务，无需登录服务器。需要 `admin:system:update` 和 `admin:system:restart` 权限。
