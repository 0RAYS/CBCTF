# 🏆 CBCTF

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache-green.svg)](LICENSE)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.20+-blue.svg)](https://kubernetes.io)
[![Redis](https://img.shields.io/badge/Redis-6.0+-red.svg)](https://redis.io)
[![MySQL](https://img.shields.io/badge/MySQL-8.0+-orange.svg)](https://mysql.com)

> **CBCTF** 是一个基于 Go 语言开发的高性能 CTF 竞赛平台，专为满足现代 CTF 竞赛需求而设计。

## ✨ 核心特性

### 🎯 多样化题目类型

#### 📝 题目类型
- **静态题目** - 队伍间共用附件，适合传统 CTF 题目
- **动态题目** - 实时生成唯一附件，确保每个队伍获得不同的挑战
- **动态容器** - 自动生成并启动队伍隔离的容器环境

#### 🚩 Flag 类型
| 类型 | 格式 | 说明 |
|------|------|------|
| **静态 Flag** | `static{固定内容}` | 每次生成的 flag 均相等 |
| **动态 Flag** | `dynamic{随机内容}` | 基于模板随机变化，保持可读性 |
| **UUID Flag** | `uuid{}` | 标准 UUID 格式 |

> 💡 **动态容器支持**：Flag 可注入至环境变量或作为文件挂载至指定路径

#### 🏅 分值系统
- **静态分数** - 每个 flag 的分数不随接触人数变化
- **线性分数** - 随着接触人数增加，等量减少分值
- **非线性分数** - 使用公式：`(MinScore - InitScore) / (Decay²) * (Solvers²) + InitScore`

> 🎁 **三血奖励**：一二三血奖励依次为初始分数的 5%、3%、1%

### 🚀 技术架构

- **🔄 动态附件生成** - 基于 Kubernetes 的容器化生成
- **🌐 自定义网络** - 基于 Kube-OVN 的 VPC 网络隔离
- **📧 邮件验证** - SMTP 邮件验证功能
- **📝 Writeup 管理** - 平台内 Writeup 收集与下载
- **📊 事件日志** - 完整的比赛期间动作事件记录
- **⚡ 高性能缓存** - 基于 Redis 的数据缓存
- **💾 数据存储** - MySQL 数据库存储
- **📁 文件存储** - 基于 NFS 的文件存储系统
- **🔥 预热机制** - 支持镜像预热、容器预热

## 🐳 动态附件系统

### 核心优势
- ✅ **环境隔离** - 采用 Docker 容器进行生成，与主机环境完全隔离
- ✅ **灵活配置** - 可基于通用 Docker 镜像，通过上传 Python 脚本进行生成
- ✅ **减轻负担** - 大幅减轻出题人的配置压力

### 📖 使用示例
详细示例请参考：[动态附件示例](example/dynamic/README.md)

## 🏗️ 动态容器系统

### 网络环境区分

后端通过以下方式区分题目网络环境：

| 环境类型 | 判断条件 | 配置要求 |
|----------|----------|----------|
| **Pod** | docker-compose 中未配置 `networks` 字段 | 使用默认网络 |
| **VPC** | docker-compose 中配置了 `networks` 字段 | 所有容器须配置 `networks`，手动指定 IP 地址 |

### 📁 配置示例

#### Pod 环境
```yaml
# 示例配置
version: '3'
services:
  web:
    image: nginx:alpine
    ports:
      - "80:80"
```
详细示例：[Pod 配置示例](example/pods/pod/docker-compose.yml)

#### VPC 环境
```yaml
# 示例配置
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
详细示例：[VPC 配置示例](example/pods/vpc/docker-compose.yml)

## 🔧 环境依赖

### 必需组件

| 组件 | 版本要求 | 说明 |
|------|----------|------|
| **[Kube-OVN](https://kubeovn.github.io/docs/stable/start/prepare/)** | 最新稳定版 | 自定义 VPC 网络支持 |
| **[Multus](https://github.com/k8snetworkplumbingwg/multus-cni/tree/master)** | 最新版本 | 多网络接口支持 |

### ⚠️ 重要提示

**Multus 插件选择建议：**

- ✅ **推荐使用 Thin Plugin** - 无需手动配置，稳定性更好
- ❌ **避免使用 Thick Plugin** - 容易发生以下问题：
  - [OOMKilled](https://github.com/k8snetworkplumbingwg/multus-cni/issues/1346)
  - [Text file busy](https://github.com/k8snetworkplumbingwg/multus-cni/issues/1221)

如需使用 Thick Plugin，请参考：
- [Issue #1346](https://github.com/k8snetworkplumbingwg/multus-cni/issues/1346#issuecomment-2644110944)
- [PR #1213](https://github.com/k8snetworkplumbingwg/multus-cni/pull/1213)

## 📄 许可证

本项目采用 [Apache 许可证](LICENSE)。

---

<div align="center">

**⭐ 如果这个项目对您有帮助，请给我们一个星标！**

</div>
