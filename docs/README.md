# CBCTF Docs

CBCTF 平台部署、运维与功能使用文档，基于 Rspress 构建。

## 开发

```bash
pnpm install
pnpm dev
```

## 构建与预览

```bash
pnpm build
pnpm preview
```

构建产物输出到 `doc_build/`。该目录是生成结果，不应提交到源码仓库。

## 部署路径

默认站点根路径为 `/`。如果部署到 GitHub Project Pages 等子路径，请在构建时设置 `DOCS_BASE`：

```bash
DOCS_BASE=/CBCTF/ pnpm build
```

`docs/public/` 中的资源用于稳定 URL，例如 `/logo.png` 和 `/img/homepage.png`。修改 `base` 后，需要在预览环境确认这些公开资源路径仍能访问。

## 文档规范

- 内容统一放在 `docs/` 下，按「开始」「部署」「功能」「管理员」组织。
- 使用 `_nav.json` 和 `_meta.json` 控制导航与侧边栏顺序，不通过文件名前缀排序。
- 每个页面需要包含 `title` 和 `description` frontmatter，用于 SEO、站内搜索和 `llms.txt`。
- 以任务为导向编写标题和示例，保持代码块最小、可复制、版本准确。
- 新增 MDX 组件前优先使用 Rspress 内置组件；全局样式放在 `styles/global.css`。

## 搜索与 AI 输出

构建时会生成本地搜索索引，并启用 Rspress `llms` 输出：

- `doc_build/llms.txt`
- `doc_build/llms-full.txt`
- 与每个页面对应的 Markdown 文件

这些文件依赖页面 frontmatter 的描述质量。新增页面时请同步补充准确的中文 `description`。

## 调试

```bash
pnpm inspect
```

当配置解析或插件行为异常时，可使用：

```bash
DEBUG=rsbuild pnpm build
```

如果页面在开发环境正常但构建失败，优先检查 SSG 兼容性、MDX frontmatter、公开资源路径和 `doc_build/.rsbuild` 中的最终配置。
