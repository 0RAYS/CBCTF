---
title: 介绍
description: 了解 CBCTF 的题型、赛事管理能力、动态靶机和附件生成，选择适合的部署方式。
---

# 介绍

CBCTF 是由 [0RAYS](https://github.com/0rays) 维护的 CTF 竞赛平台，基于 Go 语言构建，原生支持 Kubernetes
编排。平台支持动态附件生成、动态容器分发、容器与虚拟机混合部署、网络渗透场景构建等特性。

<img src="/img/homepage.png" width="100%" alt="首页" />

## 功能特性

### 题目类型

| 类型                | 说明                                  |
|-------------------|-------------------------------------|
| **静态题目**          | 所有队伍共用附件，flag 相同                    |
| **动态附件**          | 容器为每个队伍独立生成附件，flag 各不相同             |
| **动态容器 · Pod 模式** | 多容器共享同一 Pod 网络，容器间通过 `localhost` 通信 |
| **动态容器 · VPC 模式** | 每个容器独立 Pod，须分配静态 IP，适合渗透场景          |

每道题目均可配置多个 flag，每个 flag 独立计分。

<img src="/img/challenges.png" width="100%" alt="题目列表" />

### Flag 类型

flag 前缀可在赛事设置中自定义（默认 `CBCTF`）：

| 类型       | 原始值                      | 实际 Flag                                       |
|----------|--------------------------|-----------------------------------------------|
| `static` | `static{this_is_a_flag}` | `CBCTF{this_is_a_flag}`                       |
| `leet`   | `leet{this_is_a_flag}`   | `CBCTF{ThiS-ls_4-fIaG}`                       |
| `uuid`   | `uuid{}`                 | `CBCTF{1301ea62-ccd2-4543-b663-993f87b6d44a}` |

### 平台能力

- **动态分值** — 一二三血额外获得题目分值的 5% / 3% / 1%
- **Frp 内网穿透** — 容器端口转发，保留原始客户端 IP
- **SMTP 邮件验证** — 注册验证与密码找回
- **Writeup 管理** — 支持收集与批量下载
- **OAuth / OIDC** — 第三方认证，支持用户组自动分配
- **平台品牌化** — Logo、名称、首页文案等全局配置
- **可重载配置** — 多数运行配置可在线保存，并通过管理后台重启重载
- **Webhook** — GET / POST
- **国际化（i18n）** — 多语言界面支持
- **Prometheus 监控** — 完整的运行时指标暴露
- **Redis 缓存 / 任务队列** + **PostgreSQL 数据存储** + **NFS 网络存储**

<img src="/img/dashboard.png" width="100%" alt="管理后台" />

<img src="/img/contest.png" width="100%" alt="比赛详情" />

<img src="/img/scoreboard-1.png" width="100%" alt="排行榜" />

<img src="/img/scoreboard-2.png" width="100%" alt="排行榜（图表）" />

<img src="/img/contest-settings.png" width="100%" alt="比赛设置" />

<img src="/img/settings.png" width="100%" alt="系统设置" />

<img src="/img/branding.png" width="100%" alt="品牌化配置" />

<img src="/img/log.png" width="100%" alt="日志" />

## 开始使用

管理员可按[快速上手](./quick-start)准备集群、安装平台并创建第一场比赛。已加入平台的选手可直接阅读[选手操作流程](../features/playing)。

## 动态容器

### 网络模式

创建容器题时，通过 Compose 配置选择网络模式：

| 模式      | 判断条件              | 说明                              |
|---------|-------------------|---------------------------------|
| **Pod** | 未配置 `networks` 字段 | 使用默认网络，容器间可直接通信                 |
| **VPC** | 配置了 `networks` 字段 | 基于 Kube-OVN 的 VPC 网络隔离，需手动指定 IP |

### 配置示例

**Pod 模式**

```yaml
version: '3'
services:
  web:
    image: nginx:alpine
    x-kubevirt: false
    ports:
      - "80:80"
```

> 完整示例：[example/pods/pod/docker-compose.yaml][pod-compose]

**VPC 模式（含 KubeVirt 虚拟机）**

```yaml
version: '3'
services:
  web:
    # 替换为自行制作的可启动 containerDisk 镜像，普通 nginx 镜像不能启动虚拟机
    image: registry.example.com/challenges/vm-web:v1
    mem_limit: 512m
    x-kubevirt: true
    x-boot:
      bootloader: efi
      secure_boot: false
    x-cloudinit:
      users:
        - name: root
    networks:
      vpc:
        ipv4_address: 192.168.1.10
        mac_address: "00:00:00:00:01:01"
networks:
  vpc:
    ipam:
      config:
        - subnet: 192.168.1.0/24
          gateway: 192.168.1.1
```

> 完整示例：[example/pods/vpc/docker-compose.yaml][vpc-compose]

<img src="/img/docker-compose.png" width="100%" alt="容器配置" />

<img src="/img/vm.png" width="100%" alt="虚拟机" />

<img src="/img/victims-1.png" width="100%" alt="靶机列表" />

<img src="/img/victims-2.png" width="100%" alt="靶机详情" />

## 动态附件

基于 Kubernetes 容器化生成，支持上传 Python 脚本，在隔离环境中为每个队伍生成唯一附件。

**出题准备：**

- 容器必须包含 `sleep` 和 `unzip`
- 脚本路径固定为 `/root/run.sh <team_id> <base64_encoded_flags>`
- 产物须写入 `/root/mnt/attachments/{id}.zip`
- 脚本完成后，队伍可在题目页面下载生成的附件
- 建议使用固定镜像版本或 digest，保持比赛期间的环境一致

> 完整示例：[example/dynamic/README.md](https://github.com/0RAYS/CBCTF/blob/main/example/dynamic/README.md)

## Kubernetes 依赖

平台部署在 Kubernetes 集群内，额外组件按题型选择：

| 组件                                                               | 用途       |
|------------------------------------------------------------------|----------|
| [Kube-OVN](https://kubeovn.github.io/docs/stable/start/prepare/) | VPC 题目的网络隔离 |
| [Multus CNI](https://github.com/k8snetworkplumbingwg/multus-cni) | VPC/VM 题目的多网络接口 |
| [KubeVirt](https://kubevirt.io/) | VM 题目的虚拟机调度 |

## 许可证

本项目采用 [GNU Affero General Public License v3.0](https://github.com/0RAYS/CBCTF?tab=AGPL-3.0-1-ov-file) 开源协议。

[pod-compose]: https://github.com/0RAYS/CBCTF/blob/main/example/pods/pod/docker-compose.yaml
[vpc-compose]: https://github.com/0RAYS/CBCTF/blob/main/example/pods/vpc/docker-compose.yaml
