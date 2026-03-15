---
sidebar_position: 2
---

# 权限与角色

## RBAC 模型

CBCTF 使用基于角色的访问控制（RBAC）：

```
权限（Permission）→ 角色（Role）→ 用户（User）
```

- **权限**：对特定 API 路由的访问授权，由系统内置，不可自定义
- **角色**：权限的集合，可自定义创建
- **用户**：可分配多个角色，权限取所有角色权限的并集

## 内置角色

| 角色 | 说明 |
|------|------|
| `admin` | 拥有所有权限，超级管理员 |
| `organizer` | 默认无权限，需管理员手动配置权限子集 |
| `user` | 默认参赛用户权限（`self:*`、`user:*`） |

## 完整权限清单

### self — 用户自身

| 权限 | 说明 |
|------|------|
| `self:read` | 查看自身信息（登录必须，撤销后无法登录） |
| `self:update` | 更新自身信息（修改密码、头像等） |
| `self:delete` | 删除自身账号 |
| `self:activate` | 发送激活邮件（触发邮箱验证） |

### user:contest — 比赛（用户视角）

| 权限 | 说明 |
|------|------|
| `user:contest:read` | 查看比赛详情 |
| `user:contest:rank` | 查看比赛排名和计分板 |

### user:team — 队伍操作

| 权限 | 说明 |
|------|------|
| `user:team:create` | 创建队伍 |
| `user:team:join` | 加入队伍 |
| `user:team:read` | 查看队伍信息 |
| `user:team:update` | 更新队伍信息 |
| `user:team:delete` | 解散队伍 |

### user:notice — 公告

| 权限 | 说明 |
|------|------|
| `user:notice:list` | 查看比赛公告 |

### user:challenge — 题目操作

| 权限 | 说明 |
|------|------|
| `user:challenge:list` | 查看题目列表 |
| `user:challenge:read` | 查看题目详情和下载附件 |
| `user:challenge:init` | 初始化动态附件或容器靶机 |
| `user:challenge:reset` | 重置题目（重新生成 flag） |
| `user:challenge:submit` | 提交 flag |

### user:victim — 靶机控制

| 权限 | 说明 |
|------|------|
| `user:victim:control` | 启动、延长、停止靶机 |

### user:writeup — Writeup

| 权限 | 说明 |
|------|------|
| `user:writeup:upload` | 上传 Writeup |
| `user:writeup:list` | 查看 Writeup 列表 |

### admin:system — 系统管理

| 权限 | 说明 |
|------|------|
| `admin:system:status` | 查看系统健康状态 |
| `admin:system:read` | 读取系统配置 |
| `admin:system:update` | 更新系统配置 |
| `admin:system:restart` | 重启系统 |

### admin:permission — 权限管理

| 权限 | 说明 |
|------|------|
| `admin:permission:list` | 查看权限列表 |
| `admin:permission:update` | 更新权限描述 |

### admin:role — 角色管理

| 权限 | 说明 |
|------|------|
| `admin:role:create` | 创建角色 |
| `admin:role:read` | 查看角色详情 |
| `admin:role:update` | 更新角色 |
| `admin:role:delete` | 删除角色 |
| `admin:role:list` | 查看角色列表 |
| `admin:role:assign` | 为角色分配权限 |
| `admin:role:revoke` | 撤销角色权限 |

### admin:group — 用户组管理

| 权限 | 说明 |
|------|------|
| `admin:group:create` | 创建用户组 |
| `admin:group:read` | 查看用户组详情 |
| `admin:group:update` | 更新用户组 |
| `admin:group:delete` | 删除用户组 |
| `admin:group:list` | 查看用户组列表 |

### admin:user — 用户管理

| 权限 | 说明 |
|------|------|
| `admin:user:create` | 手动创建用户 |
| `admin:user:read` | 查看用户详情 |
| `admin:user:update` | 修改用户信息（含重置密码） |
| `admin:user:delete` | 删除用户 |
| `admin:user:list` | 查看用户列表 |
| `admin:user:assign` | 将用户加入分组 |
| `admin:user:revoke` | 将用户移出分组 |

### admin:oauth — OAuth 配置

| 权限 | 说明 |
|------|------|
| `admin:oauth:create/read/update/delete/list` | OAuth 提供商的 CRUD 操作 |

### admin:smtp — SMTP 配置

