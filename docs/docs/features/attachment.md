---
sidebar_position: 4
---

# 动态附件生成

## 附件类型

| 类型 | 适用题型 | 说明 |
|------|---------|------|
| 静态附件 | 问答题、静态题、容器题 | 上传 zip，所有队伍共享下载 |
| 动态附件 | 动态附件题 | 每队独立生成，需 Kubernetes |

## 动态附件工作原理

平台在 Kubernetes 集群中为每道动态附件题运行一组生成器 Pod（数量 = `节点数 × k8s.generator_worker`），每个生成器 Pod 处理不同队伍的附件生成请求。

### 生成流程

1. **生成 flag 和 team_id**：平台为该队伍生成 flag，准备 team_id

2. **启动生成器容器**：以 `sleep infinity` 替代镜像的默认 entrypoint，保持容器存活等待指令

3. **上传 generator.zip（可选）**：若出题人在平台上传了 `generator.zip`，平台将其上传至容器的 `/root` 目录并解压。适用于依赖环境相同但脚本不同的场景，避免每次改脚本都要重新构建镜像

4. **执行生成脚本**：平台在容器内执行：
   ```bash
   /root/run.sh {team_id} {base64(base64(flag1),base64(flag2),...)}
   ```
   脚本必须由出题人编写，包含生成附件的完整逻辑

5. **周期性重启容器**：定期重启生成器容器，清理状态

### `/root/run.sh` 接口约定

| 参数 | 说明 |
|------|------|
| `$1` | `team_id`，整型，队伍唯一标识符 |
| `$2` | `base64(base64(flag1) + "," + base64(flag2) + ...)` |

**输出路径（硬性要求）**：`/root/mnt/attachments/{team_id}.zip`

平台只读取此路径的文件并提供给选手下载，其他路径的文件会被忽略。

### 示例脚本

```bash
#!/bin/bash
TEAM_ID=$1
FLAGS_B64=$2

# 解码 flags（逗号分隔的 base64 列表）
FLAGS=$(echo "$FLAGS_B64" | base64 -d | tr ',' '\n' | while read f; do echo "$f" | base64 -d; done)
FLAG1=$(echo "$FLAGS" | head -1)

# 生成附件
mkdir -p /tmp/challenge
echo "$FLAG1" > /tmp/challenge/flag.txt
# ... 其他题目文件生成逻辑 ...

# 打包输出（路径固定）
mkdir -p /root/mnt/attachments
zip -j /root/mnt/attachments/${TEAM_ID}.zip /tmp/challenge/*
```

## 注意事项

1. **容器须包含 `sleep` 和 `unzip` 命令**（平台内部使用）

2. **输出路径固定**：`/root/mnt/attachments/{team_id}.zip`，不可更改

3. **脚本入口固定**：`/root/run.sh`，不可更改

4. **动态 flag 长度可变**：动态 flag 生成结果长度可能与模板不同，题目设计中不要依赖固定 flag 长度

5. **镜像 platform**：建议在 Dockerfile 中指定 `--platform linux/amd64`，避免架构不匹配问题

## 示例题目

参考官方示例：[dynamic-attachment](https://github.com/0RAYS/CBCTF/tree/main/example/dynamic)
