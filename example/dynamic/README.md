# 🔧 动态附件生成器模板

> **本模板用于构建动态题目的附件生成器 Docker 镜像**, 支持自动化部署和分发个性化附件内容。以下示例基于密码学题目的生成器构建过程。

[![Docker](https://img.shields.io/badge/Docker-20.10+-blue.svg)](https://docker.com)
[![Python](https://img.shields.io/badge/Python-3.10+-green.svg)](https://python.org)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.20+-blue.svg)](https://kubernetes.io)

---

## 🧩 工作流程

> 💡 **提示**: 生成器逻辑的 Golang 实现位于 `/internal/k8s/generator.go`, 建议结合代码阅读以加深理解。

### 📋 详细步骤

1. **🚀 初始化阶段**
   - 获取当前用户所属队伍的 `ID`（作为附件唯一标识）
   - 生成对应的 `flag`

2. **🐳 容器启动**
   - 使用 `sleep infinity` 命令保持容器持续运行
   - ⚠️ **重要**: 请确保容器中存在 `sleep` 命令

3. **📦 可选解压步骤**
   - 当存在 `generator.zip` 时: 
     - 将主机中的 `generator.zip` 拷贝至容器 `/root` 目录
     - 解压至 `/root`

4. **⚙️ 附件生成**
   - 容器内执行命令: 
     ```bash
     ./run.sh id base64(base64(flag1),base64(flag2),...)
     ```
   - 生成附件位置: `/root/mnt/attachments/{id}.zip`

5. **📤 文件复制**
   ```bash
   cp /root/mnt/attachments/{id}.zip /path/in/host/{id}.zip
   ```

6. **🧹 清理资源**
   - 删除 Pod

---

## 📦 完整示例流程

假设用户所属队伍 ID 为 `1`, 生成的 flag 为 `flag{test}`: 

```bash
# 1. 启动容器, 保持运行
sleep infinity 

# 2. generator.zip 存在时, 复制并解压
cp /host/path/generator.zip /root/generator.zip
unzip /root/generator.zip -d /root

# 3. 生成附件
./run.sh 1 ZmxhZ3R7dGVzdH0=

# 4. 拷贝回主机
cp /root/mnt/attachments/1.zip /host/path/1.zip
```

---

## 🛠️ 配置文件示例

### 📄 Dockerfile 示例

```dockerfile
FROM docker.0rays.club/test/python:3.10-slim

# 拷贝文件至镜像中
COPY ./files/* /root/
COPY files/requirements.txt /root/requirements.txt

# 安装必要工具和依赖
RUN apt-get update && apt-get install -y unzip zip \
    && pip install -r /root/requirements.txt

# 设置默认工作目录（不可修改）
WORKDIR /root

# 展示目录内容（可随意定义）
CMD ["ls", "-al", "/root"]
```

### 🏃 run.sh 示例

```bash
#!/bin/bash

# $1: team_id
# $2: base64(base64(flag1),base64(flag2),...)

# 生成附件
python3 generator.py $1 $2
```

---

## 🧳 generator.zip 文件说明

该文件为题目生成器的 **可选依赖**, 用于 **同一基础镜像下复用不同的生成脚本**, 典型用于密码学题目中快速替换 `generator.py` 和 `template.py` 等内容。

| 特性 | 说明 |
|------|------|
| **上传方式** | 由管理员完成上传 |
| **依赖关系** | 当不存在该文件时, 镜像需能独立生成附件 |
| **使用场景** | 密码学题目中快速替换生成脚本 |

---

## ⚠️ 重要注意事项

### 🔧 必需组件

#### 容器中必须存在的命令: 
- ✅ `sleep` - 用于保持运行
- ✅ `unzip` - 用于解压文件  

#### 必须存在的文件: 
- ✅ `/root/run.sh`

### 📋 run.sh 规范

| 参数 | 说明 | 示例 |
|------|------|------|
| `$1` | 队伍 ID, 用作唯一标识 | `1` |
| `$2` | 格式为 `base64(base64(flag1),base64(flag2),...)` | `ZmxhZ3R7dGVzdH0=` |

### 📁 附件生成输出要求

- **输出位置**: `/root/mnt/attachments/{id}.zip`
- **格式要求**: 所有内容必须压缩为 ZIP 格式

### ⚡ 性能要求

- ❌ **禁止**: 耗时命令
- ✅ **推荐**: 使用字节级操作（如模板替换）替代编译流程

### 🛡️ 错误处理要求

- ✅ **必须**: 附件生成流程不得产生任何错误或中断
- ✅ **必须**: 确保附件能正确处理不定长 flag

### 🚩 Flag 支持

- **变异支持**: flag 可能被变异为随机字符组合（参考 `/internal/utils/flag.go` 中 `RandFlag` 函数）
- **长度支持**: 请确保附件能正确处理 **不定长 flag**

### 🏷️ 版本号规范

❌ **禁止使用 `latest` 标签**（由于 k8s 的 `IfNotPresent` 策略）

✅ **正确示例**: 
```bash
docker build -t generator-test:20250221 .
docker tag generator-test:20250221 docker.0rays.club/test/generator:20250221
docker push docker.0rays.club/test/generator:20250221
```

### 🔄 完整流程支持

以下流程必须完整支持: 

```bash
cp /host/path/generator.zip /root/generator.zip  # 可选
unzip /root/generator.zip -d /root               # 可选
./run.sh id base64(base64(flag1),base64(flag2),...)
cp /root/mnt/attachments/{id}.zip /host/path/{id}.zip
```
