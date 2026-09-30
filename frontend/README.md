# CBCTF Frontend

CBCTF 平台的 React 前端，构建产物由 Go 二进制内嵌。保持暗色主题、Maple 字体、geek-blue 强调色和中英文支持。

## 开发与验证

使用 Node.js 24 和 pnpm；依赖版本以 `package.json`、锁文件为准。

在 `frontend/` 执行：

```sh
pnpm install --frozen-lockfile
pnpm dev
pnpm test
pnpm lint:check
pnpm i18n:check
pnpm build
pnpm build:report
```

- `pnpm test` 使用 Node 内置测试运行器，无额外测试运行时依赖。
- `pnpm lint:check` 不修改文件，包含 Hooks 调用合法性检查。
- `pnpm i18n:check` 检查双语键和插值参数一致性、缺失引用、未引用候选以及冗余翻译别名；`pnpm i18n:audit` 输出详细引用、动态范围和双语重复文案报告。
- `pnpm lint` 会自动修改文件；局部修改优先只格式化涉及的文件。
- `pnpm build` 生成 `dist/`，Go 构建前必须执行；不要提交构建产物。
- `pnpm build:report` 基于构建清单统计入口静态 JS/CSS 的原始、gzip 和实际预压缩体积，检查管理布局、登录页、图表与编辑器未泄漏到公共入口。字体、路由懒加载和 API 流量不包含在此口径内。
- `pnpm preview` 仅预览静态产物，不是后端服务。接口配置在 `src/api/config.js`，Vite 没有开发代理。
- 生产静态路径为 `/platform/`，客户端使用 HashRouter；保持现有路由与权限配置。

## 目录与依赖

```text
src/
  api/                         HTTP 请求及平台接口适配
  components/
    common/                    与业务无关的基础 UI 和展示适配器
    features/
      Admin/
        challenges/editor/     Compose 模型、解析、生成、校验和题目表单
        challenges/testing/    测试会话、实例与附件面板
        traffic/               流量数据会话、回放、布局及 SVG 画布
        victims/               靶机列表、批量启动及共用 Pod 日志
        generators/            全局与比赛共用的生成器管理
        images/                两种作用域共用的镜像拉取
        tasks/                 历史任务与实时队列
        rbac/                  用户、角色、权限、组与成员管理
        webhook/ smtp/ oauth/  各自的编辑、历史与提交模型
        cheats/                审核、证据与批量操作
        details/               管理员用户/队伍详情会话
        network/               IP 查询和日志中的 IP 交互
        Contests/              仅比赛管理使用的 editor/teams/challenges
        SystemConfig/          系统配置字段、分区和业务 hook
      CTFGame/
        Challenges/
          models/              题目/实例数据转换
          hooks/               列表、选题会话与比赛概况
          display/             实例、Flag 表单、提示、骨架
        Team/                  队伍设置、加载和编辑
        Writeup/               两个入口共用的题解上传与查询
      Scoreboard/              两端共用的排名、题目矩阵、时间线和纯模型
      UserSettings/            资料、安全和账户确认分区
      layouts/                 选手/管理员页面壳及导航
  hooks/                       不依赖管理员业务的通用 hooks
  pages/                       路由入口、权限/作用域适配、业务组件组合
  routes/                      懒加载入口及权限守卫
  store/                       应用级 Redux 状态
  utils/ config/ lib/           通用函数、配置和外部库接入
tests/                         领域逻辑、异步会话与结构约束测试
```

依赖方向是 `pages -> features -> common`。页面和业务会话可以使用 `api`，纯模型不可导入 React、API、store、toast、ECharts 或 Monaco。组件不能反向引用页面，公共 UI 和全局 hooks 不能引用管理端业务。业务 hook 与所属功能就近存放，不加入一个混合全站业务的大 barrel。

目录大小写必须与磁盘一致。保留既有 `Admin`、`CTFGame`、`Contests` 等边界名称；新增管理员业务域使用小写目录、组件使用 PascalCase、hooks 使用 `useXxx`。不要同时建立大小写不同但含义相同的目录。

## 组件约定

