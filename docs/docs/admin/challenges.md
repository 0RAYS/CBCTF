---
title: 题目管理
description: 按前端编辑器配置静态题、动态附件题和容器题，理解题库与比赛题目的关系并完成测试。
---

# 题目管理

入口：「管理后台 → 题目管理」，对应 `/platform/#/admin/challenges`。

## 题库与比赛题目

全局题库保存题目模板、附件、生成器镜像、Compose 和原始 Flag。加入比赛时，平台复制名称、描述、分类和 Flag 模板到比赛关联记录；比赛内的名称、描述、提示、标签、隐藏、提交次数与计分独立维护。

附件和底层运行模板仍引用全局题目，更新会影响引用该题目的比赛。已有比赛 Flag 是独立记录，不应假定修改题库 Flag 或 Compose 后会自动同步所有比赛的 Flag 绑定。应在加入正式比赛前完成结构调整；比赛中改动运行模板后要重新测试并检查比赛 Flag。

## 创建题目

1. 选择 `static`、`dynamic` 或 `pods` 类型。创建后编辑 API 不支持切换题型。
2. 填写名称、分类和 Markdown 描述。分类会进行标题化处理。
3. 静态题和动态附件题填写 Flag 模板列表，例如 `static{hello}`、`leet{hello_world}` 或 `uuid{}`。
4. 动态附件题必须填写生成器镜像；容器题必须提供 Compose，Flag 从环境变量、`x-volumes` 或 `x-cloudinit.write_files` 提取。
5. 保存题目后上传附件并测试。

容器编辑器支持表单和 YAML 配置；平台只转换支持的 Compose 字段，不会在节点运行 `docker compose up`。普通容器题至少暴露一个端口。网络、文件注入和 VM 限制见[动态靶机](../guide/features/container)。

## 附件

| 题型 | 上传内容与下载行为 |
| --- | --- |
| 静态题 / 容器题 | 所有队伍共享一个附件；内部路径名为 `attachment.zip`，下载尽量保留上传文件名 |
| 动态附件题 | 可选上传 ZIP 源文件，保存为 `generator.zip`，由 worker 解压到 `/root`；每队下载按 Flag 和源文件版本生成的 ZIP |

题目上传接口目前不强制静态附件后缀为 ZIP，但动态生成器会实际调用 `unzip`，其输入必须是有效 ZIP。上传字段为 `challenge`，大小受 `gin.upload.challenge` 限制。附件上传立即落库，不随编辑弹窗的取消操作回滚。

## 测试

在题目列表/详情进入测试操作：

- 静态题：检查附件下载和内容。
- 动态附件题：先在全局「生成器」页为该题启动测试生成器，再测试下载生成结果。测试使用 `team_id=0` 和题库原始 Flag 模板，与比赛队伍的最终 Flag 不同。
- 容器题：启动测试实例，等待 Running，核对端点、服务、环境变量/文件和资源限制，测试后停止。

测试实例不创建比赛队伍或正式 TeamFlag 记录，但**会注入测试 Flag**，通常为 `flag{...}`。VM 需要可启动的 containerDisk 镜像、VPC、内存、静态 IP 与 MAC，且不会根据 `ports` 生成访问地址。

测试状态与附件接口使用 `admin:challenge:read`；测试启动/停止使用 `admin:challenge:test`。动态测试还需要可用的全局生成器。

## 加入比赛与发布

在「比赛 → 题目」选择并添加题目。新关联默认隐藏，先检查 Flag 和分数，再取消隐藏。比赛题目可编辑字段为：

| 字段 | 说明 |
| --- | --- |
| `name` / `description` | 比赛内名称与 Markdown 描述 |
| `tags` / `hints` | 标签与提示字符串列表 |
| `hidden` | 是否隐藏 |
| `attempt` | 每队提交次数上限，`0` 不限制 |

计分字段使用小写 JSON 键：

```json
{"score_type": 2, "score": 1000, "min_score": 100, "decay": 50}
```

## 网络策略

`network_policies` 接收 Kubernetes NetworkPolicy **spec** 列表（JSON），平台设置所属题目的选择器。示例仅允许发往指定私网网段的出站流量：

```json
[{"policyTypes":["Egress"],"egress":[{"to":[{"ipBlock":{"cidr":"10.0.0.0/8"}}]}]}]
```

这不是通用“禁止互联网”模板：需要按集群 DNS、题目网络和依赖地址调整；实际隔离由 CNI 执行。创建/更新题目会触发镜像预热，详细机制见[工作负载调度](../deploy/workloads)。
