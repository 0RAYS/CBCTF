# 动态靶机

容器题类型为 `pods`。选手启动题目后，平台会为队伍创建独立的 Kubernetes 资源，并在页面返回可访问端点。管理员可以在后台查看、预热、停止靶机，并下载流量文件。

## 基础要求

| 要求 | 说明 |
|---|---|
| Kubernetes | 动态靶机由 Kubernetes 调度 |
| 命名空间 | `k8s.namespace` 必须存在，Helm 部署时通常为 Release namespace |
| RBAC | 应用需要创建/删除 Pod、Service、NetworkPolicy、Job、Multus、Kube-OVN 和 KubeVirt `VirtualMachine` 相关资源的权限 |
| 镜像拉取 | 集群节点必须能拉取题目镜像、流量捕获镜像和 FRP 镜像 |
| 存储 | 题目附件、流量文件、动态附件依赖 `/app/data` 数据卷 |
| 可选网络组件 | VPC 模式需要 Kube-OVN 和 Multus CNI，VM 模式额外需要 KubeVirt |

## Pod 网络模式

Pod 模式适合单 Pod 内多容器共享网络命名空间的题目。所有容器在同一个 Kubernetes Pod 中运行，容器之间通过 `localhost:port` 通信。

适用场景：

- 简单 Web、Pwn、Misc 服务
- 多个辅助容器需要共享本地网络
- 不需要为每个容器配置静态 IP

:::tip
容器间不要依赖 Docker Compose 的服务名访问，应使用 `localhost`。
:::

## VPC 网络模式

VPC 模式适合需要多子网、静态 IP、网络隔离或模拟内网拓扑的题目。平台会基于 Kube-OVN 和 Multus 为队伍创建独立网络资源。

使用 VPC 模式前需要：

- 安装 Kube-OVN
- 安装 Multus CNI

## KubeVirt VM 模式

KubeVirt VM 模式用于把某个 `docker-compose.yaml` service 按 KubeVirt `VirtualMachine` 创建，而不是创建普通 Kubernetes Pod。

在 compose 中为 service 设置 `x-kubevirt: true` 即可启用 VM 模式。代码会额外校验以下字段：

| 字段 | 要求 |
|---|---|
| `networks` | 必须选择自定义 VPC 网络，不能只使用默认网络 |
| `networks.*.ipv4_address` | 必填，且必须在对应网络 `subnet` 范围内 |
| `networks.*.mac_address` | VM 模式必填，格式必须是合法 MAC 地址 |
| `mem_limit` | VM 模式必填，会映射到 VM 内存限制 |
| `cpus` | 可选，会换算为 CPU milli limit |
| `x-boot.bootloader` | 可选，只支持 `bios` 或 `efi` |
| `x-boot.secure_boot` | 可选，仅 `bootloader: efi` 时有意义 |
| `x-cloudinit` | 可选，会渲染为 `CloudInitNoCloud` user data |

最小示例：

```yaml
services:
  vm1:
    image: registry.example.com/challenges/vm-web:latest
    cpus: 0.5
    mem_limit: 512m
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

:::warning VM 模式注意事项
- `command`、`working_dir`、普通环境变量和 `x-volumes` 不会被注入到 KubeVirt VM，请使用 `x-cloudinit.write_files` 写入文件和 Flag
- `x-cloudinit.write_files[*].content` 中出现的 `static{}`、`leet{}` 或 `uuid{}` 会被识别为 cloud-init 文件 Flag
- `ports` 不被使用，虚拟机不可直接暴露端口，访问路径需按题目设计单独设计
- VM 必须接入至少一个 VPC 网络
:::

## 端口暴露

| 模式 | 说明 |
|---|---|
| 默认暴露 | 普通 Pod service 创建 Kubernetes Service，并返回可访问端点 |
| FRP 暴露 | 当 `k8s.frp.on` 为 `true` 时，通过 FRPS 服务端和分配端口暴露普通 Pod service |
| VM service | 不创建 Service、NodePort、FRP 转发或访问端点 |

## 流量捕获

容器题可使用 `k8s.capture` 配置的镜像进行流量捕获。管理员可在靶机或队伍详情中查看并下载 pcap 文件。

## 常见问题

| 问题 | 排查方向 |
|---|---|
| 靶机启动失败 | 查看应用日志、任务日志和对应 Kubernetes Pod/Event |
| 提示权限不足 | 检查 ServiceAccount、ClusterRole、ClusterRoleBinding 是否存在 |
| 镜像拉取失败 | 检查镜像名、tag、架构、网络、私有仓库凭据和 `imagePullSecrets` |
| Pod Pending | 检查节点资源、调度约束、PVC、镜像拉取和资源限制 |
| VM 创建失败 | 检查 KubeVirt 是否安装、`virtualmachines` RBAC、节点虚拟化能力和 containerDisk 镜像 |
| 页面没有访问地址 | 普通 Pod 题检查暴露端口、Service、FRP 配置；VM 题不会由 `ports` 自动生成访问地址 |
| VPC/VM 模式失败 | 检查 Kube-OVN、Multus、KubeVirt、静态 IP 和 MAC 地址 |

## 运维建议

- 为题目容器设置合理 CPU 和内存，避免单个队伍耗尽集群资源
- 为 VM 题设置更保守的资源限制，并在比赛前压测节点可承载的 VM 数量
- 比赛前预热镜像和关键靶机，确认普通 Pod 端口、VPC 网络和 KubeVirt VM 启动正常
- 比赛期间监控 Pod Pending、ImagePullBackOff、OOMKilled 和任务队列积压
- 比赛结束后清理靶机、检查 PVC 容量并备份必要文件