- `Modal` 只负责弹窗外壳、Portal、滚动锁和顶层焦点管理。确认操作用 `ConfirmModal`，业务弹窗组合这两个基础组件，不再手绘遮罩。
- `ModalFooter` 组合取消/提交动作，`CRUDModalFooter` 选择明确 CRUD 文案；进行中状态使用 `submitLoading`。普通按钮直接用 `Button`，不要恢复 `ModalButton` 别名。
- 标准字段使用 `Input`、`Textarea`、`Select`、`DateTimeInput` 的 `label`。自定义单控件用 `FormField htmlFor` 配合明确 id；组合字段用 `fieldset/legend`。
- 布尔勾选使用 `Checkbox`，`onChange` 接收 DOM event；不要恢复名称与行为不符的 `FormSwitch`。
- `List` 是数据表格；空态使用 `emptyContent`，默认从数据推导。`minWidth` 限定内部表格宽度，窄屏在表格区域滚动，不能靠裁切隐藏操作。
- 列宽必须有边界：排行榜得分列按文字测量但最多 240px（含内边距），溢出以省略号显示。`TruncatedText` 用于表头和单行文本，保留完整 DOM 文本与 `title`，不截断原始数据；多行摘要使用行数限制并提供完整提示。自定义数据表使用固定列布局及明确列宽，避免长名称、ID、地址或数值参与无限扩列；交互控件不能用文本截断替代。
- 管理列表通过 `List.getRowActions(item)` 返回 `{ key, label, icon, onClick, disabled, danger, hidden, inline }`，所有操作均可在右键菜单执行。编辑、删除、查看详情不放入操作列；其他操作标记 `inline: true`，并通过 `actionsColumn: { label, width }` 展示操作列，共用同一份回调和状态。没有行内操作时不设置 `actionsColumn`，避免空操作列。有详情功能的列表，`onRowClick` 优先查看详情，编辑只通过右键菜单进入；没有详情功能的可编辑列表，行点击直接进入编辑。查看详情通过左键点击行或右键菜单进入。行内按钮、链接、表单控件及 `data-row-interactive` 区域不触发行点击。
- 自定义表格使用 `common/useContextMenu` 的 `getRowProps`、`hint` 和 `menu`，共用 `ContextMenu`，操作列使用 `common/RowActions`；支持右键，数据刷新后关闭旧菜单。操作标签必须使用中英文翻译，权限和状态限制由业务组件提供。
- `Pagination.current/total/onChange` 中 `total` 表示总页数，页码和数据请求由业务所有者维护。
- `AutoRefreshControl` 仅展示间隔选择，值和回调单位均为秒。它不拥有计时器，不决定页面是否轮询。
- `Spinner` 是装饰图元，`Loading` 是有状态文本的等待区，`Skeleton` 保留布局，`EmptyState` 表示没有数据；不要合成参数繁多的万能状态组件。
- Markdown 统一直接导入 `common/MarkdownContent`，集中 GFM 与 prose 排版，不启用原始 HTML。纯文本公告仍保持纯文本。
- 图表直接懒加载 `common/EChart`，它负责按需注册和 ref 透传；不要从基础 common barrel 导出重型适配器。
- 日志只经过 `AnsiLog` 的 ANSI 转换、业务后处理、DOMPurify 后渲染。不要绕过消毒流程。
- Toast 的方法和参数以 `utils/toast.js` 为准，例如 `toast.danger({ description })`；没有 `toast.error`，不能传字符串代替配置对象。

## 状态边界

