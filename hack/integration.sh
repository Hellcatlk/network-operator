#!/bin/sh

set -ue

# Create bridge and port
sudo ovs-vsctl add-br br-test
sudo ovs-vsctl add-port br-test test

# Apply switch
echo "
---
apiVersion: v1
kind: Secret
metadata:
  name: switch-example-secret
type: Opaque
data:
  username: $(echo -n $USER | base64 -)
  password: cnVubmVy

---
apiVersion: metal3.io/v1alpha1
kind: OVSSwitch
metadata:
  name: ovsswitch-sample
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
    kind: OVSSwitch
    name: ovsswitch-sample
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

# Apply switchport
switchUID=$(kubectl get switch switch-example -n network-operator-system -o yaml | grep "uid" | awk '{print $2}')
echo "
apiVersion: metal3.io/v1alpha1
kind: SwitchPort
metadata:
  name: switchport-example
  ownerReferences:
    - apiVersion: metal3.io/v1alpha1
      kind: Switch
      name: switch-example
      uid: $switchUID
spec:
  configuration:
    name: switchportconfiguration-example
" | kubectl apply -n network-operator-system -f -

kubectl get switchport -n network-operator-system
