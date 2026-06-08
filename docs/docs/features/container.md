---
sidebar_position: 5
---

# 动态靶机

容器题类型为 `pods`。选手启动题目后，平台会为队伍创建独立的 Kubernetes 资源，并在页面返回可访问端点。管理员可以在后台查看、预热、停止靶机，并下载流量文件。

## 基础要求

| 要求         | 说明                                                                                          |
|------------|---------------------------------------------------------------------------------------------|
| Kubernetes | 动态靶机由 Kubernetes 调度                                                                         |
| 命名空间       | `k8s.namespace` 必须存在，Helm 部署时通常为 Release namespace                                          |
| RBAC       | 应用需要创建/删除 Pod、Service、NetworkPolicy、Job、Multus、Kube-OVN 和 KubeVirt `VirtualMachine` 相关资源的权限 |
| 镜像拉取       | 集群节点必须能拉取题目镜像、流量捕获镜像和 FRP 镜像                                                                |
| 存储         | 题目附件、流量文件、动态附件依赖 `/app/data` 数据卷                                                            |
| 可选网络组件     | VPC 模式需要 Kube-OVN 和 Multus CNI，VM 模式额外需要 KubeVirt                                           |

## 管理员配置流程

1. 在题目管理中创建容器题。
2. 配置题目基础信息、分类、描述和 Flag。
3. 配置容器模板，包括镜像、命令、环境变量、资源、暴露端口和 Flag 注入方式。
4. 如题目需要多容器、隔离网络或虚拟机靶机，配置 Pod、VPC 或 KubeVirt VM 模式。
5. 保存后可使用题目测试启动，确认镜像、端口、Flag 和访问地址正常。
6. 将题目加入比赛，并按需在比赛管理中预热镜像或靶机。

Flag 可注入为环境变量或文件。环境变量和文件路径应与题目镜像内部程序读取方式一致。

## 选手使用流程

1. 进入比赛题目页面。
2. 对容器题点击启动靶机。
3. 等待平台创建资源并返回访问地址。
4. 通过页面展示的地址访问靶机。
5. 如支持延长时长，可在到期前延长。
6. 完成后提交 Flag，或手动停止靶机释放资源。

平台会限制每支队伍同时运行的靶机数量，限制值来自比赛配置中的靶机数量设置。达到上限时，需要停止旧靶机后再启动新的题目靶机。

## Pod 网络模式

Pod 模式适合单 Pod 内多容器共享网络命名空间的题目。所有容器在同一个 Kubernetes Pod 中运行，容器之间通过 `localhost:port`
通信。

适用场景：

- 简单 Web、Pwn、Misc 服务。
- 多个辅助容器需要共享本地网络。
- 不需要为每个容器配置静态 IP。

注意事项：

- 容器间不要依赖 Docker Compose 的服务名访问，应使用 `localhost`。
- 需要对选手开放的端口必须配置为暴露端口。
- 镜像架构应与集群节点一致，建议使用 `linux/amd64` 镜像。

## VPC 网络模式

VPC 模式适合需要多子网、静态 IP、网络隔离或模拟内网拓扑的题目。平台会基于 Kube-OVN 和 Multus 为队伍创建独立网络资源。

使用 VPC 模式前需要：

- 安装 Kube-OVN。
- 安装 Multus CNI。

如果 Kube-OVN/Multus 未安装或网络资源创建失败，VPC 靶机会启动失败或无法访问。

## KubeVirt VM 模式

KubeVirt VM 模式用于把某个 `docker-compose.yaml` service 按 KubeVirt `VirtualMachine` 创建，而不是创建普通 Kubernetes
Pod。平台使用 service 的 `image` 作为 KubeVirt `containerDisk` 镜像，并通过 Multus/Kube-OVN 将 VM 接入题目 VPC 网络。

使用 VM 模式前需要：