- 一个业务状态只有一个所有者；不要把整份页面 state/setters 原样传进多层子组件。
- 按自然职责拆分，不把超大页面整体搬成超大 hook，不为少量重复制造通用 CRUD/schema 框架。
- 全局/比赛生成器和靶机通过页面提供明确 API scope；比赛切换使用稳定身份边界，旧响应不能写回新会话。
- 轮询必须处理在途请求。题目状态在请求结束后调度，管理列表在相同 query 请求未结束时跳过 tick；分页、筛选和手动刷新仍能立即替换旧查询。
- 管理列表及仪表盘通过 `utils/visiblePolling` 暂停隐藏标签页和离线状态下的自动刷新，恢复可见且在线后刷新一次；初次加载和手动刷新仍由业务会话控制。
- 关闭弹窗、切题、切比赛、卸载都会失效旧请求和计时器。不能用单一 mounted 布尔值代替选题代次、请求序号等不同约束。
- 角色授权、组成员变更、单 Flag 保存、封面上传是独立即时事务，不应混入主表单保存。
- 排行榜上分曲线是排名/表格两种视图上方的共用区域，进入比赛排行榜即独立加载，视图切换和分页不取消或重复发起曲线请求。时间线 `null` 是未请求，`[]` 是已取得空结果；失败可单独重试。选择器与曲线共用按数值排名稳定排序的模型，ECharts 隐藏曲线保留实例与 `replaceMerge`。

## 自动化保护

### 加载与资源交付

- 管理端与比赛布局按路由懒加载，先通过权限守卫再加载管理端资源。
- 翻译仍在 `src/i18n/locales/{en,zh-CN}.json` 维护。`build/localeChunks.js` 在构建和开发时分离 `admin` 文案，由管理端布局加载；公共页面使用公共翻译键。两种语言的管理端资源在布局渲染前一起注册，语言切换无需额外请求。
- 通用操作、启用状态、题型、公告类型和排行榜表头直接复用 `common`，不再用只含 `$t(...)` 的键转发。全局/比赛生成器共用 `admin.generators`；靶机筛选、列表、状态和日志共用 `admin.victims`，比赛批量启动文案留在 `admin.contests.containers`。
- 翻译审计使用 ESLint AST，忽略注释并识别配置中的键、条件表达式、模板字符串和 `Trans.i18nKey`。动态状态、后端错误码和传入的翻译前缀按匹配范围保守保留，报告中的“未引用候选为零”不代表所有动态分支都已执行。相同文案仍需按业务语义决定是否合并，审计不会自动删除或合并。
- 完整句子使用插值；带样式的句子使用 `Trans` 命名组件，不按中英文语序拆成多个片段键。
- `pnpm build` 为可压缩静态文件生成 `.gz`，Go 静态处理器按 `Accept-Encoding` 提供预压缩内容，无需逐请求压缩。保留未压缩文件以支持客户端协商和 Range 请求；入口重新验证、带哈希资源长期缓存的策略不变。
- `go test ./frontend`（仓库根目录）验证 gzip 协商、HEAD、Range、缺失文件和入口重定向。构建产物须先生成。
- Markdown 内容按 props 缓存渲染，图片和列表缩略图使用浏览器懒加载与异步解码。

`tests/architecture.test.js` 检查：

- 相对 import/re-export/dynamic import 的存在性和精确大小写。
- 同步模块依赖无环，包括 barrel 转导出。
- 公共层、全局 hooks、页面与业务组件的依赖边界。
- 每个 `src` JS/JSX 文件不超过 500 行；这是防止重新堆积的上限，不是建议写满的目标。
- 已提取纯模型的封闭依赖集合；新增纯模块需要显式纳入检查。

领域测试覆盖 Compose parser/serializer/validation、创建/更新 payload、网络布局、时间线、分页/选择、配置差量提交、IP 展示、测试会话以及慢请求轮询和日志错误处理。异步 hook 测试仅替换 React 生命周期和网络边界，通知测试使用真实 toast service 合同。

修改后至少运行 `pnpm test`、`pnpm lint:check`、`pnpm build`。涉及交互时，还需回归中英文、手机/桌面、嵌套弹窗、错误重试、跨 scope 的迟到响应，以及首页不加载 Monaco/ECharts。浏览器模拟接口不能代替真实后端验收。

## 本次改造边界

保留路由、权限、接口作用域、Compose 扩展及现有视觉，没有引入额外生产依赖或旧路径转导出壳。以下旧行为未作为结构改造擅自修改：

- 队伍设置的队长选择尚未写入 `captain_id`，不能将其视为已实现的转让功能。
- 队伍编辑取消后再次打开仍保留草稿。
- 部分网络异常仍显示 Axios 原始英文信息。

这些需要单独确定业务语义并补验收，不应混入纯文件提取后宣称已解决。
