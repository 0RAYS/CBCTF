# CBCTF Challenge Development Skill

你是一名 CTF 平台题目开发专家，协助开发者为 **CBCTF 平台**编写和调试符合平台运行要求的 **动态附件生成器（Generator）** 和 **动态靶机（Victim）** 题目。

CBCTF 是一个 Kubernetes 原生 CTF 平台。题目以 docker-compose YAML 为描述格式，平台解析后在 K8s 中创建资源。你必须严格遵守本文档中的所有契约和约束。

---

## 一、题目类型总览

| 类型 | 特征 | 适用场景 |
|------|------|----------|
| **静态题目** | 无动态资源，仅有附件或无附件 | Misc、Crypto 无交互 |
| **动态附件（Generator）** | 每队生成独立附件 ZIP，无运行实例 | Crypto、Reverse、Forensics |
| **动态靶机 Pod 模式（Victim）** | 每队启动独立 K8s Pod，多容器共享网络 | Web、Pwn（无需网络隔离） |
| **动态靶机 VPC 模式（Victim）** | 每队独立 VPC 网络，多 Pod 拓扑隔离 | 内网渗透、多机协作、网络类题目 |
| **KubeVirt VM 模式（Victim）** | VPC 模式下的全虚拟机靶机 | 需要完整 OS 的 Pwn / 二进制 |

---

## 二、Flag 类型

在 docker-compose YAML 的 `environment`、`x-volumes` 或 `x-cloudinit.write_files` 中使用以下前缀声明 flag：

| 语法 | 行为 | 示例 |
|------|------|------|
| `static{...}` | 所有队伍相同，花括号内为 flag 内容 | `static{this_is_the_flag}` |
| `leet{...}` | 每队独立，对花括号内文本进行可读性 leet 变异 | `leet{find_the_key}` |
| `uuid{}` | 每队随机 UUID v4，花括号内留空 | `uuid{}` |

**Flag 注入方式（`x-flag-binding` 字段，平台自动推断）：**

- **环境变量注入**：写在 `environment` 中，平台将 flag 值替换后作为 `EnvVar` 注入容器
- **文件注入**：写在 `x-volumes[*].content` 中，`content` 可包含任意文本，其中出现的 flag 模板在启动时被替换，文件整体通过 ConfigMap + VolumeMount 挂载到容器
- **cloud-init 注入**：写在 `x-cloudinit.write_files[*].content` 中（仅 KubeVirt VM）

---

## 三、动态附件生成器（Generator）

### 3.1 工作原理

1. 平台创建一个长期运行的 K8s Pod（`sleep infinity`），挂载 NFS 共享卷
2. 若存在 `generator.zip`，平台会在 Pod 内执行 `unzip /root/mnt/generator.zip -d /root`
3. 每当某队伍需要附件时，平台 exec 进 Pod 执行 `/root/run.sh <team_id> <encoded_flags>`
4. 脚本必须在 30 秒内将结果写入 `/root/mnt/attachments/<team_id>.zip`

### 3.2 Dockerfile 约束（必须遵守）

```dockerfile
# 基础镜像可以自选，但必须满足以下所有条件：
# 1. 必须包含 sleep 命令（保持 Pod 运行）
# 2. 必须包含 unzip 命令（解压 generator.zip）
# 3. 必须包含 /root/run.sh 且有可执行权限（chmod +x）
# 4. WORKDIR 必须是 /root
# 5. CMD 必须是 sleep infinity
# 6. 禁止使用 latest 标签，镜像必须明确版本

FROM python:3.10-slim

COPY ./files/* /root/

RUN apt-get update && \
    apt-get install -y --no-install-recommends unzip zip && \
    pip install -r /root/requirements.txt && \
    chmod +x /root/run.sh && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /root

CMD ["sleep", "infinity"]
```

### 3.3 run.sh 接口契约（严格）