- 集群已安装 KubeVirt，且 KubeVirt CRD 和控制器正常工作。
- 题目必须使用 VPC 网络模式，并已安装 Kube-OVN 和 Multus CNI。
- 应用 ServiceAccount 具备 `kubevirt.io` 资源 `virtualmachines` 的 `create`、`get`、`deletecollection` 权限。
- VM 镜像需要是 KubeVirt 可用的 `containerDisk` 镜像，镜像内系统需要支持 cloud-init 和 virtio 网卡。
- 集群节点支持虚拟化能力；如使用嵌套虚拟化，需要提前在宿主机和集群节点上启用。

在 compose 中为 service 设置 `x-kubevirt: true` 即可启用 VM 模式。代码会额外校验以下字段：

| 字段                        | 要求                                   |
|---------------------------|--------------------------------------|
| `networks`                | 必须选择自定义 VPC 网络，不能只使用默认网络             |
| `networks.*.ipv4_address` | 必填，且必须在对应网络 `subnet` 范围内             |
| `networks.*.mac_address`  | VM 模式必填，格式必须是合法 MAC 地址               |
| `mem_limit`               | VM 模式必填，会映射到 VM 内存限制                 |
| `cpus`                    | 可选，会换算为 CPU milli limit              |
| `x-boot.bootloader`       | 可选，只支持 `bios` 或 `efi`                |
| `x-boot.secure_boot`      | 可选，仅 `bootloader: efi` 时有意义          |
| `x-cloudinit`             | 可选，会渲染为 `CloudInitNoCloud` user data |

VM 模式可以在 compose 中保留 `ports` 作为题目作者的服务端口记录，但平台不会读取这些端口创建 Service、NodePort、FRP
转发或访问端点。需要选手访问 VM 服务时，应通过题目网络拓扑、VM 内部代理或其他题目自行设计的方式提供入口；平台只负责创建 VM
和接入 VPC 网络。

最小示例：

```yaml
services:
  vm1:
    image: registry.example.com/challenges/vm-web:latest
    cpus: 0.5
    mem_limit: 512m
    ports:
      - target: 80
        published: web
        protocol: tcp
    x-kubevirt: true
    x-boot:
      bootloader: efi
      secure_boot: false
    x-cloudinit:
      users:
        - name: ctf
          groups:
            - sudo
          sudo:
            - ALL=(ALL) NOPASSWD:ALL
          shell: /bin/bash
          homedir: /home/ctf
          lock_passwd: false
          plain_text_passwd: changeme
      write_files:
        - path: /etc/motd
          content: |
            Welcome to CBCTF VM.
    networks:
      network1:
        ipv4_address: 192.168.0.2
        mac_address: "00:00:00:00:01:01"

networks:
  network1:
    ipam:
      config:
        - subnet: 192.168.0.0/24
          gateway: 192.168.0.1
```

完整 VPC/VM 示例可参考仓库中的 `example/pods/vpc/docker-compose.yaml`。

VM 模式注意事项：

- `command`、`working_dir`、普通环境变量和 `x-volumes` 不会被注入到 KubeVirt VM。需要初始化系统用户、写入文件或写入 Flag
  时，应使用 `x-cloudinit.write_files`。
- `x-cloudinit.write_files[*].content` 中出现的 `static{}`、`leet{}` 或 `uuid{}` 会被识别为 cloud-init 文件
  Flag，并在启动队伍靶机时替换为队伍实际 Flag。
- `ports` 不被使用，虚拟机不可直接暴露端口。
- VM 必须接入至少一个 VPC 网络。平台会把第一个网络设置为默认 Multus 网络，并为每张网卡生成 cloud-init network data。
- 当某个 service 设置为 `x-kubevirt: true` 时，该 service 在 VPC 模式下会以单独的 VM 运行，不会和普通容器共享同一个 Pod
  网络命名空间。

## 端口暴露

平台会根据普通 Pod service 的暴露端口生成访问端点。KubeVirt VM service 的 `ports` 不参与端点生成。

