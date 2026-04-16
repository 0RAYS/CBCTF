---
sidebar_position: 4
---

# 动态附件生成

## 附件类型

| 类型 | 适用题型 | 说明 |
|------|---------|------|
| 静态附件 | 问答题、静态题、容器题 | 上传一次, 所有队伍共享 |
| 动态附件 | 动态附件题 | 每队单独生成, 依赖 Kubernetes |

## 真实工作流程

动态附件不会根据配置自动常驻启动生成器池。当前实现是: 

1. 管理员先在后台启动生成器 Pod
2. 队伍初始化或重置题目时, 平台从可用生成器中随机选择一个
3. 平台在容器内执行 `./run.sh {team_id} {flags}`
4. 等待输出文件写入共享卷
5. 将 `{team_id}.zip` 提供给队伍下载

如果没有预热生成器, 动态附件题会因“无可用生成器”而无法初始化。

## 生成器要求

### 容器启动方式

平台会以 `sleep infinity` 覆盖镜像入口, 保持生成器容器常驻。因此镜像内必须存在 `sleep` 命令。

### 可选上传 `generator.zip`

若题目上传了 `generator.zip`, 平台会在容器内执行: 

```bash
unzip /root/mnt/generator.zip -d /root
```

因此镜像还必须包含 `unzip`。

### 执行入口

平台执行命令: 

```bash
./run.sh {team_id} {base64(base64(flag1),base64(flag2),...)}
```

参数说明: 

- `$1`: 队伍 ID
- `$2`: 多 Flag 的二次 Base64 编码结果

### 输出路径

生成结果必须写入: 

```text
/root/mnt/attachments/{team_id}.zip
```

平台只会读取这个固定路径。

## 示例脚本

```bash
#!/bin/bash
TEAM_ID=$1
FLAGS_B64=$2

FLAGS=$(echo "$FLAGS_B64" | base64 -d | tr ',' '\n' | while read f; do echo "$f" | base64 -d; done)
FLAG1=$(echo "$FLAGS" | head -1)

mkdir -p /tmp/challenge
echo "$FLAG1" > /tmp/challenge/flag.txt

mkdir -p /root/mnt/attachments
zip -j /root/mnt/attachments/${TEAM_ID}.zip /tmp/challenge/*
```

## 注意事项

1. 镜像必须包含 `sleep` 与 `unzip`。
2. 脚本入口固定为 `/root/run.sh` 或当前工作目录下的 `./run.sh`。
3. 输出文件必须是 `/root/mnt/attachments/{team_id}.zip`。
4. 动态 Flag 的实际结果可能与模板字符表现不同, 题目逻辑不要依赖固定字符形态。
5. 共享存储不可用时, 动态附件无法生成。

## 管理入口

生成器支持两种预热方式: 

- 全局生成器: `/admin/generators`
- 某场比赛专用生成器: `/admin/contests/{contestID}/generators`

建议在比赛开始前预热动态附件题所需的生成器实例。
