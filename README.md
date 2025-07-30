# CBCTF

CBCTF 是一个基于 Go 的 CTF 平台, 旨在满足常规的办赛需求

## 特性

- 可创建多种题目
  - 题目类型: 静态题目, 动态题目, 动态容器
    - 静态题目: 队伍间共用附件
    - 动态题目: 实时生成唯一附件
    - 动态容器: 自动生成并启动队伍间隔的容器环境

  - flag类型: 静态flag, 动态flag, UUID
    - 静态flag: 每次生成的flag均相等, 配置时, 格式为: static{}, `{}`中的内容固定
    - 动态flag: 每次生成的flag由基础模板变化生成, 配置时, 格式为: dynamic{}, `{}`中的内容将随机变化, 保持可读性但随机
    - UUID: 每次生成的flag为UUID格式, 配置时, 格式为: uuid{}, `{}`中的内容将被忽略
    
    动态容器支持将flag注入至环境变量或作为文件挂载至指定路径

  - 分值类型: 静态分数, 线性分数, 非线性分数
    - 静态分数: 每个flag的分数不随接触人数而变化
    - 线性分数: 随着接触人数的增加, 等量减少分值
    - 非线性分数: `(MinScore - InitScore) / (Decay ^ 2) * (Solvers^2) + InitScore`
    
    均可设置三血奖励: 一二三血奖励依次为初始分数的 5% 3% 1%

- 基于`Kubernetes`的动态附件生成, 动态容器分发
- 基于`Kube-OVN`的自定义VPC网网络
- SMTP 邮件验证功能
- 平台内 Writeup 收集, 下载
- 记录比赛期间动作事件日志
- 基于 `Redis` 的数据缓存, `MySQL` 的数据存储
- 基于 `NFS` 的文件存储
- 支持镜像预热, 容器预热

## 动态附件

- 采用Docker容器进行生成, 生成过程与主机环境隔离
- 可基于通用Docker镜像, 通过上传python脚本进行生成, 减轻出题压力
- [示例](example/dynamic/README.md)

## 动态容器

- 后端是如何区分题目网络环境为`Pod`还是`VPC`
  - 当docker-compose中配置了`networks`字段时, 为`VPC`; 此时, 所有容器须配置`networks`, 手动指定IP地址

### Pod

- [示例](example/pods/pod/docker-compose.yml)

### VPC

- [示例](example/pods/vpc/docker-compose.yml)

## 环境依赖
- [Kube-OVN](https://kubeovn.github.io/docs/stable/start/prepare/)
- [Multus](https://github.com/k8snetworkplumbingwg/multus-cni/tree/master)
  - 推荐使用`Thin Plugin`, 无需手动更改
  - 默认情况下`Thick Plugin`极易发生[OOMKilled](https://github.com/k8snetworkplumbingwg/multus-cni/issues/1346)以及[Text file busy](https://github.com/k8snetworkplumbingwg/multus-cni/issues/1221), 可参考[issue #1346](https://github.com/k8snetworkplumbingwg/multus-cni/issues/1346#issuecomment-2644110944), [pr #1213](https://github.com/k8snetworkplumbingwg/multus-cni/pull/1213)手动修正