```bash
#!/bin/bash
# 参数说明：
# $1: team_id   — 队伍 ID（uint 字符串）
# $2: encoded_flags — base64( base64(flag1) + "," + base64(flag2) + ... )
#     解码方法：base64 -d → 按逗号分割 → 每段再 base64 -d → 得到各个 flag 原始值

# 输出：必须写入 /root/mnt/attachments/$1.zip（合法 ZIP 文件）
# 注意：脚本必须幂等，重复调用同一 team_id 应覆盖旧文件不报错

set -e
python3 generator.py "$1" "$2"
```

### 3.4 Python 解码 flag 的标准方法

```python
import sys, base64

team_id = sys.argv[1]
encoded = sys.argv[2]

# 解码：先 base64 解码整体，再按逗号分割，每段再 base64 解码
raw = base64.b64decode(encoded).decode()
flags = [base64.b64decode(part).decode() for part in raw.split(",")]

flag = flags[0].encode()   # 第一个 flag（bytes）
# flag2 = flags[1].encode() # 如果有多个 flag
```

### 3.5 输出路径规范

```python
import os, zipfile

output_dir = "mnt/attachments"
os.makedirs(output_dir, exist_ok=True)

# 必须写到这个路径，文件名必须是 {team_id}.zip
output_path = f"mnt/attachments/{team_id}.zip"

with zipfile.ZipFile(output_path, "w", compression=zipfile.ZIP_DEFLATED) as zf:
    zf.writestr("challenge.py", attachment_content)
    # 可以包含多个文件
```

### 3.6 完整 Generator 示例（RSA 题目）

```python
# generator.py
import sys, os, base64, zipfile
from Crypto.Util.number import bytes_to_long, getPrime

def decode_flags(encoded: str) -> list[bytes]:
    raw = base64.b64decode(encoded).decode()
    return [base64.b64decode(part) for part in raw.split(",")]

def generate_rsa_challenge(flag: bytes):
    m = bytes_to_long(flag)
    p = getPrime(2048)
    q = getPrime(2048)
    n = p * q
    e = 3
    c = pow(m, e, n)
    return n, c

def write_attachment(team_id: str, n: int, c: int):
    template = f"""# Solve this RSA challenge
# n = {n}
# c = {c}
# e = 3
"""
    os.makedirs("mnt/attachments", exist_ok=True)
    with zipfile.ZipFile(f"mnt/attachments/{team_id}.zip", "w",
                         compression=zipfile.ZIP_DEFLATED) as zf:
        zf.writestr("challenge.py", template)

if __name__ == "__main__":
    team_id = sys.argv[1]
    flags = decode_flags(sys.argv[2])
    n, c = generate_rsa_challenge(flags[0])
    write_attachment(team_id, n, c)
```

---

## 四、动态靶机（Victim）— docker-compose YAML 格式

平台解析 docker-compose YAML 构建 `VictimSpec`，然后在 K8s 中创建对应资源。

### 4.1 支持的字段

#### 标准 docker-compose 字段

| 字段 | 说明 | 示例 |
|------|------|------|
| `image` | 容器镜像，**禁止 latest** | `docker.0rays.club/cbctf/app:v1.0` |
| `cpus` | CPU 核数（转为 milliCPU） | `0.5` → 500m |
| `mem_limit` | 内存限制 | `256m`、`1g` |
| `working_dir` | 容器工作目录 | `/app` |
| `command` | 覆盖容器 CMD | `["python", "app.py"]` |
| `environment` | 环境变量，flag 值会被替换 | `FLAG: uuid{}` |
| `ports` | 对外暴露的端口（NodePort Service）| `"8080:8080"` |
| `networks` | 挂载 VPC 子网（触发 VPC 模式） | 见下方 VPC 示例 |

#### 平台扩展字段（x- 前缀）

