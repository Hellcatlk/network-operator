#!/bin/bash

set -ue

# install openvswitch
sudo apt install openvswitch-switch
# create bridge an port
sudo ovs-vsctl add-br br-test
sudo ovs-vsctl add-port br-test port-test

# Add test user
sudo useradd test -m -s /bin/bash
sudo echo test:test | sudo chpasswd
sudo usermod -aG sudo test
sudo chmod +w /etc/sudoers
sudo sed -i '/root\tALL=(ALL:ALL)/a\test\tALL=(ALL:ALL) ALL' /etc/sudoers
sudo sed -i 's/%sudo\tALL=(ALL:ALL) ALL/%sudo\tALL=(ALL:ALL) NOPASSWD:ALL/g' /etc/sudoers
sudo chmod -w /etc/sudoers

ip=`ip addr | grep 'state UP' -m 1 -A2 | tail -n1 | awk '{print $2}' | cut -f1 -d'/'`

kind load docker-image --name chart-testing --nodes chart-testing-control-plane network-operator:dev

# Apply CR
echo "
---
apiVersion: v1
kind: Secret
metadata:
  name: switch-example-secret
type: Opaque
data:
  username: dGVzdA==
  password: dGVzdA==

---
apiVersion: metal3.io/v1alpha1
kind: SwitchResource
metadata:
  name: switchresource-example
spec:
  vlanRange: 1-1000
  tenantLimits:
    \"user-1\":
      namespace: default
      vlanRange: 1-100
---
apiVersion: metal3.io/v1alpha1
kind: AnsibleSwitch
metadata:
  name: ansible-example
spec:
  os: openvswitch
  host: $ip
  bridge: \"br-test\"
  credentials:
    name: switch-example-secret
    namespace: default
---
apiVersion: metal3.io/v1alpha1
kind: Switch
metadata:
  name: switch-example
spec:
  provider:
    kind: AnsibleSwitch
    name: ansible-example
  ports:
    \"switchport-example\":
      physicalPortName: \"port-test\"
---
apiVersion: metal3.io/v1alpha1
kind: SwitchPortConfiguration
metadata:
  name: switchportconfiguration-example
spec:
  untaggedVLAN: 11
" | kubectl apply -f -

echo "Wait for controller up..."
while kubectl get deployment network-operator-controller-manager -n network-operator-system | grep -w "0/1" >/dev/null; do
  sleep 10s
done

echo "Verify port configuration feature"
# Modify switchPort to trigger port configuration
kubectl get switchport switchport-example -o yaml > switchport.yaml && sed -i "s/^spec: {}/spec:\n  configuration:\n    name: switchportconfiguration-example/g" switchport.yaml && kubectl replace -f switchport.yaml
while [ -z "$(kubectl get switchport switchport-example -o yaml | grep "state: Active")" ]; do
  sleep 5s
done
# Verify that the port configuration is successful
if [ -z "$(sudo ovs-vsctl show | grep "tag: 11")" ]; then
  echo "Switch port configuration error"
  exit 1
fi

echo "Verify port deConfiguration feature"
# Modify switchPort to trigger port configuration
kubectl get switchport switchport-example -o yaml > switchport.yaml && sed -i '/name: switchportconfiguration-example/,+1d' switchport.yaml && kubectl replace -f switchport.yaml
sleep 10s
while [ -z "$(kubectl get switchport switchport-example -o yaml | grep "state: Idle")" ]; do
  sleep 10s
done
# Verify that the port configuration is successful
if [ -n "$(sudo ovs-vsctl show | grep "tag: 11")" ]; then
  echo "Switch port deConfiguration error"
  exit 1
fi
