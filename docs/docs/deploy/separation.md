---
title: 前后端分离
description: 配置 CBCTF 前后端分离部署、反向代理、静态资源和 API 地址。
---

# 前后端分离

CBCTF 默认将前端静态资源嵌入 Go 二进制，通过同一个服务地址提供 `/platform`。如果需要把前端单独托管到
CDN、静态站点或其他域名，可以采用前后端分离方式。

## 适用场景

- 需要 CDN 加速前端资源
- 前端需部署到独立域名
- 希望前后端分别发布

## 前端构建

```bash
git clone https://github.com/0RAYS/CBCTF.git
cd CBCTF/frontend
pnpm install
```

修改 `frontend/src/api/config.js`：

仓库默认 `BASE_URL: ''`，表示使用页面同源的 API。这里填写后端 Origin，不附加 `/platform` 或不存在的 `/api` 前缀。

```javascript
export const API_CONFIG = {
    BASE_URL: 'https://api.ctf.example.com',
};
```

然后构建：

```bash
pnpm build
```

构建完成后，将 `frontend/dist/` 内容部署到静态站点的 `/platform/` 路径。Vite 的 `base` 固定为 `/platform/`，前端路由由 HashRouter 处理，例如 `/platform/#/games`；若改部署子路径，需要同步修改 Vite base。开发服务器没有内置 API proxy。

## 后端配置

`config.yaml` 示例：

```yaml
host: https://api.ctf.example.com

gin:
  origins:
    - https://ctf.example.com
```

- `host` 必须填写后端真实对外地址，OAuth 回调与邮件链接都会使用它
- `gin.origins` 需要包含前端独立域名对应的浏览器 `Origin`，否则跨域请求和认证 cookie 可能无法正常工作
- 现有数据库上的这些值应在「系统管理」修改；CORS 需要重启生效，仅修改 Helm values 不会覆盖数据库。
- API 客户端使用 `withCredentials`，认证存储在 HttpOnly Cookie 中。允许的跨域 Origin 会使用 `SameSite=None; Secure`，需要 HTTPS；浏览器的第三方 Cookie 策略仍可能阻止真正跨站点的登录。

## OAuth 注意事项

前后端分离时，OAuth 回调链路为：

1. 第三方回调到后端 `https://api.ctf.example.com/oauth/{uri}/callback`
2. 后端完成登录后重定向到 `https://api.ctf.example.com/platform/#/oauth/callback?...`

:::warning
当前代码默认仍依赖后端提供 `/platform` 下的前端回调页。若完全拆离前端托管位置，需要同步调整 OAuth 回调后的前端跳转逻辑。
:::

邮箱验证和密码重置链接也使用 `{host}/platform/#/verify` 与 `{host}/platform/#/reset-password`。一种无需改登录逻辑的部署方式是在同一公开域名上，将 `/platform/` 转发到静态托管，将其他路径转发到后端；此时前端继续使用同源 API。

## Helm 场景

若后端仍通过 Helm 部署，只需在 `values.yaml` 中设置：

```yaml
cbctf:
  host: "https://api.ctf.example.com"
  gin:
    origins:
      - "https://ctf.example.com"
```
