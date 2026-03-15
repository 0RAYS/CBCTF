---
sidebar_position: 6
---

# 前后端分离

CBCTF 默认将前端静态资源嵌入 Go 二进制，通过同一个服务地址提供 `/platform`。如果需要把前端单独托管到 CDN、静态站点或其他域名，可以采用前后端分离方式。

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

```javascript
export const API_CONFIG = {
  BASE_URL: 'https://api.ctf.example.com',
};
```

然后构建：

```bash
pnpm build
```

构建完成后，将 `frontend/dist/` 部署到任意静态托管环境。

## 后端配置

`config.yaml` 示例：

```yaml
host: https://api.ctf.example.com

gin:
  cors:
    - https://ctf.example.com
```

说明：

- `host` 必须填写后端真实对外地址，OAuth 回调与邮件链接都会使用它
- `gin.cors` 需要包含前端独立域名

## OAuth 注意事项

前后端分离时，OAuth 回调链路为：

1. 第三方回调到后端 `https://api.ctf.example.com/oauth/{uri}/callback`
2. 后端完成登录后重定向到 `https://api.ctf.example.com/platform/#/oauth/callback?...`

也就是说，当前代码默认仍依赖后端提供 `/platform` 下的前端回调页。若完全拆离前端托管位置，需要同步调整 OAuth 回调后的前端跳转逻辑。

## Helm 场景

若后端仍通过 Helm 部署，只需在 `values.yaml` 中设置：

```yaml
cbctf:
  host: "https://api.ctf.example.com"
  gin:
    cors:
      - "https://ctf.example.com"
```
