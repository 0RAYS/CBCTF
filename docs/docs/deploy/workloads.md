---
title: 工作负载调度
description: 说明全节点镜像预热、失败节点排除、共享 informer、靶机异步就绪、附件 worker 池与资源树回收。
---

# 工作负载调度

## 镜像预热与自主调度

题目创建、更新、加入比赛及取消隐藏时，会提交镜像预热任务。任务覆盖题目镜像和实际使用的 capture、FRPC、nginx、generator worker 安装镜像。

- 预热目标为 Ready、未 cordon、无 `NoSchedule` / `NoExecute` 污点的节点，与当前工作负载不配置 toleration 的行为一致。
- 只有预热 Job 定向到节点，用于覆盖所有可调度节点；业务 Pod、generator Pod 和 VM 不设置 `nodeName`、`nodeSelector` 或正向节点偏好。
- 节点镜像清单明确包含所需镜像时，跳过对应预热 Job。镜像清单可能被 kubelet 截断，因此清单中缺失镜像不等于拉取失败。
- 实际观察到 `ErrImagePull`、`ImagePullBackOff` 或 `InvalidImageName` 后，记录该镜像在该节点的失败状态。
- 创建工作负载时，只用 `metadata.name NotIn` 排除已确认失败的节点。没有失败记录时，完全不添加 affinity，交由 Kubernetes 调度。
- 自动重试、手动预热或周期复查确认镜像可用后，清除失败记录。所有节点失败时，工作负载保持 Pending，最终由启动超时流程回收。

`WarmChallengeImages` Cron 默认每 15 分钟复查一次，也覆盖新加入的节点。预热任务在 `tasks:prepull` 队列中执行，可在任务日志中查看结果。手动预热也走同一套结果记录流程；指定 `Always` 时会重新访问镜像仓库。

预热成功不构成节点永久保留镜像的承诺：kubelet 镜像 GC、节点重建或镜像标签变化仍可能产生冷启动。比赛镜像宜使用固定版本或 digest。

## 靶机创建与就绪

```text
waiting → pending：取得跨实例锁，记录启动截止时间
        → 创建网络策略、网络、Pod/VM、Service 和 FRPC
        → 持久化对象 UID、端点信息，释放启动 worker
pending → running：共享缓存确认所有 Pod/VM 就绪，写入对外端点和开始时间
pending → terminating：启动失败或超过 4 分钟，提交停止任务
```

协调器每秒读取待启动记录，基于共享 informer 判断就绪，不为每个 Pod 建立独立 watch。它使用非阻塞 PostgreSQL advisory lock 跳过正在创建的实例，多个平台副本之间仍由锁和条件状态更新协调。进程重启后，数据库中的 Pending 记录会继续被检查。

Pod、VM、ConfigMap 根对象以及本命名空间对应的 VPC 各使用共享 informer。未用于状态观测的对象不额外建立 watch。缓存缺失不能独立证明删除完成，清理收尾保留直接 API 确认；对象 UID 变化不会被视为原实例就绪。

启动和停止有各自的 Asynq worker。`asynq.queues.victim` 默认是 `8`，分别控制两类 worker 的并发；等待就绪不会占用启动槽位。

### 对象与容器

- 每个 Pod 的 `x-volumes` 合并到一个 ConfigMap。每个容器、每个文件使用独立 key，通过 `subPath` 保持目标路径，同名文件不会互相覆盖。合并后的内容仍受 Kubernetes 单个 ConfigMap 大小限制。
- ConfigMap、Service、NetworkPolicy 和 FRPC Pod 使用由 victim ID 和对象 key 派生的确定性名称。
- 配置文件中 `k8s.capture_enabled: false` 会关闭靶机和 FRPC 的抓包容器及相关 NFS 挂载。
- FRPC 与 nginx 共用一个配置 ConfigMap；纯 UDP 暴露不启动 nginx。
- 设置了 limits 的 CPU/内存资源使用相同 requests；未指定的 requests 默认是 `100m` / `64Mi`。capture、FRPC、nginx 的 requests 为 `10m` / `32Mi`，limits 为 `500m` / `256Mi`。
- `k8s.priority_class_name` 可指定已存在的 PriorityClass，不对节点增加正向亲和性。

NetworkPolicy 先于工作负载创建，避免引入未隔离窗口。VPC 模式先创建 VPC，再并行创建依赖它的 Subnet 与 NAD。

## 网络暴露

启用 `k8s.frp.on` 后，工作负载 Service 使用 ClusterIP，FRPC 直接连接集群内端点，不再分配每题 NodePort。Service 创建后即可启动 FRPC，无需等待靶机 Pod Ready；nginx 内部监听端口按端点分配，支持不同 Service 暴露相同目标端口。

未启用 FRP 时使用 NodePort，协调器在 Pod Ready 后写入节点地址。FRP 需要真实可用的 FRPS 地址及端口池，因此默认保持关闭。

## 附件 worker 池与缓存

```text
生成团队 flag → Redis/Asynq 去重任务 → 执行时领取空闲 generator
            → HTTP 调用常驻 worker → 本地生成完整 ZIP
            → 流式传回平台 → 校验 ZIP、fsync、原子 rename → 缓存可下载
```

### 运行方式

