# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM  docker.io/centos:centos8
WORKDIR /
COPY ./bin/manager .
COPY ./bin/network-runner /usr/bin

# Prepare network runner environment
RUN dnf install epel-release sshpass python3-pip gcc python3-devel rust cargo openssl-devel openssh-clients.x86_64 -y
RUN pip3 install wheel setuptools-rust && python3 -c 'from setuptools_rust import RustExtension'
RUN pip3 install ansible networking-ansible
RUN cp -r /usr/local/lib/python3.6/site-packages/etc/ansible /etc/

ENTRYPOINT ["/manager"]