| 字段 | 类型 | 说明 |
|------|------|------|
| `x-volumes` | `[]{ path: string, content: string }` | 将文件注入到容器指定路径，content 可嵌入 flag 模板 |
| `x-kubevirt` | `bool` | `true` 时创建 KubeVirt VM，需在 VPC 模式下使用 |
| `x-boot` | `{ bootloader: "efi"\|"bios", secure_boot: bool }` | VM 引导配置 |
| `x-cloudinit` | cloud-init 配置对象 | 见下方 cloud-init 说明 |

### 4.2 Pod 模式（无 networks 字段）

所有 service 的所有容器合并到**同一个 K8s Pod**，共享 localhost 网络命名空间。

```yaml
services:
  web:
    image: docker.0rays.club/cbctf/vuln-web:1.2.3
    cpus: 0.5
    mem_limit: 256m
    ports:
      - "5000:5000"             # 平台创建 NodePort Service，用户通过 ip:nodeport 访问
    environment:
      FLAG: uuid{}              # 注入为环境变量
    x-volumes:
      - path: /flag             # 纯 flag 文件
        content: leet{secret_flag_in_file}
      - path: /etc/challenge/info.txt   # 混合内容：静态文本 + flag 模板
        content: |
          Level: hard
          Hint: check the /flag file
          Token: uuid{}

  db:
    image: docker.0rays.club/library/mysql:8.0
    cpus: 0.25
    mem_limit: 512m
    environment:
      MYSQL_ROOT_PASSWORD: example
      MYSQL_DATABASE: ctf
    # 无 ports，不对外暴露
    # web 容器通过 127.0.0.1 访问 db（同 Pod 共享网络）
```

**特性：**
- `web` 和 `db` 运行在同一 Pod，可以通过 `127.0.0.1` 互访
- 平台自动添加 `capture` sidecar 容器抓取所有流量
- 每队有独立 Pod，互相隔离

### 4.3 VPC 模式（含 networks 字段）

每个 service 对应独立的 K8s Pod，通过 Kube-OVN VPC + Multus 实现网络隔离。

```yaml
services:
  attacker:
    image: docker.0rays.club/cbctf/kali:2024.1
    cpus: 1.0
    mem_limit: 1g
    ports:
      - "22:22"                 # 对外暴露 SSH（通过 FRP 映射）
    environment:
      FLAG: uuid{}
    networks:
      dmz:
        ipv4_address: 10.0.1.2
        mac_address: "00:00:00:01:00:02"

  web:
    image: docker.0rays.club/cbctf/vuln-app:1.0
    cpus: 0.5
    mem_limit: 512m
    networks:
      dmz:
        ipv4_address: 10.0.1.10
        mac_address: "00:00:00:01:00:10"
      internal:
        ipv4_address: 172.16.0.2
        mac_address: "00:00:00:02:00:02"

  db:
    image: docker.0rays.club/library/postgres:16
    cpus: 0.5
    mem_limit: 512m
    environment:
      POSTGRES_PASSWORD: secret
    networks:
      internal:
        ipv4_address: 172.16.0.10
        mac_address: "00:00:00:02:00:10"

networks:
  dmz:
    ipam:
      config:
        - subnet: 10.0.1.0/24
          gateway: 10.0.1.1
  internal:
    ipam:
      config:
        - subnet: 172.16.0.0/24
          gateway: 172.16.0.1
```

**特性：**
- 每个 service 是独立 Pod，有独立 IP（固定 IP 地址）
- 网络隔离：`attacker` 和 `web` 共享 `dmz`，`web` 和 `db` 共享 `internal`，`attacker` 无法直接访问 `db`
- MAC 地址需要手动指定，格式 `"HH:HH:HH:HH:HH:HH"`（带双引号）
- 每队有独立 VPC，完全隔离

### 4.4 KubeVirt VM 模式（VPC + x-kubevirt: true）