| 权限 | 说明 |
|------|------|
| `admin:smtp:create/read/update/delete/list` | SMTP 服务器的 CRUD 操作 |

### admin:webhook — Webhook

| 权限 | 说明 |
|------|------|
| `admin:webhook:create/read/update/delete/list` | Webhook 的 CRUD 操作 |

### admin:challenge — 题目管理

| 权限 | 说明 |
|------|------|
| `admin:challenge:create` | 创建题目 |
| `admin:challenge:read` | 查看题目详情 |
| `admin:challenge:update` | 更新题目（含上传附件） |
| `admin:challenge:delete` | 删除题目 |
| `admin:challenge:list` | 查看题目列表 |
| `admin:challenge:test` | 测试模式启停靶机 |

### admin:contest — 比赛管理

| 权限 | 说明 |
|------|------|
| `admin:contest:create` | 创建比赛 |
| `admin:contest:read` | 查看比赛详情 |
| `admin:contest:update` | 更新比赛配置 |
| `admin:contest:delete` | 删除比赛 |
| `admin:contest:list` | 查看比赛列表 |
| `admin:contest:rank` | 查看管理视角排行榜 |

### admin:team — 队伍管理（管理视角）

| 权限 | 说明 |
|------|------|
| `admin:team:read` | 查看队伍详情、Flag 记录、流量 pcap |
| `admin:team:update` | 修改队伍信息 |
| `admin:team:delete` | 删除队伍 |
| `admin:team:list` | 查看队伍列表 |

### admin:team_writeup — Writeup 管理

| 权限 | 说明 |
|------|------|
| `admin:team_writeup:list` | 查看 Writeup 列表 |
| `admin:team_writeup:read` | 下载 Writeup 文件 |

### admin:notice — 公告管理

| 权限 | 说明 |
|------|------|
| `admin:notice:create` | 发布公告 |
| `admin:notice:update` | 修改公告 |
| `admin:notice:delete` | 删除公告 |
| `admin:notice:list` | 查看公告列表 |

### admin:cheat — 作弊记录

| 权限 | 说明 |
|------|------|
| `admin:cheat:create` | 手动创建作弊记录（重新运行检测） |
| `admin:cheat:update` | 更新作弊记录状态 |
| `admin:cheat:delete` | 删除作弊记录 |
| `admin:cheat:list` | 查看作弊记录列表 |

### admin:contest_challenge — 比赛题目关联

| 权限 | 说明 |
|------|------|
| `admin:contest_challenge:create/read/update/delete/list` | 比赛中题目的增删改查 |

### admin:contest_challenge_flag — Flag 管理

| 权限 | 说明 |
|------|------|
| `admin:contest_challenge_flag:list` | 查看比赛中各队 Flag |
| `admin:contest_challenge_flag:read` | 查看 Flag 详情 |
| `admin:contest_challenge_flag:update` | 更新 Flag 计分配置 |

### admin:image — 镜像预热

| 权限 | 说明 |
|------|------|
| `admin:image:pull` | 触发镜像预热（拉取到各节点） |

### admin:victim — 靶机控制（管理视角）

| 权限 | 说明 |
|------|------|
| `admin:victim:control` | 管理员批量启停靶机 |

### admin:file — 文件管理

| 权限 | 说明 |
|------|------|
| `admin:file:list` | 查看文件列表 |
| `admin:file:read` | 下载文件 |
| `admin:file:delete` | 删除文件 |

### admin:log — 日志

| 权限 | 说明 |
|------|------|
| `admin:log:read` | 查看系统日志 |

### admin:ip / admin:models — 搜索

| 权限 | 说明 |
|------|------|
| `admin:ip:search` | 搜索用户 IP 登录历史 |
| `admin:models:search` | 全局模型搜索 |

---

## 自定义 organizer 角色示例

创建一个只能管理比赛和题目的 organizer 角色：

```
admin:contest:create/read/update/list
admin:challenge:read/update/list
admin:contest_challenge:create/read/update/delete/list
admin:contest_challenge_flag:list/read/update
admin:team:read/list
admin:team_writeup:list/read
admin:notice:create/update/delete/list
admin:cheat:list/update
admin:image:pull
admin:victim:control
```

:::warning
撤销 `self:read` 权限将导致该用户**无法登录**，操作前请确认。
:::
