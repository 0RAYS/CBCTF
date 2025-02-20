# Generator
这是一个生成器模板, 构建一个生成器docker镜像, 用于生成动态题目的动态附件

在本示例中, 我们将会构建一个用于生成密码学题目的生成器镜像

## Workflow
Golang实现代码为 `/internal/k8s/generator.go`, 可以阅读代码辅助理解附件生成过程

1. 动态题目初始化, 获取当前用户所在队伍 `id (作为附件唯一标识)` , 并生成 `flag`
2. 启动生成器容器, 执行命令 `sleep infinity` 保持容器运行, 务必确保容器中存在 `sleep` 命令
3. 如果该题目在被设置 `generator` 镜像的同时, 还存在 `generator.zip` 文件, 则首先将主机中的 `generator.zip` 复制到容器中的 `/root` 目录下, 并解压
4. 附加至容器, 执行命令 `./run.sh id base64(flag)` 生成附件, 此处为相对路径, 默认工作路径为 `/root`, 生成的附件在容器中的绝对路径为 `/root/attachments/{id}.zip`
5. 将容器中生成的附件 `/root/attachments/{id}.zip` 复制到主机中, 并删除容器

## Example
1. 用户 user 所在队伍 ID 为 `1`, 生成的 flag 为 `flag{test}`
2. k8s 启动 pod, 并执行默认命令保持容器运行
    ```shell
    sleep infinity 
    ```
3. 题目设置了生成器镜像为 `docker.0rays.club/test/generator:latest`, 并且同时存在 `generator.zip` 文件
    ```shell
    cp /path/in/host/generator.zip /root/generator.zip
    unzip /root/generator.zip -d /root
    ```
4. 在容器中执行命令
    ```shell
    ./run.sh 1 ZmxhZ3R7dGVzdH0=
    ```
5. 此时生成的附件路径为 `/root/attachments/1.zip`, 将其复制到主机中
    ```shell
    cp /root/attachments/1.zip /path/in/host/1.zip
    ```
6. 删除 pod

## Dockerfile
```dockerfile
# 在此示例中, 使用自建的镜像站作为基础镜像, 保证镜像拉取的稳定性和速度
# 你也可以使用官方镜像, 但是最终构建的镜像一定要 push 到自建镜像站, 保证运行时镜像拉取的稳定性和速度
FROM docker.0rays.club/test/python:3.10-slim

# 将需要的文件复制到镜像中, 其中 run.sh 必须在 /root 目录下
# run.sh 将会被传入两个参数, 第一个参数是 id, 第二个参数是 base64(flag); 参数都当作字符串处理即可
COPY ./files/* /root/
COPY ./requirements.txt /root/requirements.txt

# 安装 unzip zip 以及其他依赖, 镜像中务必要有 unzip tee 命令
RUN apt-get update && apt-get install -y unzip zip && pip install -r /root/requirements.txt

# 设置工作目录, 固定值, 不可更改
WORKDIR /root

# 无所谓, 只是为了展示
CMD ["ls -al /root"]
```

## run.sh
```shell
#!/bin/bash

# $1: team_id
# $2: base64(flag)
flag=`echo $2 | base64 -d`

# 生成附件
python3 generator.py $1 $flag
```

## generator.zip
该文件存在的作用时, 允许出题人使用相同的基础镜像, 通过上传不同的 `generator.zip` 文件, 生成不同的题目, 典型的使用场景为密码学附件生成

一定程度上可以减少出题人编写Dockerfile等的工作量
```text
密码学附件通常只是简单的 .py 文件
这时候出题人可以使用相同的基础镜像, 例如 python:3.10-slim
通过上传不同的 generator.zip 文件, 其中包含 generator.py 和 template.py (本示例中), 生成不同题目的密码学附件
```

该文件由管理员上传, 作为生成器的可选依赖存在, 当该依赖不存在时, 务必保证生成器镜像可独立正确生成附件

## Attention
1. 容器中务必存在一下命令
   - `sleep`: 保持容器运行
   - `unzip`: 解压文件
   - `tee`: 容器与主机间文件复制
2. 务必存在以下文件
   - `/root/run.sh`: 附件生成入口
3. `run.sh` 务必接收并只能接收两个参数
   - `id`: 附件的唯一标识, 可认为是队伍 ID
   - `base64(flag)`: 生成附件的 flag, 以 base64 编码传入, 请确保生成的附件可解出该 flag
4. 附件的所有文件务必压缩并放置在 `/root/attachments/{id}.zip` 中, 平台将仅以此文件作为附件分发给选手
5. 请确保附件生成过程中不存在任何**耗时命令**, 尽可能加快附件的生成过程, 以免影响题目的部署速度
   - 以二进制附件为例, 如果可以, 使用字节替换的方式, 直接替换关键部分的字节, 实现生成不同的二进制文件, 代替源代码编译的过程
6. 请确保附件生成过程中不会产生任何错误, 以免影响附件生成时的正确性
7. 请确保附件可以接收**一定长度的flag**, 以免flag过长导致附件无法解出flag, flag由一定的算法在原模板的基础上进行变异, 以确保flag的随机性
   - 在本示例中, 因密钥长度较短, 当 `flag` 过长, 会导致 RSA 无法正确加密, 导致无法正常解出 `flag`
   - 变异算法为 `/internal/utils/flag.go` 中的 `RandFlag` 函数, 会在原flag模板的基础上进行字符替换, 字符重复等操作, 导致 flag 不定长
8. 构建镜像时, 因节点镜像拉取策略为 `IfNotPresent`, **严禁使用 `latest` 标签**, 请使用具体的版本号, 以确保镜像在 k8s 节点上可以正确被更新
    ```bash
    docker build -t generator-test:20250221 .
    docker tag generator-test:20250221 docker.0rays.club/test/generator:20250221
    docker push docker.0rays.club/test/generator:20250221
    ```
9. 请确保能够完成以下流程
   ```shell
   copy /path/in/host/generator.zip /root/generator.zip # 可选
   unzip /root/generator.zip -d /root # 可选
   ./run.sh id base64(flag)
   copy /root/attachments/{id}.zip /path/in/host/{id}.zip
   ```
