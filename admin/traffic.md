> For AI agents: the complete documentation index is available at /llms.txt, the full documentation bundle is available at /llms-full.txt.

# 容器流量分析

启用题目抓包后，平台通过 capture sidecar 保存原始 PCAP 与进程信息。后端分析代码位于 `internal/traffic`，支持 PCAP、PCAPNG 和混合链路类型的 PCAPNG；离线分析不需要 libpcap 或 CGO。

管理员流量回放窗口提供拓扑、抓包下载及分析证据面板。可查看 flag、攻击线索、HTTP/DNS、会话、客户端 IP 和跨队伍重叠，按当前回放窗口或选中节点/连线筛选，并跳转到证据对应时刻。分析结果分页显示，支持中英文；请求失败显示错误，不使用样例流量替代证据。报告也可通过下列管理员 API 查询。

## 接口与权限

| 接口                                                                                | 权限                           | 用途                    |
| --------------------------------------------------------------------------------- | ---------------------------- | --------------------- |
| `GET /admin/victims/:victimID/traffic?time_shift=0&duration=1000`                 | `admin:traffic:read`         | 毫秒级时间窗口拓扑             |
| `GET /admin/victims/:victimID/traffic/analysis`                                   | `admin:traffic:read`         | flag、HTTP、DNS、会话和攻击线索 |
| `GET /admin/contests/:contestID/teams/:teamID/victims/:victimID/traffic/analysis` | `admin:contest_traffic:read` | 比赛队伍靶机报告              |
| `GET /admin/contests/:contestID/traffic/overlaps`                                 | `admin:contest_traffic:read` | 同场比赛跨队伍公共访问 IP        |

比赛队伍路径同样支持 `traffic?time_shift=...&duration=...` 和 `traffic/download`。响应使用平台标准 JSON 封装。

## 跨队伍访问 IP

只有可归因于客户端的访问证据参与比较：

- 已知靶机地址收到的 TCP 初始 SYN（排除 SYN-ACK）。
- 专用 `frpc.pcap` 中 FRPC 发往 nginx 本地监听端口 `10000+` 的 PROXY v1/v2 头；支持 TCP 分段和 IPv6 客户端。

FRPC 的本地转发地址只用于解析传输封装，不作为选手地址。普通应用负载中的 PROXY 头不会作为可信客户端证据。UDP 或仅抓到连接中段而没有可信代理头的流量，不推断其真实访问者。

比较使用包中的访问时间，限定在比赛的 `[开始时间, 结束时间)` 内，并要求至少两支不同队伍。报告返回 IP、队伍 ID、首次访问时间和靶机 ID。访问同一公共 DNS、下载源等**目标地址**不会产生跨队伍访问重叠。

IPv4-mapped IPv6 统一为 IPv4。私网、回环、未指定、链路本地、组播、共享地址空间、文档/测试网段和特殊 IPv6 地址不参与比较；同时应用 `cheat.ip.whitelist` 的 IP/CIDR 白名单。结果标记为 `suspicious`，不自动判定作弊或封禁队伍。公共 NAT/VPN 仍需要管理员结合其他证据判断。

查询靶机分析接口会更新实时访问索引，15 秒内复用最近的报告；停止后的归档任务保存最终报告。重叠查询和现有 `same_victim_ip` 检测使用已生成的访问索引，不会主动遍历所有运行中的靶机文件。

## 内网拓扑回放

- 内网节点来自靶机保存的 Pod 网络与实际端点，保留题目使用的 IPv4 私网和 IPv6 ULA 地址；公网 FRP 暴露地址不当作内网节点。
- 不根据流量排名猜测靶机地址。元数据不足时，未识别的节点显示为外部节点。
- 每个窗口为 `[time_shift, time_shift + duration)`，在抓包末尾截断。`started_at` 是绝对 UTC 起点，窗口与时间轴字段使用毫秒。
- 已配置的内网节点在空闲窗口仍保留。流量分为 `ingress`、`egress`、`internal`、`external`。
- `packets` 为包数；`connections` / `total_connections` 为双向五元组数，同一五元组复用不视为独立会话。
- 不重复分析 `.enrich.pcap`。不同 Pod 在 1ms 内观测到的字节完全相同的网络包去重，同一采集点的重传保留；经路由改写或时间差较大的重复观测不会合并。
- 时间轴至多约 2000 个桶，`timeline_bucket_ms` 表示实际桶宽。Redis 使用原子快照，运行中实例缓存 15 秒，已停止实例缓存 30 分钟。

## 内容与 flag

`analysis` 返回 `report`、`accesses` 与 `archived`。其中：

| 报告字段                     | 内容                              |
| ------------------------ | ------------------------------- |
| `flags`                  | 候选值、编码路径、`verified` 与证据         |
| `http`                   | 方法、Host、URI、响应状态和 Content-Type  |
| `dns`                    | 查询域名、类型及请求/响应标记                 |
| `sessions`               | 按采集文件与五元组归并的包数、字节数、SYN/RST 数和方向 |
| `indicators`             | 规则名、线索级别、命中摘要与证据                |
| `warnings` / `truncated` | 抓包不完整、分析能力限制、资源截断情况             |

TCP 内容支持跨包拼接、乱序及重传处理；遇到序列空洞分段分析，不跨缺失字节拼出 flag。HTTP 支持 chunked、gzip 和 deflate；内容支持 URL、Base64/Base64URL、Hex、HTML 实体、JSON 字符串、gzip 和 ZIP。DNS TXT 和未加密传输负载也会参与检测。

常见 `prefix{...}` 格式只是候选结果。`verified: true` 表示与该队伍生成的 flag 或非比赛题目的静态 flag 精确匹配，也支持无花括号的已知 flag；它不表示选手已成功提交或存在作弊。

每条证据包含采集文件、源/目标 IP 和端口、协议与时间范围。TCP 内容的时间范围覆盖参与重组的方向流，不是精确的单字节到达时间。报告的 `metric_scope` 为 `capture_observations`：保留不同采集点证据，因此其包数/字节数可能高于去重后的拓扑统计。

## 攻击线索与边界

线索包括 SQL 注入、目录穿越、命令执行和模板注入的常见特征；同一源/目标在一个 UTC 分钟桶内探测至少 20 个不同端口；重复 SYN、TCP RST、较长 DNS 查询和超过 1MiB 的内网出站传输。规则用于定位值得查看的流量，不直接证明攻击成功。

每方向 TCP 流最多保留 1MiB，每文件总重组缓冲最多 32MiB，最多 4096 条方向流；各类发现/会话最多 2000 条。解码最多两层、最多 128 个待处理候选；单次解码扫描预算 4MiB。ZIP 最多读取 32 个文件，累计解压最多 1MiB，不将文件释放到磁盘。达到限制会标记截断，不能将空结果解释为没有攻击或 flag。

当前不解密 TLS/SSH，不重组 IP 分片，不跨不同采集文件重组 TCP，不保证还原缺包、复杂 TCP 重叠冲突、HTTP/2 或加密/嵌套归档。相应可识别的限制会在报告中列出。完整 PCAP 仍可下载后用于进一步取证。

## 后端验证

```bash
go test ./internal/traffic ./internal/service ./internal/db ./internal/router ./internal/utils
```

PostgreSQL 证据查询、JSONB 序列化与实时/归档写入隔离测试使用 `CBCTF_TEST_POSTGRES_DSN`。未设置时跳过该集成测试；测试通过事务内临时表验证，不创建生产表记录。
