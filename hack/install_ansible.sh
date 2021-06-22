#!/usr/bin/sh

set -ue

# Install ansible and network-runner
dnf install epel-release gcc python38-devel openssl-devel openssh-clients.x86_64 -y
dnf install sshpass -y
python3.8 -m pip install ansible networking-ansible
cp -rf /usr/local/lib/python3.8/site-packages/etc/ansible /etc/

# Clean
dnf autoremove
dnf clean all
