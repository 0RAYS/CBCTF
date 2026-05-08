---
sidebar_position: 1
---

# 项目简介

CBCTF 是一个面向 CTF 比赛的管理与参赛平台。平台提供管理端和选手端，后端使用 Go、Gin、GORM，前端使用 React、Vite，运行时依赖 PostgreSQL 和 Redis；动态附件和动态靶机功能依赖 Kubernetes。

前端构建产物默认嵌入后端服务，并通过 `/platform/` 对外提供页面。选手和管理员通常从 `/platform/#/login` 登录。

## 主要功能

| 功能 | 说明 |
|---|---|
| 用户注册和登录 | 支持本地账号，配置后可使用 OAuth/OIDC 登录 |
| 用户和队伍 | 用户资料、队伍创建、加入、成员管理、队伍头像 |
| 比赛管理 | 比赛时间、隐藏状态、队伍人数、规则、奖项、时间线、公告 |
| 题目管理 | 静态题、动态附件题、容器靶机题，支持分类和附件上传 |
| Flag 管理 | 静态 Flag、动态 Flag、UUID Flag，容器题可注入环境变量或文件 |
| 提交和排行榜 | Flag 提交、解题状态、比赛排行榜、排名时间线 |
| 动态靶机 | 基于 Kubernetes 为队伍启动独立靶机，支持启动、延长、停止 |
| 动态附件 | 使用 Kubernetes Job 运行生成器，为队伍生成附件 |
| 管理后台 | 用户、分组、角色权限、题库、比赛、公告、作弊记录、靶机、生成器、镜像预热、文件、任务、日志、系统配置 |
| 通知集成 | SMTP 邮件、Webhook、OAuth/OIDC 提供商 |

## 技术栈

| 组件 | 当前项目使用 |
|---|---|
| 后端 | Go 1.26、Gin、GORM |
| 前端 | React 19、Vite 7、JavaScript、pnpm |
| 数据库 | PostgreSQL |
| 缓存和任务队列 | Redis、Asynq |
| 靶机调度 | Kubernetes |
| 可选网络能力 | Kube-OVN、Multus CNI、FRP |
| 部署 | Docker 镜像、Helm Chart |

## 访问路径

| 页面 | 路径 |
|---|---|
| 首页 | `/platform/#/` |
| 登录 | `/platform/#/login` |
| 比赛列表 | `/platform/#/games` |
| 比赛详情 | `/platform/#/contests/:contestId` |
| 个人设置 | `/platform/#/settings` |
| 管理后台 | `/platform/#/admin` |

根路径 `/` 会重定向到 `/platform`。如果通过 Ingress、反向代理或端口转发访问，排障时应优先确认 `/platform/` 路径是否可访问。

## 环境要求

| 场景 | 要求 |
|---|---|
| 基础运行 | PostgreSQL、Redis |
| Helm 部署 | Kubernetes、Helm、可用 StorageClass |
| 文件和动态附件 | 可写的 `/app/data` 数据卷，动态附件建议使用 RWX 存储 |
| 动态靶机 | Kubernetes 命名空间、RBAC、可拉取题目镜像 |
| VPC 靶机 | Kube-OVN、Multus CNI、外部网络节点标签 |
| 本地构建 | Go 1.26、Node.js 24、pnpm、CGO、libpcap 开发库 |

如果只使用静态题，平台仍需要 PostgreSQL 和 Redis，但不一定需要 Kubernetes 靶机能力。当前 Helm Chart 默认会在 Kubernetes 中部署应用、PostgreSQL、Redis 和所需 RBAC。
