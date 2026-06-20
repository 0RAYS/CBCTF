# K8s 集群搭建

本文介绍如何搭建支持 CBCTF 动态附件题和容器题的 Kubernetes 集群。

## 硬件要求

准备多台运行高版本内核的 Linux 服务器，下文以一台 master 与 两台 worker 组成的 [K3S](https://docs.k3s.io/) 集群为例。

本教程主要关注平台所需的集群依赖，对于 `KubeVirt` / `Kube-OVN` / `Multus` 的具体使用和细节不在文档内容中涉及，更多内容请参考各自官方文档。

### KubeVirt 节点要求

VM 靶机依赖 KubeVirt。集群节点需要支持硬件虚拟化或可用的嵌套虚拟化能力，或是由 `KubeVirt` 实现**软件嵌套虚拟化**，但会影响其性能。

检查节点虚拟化能力：

```bash
egrep -c '(vmx|svm)' /proc/cpuinfo
```

返回值大于 `0` 通常表示节点 CPU 暴露了虚拟化能力。若运行在云厂商或虚拟化平台中，还需要确认该环境允许嵌套虚拟化。

### 安装 NFS 客户端

动态附件与题目文件共享依赖 RWX 存储。所有节点需要安装 NFS 客户端：

```bash
sudo apt update
sudo apt install -y nfs-common
```

## 安装 [Helm](https://helm.sh/zh/docs/intro/install/)

```bash
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-4
chmod 700 get_helm.sh
./get_helm.sh
```

## 安装 [K3S](https://docs.k3s.io/)

### Master 节点

```bash
curl -sfL https://rancher-mirror.rancher.cn/k3s/k3s-install.sh | INSTALL_K3S_MIRROR=cn sh -
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

## 安装 [Kube-OVN](https://kubeovn.github.io/docs/stable/)

提供 VPC 网络隔离功能。

为节点打上 Label，参考 [Helm chart for Kube-OVN](https://kubeovn.github.io/kube-ovn/)

```bash
kubectl label node -lbeta.kubernetes.io/os=linux kubernetes.io/os=linux --overwrite
kubectl label node -lnode-role.kubernetes.io/control-plane kube-ovn/role=master --overwrite

# 以下 label 用于 dpdk 镜像的安装，非 dpdk 情况，可以忽略
kubectl label node -lovn.kubernetes.io/ovs_dp_type!=userspace ovn.kubernetes.io/ovs_dp_type=kernel --overwrite
```

添加 Helm Repo

```bash
helm repo add kube-ovn https://kubeovn.github.io/kube-ovn/
```

根据集群网络规划，可以把 Kube-OVN 安装为主 CNI，也可以安装为副 CNI。


**作为 主 CNI 安装**

主 CNI 模式由 Kube-OVN 接管集群默认 Pod 网络，适合新建集群或计划完全使用 Kube-OVN 作为主网络插件的环境。
:::warning
主 CNI 安装会影响集群默认网络。已有生产集群切换主 CNI 风险较高，建议仅在新集群中使用。
:::
如果使用 K3S 并希望 Kube-OVN 作为主 CNI，安装 K3S 时应禁用默认 flannel 和 network-policy：
```bash
curl -sfL https://rancher-mirror.rancher.cn/k3s/k3s-install.sh | \
  INSTALL_K3S_MIRROR=cn sh - --flannel-backend=none --disable-network-policy
```
使用 kube-ovn-v2 chart
values.yaml
```yaml
# master ip
masterNodes:
    - 192.168.3.1
```
安装
```bash
helm install kube-ovn kube-ovn/kube-ovn-v2 --namespace kube-system -f values.yaml
```


**作为 副 CNI 安装**

副 CNI 模式保留集群原有主 CNI，只通过 Multus 为 CBCTF 靶机额外挂载 Kube-OVN 网络。这个模式更适合已有集群。
:::tip
如果集群已经有稳定主 CNI，优先考虑副 CNI 模式。它对现有业务网络影响更小，也更符合多网络靶机场景。
:::
使用 kube-ovn-v2 chart
values.yaml
```yaml
cni:
    configPriority: "10"
    nonPrimaryCNI: true
# master ip
masterNodes:
    - 192.168.3.1
```
安装
```bash
helm install kube-ovn kube-ovn/kube-ovn-v2 --namespace kube-system -f values.yaml
```


## 安装 [Multus CNI](https://github.com/k8snetworkplumbingwg/multus-cni)

为 VPC 网络模式靶机提供多网卡支持：

```bash
kubectl apply \
  -f https://raw.githubusercontent.com/k8snetworkplumbingwg/multus-cni/master/deployments/multus-daemonset-thick.yml
```

## 安装 [KubeVirt](https://kubevirt.io/)

为 VM 靶机提供 `VirtualMachine` 资源和运行时能力。建议按 KubeVirt 官方文档安装稳定版本，并确认 `kubevirt` 命名空间中的组件 Ready。

```bash
export RELEASE=$(curl https://storage.googleapis.com/kubevirt-prow/release/kubevirt/kubevirt/stable.txt)
kubectl apply -f https://github.com/kubevirt/kubevirt/releases/download/${RELEASE}/kubevirt-operator.yaml
kubectl apply -f https://github.com/kubevirt/kubevirt/releases/download/${RELEASE}/kubevirt-cr.yaml

# 可选，当前项目未使用 DV
export TAG=$(curl -s -w %{redirect_url} https://github.com/kubevirt/containerized-data-importer/releases/latest)
export VERSION=$(echo ${TAG##*/})
kubectl create -f https://github.com/kubevirt/containerized-data-importer/releases/download/$VERSION/cdi-operator.yaml
kubectl create -f https://github.com/kubevirt/containerized-data-importer/releases/download/$VERSION/cdi-cr.yaml
```

:::info
CBCTF 只会创建和删除 `VirtualMachine` 资源，不会自动安装 KubeVirt。VM 题目镜像需要由出题人制作成 KubeVirt 可启动的 `containerDisk` 镜像。
:::

## 配置 StorageClass

动态附件依赖支持 `ReadWriteMany` 的 PVC。K3S 默认 `local-path` 不支持 RWX，需改用 NFS 等共享存储方案。

建议流程：

1. 取消 `local-path` 的默认 StorageClass 标记
2. 安装支持 RWX 的 StorageClass
3. 将 RWX 存储类设为默认，或在 Helm `persistence.storageClass` 中显式指定

以安装 `nfs-subdir-external-provisioner` 为例

添加 Helm Repo

```bash
helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/
```

values.yaml

```yaml
nfs:
    # 共享存储路径
    path: /volume1/kubernetes
    server: 192.168.8.1
replicaCount: 3
storageClass:
    defaultClass: true
    reclaimPolicy: Retain
```

安装

```
helm install nfs-provisioner \
  nfs-subdir-external-provisioner/nfs-subdir-external-provisioner \
  --namespace kube-system \
  -f values.yaml
```