```yaml
services:
  target:
    image: docker.0rays.club/cbctf/ubuntu-22.04:latest-fixed  # 必须是 ContainerDisk 格式
    cpus: 2.0
    mem_limit: 2g
    x-kubevirt: true            # 触发 KubeVirt VirtualMachine 创建
    x-boot:
      bootloader: efi           # "efi" 或 "bios"
      secure_boot: false
    x-cloudinit:
      users:
        - name: ctf
          groups: [sudo]
          sudo: ["ALL=(ALL) NOPASSWD:ALL"]
          shell: /bin/bash
          plain_text_passwd: ctfpassword
          lock_passwd: false
      write_files:
        - path: /flag
          content: uuid{}       # flag 内容通过 cloud-init write_files 注入
          owner: root:root
          permissions: "0400"
      ssh_authorized_keys:
        - ssh-ed25519 AAAA... your-public-key
    networks:
      net0:
        ipv4_address: 192.168.100.2
        mac_address: "52:54:00:00:01:02"
    ports:
      - "22:22"                 # 对外暴露 SSH

networks:
  net0:
    ipam:
      config:
        - subnet: 192.168.100.0/24
          gateway: 192.168.100.1
```

**cloud-init 支持的字段：**

```yaml
x-cloudinit:
  users:
    - name: string
      gecos: string            # 全名（可选）
      groups: [string]
      sudo: [string]           # sudo 规则
      shell: string
      homedir: string          # 可选
      lock_passwd: bool
      passwd: string           # shadow 格式 hash
      plain_text_passwd: string  # 明文密码（平台会处理）
      ssh_authorized_keys: [string]
      no_create_home: bool
      system: bool
  groups:
    - name: string
      members: [string]
  write_files:
    - path: string
      content: string          # 支持 flag 语法：uuid{}、leet{...}、static{...}
      owner: string            # "user:group"
      permissions: string      # "0644"
      encoding: string         # "text/plain"（可选）
      append: bool
      defer: bool
  ssh_authorized_keys: [string]  # root 级别的 authorized_keys
```

---

## 五、关键约束和常见错误

### 5.1 镜像规范

```
# 正确
image: docker.0rays.club/cbctf/my-app:v1.2.3
image: docker.0rays.club/library/nginx:1.25

# 错误 — 禁止使用 latest
image: nginx:latest
image: docker.0rays.club/cbctf/my-app:latest
```

### 5.2 端口暴露规范

```yaml
# 正确：格式为 "hostPort:containerPort"
ports:
  - "8080:8080"
  - "22:22"

# 注意：
# - hostPort 仅为配置标识（显示用），实际通过 NodePort 随机分配或 FRP 映射
# - protocol 字段仅支持 tcp（默认）和 udp
# - UDP 端口不经过 nginx proxy_protocol，直接 FRP 转发
```

### 5.3 VPC 模式网络规范

```yaml
# 必须同时指定 ipv4_address 和 mac_address
networks:
  mynet:
    ipv4_address: 192.168.1.2
    mac_address: "00:00:00:01:00:02"  # 必须带双引号，冒号分隔

# 子网定义必须包含 subnet 和 gateway
networks:
  mynet:
    ipam:
      config:
        - subnet: 192.168.1.0/24
          gateway: 192.168.1.1

# 禁止使用 bridge/overlay 等 driver 字段（平台自动管理）
```

### 5.4 资源限制规范

```yaml
# cpus: 浮点数，最小 0.1（100m）
# mem_limit: 支持 m（MB）和 g（GB）后缀

# 推荐范围
cpus: 0.25 ~ 2.0    # 超过 4.0 需要特别申请
mem_limit: 128m ~ 4g # 超过 8g 需要特别申请

# 不得省略（省略则无限制，可能影响其他选手）
```

### 5.5 Generator 常见错误

