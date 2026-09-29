---
title: 作弊检测
description: 查看和处理 CBCTF IP、跨队 Flag 和异常提交记录。
---

# 作弊检测

CBCTF 内置多维度自动作弊检测，管理员可查看、确认或驳回检测结果。

## 检测机制

### 1. `same_web_ip` — 多队伍共用 Web IP

多支队伍的 Web 请求（登录、提交等）来自同一公网 IP，可能存在同场地协作或账号共享。

### 2. `same_victim_ip` — 多队伍共用靶机访问 IP

多支队伍访问靶机的 IP 相同（capture 捕获），需要结合公共出口、代理及比赛场地判断。关闭抓包后不会获得这类实例访问流量证据。

检测比较实际访问时间内的公共客户端源 IP，排除特殊 IPv4/IPv6 地址与白名单，不比较公共 DNS 等目标地址。证据来源、限制和重叠查询接口见[容器流量分析](./traffic)。

### 3. `wrong_flag` — 跨队 flag 提交

队伍提交了属于另一支队伍的动态 flag，表明存在 flag 共享行为。

## 作弊记录字段

| 字段            | 类型       | 说明                                                                                      |
|---------------|----------|-----------------------------------------------------------------------------------------|
| `type`        | string   | 处理状态（`suspicious` / `cheater` / `pass`）                                                 |
| `reason`      | string   | 详细原因描述                                                                                  |
| `reason_type` | string   | 检测规则类型（`same_web_ip` / `same_victim_ip` / `wrong_flag`） |
| `ip`          | string   | 涉及的 IP 地址                                                                               |
| `time`        | datetime | 检测时间                                                                                    |
| `checked`     | bool     | 管理员是否已审核                                                                                |
| `comment`     | string   | 管理员备注                                                                                   |

## 作弊状态

| 状态           | 说明           |
|--------------|--------------|
| `suspicious` | 系统自动检测，待人工确认 |
| `cheater`    | 确认为作弊        |
| `pass`       | 确认为误报        |

## 管理员处理流程

1. 访问比赛的作弊记录列表（`admin:cheat:list`）
2. 查看各条记录的 `type`、`reason`、相关队伍信息
3. 结合 IP 地理查询（`admin:ip:search`）、相关队伍提交及流量捕获综合判断；IP 查询接口不提供登录历史
4. 更新 `type` 为 `cheater` 或 `pass`（`admin:cheat:update`）
5. 在 `comment` 中填写处置说明
6. 标记 `checked: true`

处理作弊记录不自动封禁队伍。需要限制参赛时，在比赛队伍管理中另行设置封禁。

## 重新运行检测

```bash
POST /admin/contests/:contestID/cheats
```

需要 `admin:cheat:create` 权限。对当前比赛所有数据重新运行全量作弊检测，适用于比赛结束后的全面审查。

## 批量删除

```bash
DELETE /admin/contests/:contestID/cheats
```

需要 `admin:cheat:delete` 权限。删除该比赛的所有作弊记录（谨慎使用）。

## IP 白名单

首次部署可在 `config.yaml` 中配置 IP 白名单；已有平台在「系统管理」中更新数据库设置。白名单内的 IP 不触发 IP 类作弊检测：

```yaml
cheat:
  ip:
    whitelist:
      - 127.0.0.1
      - ::1
      - 10.0.0.0/8
      - 192.168.0.0/16
      - 172.16.0.0/12
      - 100.64.0.0/10
```

## 误报场景说明

以下场景可能导致误报，管理员应结合实际情况判断：

| 场景              | 可能误报的类型                        |
|-----------------|--------------------------------|
| 大学/公司内网（NAT）    | `same_web_ip`、`same_victim_ip` |
| 商业 VPN 服务       | `same_web_ip`                  |
| 比赛现场同一 WiFi     | `same_web_ip`                  |

:::tip
建议将竞赛现场的出口 IP 加入白名单；Helm 初始配置对应 `cbctf.cheat.ip.whitelist`，已有数据库不会因升级 values 自动覆盖。
:::
