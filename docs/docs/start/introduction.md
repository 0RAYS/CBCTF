---
sidebar_position: 1
---

# 简介

[![Go Version](https://img.shields.io/badge/Go-1.26-blue.svg)](https://golang.org)
[![MySQL](https://img.shields.io/badge/MySQL-8.0+-orange.svg)](https://mysql.com)
[![Redis](https://img.shields.io/badge/Redis-6.0+-red.svg)](https://redis.io)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.20+-blue.svg)](https://kubernetes.io)
[![License](https://img.shields.io/badge/License-Apache_2.0-green.svg)](https://github.com/0RAYS/CBCTF/blob/main/LICENSE)

**CBCTF** 是由 0RAYS 开发的 Kubernetes 原生 CTF 竞赛平台。平台使用 Go 语言编写，将 React 前端嵌入单一二进制，部署简单，同时深度集成 Kubernetes，支持动态容器靶机、VPC 网络隔离、流量捕获等高级功能。

## 核心功能

| 功能 | 问答题 | 静态题 | 动态附件题 | 容器题 |
|------|:------:|:------:|:----------:|:------:|
| 共享附件 | ✓ | ✓ | — | — |
| 动态附件（per-team） | — | — | ✓ | — |
| 静态 Flag | ✓ | ✓ | ✓ | ✓ |
| 动态 Flag | — | — | ✓ | ✓ |
| UUID Flag | — | — | ✓ | ✓ |
| 容器靶机 | — | — | — | ✓ |
| Pod 网络 | — | — | — | ✓ |
| VPC 网络隔离 | — | — | — | ✓ |
| 流量捕获 | — | — | — | ✓ |
| 需要 Kubernetes | — | — | ✓ | ✓ |

## 架构

CBCTF 后端采用 Go 1.26 + Gin 框架，前端使用 React 19 + Vite 7，生产构建产物通过 `go:embed` 嵌入 Go 二进制。平台运行时依赖 MySQL 存储数据、Redis 缓存与任务队列，并通过 client-go 与 Kubernetes 集群交互以管理容器靶机。

## 依赖说明

| 组件 | 版本要求 | 用途 |
|------|---------|------|
| MySQL | 8.0+ | 主数据存储 |
| Redis | 6.0+ | 缓存、Asynq 任务队列 |
| Kubernetes | 1.20+ | 动态附件生成、容器靶机 |
| KubeOVN | 推荐 v1.14.5 | VPC 网络隔离（可选） |
| Multus CNI | 最新稳定版 | 多网卡支持，VPC 模式必须（可选） |

> Kubernetes、KubeOVN 和 Multus 仅在使用动态附件题或容器题时需要。若仅使用问答题和静态题，只需 MySQL 和 Redis。

## 技术栈

| 层 | 技术 |
|----|------|
| 后端语言 | Go 1.26 |
| Web 框架 | Gin v1.11 |
| 前端框架 | React 19 + Vite 7 |
| 任务队列 | Asynq v0.26（基于 Redis） |
| ORM | GORM v1.31 + MySQL 驱动 |
| K8s SDK | client-go v0.35 |
| 网络隔离 | KubeOVN + Multus CNI |
| 认证 | JWT（golang-jwt/jwt v5）+ 设备指纹 |
| 流量捕获 | gopacket v1.5（需 libpcap） |
| 地理信息 | MaxMind GeoLite2-City |
