FROM  docker.io/centos:centos8

# Prepare running environment
WORKDIR /
COPY ./bin/manager .
COPY ./hack/install_ansible.sh .
COPY ./bin/network-runner /usr/bin
RUN ./install_ansible.sh

ENTRYPOINT ["/manager"]
