> For AI agents: the complete documentation index is available at /llms.txt, the full documentation bundle is available at /llms-full.txt.

# 前后端分离

CBCTF 默认通过同一个服务地址提供页面和 API。如果需要将页面托管到 CDN、静态站点或其他域名，可以采用前后端分离部署。

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

`BASE_URL: ''` 表示页面与 API 使用同一域名。独立部署时填写 API 服务的完整域名，例如 `https://api.ctf.example.com`，不附加路径。

```javascript
export const API_CONFIG = {
    BASE_URL: 'https://api.ctf.example.com',
};
```

然后构建：

```bash
pnpm build
```

构建完成后，将 `frontend/dist/` 内容部署到静态站点的 `/platform/` 路径。比赛列表地址为 `/platform/#/games`；请保留这一部署路径，以便页面资源和链接正常加载。

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
- 页面和 API 均需使用 HTTPS，并允许浏览器携带登录 Cookie。浏览器限制第三方 Cookie 时，跨站点登录可能失败，建议使用同一公开域名提供页面和 API。

## OAuth 注意事项

前后端分离时，OAuth 回调链路为：

1. 第三方回调到后端 `https://api.ctf.example.com/oauth/{uri}/callback`
2. 后端完成登录后重定向到 `https://api.ctf.example.com/platform/#/oauth/callback?...`

:::warning
`host` 域名下的 `/platform/` 页面仍需可访问，第三方登录完成后会返回这里。可通过反向代理将此路径转发到静态站点。
:::

邮箱验证和密码重置也使用 `host` 域名下的页面。推荐在同一公开域名上，将 `/platform/` 转发到静态站点，其他路径转发到 API 服务；此时 `BASE_URL` 保持为空即可。

## Helm 场景

若后端仍通过 Helm 部署，只需在 `values.yaml` 中设置：

```yaml
cbctf:
  host: "https://api.ctf.example.com"
  gin:
    origins:
      - "https://ctf.example.com"
```
