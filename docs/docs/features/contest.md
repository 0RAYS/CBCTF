---
sidebar_position: 6
---

# 比赛配置

比赛（Contest）是 CBCTF 的核心组织单元，每场比赛包含独立的题目集、队伍、排行榜和作弊检测记录。

## 比赛状态

| 状态 | 条件 | 说明 |
|------|------|------|
| `coming` | 当前时间 < `start` | 比赛未开始，参赛者可加入队伍但无法看到题目 |
| `running` | `start` ≤ 当前时间 < `start + duration` | 比赛进行中，题目可见（非隐藏比赛） |
| `over` | 当前时间 ≥ `start + duration` | 比赛已结束，不再接受提交 |

## 字段参考

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 比赛名称 |
| `description` | string | 比赛介绍（支持 Markdown） |
| `prefix` | string | Flag 前缀，替换 `static{}/dynamic{}/uuid{}` 的输出前缀，如 `flag`、`CTF`、`CBCTF` |
| `size` | int | 最大队伍人数 |
| `start` | datetime | 比赛开始时间 |
| `duration` | int | 比赛持续时长（秒） |
| `blood` | bool | 是否启用三血奖励 |
| `hidden` | bool | 隐藏比赛（不出现在公开列表，仅邀请制） |
| `captcha` | bool | 加入队伍时是否需要验证码 |
| `victims_count` | int | 每队最大同时运行靶机数 |
| `rules` | string | 比赛规则（支持 Markdown） |
| `prizes` | []object | 奖励配置，`[{amount, description}]` |
| `timelines` | []object | 时间线，`[{date, title, description}]` |

## Flag 前缀

`prefix` 字段决定比赛中所有 flag 的实际输出格式：

```
static{content}  →  {prefix}{content}
dynamic{content} →  {prefix}{随机内容}
uuid{}           →  {prefix}{uuid}
```

例如，`prefix: "flag"` 则选手提交 `flag{...}`；`prefix: "CBCTF"` 则提交 `CBCTF{...}`。

## 公告系统

管理员可发布公告，平台通过 WebSocket 实时推送给所有在线参赛者，无需刷新页面即可收到通知。

## 奖励配置

```json
[
  {"amount": "一等奖", "description": "¥5000"},
  {"amount": "二等奖", "description": "¥3000"}
]
```

在比赛页面展示奖励信息。

## 时间线

```json
[
  {"date": "2025-01-01T09:00:00Z", "title": "比赛开始", "description": ""},
  {"date": "2025-01-01T12:00:00Z", "title": "线上破题截止", "description": ""}
]
```

在比赛详情页以时间轴形式展示里程碑事件。

## 排行榜

比赛提供三种视图：

- **总分排名**：实时总分排序，支持队伍搜索
- **时间线视图**：各队分数随时间变化的折线图
- **计分板（Scoreboard）**：横轴为队伍，纵轴为题目，展示每道题的完成情况

## Writeup 管理

比赛结束后，参赛队伍可上传 Writeup（PDF 或 ZIP 格式）。管理员可从后台按队下载 Writeup。

## 题目管理

在比赛中，可从全局题库添加题目。添加时可覆盖以下字段（仅对本场比赛生效，不修改全局题库）：

- `name`：比赛中显示的题目名称
- `description`：比赛中显示的题目描述
- `category`：题目分类
- `tags`：题目标签
- `hints`：提示内容
- `hidden`：是否在比赛中隐藏该题目
- `attempt_limit`：提交尝试次数限制（0 表示无限制）

每个 flag 的计分参数（`score_type`、`Score`、`MinScore`、`Decay`）也在比赛层级独立配置。
