# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM  docker.io/centos:centos8
WORKDIR /
COPY ./bin/manager .
COPY ./bin/network-runner /usr/bin

# Prepare network runner environment
RUN dnf install epel-release gcc python38-devel openssl-devel openssh-clients.x86_64 -y
RUN dnf install sshpass -y
RUN python3.8 -m pip install ansible networking-ansible
RUN cp -rf /usr/local/lib/python3.8/site-packages/etc/ansible /etc/

# Clean image
RUN dnf autoremove
RUN dnf clean all

ENTRYPOINT ["/manager"]
