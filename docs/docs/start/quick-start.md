---
sidebar_position: 2
---

# 快速上手

## 部署方式

| 方式 | 适用场景 | 优点 | 注意事项 |
|------|---------|------|--------|
| **Helm（推荐）** | 当前唯一支持的部署方式 | 一次性部署应用、PostgreSQL、Redis 与所需 RBAC | 需要 Kubernetes 与 RWX 存储 |

CBCTF **不支持 Docker 部署应用本体**。仓库根目录的依赖示例仅用于本地启动 PostgreSQL 和 Redis。

## Helm 快速部署

详细说明见 [Helm 部署](../deploy/helm.md): 

```bash
helm repo add 0rays https://cbctf.0rays.club/CBCTF
helm repo update

helm install cbctf 0rays/cbctf \
  --namespace cbctf \
  --create-namespace \
  --set cbctf.host=https://your.domain.com \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=your.domain.com

kubectl logs -n cbctf deployment/cbctf | grep "Init Admin"
```

## 初始登录

平台首次启动且数据库中不存在管理员账号时, 会在日志中输出初始管理员凭据: 

```text
Init Admin: Admin{ name: admin, password: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx, email: admin@0rays.club}
```

- Helm: `kubectl logs -n cbctf deployment/cbctf | grep "Init Admin"`

前端访问地址为 `https://your.domain/platform/#/login`。

:::warning
首次登录后请立即修改管理员密码, 并替换 `gin.jwt.secret`、数据库密码与 Redis 密码等默认值。
:::
