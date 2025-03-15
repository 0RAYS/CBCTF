FROM alpine:latest

ENV FRP_VERSION=0.61.2
ENV FRP_TAR=frp_${FRP_VERSION}_linux_amd64.tar.gz

RUN apk add --no-cache wget tar gettext && \
    wget https://github.com/fatedier/frp/releases/download/v${FRP_VERSION}/${FRP_TAR} && \
    tar -xzf ${FRP_TAR} && \
    mv frp_${FRP_VERSION}_linux_amd64 /frpc && \
    rm -f ${FRP_TAR}

WORKDIR /frpc

# serverAddr = "${serverAddr}"
# serverPort = ${serverPort}
# auth.token = "${token}"

# [[proxies]]
# name = "${name}"
# type = "${type}"
# localIP = "${localIP}"
# localPort = ${localPort}
# remotePort = ${remotePort}
RUN echo "c2VydmVyQWRkciA9ICIke3NlcnZlckFkZHJ9Ig0Kc2VydmVyUG9ydCA9ICR7c2VydmVyUG9ydH0NCmF1dGgudG9rZW4gPSAiJHt0b2tlbn0iDQoNCltbcHJveGllc11dDQpuYW1lID0gIiR7bmFtZX0iDQp0eXBlID0gIiR7dHlwZX0iDQpsb2NhbElQID0gIiR7bG9jYWxJUH0iDQpsb2NhbFBvcnQgPSAke2xvY2FsUG9ydH0NCnJlbW90ZVBvcnQgPSAke3JlbW90ZVBvcnR9DQo=" | base64 -d > /frpc/frpc.template.toml

# 运行 frpc 客户端
CMD ["/bin/sh", "-c", "envsubst < /frpc/frpc.template.toml > /frpc/frpc.toml && /frpc/frpc -c /frpc/frpc.toml"]
