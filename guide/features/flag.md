***

title: Flag 规则
description: 说明 CBCTF 的 static、leet 和 uuid Flag 类型及多队伍动态 Flag 行为。
-----------------------------------------------------------------

# Flag 规则

CBCTF 支持三种 flag 类型，每道题可配置多个 flag，每个 flag 独立计分。

## Flag 类型

| 类型   | 配置格式              | 生成结果（前缀 `CBCTF`）                              | 特点              |
| ---- | ----------------- | --------------------------------------------- | --------------- |
| 静态   | `static{content}` | `CBCTF{content}`                              | 所有队伍相同          |
| 动态   | `leet{template}`  | `CBCTF{ThiliS-iS-4_Dyn4MIC_FLaG}`             | 每队不同，**长度可变**   |
| UUID | `uuid{}`          | `CBCTF{550e8400-e29b-41d4-a716-446655440000}` | 每队不同，标准 UUID 格式 |

:::info
flag 的实际前缀由比赛配置中的 `prefix` 字段决定，而非固定为 `CBCTF`。
:::

## 静态 Flag

```text
配置: static{this_is_a_static_flag}
生成: CBCTF{this_is_a_static_flag}
```

所有队伍的 flag 内容相同。适用于所有题目类型。

## 动态 Flag

```text
配置: leet{this_is_a_leet_flag}
生成: CBCTF{ThiliS-iS-4_Dyn4MIC_FLaG}
```

基于模板随机替换字符，保持可读性。**生成结果的长度可能与模板长度不同**，题目设计中不得依赖 flag 长度。

## UUID Flag

```text
配置: uuid{}
生成: CBCTF{550e8400-e29b-41d4-a716-446655440000}
```

标准 UUID v4 格式，长度固定，每队不同。

## Flag 注入方式


**动态附件题**

### 动态附件题
平台自动调用生成器容器中的 `/root/run.sh` 脚本，将 flag 作为参数传入：
```bash
/root/run.sh {team_id} {base64(base64(flag1),base64(flag2),...)}
```
示例（team\_id=1，两个 flag）：
```bash
/root/run.sh 1 UTBKRFZFWjdabXhoWnpGOSxRMEpEVkVaN1pteGhaejo5
```
其中第二个参数为 `base64(base64("CBCTF{flag1}") + "," + base64("CBCTF{flag2}"))` 的结果。


**容器题：环境变量**

### 容器题：环境变量注入
环境变量名必须以 `FLAG_` 为前缀，平台在启动容器时替换 flag 值：
```yaml
services:
  web:
    image: nginx:latest
    environment:
      - FLAG_1=leet{this_is_a_dynamic_flag}
      - FLAG_2=static{this_is_a_static_flag}
      - FLAG_3=uuid{}
    ports:
      - "80:80"
```


**容器题：文件**

### 容器题：文件注入
通过 service 下的 `x-volumes` 扩展字段配置，平台将文件内容写入容器内的指定路径（通过 Kubernetes ConfigMap +
VolumeMount 实现）。`content` 中可嵌入任意数量的 flag 模板，也可包含其他非 flag 内容：
```yaml
services:
  web:
    image: nginx:1.25
    x-kubevirt: false
    x-volumes:
      - path: /flag
        content: uuid{}
      - path: /etc/app/config.env
        content: |
          APP_MODE=prod
          SECRET_KEY=static{shared_secret}
          FLAG=leet{another_flag}
      - path: /home/ctf/readme.txt
        content: |
          Welcome to the challenge!
          The flag is hidden somewhere on this machine.
    ports:
      - "80:80"
```
每个 `x-volumes` 条目对应容器内一个文件。`content` 中出现的 `static{}`、`leet{}` 或 `uuid{}` 会被识别为 flag
模板并在启动时替换为实际 flag 值，其余内容原样写入文件。


## Flag 生成时机

队伍的 flag **仅在初始化（init）或重置（reset）时生成**，容器重启不重新生成 flag。

- 选手重启容器不会改变 flag
- 选手执行「重置题目」时会生成新 flag，旧 flag 失效
- 题目测试模式（admin）不产生真实 flag 记录

## 多 Flag 配置

每道题可配置多个 flag，每个 flag 独立计分。选手提交其中任一 flag 即可获得对应分数，提交所有 flag 可获得该题目的全部分数。

Flag 前缀由所在比赛的 `prefix` 字段配置，不同比赛可以使用不同前缀（如 `flag`、`CTF`、`CBCTF`）。
