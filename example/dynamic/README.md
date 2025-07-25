# 🔧 Generator 镜像模板说明

本模板用于构建 **动态题目的附件生成器 Docker 镜像**，以支持自动化部署和分发个性化附件内容。以下示例基于 **密码学题目** 的生成器构建过程。

---

## 🧩 工作流程

生成器逻辑的 Golang 实现位于 `/internal/k8s/generator.go`，建议结合代码阅读，以加深理解。

1. **初始化：**
   - 获取当前用户所属队伍的 `ID`（作为附件唯一标识）
   - 生成对应的 `flag`

2. **启动容器：**
   - 使用 `sleep infinity` 命令保持容器持续运行
   - 请确保容器中存在 `sleep` 命令

3. **可选解压步骤（若存在 `generator.zip`）：**
   - 将主机中的 `generator.zip` 拷贝至容器 `/root` 目录
   - 解压至 `/root`

4. **生成附件：**
   - 容器内执行命令：
     ```bash
     ./run.sh id base64(base64(flag1),base64(flag2),...)
     ```
   - 生成附件位置为：`/root/mnt/attachments/{id}.zip`

5. **将附件复制回主机：**
   ```bash
   cp /root/mnt/attachments/{id}.zip /path/in/host/{id}.zip
   ```

6. **删除 Pod**

---

## 📦 示例流程

假设用户所属队伍 ID 为 `1`，生成的 flag 为 `flag{test}`：

```bash
# 启动容器，保持运行
sleep infinity 

# generator.zip 存在时，复制并解压
cp /host/path/generator.zip /root/generator.zip
unzip /root/generator.zip -d /root

# 生成附件
./run.sh 1 ZmxhZ3R7dGVzdH0=

# 拷贝回主机
cp /root/mnt/attachments/1.zip /host/path/1.zip
```

---

## 🛠️ Dockerfile 示例

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

---

## 🏃 run.sh 示例

```bash
#!/bin/bash

# $1: team_id
# $2: base64(base64(flag1),base64(flag2),...)

# 生成附件
python3 generator.py $1 $2
```

---

## 🧳 generator.zip 文件说明

该文件为题目生成器的可选依赖，用于 **同一基础镜像下复用不同的生成脚本**，典型用于密码学题目中快速替换 `generator.py` 和 `template.py` 等内容。

- 上传此文件由管理员完成
- 当不存在该文件时，镜像需能 **独立生成附件**

---

## ⚠️ 注意事项

1. **容器中必须存在的命令：**
   - `sleep`：用于保持运行
   - `unzip`：用于解压文件
   - `tee`：用于复制文件至主机

2. **必须存在的文件：**
   - `/root/run.sh`

3. **`run.sh` 规范：**
   - 接收两个参数：
      - `$1`：队伍 ID，用作唯一标识
      - `$2`：格式为 `base64(base64(flag1),base64(flag2),...)`

4. **附件生成输出要求：**
   - 所有内容必须压缩为 `/root/mnt/attachments/{id}.zip`

5. **性能要求：**
   - 不得存在耗时命令
   - 推荐使用字节级操作（如模板替换）替代编译流程

6. **错误处理要求：**
   - 附件生成流程不得产生任何错误或中断

7. **flag 支持：**
   - flag 可能被变异为随机字符组合（参考 `/internal/utils/flag.go` 中 `RandFlag` 函数）
   - 请确保附件能正确处理 **不定长 flag**

8. **版本号规范：**
   - 禁止使用 `latest` 标签（由于 k8s 的 `IfNotPresent` 策略）
   - 正确示例：
     ```bash
     docker build -t generator-test:20250221 .
     docker tag generator-test:20250221 docker.0rays.club/test/generator:20250221
     docker push docker.0rays.club/test/generator:20250221
     ```

9. **流程必须完整支持：**
   ```bash
   cp /host/path/generator.zip /root/generator.zip  # 可选
   unzip /root/generator.zip -d /root               # 可选
   ./run.sh id base64(base64(flag1),base64(flag2),...)
   cp /root/mnt/attachments/{id}.zip /host/path/{id}.zip
   ```
