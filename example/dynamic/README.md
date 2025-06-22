# 🔧 Generator 镜像模板说明

本模板用于构建 **动态题目的附件生成器 Docker 镜像**, 以支持自动化部署和分发个性化附件内容, 以下示例为密码学题目附件的生成器构建过程

---

## 🧩 工作流程

生成器逻辑的 Golang 实现位于 `/internal/k8s/generator.go`

1. 初始化:
   - 获取当前用户所属队伍的 `ID`（作为附件唯一标识）
   - 生成对应的 `flag`

2. 启动容器:
   - 使用 `sleep infinity` 命令保持容器持续运行
   - 请确保容器中存在 `sleep` 命令

3. 可选解压步骤（若存在 `generator.zip`）:
   - 将主机中的 `generator.zip` 拷贝至容器 `/root` 目录
   - 解压至 `/root`

4. 生成附件:
   - 由服务器向 `http://pod:8000/gen` 发起 `GET` 请求, 携带参数 `id` 和 `flags`, pwd参数由服务器自动生成, 无需关注
     ```bash
     /gen?id=123&flags=`base64(base64(flag1),base64(flag2),...)`&pwd=pwd
     ```
   - 返回压缩附件字节流

---

## 📦 示例流程

假设用户所属队伍 ID 为 `1`, 生成的 flag 为 `flag{test}`: 

```bash
# 启动容器, 保持运行
sleep infinity 

# generator.zip 存在时, 复制并解压
cp /host/path/generator.zip /root/generator.zip
unzip /root/generator.zip -d /root

# 下载
GET /gen?id=1&flags=ZmxhZ3R7dGVzdH0= HTTP/1.1
Host: pod:8000

```

---

## 🛠️ Dockerfile 示例

```dockerfile
FROM python:3.10-slim

COPY ./files/* /root/

RUN apt-get update &&  \
    apt-get install -y unzip &&  \
    pip install -r /root/requirements.txt  \

EXPOSE 8000

WORKDIR /root

CMD ["python", "app.py"]

```

---

## 🧳 generator.zip 文件说明

该文件为题目生成器的可选依赖, 用于 **同一基础镜像下复用不同的生成脚本**, 典型用于密码学题目中快速替换 `generator.py` 和 `template.py` 等内容

- 上传此文件由管理员完成
- 当不存在该文件时, 镜像需能 **独立生成附件**

---

## ⚠️ 注意事项

1. 容器中必须存在的命令:
   - `sleep`: 用于保持运行
   - `unzip`: 用于解压文件
   - `tee`: 用于复制文件至主机

2. 必须存在的文件:
   - `/root/run.sh`

3. `app.py` 规范:
   - 监听端口: 8000
   - 处理请求 `GET /gen` 携带三个参数
      - id: 队伍 ID, 用作唯一标识
      - flags: 格式为 `base64(base64(flag1),base64(flag2),...)`
      - pwd: 自动生成的密码参数（无需关注）, 由服务器自动生成, 与服务器注入环境变量的 `generator_pwd` 进行比对

4. 性能要求:
   - 尽可能提升附件生成速度
   - 推荐使用字节级操作（如模板替换）替代编译流程

5. 错误处理要求:
   - 附件生成流程不得产生任何错误或中断

6. flag 支持:
   - flag 可能被变异为随机字符组合（参考 `/internal/utils/flag.go` 中 `RandFlag` 函数）
   - 请确保附件能正确处理 **不定长 flag**

7. 版本号规范:
   - 禁止使用 `latest` 标签（由于 k8s 的 `IfNotPresent` 策略）
   - 正确示例: 
     ```bash
     docker build -t generator-test:20250221 .
     docker tag generator-test:20250221 docker.0rays.club/test/generator:20250221
     docker push docker.0rays.club/test/generator:20250221
     ```