| 模式         | 说明                                                         |
|------------|------------------------------------------------------------|
| 默认暴露       | 普通 Pod service 创建 Kubernetes Service，并返回可访问端点              |
| FRP 暴露     | 当 `k8s.frp.on` 为 `true` 时，通过 FRPS 服务端和分配端口暴露普通 Pod service |
| VM service | 不创建 Service、NodePort、FRP 转发或访问端点                           |

使用 FRP 时需要配置 `k8s.frp.frps`，包括 FRPS 地址、端口、token 和可用端口池。端口池不足会导致新靶机无法分配访问端口。

## 生命周期和回收

选手端支持启动、延长和停止靶机。管理员端支持查看全局或比赛内靶机列表，并批量停止靶机。平台也有 Cron 任务用于后台维护，具体执行状态可在管理后台任务和
CronJobs 页面查看。

如果发现靶机未及时释放，可在管理后台靶机页面停止，或排查后台任务队列和 Kubernetes 资源状态。

## 流量捕获

容器题可使用 `k8s.capture` 配置的镜像进行流量捕获。管理员可在靶机或队伍详情中查看并下载 pcap 文件。

排查流量文件为空或下载失败时，检查：

- `k8s.capture` 镜像是否可拉取。
- `/app/data` 数据卷是否可写。
- 靶机是否实际产生流量。
- 管理后台文件和流量接口权限是否正常。

## 镜像预热和生成器

管理后台提供镜像预热、生成器管理和比赛内预热入口。比赛开始前建议预热大型题目镜像，降低选手首次启动等待时间。

动态附件题会使用生成器镜像和上传的生成器文件创建 Kubernetes Job。生成失败时，优先查看生成器任务、Pod 日志、共享 PVC 和镜像拉取状态。

## 常见问题

| 问题          | 排查方向                                                              |
|-------------|-------------------------------------------------------------------|
| 靶机启动失败      | 查看应用日志、任务日志和对应 Kubernetes Pod/Event                               |
| 提示权限不足      | 检查 ServiceAccount、ClusterRole、ClusterRoleBinding 是否存在             |
| 命名空间不存在     | 创建 `k8s.namespace` 对应命名空间，Helm 部署通常自动使用 Release namespace         |
| 镜像拉取失败      | 检查镜像名、tag、架构、网络、私有仓库凭据和 `imagePullSecrets`                        |
| Pod Pending | 检查节点资源、调度约束、PVC、镜像拉取和资源限制                                         |
| VM 创建失败     | 检查 KubeVirt 是否安装、`virtualmachines` RBAC、节点虚拟化能力和 containerDisk 镜像 |
| 页面没有访问地址    | 普通 Pod 题检查暴露端口、Service、FRP 配置和应用日志；VM 题不会由 `ports` 自动生成访问地址       |
| 访问地址不可达     | 检查 Ingress/Service/NodePort/FRP、网络策略、防火墙和 `cbctf.host`            |
| VPC/VM 模式失败 | 检查 Kube-OVN、Multus、KubeVirt、静态 IP 和 MAC 地址                        |
| 动态附件失败      | 检查共享 PVC、生成器镜像、生成器文件和 Job 日志                                      |
| 无法停止靶机      | 检查 Kubernetes 删除权限、任务队列和残留资源标签                                    |

## 运维建议

- 为题目容器设置合理 CPU 和内存，避免单个队伍耗尽集群资源。
- 为 VM 题设置更保守的 CPU 和内存限制，并在比赛前压测节点可承载的 VM 数量。
- 将比赛靶机运行在独立命名空间，并限制访问生产内网。
- 私有镜像使用专用拉取凭据，不要把凭据写入公开仓库。
- 比赛前预热镜像和关键靶机，确认普通 Pod 端口、VPC 网络和 KubeVirt VM 启动正常；VM 访问路径需按题目设计单独验证。
- 比赛期间监控 Pod Pending、ImagePullBackOff、OOMKilled 和任务队列积压。
- 比赛结束后清理靶机、检查 PVC 容量并备份必要文件。
