---
title: 常见问题
description: 排查 CBCTF 部署、集群组件、容器靶机、动态附件和网络访问问题。
---

# 常见问题

## 服务启动失败

排查顺序：

1. 查看应用日志：`kubectl logs -n cbctf deployment/cbctf`
2. 检查 PostgreSQL 和 Redis Pod 是否 Ready
3. 检查 ConfigMap 中的数据库、Redis、JWT 配置是否正确
4. 检查应用 ServiceAccount 是否同时绑定 Chart 的 Role 和 ClusterRole
5. 如果使用 VPC 网络，检查 Kube-OVN 和 Multus
6. 如果使用 VM 靶机，检查 KubeVirt 组件、CRD 和节点虚拟化能力

## 数据库连接失败

关闭内置 PostgreSQL/Redis 后，请填写相应的 `externalHost`，并检查应用配置中的地址能否从平台 Pod 访问。

常见原因：

- PostgreSQL Pod 未启动或 Service 名称不可解析
- `postgres.auth.password` 与 ConfigMap 中的密码不一致
- 数据库 PVC 损坏、容量不足或权限异常
- 修改了 values 密码，但已有 PostgreSQL 数据目录中的角色密码未变更
- 外部数据库要求 SSL，但 `gorm.postgres.sslmode` 未设置为 `true`

应用启动时会自动迁移数据库。如果迁移失败，服务会退出，需要先修复数据库连接或权限问题。

## Redis 连接失败

Redis 用于缓存和 Asynq 任务队列。连接失败会影响登录状态、限流、后台任务、邮件、Webhook、动态附件和靶机操作。

排查方向：

- Redis Pod、Service 是否正常
- `redis.auth.password` 是否与运行中 Redis 一致
- 网络策略或命名空间 DNS 是否阻断访问
- Redis PVC 是否容量不足

## 前端页面无法访问

平台页面位于 `/platform/`，登录页地址为 `/platform/#/login`。请确认反向代理正确转发该路径及其静态资源。

排查方向：

- 访问路径是否包含 `/platform/`
- Ingress 是否把 `/` 路径转发到应用 Service
- `cbctf.host` 是否为实际外部地址，且没有尾部 `/`
- 浏览器控制台是否有 CORS 或资源加载错误

## Helm 部署后无法访问

排查方向：

- `kubectl get ingress -n cbctf` 查看 Ingress 是否生成
- 检查 IngressClass、域名解析、TLS Secret 和负载均衡地址
- 检查 Service 端口是否为 `8000`，应用容器端口是否为 `cbctf.gin.port`
- 使用 `kubectl port-forward -n cbctf svc/cbctf 8000:8000` 判断应用本身是否正常

## 登录或 OAuth 跳转异常

排查方向：

- `cbctf.host` 必须是用户浏览器可访问的外部地址
- OAuth 提供商回调地址应指向后端的 `/oauth/{provider}/callback`
- 如果前后端分离，当前 OAuth 登录完成后仍默认跳回后端 `/platform/#/oauth/callback` 页面
- 检查 `gin.origins` 是否包含前端页面发起请求时的 `Origin`
- 检查浏览器是否保存认证 Cookie：同源 HTTP 应为 `SameSite=Lax`，跨域登录为 `SameSite=None; Secure`，后者需要 HTTPS
- 新用户登录后接口提示权限不足：检查是否加入了带 `user` 角色的分组；`registration.default_group: 0` 不自动授权
- 注册或加入队伍失败：检查账号是否验证邮箱，SMTP 是否可用

## 修改 values 后配置没有变化

`host`、JWT、注册策略、Gin 运行策略及 Kubernetes 运行参数在首次启动后保存在数据库。检查「系统管理」中读取的实际值，保存后按[生效时机](./settings)重启。数据库连接、监听地址、数据路径和共享 PVC 来自部署文件，应重建 Pod。

## 上传失败

排查方向：

- 对应的 `gin.upload.picture` / `gin.upload.challenge` / `gin.upload.writeup` 是否小于上传请求大小（MiB）
- Ingress 或反向代理是否有更小的 body size 限制
- `/app/data` PVC 是否可写、容量是否充足
- 管理员或用户是否具备对应上传权限
- Writeup 只接受 `.pdf`、`.doc`、`.docx`；ZIP 是管理员导出整场题解的格式

## 靶机无法启动

排查方向：

