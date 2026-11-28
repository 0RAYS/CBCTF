---
title: 渗透靶场
description: 说明 CBCTF 渗透靶场的实例池、租约、玩家网关、网络拓扑与出题人配置方式。
---

# 渗透靶场

渗透靶场是容器题的一种形态：一整套复杂网络拓扑打包成一个实例，环境完全相同，用完即删。出题人只负责设计环境与拓扑，平台负责批量启停实例并把环境分发给选手。

与普通容器题的区别：

| 维度     |   普通容器题    |   渗透靶场    |
|---------|:-----------:|:----------:|
| 实例数量    | 每队一份，按需创建    | 常驻一批相同实例，队伍领取 |
| Flag     | 写进实例环境          | 平台下发到玩家网关，环境里没有 |
| 对外入口    | 每个暴露端口一个端点       | 只有玩家网关一个入口    |
| 使用后     | 停止即销毁            | 租约到期销毁，不回收复用   |
| 网络      | 可选 VPC 模式         | 多网段 + 三层互通 + 静态地址 |

## 概念

| 概念     | 载体                | 说明                             |
|---------|-------------------|--------------------------------|
| 靶场题    | `challenges.range_config` | 出题人声明的节点、网段、互通关系与玩家网关   |
| 靶场实例   | `victims`，`kind=range` | 一套跑起来的拓扑，每个实例的内容完全相同       |
| 玩家网关   | 租约级 Pod           | 每队一台，选手唯一拿到的入口，承载该队的 Flag |
| 租约     | `victims.team_id` 与 `lease_expire_at` | 一个实例同时只属于一个队伍，到期回收 |

实例用完即销毁，不做回收复用：渗透过程会改环境（留后门、改配置），复用就不再是同一套环境。租约结束后平台补建一个新的同构实例。

## 出题人配置

靶场题仍是 `pods` 类型，在 Compose 顶层用 `x-range` 声明调度与接入参数，在 service 上用 `x-role` 声明节点角色。

```yaml
x-range:
  pool_size: 4            # 常驻就绪的实例数，0 表示不预热
  max_leased: 20          # 同时租出的上限，0 表示不限
  lease_duration: 7200    # 单次租约时长（秒）
  gateway: jump           # 玩家网关对应的 service
  access_network: mgmt    # 网关接入的网段，省略时取网关唯一的网段
  routes:                 # 网段间的三层互通声明
    - {from: mgmt, to: internal}
    - {from: internal, to: dmz}

services:
  jump:
    image: registry.example.com/challenges/range-gateway:latest
    x-role: gateway
    ports:
      - target: 22
        published: ssh
    networks:
      mgmt:
        ipv4_address: 192.168.99.254

  web:
    image: registry.example.com/challenges/web:latest
    x-role: node
    networks:
      internal:
        ipv4_address: 10.10.0.2

  db:
    image: registry.example.com/challenges/db:latest
    x-role: node
    networks:
      dmz:
        ipv4_address: 10.20.0.2

networks:
  mgmt:
    ipam:
      config:
        - subnet: 192.168.99.0/24
          gateway: 192.168.99.1
  internal:
    ipam:
      config:
        - subnet: 10.10.0.0/24
          gateway: 10.10.0.1
  dmz:
    ipam:
      config:
        - subnet: 10.20.0.0/24
          gateway: 10.20.0.1
```

`x-range` 字段：

| 字段             | 要求                                                    |
|----------------|-------------------------------------------------------|
| `pool_size`    | 大于等于 0，常驻就绪的实例数，选手领取后平台自动补建               |
| `max_leased`   | 大于等于 0，同时租出的上限，0 表示不限                       |
| `lease_duration` | 大于 0，秒；省略时为 7200                                 |
| `gateway`      | 必填，必须指向存在的 service，该 service 需要有端口且不能是 `x-kubevirt` |
| `access_network` | 必须是已定义的网段；省略时取网关唯一的网段，因此网关只能接一个网段        |
| `routes`       | 每条的两个网段都必须已定义、不能相同、不能重复；省略表示各网段之间不互通   |

接入网段是选手进入拓扑的落脚点，节点在别的网段时必须用互通声明把接入网段连过去，否则选手够不到节点，平台会直接拒绝这份配置。互通是逐跳的：接入网段只连到 `internal` 时，选手到 `dmz` 要经 `internal` 转发。

`x-role` 取值为 `node`、`router`、`gateway`，用于拓扑图展示与角色约束：`gateway` 角色只能标在 `x-range.gateway` 指向的 service 上，声明了 `x-role` 却没有 `x-range` 的题目会被拒绝。靶场题的节点不需要暴露端口，端口由玩家网关承担。

## 虚拟机节点

节点可以是容器，也可以是虚拟机。声明 `x-kubevirt: true` 的 service 用 KubeVirt 起一台虚拟机，镜像需要是容器磁盘格式的 guest 系统，例如 `quay.io/kubevirt/fedora-cloud-container-disk-demo`。虚拟机节点的端口不对外暴露，出题人给它的端口声明会被忽略；选手进入虚拟机的唯一路径是经过玩家网关。

