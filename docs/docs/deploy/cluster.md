---
sidebar_position: 1
---

# K8s 集群搭建

本文介绍如何搭建支持 CBCTF 动态附件题和容器题的 Kubernetes 集群。由于项目当前仅支持 Helm 部署, 若仅使用静态题, 也建议至少准备基础的 Kubernetes、PostgreSQL 和 Redis 运行环境。

## 硬件要求

准备多台运行 Ubuntu 22.04+ 的 Linux 服务器, 分为: 

- **Master 节点** × 1: 管理 K8s 控制平面
- **Worker 节点** × N: 运行题目容器, 建议至少 2 台

### VMware ESXi 说明

若使用 ESXi 虚拟机, 需在 vSwitch 上开启混杂模式（`Promiscuous Mode`）和 MAC 地址更改（`MAC Address Changes`）, 否则 VPC 网络无法正常工作。

![vswitch.png](img/vswitch.png)

> 不支持 macvlan 的云服务商通常无法使用 VPC 网络模式, 只能运行 Pod 网络模式的容器题。

### 外部网络网卡标签

CBCTF 支持为不同节点配置不同的外部网卡。Helm values 中的 `cbctf.k8s.external_networks.interfaces` 定义每个外部网络使用的宿主机网卡名、CIDR 和网关。后端启动时会为每一项创建对应的 Kube-OVN Subnet 与 Multus NAD。

需要为可承载该外部网络的节点添加固定标签, 标签值必须等于配置中的 `interface`: 

```bash
kubectl label node <node-with-eth0> node.cbctf.io/external-network=eth0
kubectl label node <node-with-ens192> node.cbctf.io/external-network=ens192
```

`VpcNatGateway` 会根据所选外部网络自动添加节点亲和性, 仅调度到匹配 `node.cbctf.io/external-network=<interface>` 的节点。

## 安装 NFS 客户端

动态附件与题目文件共享依赖 RWX 存储。所有节点需要安装 NFS 客户端: 

```bash
sudo apt update
sudo apt install -y nfs-common
```

## 安装 K3S

CBCTF 推荐使用 [K3S](https://docs.k3s.io/)。

### Master 节点

```bash
curl -sfL https://rancher-mirror.rancher.cn/k3s/k3s-install.sh | \
  INSTALL_K3S_MIRROR=cn sh - --flannel-backend=none --disable-network-policy
```

安装完成后: 

- kubeconfig 位于 `/etc/rancher/k3s/k3s.yaml`
- 节点 Token 位于 `/var/lib/rancher/k3s/server/node-token`

### Worker 节点

```bash
curl -sfL https://rancher-mirror.rancher.cn/k3s/k3s-install.sh | \
  INSTALL_K3S_MIRROR=cn \
  K3S_URL=https://myserver:6443 \
  K3S_TOKEN=mynodetoken \
  sh
```

## 安装 Multus CNI

[Multus CNI](https://github.com/k8snetworkplumbingwg/multus-cni) 为 VPC 网络模式提供多网卡支持: 

```bash
kubectl apply -f https://raw.githubusercontent.com/k8snetworkplumbingwg/multus-cni/master/deployments/multus-daemonset.yml
```

若 K3S + Flannel 环境出现 `cannot find valid master CNI config`, 需按 Multus 官方文档调整 `--multus-kubeconfig-file-host` 与 CNI 目录。

## 安装 Kube-OVN

[Kube-OVN](https://kubeovn.github.io/docs/stable/) 提供 VPC 网络隔离功能。推荐参考其官方安装文档部署稳定版本。

示例: 

```bash
wget https://raw.githubusercontent.com/kubeovn/kube-ovn/refs/tags/v1.16.0/dist/images/install.sh
bash install.sh
```

## 配置 StorageClass

动态附件依赖支持 `ReadWriteMany` 的 PVC。K3S 默认 `local-path` 不支持 RWX, 需改用 NFS 等共享存储方案。

建议流程: 

1. 取消 `local-path` 的默认 StorageClass 标记。
2. 安装支持 RWX 的 StorageClass, 例如 `nfs-subdir-external-provisioner`。
3. 将 RWX 存储类设为默认, 或在 Helm `persistence.storageClass` 中显式指定。

## 启动时资源检查与创建

Helm 安装后, 应用启动时会检查或创建以下资源: 

- 命名空间: `{namespace}`
- 共享存储 PVC: `{namespace}-shared-volume`
- 外部网络 Subnet: `{namespace}-external-network-{interface}`
- 外部网络 NAD: `kube-system/{namespace}-external-network-{interface}`

其中 PVC 缺失会导致动态附件不可用。外部网络资源由后端按 `external_networks.interfaces` 创建; 若配置为空、创建失败或没有匹配节点标签, VPC 模式不可用。

## 跨云厂商节点

若部分节点不支持 VPC 网络, 可为其添加节点标签, 阻止 `CreateVPCNatGateway` 生成的 Pod 调度到这些节点: 

```bash
kubectl label node <node-name> node.cbctf.io/vpc-unsupported=true
```

该标签表示节点不支持 VPC。`CreateVPCNatGateway` 会通过节点亲和性排除这些节点, 不会影响其他服务 Pod 的调度。
