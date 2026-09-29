---
title: 后端上下文与 JSON 数据模型
description: CBCTF 后端的请求取消、后台任务生命周期与 GORM JSONB 序列化约定。
---

# 后端上下文与 JSON 数据模型

## Gin Context 与标准 Context

`*gin.Context` 管理 HTTP 参数、认证信息、响应和中间件状态。它在请求结束后可能被 Gin 复用，不应被业务层、并发工作函数或后台任务长期持有，也不能并发写响应。

`ctx.Request.Context()` 是标准 `context.Context`，可并发读取，随客户端断开、请求取消或父级超时而结束。Gin 的 `ContextWithFallback` 决定 `*gin.Context` 是否转发标准 Context 的取消/截止时间；业务代码不应依赖这个全局开关。

项目采用如下边界：

| 场景 | 传递方式 |
| --- | --- |
| 绑定参数、读写 Cookie、取中间件模型、输出 JSON | 当前请求的 `*gin.Context` |
| Kubernetes、HTTP 客户端、耗时分析、请求内 Redis 操作 | `ctx.Request.Context()`，必要时派生更短超时 |
| 请求内 GORM 操作 | `db.DB.WithContext(ctx.Request.Context())`；事务沿用传入连接的 Context |
| Asynq worker | 使用 worker 收到的 Context，同时传给数据库和外部调用 |
| 请求后审计、解锁、失败清理 | 有明确时限的独立 Context；需要保留追踪值时使用 `context.WithoutCancel(parent)` 再加超时 |
| 启动初始化、长期控制器、定时任务 | 独立生命周期；不能继承某个 HTTP 请求 |

不要在 goroutine 中给外部共享的 `ctx` / `cancel` 重新赋值。为 errgroup 获取一个派生上下文并只读传递；`Wait()` 后这个派生上下文已经取消，后续操作应使用仍有效的父级上下文。

批量启停接口在请求中完成数据库登记及任务入队，实际 Kubernetes 工作由 Asynq 执行，不再通过裸 goroutine 提前返回成功。flag 提交/重置后的清理只同步执行短时入队操作，避免将请求事务交给生命周期不明的 goroutine。

## JSONB 模型字段

纯 JSON 字段使用 GORM 内置 serializer，例如：

```go
Rules    []string          `gorm:"serializer:json;type:jsonb;default:'[]'"`
Headers  map[string]string `gorm:"serializer:json;type:jsonb;default:'{}'"`
Spec     VictimSpec        `gorm:"serializer:json;type:jsonb;default:'{}'"`
Payload  any               `gorm:"serializer:json;type:jsonb"`
```

已经迁移的类型包括集合、赛事奖项/时间线、多语言品牌内容、flag 绑定、题目模板、Pod/靶机规格、运行资源、作弊引用与任务负载。嵌套结构由外层 JSON 字段一次序列化，无需为每个内层结构实现 `Value()` / `Scan()`。

无独立行为的集合包装类型已移除，模型、DTO 和服务直接使用 `[]string`、`map[string]string`、`map[string]uint`、`map[string][]uint`、`[]Prize`、`[]Timeline` 和 `[]XVolume`。不保留旧类型别名或转接层。

空值约定：

- 创建时省略字段，使用数据库 tag 中的 `[]` / `{}` 默认值。
- 显式的空切片、空 map 分别保存为 JSON `[]`、`{}`，不再统一折叠为 SQL NULL。
- 显式以 `nil` 更新可空字段时保存为 SQL NULL；读取 NULL 会清空原对象中的旧值。
- 列表/对象默认值只影响 INSERT 时省略的字段，不代表 UPDATE 会自动替换 nil。

### Map 更新

GORM 的结构体写入会使用字段 serializer，但 `Updates(map[string]any)` 默认绑定的是原始值，切片可能展开成 SQL 元组。所有生产数据库连接池统一安装 `db.JSONUpdates` 插件，让 map 更新沿用同一字段声明的 serializer。

插件保留显式 SQL 表达式、GORM 已包装的值以及 `Select` / `Omit` 语义，不修改调用者的 map，避免乐观锁重试时重复编码。单独创建测试连接池并测试 map 更新时，也应执行 `database.Use(db.JSONUpdates{})`。

### 保留自定义转换的类型

`FileURL` / `FilePath` 带有路径转换；`SettingValue` 使用 `json.Decoder.UseNumber()` 并保存包装对象的实际值，避免系统设置数字精度丢失；网络策略、网络附件、暴露端口、VPC 和端点类型包含配置清理逻辑。这些并非单纯 JSON 编解码，仍保留其必要的转换方法。

## 验证

```bash
go test ./...
go vet ./...
```

Context 回归测试覆盖关闭 Gin 转发时，请求取消仍可传到 Redis/数据库，以及 OAuth HTTP 请求取消和 errgroup 取消隔离。JSON 测试覆盖结构体/map 写入一致性、空集合、NULL、嵌套规格、无效 JSON、表达式及零值设置。

真实 PostgreSQL 往返测试需要 `CBCTF_TEST_POSTGRES_DSN`；未设置时跳过。支持 CGO 和 C 编译器的环境还可运行 `go test -race ./...` 检查动态竞争。
