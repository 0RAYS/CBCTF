---
sidebar_position: 2
---

# 快速上手

推荐使用 Helm 将 CBCTF 部署到 Kubernetes。仓库中的 Dockerfile 用于构建应用镜像；如果需要动态附件或动态靶机，运行环境仍需要 Kubernetes 权限、共享存储和可拉取的题目镜像。

## Helm 快速部署

Chart 位于仓库根目录的 `chart/`。

```bash
helm install cbctf ./chart -n cbctf --create-namespace
kubectl get pods -n cbctf
```

建议复制默认 values 后再安装：

```bash
cp chart/values.yaml my-values.yaml
helm install cbctf ./chart -n cbctf --create-namespace -f my-values.yaml
```

如果暂时没有 Ingress，可先端口转发访问：

```bash
kubectl port-forward -n cbctf svc/cbctf 8000:8000
```

访问地址：`http://127.0.0.1:8000/platform/#/login`。

## 生产部署前建议修改

| 配置项 | 说明 | 示例 |
|---|---|---|
| `cbctf.host` | 平台外部访问地址 | `https://ctf.example.com` |
| `ingress.hosts[0].host` | Ingress 域名 | `ctf.example.com` |
| `cbctf.gin.cors` | 允许的前端访问地址 | `https://ctf.example.com` |
| `cbctf.gin.proxies` | Ingress 或反向代理来源 | `10.244.0.0/16` |
| `cbctf.gin.jwt.secret` | JWT 密钥，生产环境建议显式设置 | `change-me-long-random` |
| `postgres.auth.password` | PostgreSQL 密码，留空会自动生成 | `example-postgres-password` |
| `redis.auth.password` | Redis 密码，留空会自动生成 | `example-redis-password` |
| `persistence.storageClass` | 共享数据卷 StorageClass | `nfs-client` |

不要在 values 文件中写入真实生产密钥后提交到仓库。

## 初始管理员

应用启动时会自动迁移数据库。首次启动且管理员组中没有用户时，会自动创建 `admin` 用户，并把初始密码打印到日志：

```text
Init Admin: Admin{ name: admin, password: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx, email: admin@0rays.club}
```

查看日志：

```bash
kubectl logs -n cbctf deployment/cbctf | grep "Init Admin"
```

首次登录后请立即修改管理员密码。管理后台入口为 `/platform/#/admin`。

## 升级和卸载

```bash
helm upgrade cbctf ./chart -n cbctf -f my-values.yaml
helm uninstall cbctf -n cbctf
```

Chart 创建的共享 PVC 带有保留策略，卸载后不会自动删除数据。如需清空数据，请确认备份后手动删除 PVC。

## 本地开发运行

本地运行需要先准备 PostgreSQL 和 Redis。

```bash
cd frontend
pnpm install
pnpm run build
cd ..
go run .
```

首次运行如果没有 `config.yaml`，程序会生成默认配置并退出。修改数据库、Redis、监听地址和 JWT 密钥后再次启动。

后端构建命令：

```bash
CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o CBCTF .
```

流量捕获能力依赖 CGO 和 libpcap；如果只做文档或前端调试，可优先使用 Helm 或已有镜像。
