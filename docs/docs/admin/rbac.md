---
title: 权限与角色
description: 按 CBCTF 的用户组、角色和 API 权限模型配置后台访问与选手授权。
---

# 权限与角色

## 授权模型

```text
权限 Permission → 角色 Role → 分组 Group → 用户 User
```

- 权限点由后端路由映射内置，后台可查看和编辑描述。
- 每个分组关联一个角色；一个角色包含多个权限。
- 一个用户可加入多个分组，最终权限取各组角色的并集。界面不直接向用户分配多个角色。
- 比赛中的队伍与 RBAC 分组是不同概念，加入参赛队伍不会自动授予后台权限。

## 内置角色与分组

| 名称 | 默认用途 |
| --- | --- |
| `admin` | 全部平台管理和选手权限 |
| `organizer` | 赛事、题目、队伍、公告、作弊以及比赛内镜像/靶机/生成器管理，同时具备选手权限 |
| `user` | 个人设置、参赛、队伍、题目、提交、靶机和题解权限 |

初始化会创建同名默认分组并绑定角色；每次启动会补齐默认角色映射中缺失的权限。需要长期限制权限时创建自定义角色和分组，不要依赖从默认角色删去权限后永远保持不变。

角色权限是接口级权限，`organizer` 不等于“仅能管理自己创建的比赛”；当前路由没有比赛所有者授权模型。

## 前端操作流程

1. 在「RBAC → 角色」创建自定义角色。
2. 打开角色的权限管理，选择需要的权限点。
3. 在「分组」创建组并选择该角色。
4. 在组成员管理中添加用户。
5. 使用目标用户验证后台导航与具体操作。

导航按 `GET /me/permissions` 返回的 **HTTP 方法 + 路由** 过滤；某页可见不表示所有按钮对应的接口均有权限。RBAC 页面入口本身要求 `GET /admin/roles`，自定义管理角色还需配齐其实际使用的列表、读取和修改权限。

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

完整权限名称和接口映射以后台「权限」页及 `internal/model/permission.go` 为准。撤销 `self:read` 会使前端无法加载 `/me` 与权限信息，不能正常完成登录后的会话初始化。
