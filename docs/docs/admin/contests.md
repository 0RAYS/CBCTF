---
sidebar_position: 5
---

# 比赛管理

## 比赛生命周期

```
创建 → 配置题目 → 参赛者注册 → 比赛进行中 → 结束后处理
```

1. **创建比赛**：配置名称、时间、flag 前缀、队伍规模等基础信息
2. **添加题目**：从全局题库选择并加入比赛，可覆盖比赛特定参数
3. **配置计分**：为每个 flag 设置计分类型和参数
4. **镜像预热**：比赛开始前 30-60 分钟拉取镜像到所有节点
5. **参赛管理**：比赛中管理队伍、发布公告、处理作弊
6. **赛后处理**：收集 Writeup、导出数据

## 比赛状态切换

平台根据 `start` 和 `duration` 字段自动切换状态，无需手动操作：

| 时间条件 | 状态 |
|---------|------|
| 当前时间 < start | `coming` |
| start ≤ 当前时间 < start + duration | `running` |
| 当前时间 ≥ start + duration | `over` |

## 队伍管理

通过 `GET /admin/contests/:contestID/teams`（`admin:team:list`）查看所有队伍。

管理操作：

| 操作 | API | 权限 |
|------|-----|------|
| 查看队伍详情 | `GET /admin/contests/:contestID/teams/:teamID` | `admin:team:read` |
| 查看队伍成员 | `GET .../users` | `admin:team:read` |
| 修改队伍信息 | `PUT .../teams/:teamID` | `admin:team:update` |
| 踢出成员 | `POST .../kick` | `admin:team:update` |
| 删除队伍 | `DELETE .../teams/:teamID` | `admin:team:delete` |
| 查看提交记录 | `GET .../submissions` | `admin:team:read` |
| 查看 Flag | `GET .../flags` | `admin:team:read` |

## 批量靶机控制

通过 `POST /admin/contests/:contestID/victims` 和 `DELETE /admin/contests/:contestID/victims`（`admin:victim:control`）批量启动/停止该比赛的所有靶机。

## 流量捕获

管理员可为指定队伍的靶机下载流量捕获文件：

```
GET /admin/contests/:contestID/teams/:teamID/victims/:victimID/traffic/download
```

需要 `admin:team:read` 权限。下载的 pcap 文件可用 Wireshark 分析。

## 公告发布

通过 `POST /admin/contests/:contestID/notices`（`admin:notice:create`）发布公告，平台通过 WebSocket 实时推送给所有在线的参赛者。

## 镜像预热

通过 `POST /admin/contests/:contestID/images`（`admin:image:pull`）触发将该比赛所有题目的容器镜像预拉取到各 Kubernetes 节点，避免比赛开始时因拉取镜像导致延迟。

**建议在比赛开始前 30-60 分钟执行。**

## Writeup 收集

选手上传 Writeup 后，管理员可通过以下方式访问：

```
GET /admin/contests/:contestID/teams/:teamID/writeups        # 列表（需 admin:team_writeup:list）
GET /admin/contests/:contestID/teams/:teamID/writeups/:fileID  # 下载（需 admin:team_writeup:read）
```

## 作弊检测

详见 [作弊检测](/docs/admin/cheat)。通过 `POST /admin/contests/:contestID/cheats` 可手动重新运行一次全量作弊检测。

## 排行榜

管理员视角排行榜通过 `GET /admin/contests/:contestID/rank` 查看（`admin:contest:rank`），功能与选手视图相同但不受隐藏比赛限制。
