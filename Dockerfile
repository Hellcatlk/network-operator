# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM  docker.io/centos:centos8
WORKDIR /
COPY ./bin/manager .
COPY ./hack/install_ansible.sh .
COPY ./bin/network-runner /usr/bin

# Prepare running environment
RUN ./install_ansible.sh

ENTRYPOINT ["/manager"]
