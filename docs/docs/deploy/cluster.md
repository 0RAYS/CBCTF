# K8s 集群搭建

本文介绍如何搭建支持 CBCTF 动态附件题和容器题的 Kubernetes 集群。

## 硬件要求

准备多台运行 Ubuntu 22.04+ 的 Linux 服务器，分为：

- **Master 节点** × 1：管理 K8s 控制平面
- **Worker 节点** × N：运行题目容器，建议至少 2 台

### KubeVirt 节点要求

VM 靶机依赖 KubeVirt。集群节点需要支持硬件虚拟化或可用的嵌套虚拟化能力，并且 KubeVirt 控制组件、CRD 和 virt-launcher Pod
能正常运行。

检查节点虚拟化能力：

```bash
egrep -c '(vmx|svm)' /proc/cpuinfo
```

返回值大于 `0` 通常表示节点 CPU 暴露了虚拟化能力。若运行在云厂商或虚拟化平台中，还需要确认该环境允许嵌套虚拟化。

## 安装 NFS 客户端

动态附件与题目文件共享依赖 RWX 存储。所有节点需要安装 NFS 客户端：

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

安装完成后：

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

[Multus CNI](https://github.com/k8snetworkplumbingwg/multus-cni) 为 VPC 网络模式提供多网卡支持：

```bash
kubectl apply -f https://raw.githubusercontent.com/k8snetworkplumbingwg/multus-cni/master/deployments/multus-daemonset.yml
```

:::tip
若 K3S + Flannel 环境出现 `cannot find valid master CNI config`，需按 Multus 官方文档调整 `--multus-kubeconfig-file-host`
与 CNI 目录。
:::

## 安装 Kube-OVN

[Kube-OVN](https://kubeovn.github.io/docs/stable/) 提供 VPC 网络隔离功能。推荐参考其官方安装文档部署稳定版本。

```bash
wget https://raw.githubusercontent.com/kubeovn/kube-ovn/refs/tags/v1.16.0/dist/images/install.sh
bash install.sh
```

## 安装 KubeVirt

[KubeVirt](https://kubevirt.io/) 为 VM 靶机提供 `VirtualMachine` 资源和运行时能力。建议按 KubeVirt 官方文档安装稳定版本，并确认
`kubevirt` 命名空间中的组件 Ready。

```bash
kubectl get kubevirt -A
kubectl get pods -n kubevirt
kubectl api-resources | grep virtualmachines
```

:::info
CBCTF 只会创建和删除 `VirtualMachine` 资源，不会自动安装 KubeVirt。VM 题目镜像需要由出题人制作成 KubeVirt 可启动的
`containerDisk` 镜像。
:::

## 配置 StorageClass

动态附件依赖支持 `ReadWriteMany` 的 PVC。K3S 默认 `local-path` 不支持 RWX，需改用 NFS 等共享存储方案。

建议流程：

1. 取消 `local-path` 的默认 StorageClass 标记
2. 安装支持 RWX 的 StorageClass，例如 `nfs-subdir-external-provisioner`
3. 将 RWX 存储类设为默认，或在 Helm `persistence.storageClass` 中显式指定

## 启动时资源检查与创建

Helm 安装后，应用启动时会检查或创建以下资源：

- 命名空间：`{namespace}`
- 共享存储 PVC：`{namespace}-shared-volume`

:::warning
PVC 缺失会导致动态附件不可用。KubeVirt 资源不会在启动时创建，只有启动包含 `x-kubevirt: true` 的 VM 靶机时才会创建对应
`VirtualMachine`。
:::
