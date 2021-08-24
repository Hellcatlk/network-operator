FROM  alpine:latest

# Copy file
WORKDIR /
COPY ./bin/manager .
COPY ./bin/network-runner /usr/bin

# Prepare running environment
RUN apk add ansible openssh sshpass py3-pip gcc g++ --no-cache && \
    apk add python3-dev libc-dev linux-headers --no-cache && \
    # Install network-runner
    pip3 install networking-ansible && \
    ln -s /usr/lib/python3.9/site-packages/etc/ansible /etc/ && \
    echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config && \
    # Install fos ansible plugin
    ansible-galaxy collection install fujitsu.fos && \
    ln -s ~/.ansible/collections/ansible_collections/fujitsu/fos/plugins ~/.ansible/

ENTRYPOINT ["/manager"]
