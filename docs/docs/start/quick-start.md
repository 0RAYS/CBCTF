---
sidebar_position: 2
---

# 快速上手

## 部署方式对比

| 方式 | 适用场景 | 优点 | 缺点 |
|------|---------|------|------|
| **Docker** | 本地体验、小规模测试 | 零依赖、5 分钟上线 | 不支持高可用，资源有限 |
| **Helm（推荐）** | 生产级 Kubernetes 部署 | 自动管理 MySQL/Redis，支持完整功能 | 需要 K8s 集群 |
| **二进制** | 自定义部署、已有基础设施 | 最灵活 | 需手动维护依赖 |

如果您只是想快速体验平台功能，选择 **Docker**；如果需要承办正式比赛，推荐使用 **Helm** 部署在 Kubernetes 集群上。

---

## Docker 快速部署

详细说明见 [Docker 部署](/docs/deploy/docker)，以下为极简流程：

```bash
# 1. 准备 docker-compose.yaml 和 config.yaml（见 Docker 部署文档）
# 2. 复制 K8s kubeconfig（如使用容器题功能）
cp ~/.kube/config ./admin.yaml
# 3. 启动服务
docker compose up -d
# 4. 查看初始管理员密码
docker compose logs cbctf | grep "Init Admin"
# 5. 浏览器访问
#    http://localhost:8000/platform/#/login
```

---

## Helm 快速部署

详细说明见 [Helm 部署](/docs/deploy/helm)，以下为极简流程：

```bash
# 1. 添加 Chart 仓库
helm repo add 0rays https://0rays.github.io/CBCTF-charts
helm repo update
# 2. 安装（最小配置）
helm install cbctf 0rays/cbctf \
  --set cbctf.host=https://your.domain.com \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=your.domain.com
# 3. 查看初始管理员密码
kubectl logs deployment/cbctf | grep "Init Admin"
```

---

## 初始登录

所有部署方式在首次启动时，若数据库中没有管理员账号，平台会自动创建并将凭据打印到日志：

```
Init Admin: Admin{ name: admin, password: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx, email: admin@0rays.club}
```

- **Docker**：`docker compose logs cbctf | grep "Init Admin"`
- **Helm**：`kubectl logs deployment/cbctf | grep "Init Admin"`
- **二进制**：查看终端输出或 `./logs/` 目录下的日志文件

登录地址：`https://your.domain/platform/#/login`

:::warning
首次登录后请立即修改初始密码，并检查 `gin.jwt.secret` 配置。
:::
