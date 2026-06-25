# 快速上手

推荐使用 Helm 将 CBCTF 部署到 Kubernetes。

:::tip
须提前安装好所需的 k8s 组件，参考[部署](/deploy/cluster.md)

| 组件                                                               | 用途       |
| ---------------------------------------------------------------- | -------- |
| [Kube-OVN](https://kubeovn.github.io/docs/stable/start/prepare/) | VPC 网络隔离 |
| [Multus CNI](https://github.com/k8snetworkplumbingwg/multus-cni) | 多网络接口    |
| [KubeVirt](https://kubevirt.io/)                                 | 虚拟机调度    |

:::

## Helm 快速部署

:::warning
没有快速部署，必炸，请仔细阅读[部署](/deploy/cluster.md)
:::

## 初始管理员

应用启动时会自动迁移数据库。首次启动且管理员组中没有用户时，会自动创建 `admin` 用户，并把初始密码打印到日志：

```text
Init Admin: Admin{ name: admin, password: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx, email: admin@0rays.club}
```

查看日志：

```bash
kubectl logs -n cbctf deployment/cbctf | grep "Init Admin"
```

首次登录后请立即修改管理员密码。管理后台为 `/platform/#/admin`。

## 本地开发运行

本地运行需要先准备 PostgreSQL 和 Redis。


### 构建前端
```bash
cd frontend
pnpm install
pnpm run build
```
### 启动后端
```bash
cd ..
go run .
```

如果没有 `config.yaml`，程序会使用内置默认配置。系统配置页不会修改 PostgreSQL/GORM、Redis、数据目录、Gin 监听地址和监听端口。

后端构建命令：

```bash
CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o CBCTF .
```

:::tip
流量捕获能力依赖 CGO 和 libpcap；如果只做文档或前端调试，可优先使用 Helm 或已有镜像。

可通过修改 `frontend/src/api/config.js` 快速在本地实现前端调试开发

```javascript
export const API_CONFIG = {
    // 已部署服务
    BASE_URL: 'https://ctf.example.com',
};
```

需要在 `config.yaml` 或 Helm values 中添加对应的 `gin.cors` 白名单配置
:::
