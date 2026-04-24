# CBCTF

[![Go Version](https://img.shields.io/badge/Go-1.26+-blue.svg)](https://golang.org)
[![React](https://img.shields.io/badge/React-19-blue.svg)](https://react.dev)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.20+-blue.svg)](https://kubernetes.io)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-blue.svg)](https://www.postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-6.0+-red.svg)](https://redis.io)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)

CBCTF 是一个基于 Go 语言开发的高性能 CTF 竞赛平台, 前后端一体化编译, 支持动态容器、动态附件、VPC 网络隔离等特性。

## 特性概览

**题目系统**
- 静态题目 / 动态题目 / 动态容器, 支持问答类题目
- Flag 类型: 静态 Flag、动态 Flag、UUID Flag, 可注入环境变量或文件挂载
- 分值系统: 静态 / 线性 / 非线性动态分值, 支持三血奖励

**竞赛管理**
- 团队管理、实时排行榜、公告系统
- Writeup 收集与批量下载
- 完整的比赛事件日志

**动态容器**
- 基于 Kubernetes 的队伍隔离容器环境
- Pod / VPC 两种网络模式, VPC 基于 Kube-OVN 实现网络隔离
- 支持 Frp 端口转发、镜像预热、容器预热、流量抓取

**动态附件生成**
- 基于 Kubernetes 的容器化隔离生成
- 可基于通用 Docker 镜像, 通过上传 Python 脚本进行生成

**安全特性**
- 设备指纹 (FingerprintJS) + JWT 双重认证
- 作弊检测（Flag 共享、设备指纹异常）
- 基于 Redis 的速率限制, 支持全局及 IP 白名单

**其他**
- 异步任务队列（邮件、Webhook、附件生成、图片处理）
- OAuth / OIDC 第三方认证
- 国际化 (i18n) 支持
- Prometheus 监控指标

## 技术栈

### 后端

| 组件        | 技术                | 版本    |
|-----------|-------------------|-------|
| 语言        | Go                | 1.26  |
| Web 框架    | Gin               | 1.12  |
| ORM       | GORM              | 1.31  |
| 数据库       | PostgreSQL        | 16+   |
| 缓存 / 消息队列 | Redis             | 6.0+  |
| 异步任务      | Asynq             | 0.26  |
| 定时任务      | robfig/cron       | v3    |
| 容器编排      | Kubernetes        | 1.20+ |
| 网络插件      | Kube-OVN / Multus | -     |
| 虚拟机编排      | KubeVirt          | -     |
| JWT       | golang-jwt        | v5    |
| 配置管理      | Viper             | 1.21  |

### 前端

| 组件       | 技术            | 版本   |
|----------|---------------|------|
| 框架       | React         | 19   |
| 构建工具     | Vite          | 7    |
| 状态管理     | Redux Toolkit | 2    |
| HTTP 客户端 | Axios         | 1.15 |
| CSS      | Tailwind CSS  | 4    |
| 路由       | React Router  | 7    |
| 国际化      | i18next       | 25   |
| 图表       | ECharts       | 6    |
| 代码编辑器    | Monaco Editor | 4    |
| 设备指纹     | FingerprintJS | 5    |

## 快速开始

### 环境要求

- Go 1.26+
- Node.js 24+ / pnpm
- PostgreSQL 16+
- Redis 6.0+
- Kubernetes 1.20+（动态容器功能需要）

### 构建

```bash
# 构建前端
cd frontend && pnpm install && pnpm run build && cd ..

# 构建后端（前端静态文件会被嵌入二进制，流量抓取功能依赖 libpcap）
CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o CBCTF .
```

### 运行

首次运行会自动生成 `config.yaml` 配置文件, 修改配置后重新启动: 

```bash
./CBCTF
```

## 配置说明

配置文件 `config.yaml` 在首次运行时从内置默认配置自动生成。主要配置段: 

| 配置段          | 说明                                    |
|--------------|---------------------------------------|
| `host`       | 后端服务对外地址                              |
| `path`       | 数据存储路径                                |
| `log`        | 日志级别 (DEBUG/INFO/WARNING/ERROR)、是否持久化 |
| `gin`        | 服务监听地址/端口、上传限制、速率限制、CORS、JWT          |
| `gorm.postgres` | PostgreSQL 连接配置、连接池参数                 |
| `redis`      | Redis 连接配置                            |
| `k8s`        | Kubeconfig 路径、命名空间、Tcpdump 镜像、Frp 配置   |
| `cheat`      | 作弊检测 IP 白名单                           |
| `webhook`    | Webhook 目标地址黑名单                       |
| `asynq`      | 异步任务并发数                               |
| `registration` | 是否允许注册、新用户默认分组 ID               |
| `geocity_db` | GeoIP 数据库路径（GeoLite2-City.mmdb）       |

支持环境变量覆盖, 前缀为 `CBCTF_`, 例如 `CBCTF_GIN_PORT=9000`。

## 题目系统

### 题目类型

- **静态题目** - 队伍间共用附件, 适合传统 CTF 题目
- **动态题目** - 实时生成唯一附件, 确保每个队伍获得不同的挑战
- **动态容器** - 自动生成并启动队伍隔离的容器环境

### Flag 类型

| 类型        | 格式              | 说明             |
|-----------|-----------------|----------------|
| 静态 Flag   | `static{固定内容}`  | 每次生成的 Flag 均相等 |
| 动态 Flag   | `dynamic{随机内容}` | 基于模板随机变化, 保持可读性 |
| UUID Flag | `uuid{}`        | 标准 UUID 格式     |

动态容器支持将 Flag 注入至环境变量或作为文件挂载至指定路径。

### 分值系统

- **静态分数** - 分值固定, 不随解题人数变化
- **线性分数** - 随解题人数增加等量递减
- **非线性分数** - 指数衰减公式: `(Score - MinScore) × e^(-5/Decay × Solvers) + MinScore`

三血奖励: 一血 / 二血 / 三血分别额外获得初始分数的 5% / 3% / 1%。

## 动态容器系统

### 网络模式

后端通过 docker-compose 配置自动区分网络环境: 

| 模式  | 判断条件              | 说明                              |
|-----|-------------------|---------------------------------|
| Pod | 未配置 `networks` 字段 | 使用默认网络, 容器间可直接通信                 |
| VPC | 配置了 `networks` 字段 | 基于 Kube-OVN 的 VPC 网络隔离, 需手动指定 IP |

### 配置示例

**Pod 模式: **

```yaml
version: '3'
services:
  web:
    image: nginx:alpine
    ports:
      - "80:80"
```

详细示例: [Pod 配置示例](example/pods/pod/docker-compose.yaml)

**VPC 模式: **

```yaml
version: '3'
services:
  web:
    image: nginx:alpine
    networks:
      vpc:
        ipv4_address: 192.168.1.10
networks:
  vpc:
    external: true
```

详细示例: [VPC 配置示例](example/pods/vpc/docker-compose.yaml)

## 动态附件系统

基于 Kubernetes 容器化生成, 支持上传 Python 脚本在隔离环境中为每个队伍生成唯一附件。

详细示例: [动态附件示例](example/dynamic/README.md)

## 环境依赖（Kubernetes）

动态容器和动态附件功能依赖以下 Kubernetes 组件: 

| 组件                                                               | 说明         |
|------------------------------------------------------------------|------------|
| [Kube-OVN](https://kubeovn.github.io/docs/stable/start/prepare/) | VPC 网络隔离支持 |
| [Multus](https://github.com/k8snetworkplumbingwg/multus-cni)     | 多网络接口支持    |
| [KubeVirt](https://kubevirt.io/)                                 | 虚拟机编排（实验性） |

**Multus 插件选择建议: **

推荐使用 **Thin Plugin**, 无需手动配置且稳定性更好。避免使用 Thick Plugin, 已知存在以下问题: 

- [OOMKilled](https://github.com/k8snetworkplumbingwg/multus-cni/issues/1346)
- [Text file busy](https://github.com/k8snetworkplumbingwg/multus-cni/issues/1221)

如需使用 Thick Plugin, 请参考: 
- [Issue #1346](https://github.com/k8snetworkplumbingwg/multus-cni/issues/1346#issuecomment-2644110944)
- [PR #1213](https://github.com/k8snetworkplumbingwg/multus-cni/pull/1213)

## 许可证

本项目采用 [Apache License 2.0](LICENSE)。