独立镜像 `ghcr.io/0rays/cbctf-worker` 基于 `scratch`，只包含静态二进制 `/app/worker`。generator Pod 的 init container 使用此镜像，将二进制复制到 EmptyDir 的 `/worker/worker`，题目容器运行该 worker。worker 在容器内部解压 `generator.zip` 并调用：

```text
/root/run.sh <team_id> <base64_encoded_flags>
```

flag 参数仍是「各 flag 分别 Base64 编码、逗号拼接、整体再 Base64 编码」。脚本同步完成，并写入 `/root/mnt/attachments/{id}.zip`。这个输出目录是 Pod 本地 EmptyDir；`generator.zip` 等共享输入继续从共享 PVC 读取。脚本不应把尚未完成的生成过程交给后台进程。

worker 使用每实例 token 校验请求，串行执行生成任务；题目容器不持有平台 Redis 或 Kubernetes 凭据。平台需要能访问 generator Pod IP 的 TCP 8080 端口。附件生成和生成器初始化都不再调用 Kubernetes Exec。

Redis/Asynq 由平台负责队列分发，worker 接收执行请求，无需题目镜像自带 Redis 客户端。单次脚本执行限制为 1 分钟。ZIP 通过完成响应流式返回，未完整传输或格式无效的文件不会发布到缓存。

### 池容量与缓存身份

- `k8s.generator_pool_size` 默认 `2`，为每个比赛、每个动态题维护启动时的目标池容量。题目加入比赛和发布时预建；排队任务遇到空池时也会补齐。设为 `0` 可仅使用管理员手动启动的实例。
- 可以在管理后台启动更多 generator。附件任务领取实例发生在执行阶段，排队中的任务不占用 generator 租约。
- `asynq.queues.attachment` 默认 `32`，是平台分发并发上限，独立于每题 Pod 数。池繁忙时延后任务，释放 worker 槽位。
- 同一题、同一团队、同一组有序 flag 和源文件版本具有相同任务身份。重复请求复用队列任务或缓存结果。
- 平台缓存路径为 `{path}/challenges/{challenge_id}/cache/{team_id}-{hash}.zip`。hash 包含题目版本、镜像、源文件版本和 flag；重置 flag 的新任务不会与旧任务写入同一个文件。
- 下载和题目状态查询按照当前团队 flag 解析缓存路径，过期版本不会被当成当前附件。

### NFS 参数

生成完成信号来自 worker 响应，生成端的产物位于 EmptyDir，因此不依赖 NFS `Stat` / `ReadDir` 轮询。

多个平台副本共享最终缓存时，仍应按存储系统配置 NFS 一致性。挂载参数属于 PV 或 StorageClass，而不是 PVC。可在 NFS provisioner 的 StorageClass 中设置：

```yaml
mountOptions:
  - actimeo=1
  - lookupcache=positive
```

`noac` 会影响整个挂载的属性缓存和写入行为，需要结合附件大小与存储吞吐选择。Chart 使用已有 StorageClass/PVC，不会修改其挂载策略。

## 有序资源树回收

每个 victim 有两个 namespaced ConfigMap 根对象：

1. `victim-{id}-workloads`：持有 Pod、VM、Service、文件配置和 FRPC 配置。
2. `victim-{id}-network`：持有 NetworkPolicy 与 NAD。

停止时先 foreground 删除工作负载根，确认全部 Pod/VM 消失，再删除网络根。集群级 Subnet 由集群级 VPC 持有；工作负载消失后清理 CNI 的 IP 记录，再 foreground 删除 VPC。namespaced ConfigMap 不作为集群级对象的 owner。

根对象包含身份标签与网络清理信息。孤儿清理从共享 ConfigMap 缓存发现资源树，即使创建过程尚未产生 Pod，也能够发现残留。删除使用观测到的 UID 作为前置条件，清理失败保留记录供重试。

## Helm 配置示例

```yaml
image:
  repository: ghcr.io/0rays/cbctf
  tag: your-build-tag

cbctf:
  k8s:
    captureEnabled: false
    priorityClassName: ""
    workerImage: ghcr.io/0rays/cbctf-worker:your-worker-tag
    generatorPoolSize: 4
  asynq:
    queues:
      victim: 8
      generator: 3
      attachment: 32
```

Chart 将 `cbctf.k8s.workerImage` 写入 `k8s.worker_image`，worker 的仓库和版本独立于主程序镜像。直接部署时在配置文件中设置 `k8s.worker_image`。上述抓包开关、PriorityClass、worker 镜像和池容量是部署配置，修改后重启平台生效。

## 独立构建 worker 镜像

从仓库根目录构建：

```bash
docker build -f worker/Dockerfile -t ghcr.io/0rays/cbctf-worker:your-worker-tag .
```

该构建只编译 `worker/main.go` 和 `internal/worker`，不构建前端或主程序，也不安装 CGO/libpcap。主程序 Dockerfile 不再打包 worker 二进制。

`.github/workflows/worker.yaml` 独立发布 worker 镜像，监听 worker 源码和构建配置变化，生成时间戳标签与 `latest` 标签。主程序和 worker 可以分别构建、发布和指定版本。

秒级启动与附件生成时长需要在实际集群中测量；创建阶段解耦提升吞吐，但不会缩短题目进程自身的初始化或脚本执行时间。
