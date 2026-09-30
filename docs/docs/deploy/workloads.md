---
title: 工作负载调度
description: 配置镜像预热、靶机启动、附件生成器容量和资源清理，排查等待、超时与生成失败。
---

# 工作负载调度

## 镜像预热

镜像预热将题目及辅助服务所需镜像提前拉取到节点，可减少开赛时的等待。题目创建、更新、加入比赛及发布时会自动提交预热任务；也可在「镜像」页面手动预热。

1. 选择需要预热的节点和镜像，或输入镜像名称。
2. 选择拉取策略：`IfNotPresent` 拉取缺失镜像，`Always` 重新检查镜像仓库，`Never` 不创建预热任务。
3. 提交后查看逐节点结果。部分提交失败时，核对失败和未执行项后重试。
4. 在「任务队列」查看预热执行结果，并在镜像页刷新节点清单。

预热面向 Ready、未暂停调度且没有 `NoSchedule` / `NoExecute` 污点的节点。`WarmChallengeImages` 定时任务默认每 15 分钟复查未结束比赛的题目，包括尚未开始的比赛和隐藏题目，也会覆盖新加入的节点。

出现 `ErrImagePull`、`ImagePullBackOff` 或 `InvalidImageName` 时，检查镜像名称、仓库连通性和拉取凭据。平台会避开该镜像已确认拉取失败的节点；修复后手动预热或等待周期复查，成功后节点可重新参与调度。所有节点均不可用时，靶机会等待至启动超时。

预热任务最长等待约 10 分钟，结束后会清理临时 Job 和 Pod，历史结果可在任务队列查看。节点磁盘清理、重建或镜像版本变化后可能需要重新拉取；比赛镜像建议固定版本或 digest。

## 靶机创建与就绪

| 状态 | 含义与操作 |
| --- | --- |
| `waiting` | 启动任务正在排队，可查看任务队列是否积压 |
| `pending` | 正在准备网络、拉取镜像或启动服务，可打开实例状态和日志查看进度 |
| `running` | 实例已就绪，使用页面提供的地址访问 |
| `terminating` | 正在停止并释放资源，等待清理完成 |

启动失败或等待超过约 4 分钟后，平台会停止该实例。长时间处于 `pending` 时，依次检查镜像拉取、节点剩余资源、存储挂载以及题目进程日志。平台重启后会继续检查未完成的启动。

`asynq.queues.victim` 默认是 `8`，分别控制启动和停止任务的并发数。大量选手同时启动时，可结合节点容量调整；提高并发不会加快题目进程自身的初始化。

### 资源与抓包配置

- 在题目模板中填写 CPU、内存限制，调度时也会预留相同资源；未指定时默认预留 `100m` CPU 和 `64Mi` 内存。
- 抓包、FRPC、nginx 辅助容器各预留 `10m` CPU 和 `32Mi` 内存，上限为 `500m` / `256Mi`，规划容量时需一并计入。
- `k8s.capture_enabled: false` 关闭新实例的抓包。已有实例保持原设置。
- `k8s.priority_class_name` 可指定集群中已有的 PriorityClass，适用于普通 Pod，不应用到 VM。
- `x-volumes` 适合小型配置和 Flag 文件；大型文件请放入镜像或题目附件。

## 网络暴露

未启用 FRP 时，普通容器通过 NodePort 暴露服务。启用 `k8s.frp.on` 后，访问地址使用 FRPS 的地址和分配端口；需事先部署 FRPS，并配置服务器地址、认证 token 和允许分配的端口范围。

靶机就绪后，以页面显示的访问地址为准。VM 的访问入口需在题目中单独设计，详见[动态靶机](../guide/features/container)。

## 附件生成器容量

动态附件题加入比赛或发布时，会提前启动生成器。选手初始化题目后，由空闲生成器制作该队附件；繁忙时任务排队等待。

| 设置 | 用途 |
| --- | --- |
| `k8s.generator_pool_size` | 每场比赛、每道动态题预建的生成器数量，默认 `2`；设为 `0` 时由管理员手动启动 |
| `asynq.queues.generator` | 生成器启停任务并发数，默认 `3` |
| `asynq.queues.attachment` | 平台同时处理附件生成任务的上限，默认 `32` |
| `k8s.worker_image` | 附件执行服务镜像，默认使用 `ghcr.io/0rays/cbctf-worker` |

比赛前在「比赛 → 生成器」确认实例已就绪，并根据脚本耗时和参赛队伍数增加实例。全局生成器用于题库测试，不能代替比赛专用生成器。

生成失败时，检查：

1. 平台和生成器使用同一共享 PVC。Helm 通过 `persistence.existingClaim` 指定已有卷；默认卷名为 `{namespace}-shared-volume`。
2. 节点能拉取题目镜像和 worker 镜像，平台可访问生成器 Pod 的 TCP 8080 端口。
3. 生成脚本在 1 分钟内完成，并输出有效 ZIP。入口和参数见[动态附件生成](../guide/features/attachment)。

同一队伍重复请求可复用已有附件。重置 Flag 或更新题目源文件后会生成相应的新附件，选手应重新下载。

### NFS 参数

多个平台副本需要共享题目源文件和最终附件。使用 NFS 时，在 PV 或 StorageClass 中配置挂载参数，而不是在 PVC 中配置。例如：

```yaml
mountOptions:
  - actimeo=1
  - lookupcache=positive
```

`noac` 会影响整个挂载的属性缓存和写入行为，需要结合附件大小与存储吞吐选择。Chart 使用已有 StorageClass/PVC，不会修改其挂载策略。

## 停止与清理

停止实例后，平台先关闭 Pod/VM，再释放网络和地址。列表显示 `terminating` 时，清理尚未完成。

长时间未结束时，查看任务日志、Pod/VM 的删除状态以及 CNI 组件状态。清理失败的资源会保留待重试；定时任务也会检查超时或失控实例。排障方法见[故障排查](./troubleshooting)。

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

Chart 将 `cbctf.k8s.workerImage` 写入 `k8s.worker_image`，worker 的仓库和版本独立于主程序镜像。配置文件提供初始值，随后以数据库设置为准。抓包开关、Pod PriorityClass、worker 镜像和池容量可在「系统管理 → Kubernetes 配置」中修改，对后续新建实例及补池生效；不会自动重建或停止已有实例。其他平台副本需重新加载配置。

开赛前建议按预计并发启动一批测试实例，检查镜像拉取、附件耗时和资源占用，再确定池容量与队列并发。