guest 里的初始状态由 `x-cloudinit` 声明：

| 字段 | 作用 |
| --- | --- |
| `users` / `groups` | guest 里的用户、用户组、sudo 权限与登录口令 |
| `write_files` | 写进 guest 的文件，路径、内容、权限与属主 |
| `ssh_authorized_keys` | guest 用户的公钥 |

平台按出题人规划的静态地址与路由渲染 cloud-init 的 network-data：网卡按 MAC 匹配，地址与拓扑一致，到互通网段的下一跳是本网段网关。guest 里只有靶场网卡，没有集群默认路由可用，跨网段可达性完全取决于这些静态路由。

写文件的属主由平台放到 cloud-init 的 `runcmd` 里用 `chown` 设置，而不是写进 `write_files`：`write-files` 在部分镜像里先于 `users-groups` 执行，那时目标用户还不存在，cloud-init 会因为改不了属主让整个模块报错，连用户都不会被创建。用户主目录同样会补一次 `chown`，避免被提前建出来的目录留在 root 名下。

三层递进的完整例子：玩家从接入网段的网关进对外层跳板机，跳板机给出内网堡垒机的口令，堡垒机再给出核心机的口令，每层只在上一层拿到的网段里可达。

```yaml
x-range:
  pool_size: 1
  max_leased: 2
  lease_duration: 3600
  gateway: jump
  access_network: access
  routes:
    - {from: access, to: dmz}
    - {from: dmz, to: core}
    - {from: core, to: vault}

services:
  jump:
    image: nicolaka/netshoot:latest
    x-role: gateway
    ports:
      - target: 22
        published: ssh
    networks:
      access:
        ipv4_address: 192.168.99.254

  edge:
    image: quay.io/kubevirt/fedora-cloud-container-disk-demo:latest
    x-kubevirt: true
    x-cloudinit:
      users:
        - name: ctf
          groups: [sudo]
          sudo: ["ALL=(ALL) NOPASSWD:ALL"]
          shell: /bin/bash
          plain_text_passwd: edge-pass
      write_files:
        - path: /home/ctf/hint.txt
          owner: ctf:ctf
          permissions: "0644"
          content: |
            第二层在内网段 10.20.0.0/24,堡垒机 10.20.0.2,用户 ctf,口令 bastion-pass
    networks:
      dmz:
        ipv4_address: 10.10.0.2
        mac_address: "00:00:5e:00:02:01"

  bastion:
    image: quay.io/kubevirt/fedora-cloud-container-disk-demo:latest
    x-kubevirt: true
    x-cloudinit:
      users:
        - name: ctf
          plain_text_passwd: bastion-pass
      write_files:
        - path: /home/ctf/hint.txt
          owner: ctf:ctf
          permissions: "0644"
          content: |
            核心层 10.30.0.0/24,核心机 10.30.0.2,用户 ctf,口令 vault-pass
    networks:
      core:
        ipv4_address: 10.20.0.2
        mac_address: "00:00:5e:00:03:01"

  vault:
    image: quay.io/kubevirt/fedora-cloud-container-disk-demo:latest
    x-kubevirt: true
    x-cloudinit:
      users:
        - name: ctf
          plain_text_passwd: vault-pass
      write_files:
        - path: /home/ctf/asset.txt
          owner: ctf:ctf
          permissions: "0644"
          content: |
            内网核心资产:core-asset-9f3a
    networks:
      vault:
        ipv4_address: 10.30.0.2
        mac_address: "00:00:5e:00:04:01"

networks:
  access:
    ipam:
      config:
        - subnet: 192.168.99.0/24
          gateway: 192.168.99.1
  dmz:
    ipam:
      config:
        - subnet: 10.10.0.0/24
          gateway: 10.10.0.1
  core:
    ipam:
      config:
        - subnet: 10.20.0.0/24
          gateway: 10.20.0.1
  vault:
    ipam:
      config:
        - subnet: 10.30.0.0/24
          gateway: 10.30.0.1
```

这份配置里 `dmz` 与 `core` 互通、`core` 与 `vault` 互通，但 `dmz` 与 `vault` 之间没有声明互通，因此两层之间只能一跳一跳地走。

虚拟机起得比容器慢，实例就绪的超时按 `k8s.range.ready_timeout` 计。

## 领取与回收

1. 队伍在题目页点「启动」。平台先看实例池里有没有空闲的已就绪实例，有就直接租给队伍；没有就现场启动一个。
2. 实例就绪后平台创建玩家网关，把该队的 Flag 渲染后挂到网关的 `/run/range/payload`，并为网关开一个 NodePort，页面返回的就是这个地址。
3. 队伍用网关地址进入拓扑。拓扑节点本身不对外暴露，选手只能从网关访问。
4. 停止或租约到期后平台销毁网关与实例，实例池补建一个新的同构实例。

