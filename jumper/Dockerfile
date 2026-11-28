FROM alpine:3.22

RUN apk add --no-cache openssh-server openssh-keygen sudo bash \
    && mkdir -p /run/sshd \
    && ssh-keygen -A

COPY entrypoint.sh /usr/local/bin/range-jumper-entrypoint.sh
RUN chmod +x /usr/local/bin/range-jumper-entrypoint.sh

EXPOSE 22
ENTRYPOINT ["/usr/local/bin/range-jumper-entrypoint.sh"]
