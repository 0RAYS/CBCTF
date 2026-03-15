---
sidebar_position: 6
---

# 作弊检测

CBCTF 内置多维度自动作弊检测，管理员可查看、确认或驳回检测结果。

## 检测机制

### 1. `token_magic` — 设备指纹不匹配

同一账号的 JWT 中绑定的设备指纹（FingerprintJS 生成）与当前请求的设备指纹不一致，说明该账号在不同设备上使用，可能存在账号共享行为。

### 2. `same_device` — 多账号共用设备

多个不同账号（跨队伍）使用相同的设备指纹，说明可能是同一台设备用多个账号参赛。

### 3. `same_web_ip` — 多队伍共用 Web IP

多支队伍的 Web 请求（登录、提交等）来自同一公网 IP，可能存在同场地协作或账号共享。

### 4. `same_victim_ip` — 多队伍共用靶机访问 IP

多支队伍访问靶机的 IP 相同（tcpdump 捕获），强烈暗示存在跨队协作。

### 5. `wrong_flag` — 跨队 flag 提交

队伍提交了属于另一支队伍的动态 flag，表明存在 flag 共享行为，记录为 `wrong_flag` 类型。

## 作弊记录字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `type` | string | 检测类型（见上） |
| `reason` | string | 详细原因描述 |
| `reason_type` | string | 触发类型（`auto` 自动 / `manual` 手动） |
| `magic` | string | 涉及的设备指纹（仅设备类检测） |
| `ip` | string | 涉及的 IP 地址 |
| `time` | datetime | 检测时间 |
| `checked` | bool | 管理员是否已审核 |
| `comment` | string | 管理员备注 |

## 作弊状态

| 状态 | 说明 |
|------|------|
| `suspicious` | 系统自动检测，待人工确认 |
| `cheater` | 确认为作弊 |
| `pass` | 确认为误报 |

## 管理员处理流程

1. 访问比赛的作弊记录列表（`admin:cheat:list`）
2. 查看各条记录的 `type`、`reason`、相关队伍信息
3. 结合 IP 历史（`admin:ip:search`）和流量捕获进行综合判断
4. 更新 `type` 为 `cheater` 或 `pass`（`admin:cheat:update`）
5. 在 `comment` 中填写处置说明
6. 标记 `checked: true`

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

在 `config.yaml` 中配置 IP 白名单，白名单内的 IP 不触发 IP 类作弊检测：

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

| 场景 | 可能误报的类型 |
|------|--------------|
| 大学/公司内网（NAT） | `same_web_ip`、`same_victim_ip` |
| 商业 VPN 服务 | `same_web_ip` |
| 参赛者使用多个浏览器或无痕模式 | `token_magic` |
| 比赛现场同一 WiFi | `same_web_ip` |

建议将竞赛现场的出口 IP 加入白名单，或在 Helm values 中配置 `cbctf.cheat.ip.whitelist`。
