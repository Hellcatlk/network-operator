FROM  docker.io/centos:centos8

# Copy file
WORKDIR /
COPY ./bin/manager .
COPY ./bin/network-runner /usr/bin

# Prepare running environment
RUN dnf install gcc rust cargo epel-release python3-devel openssl-devel -y && \
    dnf install openssh-clients.x86_64 sshpass -y && \
    pip3 install wheel setuptools-rust && \
    python3 -c 'from setuptools_rust import RustExtension' && \
    pip3 install ansible networking-ansible && \
    cp -rf /usr/local/lib/python3.6/site-packages/etc/ansible /etc/ && \
    dnf autoremove && \
    dnf clean all

ENTRYPOINT ["/manager"]