- 应用日志和任务日志是否有 Kubernetes API 错误
- ServiceAccount 是否有创建、读取、监听和清理 Pod，以及管理 Service、NetworkPolicy、Job、Multus、KubeVirt、Kube-OVN 资源的权限
- 题目镜像是否存在、tag 是否正确、节点是否能拉取
- 集群 CPU、内存、存储或配额是否不足
- 题目端口、资源、环境变量和 Flag 注入配置是否正确
- VPC 模式是否安装 Kube-OVN/Multus，并配置节点标签
- VM 模式是否安装 KubeVirt，service 是否配置 `x-kubevirt: true`、`mem_limit`、VPC 网络、静态 IP 和 MAC 地址
- FRP 模式是否配置可用 FRPS、token 和端口池

## KubeVirt VM 启动失败

排查方向：

- `kubectl get vm -n cbctf` 查看 `VirtualMachine` 是否创建
- `kubectl describe vm <name> -n cbctf` 和 `kubectl describe vmi <name> -n cbctf` 查看 KubeVirt 事件
- 检查 `kubevirt.io/virtualmachines` RBAC 是否包含 `create`、`get`、`deletecollection`
- 检查 VM 镜像是否是可启动的 KubeVirt `containerDisk` 镜像
- 检查节点是否支持虚拟化，KubeVirt 组件和 `virt-launcher` Pod 是否正常
- 检查 compose 是否使用 VPC 网络，并为每张 VM 网卡配置合法的 `ipv4_address` 和 `mac_address`
- 如果 cloud-init 未生效，检查镜像内是否安装并启用 cloud-init，`x-cloudinit.write_files` 路径和权限是否正确

## 启动时报缺少 Kubernetes 权限

如果启动日志出现 `Missing K8s permissions` 或 `Insufficient K8s permissions`，请检查平台 ServiceAccount 是否具有所选题型需要的集群权限。

排查方向：

- 使用 Helm Chart 的 ServiceAccount、Role/RoleBinding 和 ClusterRole/ClusterRoleBinding，或完整同步两组规则到自定义 RBAC
- 确认 `pods` 包含 `watch`，否则创建 Pod 后无法等待 Running 状态
- 确认 `authorization.k8s.io/selfsubjectaccessreviews` 包含 `create`，否则启动自检本身无法执行
- VPC 模式需要 `k8s.cni.cncf.io/network-attachment-definitions` 和 `kubeovn.io/subnets`、`vpcs`、`ips` 权限
- VM 模式需要 `kubevirt.io/virtualmachines` 的 `create`、`get`、`list`、`watch`、`deletecollection`

## 镜像拉取失败

排查方向：

- 镜像仓库地址、tag、平台架构是否正确
- 平台 Pod 的私有仓库检查 `imagePullSecrets` / `imageCredentials`；题目、生成器和预热 Pod 检查运行命名空间的 `default` ServiceAccount 凭据，二者不自动继承
- 节点是否可以访问镜像仓库
- 是否达到镜像仓库限流

## 邮件发送失败

排查方向：

- 后台 SMTP 配置是否正确
- Redis 和 Asynq 邮件队列是否正常
- SMTP 服务是否允许集群出口访问
- 发件账号、密码、TLS/端口设置是否符合服务商要求

## 动态附件一直未就绪

- 在比赛的「生成器」页面确认实例 Ready；全局测试生成器不供比赛队伍使用
- 平台与生成器必须挂载同一个 PVC；使用已有卷时，检查 `persistence.existingClaim`，并确认 Chart 与应用镜像版本匹配
- 平台必须能够连接 generator Pod IP 的 TCP 8080
- 查看生成器日志与任务日志，检查 `/root/run.sh` 可执行、依赖已安装，且在一分钟内同步写出有效 ZIP
- `generator.zip` 必须以脚本/资源为根目录，而不是再套一层目录；重复测试可清理脚本自己的临时文件，避免复用上个团队产物

## Webhook 失败

排查方向：

- Webhook 目标是否在 `webhook.whitelist` 允许范围内
- 集群是否能访问目标地址
- 目标服务是否返回非 2xx 状态
- 查看后台 Webhook 历史记录和任务日志

## 安全建议

:::warning

- 首次登录后立即修改 `admin` 密码
- 生产环境设置固定、足够长的 `gin.jwt.secret`
- 数据库、Redis、SMTP、OAuth、FRP 使用强密码或随机密钥
- 使用 HTTPS，并正确配置 Ingress TLS
- 限制管理后台访问来源，按角色授予最小权限
- Kubernetes RBAC 不要授予超过平台运行所需的权限
- 为题目容器设置 CPU、内存、网络和命名空间隔离
- 定期备份 PostgreSQL 和 `/app/data`
- 不要公开 `config.yaml`、kubeconfig、私有镜像凭据和含敏感信息的日志
  :::
