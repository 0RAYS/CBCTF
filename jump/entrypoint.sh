#!/bin/sh
# 玩家跳板容器的入口脚本。
# 登录用户与公钥由平台挂载在 RANGE_SSH_DIR 下，容器只负责把公钥落到账号上并常驻 sshd。
set -eu

ssh_dir="${RANGE_SSH_DIR:-/run/range/ssh}"
ssh_port="${RANGE_SSH_PORT:-22}"
ssh_user="$(cat "$ssh_dir/user")"

if ! id "$ssh_user" >/dev/null 2>&1; then
    adduser -D -s /bin/bash -h "/home/$ssh_user" "$ssh_user"
fi
home_dir="$(getent passwd "$ssh_user" | cut -d: -f6)"

mkdir -p "$home_dir/.ssh"
cp "$ssh_dir/authorized_keys" "$home_dir/.ssh/authorized_keys"
chmod 700 "$home_dir/.ssh"
chmod 600 "$home_dir/.ssh/authorized_keys"
chown -R "$ssh_user" "$home_dir/.ssh"

printf '%s ALL=(ALL) NOPASSWD:ALL\n' "$ssh_user" > "/etc/sudoers.d/$ssh_user"
chmod 440 "/etc/sudoers.d/$ssh_user"

mkdir -p /run/sshd
ssh-keygen -A >/dev/null 2>&1 || true

cat > /etc/ssh/sshd_config <<EOF
Port $ssh_port
ListenAddress 0.0.0.0
PermitRootLogin no
PasswordAuthentication no
KbdInteractiveAuthentication no
PubkeyAuthentication yes
AuthorizedKeysFile .ssh/authorized_keys
UsePAM no
X11Forwarding no
PrintMotd no
AllowTcpForwarding yes
AllowAgentForwarding yes
Subsystem sftp /usr/lib/ssh/sftp-server
EOF

exec /usr/sbin/sshd -D -e
