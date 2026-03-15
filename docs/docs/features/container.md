---
sidebar_position: 5
---

# 容器靶机

容器题（`pods` 类型）为每支队伍部署独立的靶机容器，出题人通过编写 `docker-compose.yaml` 配置容器环境。

## 网络模式对比

| 特性 | Pod 网络模式 | VPC 网络模式 |
|------|:----------:|:-----------:|
| docker-compose 条件 | 无 `networks` 字段 | 有 `networks` 字段 |
| 容器间通信 | `localhost:port` | 静态 IP 地址 |
| 子网隔离 | 不支持 | 支持多子网 |
| 外部网络访问 | 默认允许 | 按 `external` 字段控制 |
| KubeOVN 依赖 | 不需要 | **必须** |
| Multus CNI 依赖 | 不需要 | **必须** |

## Pod 网络模式

当 `docker-compose.yaml` 中**不含 `networks` 字段**时，平台使用 Pod 网络模式。

所有服务被部署在同一个 K8s Pod 中，共享网络命名空间：

```yaml
services:
  web:
    image: myapp:latest
    ports:
      - "5000:5000"
    environment:
      - FLAG_0=uuid{}
      - DB_HOST=localhost    # 注意：使用 localhost，而非服务名
      - DB_PORT=3306

  db:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: example
      MYSQL_DATABASE: testdb
```

:::warning 重要区别
在 Kubernetes Pod 中，**容器间通信使用 `localhost:port`**，而非 docker-compose 中的 `服务名:port`。出题人必须确保应用代码使用 `localhost` 访问同 Pod 内的其他服务。
:::

通过 `ports` 字段暴露的端口对选手可见，平台通过 FRP（如已启用）或 NodePort 转发给选手。

## VPC 网络模式

当 `docker-compose.yaml` 中**含有 `networks` 字段**时，平台使用 VPC 网络模式。

平台通过 KubeOVN 为每支队伍创建独立的 VPC 和子网，每个服务分配静态 IP。

### 多子网拓扑示例

```yaml
services:
  web1:
    ports:
      - "5000:5000"
    environment:
      FLAG_1: uuid{}
    image: myapp:latest
    networks:
      network1:
        ipv4_address: 192.168.0.2
      network2:
        ipv4_address: 192.168.1.2

  web2:
    environment:
      FLAG_2: uuid{}
    image: myapp:latest
    networks:
      network2:
        ipv4_address: 192.168.1.3
      network3:
        ipv4_address: 192.168.2.3

  web3:
    environment:
      FLAG_3: uuid{}
    image: myapp:latest
    networks:
      network3:
        ipv4_address: 192.168.2.4

networks:
  network1:
    ipam:
      config:
        - subnet: 192.168.0.0/24
          gateway: 192.168.0.1
    external: true     # 此网络可访问外部互联网
  network2:
    ipam:
      config:
        - subnet: 192.168.1.0/24
          gateway: 192.168.1.1
  network3:
    ipam:
      config:
        - subnet: 192.168.2.0/24
          gateway: 192.168.2.1
```

**网络连通性**：

- `web1` ↔ `web2`：可以（均在 network2，通过 `192.168.1.x` 通信）
- `web2` ↔ `web3`：可以（均在 network3，通过 `192.168.2.x` 通信）
- `web1` ↔ `web3`：**不可以**（无共同子网）
- `network1` 中的容器：可以访问外部互联网（`external: true`）
- `network2`、`network3`：仅内网通信，不可访问外部互联网

## 容器生命周期

1. 选手触发 **init**：平台生成 flag → 启动 Pod → 注入 flag
2. 选手**延长时长**：重置计时器，不重启容器
3. 选手触发 **stop** 或时间到期：停止并删除 Pod

## 并发限制

`victims_count` 字段（在比赛配置中）限制每支队伍同时运行的靶机数量。选手需要停止旧实例才能启动新实例。

## 流量捕获

平台使用 `k8s.tcpdump` 配置的镜像作为 sidecar 容器，自动捕获靶机流量并保存为 pcap 文件。管理员可从后台下载。需要在 Helm values 或 config.yaml 中配置对应镜像。

## FRP 端口转发

当 `k8s.frp.on: true` 时，容器的 `ports` 字段中定义的端口通过 frps 服务器对外暴露，选手通过 frps 服务器地址和分配的端口号访问靶机。

## 网络策略（NetworkPolicy）

可在题目配置的 `NetworkPolicies` 字段中配置 Kubernetes NetworkPolicy，限制靶机的入/出站流量，例如：

- 禁止出站访问互联网（防止反弹 shell 外联）
- 仅允许指定 CIDR 的选手 IP 入站

## docker-compose.yaml 编写要求

- 建议在服务中指定 `platform: linux/amd64`，避免镜像架构问题
- 静态 IP 地址需在各服务的 `networks` 下通过 `ipv4_address` 指定
- 所有服务必须归属至少一个 network（VPC 模式下）

## 官方示例

- [Pod 网络示例](https://github.com/0RAYS/CBCTF/blob/main/example/pods/pod/docker-compose.yml)
- [VPC 网络示例](https://github.com/0RAYS/CBCTF/blob/main/example/pods/vpc/docker-compose.yml)
