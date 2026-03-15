---
sidebar_position: 2
---

# Docker 部署

Docker 部署适用于本地体验和小规模测试，无需 Kubernetes 集群（若不使用动态附件题和容器题）。

## 前置要求

- Docker 24+
- Docker Compose v2（`docker compose` 命令）
- （可选）Kubernetes 集群 kubeconfig，用于动态附件题目和容器题目

## 准备文件

创建工作目录，并准备以下文件：

### docker-compose.yaml

```yaml
services:
  cbctf:
    image: ghcr.io/0rays/cbctf:latest
    ports:
      - "8000:8000"
    volumes:
      - ./config.yaml:/app/config.yaml
      - ./data:/app/data
      - ./admin.yaml:/app/admin.yaml   # K8s kubeconfig，不使用容器题可省略
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_started
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: cbctf
      MYSQL_USER: cbctf
      MYSQL_PASSWORD: cbctf_password
    volumes:
      - mysql_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass redis_password
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  mysql_data:
  redis_data:
```

### config.yaml

```yaml
host: http://your.domain.com    # 平台对外访问地址

path: /app/data

gin:
  host: 0.0.0.0
  port: 8000
  cors:
    - http://your.domain.com
  jwt:
    secret: change-this-to-a-random-string   # 必须修改！

gorm:
  mysql:
    host: mysql
    port: 3306
    user: cbctf
    pwd: cbctf_password
    db: cbctf

redis:
  host: redis
  port: 6379
  pwd: redis_password

k8s:
  config: /app/admin.yaml
  namespace: cbctf
```

完整配置说明见 [配置说明](/docs/deploy/settings)。

## 挂载 kubeconfig

若需使用动态附件题或容器题，将 Kubernetes 集群的 kubeconfig 复制为工作目录下的 `admin.yaml`：

```bash
cp ~/.kube/config ./admin.yaml
```

若不使用容器相关功能，可从 `docker-compose.yaml` 的 volumes 中删除 `./admin.yaml:/app/admin.yaml` 行。

## 启动

```bash
docker compose up -d
```

## 查看初始密码

```bash
docker compose logs cbctf | grep "Init Admin"
```

输出示例：

```
Init Admin: Admin{ name: admin, password: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx, email: admin@0rays.club}
```

登录地址：`http://your.domain.com/platform/#/login`

## 常见问题

**数据目录权限问题**

若 cbctf 容器无法写入 `./data` 目录，检查宿主机目录权限：

```bash
chmod 777 ./data
```

**kubeconfig 路径错误**

确认 `config.yaml` 中 `k8s.config` 路径与挂载路径一致，默认为 `/app/admin.yaml`。
