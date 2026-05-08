---
sidebar_position: 7
---

# 常见问题

## 服务启动失败

排查顺序：

1. 查看应用日志：`kubectl logs -n cbctf deployment/cbctf`。
2. 检查 PostgreSQL 和 Redis Pod 是否 Ready。
3. 检查 Secret 中的数据库、Redis、JWT 配置是否正确。
4. 检查应用 ServiceAccount 是否绑定了 Chart 创建的 ClusterRole。
5. 如果启用 VPC 外部网络，检查 Kube-OVN、Multus 和节点标签。

## 数据库连接失败

常见原因：

- PostgreSQL Pod 未启动或 Service 名称不可解析。
- `postgres.auth.password` 与 Secret 中的密码不一致。
- 数据库 PVC 损坏、容量不足或权限异常。
- 外部数据库要求 SSL，但 `gorm.postgres.sslmode` 未设置为 `true`。

应用启动时会自动迁移数据库。如果迁移失败，服务会退出，需要先修复数据库连接或权限问题。

## Redis 连接失败

Redis 用于缓存和 Asynq 任务队列。连接失败会影响登录状态、限流、后台任务、邮件、Webhook、动态附件和靶机操作。

排查方向：

- Redis Pod、Service、Secret 是否正常。
- `redis.auth.password` 是否与运行中 Redis 一致。
- 网络策略或命名空间 DNS 是否阻断访问。
- Redis PVC 是否容量不足。

## 前端页面无法访问

平台页面默认在 `/platform/` 下，根路径 `/` 会重定向到 `/platform`。前端使用 HashRouter，登录页路径是 `/platform/#/login`。

排查方向：

- 访问路径是否包含 `/platform/`。
- Ingress 是否把 `/` 路径转发到应用 Service。
- `cbctf.host` 是否为实际外部地址，且没有尾部 `/`。
- 浏览器控制台是否有 CORS 或资源加载错误。

## Helm 部署后无法访问

排查方向：

- `kubectl get ingress -n cbctf` 查看 Ingress 是否生成。
- 检查 IngressClass、域名解析、TLS Secret 和负载均衡地址。
- 检查 Service 端口是否为 `8000`，应用容器端口是否为 `cbctf.gin.port`。
- 使用 `kubectl port-forward -n cbctf svc/cbctf 8000:8000` 判断应用本身是否正常。

## 登录或 OAuth 跳转异常

排查方向：

- `cbctf.host` 必须是用户浏览器可访问的外部地址。
- OAuth 提供商回调地址应指向后端的 `/oauth/{provider}/callback`。
- 如果前后端分离，当前 OAuth 登录完成后仍默认跳回后端 `/platform/#/oauth/callback` 页面。
- 检查 `gin.cors` 是否包含前端来源。

## 上传失败

排查方向：

- `gin.upload.max` 是否小于上传文件大小。
- Ingress 或反向代理是否有更小的 body size 限制。
- `/app/data` PVC 是否可写、容量是否充足。
- 管理员或用户是否具备对应上传权限。

## 靶机无法启动

排查方向：

- 应用日志和任务日志是否有 Kubernetes API 错误。
- ServiceAccount 是否有创建 Pod、Service、NetworkPolicy、Job、Multus、Kube-OVN 资源的权限。
- 题目镜像是否存在、tag 是否正确、节点是否能拉取。
- 集群 CPU、内存、存储或配额是否不足。
- 题目端口、资源、环境变量和 Flag 注入配置是否正确。
- VPC 模式是否安装 Kube-OVN/Multus，并配置节点标签。
- FRP 模式是否配置可用 FRPS、token 和端口池。

## 镜像拉取失败

排查方向：

- 镜像仓库地址、tag、平台架构是否正确。
- 私有仓库是否配置 `imagePullSecrets` 或 `imageCredentials`。
- 节点是否可以访问镜像仓库。
- 是否达到镜像仓库限流。

## 邮件发送失败

排查方向：

- 后台 SMTP 配置是否正确。
- Redis 和 Asynq 邮件队列是否正常。
- SMTP 服务是否允许集群出口访问。
- 发件账号、密码、TLS/端口设置是否符合服务商要求。

## Webhook 失败

排查方向：

- Webhook 目标是否在 `webhook.whitelist` 允许范围内。
- 集群是否能访问目标地址。
- 目标服务是否返回非 2xx 状态。
- 查看后台 Webhook 历史记录和任务日志。

## 安全建议

- 首次登录后立即修改 `admin` 密码。
- 生产环境设置固定、足够长的 `gin.jwt.secret`。
- 数据库、Redis、SMTP、OAuth、FRP 使用强密码或随机密钥。
- 使用 HTTPS，并正确配置 Ingress TLS。
- 限制管理后台访问来源，按角色授予最小权限。
- Kubernetes RBAC 不要授予超过平台运行所需的权限。
- 为题目容器设置 CPU、内存、网络和命名空间隔离。
- 定期备份 PostgreSQL 和 `/app/data`。
- 不要公开 `config.yaml`、kubeconfig、Secret、私有镜像凭据和含敏感信息的日志。
