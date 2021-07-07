#!/bin/sh

set -ue

# Apply switch
echo "
---
apiVersion: v1
kind: Secret
metadata:
  name: switch-example-secret
type: Opaque
data:
  username: cnVubmVy
  password: cnVubmVy
---
apiVersion: metal3.io/v1alpha1
kind: Ansible
metadata:
  name: ansible-example
spec:
  host: <host-ip>
  bridge: \"br-test\"
  secret:
    name: switch-example-secret
    namespace: default
---
apiVersion: metal3.io/v1alpha1
kind: Switch
metadata:
  name: switch-example
spec:
  providerSwitch:
    kind: Ansible
    name: ansible-example
  ports:
    \"switchport-example\":
      name: \"test\"
" | kubectl apply -n network-operator-system -f -


# Apply switchportconifguration
echo "
apiVersion: metal3.io/v1alpha1
kind: SwitchPortConfiguration
metadata:
  name: switchportconfiguration-example
spec:
  untaggedVLAN: 10
" | kubectl apply -n network-operator-system -f -

kubectl get switch -n network-operator-system
kubectl get switchport -n network-operator-system
