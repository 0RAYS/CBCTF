# CBCTF

基于 Kubernetes 的现代化 CTF 竞赛平台

> 动态容器 · 虚拟机靶机 · 动态附件 · 网络渗透场景构建

[快速上手](/guide/start/quick-start) | [查看介绍](/guide/start/introduction)

## Features

- [🐳 **动态容器与虚拟机**](/guide/features/container): 支持 Pod、VPC 网络隔离模式以及 KubeVirt VM 靶机，覆盖 Web、Pwn、渗透等各类场景。
- [📦 **动态附件生成**](/guide/features/attachment): 每支队伍独立生成含唯一 Flag 的附件，防止抄答案，支持 Python 自定义生成逻辑。
- [🔍 **多维度作弊检测**](/admin/cheat): 内置设备指纹、IP、跨队 Flag 提交等五种检测机制，自动标记可疑行为。
- [📊 **灵活计分系统**](/guide/features/scoring): 静态、线性、对数三种计分类型，支持三血奖励和多 Flag 独立计分。
- [⚡ **Helm 一键部署**](/deploy/helm): Chart 内置 PostgreSQL 和 Redis，支持 PVC 持久化、Ingress TLS 和镜像拉取凭据。
- [🔐 **OAuth / OIDC 认证**](/guide/features/auth): 支持多个第三方认证提供商，可配置用户组自动分配，兼容 GitHub 等标准 OAuth 流程。
