FROM  alpine:latest

# Copy file
WORKDIR /
COPY ./bin/manager .
COPY ./bin/network-runner /usr/bin

# Prepare running environment
RUN apk add ansible openssh sshpass py3-pip gcc g++ --no-cache git && \
    apk add python3-dev libc-dev linux-headers --no-cache && \
    pip3 install networking-ansible && \
    echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config && \
    ansible-galaxy collection install fujitsu.fos && \
    ln -s ~/.ansible/collections/ansible_collections/fujitsu/fos/plugins ~/.ansible/ && \
    git clone https://github.com/ansible-network/network-runner.git && \
    rm -rf /usr/lib/python3.9/site-packages/network_runner && \
    mv -f ./network-runner/network_runner /usr/lib/python3.9/site-packages/ && \
    cp -rf ./network-runner/etc/ansible /etc/ && \
    rm -rf ./network-runner


ENTRYPOINT ["/manager"]