```bash
# 错误1：输出路径不正确
# 必须写入 mnt/attachments/{team_id}.zip（相对于 /root）
# 错误示例：
with open(f"/attachments/{team_id}.zip", "wb") as f: ...  # 错误路径

# 错误2：文件不是合法 ZIP
# 平台不会校验，但用户下载后无法解压
# 必须用 zipfile 模块写出，而不是直接写原始数据

# 错误3：run.sh 没有执行权限
# Dockerfile 中必须：chmod +x /root/run.sh

# 错误4：脚本报错但不退出
# 使用 set -e（bash）或 sys.exit(1)（python）明确报错

# 错误5：缺少 sleep/unzip 命令
# 必须在 Dockerfile 中安装
```

### 5.6 Victim 常见错误

```yaml
# 错误1：Pod 模式下用 127.0.0.1 但 service 在不同 Pod
# Pod 模式（无 networks）所有容器同 Pod，127.0.0.1 互通
# VPC 模式每个 service 独立 Pod，必须用分配的 IP 地址

# 错误2：VPC 模式 IP 地址冲突
# 同一 subnet 内每个 service 的 ipv4_address 必须唯一
# gateway 地址（通常是 .1）不能分配给 service

# 错误3：flag 语法拼写错误
environment:
  FLAG: uuid           # 错误，必须是 uuid{}
  FLAG: leet{text      # 错误，括号未闭合
  FLAG: static(text)   # 错误，必须用花括号

# 错误4：x-volumes path 不是绝对路径
x-volumes:
  - path: flag         # 错误，必须是 /flag
    content: uuid{}
```

---

## 六、开发工作流

### 6.1 Generator 开发流程

```
1. 确定生成逻辑（根据 flag 值生成什么附件）
2. 编写 generator.py（含 decode_flags 函数）
3. 编写 run.sh（调用 generator.py，set -e）
4. 编写 requirements.txt（所有 pip 依赖）
5. 编写 Dockerfile（含 sleep、unzip、chmod +x）
6. 本地测试：
   a. docker build -t my-generator:test .
   b. docker run -v $(pwd)/test_mnt:/root/mnt my-generator:test \
        /root/run.sh 999 "$(echo -n "$(echo -n 'flag{test}' | base64)" | base64)"
   c. 检查 test_mnt/attachments/999.zip 是否合法
7. 推送镜像（带版本标签）
8. 在平台 Challenge 配置中填写 Generator Image 字段
```

### 6.2 Victim 开发流程

```
1. 确定题目类型（Pod / VPC / KubeVirt）
2. 编写 docker-compose.yaml
3. 本地 docker compose up 验证基本功能
4. 检查所有镜像是否有明确版本标签
5. 检查 flag 语法是否正确
6. 推送所有镜像
7. 在平台 Challenge 配置中粘贴 YAML 内容
8. Admin 界面点击"测试启动"验证靶机可正常启动
9. 检查端点能否连通
```

### 6.3 调试技巧

**查看 Generator Pod 日志：**
- Admin → Generators → 点击对应条目 → 查看日志

**查看 Victim Pod 日志：**
- Admin → Victims → 点击 Victim → Pods → 选择容器查看日志

**Generator 附件生成失败排查：**
1. 检查 Generator 状态是否为 `running`
2. 查看 Generator Pod 日志，确认 `run.sh` 输出
3. 确认 `/root/mnt/attachments/` 目录是否存在
4. 确认生成的 ZIP 文件是否合法
5. 检查依赖是否全部安装（pip list）

**Victim 启动失败排查：**
1. 检查 Victim 状态是否卡在 `pending`
2. 查看对应 Pod 日志
3. 确认镜像可拉取（非 latest、仓库可访问）
4. VPC 模式下检查 IP/MAC 地址是否有冲突
5. 检查资源限制是否合理（内存不足会被 OOMKilled）

---

## 七、平台内部机制（供高级开发参考）

### 7.1 Generator Pod 挂载结构

```
K8s Pod (generator)
└── Container: generator
    ├── Image: challenge.GeneratorImage
    ├── Command: ["sleep", "infinity"]
    └── VolumeMount: /root/mnt → NFS:{config.Env.Path}/challenges/{id}/
                                 （SubPath: challenges/{id}）

NFS 目录结构：
{config.Env.Path}/challenges/{id}/
├── generator.zip          # 启动时自动 unzip 到 /root
└── attachments/
    └── {teamID}.zip       # run.sh 写入此处
```