实例池由定时任务维护，每分钟检查一次，把就绪实例补到 `pool_size`。后台的批量开题对靶场题同样生效：平台给每队租一个实例，池里没有空闲实例就现场启动一个；靶机列表里标出「靶场」，鼠标悬停可以看到该实例的拓扑指纹，同一题目的实例指纹相同。

后台「靶机」页顶部的实例池面板列出每个靶场题的池容量、已租出数量与每个实例的状态、占用队伍、租约到期时间和拓扑指纹，并提供「预热」按钮把池补齐到 `pool_size`。对应接口是 `GET /admin/contests/:contestID/ranges` 与 `POST /admin/contests/:contestID/ranges`，权限 `admin:contest_range:control`。面板只做观测与手动触发，实例的启停仍由平台调度。

## Flag 下发

靶场环境本体不携带 Flag，实例因此与队伍无关，池中的实例才是真正相同的。Flag 在领取时按队伍渲染，随网关的载荷一起下发：

| 载荷              | 内容                                          |
|-----------------|---------------------------------------------|
| `/run/range/payload/flags/*` | 该队的 Flag 文件                       |
| `/run/range/payload/payload.json` | 凭据、环境变量与文件清单                       |

Flag 的取值规则与普通题目一致：`static{}` 所有队伍相同，`leet{}` 与 `uuid{}` 每队不同。

网关容器同时带着出题人在 compose 里声明的环境变量与按队伍渲染后的 Flag，同名键以渲染结果为准；环境本体里不带 Flag，因此同一题目的池实例内容完全相同。

## 网络拓扑

靶场实例的每个网段是一个 Kube-OVN 子网，每个节点在网段里有一个静态地址。跨网段的可达性由两层机制共同决定：

| 机制         | 作用                                                  |
|------------|-----------------------------------------------------|
| 静态路由        | 节点到互通网段的下一跳是本网段的网关（VPC 逻辑路由器）              |
| NetworkPolicy | 闸门在策略层：只有声明互通的网段 CIDR 才会被双向放行，其余一律拒绝        |

平台不给靶场子网设置 ACL，跨网段转发由 VPC 逻辑路由器完成；真正挡住不该互通的方向是平台生成的 NetworkPolicy。删掉 `x-range.routes` 里的声明后重建实例，节点之间就不再互通。

玩家接入网段也是实例的一个子网，由实例创建、随实例销毁，网关在租约时挂到上面。网关自己带一份出站策略，出站范围与出题人给该 service 声明的策略一致。

靶场网段一律以附加网卡的形式接入，默认路由仍属于集群网络，跨网段只依赖上面那条静态路由。因此 Kube-OVN 作主 CNI 或副 CNI 时行为一致，靶场节点也不会因为网络实现不同而改变可达范围。平台在实例就绪前会比对 Pod 各网卡的实际地址，与出题人规划的静态地址不符即判定启动失败，避免拓扑静默跑偏。

## 容量与调参

| 配置                                       | 默认值  | 说明                              |
|------------------------------------------|-------|---------------------------------|
| `k8s.range.node_concurrency`              | 8     | 单个实例同时创建的节点数上限                 |
| `k8s.range.ready_timeout`                 | 600   | 实例就绪超时（秒），含 VM 启动时间             |
| `asynq.queues.range`                      | 2     | 靶场启停队列的并发度，与每队靶机队列分开          |

Helm 对应 `cbctf.k8s.range.node_concurrency`、`cbctf.k8s.range.ready_timeout` 与 `cbctf.asynq.queues.range`。

实例池的常驻数量与同时租出上限共同决定集群压力：每个实例的节点数、CPU 与内存之和乘以 `pool_size` 就是常驻开销，`max_leased` 决定比赛高峰期的峰值。

## 排查

| 问题            | 排查方向                                                             |
|---------------|--------------------------------------------------------------------|
| 队伍启动后一直排队   | 实例池为空时平台会现场启动一个，检查就绪超时与节点资源；`pool_size` 偏小会明显增加等待           |
| 网关地址打不开      | 检查网关 Pod 是否就绪、NodePort 是否被节点放行、集群 CNI 是否支持 NodePort 到靶场网段的转发        |
| 节点之间不通       | 检查 `x-range.routes` 是否声明了这两个网段；互通是显式声明的，不是自动的                 |
| 实例启动失败      | 检查应用日志中的地址校验结果，静态地址没被 CNI 采纳时实例不会进入就绪态                   |
| Flag 拿不到      | Flag 由平台下发到玩家网关，检查网关的 `/run/range/payload` 是否挂载成功                    |
| 虚拟机节点起不来    | 查虚拟机 Pod 的 `guest-console-log` 容器日志，guest 的启动过程都在这里              |
| guest 里用户没建出来  | 查 guest 的 `/var/log/cloud-init.log`，`config-write-files` 报错说明模块顺序与属主声明冲突 |
