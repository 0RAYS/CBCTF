---
sidebar_position: 4
---

# 二进制部署

二进制部署适用于已有基础设施（MySQL、Redis、K8s）且需要最大灵活性的场景。

## 前置要求

| 组件 | 版本要求 | 说明 |
|------|---------|------|
| Go | 1.26+ | 编译后端 |
| Node.js | 22+ | 编译前端 |
| pnpm | 10+ | 前端包管理器 |
| libpcap-dev | 任意 | 流量捕获功能（不需要时可跳过） |
| MySQL | 8.0+ | 主数据存储 |
| Redis | 6.0+ | 缓存与任务队列 |
| Kubernetes | 1.20+ | 动态附件题目和容器题目（可选） |

## 编译

### 1. 克隆仓库

```bash
git clone https://github.com/0RAYS/CBCTF.git
cd CBCTF
```

### 2. 编译前端

```bash
cd frontend
pnpm install
pnpm build
cd ..
```

构建产物位于 `frontend/dist/`，将通过 `go:embed` 嵌入 Go 二进制。

### 3. 编译后端

**含流量捕获功能（推荐，需 libpcap）**

```bash
# Ubuntu/Debian
sudo apt install -y libpcap-dev

# Alpine
apk add libpcap-dev build-base

CGO_ENABLED=1 go build \
  -ldflags="-linkmode external -extldflags '-static' -s -w" \
  -trimpath -o cbctf .
```

**不含流量捕获功能（无需 CGO）**

```bash
go build -ldflags="-s -w" -trimpath -o cbctf .
```

## 配置

### 首次运行

首次运行时，若当前目录下不存在 `config.yaml`，程序会自动生成默认配置并退出：

```bash
./cbctf
```

按照实际环境编辑生成的 `config.yaml`，然后重新运行。

### 最小配置示例

```yaml
host: https://ctf.example.com    # 平台对外访问地址

path: ./data                      # 数据存储目录

gin:
  host: 0.0.0.0
  port: 8000
  cors:
    - https://ctf.example.com
  jwt:
    secret: change-this-to-a-random-string   # 必须修改！

gorm:
  mysql:
    host: 127.0.0.1
    port: 3306
    user: cbctf
    pwd: your-mysql-password
    db: cbctf

redis:
  host: 127.0.0.1
  port: 6379
  pwd: your-redis-password

k8s:
  config: ./admin.yaml      # K8s kubeconfig，不使用容器题可省略
  namespace: cbctf
```

完整配置字段说明见 [配置说明](/docs/deploy/settings)。

## 运行

```bash
./cbctf
```

查看初始管理员密码：

```bash
./cbctf 2>&1 | grep "Init Admin"
# 或查看日志文件（log.save: true 时）
grep "Init Admin" ./logs/*.log
```

## 后台运行（systemd）

创建服务单元文件 `/etc/systemd/system/cbctf.service`：

```ini
[Unit]
Description=CBCTF Platform
After=network.target mysql.service redis.service

[Service]
Type=simple
User=cbctf
WorkingDirectory=/opt/cbctf
ExecStart=/opt/cbctf/cbctf
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

启用并启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now cbctf
# 查看日志
sudo journalctl -u cbctf -f
```