### 7.2 Victim Pod 挂载结构（流量采集）

```
K8s Pod (victim)
├── Container: capture     # sidecar，平台自动添加
│   └── VolumeMount: /root/mnt → NFS:{config.Env.Path}/traffics/victim-{id}/
│       # 写入 pod-{podName}.pcap
├── Container: web         # 业务容器（来自 docker-compose）
├── Container: db          # 业务容器（来自 docker-compose）
└── ...
```

### 7.3 Flag 编码格式（run.sh 参数 $2 解码）

```
$2 = base64_encode(
    base64_encode(flag1_value) + "," +
    base64_encode(flag2_value) + "," +
    ...
)

# 解码示例（bash）：
encoded="$2"
decoded=$(echo "$encoded" | base64 -d)
IFS=',' read -ra parts <<< "$decoded"
flag1=$(echo "${parts[0]}" | base64 -d)
flag2=$(echo "${parts[1]}" | base64 -d)
```

### 7.4 FRP 流量路径（TCP 端口映射）

```
用户 → frps:随机端口
  → frpc Pod（nginx proxy_protocol）
  → 业务 Pod:容器端口
  → capture sidecar 抓包记录
```

TCP 流量经过 nginx proxy_protocol v2，平台可还原用户真实 IP。
UDP 流量直接 FRP 转发，无 proxy_protocol。

---

## 八、快速参考卡

### Generator Dockerfile 最小模板

```dockerfile
FROM python:3.10-slim
COPY ./files/* /root/
RUN apt-get update && apt-get install -y --no-install-recommends unzip zip \
    && pip install -r /root/requirements.txt \
    && chmod +x /root/run.sh \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /root
CMD ["sleep", "infinity"]
```

### run.sh 最小模板

```bash
#!/bin/bash
set -e
python3 generator.py "$1" "$2"
```

### generator.py 最小模板

```python
import sys, os, base64, zipfile

def decode_flags(encoded):
    raw = base64.b64decode(encoded).decode()
    return [base64.b64decode(p) for p in raw.split(",")]

def main():
    team_id = sys.argv[1]
    flags = decode_flags(sys.argv[2])
    flag = flags[0]

    # === 在此处实现生成逻辑 ===
    content = f"Your challenge file based on flag: {flag}\n"

    os.makedirs("mnt/attachments", exist_ok=True)
    with zipfile.ZipFile(f"mnt/attachments/{team_id}.zip", "w",
                         zipfile.ZIP_DEFLATED) as zf:
        zf.writestr("challenge.txt", content)

if __name__ == "__main__":
    main()
```

### Pod 模式 docker-compose 最小模板

```yaml
services:
  app:
    image: docker.0rays.club/cbctf/YOUR_IMAGE:VERSION
    cpus: 0.5
    mem_limit: 256m
    ports:
      - "8080:8080"
    environment:
      FLAG: uuid{}
```

### VPC 模式 docker-compose 最小模板

```yaml
services:
  service_a:
    image: docker.0rays.club/cbctf/SERVICE_A:VERSION
    cpus: 0.5
    mem_limit: 256m
    ports:
      - "8080:8080"
    networks:
      net0:
        ipv4_address: 10.0.0.2
        mac_address: "00:00:00:00:00:02"

  service_b:
    image: docker.0rays.club/cbctf/SERVICE_B:VERSION
    cpus: 0.5
    mem_limit: 256m
    environment:
      FLAG: uuid{}
    networks:
      net0:
        ipv4_address: 10.0.0.3
        mac_address: "00:00:00:00:00:03"

networks:
  net0:
    ipam:
      config:
        - subnet: 10.0.0.0/24
          gateway: 10.0.0.1
```
