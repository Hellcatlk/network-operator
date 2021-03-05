# API and Resource Definitions

## Port

**Port** CR represents a specific port of a network device, including port information,
the performance of the network device to which it belongs, and the performance of
the connected network interface card.

### Port Spec

The *Port Spec* defines the port on which network device and what configuration should be configured.

#### portID

The `portID` field is the port's ID on network device, it may look like `interface 0/14`.

#### portConfigurationRef

A reference to define which configuration should be configured.

#### deviceRef

A reference to define this port on which network device.

### Port Status

#### state

The `state` shows the progress of configuring the network.

* *\<empty string\>* -- means we haven't do any thing.
* *Created* -- Means the port can be configured.
* *Configuring* -- Means we are configuring configuration for the port.
* *Configured* -- Means the port have been configured, you can use it now.
* *deleting* -- Means we are removing configuration from the port.
* *deleted* -- Means now configuration of the port have been removed.

#### configurationRef

<!-- TODO -->

### Port Example

```yaml
metaData:
  name: port0
  ownerRef:
    name: bm0
    kind: BareMetalHost
    namespace: default
  finalizers:
spec:
  portID: 0
  portConfigurationRef:
    name: sc1
    kind: SwitchPortConfiguration
    namespace: default
  deviceRef:
    name: switch0
    kind: Switch
    namespace: default
status:
  state: Configured
  configurationRef:
    name: sc1
    kind: SwitchPortConfiguration
    namespace: default
```

## SwitchPortConfiguration

The **SwitchPortConfiguration** is a kind configure for port on switch.

### SwitchPortConfiguration Spec

The *SwitchPortConfiguration Spec* defines details of configuration.

#### acl

The `acl` defines access control list of switch's port.

The sub-fields are

<!-- TODO -->
* *type* --
* *action* --
* *protocol* --
* *src* --
* *srcPortRange* --
* *des* --
* *desPortRange* --

#### type

Indicates which mode this port should be set to. Valid values are *access*, *trunk* or *hybrid*,
 default value is *access*.

#### vlans

The `vlans` define port that use the configure should be configured to which vlans.

### SwitchPortConfiguration Example

```yaml
apiVersion: v1alpha1
metaData:
  name: sc1
  namespace: default
  finalizers:
    - m3m1
spec:
  acl:
    - type: // ipv4, ipv6
      action: // allow, deny
      protocol: // TCP, UDP, ICMP, ALL
      src: // xxx.xxx.xxx.xxx/xx
      srcPortRange: // 22, 22-30
      des: // xxx.xxx.xxx.xxx/xx
      desPortRange: // 22, 22-30
  type: accesss
  vlans:
    - id: 2
    - id: 3
```

## Switch

The **Switch** describes a switch.

### Switch spec

The `spec` contains the connection information for the switch.

#### os

The `os` is operator system of switch.

#### ip

The `ip` is ipv4 address of the switch.

#### port

The `port` is which port we can `ssh` to the switch.

#### secret

The `secret` is a secret resource contains username and password for the switch.

#### restrictedPorts

The `restrictedPorts` an restrict the use of ports on the switch.

The sub-fields are

* *portID* -- The port's ID on the switch.
* *disabled* -- The port can't be used or not.
* *vlanRange* -- The port can be divided into which vlans.
* *trunkDisable* -- The port can't be set to trunk port or not.

### Switch Example
```yaml
apiVersion: v1alpha1
kind: Switch
metadata:
  name: switch0
  namespace: default
spec:
  os: fos
  ip: 192.168.0.1
  port: 22
  secret:
    name: switch0-secret
    namespace: default
  restrictedPorts:
    - portID: 0
      disabled: false
      vlanRange: 1, 6-10
      trunkDisable: false
---
apiVersion: v1
kind: Secret
metadata:
  name: switch0-secret
  namespace: default
type: Opaque
data:
  username: YWRtaW4=
  password: cGFzc3dvcmQ=
```
