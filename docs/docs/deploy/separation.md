---
sidebar_position: 6
---

# 前后端分离

CBCTF 默认将前端静态资源嵌入 Go 二进制，通过同一端口提供服务。如需将前端托管至 CDN、GitHub Pages 或其他静态托管服务，可采用前后端分离部署方式。

## 使用场景

- 需要通过 CDN 加速前端资源
- 前端部署在 GitHub Pages 等静态托管服务
- 前后端使用不同域名

## 编译前端

```bash
git clone https://github.com/0RAYS/CBCTF.git
cd CBCTF/frontend
pnpm install
```

编辑 `src/api/config.js`，将 `BASE_URL` 修改为后端 API 地址（**不含末尾 `/`**）：

```javascript
export const API_CONFIG = {
  BASE_URL: 'https://api.ctf.example.com',   // 修改为后端地址
  TIMEOUT: 10000,
  // ...
};
```

然后构建：

```bash
pnpm build
```

构建产物位于 `dist/` 目录，将其部署至任意静态托管服务。

## 后端配置

编辑 `config.yaml`：

```yaml
host: https://api.ctf.example.com    # 后端 API 地址

gin:
  cors:
    - https://ctf.example.com        # 前端部署地址（允许跨域）
```

:::info
分离部署时：
- `host` 填写**后端** API 服务域名，用于 OAuth 回调等
- `gin.cors` 填写**前端**部署域名，允许跨域请求
:::

## 使用 Helm 时的分离部署

在 `values.yaml` 中配置：

```yaml
cbctf:
  host: "https://api.ctf.example.com"
  gin:
    cors:
      - "https://ctf.example.com"
```
