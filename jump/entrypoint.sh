#!/bin/sh
# 玩家跳板容器的入口脚本。
# 登录账号与公钥由平台用环境变量下发,容器只负责把公钥落到账号上并常驻 sshd。
set -eu

ssh_user="${RANGE_SSH_USER:?RANGE_SSH_USER is required}"
ssh_key="${RANGE_SSH_KEY:?RANGE_SSH_KEY is required}"

if ! id "$ssh_user" >/dev/null 2>&1; then
    adduser -D -s /bin/bash -h /home "$ssh_user"
fi
# adduser 建出来的账号是锁定的,锁定状态下 sshd 会拒绝公钥登录,
# 把锁定标记换成不可用口令,账号即可用公钥登录。
sed -i -E "s/^(${ssh_user}):[^:]*!/\1:*/" /etc/shadow

mkdir -p /home/.ssh
printf '%s\n' "$ssh_key" > /home/.ssh/authorized_keys
chmod 700 /home/.ssh
chmod 600 /home/.ssh/authorized_keys
chown -R "$ssh_user" /home/.ssh

printf '%s ALL=(ALL) NOPASSWD:ALL\n' "$ssh_user" > "/etc/sudoers.d/$ssh_user"
chmod 440 "/etc/sudoers.d/$ssh_user"

mkdir -p /run/sshd
ssh-keygen -A >/dev/null 2>&1 || true

cat > /etc/ssh/sshd_config <<EOF
Port 22
ListenAddress 0.0.0.0
PermitRootLogin no
PasswordAuthentication no
KbdInteractiveAuthentication no
PubkeyAuthentication yes
AuthorizedKeysFile .ssh/authorized_keys
X11Forwarding no
PrintMotd no
AllowTcpForwarding yes
AllowAgentForwarding yes
Subsystem sftp /usr/lib/ssh/sftp-server
EOF

cd /home
exec /usr/sbin/sshd -D -e
