---
sidebar_position: 4
---

# 题目管理

## 题库与比赛题目

CBCTF 采用两层题目管理：

| 层级 | 说明 |
|------|------|
| **全局题库** | 题目模板，不直接对选手可见 |
| **比赛题目（ContestChallenge）** | 将全局题目关联到比赛，可覆盖部分字段 |

对全局题库的修改（如更新附件）会影响所有使用该题目的比赛。对比赛题目的覆盖配置（name、description 等）仅影响该场比赛。

## 创建题目

必填字段：`name`（题目名称）、`type`（题目类型）、`category`（分类）。

可选字段：
- `description`：题目描述（Markdown）
- `generator_image`：动态附件题的生成器容器镜像
- `network_policies`：Kubernetes NetworkPolicy 配置（JSON）

## 上传附件

- **静态/问答/容器题**：上传 `attachment.zip`，所有队伍共享下载
- **动态附件题**：上传 `generator.zip`，平台在生成时上传到生成器容器的 `/root` 目录（可选）

## 题目测试模式

需要 `admin:challenge:test` 权限。管理员可在题目加入比赛前，通过测试模式验证容器靶机配置：

1. 访问题目详情，点击「测试」
2. 平台以管理员身份启动靶机容器（不绑定任何队伍）
3. 验证端口连通性、flag 注入是否正确
4. 测试完成后停止靶机

测试模式下不产生实际的 flag 和队伍记录。

## 加入比赛

通过 `POST /admin/contests/:contestID/challenges`（`admin:contest_challenge:create`）将题目加入比赛。

加入时可覆盖以下字段（仅对本场比赛生效）：

| 字段 | 说明 |
|------|------|
| `name` | 比赛中显示的题目名称 |
| `description` | 比赛中显示的题目描述 |
| `category` | 分类 |
| `tags` | 标签 |
| `hints` | 提示列表（选手可查看） |
| `hidden` | 是否在比赛中隐藏该题目 |
| `attempt_limit` | 提交次数限制（0 表示无限制） |

## 管理比赛 Flag 分数

通过 `PUT /admin/contests/:contestID/challenges/:challengeID/flags/:flagID`（`admin:contest_challenge_flag:update`）为每个 flag 独立配置计分：

```json
{
  "score_type": 2,
  "score": 1000,
  "min_score": 100,
  "decay": 50
}
```

详见 [计分系统](/docs/features/scoring)。

## NetworkPolicy

在题目的 `network_policies` 字段中配置 Kubernetes NetworkPolicy（JSON 格式），对该题目的所有靶机生效。

示例：禁止靶机出站访问互联网：

```json
[{
  "policyTypes": ["Egress"],
  "egress": [{
    "to": [{"ipBlock": {"cidr": "10.0.0.0/8"}}]
  }]
}]
```

## 文件管理

通过 `GET /admin/files`（`admin:file:list`）查看所有已上传的题目附件，支持下载和删除。流量 pcap 文件和 Writeup 文件也在此处统一管理。
