---
title: 权限与角色
description: 通过用户分组和角色配置后台访问、比赛管理与选手操作权限。
---

# 权限与角色

## 授权模型

```text
权限 Permission → 角色 Role → 分组 Group → 用户 User
```

- 平台提供内置权限点，可在后台查看用途和编辑描述。
- 每个分组关联一个角色；一个角色包含多个权限。
- 一个用户可加入多个分组，最终权限取各组角色的并集。界面不直接向用户分配多个角色。
- 比赛中的队伍与 RBAC 分组是不同概念，加入参赛队伍不会自动授予后台权限。

## 内置角色与分组

| 名称 | 默认用途 |
| --- | --- |
| `admin` | 全部平台管理和选手权限 |
| `organizer` | 赛事、题目、队伍、公告、作弊以及比赛内镜像/靶机/生成器管理，同时具备选手权限 |
| `user` | 个人设置、参赛、队伍、题目、提交、靶机和题解权限 |

平台提供与内置角色同名的默认分组。内置角色的权限会在启动时补齐；如需精细控制授权，请创建自定义角色和分组。

`organizer` 可管理平台上的比赛，授权范围不限于该用户创建的比赛。分配此角色前请确认管理范围。

## 配置步骤

1. 在「RBAC → 角色」创建自定义角色。
2. 打开角色的权限管理，选择需要的权限点。
3. 在「分组」创建组并选择该角色。
4. 在组成员管理中添加用户。
5. 使用目标用户验证后台导航与具体操作。

导航会显示用户有权访问的页面。查看列表、读取详情和修改内容可能需要不同权限；配置自定义角色时，请按实际操作同时勾选。访问 RBAC 页面还需要查看角色列表的权限。

## 常用权限范围

以下 `*` 只是说明一组权限，不是可提交的通配符：

| 场景 | 主要权限 |
| --- | --- |
| 个人信息与邮箱验证 | `self:read`、`self:update`、`self:activate` |
| 参赛 | `user:contest:*`、`user:team:*`、`user:notice:list` |
| 解题与题解 | `user:challenge:*`、`user:victim:control`、`user:writeup:*` |
| 题库与测试 | `admin:challenge:*` |
| 比赛与比赛题目 | `admin:contest:*`、`admin:contest_challenge:*`、`admin:contest_challenge_flag:*` |
| 队伍与题解 | `admin:team:*`、`admin:team_writeup:list/read`、`admin:contest_writeup:export` |
| 全局资源 | `admin:image:pull`、`admin:victim:control`、`admin:generator:control`、`admin:traffic:read` |
| 比赛资源 | `admin:contest_image:pull`、`admin:contest_victim:control`、`admin:contest_generator:control`、`admin:contest_traffic:read` |
| SMTP | `admin:smtp:*`，其中发测试邮件为 `admin:smtp:test` |
| 定时任务 | `admin:cronjob:list`、`admin:cronjob:update` |
| 系统配置 | `admin:system:read/update/restart` |

完整权限名称和用途可在后台「权限」页查看。登录用户需保留 `self:read`，用于读取个人信息和正常使用平台。
